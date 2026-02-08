# Enhance GetReport with Date Range Query Parameters

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Enable `GetReport` endpoint to accept optional `start_date` and `end_date` query parameters, defaulting to today's report when no parameters are provided.

**Architecture:** Extend the existing `GetReport` handler to parse and validate query parameters. When both parameters are present, delegate to the existing `GetReportByDateRange` service method. When neither is present, delegate to `GetTodayReport`. Implement specific error codes for validation failures.

**Tech Stack:** Go 1.x, net/http, testify for testing

---

## Task 1: Add Custom Error Types to Service Layer

**Files:**
- Modify: `internal/services/report_service.go`

**Step 1: Define custom error types**

Add these custom error types at the top of the service file after imports:

```go
// Custom error types for report validation
var (
	ErrInvalidDateFormat = errors.New("invalid date format")
	ErrInvalidDateRange  = errors.New("end_date must be after start_date")
)
```

**Step 2: Import errors package if not already present**

Ensure `errors` is imported:

```go
import (
    "context"
    "errors"
    "fajar7xx/go-kasir-umam-ds/internal/repositories"
    "fajar7xx/go-kasir-umam-ds/models"
    "fmt"
    "time"
)
```

**Step 3: Update GetReportByDateRange to return custom errors**

Modify the error handling to wrap custom errors:

```go
func (s *reportService) GetReportByDateRange(ctx context.Context, startDate, endDate string) (*models.ReportResponse, error) {
	layout := "2006-01-02"

	start, err := time.Parse(layout, startDate)
	if err != nil {
		return nil, fmt.Errorf("%w: start_date must be YYYY-MM-DD format", ErrInvalidDateFormat)
	}

	end, err := time.Parse(layout, endDate)
	if err != nil {
		return nil, fmt.Errorf("%w: end_date must be YYYY-MM-DD format", ErrInvalidDateFormat)
	}

	if end.Before(start) {
		return nil, ErrInvalidDateRange
	}

	end = end.Add(24 * time.Hour)

	return s.reportRepo.GetReportByDateRange(ctx, start, end)
}
```

**Step 4: Commit service layer changes**

```bash
git add internal/services/report_service.go
git commit -m "feat: add custom error types for report date validation"
```

---

## Task 2: Write Tests for GetReport Handler

**Files:**
- Modify: `handlers/report_handler_test.go`

**Step 1: Write test for GetReport with no query parameters (defaults to today)**

Add this test function after existing tests:

```go
func TestReportHandler_GetReport_NoParams_DefaultsToToday(t *testing.T) {
	mockService := new(mocks.ReportServiceMock)
	handler := NewReportHandler(mockService)

	expectedReport := &models.ReportResponse{
		TotalRevenue:      250000.0,
		TotalTransactions: 10,
		BestSeller: &models.BestSeller{
			Name:         "Nasi Goreng",
			QuantitySold: 20,
		},
	}

	mockService.On("GetTodayReport", mock.Anything).Return(expectedReport, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/reports", nil)
	w := httptest.NewRecorder()

	handler.GetReport(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)

	data := response["data"].(map[string]interface{})
	assert.Equal(t, 250000.0, data["total_revenue"])
	assert.Equal(t, float64(10), data["total_transactions"])

	mockService.AssertExpectations(t)
}
```

**Step 2: Run test to verify it fails**

```bash
go test -v ./handlers -run TestReportHandler_GetReport_NoParams_DefaultsToToday
```

Expected: FAIL - GetReport doesn't call GetTodayReport yet

**Step 3: Write test for GetReport with both parameters**

Add this test:

```go
func TestReportHandler_GetReport_WithBothParams_Success(t *testing.T) {
	mockService := new(mocks.ReportServiceMock)
	handler := NewReportHandler(mockService)

	expectedReport := &models.ReportResponse{
		TotalRevenue:      1500000.0,
		TotalTransactions: 45,
		BestSeller: &models.BestSeller{
			Name:         "Mie Ayam",
			QuantitySold: 60,
		},
	}

	mockService.On("GetReportByDateRange", mock.Anything, "2026-01-01", "2026-02-01").
		Return(expectedReport, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/reports?start_date=2026-01-01&end_date=2026-02-01", nil)
	w := httptest.NewRecorder()

	handler.GetReport(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)

	data := response["data"].(map[string]interface{})
	assert.Equal(t, 1500000.0, data["total_revenue"])
	assert.Equal(t, float64(45), data["total_transactions"])

	mockService.AssertExpectations(t)
}
```

