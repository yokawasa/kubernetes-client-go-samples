
## How to develop informer sample code
```
mkdir informer
cd informer
vi main.go
go mod init informer

# get modules
go get k8s.io/client-go/informers
go get k8s.io/client-go/informers/core/v1
go get k8s.io/client-go/kubernetes
go get k8s.io/client-go/tools/cache
go get k8s.io/client-go/tools/clientcmd
go get k8s.io/klog/v2
go get k8s.io/kubectl/pkg/util/logs
go get k8s.io/api/core/v1@v0.21.1

# build
go build 
```
