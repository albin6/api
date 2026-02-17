package toolapi

import "time"

// AuthRequest represents the request body for verifySignIn
type AuthRequest struct {
	Mobile   string `json:"mobile"`
	Password string `json:"password"`
	Token    string `json:"token"`
}

// AuthResponse represents the response from verifySignIn
type AuthResponse struct {
	Data struct {
		Token string `json:"token"`
		User  struct {
			Name       string `json:"name"`
			Mobile     string `json:"mobile"`
			Role       string `json:"role"`
			Email      string `json:"email"`
			EmployeeID string `json:"employeeId"`
		} `json:"user"`
	} `json:"data"`
}

// Student represents a student record from the external API
type Student struct {
	ID            string    `json:"id"`
	AdmissionID   string    `json:"admissionId"`
	Name          string    `json:"name"`
	Email         string    `json:"email"`
	Mobile        string    `json:"mobile"`
	BatchName     string    `json:"batchName"`
	BatchID       string    `json:"batchId"`
	Course        string    `json:"course"`
	Status        string    `json:"status"`
	Stage         string    `json:"stage"`
	DomainName    string    `json:"domainName"`
	HubName       string    `json:"hubName"`
	AdvisorName   string    `json:"advisorName"`
	PaymentMethod string    `json:"paymentMethod"`
	CreatedOn     time.Time `json:"createdOn"`
	ProfileImage  string    `json:"profileImage,omitempty"` // Assuming image might be there or not
	TotalCount    string    `json:"totalCount"`
}

// StudentDTO represents the lean student record for the frontend
type StudentDTO struct {
	ID           string    `json:"id"`
	AdmissionID  string    `json:"admissionId"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	Mobile       string    `json:"mobile"`
	BatchName    string    `json:"batchName"`
	BatchID      string    `json:"batchId"`
	Course       string    `json:"course"`
	DomainName   string    `json:"domainName"`
	Status       string    `json:"status"`
	CreatedOn    time.Time `json:"createdOn"`
	ProfileImage string    `json:"profileImage,omitempty"`
}

// StudentListResponse represents the full response from external API
type StudentListResponse struct {
	Data       []Student `json:"data"`
	TotalCount string    `json:"totalCount"`
}

// StudentListDTOResponse represents the lean response for the frontend
type StudentListDTOResponse struct {
	Data       []StudentDTO `json:"data"`
	TotalCount string       `json:"totalCount"`
}

// FilterOption represents a generic filter option
type FilterOption struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Batch represents a batch filter option
type Batch struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	// Add other fields if necessary
}

// BatchResponse represents the response for getBatches
type BatchResponse struct {
	Data       []Batch `json:"data"`
	TotalCount int     `json:"totalCount"`
}

// Course represents a course filter option
type Course struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// CourseResponse represents the response for getCourses
type CourseResponse struct {
	Data       []Course `json:"data"`
	TotalCount int      `json:"totalCount"`
}

// Domain represents a domain filter option
type Domain struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// DomainResponse represents the response for getDomains
type DomainResponse struct {
	Data       []Domain `json:"data"`
	TotalCount int      `json:"totalCount"`
}

// Employee represents an employee (e.g. Batch Associate)
type Employee struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Role string `json:"role"`
}

// EmployeeResponse represents the response for getEmployees
type EmployeeResponse struct {
	Data []Employee `json:"data"`
}

// StatusOption represents a status option
type StatusOption struct {
	ID      string `json:"id"`
	Code    int    `json:"code"`
	Meaning string `json:"meaning"`
}

// StatusResponse represents the response for getStatus
type StatusResponse struct {
	Data []StatusOption `json:"data"`
}

// AvailableFilter represents an available filter key
type AvailableFilter string

// PageFiltersResponse represents the available filters for a page
type PageFiltersResponse struct {
	Data       []AvailableFilter `json:"data"`
	TotalCount int               `json:"totalCount"`
}
