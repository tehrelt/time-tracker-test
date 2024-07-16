package repositories

import (
	"database/sql"
	"em-test/internal/domain"
	"em-test/internal/lib/dto"
	"em-test/internal/lib/filters"
	"errors"
	"log/slog"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/rcmonitor/pginterval"
)

// var _ services.ActivityRepository = (*ActivityRepository)(nil)

type ActivityRepository struct {
	db *sqlx.DB
}

func (a *ActivityRepository) Create(activity *dto.SaveActivity) error {

	fn := "ActivityRepository.Create"
	logger := slog.With(slog.String("fn", fn))

	sql, args, err := sq.Insert(ACTIVITY_TABLE).
		Columns("user_id", "start_time").
		Values(activity.UserId, activity.StartTime).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		logger.Error("failed to build sql", slog.String("err", err.Error()))
		return err
	}

	logger.Debug("executing query", slog.String("sql", sql), slog.Any("args", args))

	if _, err := a.db.Exec(sql, args...); err != nil {
		logger.Error("failed to execute query", slog.String("err", err.Error()))
		return err
	}

	return nil
}

func (a *ActivityRepository) IsActive(userId string) (bool, error) {
	fn := "ActivityRepository.IsActive"
	logger := slog.With(slog.String("fn", fn))

	query, args, err := sq.
		Select(`a.id, a.start_time, a.end_time, u.id as "user.id", u.surname as "user.surname", u.name as "user.name", u.patronymic as "user.patronymic", u.passport_serie as "user.passport_serie", u.passport_number as "user.passport_number"`).
		From(ACTIVITY_TABLE + " a").
		Join("users u ON u.id = a.user_id").
		Where(sq.And{
			sq.Eq{"a.user_id": userId},
			sq.Eq{"a.end_time": nil},
		}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		logger.Error("failed to build sql", slog.String("err", err.Error()))
		return false, err
	}

	logger.Debug("executing query", slog.String("sql", query), slog.Any("args", args))

	var res domain.ActivityRecord
	if err := a.db.Get(&res, query, args...); err != nil {
		logger.Error("failed to execute query", slog.String("err", err.Error()))
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}

	return true, err
}

func (a *ActivityRepository) PatchEndTime(d *dto.StopActivityDto) error {
	fn := "ActivityRepository.PatchEndTime"
	logger := slog.With(slog.String("fn", fn))

	sql, args, err := sq.Update(ACTIVITY_TABLE+" a").
		Set("end_time", d.EndTime).
		Where(sq.And{
			sq.Eq{"a.user_id": d.UserId},
			sq.Eq{"a.end_time": nil},
		}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		logger.Error("failed to build sql", slog.String("err", err.Error()))
		return err
	}

	logger.Debug("executing query", slog.String("sql", sql), slog.Any("args", args))

	if _, err := a.db.Exec(sql, args...); err != nil {
		logger.Error("failed to execute query", slog.String("err", err.Error()))
		return err
	}

	return nil
}

func (a *ActivityRepository) GetSessions(f *filters.Activity) ([]*domain.Session, error) {
	fn := "ActivityRepository.GetSessions"
	logger := slog.With(slog.String("fn", fn), slog.Any("filters", f))

	builder := sq.Select("start_time", "end_time").
		From(ACTIVITY_TABLE).
		Where(sq.Eq{"user_id": f.UserId}).
		PlaceholderFormat(sq.Dollar)

	if f.StartTime != nil {
		builder = builder.Where(sq.GtOrEq{"start_time": f.StartTime})
	}

	if f.EndTime != nil {
		builder = builder.Where(sq.LtOrEq{"end_time": f.EndTime})
	}

	query, args, err := builder.ToSql()
	if err != nil {
		logger.Error("failed to build sql", slog.String("err", err.Error()))
		return nil, err
	}

	logger.Debug("executing query", slog.String("sql", query), slog.Any("args", args))

	var res []*domain.Session
	if err := a.db.Select(&res, query, args...); err != nil {
		logger.Error("failed to execute query", slog.String("err", err.Error()))
		return nil, err
	}

	return res, nil
}

func (a *ActivityRepository) GetSummary(f *filters.Activity) (time.Duration, int, error) {
	fn := "ActivityRepository.GetSummary"
	logger := slog.With(slog.String("fn", fn), slog.Any("filters", f))

	builder := sq.Select("SUM(end_time - start_time) as duration, COUNT(*) as total").
		From(ACTIVITY_TABLE).
		Where(sq.Eq{"user_id": f.UserId}).
		PlaceholderFormat(sq.Dollar)

	if f.StartTime != nil {
		builder = builder.Where(sq.GtOrEq{"start_time": f.StartTime})
	}

	if f.EndTime != nil {
		builder = builder.Where(sq.LtOrEq{"end_time": f.EndTime})
	}

	query, args, err := builder.ToSql()
	if err != nil {
		logger.Error("failed to build sql", slog.String("err", err.Error()))
		return 0, 0, err
	}

	logger.Debug("executing query", slog.String("sql", query), slog.Any("args", args))

	var durationString string
	var total int

	if err := a.db.QueryRow(query, args...).Scan(
		&durationString,
		&total,
	); err != nil {
		logger.Error("failed to scan result", slog.String("err", err.Error()))
		return 0, 0, err
	}

	duration, err := pginterval.FParse(durationString)
	if err != nil {
		logger.Error("failed to parse duration", slog.String("err", err.Error()))
		return 0, 0, err
	}

	logger.Debug("row scanned", slog.Any("duration", duration), slog.Int("total", total), slog.String("durationString", durationString))

	return *duration, total, nil
}

func NewActivityRepository(db *sqlx.DB) *ActivityRepository {
	return &ActivityRepository{db: db}
}
