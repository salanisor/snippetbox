####  Test builds on OpenShift v4.x

###### Requirements:
* Access to an OpenShift cluster with the ability to create new namespaces.
* Internet access to github.com for git-pulls.

```
.
├── build
│   └── Dockerfile
├── cmd
├── deploy
│   ├── 000-dummy-app-buildconfig.yaml
│   └── 001-pod-deployment.yaml
├── go.mod
├── go.sum
├── main.go
├── README.md
└── vendor
```

###### How-to deploy:

1) Authenticate to an OpenShift cluster
2) `oc apply -f deploy/000-dummy-app-buildconfig.yaml; oc get po -w -n arch-dummy` - wait for the build pod to complete.
3) `oc apply -f deploy/001-pod-deployment.yaml; oc get po -w  -n arch-dummy` - wait for the cpu-dummy pod to run.

Cloned from [Arango_Gutierrez](https://github.com/ArangoGutierrez/Arch-Dummy/)
