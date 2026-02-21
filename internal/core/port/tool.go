package port

import (
	"context"

	"github.com/albin6/api/pkg/toolapi"
)

type ToolService interface {
	GetStudents(ctx context.Context, params map[string]string) (*toolapi.StudentListDTOResponse, error)
	GetPageFilters(ctx context.Context, pageName string) (*toolapi.PageFiltersResponse, error)
	GetBatches(ctx context.Context) (*toolapi.BatchResponse, error)
	GetCourses(ctx context.Context) (*toolapi.CourseResponse, error)
	GetDomains(ctx context.Context) (*toolapi.DomainResponse, error)
	GetEmployees(ctx context.Context, roles string) (*toolapi.EmployeeResponse, error)
	GetStatusOptions(ctx context.Context, category string) (*toolapi.StatusResponse, error)
}
