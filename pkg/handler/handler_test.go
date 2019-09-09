package handler

import (
	"testing"

	"github.com/krijohs/dcc/pkg/config"
	"github.com/krijohs/dcc/pkg/logger"
	"github.com/krijohs/dcc/pkg/store/inmem"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
)

func TestNew(t *testing.T) {
	type args struct {
		log    logger.Logger
		config Config
		client kubernetes.Interface
		store  Storer
	}
	tests := []struct {
		name string
		args args
		want *Handler
	}{
		{
			name: "Initialize a new Handler",
			args: args{},
			want: &Handler{config: Config{GracePeriod: 0}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := New(tt.args.log, tt.args.config, tt.args.client, tt.args.store)
			assert.NotNil(t, got)
			assert.True(t, got.config.GracePeriod == 5)
		})
	}
}

func TestHandler_HandleItem(t *testing.T) {
	type fields struct {
		log    logger.Logger
		config Config
		client kubernetes.Interface
		store  Storer
	}
	type args struct {
		action int
		ns     string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Handle add an item",
			fields: fields{
				log: logrus.NewEntry(&logrus.Logger{}),
				config: Config{
					GracePeriod: 5,
					Registries:  []config.DockerRegistry{{Name: "test", Config: "config"}},
				},
				client: fake.NewSimpleClientset(),
				store:  inmem.NewStore(),
			},
			args:    args{action: add, ns: "default"},
			wantErr: false,
		},
		{
			name: "Handle update an item",
			fields: fields{
				log: logrus.NewEntry(&logrus.Logger{}),
				config: Config{
					GracePeriod: 5,
					Registries:  []config.DockerRegistry{{Name: "test", Config: "config"}},
				},
				client: fake.NewSimpleClientset(),
				store:  inmem.NewStore(),
			},
			args:    args{action: update, ns: "default"},
			wantErr: true,
		},
		{
			name: "Handle delete an item",
			fields: fields{
				log: logrus.NewEntry(&logrus.Logger{}),
				config: Config{
					GracePeriod: 5,
					Registries:  []config.DockerRegistry{{Name: "test", Config: "config"}},
				},
				client: fake.NewSimpleClientset(),
				store:  inmem.NewStore(),
			},
			args:    args{action: remove, ns: "default"},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{
				log:    tt.fields.log,
				config: tt.fields.config,
				client: tt.fields.client,
				store:  tt.fields.store,
			}
			err := h.HandleItem(tt.args.action, tt.args.ns)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
		})
	}
}

func TestHandler_handleAdd(t *testing.T) {
	type fields struct {
		log    logger.Logger
		config Config
		client kubernetes.Interface
		store  Storer
	}
	type args struct {
		ns         string
		registries []config.DockerRegistry
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{
				log:    tt.fields.log,
				config: tt.fields.config,
				client: tt.fields.client,
				store:  tt.fields.store,
			}
			if err := h.handleAdd(tt.args.ns, tt.args.registries); (err != nil) != tt.wantErr {
				t.Errorf("Handler.handleAdd() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_shouldHandle(t *testing.T) {
	type args struct {
		ns  string
		reg config.DockerRegistry
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := shouldHandle(tt.args.ns, tt.args.reg); got != tt.want {
				t.Errorf("shouldHandle() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHandler_handleUpdate(t *testing.T) {
	type fields struct {
		log    logger.Logger
		config Config
		client kubernetes.Interface
		store  Storer
	}
	type args struct {
		ns         string
		registries []config.DockerRegistry
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{
				log:    tt.fields.log,
				config: tt.fields.config,
				client: tt.fields.client,
				store:  tt.fields.store,
			}
			if err := h.handleUpdate(tt.args.ns, tt.args.registries); (err != nil) != tt.wantErr {
				t.Errorf("Handler.handleUpdate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHandler_handleRemove(t *testing.T) {
	type fields struct {
		log    logger.Logger
		config Config
		client kubernetes.Interface
		store  Storer
	}
	type args struct {
		ns         string
		registries []config.DockerRegistry
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{
				log:    tt.fields.log,
				config: tt.fields.config,
				client: tt.fields.client,
				store:  tt.fields.store,
			}
			h.handleRemove(tt.args.ns, tt.args.registries)
		})
	}
}
