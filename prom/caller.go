package prom

import (
	"path"
	"runtime"
)

// callerFunc returns trimmed path and name of parent function.
func callerFunc(skip ...int) string {
	skipFrames := 2
	if len(skip) == 1 {
		skipFrames = skip[0]
	}

	pc, _, _, ok := runtime.Caller(skipFrames)
	if !ok {
		return ""
	}

	f := runtime.FuncForPC(pc)

	pathName := path.Base(path.Dir(f.Name())) + "/" + path.Base(f.Name())

	return pathName
}
