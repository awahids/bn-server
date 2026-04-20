package publichandler

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/awahids/bn-server/internal/delivery/data/response"
	"github.com/awahids/bn-server/internal/domain/services/serviceinterface"

	"github.com/gin-gonic/gin"
)

const (
	publicRequestTimeout = 15 * time.Second
	audioUserAgent       = "Hijaiyah-App/1.0"
)

var allowedAudioDomains = []string{"localhost", "127.0.0.1"}
var dateRegex = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)

type PublicHandler struct {
	httpClient    *http.Client
	publicService serviceinterface.PublicService
}

func NewPublicHandler(publicService serviceinterface.PublicService) *PublicHandler {
	return &PublicHandler{
		httpClient:    &http.Client{Timeout: 12 * time.Second},
		publicService: publicService,
	}
}

// GetAudioProxy godoc
// @Summary Audio proxy
// @Description Proxy local/allowed audio URL for client playback.
// @Tags Public
// @Produce octet-stream
// @Param url query string true "Audio source URL or local path"
// @Success 200 {file} file
// @Failure 400 {object} response.APIResponse
// @Failure 403 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /audio-proxy [get]
func (h *PublicHandler) GetAudioProxy(c *gin.Context) {
	audioURL := strings.TrimSpace(c.Query("url"))
	if audioURL == "" {
		response.Failed(c, http.StatusBadRequest, "audio URL is required", "audio URL is required")
		return
	}

	isLocalFile := strings.HasPrefix(audioURL, "/") || strings.HasPrefix(audioURL, "./") || strings.HasPrefix(audioURL, "../")
	if isLocalFile {
		filePath := audioURL
		if !strings.HasPrefix(filePath, "/") {
			filePath = "/" + filePath
		}
		if strings.Contains(filePath, "..") || strings.Contains(filePath, "~") {
			response.Failed(c, http.StatusBadRequest, "invalid file path", "invalid file path")
			return
		}
		c.Redirect(http.StatusTemporaryRedirect, filePath)
		return
	}

	parsedURL, err := url.Parse(audioURL)
	if err != nil || parsedURL.Scheme == "" || parsedURL.Host == "" {
		response.Failed(c, http.StatusBadRequest, "invalid URL format", "invalid URL format")
		return
	}
	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		response.Failed(c, http.StatusBadRequest, "invalid URL format", "URL must use http or https")
		return
	}

	if !isAllowedAudioHost(parsedURL.Hostname()) {
		response.Failed(c, http.StatusForbidden, "domain not allowed", "domain not allowed")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), publicRequestTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, audioURL, nil)
	if err != nil {
		response.Failed(c, http.StatusBadRequest, "invalid URL format", err.Error())
		return
	}
	req.Header.Set("User-Agent", audioUserAgent)

	res, err := h.httpClient.Do(req)
	if err != nil {
		response.Failed(c, http.StatusInternalServerError, "failed to fetch audio file", err.Error())
		return
	}
	defer res.Body.Close()

	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusMultipleChoices {
		response.Failed(c, res.StatusCode, "failed to fetch audio file", "failed to fetch audio file")
		return
	}

	contentType := strings.TrimSpace(res.Header.Get("Content-Type"))
	if contentType == "" {
		contentType = "audio/mpeg"
	}

	c.Header("Content-Type", contentType)
	c.Header("Cache-Control", "public, max-age=86400")
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "GET, OPTIONS")
	c.Header("Access-Control-Allow-Headers", "Content-Type")
	c.Status(http.StatusOK)

	if _, err := io.Copy(c.Writer, res.Body); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
}

// GetPrayerTimes godoc
// @Summary Get prayer times
// @Description Get prayer times by coordinates and optional date/method.
// @Tags Public
// @Produce json
// @Param lat query number true "Latitude"
// @Param lng query number true "Longitude"
// @Param date query string false "Date in YYYY-MM-DD (default today)"
// @Param method query string false "Calculation method (default 2)"
// @Success 200 {object} map[string]any
// @Failure 400 {object} response.APIResponse
// @Router /prayer-times [get]
func (h *PublicHandler) GetPrayerTimes(c *gin.Context) {
	if h.publicService == nil {
		response.Failed(c, http.StatusInternalServerError, "prayer service is not configured", "prayer service is not configured")
		return
	}

	latParam := strings.TrimSpace(c.Query("lat"))
	lngParam := strings.TrimSpace(c.Query("lng"))
	if latParam == "" || lngParam == "" {
		response.Failed(c, http.StatusBadRequest, "latitude and longitude are required", "latitude and longitude are required")
		return
	}

	latitude, err := strconv.ParseFloat(latParam, 64)
	if err != nil {
		response.Failed(c, http.StatusBadRequest, "invalid latitude or longitude", "invalid latitude")
		return
	}
	longitude, err := strconv.ParseFloat(lngParam, 64)
	if err != nil {
		response.Failed(c, http.StatusBadRequest, "invalid latitude or longitude", "invalid longitude")
		return
	}

	if latitude < -90 || latitude > 90 || longitude < -180 || longitude > 180 {
		response.Failed(c, http.StatusBadRequest, "invalid latitude or longitude", "latitude must be between -90 and 90, longitude between -180 and 180")
		return
	}

	dateValue := strings.TrimSpace(c.Query("date"))
	if dateValue == "" {
		dateValue = time.Now().UTC().Format("2006-01-02")
	}
	if !dateRegex.MatchString(dateValue) {
		response.Failed(c, http.StatusBadRequest, "invalid date format", "date must be in YYYY-MM-DD format")
		return
	}

	method := strings.TrimSpace(c.Query("method"))
	if method == "" {
		method = "2"
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), publicRequestTimeout)
	defer cancel()

	prayerTimes, meta := h.publicService.GetPrayerTimes(ctx, latitude, longitude, dateValue, method)
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    prayerTimes,
		"meta":    meta,
	})
}

func isAllowedAudioHost(hostname string) bool {
	host := strings.ToLower(strings.TrimSpace(hostname))
	if host == "" {
		return false
	}

	for _, allowed := range allowedAudioDomains {
		candidate := strings.ToLower(strings.TrimSpace(allowed))
		if candidate == "" {
			continue
		}

		// Support wildcard rules like *.example.com without substring bypass.
		if strings.HasPrefix(candidate, "*.") {
			suffix := strings.TrimPrefix(candidate, "*.")
			if suffix != "" && (host == suffix || strings.HasSuffix(host, "."+suffix)) {
				return true
			}
			continue
		}

		if host == candidate {
			return true
		}
	}

	return false
}
