package helpers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"strings"
	"testing"

	imageHandler "github.com/Leviosa-care/leviosa/backend/internal/catalog/interface/image"
	ck "github.com/Leviosa-care/leviosa/backend/internal/common/auth/cookies"
	"github.com/stretchr/testify/require"
)

// ============================================================================
// Admin-Only Image Endpoints (All image endpoints require auth)
// ============================================================================

// NewUploadImageRequest creates a request to upload an image (admin endpoint)
// accessToken is optional - if empty, no auth cookie is added (for testing unauthorized access)
// fileContent is the image file bytes, fieldName is the form field name, fileName is the file name
func NewUploadImageRequest(t *testing.T, ctx context.Context, serverURL string, fileContent []byte, fieldName string, fileName string, additionalFields map[string]string, accessToken string) *http.Request {
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	// Determine content type from filename extension
	var contentType string
	lowerFileName := strings.ToLower(fileName)
	switch {
	case strings.HasSuffix(lowerFileName, ".jpg"), strings.HasSuffix(lowerFileName, ".jpeg"):
		contentType = "image/jpeg"
	case strings.HasSuffix(lowerFileName, ".png"):
		contentType = "image/png"
	case strings.HasSuffix(lowerFileName, ".gif"):
		contentType = "image/gif"
	default:
		contentType = "image/jpeg" // Default fallback
	}

	// Create multipart header with proper Content-Type
	// This is critical: Go's CreateFormFile() defaults to application/octet-stream
	// which fails server-side validation. We must use CreatePart() with explicit headers.
	header := make(textproto.MIMEHeader)
	header.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"; filename="%s"`, fieldName, fileName))
	header.Set("Content-Type", contentType)

	// Write file part with correct content type
	part, err := writer.CreatePart(header)
	require.NoError(t, err)
	_, err = io.Copy(part, bytes.NewReader(fileContent))
	require.NoError(t, err)

	// Add additional form fields (e.g., entity_type, entity_id)
	for key, val := range additionalFields {
		err = writer.WriteField(key, val)
		require.NoError(t, err)
	}

	err = writer.Close()
	require.NoError(t, err)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, serverURL+imageHandler.UploadImageEndpoint, &body)
	require.NoError(t, err)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	if accessToken != "" {
		cookie := &http.Cookie{
			Name:  ck.AccessTokenCookieName,
			Value: accessToken,
		}
		req.AddCookie(cookie)
	}

	return req
}

// NewRemoveImageRequest creates a request to remove an image (admin endpoint)
// accessToken is optional - if empty, no auth cookie is added (for testing unauthorized access)
func NewRemoveImageRequest(t *testing.T, ctx context.Context, serverURL string, requestBody interface{}, accessToken string) *http.Request {
	jsonBody, err := json.Marshal(requestBody)
	require.NoError(t, err)

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, serverURL+imageHandler.RemoveImageEndpoint, bytes.NewReader(jsonBody))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	if accessToken != "" {
		cookie := &http.Cookie{
			Name:  ck.AccessTokenCookieName,
			Value: accessToken,
		}
		req.AddCookie(cookie)
	}

	return req
}

// NewSetActiveImageRequest creates a request to set an active image for a resource (admin endpoint)
// accessToken is optional - if empty, no auth cookie is added (for testing unauthorized access)
func NewSetActiveImageRequest(t *testing.T, ctx context.Context, serverURL string, requestBody interface{}, accessToken string) *http.Request {
	jsonBody, err := json.Marshal(requestBody)
	require.NoError(t, err)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, serverURL+imageHandler.SetActiveImageEndpoint, bytes.NewReader(jsonBody))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	if accessToken != "" {
		cookie := &http.Cookie{
			Name:  ck.AccessTokenCookieName,
			Value: accessToken,
		}
		req.AddCookie(cookie)
	}

	return req
}
