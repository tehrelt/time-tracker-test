package repositories

import (
	"em-test/internal/domain"
	"em-test/internal/lib/dto"
	"em-test/internal/lib/filters"
	"log/slog"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

// var _ services.UserRepository = (*UsersRepository)(nil)

type UsersRepository struct {
	db *sqlx.DB
}

func NewUsersRepository(db *sqlx.DB) *UsersRepository {
	return &UsersRepository{
		db: db,
	}
}

func (u *UsersRepository) Add(dto dto.SaveUserDto) (*domain.User, error) {
	fn := "UsersRepository.Add"
	logger := slog.With(slog.String("fn", fn))

	id := uuid.New()

	query, args, err := sq.Insert(USERS_TABLE).
		Columns("id", "surname", "name", "patronymic", "address", "passport_serie", "passport_number").
		Values(id.String(), dto.Surname, dto.Name, dto.Patronymic, dto.Address, dto.PassportSerie, dto.PassportNumber).
		Suffix("RETURNING *").
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		slog.Error("error formatting query", slog.String("err", err.Error()))
		return nil, err
	}

	logger.Debug("executing query", slog.String("query", query), slog.Any("args", args))

	var user domain.User
	if err = u.db.Get(&user, query, args...); err != nil {
		slog.Error("error executing query", slog.String("err", err.Error()))
		if e, ok := err.(*pq.Error); ok {
			if e.Code == "23505" {
				return nil, domain.ErrUserAlreadyExists
			}
		}
		return nil, err
	}

	return &user, nil
}

func (u *UsersRepository) Read(id string) (*domain.User, error) {
	fn := "UsersRepository.Read"
	logger := slog.With(slog.String("fn", fn))

	query, args, err := sq.Select("*").
		From(USERS_TABLE).
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		slog.Error("error formatting query", slog.String("err", err.Error()))
		return nil, err
	}

	logger.Debug("executing query", slog.String("query", query), slog.Any("args", args))

	var users domain.User
	if err = u.db.Get(&users, query, args...); err != nil {
		slog.Error("error executing query", slog.String("err", err.Error()))
		return nil, err
	}

	return &users, nil
}

func (u *UsersRepository) ReadMany(filters *filters.UsersFilters) ([]*domain.User, int64, error) {
	fn := "UsersRepository.ReadMany"
	logger := slog.With(slog.String("fn", fn))

	builder := sq.Select("*").
		From(USERS_TABLE).
		OrderBy("id ASC").
		PlaceholderFormat(sq.Dollar)

	if filters != nil {
		if filters.Surname != nil {
			builder = builder.Where(sq.ILike{"surname": *filters.Surname + "%"})
		}

		if filters.Name != nil {
			builder = builder.Where(sq.ILike{"name": *filters.Name + "%"})
		}

		if filters.Patronymic != nil {
			builder = builder.Where(sq.ILike{"patronymic": *filters.Patronymic + "%"})
		}

		if filters.Address != nil {
			builder = builder.Where(sq.ILike{"address": *filters.Address + "%"})
		}

		if filters.Limit != nil {
			builder = builder.Limit(uint64(*filters.Limit))
		} else {
			builder = builder.Limit(100)
		}

		if filters.Offset != nil {
			builder = builder.Offset(uint64(*filters.Offset))
		}
	}

	sql, args, err := builder.ToSql()
	if err != nil {
		logger.Error("error formatting query", slog.String("err", err.Error()))
		return nil, 0, err
	}

	logger.Debug("executing query", slog.String("query", sql), slog.Any("args", args))

	users := make([]*domain.User, 0)
	if err = u.db.Select(&users, sql, args...); err != nil {
		logger.Error("error executing query", slog.String("err", err.Error()))
		return nil, 0, err
	}

	total, err := u.Count(filters)
	if err != nil {
		logger.Error("error with couning records in db by filters")
		return nil, 0, err
	}

	return users, total, nil
}

func (u *UsersRepository) Count(filters *filters.UsersFilters) (int64, error) {
	fn := "UsersRepository.Count"
	logger := slog.With(slog.String("fn", fn))

	builder := sq.Select("COUNT(*)").
		From(USERS_TABLE).
		PlaceholderFormat(sq.Dollar)

	if filters != nil {
		if filters.Surname != nil {
			builder = builder.Where(sq.ILike{"surname": *filters.Surname + "%"})
		}

		if filters.Name != nil {
			builder = builder.Where(sq.ILike{"name": *filters.Name + "%"})
		}

		if filters.Patronymic != nil {
			builder = builder.Where(sq.ILike{"patronymic": *filters.Patronymic + "%"})
		}

		if filters.Address != nil {
			builder = builder.Where(sq.ILike{"address": *filters.Address + "%"})
		}
	}

	sql, args, err := builder.ToSql()
	if err != nil {
		logger.Error("error formatting query", slog.String("err", err.Error()))
		return 0, err
	}

	logger.Debug("executing query", slog.String("query", sql), slog.Any("args", args))

	var total int64
	if err := u.db.Get(&total, sql, args...); err != nil {
		logger.Error("error executing query", slog.String("err", err.Error()))
		return 0, err
	}

	return total, nil
}
