package handlers

import (
	"context"
	"errors"
	"fajar7xx/go-kasir-umam-ds/internal/services"
	"fajar7xx/go-kasir-umam-ds/utils"
	"net/http"
	"strings"
)

type ReportHandler struct {
	reportService services.ReportService
}

func NewReportHandler(reportService services.ReportService) *ReportHandler {
	return &ReportHandler{
		reportService: reportService,
	}
}

func (h *ReportHandler) HandleReport(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.GetReport(w, r)
	default:
		utils.SendError(w, "METHOD_NOT_ALLOWED", "Method not allowe", http.StatusMethodNotAllowed)
	}
}

func (h *ReportHandler) HandleReportByPeriod(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.GetReportByPeriod(w, r)
	default:
		utils.SendError(w, "METHOD_NOT_ALLOWED", "Method not allowe", http.StatusMethodNotAllowed)
	}
}

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

func (h *ReportHandler) GetReportByPeriod(w http.ResponseWriter, r *http.Request) {
	period := r.PathValue("period")

	ctx, cancel := context.WithTimeout(r.Context(), requestTimeout)
	defer cancel()

	// Handle "hari-ini" atau "today"
	if period == "hari-ini" || period == "today" {
		report, err := h.reportService.GetTodayReport(ctx)
		if err != nil {
			utils.SendError(w, "INTERNAL_ERROR", err.Error(), http.StatusInternalServerError)
			return
		}
		utils.SendSuccess(w, report, http.StatusOK)
		return
	}

	// Future: bisa extend untuk "week", "month", "year"
	utils.SendError(w, "NOT_FOUND", "Period not supported. Use: hari-ini, today", http.StatusNotFound)
}

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
