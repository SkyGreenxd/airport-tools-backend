package parse

import (
	"airport-tools-backend/internal/domain"
	"airport-tools-backend/pkg/e"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type CommonFilters struct {
	EmployeeId      *string
	StartDate       *time.Time
	EndDate         *time.Time
	Limit           *int
	AvgWorkDuration bool
	ErrorType       *string
}

func NewCommonFilters(employeeId *string, startDate, endDate *time.Time, limit *int, avgWorkDuration bool, errorType *string) *CommonFilters {
	return &CommonFilters{
		EmployeeId:      employeeId,
		StartDate:       startDate,
		EndDate:         endDate,
		Limit:           limit,
		AvgWorkDuration: avgWorkDuration,
		ErrorType:       errorType,
	}
}

// parseCommonFilters извлекает общие параметры фильтрации из запроса.
func parseCommonFilters(c *gin.Context) (*CommonFilters, error) {
	const op = "parse.parseCommonFilters"

	var startDate, endDate *time.Time
	var avgWorkDuration bool
	var employeeId *string
	var limit *int
	var errorType *string

	employeeIdStr := c.Query("employee_id")
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")
	limitStr := c.Query("limit")
	avgWorkDurationStr := c.Query("avg_work_duration")
	errorTypeStr := c.Query("error_type")

	validErrorType := func(e string) bool {
		return e == string(domain.ModelError) || e == string(domain.HumanError)
	}

	if validErrorType(errorTypeStr) {
		errorType = &errorTypeStr
	}

	if avgWorkDurationStr == "true" {
		avgWorkDuration = true
	} else if avgWorkDurationStr == "false" || avgWorkDurationStr == "" {
		avgWorkDuration = false
	} else {
		return nil, e.Wrap(op, e.ErrInvalidRequestBody)
	}

	if employeeIdStr != "" {
		employeeId = &employeeIdStr
	}

	if startDateStr != "" {
		t, err := time.Parse("02-01-2006", startDateStr)
		if err != nil {
			return nil, e.Wrap(op, err)
		}
		startDate = &t
	}

	if endDateStr != "" {
		t, err := time.Parse("02-01-2006", endDateStr)
		if err != nil {
			return nil, e.Wrap(op, err)
		}
		endDate = &t
	}

	if limitStr != "" {
		n, err := strconv.Atoi(limitStr)
		if err != nil || n <= 0 {
			return nil, e.Wrap(op, err)
		}
		limit = &n
	}

	return NewCommonFilters(employeeId, startDate, endDate, limit, avgWorkDuration, errorType), nil
}
