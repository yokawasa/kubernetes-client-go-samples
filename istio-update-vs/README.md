
# Usage

Update a VirtualService like this

- Set weight of your desitnation to 100%
- Set the rest of the destination to 0%

```
istio-update-vs -s <your-name> -n <your-namespace> -desthost <your-destination-host> -destsubset <your-destination-subset>
```

Let's go through an example senario

You have a VirtualService like this

```yaml
$ kubectl get vs hoge-api -n testns -o yaml
```
> Output
```
apiVersion: networking.istio.io/v1alpha3
kind: VirtualService
metadata:
  name: hoge-api
spec:
  hosts:
  - hoge-api.testns.svc.cluster.local
  http:
  - route:
    - destination:
        host: hoge-api.testns01.svc.cluster.local
        port:
          number: 8000
        subset: hoge-api
      weight: 33
    - destination:
        host: hoge-api.testns02.svc.cluster.local
        port:
          number: 8000
        subset: hoge-api
      weight: 33
    - destination:
        host: hoge-api.testns03.svc.cluster.local
        port:
          number: 8000
        subset: hoge-api
      weight: 34
    retries:
      attempts: 3
      perTryTimeout: 2s
      retryOn: 5xx,connect-failure
    timeout: 3s
```

Then, you update the VirtualService using the tool
```
istio-update-vs -s hoge-api -n testns -desthost hoge-api.testns01.svc.cluster.local -destsubset hoge-api
```
Let's check the updated VirtualService manifest
```yaml
$ kubectl get vs hoge-api -n testns -o yaml
```
> Output
```
apiVersion: networking.istio.io/v1alpha3
kind: VirtualService
metadata:
  name: hoge-api
spec:
  hosts:
  - hoge-api.testns.svc.cluster.local
  http:
  - route:
    - destination:
        host: hoge-api.testns01.svc.cluster.local
        port:
          number: 8000
        subset: hoge-api
      weight: 100
    - destination:
        host: hoge-api.testns02.svc.cluster.local
        port:
          number: 8000
        subset: hoge-api
    - destination:
        host: hoge-api.testns03.svc.cluster.local
        port:
          number: 8000
        subset: hoge-api
    retries:
      attempts: 3
      perTryTimeout: 2s
      retryOn: 5xx,connect-failure
    timeout: 3s
```



# References
- https://istio.io/latest/blog/2019/announcing-istio-client-go/
- https://github.com/istio/client-go/blob/release-1.9/cmd/example/client.go
- https://github.com/Michael754267513/k8s-client-go/blob/k8s-note/istio/VirtualServices/main.go
- https://pkg.go.dev/istio.io/client-go/pkg/clientset/versioned
