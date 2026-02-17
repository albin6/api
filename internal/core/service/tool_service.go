package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"sort"
	"strings"
	"time"

	"github.com/albin6/api/internal/core/port"
	"github.com/albin6/api/pkg/toolapi"
	"github.com/redis/go-redis/v9"
)

type ToolService struct {
	client *toolapi.Client
	redis  *redis.Client
}

func NewToolService(client *toolapi.Client, rdb *redis.Client) port.ToolService {
	return &ToolService{
		client: client,
		redis:  rdb,
	}
}

func (s *ToolService) GetStudents(ctx context.Context, params map[string]string) (*toolapi.StudentListDTOResponse, error) {
	// Generate cache key based on params
	cacheKey := "tool:students:dto:" + generateCacheKey(params)

	// Try getting from cache
	val, err := s.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		var resp toolapi.StudentListDTOResponse
		if err := json.Unmarshal([]byte(val), &resp); err == nil {
			return &resp, nil
		}
	}

	// Fetch from API
	fullResp, err := s.client.GetStudents(params)
	if err != nil {
		return nil, err
	}

	// Map to DTO
	dtos := make([]toolapi.StudentDTO, len(fullResp.Data))
	for i, student := range fullResp.Data {
		dtos[i] = toolapi.StudentDTO{
			ID:           student.ID,
			AdmissionID:  student.AdmissionID,
			Name:         student.Name,
			Email:        student.Email,
			Mobile:       student.Mobile,
			BatchName:    student.BatchName,
			BatchID:      student.BatchID,
			Course:       student.Course,
			DomainName:   student.DomainName,
			Status:       student.Status,
			CreatedOn:    student.CreatedOn,
			ProfileImage: student.ProfileImage,
		}
	}

	resp := &toolapi.StudentListDTOResponse{
		Data:       dtos,
		TotalCount: fullResp.TotalCount,
	}

	// Cache result (1 minute)
	if data, err := json.Marshal(resp); err == nil {
		s.redis.Set(ctx, cacheKey, data, 1*time.Minute)
	}

	return resp, nil
}

func (s *ToolService) GetPageFilters(ctx context.Context, pageName string) (*toolapi.PageFiltersResponse, error) {
	cacheKey := "tool:filters:" + pageName

	val, err := s.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		var resp toolapi.PageFiltersResponse
		if err := json.Unmarshal([]byte(val), &resp); err == nil {
			return &resp, nil
		}
	}

	resp, err := s.client.GetPageFilters(pageName)
	if err != nil {
		return nil, err
	}

	if data, err := json.Marshal(resp); err == nil {
		s.redis.Set(ctx, cacheKey, data, 15*time.Minute)
	}

	return resp, nil
}

func (s *ToolService) GetBatches(ctx context.Context) (*toolapi.BatchResponse, error) {
	cacheKey := "tool:batches"

	val, err := s.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		var resp toolapi.BatchResponse
		if err := json.Unmarshal([]byte(val), &resp); err == nil {
			return &resp, nil
		}
	}

	resp, err := s.client.GetBatches()
	if err != nil {
		return nil, err
	}

	if data, err := json.Marshal(resp); err == nil {
		s.redis.Set(ctx, cacheKey, data, 15*time.Minute)
	}

	return resp, nil
}

func (s *ToolService) GetCourses(ctx context.Context) (*toolapi.CourseResponse, error) {
	cacheKey := "tool:courses"

	val, err := s.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		var resp toolapi.CourseResponse
		if err := json.Unmarshal([]byte(val), &resp); err == nil {
			return &resp, nil
		}
	}

	resp, err := s.client.GetCourses()
	if err != nil {
		return nil, err
	}

	if data, err := json.Marshal(resp); err == nil {
		s.redis.Set(ctx, cacheKey, data, 15*time.Minute)
	}

	return resp, nil
}

func (s *ToolService) GetDomains(ctx context.Context) (*toolapi.DomainResponse, error) {
	cacheKey := "tool:domains"

	val, err := s.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		var resp toolapi.DomainResponse
		if err := json.Unmarshal([]byte(val), &resp); err == nil {
			return &resp, nil
		}
	}

	resp, err := s.client.GetDomains()
	if err != nil {
		return nil, err
	}

	if data, err := json.Marshal(resp); err == nil {
		s.redis.Set(ctx, cacheKey, data, 15*time.Minute)
	}

	return resp, nil
}

func (s *ToolService) GetEmployees(ctx context.Context, roles string) (*toolapi.EmployeeResponse, error) {
	cacheKey := "tool:employees:" + roles

	val, err := s.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		var resp toolapi.EmployeeResponse
		if err := json.Unmarshal([]byte(val), &resp); err == nil {
			return &resp, nil
		}
	}

	resp, err := s.client.GetEmployees(roles)
	if err != nil {
		return nil, err
	}

	if data, err := json.Marshal(resp); err == nil {
		s.redis.Set(ctx, cacheKey, data, 15*time.Minute)
	}

	return resp, nil
}

func (s *ToolService) GetStatusOptions(ctx context.Context, category string) (*toolapi.StatusResponse, error) {
	cacheKey := "tool:status:" + category

	val, err := s.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		var resp toolapi.StatusResponse
		if err := json.Unmarshal([]byte(val), &resp); err == nil {
			return &resp, nil
		}
	}

	resp, err := s.client.GetStatusOptions(category)
	if err != nil {
		return nil, err
	}

	if data, err := json.Marshal(resp); err == nil {
		s.redis.Set(ctx, cacheKey, data, 15*time.Minute)
	}

	return resp, nil
}

// Helper to generate a consistent cache key from a map
func generateCacheKey(params map[string]string) string {
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var builder strings.Builder
	for _, k := range keys {
		builder.WriteString(k)
		builder.WriteString("=")
		builder.WriteString(params[k])
		builder.WriteString("&")
	}

	hash := sha256.Sum256([]byte(builder.String()))
	return hex.EncodeToString(hash[:])
}
