package debug

import "runtime"

// GetStacktrace returns formatted string of stacktrace
func GetStacktrace(isForAll bool) string {
	outBuffer := make([]byte, 1024)

	for {
		n := runtime.Stack(outBuffer, isForAll)
		if n < len(outBuffer) {
			break
		}
		outBuffer = make([]byte, 2*len(outBuffer))
	}

	return string(outBuffer)
}
