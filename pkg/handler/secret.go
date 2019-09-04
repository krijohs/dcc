package handler

import (
	"bytes"

	"github.com/pkg/errors"
	v1 "k8s.io/api/core/v1"
	k8serr "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (h *Handler) newSecret(namespace, name string, b []byte) (v1.Secret, error) {
	secret, err := h.client.CoreV1().Secrets(namespace).Create(&v1.Secret{
		TypeMeta: metav1.TypeMeta{
			Kind: "Secret",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Data: makeBytesDockerCfg(b),
	})
	if err != nil {
		return v1.Secret{}, errors.Wrapf(err, "unable to create secret in namespace %s", namespace)
	}

	return *secret, nil
}

func (h *Handler) updateSecret(namespace, name string, b []byte) (v1.Secret, error) {
	secret, err := h.client.CoreV1().Secrets(namespace).Update(&v1.Secret{
		TypeMeta: metav1.TypeMeta{
			Kind: "Secret",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Data: makeBytesDockerCfg(b),
	})
	if err != nil {
		return v1.Secret{}, errors.Wrap(err, "unable to update secret in namespace")
	}

	return *secret, nil
}

func (h *Handler) secretExists(namespace, name string, data []byte) (bool, error) {
	secret, err := h.client.CoreV1().Secrets(namespace).Get(name, metav1.GetOptions{})
	if err != nil {
		if kErr, ok := err.(*k8serr.StatusError); ok && kErr.Status().Code == 404 {
			return false, nil
		}
		return false, errors.Wrap(err, "unable to get secret")
	}

	return bytes.Equal(secret.Data[dockerConfigJSON], data), nil
}

func makeBytesDockerCfg(b []byte) map[string][]byte {
	data := make(map[string][]byte)
	data[dockerConfigJSON] = b

	return data
}
