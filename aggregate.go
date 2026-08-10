package errx

import "github.com/lcylpzls/errx/internal/core"

type Aggregate = core.Aggregate

func Join(errs ...error) error { return core.Join(errs...) }
