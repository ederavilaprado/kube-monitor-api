# Kube Monitor API

API to help with some useful informations about kubernetes cluster and everything around it.

## Quick Start

For this example, first you need docker up and running.
```
# building container...
$ docker build -t kube-monitor:v1 .
# running container... (change env vars bellow first)
$ docker run -d -p 5000:5000 -e 'PASSWORD=pass' -e 'PORT=5000' -e 'K8S_HOST=https://kubernetes.local.com' -e 'K8S_USERNAME=admin' -e 'K8S_PASSWORD=mystrangepassword' -e 'K8S_INSECURE=false' kube-monitor:v1
# testing...
$ curl -X GET -H "Authorization: Basic bGVyb3k6dHNHSUVYYU9UYnlROHR5dg==" -H "Cache-Control: no-cache" "http://localhost:5000"
```

## Config

This env vars bellow will be used to configure the API

Kubernetes configs...
>
- K8S_HOST
- K8S_USERNAME
- K8S_PASSWORD
- K8S_INSECURE

Api configs...
>
- PASSWORD: password for the user leroy (BasicAuth)
- PORT: the api will be launch on this port
