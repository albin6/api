package domain

type FollowUpStage string

const (
	StageContactPending   FollowUpStage = "CONTACT_PENDING"
	StageContactCompleted FollowUpStage = "CONTACT_COMPLETED"
	StageMeetingScheduled FollowUpStage = "MEETING_SCHEDULED"
	StageMeetingCompleted FollowUpStage = "MEETING_COMPLETED"
	StageSelected         FollowUpStage = "SELECTED"
	StageRejected         FollowUpStage = "REJECTED"
)

func (s FollowUpStage) String() string {
	return string(s)
}

// IsValid checks if the stage is a valid FollowUpStage
func (s FollowUpStage) IsValid() bool {
	switch s {
	case StageContactPending, StageContactCompleted, StageMeetingScheduled,
		StageMeetingCompleted, StageSelected, StageRejected:
		return true
	}
	return false
}
