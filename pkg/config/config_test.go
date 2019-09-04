package config

import (
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoad(t *testing.T) {
	tests := []struct {
		name    string
		want    Config
		wantErr bool
	}{
		{
			name: "Load config",
			want: Config{
				KubeConf: "",
				Registries: []DockerRegistry{
					DockerRegistry{
						Name:   "testregistry",
						Config: "{\n  \"auths\": {\n    \"test.registry.com:80\": {\n      \"auth\": \"lf8pc3RvZmZlci5qb2hhbnNzb246MzU3N2FsYmlu\"\n    }\n  },\n  \"HttpHeaders\": {\n    \"User-Agent\": \"Docker-Client/19.03.1-ce (linux)\"\n  }\n}", Exclude: []string{"kube-public", "kube-system"}, Include: []string(nil)}},
				LogFormat: "text", LogLevel: "info", LogFile: "stdout",
			},
			wantErr: false,
		},
	}

	viper.AddConfigPath("testdata/")

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Load()
			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