**Step 4: Run test to verify it fails**

```bash
go test -v ./handlers -run TestReportHandler_GetReport_WithBothParams_Success
```

Expected: FAIL - GetReport doesn't handle query parameters yet

**Step 5: Write test for missing one parameter**

Add this test:

```go
func TestReportHandler_GetReport_MissingEndDate_Error(t *testing.T) {
	mockService := new(mocks.ReportServiceMock)
	handler := NewReportHandler(mockService)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/reports?start_date=2026-01-01", nil)
	w := httptest.NewRecorder()

	handler.GetReport(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	var response utils.ErrorResponse
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, "MISSING_PARAMETERS", response.Error.Code)
	assert.Contains(t, response.Error.Message, "Both start_date and end_date are required")

	mockService.AssertNotCalled(t, "GetTodayReport")
	mockService.AssertNotCalled(t, "GetReportByDateRange")
}

func TestReportHandler_GetReport_MissingStartDate_Error(t *testing.T) {
	mockService := new(mocks.ReportServiceMock)
	handler := NewReportHandler(mockService)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/reports?end_date=2026-02-01", nil)
	w := httptest.NewRecorder()

	handler.GetReport(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	var response utils.ErrorResponse
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, "MISSING_PARAMETERS", response.Error.Code)

	mockService.AssertNotCalled(t, "GetTodayReport")
	mockService.AssertNotCalled(t, "GetReportByDateRange")
}
```

**Step 6: Run tests to verify they fail**

```bash
go test -v ./handlers -run "TestReportHandler_GetReport_Missing"
```

Expected: FAIL - GetReport doesn't validate parameters yet

**Step 7: Write test for invalid date format**

Add this test:

```go
func TestReportHandler_GetReport_InvalidDateFormat_Error(t *testing.T) {
	mockService := new(mocks.ReportServiceMock)
	handler := NewReportHandler(mockService)

	mockService.On("GetReportByDateRange", mock.Anything, "invalid-date", "2026-02-01").
		Return(nil, fmt.Errorf("%w: start_date must be YYYY-MM-DD format", errors.New("invalid date format")))

	req := httptest.NewRequest(http.MethodGet, "/api/v1/reports?start_date=invalid-date&end_date=2026-02-01", nil)
	w := httptest.NewRecorder()

	handler.GetReport(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	var response utils.ErrorResponse
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, "INVALID_DATE_FORMAT", response.Error.Code)
	assert.Contains(t, response.Error.Message, "start_date must be YYYY-MM-DD format")

	mockService.AssertExpectations(t)
}
```

**Step 8: Write test for invalid date range**

Add this test:

```go
func TestReportHandler_GetReport_InvalidDateRange_Error(t *testing.T) {
	mockService := new(mocks.ReportServiceMock)
	handler := NewReportHandler(mockService)

	mockService.On("GetReportByDateRange", mock.Anything, "2026-02-01", "2026-01-01").
		Return(nil, errors.New("end_date must be after start_date"))

	req := httptest.NewRequest(http.MethodGet, "/api/v1/reports?start_date=2026-02-01&end_date=2026-01-01", nil)
	w := httptest.NewRecorder()

	handler.GetReport(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	var response utils.ErrorResponse
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, "INVALID_DATE_RANGE", response.Error.Code)
	assert.Contains(t, response.Error.Message, "end_date must be after start_date")

	mockService.AssertExpectations(t)
}
```

**Step 9: Run all new tests to verify they fail**

```bash
go test -v ./handlers -run TestReportHandler_GetReport
```

Expected: Multiple FAIL - implementation not done yet

**Step 10: Commit test files**

```bash
git add handlers/report_handler_test.go
git commit -m "test: add comprehensive tests for GetReport with date range params"
```

---

## Task 3: Implement GetReport Handler Logic

