package testutil

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const redactedValue = "[REDACTED]"

// RecordBEAPIResponse writes the primary asserted backend HTTP response for a
// testcase to .ctrlc/e2e/api-responses/be_<tc_id>.json.
//
// It mirrors the frontend capture format:
//
//	{"status": <http status>, "body": <json object/array or text>}
//
// The helper is intentionally test-only in usage but lives in a normal Go file
// so generated tests in any package can import it.
func RecordBEAPIResponse(t testing.TB, tcID string, recorder *httptest.ResponseRecorder) {
	t.Helper()
	if recorder == nil {
		t.Fatal("RecordBEAPIResponse: nil response recorder")
	}

	capture := map[string]any{
		"status": recorder.Code,
		"body":   parseAndRedactBody(recorder.Body.Bytes()),
	}

	out, err := json.MarshalIndent(capture, "", "  ")
	if err != nil {
		t.Fatalf("RecordBEAPIResponse: marshal capture: %v", err)
	}

	dir := filepath.Join(projectRoot(t), ".ctrlc", "e2e", "api-responses")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatalf("RecordBEAPIResponse: create capture dir: %v", err)
	}

	path := filepath.Join(dir, "be_"+safeTCID(tcID)+".json")
	if err := os.WriteFile(path, append(out, '\n'), 0o644); err != nil {
		t.Fatalf("RecordBEAPIResponse: write %s: %v", path, err)
	}
}

func parseAndRedactBody(body []byte) any {
	body = bytes.TrimSpace(body)
	if len(body) == 0 {
		return ""
	}

	var value any
	dec := json.NewDecoder(bytes.NewReader(body))
	dec.UseNumber()
	if err := dec.Decode(&value); err != nil {
		return string(body)
	}
	return redactJSON(value)
}

func redactJSON(value any) any {
	switch v := value.(type) {
	case map[string]any:
		out := make(map[string]any, len(v))
		for key, child := range v {
			if isSensitiveKey(key) {
				out[key] = redactedValue
				continue
			}
			out[key] = redactJSON(child)
		}
		return out
	case []any:
		out := make([]any, len(v))
		for i, child := range v {
			out[i] = redactJSON(child)
		}
		return out
	default:
		return v
	}
}

func isSensitiveKey(key string) bool {
	normalized := strings.NewReplacer("_", "", "-", "", ".", "").Replace(strings.ToLower(key))
	switch normalized {
	case "accesstoken",
		"refreshtoken",
		"idtoken",
		"token",
		"tokenhash",
		"otp",
		"otpcode",
		"otpcodehash",
		"pin",
		"pinconfirm",
		"password",
		"secret",
		"authorization",
		"cookie",
		"setcookie",
		"jwt",
		"apikey":
		return true
	default:
		return false
	}
}

func safeTCID(tcID string) string {
	tcID = strings.TrimSpace(tcID)
	if tcID == "" {
		return "UNKNOWN_TC"
	}
	var b strings.Builder
	for _, r := range tcID {
		switch {
		case r >= 'a' && r <= 'z':
			b.WriteRune(r)
		case r >= 'A' && r <= 'Z':
			b.WriteRune(r)
		case r >= '0' && r <= '9':
			b.WriteRune(r)
		case r == '_' || r == '-':
			b.WriteRune(r)
		default:
			b.WriteRune('_')
		}
	}
	return b.String()
}

func projectRoot(t testing.TB) string {
	t.Helper()
	if override := strings.TrimSpace(os.Getenv("CTRLC_PROJECT_ROOT")); override != "" {
		return override
	}

	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("RecordBEAPIResponse: get working directory: %v", err)
	}

	moduleRoot := findModuleRoot(wd)
	if filepath.Base(moduleRoot) == "backend" {
		return filepath.Dir(moduleRoot)
	}
	return moduleRoot
}

func findModuleRoot(start string) string {
	dir := start
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return start
		}
		dir = parent
	}
}
