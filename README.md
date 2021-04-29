# kubernetes-client-go-samples
A collection of kubernetes client-go samples

## Sample Lists

- [List Services](list-svc)
- [List Nodes](list-nodes)
- [List Pods](list-pods)
- [List Pods in Service](pods-in-svc)
- [Get Istio Virtual Services](istio-get-vs)


## Develop modules
### Build from source
```
cd $GOPATH/src/github.com/
mkdir -p yokawasa
cd yokawasa/
git clone git@github.com:yokawasa/kubernetes-client-go-samples.git
cd k8s-client-sandbox
make
```

### Add new module

```
git clone git@github.com:yokawasa/kubernetes-client-go-samples.git
cd k8s-client-sandbox
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
  list-services \
  list-nodes \
  list-pods \
  ...snip..
  <new-module> \
```
Then, the module can be built with the Makefile like this
```
make

...snip...
cd $GOPATH/src/github.com/yokawasa/kubernetes-client-go-samples/<new-module> && GO111MODULE=on go build -o $GOPATH/src/github.com/yokawasa/kubernetes-client-go-samples/dist/<new-module> main.go
```
