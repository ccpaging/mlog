// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See License.txt for license information.

package log

import (
	"fmt"
	"os"
)

// defaultLog manually encodes the log to STDERR, providing a basic, default logging implementation
// before elog is fully configured.
func defaultLog(level string, v ...any) {
	fmt.Fprint(os.Stderr, level+fmt.Sprint(v...)+"\n")
}

func defaultDebugLog(v ...any) {
	if DEBUG {
		defaultLog("DEBG ", v...)
	}
}

func defaultTraceLog(v ...any) {
	defaultLog("TRAC ", v...)
}

func defaultInfoLog(v ...any) {
	defaultLog("INFO ", v...)
}

func defaultWarnLog(v ...any) {
	defaultLog("WARN ", v...)
}

func defaultErrorLog(v ...any) {
	defaultLog("EROR ", v...)
}

func defaultFatalLog(v ...any) {
	defaultLog("FATAL ", v...)
}
