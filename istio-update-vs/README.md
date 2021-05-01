
# Usage

Update a VirtualService like this

- Set weight of your desitnation to 100%
- Set the rest of the destination to 0%

```
istio-update-vs -name <your-name> -namespace <your-namespace> -destination <your-destination>
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
istio-update-vs -name hoge-api -namespace testns -destination hoge-api.testns01.svc.cluster.local
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