**Files:**
- Modify: `handlers/report_handler.go`

**Step 1: Add imports for errors package**

Ensure `errors` and `strings` are imported at the top:

```go
import (
	"context"
	"errors"
	"fajar7xx/go-kasir-umam-ds/internal/services"
	"fajar7xx/go-kasir-umam-ds/utils"
	"net/http"
	"strings"
)
```

**Step 2: Implement error handler helper method**

Add this private method after the handler functions:

```go
func (h *ReportHandler) handleReportError(w http.ResponseWriter, err error) {
	// Check for custom service errors
	if errors.Is(err, services.ErrInvalidDateFormat) {
		utils.SendError(w, "INVALID_DATE_FORMAT", err.Error(), http.StatusBadRequest)
		return
	}
	
	if errors.Is(err, services.ErrInvalidDateRange) {
		utils.SendError(w, "INVALID_DATE_RANGE", err.Error(), http.StatusBadRequest)
		return
	}

	// Check for specific error messages in error string as fallback
	errMsg := err.Error()
	if strings.Contains(errMsg, "invalid date format") || strings.Contains(errMsg, "YYYY-MM-DD format") {
		utils.SendError(w, "INVALID_DATE_FORMAT", errMsg, http.StatusBadRequest)
		return
	}

	if strings.Contains(errMsg, "end_date must be after start_date") {
		utils.SendError(w, "INVALID_DATE_RANGE", errMsg, http.StatusBadRequest)
		return
	}

	// Generic internal error
	utils.SendError(w, "INTERNAL_ERROR", err.Error(), http.StatusInternalServerError)
}
```

**Step 3: Implement GetReport handler logic**

Replace the existing `GetReport` function with:

```go
func (h *ReportHandler) GetReport(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), requestTimeout)
	defer cancel()

	queryParams := r.URL.Query()
	startDate := queryParams.Get("start_date")
	endDate := queryParams.Get("end_date")

	// Scenario 1: No parameters → default to today's report
	if startDate == "" && endDate == "" {
		report, err := h.reportService.GetTodayReport(ctx)
		if err != nil {
			utils.SendError(w, "INTERNAL_ERROR", err.Error(), http.StatusInternalServerError)
			return
		}
		utils.SendSuccess(w, report, http.StatusOK)
		return
	}

	// Scenario 2: Only one parameter → error
	if startDate == "" || endDate == "" {
		utils.SendError(w, "MISSING_PARAMETERS",
			"Both start_date and end_date are required",
			http.StatusBadRequest)
		return
	}

	// Scenario 3: Both parameters provided → get report by date range
	report, err := h.reportService.GetReportByDateRange(ctx, startDate, endDate)
	if err != nil {
		h.handleReportError(w, err)
		return
	}

	utils.SendSuccess(w, report, http.StatusOK)
}
```

**Step 4: Run all GetReport tests**

```bash
go test -v ./handlers -run TestReportHandler_GetReport
```

Expected: All tests PASS

**Step 5: Run full handler test suite**

```bash
go test -v ./handlers
```

Expected: All tests PASS

**Step 6: Commit handler implementation**

```bash
git add handlers/report_handler.go
git commit -m "feat: implement GetReport with date range query parameters"
```

---

## Task 4: Manual Testing and Documentation

**Files:**
- Create: `docs/api/report-endpoints.md`

**Step 1: Start the application**

```bash
go run main.go
```

Expected: Server starts on configured port

**Step 2: Test endpoint with no parameters (default to today)**

```bash
curl -X GET "http://localhost:8080/api/v1/reports"
```

Expected: JSON response with today's report data

**Step 3: Test endpoint with valid date range**

```bash
curl -X GET "http://localhost:8080/api/v1/reports?start_date=2026-01-01&end_date=2026-02-01"
```

Expected: JSON response with report data for January 2026

**Step 4: Test endpoint with missing parameter**

```bash
curl -X GET "http://localhost:8080/api/v1/reports?start_date=2026-01-01"
```

Expected: 400 error with code "MISSING_PARAMETERS"

**Step 5: Test endpoint with invalid date format**

```bash
curl -X GET "http://localhost:8080/api/v1/reports?start_date=01-01-2026&end_date=2026-02-01"
```

