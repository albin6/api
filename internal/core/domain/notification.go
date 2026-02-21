package domain

import "time"

type NotificationType string

const (
	NotificationFollowUpAssigned NotificationType = "FOLLOWUP_ASSIGNED"
)

type Notification struct {
	ID        string           `json:"id"`
	Type      NotificationType `json:"type"`
	Title     string           `json:"title"`
	Message   string           `json:"message"`
	Data      interface{}      `json:"data"`
	CreatedAt time.Time        `json:"created_at"`
}

type FollowUpAssignedData struct {
	FollowUpID  uint   `json:"followup_id"`
	StudentName string `json:"student_name"`
	AssignedBy  string `json:"assigned_by"`
}
