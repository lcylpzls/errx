package errx

import "github.com/lcylpzls/errx/internal/core"

type MetricsHook = core.MetricsHook

func SetMetricsHook(hook MetricsHook) { core.SetMetricsHook(hook) }
func ResetMetricsHook()               { core.ResetMetricsHook() }
