package logs

import (
	"runtime"
	"strings"
)

type fileInfo struct {
	file     string
	line     int
	funcName string
}

func newFileInfo(callDepth int) *fileInfo {
	// Inspect runtime call stack
	pc := make([]uintptr, callDepth)
	runtime.Callers(callDepth, pc)

	f := runtime.FuncForPC(pc[callDepth-1])
	file, line := f.FileLine(pc[callDepth-1])

	funcName := f.Name()
	if slash := strings.LastIndex(funcName, "."); slash >= 0 {
		funcName = funcName[slash+1:]
	}

	return &fileInfo{
		file:     file,
		line:     line,
		funcName: funcName,
	}
}
