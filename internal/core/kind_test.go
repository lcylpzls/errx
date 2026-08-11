package core

import (
	"testing"
)

func TestKindString(t *testing.T) {
	cases := map[Kind]string{
		KindUnknown:          "unknown",
		KindInvalid:          "invalid_argument",
		KindNotFound:         "not_found",
		KindAlreadyExists:    "already_exists",
		KindUnauthorized:     "unauthorized",
		KindForbidden:        "forbidden",
		KindConflict:         "conflict",
		KindCancelled:        "cancelled",
		KindDeadlineExceeded: "deadline_exceeded",
		KindTimeout:          "timeout",
		KindRateLimited:      "rate_limited",
		KindQuotaExceeded:    "quota_exceeded",
		KindUnavailable:      "unavailable",
		KindInternal:         "internal",
		KindNotImplemented:   "not_implemented",
		KindDataLoss:         "data_loss",
		KindBusiness:         "business",
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

func TestKindCategory(t *testing.T) {
	cases := map[Kind]Category{
		KindInvalid:          CatInput,
		KindUnauthorized:     CatAuth,
		KindForbidden:        CatAuth,
		KindNotFound:         CatState,
		KindAlreadyExists:    CatState,
		KindConflict:         CatState,
		KindCancelled:        CatState,
		KindDeadlineExceeded: CatDependency,
		KindTimeout:          CatDependency,
		KindRateLimited:      CatDependency,
		KindQuotaExceeded:    CatDependency,
		KindUnavailable:      CatDependency,
		KindInternal:         CatSystem,
		KindNotImplemented:   CatSystem,
		KindDataLoss:         CatSystem,
		KindBusiness:         CatBusiness,
		KindUnknown:          CatSystem,
		Kind(255):            CatSystem,
	}
	for k, want := range cases {
		if got := k.Category(); got != want {
			t.Errorf("Kind(%d) Category 不符：got %v, want %v", k, got, want)
		}
	}
}

func TestCategoryString(t *testing.T) {
	if got := CatInput.String(); got != "输入与参数" {
		t.Errorf("CatInput String 不符：%s", got)
	}
	if got := Category(255).String(); got != "未知" {
		t.Errorf("未知 Category 应为 未知：%s", got)
	}
}

func TestKindPolicy(t *testing.T) {
	cases := map[Kind]Policy{
		KindTimeout:          {Retryable: true, UserVisible: true},
		KindRateLimited:      {Retryable: true, UserVisible: true},
		KindQuotaExceeded:    {Retryable: true, UserVisible: true},
		KindUnavailable:      {Retryable: true, Alert: true, UserVisible: true},
		KindInternal:         {Alert: true},
		KindDataLoss:         {Alert: true},
		KindUnknown:          {},
		KindBusiness:         {UserVisible: true},
		KindDeadlineExceeded: {UserVisible: true},
	}
	for k, want := range cases {
		got := k.Policy()
		if got != want {
			t.Errorf("Kind(%d) Policy 不符：got %+v, want %+v", k, got, want)
		}
	}
}

func TestKindRetryable(t *testing.T) {
	retryable := map[Kind]bool{
		KindUnknown:          false,
		KindInvalid:          false,
		KindNotFound:         false,
		KindAlreadyExists:    false,
		KindUnauthorized:     false,
		KindForbidden:        false,
		KindConflict:         false,
		KindCancelled:        false,
		KindDeadlineExceeded: false,
		KindTimeout:          true,
		KindRateLimited:      true,
		KindQuotaExceeded:    true,
		KindUnavailable:      true,
		KindInternal:         false,
		KindNotImplemented:   false,
		KindDataLoss:         false,
		KindBusiness:         false,
		Kind(255):            false,
	}
	for k, want := range retryable {
		if got := k.Retryable(); got != want {
			t.Errorf("Kind(%d) Retryable 不符：got %v, want %v", k, got, want)
		}
	}
}
