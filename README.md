
export KUBECONFIG=~/.kube/kubeconfig-k3s 

kubectl apply -f daemonset.yml

kubectl apply -f ingress.yml

kubectl apply -f service.yml

kubectl taint node master0  node-role.kubernetes.io/master=true:NoSchedule