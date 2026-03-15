package handlers

import (
	"net/http"
	"strconv"
	"time"

	"habitflow/internal/services"

	"github.com/gin-gonic/gin"
)

// ReportHandler handles report and score requests.
type ReportHandler struct {
	ReportService      *services.ReportService
	ScoreService       *services.ScoreService
	InsightService     *services.InsightService
	DailyReportService *services.DailyReportService
}

// NewReportHandler creates a new ReportHandler.
func NewReportHandler(reportService *services.ReportService, scoreService *services.ScoreService, insightService *services.InsightService, dailyReportService *services.DailyReportService) *ReportHandler {
	return &ReportHandler{
		ReportService:      reportService,
		ScoreService:       scoreService,
		InsightService:     insightService,
		DailyReportService: dailyReportService,
	}
}

// Weekly handles GET /api/v1/reports/weekly
func (h *ReportHandler) Weekly(c *gin.Context) {
	userID := c.GetUint("user_id")

	var report *services.WeeklyReport
	var err error

	startQ := c.Query("start")
	endQ := c.Query("end")
	if startQ != "" && endQ != "" {
		report, err = h.ReportService.GenerateWeeklyForPeriod(userID, startQ, endQ)
	} else {
		report, err = h.ReportService.GenerateWeekly(userID)
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to generate weekly report: " + err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Weekly report generated",
		"data":    report,
	})
}

// Score handles GET /api/v1/reports/score
func (h *ReportHandler) Score(c *gin.Context) {
	userID := c.GetUint("user_id")

	days := 7
	if d := c.Query("days"); d != "" {
		if parsed, err := strconv.Atoi(d); err == nil && parsed > 0 && parsed <= 365 {
			days = parsed
		}
	}

	score, err := h.ScoreService.Calculate(userID, days)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to calculate score: " + err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Consistency score calculated",
		"data":    score,
	})
}

// Insights handles GET /api/v1/reports/insights
func (h *ReportHandler) Insights(c *gin.Context) {
	userID := c.GetUint("user_id")

	now := time.Now()
	weekday := int(now.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	monday := now.AddDate(0, 0, -(weekday - 1))
	startStr := monday.Format("2006-01-02")
	endStr := now.Format("2006-01-02")

	insights, err := h.InsightService.Generate(userID, startStr, endStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to generate insights: " + err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Insights generated",
		"data":    insights,
	})
}

// Daily handles GET /api/v1/reports/daily
func (h *ReportHandler) Daily(c *gin.Context) {
	userID := c.GetUint("user_id")

	var date time.Time
	if q := c.Query("date"); q != "" {
		parsed, err := time.ParseInLocation("2006-01-02", q, services.WIB)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Format tanggal tidak valid (gunakan YYYY-MM-DD)",
				"data":    nil,
			})
			return
		}
		date = parsed
	} else {
		date = time.Now().In(services.WIB)
	}

	report, err := h.DailyReportService.Generate(userID, date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal membuat laporan harian: " + err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Laporan harian berhasil dibuat",
		"data":    report,
	})
}
