package redis

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/albin6/api/internal/core/domain"
	goredis "github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// userEventPayload mirrors the UserEventPayload published by the backend.
// Fields must match the JSON keys in backend/internal/events/payload.go exactly.
// If the backend adds new fields, add them here too (they will default to zero-value
// on old code, which is safe — additive changes only).
type userEventPayload struct {
	ID         string  `json:"id"`
	Name       string  `json:"name"`
	Email      string  `json:"email"`
	Role       string  `json:"role"`
	IsApproved bool    `json:"is_approved"`
	AvatarURL  *string `json:"avatar_url"`
}

// eventEnvelope mirrors the Envelope struct published by the backend on versioned channels.
type eventEnvelope struct {
	Version    string           `json:"version"`
	Source     string           `json:"source"`
	Action     string           `json:"action"`
	OccurredAt time.Time        `json:"occurred_at"`
	Payload    userEventPayload `json:"payload"`
}

// UserSubscriber listens on both the legacy short channels and the new versioned
// channels, performing an upsert into the local DB on each event.
type UserSubscriber struct {
	rdb *goredis.Client
	db  *gorm.DB
}

func NewUserSubscriber(rdb *goredis.Client, db *gorm.DB) *UserSubscriber {
	return &UserSubscriber{rdb: rdb, db: db}
}

// SubscribeToUserUpdates subscribes to the legacy and versioned user channels.
// Called as a goroutine from main.
func (s *UserSubscriber) SubscribeToUserUpdates() {
	ctx := context.Background()

	// Subscribe to both old (legacy) and new (versioned) channels so that
	// the API service keeps working during the migration window.
	pubsub := s.rdb.Subscribe(ctx,
		// Legacy channels (backend still publishes on these)
		"user:created",
		"user:updated",
		// Versioned channels
		"taskflow:v1:user:created",
		"taskflow:v1:user:updated",
		"taskflow:v1:user:approved",
		"taskflow:v1:user:role_changed",
		// user:deleted — handled by marking inactive, not implemented here yet
	)
	defer pubsub.Close()

	ch := pubsub.Channel()
	log.Println("UserSubscriber: Listening for user updates on Redis...")

	for msg := range ch {
		s.handleMessage(msg)
	}
}

func (s *UserSubscriber) handleMessage(msg *goredis.Message) {
	// Determine if this is a versioned envelope or a legacy flat payload.
	var payload userEventPayload

	if isVersionedChannel(msg.Channel) {
		var env eventEnvelope
		if err := json.Unmarshal([]byte(msg.Payload), &env); err != nil {
			log.Printf("UserSubscriber: Failed to unmarshal envelope on %s: %v", msg.Channel, err)
			return
		}
		payload = env.Payload
		log.Printf("UserSubscriber: [%s] action=%s user=%s", msg.Channel, env.Action, payload.Email)
	} else {
		// Legacy flat payload
		if err := json.Unmarshal([]byte(msg.Payload), &payload); err != nil {
			log.Printf("UserSubscriber: Failed to unmarshal legacy event on %s: %v", msg.Channel, err)
			return
		}
		log.Printf("UserSubscriber: [legacy %s] user=%s", msg.Channel, payload.Email)
	}

	if payload.ID == "" {
		log.Printf("UserSubscriber: Skipping event with empty ID on channel %s", msg.Channel)
		return
	}

	s.upsertUser(payload)
}

func (s *UserSubscriber) upsertUser(p userEventPayload) {
	var existing domain.User
	if err := s.db.First(&existing, "id = ?", p.ID).Error; err == nil {
		// Update
		existing.Name = p.Name
		existing.Email = p.Email
		existing.Role = domain.Role(p.Role)
		if err := s.db.Save(&existing).Error; err != nil {
			log.Printf("UserSubscriber: Failed to update user %s: %v", p.ID, err)
		} else {
			log.Printf("UserSubscriber: User updated — %s (%s)", p.Email, p.ID)
		}
	} else {
		// Create (snapshot — password not synced)
		newUser := domain.User{
			ID:       p.ID,
			Name:     p.Name,
			Email:    p.Email,
			Role:     domain.Role(p.Role),
			Password: "",
		}
		if err := s.db.Create(&newUser).Error; err != nil {
			log.Printf("UserSubscriber: Failed to create user %s: %v", p.ID, err)
		} else {
			log.Printf("UserSubscriber: User created — %s (%s)", p.Email, p.ID)
		}
	}
}

// isVersionedChannel returns true for taskflow:v1:* channels.
func isVersionedChannel(channel string) bool {
	return len(channel) > 11 && channel[:11] == "taskflow:v1"
}
