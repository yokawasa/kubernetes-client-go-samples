

```
kubens default
kubectl run -it --rm=true busybox --image=yauritux/busybox-curl --restart=Never
./podlist
kubectl delete pod busybox
./podlist
```
