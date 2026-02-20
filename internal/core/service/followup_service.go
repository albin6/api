package service

import (
	"context"
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/albin6/api/internal/core/domain"
	"github.com/albin6/api/internal/core/port"
)

type FollowUpService struct {
	followUpRepo port.FollowUpRepository
	contactRepo  port.ContactLogRepository
	meetingRepo  port.MeetingRepository
	outcomeRepo  port.MeetingOutcomeRepository
	reminderRepo port.ReminderRepository
	hub          WebSocketHub
}

type WebSocketHub interface {
	BroadcastToUser(userID string, notification *domain.Notification)
}

func NewFollowUpService(
	followUpRepo port.FollowUpRepository,
	contactRepo port.ContactLogRepository,
	meetingRepo port.MeetingRepository,
	outcomeRepo port.MeetingOutcomeRepository,
	reminderRepo port.ReminderRepository,
	hub WebSocketHub,
) *FollowUpService {
	return &FollowUpService{
		followUpRepo: followUpRepo,
		contactRepo:  contactRepo,
		meetingRepo:  meetingRepo,
		outcomeRepo:  outcomeRepo,
		reminderRepo: reminderRepo,
		hub:          hub,
	}
}

func (s *FollowUpService) CreateFollowUp(ctx context.Context, studentID uint, assignedTo, createdBy string) (*domain.StudentFollowUp, error) {
	// User role validation transferred to Gateway/Auth service

	followUp := &domain.StudentFollowUp{
		StudentID:  studentID,
		AssignedTo: assignedTo,
		Stage:      domain.StageContactPending,
	}

	if err := s.followUpRepo.Create(ctx, followUp); err != nil {
		return nil, err
	}

	result, err := s.followUpRepo.GetByID(ctx, followUp.ID)
	if err != nil {
		return nil, err
	}

	// Notification disabled temporarily
	// go s.notifyFollowUpAssigned(ctx, result, creator)

	return result, nil
}

func (s *FollowUpService) GetFollowUp(ctx context.Context, id uint, requestingUserID string) (*domain.StudentFollowUp, error) {
	followUp, err := s.followUpRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if followUp.AssignedTo != requestingUserID {
		return nil, errors.New("unauthorized: you can only view follow-ups assigned to you")
	}

	return followUp, nil
}

