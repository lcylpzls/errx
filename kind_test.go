package errx

import (
	"testing"
)

func TestKindString(t *testing.T) {
	cases := map[Kind]string{
		KindUnknown:      "unknown",
		KindInvalid:      "invalid_argument",
		KindNotFound:     "not_found",
		KindUnauthorized: "unauthorized",
		KindConflict:     "conflict",
		KindTimeout:      "timeout",
		KindRateLimited:  "rate_limited",
		KindUnavailable:  "unavailable",
		KindInternal:     "internal",
		KindBusiness:     "business",
	}
	for k, want := range cases {
		if got := k.String(); got != want {
			t.Errorf("Kind(%d) String 不符：got %s, want %s", k, got, want)
		}
	}
	if got := Kind(255).String(); got != "unknown" {
		t.Errorf("未知 Kind 应为 unknown：%s", got)
	}
}

func TestKindRetryable(t *testing.T) {
	retryable := map[Kind]bool{
		KindUnknown:      false,
		KindInvalid:      false,
		KindNotFound:     false,
		KindUnauthorized: false,
		KindConflict:     false,
		KindTimeout:      true,
		KindRateLimited:  true,
		KindUnavailable:  true,
		KindInternal:     false,
		KindBusiness:     false,
		Kind(255):        false,
	}
	for k, want := range retryable {
		if got := k.Retryable(); got != want {
			t.Errorf("Kind(%d) Retryable 不符：got %v, want %v", k, got, want)
		}
	}
}
