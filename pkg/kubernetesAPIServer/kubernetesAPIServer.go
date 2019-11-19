package kubernetesAPIServer

import (
	"k8s.io/client-go/kubernetes"
)

//This is interface for Kubernetes API Server
type KubernetesAPIServer struct {
	Suffix string
	Client kubernetes.Interface
}

func New() *KubernetesAPIServer {
	return &KubernetesAPIServer{}
}

func (kAS *KubernetesAPIServer) SetClient(client kubernetes.Interface) {
	kAS.Client = client
}
