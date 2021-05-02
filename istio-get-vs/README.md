
# Usage

Get a VirtualService named `your-name` in namespace of `your-namespace`
```
istio-get-vs -s <your-name> -n <your-namespace>
```

List VirtualServices  in namespace of `your-namespace`
```
istio-get-vs -n <your-namespace>
```

# References
- https://istio.io/latest/blog/2019/announcing-istio-client-go/
- https://github.com/istio/client-go/blob/release-1.9/cmd/example/client.go
- https://pkg.go.dev/istio.io/client-go/pkg/clientset/versioned
