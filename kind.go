package errx

import "github.com/lcylpzls/errx/internal/core"

type Kind = core.Kind
type Category = core.Category
type Policy = core.Policy

const (
	KindUnknown          = core.KindUnknown
	KindInvalid          = core.KindInvalid
	KindNotFound         = core.KindNotFound
	KindAlreadyExists    = core.KindAlreadyExists
	KindUnauthorized     = core.KindUnauthorized
	KindForbidden        = core.KindForbidden
	KindConflict         = core.KindConflict
	KindCancelled        = core.KindCancelled
	KindDeadlineExceeded = core.KindDeadlineExceeded
	KindTimeout          = core.KindTimeout
	KindRateLimited      = core.KindRateLimited
	KindQuotaExceeded    = core.KindQuotaExceeded
	KindUnavailable      = core.KindUnavailable
	KindInternal         = core.KindInternal
	KindNotImplemented   = core.KindNotImplemented
	KindDataLoss         = core.KindDataLoss
	KindBusiness         = core.KindBusiness
)

const (
	CatInput      = core.CatInput
	CatAuth       = core.CatAuth
	CatState      = core.CatState
	CatDependency = core.CatDependency
	CatSystem     = core.CatSystem
	CatBusiness   = core.CatBusiness
)