Expected: 400 error with code "INVALID_DATE_FORMAT"

**Step 6: Test endpoint with invalid date range**

```bash
curl -X GET "http://localhost:8080/api/v1/reports?start_date=2026-02-01&end_date=2026-01-01"
```

Expected: 400 error with code "INVALID_DATE_RANGE"

**Step 7: Create API documentation**

Create file `docs/api/report-endpoints.md`:

```markdown
# Report API Endpoints

## Get Report

Retrieve sales report data with optional date range filtering.

### Endpoint

```
GET /api/v1/reports
```

### Query Parameters

| Parameter | Type | Required | Format | Description |
|-----------|------|----------|--------|-------------|
| start_date | string | conditional | YYYY-MM-DD | Start date for report range |
| end_date | string | conditional | YYYY-MM-DD | End date for report range |

**Note:** Both `start_date` and `end_date` must be provided together, or neither.

### Behavior

1. **No parameters**: Returns today's report
2. **Both parameters**: Returns report for specified date range
3. **One parameter**: Returns error (both required)

### Response Success (200 OK)

```json
{
  "data": {
    "total_revenue": 1500000.0,
    "total_transactions": 45,
    "best_seller": {
      "name": "Mie Ayam",
      "quantity_sold": 60
    }
  }
}
```

### Response Errors

#### Missing Parameters (400 Bad Request)

```json
{
  "error": {
    "code": "MISSING_PARAMETERS",
    "message": "Both start_date and end_date are required"
  }
}
```

#### Invalid Date Format (400 Bad Request)

```json
{
  "error": {
    "code": "INVALID_DATE_FORMAT",
    "message": "invalid date format: start_date must be YYYY-MM-DD format"
  }
}
```

#### Invalid Date Range (400 Bad Request)

```json
{
  "error": {
    "code": "INVALID_DATE_RANGE",
    "message": "end_date must be after start_date"
  }
}
```

### Examples

#### Get today's report
```bash
curl -X GET "http://localhost:8080/api/v1/reports"
```

#### Get report for specific date range
```bash
curl -X GET "http://localhost:8080/api/v1/reports?start_date=2026-01-01&end_date=2026-02-01"
```

## Get Report By Period

Retrieve sales report for predefined periods.

### Endpoint

```
GET /api/v1/reports/{period}
```

### Path Parameters

| Parameter | Type | Required | Values | Description |
|-----------|------|----------|--------|-------------|
| period | string | yes | hari-ini, today | Report period |

### Response

Same as Get Report endpoint.

### Example

```bash
curl -X GET "http://localhost:8080/api/v1/reports/hari-ini"
```
```

**Step 8: Update docs/README.md to reference new API doc**

Add to the API section:

```markdown
- [Report Endpoints](api/report-endpoints.md) - Report API documentation
```

**Step 9: Commit documentation**

```bash
git add docs/api/report-endpoints.md docs/README.md
git commit -m "docs: add comprehensive report endpoints documentation"
```

---

## Task 5: Final Verification

**Step 1: Run full test suite**

```bash
go test ./... -v
```

Expected: All tests PASS

**Step 2: Run test coverage check**

```bash
go test ./handlers -coverprofile=coverage.out
go tool cover -func=coverage.out
```

Expected: Coverage for report_handler.go shows high coverage

**Step 3: Check code formatting**

```bash
go fmt ./...
```

Expected: No changes needed

**Step 4: Run go vet**

```bash
go vet ./...
```

Expected: No issues reported

**Step 5: Final commit and summary**

```bash
git log --oneline -5
```

Expected: See all 5 commits from this implementation

---

## Summary

This implementation adds flexible date range query parameters to the `GetReport` endpoint:

- ✅ No parameters → defaults to today's report
- ✅ Both parameters → returns report for date range
- ✅ One parameter → returns validation error
- ✅ Invalid format → returns INVALID_DATE_FORMAT error
- ✅ Invalid range → returns INVALID_DATE_RANGE error
- ✅ Comprehensive test coverage
- ✅ API documentation

**Total Commits:** 5
**Files Modified:** 3 (service, handler, test)
**Files Created:** 1 (API docs)
