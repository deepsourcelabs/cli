package version

import (
	"testing"
	"time"

	"github.com/deepsourcelabs/cli/version"
)

func TestOptions_Run(t *testing.T) {
	date, _ := time.Parse("2006-01-02", "2021-01-21")

	getBuildInfo = func() *version.BuildInfo {
		return &version.BuildInfo{
			Version: "1.5.0",
			Date:    date,
		}
	}

	tests := []struct {
		name string
		o    Options
		want string
	}{
		{
			name: "must return the string output for command",
			o:    Options{},
			want: "DeepSource CLI version 1.5.0 (2021-01-21)",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := Options{}
			if got := o.Run(); got != tt.want {
				t.Errorf("Options.Run() = %v, want %v", got, tt.want)
			}
		})
	}
}
