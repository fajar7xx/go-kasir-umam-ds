package handlers

import (
	"context"
	"fajar7xx/go-kasir-umam-ds/internal/services"
	"fajar7xx/go-kasir-umam-ds/utils"
	"net/http"
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
	utils.SendSuccess(w, map[string]string{
		"message": "Report fetched successfully",
	}, http.StatusOK)
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
