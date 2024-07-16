package services

import (
	"em-test/internal/domain"
	"em-test/internal/lib/dto"
	"em-test/internal/lib/filters"
	"log/slog"
	"time"
)

type ActivityRepository interface {
	Create(*dto.SaveActivity) error
	IsActive(userId string) (bool, error)
	PatchEndTime(*dto.StopActivityDto) error

	GetSessions(*filters.Activity) ([]*domain.Session, error)
	GetSummary(f *filters.Activity) (duration time.Duration, total int, err error)
}

type ActivityService struct {
	activityRepository ActivityRepository
}

func NewActivityService(activityRepository ActivityRepository) *ActivityService {
	return &ActivityService{
		activityRepository: activityRepository,
	}
}

func (s *ActivityService) Start(userId string) error {

	fn := "ActivityService.Start"
	logger := slog.With(slog.String("fn", fn), slog.String("userId", userId))

	logger.Debug("checking active record")
	isActive, err := s.activityRepository.IsActive(userId)
	if err != nil {
		logger.Error("checking activity error", slog.String("err", err.Error()))
		return err
	}

	if isActive {
		logger.Debug("found not finished activity")
		return domain.ErrUserAlreadyWorking
	}

	saveDto := &dto.SaveActivity{
		UserId:    userId,
		StartTime: time.Now(),
	}
	logger.Debug("creating activity", slog.Any("dto", saveDto))
	return s.activityRepository.Create(saveDto)
}

func (s *ActivityService) Stop(userId string) error {

	fn := "ActivityService.Stop"
	logger := slog.With(slog.String("fn", fn), slog.String("userId", userId))

	isActive, err := s.activityRepository.IsActive(userId)
	if err != nil {
		logger.Error("checking activity error", slog.String("err", err.Error()))
		return err
	}

	if !isActive {
		logger.Debug("user not working at this moment")
		return domain.ErrUserNotWorking
	}

	d := &dto.StopActivityDto{
		UserId:  userId,
		EndTime: time.Now(),
	}
	logger.Debug("patching end time", slog.Any("dto", d))
	return s.activityRepository.PatchEndTime(d)
}

func (s *ActivityService) GetSummary(f *filters.Activity) (*domain.ActivitySummary, error) {
	fn := "ActivityService.GetSummary"
	logger := slog.With(slog.String("fn", fn), slog.Any("filters", f))

	isActive, err := s.activityRepository.IsActive(f.UserId)
	if err != nil {
		logger.Error("checking activity error", slog.String("err", err.Error()))
		return nil, err
	}

	sessions, err := s.activityRepository.GetSessions(f)
	if err != nil {
		logger.Error("getting sessions error", slog.String("err", err.Error()))
		return nil, err
	}

	duration, total, err := s.activityRepository.GetSummary(f)
	if err != nil {
		logger.Error("getting summary error", slog.String("err", err.Error()))
		return nil, err
	}

	summary := &domain.ActivitySummary{
		UserId:      f.UserId,
		IsActiveNow: isActive,
		Sessions:    sessions,
		TotalTime:   duration,
		TotalCount:  total,
	}

	logger.Debug("calculated summary", slog.Any("summary", summary))
	return summary, nil
}
