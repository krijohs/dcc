# Docker Config Controller for Kubernetes

[![Build Status](https://travis-ci.org/krijohs/dcc.svg?branch=master)](https://travis-ci.org/krijohs/dcc)
[![codecov](https://codecov.io/gh/krijohs/dcc/branch/master/graph/badge.svg)](https://codecov.io/gh/krijohs/dcc)

This is a controller for Kubernetes which creates a Docker config secret in namespaces which later can be used for authentication when using private Docker repositories.
The secret will be created in all namespaces except excluded ones, and the controller will listen on namespace state changes and automatically create the secret on a newly added namespace.

## Usage

The application can be deployed using the Helm chart located in 'deployments/charts/dockerconfig-controller', use the values.yaml file to add configuration for your private Docker registries.

The config key in values.yaml file is required and should have the following structure.
```yaml
kubeconf: "" # empty when run inside kubernetes
registries:
  - name: privateregistry
    exclude:
      - kube-public
      - kube-system
    config: |
      {
        "auths": {
          "private.docker.reg": {
            "auth": "xxx"
          }
        },
        "HttpHeaders": {
          "User-Agent": "Docker-Client/19.03.1-ce (linux)"
        }
      }
```
