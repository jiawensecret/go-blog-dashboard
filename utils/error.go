package utils

import (
	"fmt"
	"math"
	"platform/global"
	"runtime"
)

func ErrorException() {
	if err := recover(); err != nil {
		trace := make([]byte, 1<<16)
		n := runtime.Stack(trace, true)
		s := fmt.Sprintf("panic: '%v'\n, Stack Trace:\n %s", err, string(trace[:int(math.Min(float64(n), float64(7000)))]))
		global.Log.Error(s)
	}
}
