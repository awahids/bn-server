package schoolhandler

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/awahids/bn-server/internal/delivery/data/request/schoolreq"
	"github.com/awahids/bn-server/internal/delivery/data/response"
	"github.com/awahids/bn-server/internal/delivery/handlers/handlerutil"
	"github.com/awahids/bn-server/internal/domain/models"
	appservice "github.com/awahids/bn-server/internal/domain/services/appservice"
	"github.com/awahids/bn-server/internal/domain/services/serviceinterface"

	"github.com/gin-gonic/gin"
)

type SchoolHandler struct {
	appService serviceinterface.AppService
}

type SchoolCreatedByResponse struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Username *string `json:"username,omitempty"`
}

type SchoolResponse struct {
	ID            string                  `json:"id"`
	Name          string                  `json:"name"`
	Location      string                  `json:"location"`
	Jenjang       string                  `json:"jenjang"`
	StatusSekolah string                  `json:"statusSekolah"`
	MonthlyFee    int                     `json:"monthlyFee"`
	MapURL        string                  `json:"mapUrl"`
	Contact       string                  `json:"contact"`
	Description   string                  `json:"description"`
	CreatedBy     SchoolCreatedByResponse `json:"createdBy"`
	CreatedAt     string                  `json:"createdAt"`
	UpdatedAt     string                  `json:"updatedAt"`
}

var schoolJenjangAliases = map[string]string{
	"tk":      "TK",
	"sd":      "SD",
	"smp":     "SMP",
	"sma":     "SMA",
	"smk":     "SMK",
	"mi":      "MI",
	"mts":     "MTs",
	"ma":      "MA",
	"lainnya": "Lainnya",
}

func NewSchoolHandler(appService serviceinterface.AppService) *SchoolHandler {
	return &SchoolHandler{appService: appService}
}

// GetSchools godoc
// @Summary Get schools
// @Description Get public list of schools.
// @Tags School
// @Produce json
// @Success 200 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /schools [get]
func (h *SchoolHandler) GetSchools(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), handlerutil.RequestTimeout)
	defer cancel()

	schools, err := h.appService.GetSchools(ctx)
	if err != nil {
		response.Failed(c, http.StatusInternalServerError, "failed to get schools", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "ok", toSchoolResponses(schools))
}

// PostSchool godoc
// @Summary Create school
// @Description Create a new school by authenticated user.
// @Tags School
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param payload body schoolreq.CreateSchoolRequest true "Create school payload"
// @Success 201 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /schools [post]
func (h *SchoolHandler) PostSchool(c *gin.Context) {
	userID, err := handlerutil.GetUserID(c)
	if err != nil {
		handlerutil.FailUnauthorized(c, err)
		return
	}

	var req schoolreq.CreateSchoolRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Failed(c, http.StatusBadRequest, "invalid request payload", err.Error())
		return
	}

	name := strings.TrimSpace(req.Name)
	if name == "" || len(name) > 191 {
		response.Failed(c, http.StatusBadRequest, "invalid request data", "name must be 1-191 characters")
		return
	}

	location := strings.TrimSpace(req.Location)
	if location == "" || len(location) > 255 {
		response.Failed(c, http.StatusBadRequest, "invalid request data", "location must be 1-255 characters")
		return
	}

	jenjang, jenjangOK := normalizeJenjang(req.Jenjang)
	if !jenjangOK {
		response.Failed(c, http.StatusBadRequest, "invalid request data", "jenjang must be one of TK, SD, SMP, SMA, SMK, MI, MTs, MA, Lainnya")
		return
	}

	statusSekolah, statusOK := normalizeStatusSekolah(req.StatusSekolah)
	if !statusOK {
		response.Failed(c, http.StatusBadRequest, "invalid request data", "statusSekolah must be one of: negeri, swasta")
		return
	}

	if req.MonthlyFee == nil || *req.MonthlyFee < 0 {
		response.Failed(c, http.StatusBadRequest, "invalid request data", "monthlyFee is required and must be >= 0")
		return
	}

	mapURL := strings.TrimSpace(req.MapURL)
	if !isValidMapURL(mapURL) {
		response.Failed(c, http.StatusBadRequest, "invalid request data", "mapUrl must be a valid http/https URL")
		return
	}
	if len(mapURL) > 1024 {
		response.Failed(c, http.StatusBadRequest, "invalid request data", "mapUrl must be at most 1024 characters")
		return
	}

	var contact *string
	if req.Contact != nil {
		trimmed := strings.TrimSpace(*req.Contact)
		if len(trimmed) > 100 {
			response.Failed(c, http.StatusBadRequest, "invalid request data", "contact must be at most 100 characters")
			return
		}
		contact = &trimmed
	}

	var description *string
	if req.Description != nil {
		trimmed := strings.TrimSpace(*req.Description)
		if len(trimmed) > 1000 {
			response.Failed(c, http.StatusBadRequest, "invalid request data", "description must be at most 1000 characters")
			return
		}
		description = &trimmed
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), handlerutil.RequestTimeout)
	defer cancel()

	school, err := h.appService.CreateSchool(ctx, userID, serviceinterface.CreateSchoolInput{
		Name:          name,
		Location:      location,
		Jenjang:       jenjang,
		StatusSekolah: statusSekolah,
		MonthlyFee:    *req.MonthlyFee,
		MapURL:        mapURL,
		Contact:       contact,
		Description:   description,
	})
	if err != nil {
		if errors.Is(err, appservice.ErrSchoolInvalidData) {
			response.Failed(c, http.StatusBadRequest, "invalid request data", err.Error())
			return
		}
		response.Failed(c, http.StatusInternalServerError, "failed to create school", err.Error())
		return
	}

	response.Success(c, http.StatusCreated, "created", toSchoolResponse(*school))
}

func isValidMapURL(value string) bool {
	if value == "" {
		return false
	}
	parsed, err := url.ParseRequestURI(value)
	if err != nil {
		return false
	}
	return (parsed.Scheme == "http" || parsed.Scheme == "https") && parsed.Host != ""
}

func toSchoolResponses(items []models.School) []SchoolResponse {
	result := make([]SchoolResponse, 0, len(items))
	for _, item := range items {
		result = append(result, toSchoolResponse(item))
	}
	return result
}

func toSchoolResponse(item models.School) SchoolResponse {
	createdByID := item.UserID
	if item.User.ID != "" {
		createdByID = item.User.ID
	}

	createdByName := strings.TrimSpace(item.User.Name)
	if createdByName == "" {
		createdByName = "Unknown User"
	}

	return SchoolResponse{
		ID:            item.ID,
		Name:          item.Name,
		Location:      item.Location,
		Jenjang:       item.Jenjang,
		StatusSekolah: item.StatusSekolah,
		MonthlyFee:    item.MonthlyFee,
		MapURL:        item.MapURL,
		Contact:       item.Contact,
		Description:   item.Description,
		CreatedBy: SchoolCreatedByResponse{
			ID:       createdByID,
			Name:     createdByName,
			Username: item.User.Username,
		},
		CreatedAt: item.CreatedAt.Format(time.RFC3339),
		UpdatedAt: item.UpdatedAt.Format(time.RFC3339),
	}
}

func normalizeJenjang(value string) (string, bool) {
	key := strings.ToLower(strings.TrimSpace(value))
	normalized, ok := schoolJenjangAliases[key]
	return normalized, ok
}

func normalizeStatusSekolah(value string) (string, bool) {
	status := strings.ToLower(strings.TrimSpace(value))
	switch status {
	case "negeri", "swasta":
		return status, true
	default:
		return "", false
	}
}
