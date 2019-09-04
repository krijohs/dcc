package k8sclient

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
)

func TestNew(t *testing.T) {
	type args struct {
		kubeConfPath string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Initialize a new Client",
			args:    args{kubeConfPath: "testdata/config.yaml"},
			wantErr: false,
		},
		{
			name:    "Initialize a new Client without config file returns an error",
			args:    args{kubeConfPath: ""},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.args.kubeConfPath)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.NotNil(t, got)
		})
	}
}
