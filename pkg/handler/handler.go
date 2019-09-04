package handler

import (
	"bytes"

	"github.com/krijohs/dcc/pkg/config"
	"github.com/krijohs/dcc/pkg/logger"

	"github.com/pkg/errors"
	"k8s.io/client-go/kubernetes"
)

const (
	dockerConfigJSON = ".dockerconfigjson"
)

type Handler struct {
	log    logger.Logger
	config Config
	client kubernetes.Interface
	store  Storer
}

type Config struct {
	GracePeriod int64
	Registries  []config.DockerRegistry
}

type Storer interface {
	Insert(key string, value []byte)
	Delete(key string)
	Read(key string) ([]byte, error)
}

// New initializes a Handler.
func New(log logger.Logger, config Config, client kubernetes.Interface, store Storer) *Handler {
	if config.GracePeriod == 0 {
		config.GracePeriod = 5
	}
	return &Handler{log: log, config: config, client: client, store: store}
}

const (
	add = iota
	update
	remove
)

// HandleItem adds, updates or deletes a secret in the passed in namespace.
func (h *Handler) HandleItem(action int, ns string) error {
	h.log.Debugf("handling item for namespace %s", ns)

	switch action {
	case add:
		if err := h.handleAdd(ns, h.config.Registries); err != nil {
			return errors.Wrapf(err, "unable to handle add of registries in namespace %s", ns)
		}

	case update:
		if err := h.handleUpdate(ns, h.config.Registries); err != nil {
			return errors.Wrapf(err, "unable to handle update of registries in namespace %s", ns)
		}

	case remove:
		h.handleRemove(ns, h.config.Registries)

	default:
		return errors.Errorf("invalid action %d passed in as argument", action)
	}

	return nil
}

func (h *Handler) handleAdd(ns string, registries []config.DockerRegistry) error {
	for _, reg := range registries {
		if !shouldHandle(ns, reg) {
			h.log.Infof("skipping namespace %s since it is filtered out", ns)
			return nil
		}

		registryCfg := []byte(reg.Config)
		// check if namespace already has been processed, eg its already in the store
		// also check if the stored secret data is the same as the new data
		storedSecret, err := h.store.Read(ns)
		if err == nil {
			if bytes.Equal(storedSecret, registryCfg) {
				h.log.Debugf("namespace %s and docker registry secret %s already exists in store", ns, reg.Name)
				return nil
			}
		}

		exists, err := h.secretExists(ns, reg.Name, registryCfg)
		if err != nil {
			return errors.Wrapf(err, "unable to check if docker registry secret %s exists", reg.Name)
		}
		if exists {
			h.log.Debugf("secret %s already exists in namespace %s", reg.Name, ns)
			return nil
		}

		if _, err = h.newSecret(ns, reg.Name, registryCfg); err != nil {
			return errors.Wrapf(err, "unable to create docker registry secret %s in namespace", reg.Name)
		}

		h.store.Insert(ns, registryCfg)
		h.log.Infof("create docker registry secret %s in namespace %s", reg.Name, ns)
	}

	return nil
}

func shouldHandle(ns string, reg config.DockerRegistry) bool {
	wildcard := "*"

	for _, incNs := range reg.Include {
		if incNs == ns {
			return true
		}
	}

	for _, excNs := range reg.Exclude {
		if excNs == ns || excNs == wildcard {
			return false
		}
	}

	return true
}

func (h *Handler) handleUpdate(ns string, registries []config.DockerRegistry) error {
	for _, reg := range registries {
		if !shouldHandle(ns, reg) {
			h.log.Infof("skipping namespace %s since it is filtered out", ns)
			return nil
		}

		registryCfg := []byte(reg.Config)
		exists, err := h.secretExists(ns, reg.Name, registryCfg)
		if err != nil {
			return errors.Wrap(err, "unable to check if docker registry secret exists")
		}
		if exists {
			return nil
		}

		if _, err := h.updateSecret(ns, reg.Name, registryCfg); err != nil {
			return errors.Wrap(err, "unable to create docker registry secret in namespace")
		}

		h.store.Insert(ns, registryCfg)
		h.log.Infof("updated docker registry secret %s in namespace %s", reg.Name, ns)
	}

	return nil
}

func (h *Handler) handleRemove(ns string, registries []config.DockerRegistry) {
	for _, reg := range registries {
		if !shouldHandle(ns, reg) {
			h.log.Infof("skipping namespace %s since it is filtered out", ns)
		}

		h.store.Delete(ns)
		h.log.Infof("deleted docker registry secret %s in namespace %s", reg.Name, ns)
	}
}
