package redis

import (
	"context"
	"encoding/json"
	"log"

	"github.com/albin6/api/internal/core/domain"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type UserSubscriber struct {
	rdb *redis.Client
	db  *gorm.DB
}

func NewUserSubscriber(rdb *redis.Client, db *gorm.DB) *UserSubscriber {
	return &UserSubscriber{rdb: rdb, db: db}
}

func (s *UserSubscriber) SubscribeToUserUpdates() {
	ctx := context.Background()
	// Subscribe to both created and updated events
	pubsub := s.rdb.Subscribe(ctx, "user:created", "user:updated")
	defer pubsub.Close()

	ch := pubsub.Channel()

	log.Println("UserSubscriber: Listening for user updates on Redis...")

	for msg := range ch {
		var event struct {
			ID        string  `json:"id"`
			Name      string  `json:"name"`
			Email     string  `json:"email"`
			Role      string  `json:"role"`
			AvatarURL *string `json:"avatar_url"`
		}

		if err := json.Unmarshal([]byte(msg.Payload), &event); err != nil {
			log.Printf("UserSubscriber: Failed to unmarshal user event: %v", err)
			continue
		}

		log.Printf("UserSubscriber: Received user event for %s (%s)", event.Email, event.ID)

		// Upsert logic
		var existing domain.User
		if err := s.db.First(&existing, "id = ?", event.ID).Error; err == nil {
			// Update
			existing.Name = event.Name
			existing.Email = event.Email
			existing.Role = domain.Role(event.Role)
			// Avatar is not in domain.User yet?
			if err := s.db.Save(&existing).Error; err != nil {
				log.Printf("UserSubscriber: Failed to update user: %v", err)
			} else {
				log.Println("UserSubscriber: User updated successfully")
			}
		} else {
			// Create
			newUser := domain.User{
				ID:       event.ID,
				Name:     event.Name,
				Email:    event.Email,
				Role:     domain.Role(event.Role),
				Password: "", // Password not synced, or use placeholder?
			}
			if err := s.db.Create(&newUser).Error; err != nil {
				log.Printf("UserSubscriber: Failed to create user: %v", err)
			} else {
				log.Println("UserSubscriber: User created successfully")
			}
		}
	}
}
