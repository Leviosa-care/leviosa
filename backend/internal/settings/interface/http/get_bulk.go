package http

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/Leviosa-care/leviosa/backend/internal/common/contracts/settings"
	"github.com/Leviosa-care/leviosa/backend/internal/common/ctxutil"
	"github.com/Leviosa-care/leviosa/backend/internal/common/httpx"
	"github.com/hengadev/errsx"
)

func (h *handler) BulkSettingsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	logger.InfoContext(ctx, "Handler: Processing get bulk settings request",
		"operation", "get_bulk_settings",
		"method", r.Method,
		"path", r.URL.Path)

	keysParam := r.URL.Query().Get("keys")
	if keysParam == "" {
		http.Error(w, "keys parameter required", http.StatusBadRequest)
		return
	}

	keys := strings.Split(keysParam, ",")

	var response []*settings.SettingDTO

	var errs errsx.Map

	for _, key := range keys {
		switch key {
		case settings.CompanyName:
			res, err := h.svc.GetCompanyName(ctx)
			if err != nil {
				errs.Set(settings.CompanyName, err)
				continue
			}
			response = append(response, &settings.SettingDTO{
				Key:   key,
				Value: res.Name,
			})

		case settings.CompanyLogo:
			res, err := h.svc.GetCompanyLogo(ctx)
			if err != nil {
				errs.Set(settings.CompanyLogo, err)
				continue
			}
			response = append(response, &settings.SettingDTO{
				Key:   key,
				Value: res.LogoURL,
			})

		case settings.CompanyEmail:
			res, err := h.svc.GetCompanyEmail(ctx)
			if err != nil {
				errs.Set(settings.CompanyEmail, err)
				continue
			}
			response = append(response, &settings.SettingDTO{
				Key:   key,
				Value: res.Email,
			})

		case settings.CompanyPhone:
			res, err := h.svc.GetCompanyTelephone(ctx)
			if err != nil {
				errs.Set(settings.CompanyPhone, err)
				continue
			}
			response = append(response, &settings.SettingDTO{
				Key:   key,
				Value: res.Telephone,
			})

		case settings.CompanyLegalAddress:
			res, err := h.svc.GetCompanyLegalAddress(ctx)
			if err != nil {
				errs.Set(settings.CompanyLegalAddress, err)
				continue
			}
			response = append(response, &settings.SettingDTO{
				Key:   key,
				Value: res.Address,
			})

		case settings.CompanyInstagram:
			res, err := h.svc.GetCompanyInstagram(ctx)
			if err != nil {
				errs.Set(settings.CompanyInstagram, err)
				continue
			}
			response = append(response, &settings.SettingDTO{
				Key:   key,
				Value: res.Instagram,
			})

		case settings.OTPDuration:
			res, err := h.svc.GetOTPDuration(ctx)
			if err != nil {
				errs.Set(settings.OTPDuration, err)
				continue
			}
			response = append(response, &settings.SettingDTO{
				Key:   key,
				Value: strconv.Itoa(res.Duration),
			})

		case settings.OTPLength:
			res, err := h.svc.GetOTPLength(ctx)
			if err != nil {
				errs.Set(settings.OTPLength, err)
				continue
			}
			response = append(response, &settings.SettingDTO{
				Key:   key,
				Value: strconv.Itoa(res.Length),
			})

		case settings.OTPMaxAttempts:
			res, err := h.svc.GetOTPMaxAttempts(ctx)
			if err != nil {
				errs.Set(settings.OTPMaxAttempts, err)
				continue
			}
			response = append(response, &settings.SettingDTO{
				Key:   key,
				Value: strconv.Itoa(res.MaxAttempts),
			})
		default:
			errs.Set(key, fmt.Errorf("invalid key: %s", key))
			continue
		}

	}

	if !errs.IsEmpty() {
		logger.InfoContext(ctx, "Handler: Get bulk settings completed with partial errors",
			"operation", "get_bulk_settings",
			"status_code", http.StatusMultiStatus,
			"success_count", len(response),
			"error_count", len(errs))

		httpx.RespondWithJSON(w, map[string]any{
			"data":   response,
			"errors": errs,
		}, http.StatusMultiStatus)
		return
	}

	logger.InfoContext(ctx, "Handler: Get bulk settings completed",
		"operation", "get_bulk_settings",
		"status_code", http.StatusOK,
		"settings_count", len(response))

	httpx.RespondWithJSON(w, response, http.StatusOK)
}
