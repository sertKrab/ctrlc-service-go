package testutil

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestRecordBEAPIResponseWritesRedactedCapture(t *testing.T) {
	root := t.TempDir()
	t.Setenv("CTRLC_PROJECT_ROOT", root)

	rec := httptest.NewRecorder()
	rec.Code = http.StatusOK
	rec.Body.WriteString(`{
		"success": true,
		"data": {
			"access_token": "secret-access",
			"refresh_token": "secret-refresh",
			"token_type": "Bearer",
			"nested": {"otp_code_hash": "secret-hash", "pin": "123456"}
		},
		"error": null
	}`)

	RecordBEAPIResponse(t, "AUTH_TC_034", rec)

	path := filepath.Join(root, ".ctrlc", "e2e", "api-responses", "be_AUTH_TC_034.json")
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("expected capture file at %s: %v", path, err)
	}

	var capture map[string]any
	if err := json.Unmarshal(raw, &capture); err != nil {
		t.Fatalf("capture is not valid JSON: %v", err)
	}

	if got := int(capture["status"].(float64)); got != http.StatusOK {
		t.Fatalf("status = %d, want %d", got, http.StatusOK)
	}

	body := capture["body"].(map[string]any)
	data := body["data"].(map[string]any)
	if data["access_token"] != redactedValue {
		t.Fatalf("access_token was not redacted: %#v", data["access_token"])
	}
	if data["refresh_token"] != redactedValue {
		t.Fatalf("refresh_token was not redacted: %#v", data["refresh_token"])
	}
	if data["token_type"] != "Bearer" {
		t.Fatalf("token_type = %#v, want Bearer", data["token_type"])
	}

	nested := data["nested"].(map[string]any)
	if nested["otp_code_hash"] != redactedValue {
		t.Fatalf("otp_code_hash was not redacted: %#v", nested["otp_code_hash"])
	}
	if nested["pin"] != redactedValue {
		t.Fatalf("pin was not redacted: %#v", nested["pin"])
	}
}