func (s *FollowUpService) ListFollowUps(ctx context.Context, stage *domain.FollowUpStage, assignedTo *string, page, limit int) (map[string]interface{}, error) {
	const maxLimit = 100
	if limit > maxLimit {
		limit = maxLimit
	}
	if limit < 1 {
		limit = 10
	}
	if page < 1 {
		page = 1
	}

	followUps, total, err := s.followUpRepo.GetAll(ctx, stage, assignedTo, page, limit)
	if err != nil {
		return nil, err
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	return map[string]interface{}{
		"data": followUps,
		"pagination": map[string]interface{}{
			"total":       total,
			"page":        page,
			"limit":       limit,
			"total_pages": totalPages,
		},
	}, nil
}

func (s *FollowUpService) AddContactLog(ctx context.Context, followUpID uint, userID string, successful bool, notes string) error {
	followUp, err := s.followUpRepo.GetByID(ctx, followUpID)
	if err != nil {
		return err
	}

	if followUp.AssignedTo != userID {
		return errors.New("unauthorized: only the assigned member can add contact logs")
	}

	if followUp.Stage != domain.StageContactPending {
		return fmt.Errorf("invalid state: cannot add contact log in stage %s", followUp.Stage)
	}

	contactLog := &domain.ContactLog{
		FollowUpID:  followUpID,
		Successful:  successful,
		Notes:       notes,
		ContactedAt: time.Now(),
	}

	if err := s.contactRepo.Create(ctx, contactLog); err != nil {
		return err
	}

	if successful {
		if err := s.followUpRepo.UpdateStage(ctx, followUpID, domain.StageContactCompleted); err != nil {
			return err
		}
	}

	return nil
}

func (s *FollowUpService) ScheduleMeeting(ctx context.Context, followUpID uint, userID string, scheduledAtStr string, meetingLink string, participantIDs []string) (*domain.Meeting, error) {
	followUp, err := s.followUpRepo.GetByID(ctx, followUpID)
	if err != nil {
		return nil, err
	}

	if followUp.AssignedTo != userID {
		return nil, errors.New("unauthorized: only the assigned member can schedule meetings")
	}

	if followUp.Stage != domain.StageContactCompleted {
		return nil, errors.New("invalid state: cannot schedule meeting before contact completion")
	}

	scheduledAt, err := time.Parse(time.RFC3339, scheduledAtStr)
	if err != nil {
		return nil, errors.New("invalid scheduled_at format, use RFC3339 (e.g., 2026-02-15T10:00:00Z)")
	}

	if scheduledAt.Before(time.Now()) {
		return nil, errors.New("scheduled_at must be in the future")
	}

	if meetingLink == "" {
		return nil, errors.New("meeting_link is required")
	}

	if len(participantIDs) < 1 {
		return nil, errors.New("at least one participant is required")
	}

	meeting := &domain.Meeting{
		FollowUpID:  followUpID,
		ScheduledAt: scheduledAt,
		MeetingLink: meetingLink,
		Status:      domain.MeetingScheduled,
		CreatedBy:   userID,
	}

	if err := s.meetingRepo.Create(ctx, meeting); err != nil {
		return nil, err
	}

	if err := s.meetingRepo.AddParticipants(ctx, meeting.ID, participantIDs); err != nil {
		return nil, err
	}

	if err := s.followUpRepo.UpdateStage(ctx, followUpID, domain.StageMeetingScheduled); err != nil {
		return nil, err
	}

	return s.meetingRepo.GetByID(ctx, meeting.ID)
}

func (s *FollowUpService) ListMeetings(ctx context.Context, followUpID uint, userID string) ([]domain.Meeting, error) {
	followUp, err := s.followUpRepo.GetByID(ctx, followUpID)
	if err != nil {
		return nil, err
	}

	if followUp.AssignedTo != userID {
		return nil, errors.New("unauthorized: only the assigned member can view meetings")
	}

	return s.meetingRepo.GetByFollowUpID(ctx, followUpID)
}

func (s *FollowUpService) CompleteMeeting(ctx context.Context, meetingID uint, userID string) error {
	meeting, err := s.meetingRepo.GetByID(ctx, meetingID)
	if err != nil {
		return err
	}

	followUp, err := s.followUpRepo.GetByID(ctx, meeting.FollowUpID)
	if err != nil {
		return err
	}

	if followUp.AssignedTo != userID {
		return errors.New("unauthorized: only the assigned member can mark meetings as complete")
	}

	if followUp.Stage != domain.StageMeetingScheduled {
		return fmt.Errorf("invalid state: cannot complete meeting in stage %s", followUp.Stage)
	}

	if meeting.Status != domain.MeetingScheduled {
		return errors.New("meeting is not in scheduled state")
	}

	if err := s.meetingRepo.UpdateStatus(ctx, meetingID, domain.MeetingCompleted); err != nil {
		return err
	}

	return s.followUpRepo.UpdateStage(ctx, meeting.FollowUpID, domain.StageMeetingCompleted)
}

func (s *FollowUpService) SubmitOutcome(ctx context.Context, meetingID uint, userID string, status domain.OutcomeStatus, remarks, recordingURL string, nextFollowUpAtStr *string) error {
	meeting, err := s.meetingRepo.GetByID(ctx, meetingID)
	if err != nil {
		return err
	}

	followUp, err := s.followUpRepo.GetByID(ctx, meeting.FollowUpID)
	if err != nil {
		return err
	}

	if followUp.AssignedTo != userID {
		return errors.New("unauthorized: only the assigned member can submit outcomes")
	}

	if followUp.Stage != domain.StageMeetingCompleted {
		return errors.New("invalid state: cannot submit outcome before meeting completion")
	}

	if status != domain.OutcomeSelected && status != domain.OutcomeRejected {
		return errors.New("status must be SELECTED or REJECTED")
	}
	if remarks == "" {
		return errors.New("remarks are required")
	}
	if recordingURL == "" {
		return errors.New("recording_url is required")
	}

	var nextFollowUpAt *time.Time
	if status == domain.OutcomeRejected {
		if nextFollowUpAtStr == nil || *nextFollowUpAtStr == "" {
			return errors.New("next_follow_up_at is required for rejected outcomes")
		}
		parsed, err := time.Parse(time.RFC3339, *nextFollowUpAtStr)
		if err != nil {
			return errors.New("invalid next_follow_up_at format, use RFC3339")
		}
		if parsed.Before(time.Now()) {
			return errors.New("next_follow_up_at must be in the future")
		}
		nextFollowUpAt = &parsed
	}

	outcome := &domain.MeetingOutcome{
		MeetingID:      meetingID,
		Status:         status,
		Remarks:        remarks,
		RecordingURL:   recordingURL,
		NextFollowUpAt: nextFollowUpAt,
	}

	if err := s.outcomeRepo.Create(ctx, outcome); err != nil {
		return err
	}

	var newStage domain.FollowUpStage
	if status == domain.OutcomeSelected {
		newStage = domain.StageSelected
	} else {
		newStage = domain.StageRejected

		reminder := &domain.FollowUpReminder{
			FollowUpID: followUp.ID,
			RemindAt:   *nextFollowUpAt,
			Sent:       false,
		}
		if err := s.reminderRepo.Create(ctx, reminder); err != nil {
			return err
		}
	}

	return s.followUpRepo.UpdateStage(ctx, followUp.ID, newStage)
}

func (s *FollowUpService) RestartFollowUp(ctx context.Context, followUpID uint, userID string) error {
	followUp, err := s.followUpRepo.GetByID(ctx, followUpID)
	if err != nil {
		return err
	}

	if followUp.AssignedTo != userID {
		return errors.New("unauthorized: only the assigned member can restart follow-ups")
	}

	if followUp.Stage != domain.StageRejected {
		return errors.New("can only restart rejected follow-ups")
	}

	return s.followUpRepo.UpdateStage(ctx, followUpID, domain.StageContactPending)
}

/*
func (s *FollowUpService) notifyFollowUpAssigned(ctx context.Context, followUp *domain.StudentFollowUp, creator *domain.User) {
	if s.hub == nil {
		return
	}


	notification := &domain.Notification{
		ID:      uuid.New().String(),
		Type:    domain.NotificationFollowUpAssigned,
		Title:   "New Follow-Up Assigned",
		Message: fmt.Sprintf("You have been assigned to follow up with %s", followUp.Student.FullName),
		Data: domain.FollowUpAssignedData{
			FollowUpID:  followUp.ID,
			StudentName: followUp.Student.FullName,
			AssignedBy:  creator.Name,
		},
		CreatedAt: time.Now(),
	}


	fmt.Printf("[FollowUpService] Notification generated for user %d: %+v\n", followUp.AssignedTo, notification)
	s.hub.BroadcastToUser(followUp.AssignedTo, notification)
}
*/
