# Docker Config Controller for Kubernetes

[![Build Status](https://travis-ci.org/krijohs/dcc.svg?branch=master)](https://travis-ci.org/krijohs/dcc)
[![codecov](https://codecov.io/gh/krijohs/dcc/branch/master/graph/badge.svg)](https://codecov.io/gh/krijohs/dcc)

This is an controller for Kubernetes which creates a Docker config secret to use for authentication for private Docker repositories.
The secret will be created in all namespaces except excluded ones, and the controller will listen on namespace state changes and automatically create the secret on a newly added namespace.
The controller can be used to just replicate secrets over namespaces and it is not bound to just to be used with Docker registries even though it was the intention when creating it.

## Usage

A config.yaml is required which should have the following structure.
```yaml
kubeconf: "/path/to/kubeconf"
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
If the kubeconf is empty, the application assumes that it's run inside Kubernetes. Multiple registries can be defined by creating a new item below registries key.

This file can be placed in the same directory the binary file is located, in $HOME/.dockerconfig or in /etc/dockerconfig.

## TODO
Example Kubernetes deployment