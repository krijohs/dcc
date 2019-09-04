package k8sclient

import (
	"github.com/pkg/errors"
	"k8s.io/client-go/kubernetes"

	// Kubernetes authentication modules for handling auth against for example azure or gcp.
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/client-go/tools/clientcmd"
)

// New initializes a new kubernetes ClientSet using a kubeconf or incluster config.
func New(kubeConfPath string) (*kubernetes.Clientset, error) {
	kubeconf, err := clientcmd.BuildConfigFromFlags("", kubeConfPath)
	if err != nil {
		return nil, errors.Wrap(err, "unable to build kubeconfig")
	}

	client, err := kubernetes.NewForConfig(kubeconf)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create clientset")
	}

	return client, nil
}
