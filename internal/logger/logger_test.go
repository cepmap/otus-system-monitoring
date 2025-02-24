package logger

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLogger(t *testing.T) {
	type testCase struct {
		loggerLevel  string
		messageLevel string
		shouldPrint  bool
	}
	testCases := []testCase{
		// Debug
		{
			loggerLevel:  DebugLevel,
			messageLevel: DebugLevel,
			shouldPrint:  true,
		},
		{
			loggerLevel:  DebugLevel,
			messageLevel: InfoLevel,
			shouldPrint:  true,
		},
		{
			loggerLevel:  DebugLevel,
			messageLevel: WarnLevel,
			shouldPrint:  true,
		},
		{
			loggerLevel:  DebugLevel,
			messageLevel: ErrorLevel,
			shouldPrint:  true,
		},
		// Info
		{
			loggerLevel:  InfoLevel,
			messageLevel: DebugLevel,
			shouldPrint:  false,
		},
		{
			loggerLevel:  InfoLevel,
			messageLevel: InfoLevel,
			shouldPrint:  true,
		},
		{
			loggerLevel:  InfoLevel,
			messageLevel: WarnLevel,
			shouldPrint:  true,
		},
		{
			loggerLevel:  InfoLevel,
			messageLevel: ErrorLevel,
			shouldPrint:  true,
		},
		// Warning
		{
			loggerLevel:  WarnLevel,
			messageLevel: DebugLevel,
			shouldPrint:  false,
		},
		{
			loggerLevel:  WarnLevel,
			messageLevel: InfoLevel,
			shouldPrint:  false,
		},
		{
			loggerLevel:  WarnLevel,
			messageLevel: WarnLevel,
			shouldPrint:  true,
		},
		{
			loggerLevel:  WarnLevel,
			messageLevel: ErrorLevel,
			shouldPrint:  true,
		},
		// Error
		{
			loggerLevel:  ErrorLevel,
			messageLevel: DebugLevel,
			shouldPrint:  false,
		},
		{
			loggerLevel:  ErrorLevel,
			messageLevel: InfoLevel,
			shouldPrint:  false,
		},
		{
			loggerLevel:  ErrorLevel,
			messageLevel: WarnLevel,
			shouldPrint:  false,
		},
		{
			loggerLevel:  WarnLevel,
			messageLevel: ErrorLevel,
			shouldPrint:  true,
		},
	}
	for _, tc := range testCases {
		var msg string
		if tc.shouldPrint {
			msg = fmt.Sprintf("%s logger should print %s message", tc.loggerLevel, tc.messageLevel)
		} else {
			msg = fmt.Sprintf("%s logger should not print %s message", tc.loggerLevel, tc.messageLevel)
		}
		t.Run(msg, func(t *testing.T) {
			w := &bytes.Buffer{}

			SetLogLevel(tc.loggerLevel)
			SetWriter(w)

			switch tc.messageLevel {
			case DebugLevel:
				Debug(msg)
			case InfoLevel:
				Info(msg)
			case WarnLevel:
				Warn(msg)
			case ErrorLevel:
				Error(msg)
			default:
				require.Failf(t, "undefined message level: %s", tc.messageLevel)
			}

			fmt.Println(w.String(), msg)

			if tc.shouldPrint {
				require.Contains(t, w.String(), msg)
			} else {
				require.NotContains(t, w.String(), msg)
			}
		})
	}
}
