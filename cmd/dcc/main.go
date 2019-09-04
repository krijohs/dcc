package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/krijohs/dcc/pkg/config"
	"github.com/krijohs/dcc/pkg/controller"
	"github.com/krijohs/dcc/pkg/handler"
	"github.com/krijohs/dcc/pkg/k8sclient"
	"github.com/krijohs/dcc/pkg/logger"
	"github.com/krijohs/dcc/pkg/store/inmem"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	loggr, err := logger.Setup(cfg.LogFormat, cfg.LogLevel, cfg.LogFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	log := loggr.WithField("app", "docker-config-controller")

	client, err := k8sclient.New(cfg.KubeConf)
	if err != nil {
		log.Fatal(err)
	}

	inMemStore := inmem.NewStore()

	handlerCfg := handler.Config{Registries: cfg.Registries}
	h := handler.New(log, handlerCfg, client, inMemStore)

	ctrlCfg := controller.Config{
		KubeConf:     cfg.KubeConf,
		ResyncPeriod: time.Second * 5,
	}
	c := controller.New(log, ctrlCfg, client, h)

	stopCh := make(chan struct{})

	go signaler(stopCh)
	go c.Watch(stopCh)

	<-stopCh

	log.Error("terminated")
}

func signaler(stopCh chan struct{}) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	<-c
	stopCh <- struct{}{}
}
