package httpx

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/lcylpzls/errx"
)

func TestStatus(t *testing.T) {
	cases := []struct {
		name string
		err  error
		want int
	}{
		{"not found", errx.New(errx.KindNotFound, "N", "n"), http.StatusNotFound},
		{"timeout", errx.New(errx.KindTimeout, "T", "t"), http.StatusGatewayTimeout},
		{"plain", errors.New("plain"), http.StatusInternalServerError},
		{"nil", nil, http.StatusInternalServerError},
		{"aggregate", errx.Join(
			errx.New(errx.KindBusiness, "B", "b"),
			errx.New(errx.KindUnauthorized, "U", "u"),
		), http.StatusUnprocessableEntity},
	}
	for _, tc := range cases {
		if got := Status(tc.err); got != tc.want {
			t.Errorf("%s：Status 不符：got %d, want %d", tc.name, got, tc.want)
		}
	}
}

func TestWriteJSON(t *testing.T) {
	rec := httptest.NewRecorder()
	err := errx.New(errx.KindBusiness, "ORDER_FAIL", "下单失败").WithField("order_id", "1")
	WriteJSON(rec, err)

	if rec.Code != http.StatusUnprocessableEntity {
		t.Errorf("状态码不符：%d", rec.Code)
	}
	if got := rec.Header().Get("Content-Type"); got != "application/json; charset=utf-8" {
		t.Errorf("Content-Type 不符：%s", got)
	}
	var body map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("响应体非法 JSON：%v", err)
	}
	if body["code"] != "ORDER_FAIL" || body["kind"] != "business" || body["message"] != "ORDER_FAIL: 下单失败" {
		t.Errorf("响应体不符：%v", body)
	}
}

func TestWriteJSONPlainError(t *testing.T) {
	rec := httptest.NewRecorder()
	WriteJSON(rec, errors.New("boom"))
	if rec.Code != http.StatusInternalServerError {
		t.Errorf("普通错误状态码不符：%d", rec.Code)
	}
	var body map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("响应体非法 JSON：%v", err)
	}
	if body["code"] != string(errx.CodeUnknown) {
		t.Errorf("普通错误 code 应回退为 UNKNOWN：%v", body)
	}
	if body["kind"] != "unknown" {
		t.Errorf("普通错误 kind 不符：%v", body)
	}
}

func TestWriteJSONNil(t *testing.T) {
	rec := httptest.NewRecorder()
	WriteJSON(rec, nil)
	if rec.Code != http.StatusInternalServerError {
		t.Errorf("nil 错误状态码不符：%d", rec.Code)
	}
	if got := rec.Header().Get("Content-Type"); got != "application/json; charset=utf-8" {
		t.Errorf("Content-Type 不符：%s", got)
	}
	var body map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("响应体非法 JSON：%v", err)
	}
	if body["code"] != string(errx.CodeUnknown) || body["kind"] != "unknown" {
		t.Errorf("nil 错误响应体不符：%v", body)
	}
}
