package logger

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetup(t *testing.T) {
	type args struct {
		logFormat string
		logLevel  string
		logFile   string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			"JSON, debug loglevel, stdout output",
			args{"json", "debug", "stdout"},
		},
		{
			"JSON, warning loglevel, file and stdout output",
			args{"json", "warning", "testdata/logtest.log"},
		},
		{
			"Text, info loglevel, discard output",
			args{"text", "info", "null"},
		},
		{
			"Unset logformat, error loglevel, invalid output",
			args{logLevel: "error"},
		},
		{
			"Unset logformat, unset loglevel, unset output",
			args{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Setup(tt.args.logFormat, tt.args.logLevel, tt.args.logFile)
			assert.NoError(t, err)

			if _, err := os.Stat(tt.args.logFile); err == nil {
				if err = os.Remove(tt.args.logFile); err != nil {
					fmt.Println("unable to remove temporary logfile")
				}
			}
		})
	}
}
