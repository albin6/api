package toolapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/albin6/api/config"
)

type Client struct {
	httpClient *http.Client
	config     *config.Config
	token      string
	mu         sync.RWMutex
}

func NewClient(cfg *config.Config) *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		config: cfg,
	}
}

func (c *Client) Authenticate() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	reqBody := AuthRequest{
		Mobile:   c.config.ToolAPIMobile,
		Password: c.config.ToolAPIPassword,
		Token:    c.config.ToolAPIToken,
	}

	body, _ := json.Marshal(reqBody)
	req, err := http.NewRequest("POST", c.config.ToolAPIBaseURL+"/auth/verifySignIn", bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-transaction-id", c.config.ToolAPITransactionID)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		fmt.Printf("DEBUG: Authenticate request failed: %v\n", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("DEBUG: Authenticate failed with status: %d\n", resp.StatusCode)
		return fmt.Errorf("authentication failed: status %d", resp.StatusCode)
	}

	var authResp AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		return err
	}

	c.token = authResp.Data.Token
	return nil
}

func (c *Client) doRequest(method, path string, params map[string]string) (*http.Response, error) {
	c.mu.RLock()
	token := c.token
	c.mu.RUnlock()

	if token == "" {
		// Just in case authentication hasn't happened yet, though Service should handle this
		if err := c.Authenticate(); err != nil {
			return nil, err
		}
		c.mu.RLock()
		token = c.token
		c.mu.RUnlock()
	}

	req, err := http.NewRequest(method, c.config.ToolAPIBaseURL+path, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("x-transaction-id", c.config.ToolAPITransactionID)

	q := req.URL.Query()
	for k, v := range params {
		q.Add(k, v)
	}
	// Add transactionId if not present (using current time millis)
	if q.Get("transactionId") == "" {
		q.Add("transactionId", fmt.Sprintf("%d", time.Now().UnixMilli()))
	}
	req.URL.RawQuery = q.Encode()

	resp, err := c.httpClient.Do(req)
	if err != nil {
		fmt.Printf("DEBUG: Request to %s failed: %v\n", path, err)
		return nil, err
	}

	// Simple token expiration handling (if 401, re-auth and retry once)
	if resp.StatusCode == http.StatusUnauthorized {
		resp.Body.Close()
		fmt.Println("DEBUG: Received 401 Unauthorized, re-authenticating...")
		if err := c.Authenticate(); err != nil {
			return nil, err
		}

		c.mu.RLock()
		token = c.token
		c.mu.RUnlock()

		req.Header.Set("Authorization", "Bearer "+token)
		return c.httpClient.Do(req)
	}

	if resp.StatusCode >= 400 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		resp.Body.Close()                                    // Close original body
		resp.Body = io.NopCloser(bytes.NewBuffer(bodyBytes)) // Restore body for caller

		fmt.Printf("DEBUG: Request to %s failed with status %d. Body: %s\n", path, resp.StatusCode, string(bodyBytes))
	}

	return resp, nil
}

func (c *Client) GetStudents(params map[string]string) (*StudentListResponse, error) {
	resp, err := c.doRequest("GET", "/student", params)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get students: status %d", resp.StatusCode)
	}

	var result StudentListResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *Client) GetBatches() (*BatchResponse, error) {
	params := map[string]string{"batchTypes": "1"}
	resp, err := c.doRequest("GET", "/batch", params)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get batches: status %d", resp.StatusCode)
	}

	var result BatchResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) GetCourses() (*CourseResponse, error) {
	resp, err := c.doRequest("GET", "/course", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get courses: status %d", resp.StatusCode)
	}

	var result CourseResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) GetDomains() (*DomainResponse, error) {
	resp, err := c.doRequest("GET", "/domain", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get domains: status %d", resp.StatusCode)
	}

	var result DomainResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) GetPageFilters(pageName string) (*PageFiltersResponse, error) {
	params := map[string]string{"pageName": pageName}
	resp, err := c.doRequest("GET", "/common/page-filters", params)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get page filters: status %d", resp.StatusCode)
	}

	var result PageFiltersResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) GetStatusOptions(category string) (*StatusResponse, error) {
	params := map[string]string{"category": category}
	resp, err := c.doRequest("GET", "/common/status", params)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get status options: status %d", resp.StatusCode)
	}

	var result StatusResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) GetEmployees(roles string) (*EmployeeResponse, error) {
	params := map[string]string{"roles": roles}
	resp, err := c.doRequest("GET", "/employee/roles", params)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get employees: status %d", resp.StatusCode)
	}

	var result EmployeeResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return &result, nil
}
