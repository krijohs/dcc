package controller

import (
	"testing"
	"time"

	"github.com/krijohs/dcc/pkg/logger"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
)

func TestNew_and_Controller_initController(t *testing.T) {
	type args struct {
		log     logger.Logger
		config  Config
		client  kubernetes.Interface
		handler Handlerer
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Initialize and call initController",
			args: args{
				log:     logrus.NewEntry(&logrus.Logger{}),
				config:  Config{},
				client:  fake.NewSimpleClientset(),
				handler: &mockHandler{err: nil},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := New(tt.args.log, tt.args.config, tt.args.client, tt.args.handler)
			c.initController()
		})
	}
}

func TestController_Watch(t *testing.T) {
	type controllerArgs struct {
		log     logger.Logger
		config  Config
		client  kubernetes.Interface
		handler Handlerer
	}
	type args struct {
		stopCh chan struct{}
	}
	tests := []struct {
		name           string
		controllerArgs controllerArgs
		args           args
		sleep          time.Duration
	}{
		{
			name: "Watch for namespace state changes",
			controllerArgs: controllerArgs{
				log:     logrus.NewEntry(&logrus.Logger{}),
				config:  Config{},
				client:  fake.NewSimpleClientset(),
				handler: &mockHandler{err: nil},
			},
			args:  args{stopCh: make(chan struct{}, 1)},
			sleep: time.Second * 2,
		},
		{
			name: "Watch for namespace state changes, test processQueue handling of error from mock HandleItem",
			controllerArgs: controllerArgs{
				log:     logrus.NewEntry(&logrus.Logger{}),
				config:  Config{},
				client:  fake.NewSimpleClientset(),
				handler: &mockHandler{err: errors.New("err")},
			},
			args:  args{stopCh: make(chan struct{}, 1)},
			sleep: time.Second * 15,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := New(tt.controllerArgs.log, tt.controllerArgs.config, tt.controllerArgs.client, tt.controllerArgs.handler)
			c.initController()

			tt.controllerArgs.client.CoreV1().Namespaces().Create(&v1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "testns"}})

			go func() {
				time.Sleep(tt.sleep)
				tt.args.stopCh <- struct{}{}
				close(tt.args.stopCh)
			}()

			c.Watch(tt.args.stopCh)
		})
	}
}

type mockHandler struct {
	err error
}

func (mh *mockHandler) HandleItem(action int, ns string) error {
	return mh.err
}
