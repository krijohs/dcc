# Docker Config Controller for Kubernetes

[![Build Status](https://travis-ci.org/krijohs/dcc.svg?branch=master)](https://travis-ci.org/krijohs/dcc)
[![codecov](https://codecov.io/gh/krijohs/dcc/branch/master/graph/badge.svg)](https://codecov.io/gh/krijohs/dcc)

This is an controller for Kubernetes which creates Docker config in all namespaces except excluded ones.
This makes it possible to define a Docker config for a private repository in a config.yaml file which then will be created
in all namespaces and can be used as to pull Docker images from private repositories.

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
If the kubeconf is empty, the application assumes that it's run in Kubernetes. Multiple registries can be defined by creating a new item below registries key.

This file can be placed in the same directory the binary file is located, in $HOME/.dockerconfig or in /etc/dockerconfig.


