kubectl delete serviceaccount collocation-scheduler --namespace=kube-system
kubectl delete clusterrolebinding collocation-scheduler-as-kube-scheduler
kubectl delete clusterrolebinding collocation-scheduler-as-volume-scheduler
kubectl delete rolebinding collocation-scheduler-as-kube-scheduler --namespace=kube-system
kubectl delete deployment collocation-scheduler --namespace=kube-system