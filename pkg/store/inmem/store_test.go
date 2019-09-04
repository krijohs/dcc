package inmem

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewStore(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "Initialize a new Store",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.NotNil(t, NewStore())
		})
	}
}

func TestStore_Operations(t *testing.T) {
	type args struct {
		key   string
		value []byte
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Insert, Read, Delete store operations",
			args: args{
				key:   "default",
				value: []byte("value"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewStore()
			m.Insert(tt.args.key, tt.args.value)

			got, err := m.Read(tt.args.key)
			require.NoError(t, err)

			assert.Equal(t, tt.args.value, got)

			m.Delete(tt.args.key)

			_, err = m.Read(tt.args.key)
			assert.Error(t, err)
		})
	}
}
