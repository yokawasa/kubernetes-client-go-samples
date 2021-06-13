# kubernetes-client-go-samples
A collection of Kubernetes and Istio client-go sample code

- [kubernetes/client-go](https://github.com/kubernetes/client-go)
- [istio/client-go](https://github.com/istio/client-go)

## Sample Lists

- [List Services](svclist)
- [List Nodes](nodelist)
- [List Pods](podlist)
- [List Pods in Service](podlist-in-svc)
- [Infomrer](informer)
- [Get Istio VirtualServices](istio-get-vs)
- [Update Istio VirtualService](istio-update-vs)

## Quickstart

### Local build and run

To Locally build binary, run the following command. Compbiled binaries are created under dist directory
```bash
git clone git@github.com:yokawasa/kubernetes-client-go-samples.git
cd kubernetes-client-go-samples
make
```

Then, run the compiled binaries
```bash
./dist/svclist
./dist/nodelist
./dist/podlist
./dist/podlist-in-svc
./dist/istio-get-vs -s <virtual service name> -n <namespace>
./dist/istio-update-vs -s <virtual service name> -n <namespace> -desthost <destination host> -destsubset <destination subset>
```

### Docker

Docker build with the Makefile like this
```bash
git clone git@github.com:yokawasa/kubernetes-client-go-samples.git
cd kubernetes-client-go-samples
# docker build using make
make docker-build
```
Or you can docker build using docker command
```bash
# docker build using docker command
export docker_image_repo=kubernetes-client-go-samples
export docker_image_tag=latest
docker build -t ${docker_image_repo}:${docker_image_tag} . --target executor
```

Then, run the commands like this
```bash
export docker_image_repo=kubernetes-client-go-samples
export docker_image_tag=latest

docker run --rm -it -v ${HOME}/.kube/config:/root/.kube/config ${docker_image_repo}:${docker_image_tag} sh -c "/svclist"
docker run --rm -it -v ${HOME}/.kube/config:/root/.kube/config ${docker_image_repo}:${docker_image_tag} sh -c "/nodelist"
docker run --rm -it -v ${HOME}/.kube/config:/root/.kube/config ${docker_image_repo}:${docker_image_tag} sh -c "/podlist"
docker run --rm -it -v ${HOME}/.kube/config:/root/.kube/config ${docker_image_repo}:${docker_image_tag} sh -c "/podlist-in-svc"
docker run --rm -it -v ${HOME}/.kube/config:/root/.kube/config ${docker_image_repo}:${docker_image_tag} sh -c "/informer"
docker run --rm -it -v ${HOME}/.kube/config:/root/.kube/config ${docker_image_repo}:${docker_image_tag} sh -c "/istio-get-vs -s <virtual servide name> -n <namespace>"
docker run --rm -it -v ${HOME}/.kube/config:/root/.kube/config ${docker_image_repo}:${docker_image_tag} sh -c "/istio-update-vs -s <virtual service name> -n <namespace> -desthost <destination host> -destsubset <destination subset>"

# You might need to volume mount your .aws dir if you are accessing to AWS EKS Kubernets cluster 
docker run --rm -it -v ${HOME}/.kube/config:/root/.kube/config -v ${HOME}/.aws:/root/.aws ${docker_image_repo}:${docker_image_tag} sh -c "/svclist"
docker run --rm -it -v ${HOME}/.kube/config:/root/.kube/config -v ${HOME}/.aws:/root/.aws ${docker_image_repo}:${docker_image_tag} sh -c "/nodelist"
docker run --rm -it -v ${HOME}/.kube/config:/root/.kube/config -v ${HOME}/.aws:/root/.aws ${docker_image_repo}:${docker_image_tag} sh -c "/podlist"
docker run --rm -it -v ${HOME}/.kube/config:/root/.kube/config -v ${HOME}/.aws:/root/.aws ${docker_image_repo}:${docker_image_tag} sh -c "/podlist-in-svc"
docker run --rm -it -v ${HOME}/.kube/config:/root/.kube/config -v ${HOME}/.aws:/root/.aws ${docker_image_repo}:${docker_image_tag} sh -c "/informer"
docker run --rm -it -v ${HOME}/.kube/config:/root/.kube/config -v ${HOME}/.aws:/root/.aws ${docker_image_repo}:${docker_image_tag} sh -c "/istio-get-vs -s <virtual servide name> -n <namespace>"
docker run --rm -it -v ${HOME}/.kube/config:/root/.kube/config -v ${HOME}/.aws:/root/.aws ${docker_image_repo}:${docker_image_tag} sh -c "/istio-update-vs -s <virtual service name> -n <namespace> -desthost <destination host> -destsubset <destination subset>"
```

## Develop modules
### Build from source
```
cd $GOPATH/src/github.com/
mkdir -p yokawasa
cd yokawasa/
git clone git@github.com:yokawasa/kubernetes-client-go-samples.git
cd kubernetes-client-go-samples
make
```

### Add new module

```
git clone git@github.com:yokawasa/kubernetes-client-go-samples.git
cd kubernetes-client-go-samples
mkdir new-module
vi main.go
go mod init new-module
go build
```

Once you confirm the module build, add the name of module to `Makefile` on the project root

```
.PHONY: clean all list-services list-nodes list-pods ... <new-module>
...snip...
TARGETS=\
  svclist \
  nodelist \
  podlist \
  ...snip..
  <new-module> \

```

Then, the module can be built with the Makefile like this


```
make

...snip...
cd $GOPATH/src/github.com/yokawasa/kubernetes-client-go-samples/<new-module> && GO111MODULE=on go build -o $GOPATH/src/github.com/yokawasa/kubernetes-client-go-samples/dist/<new-module> main.go
```
