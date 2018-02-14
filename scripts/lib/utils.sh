# copy fetched kubeconfig to the default location
set_kubeconfig() {
  mkdir -p ~/.kube
  cp gcs-kubeconfig/config ~/.kube/config
}
