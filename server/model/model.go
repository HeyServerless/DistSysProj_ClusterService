package model

import (
	// "encoding/json"

	// appsv1 "k8s.io/api/core/v1"
	// metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	corev1 "k8s.io/api/core/v1"
	// "k8s.io/client-go/util/retry"
	appsv1 "k8s.io/api/apps/v1"
	// "k8s.io/apimachinery/pkg/runtime/schema"
	// "k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type KubeDeployment struct {
	// kubernete deployment request model
	Title         string `json:"title"`
	Replicas      int    `json:"replicas"`
	ContainerName string `json:"container_name"`
}

type KubeService struct {
	// kubernete service request model
	Title         string `json:"title"`
	ContainerName string `json:"container_name"`
}

type KubePods struct {
	// kubernete pods request model
	Title         string `json:"title"`
	ContainerName string `json:"container_name"`
}

type KubeIngress struct {
	// kubernete ingress request model
	Title         string `json:"title"`
	ContainerName string `json:"container_name"`
}

type KubeSecret struct {
	// kubernete secret request model
	Title         string `json:"title"`
	ContainerName string `json:"container_name"`
}

type CreateDeploymentRequest struct {
	Name       string
	Replicas   int32
	Selector   map[string]string
	Labels     map[string]string
	Containers []corev1.Container
}

type CreateDeploymentResponse struct {
	Deployment *appsv1.Deployment
}

type UpdateDeploymentRequest struct {
	Name      string
	Namespace string
	Replicas  int32
}

type UpdateDeploymentResponse struct {
	Deployment *appsv1.Deployment
}

type DeploymentList struct {
	Deployments *appsv1.DeploymentList
}

type CreatePodRequest struct {
	Name       string
	Containers []corev1.Container
}

type CreatePodResponse struct {
	Pod *corev1.Pod
}
type DeletePodRequest struct {
	Name string
}

type DeletePodResponse struct {
	Pod *corev1.Pod
}

type GetPodRequest struct {
	Name string
}

type GetPodResponse struct {
	Pod *corev1.Pod
}

type UpdatePodRequest struct {
	Name       string
	Containers []corev1.Container
}

type UpdatePodResponse struct {
	Pod *corev1.Pod
}

type PatchPodRequest struct {
	Name       string
	Containers []corev1.Container
}

type PatchPodResponse struct {
	Pod *corev1.Pod
}

type ReplacePodRequest struct {
	Name       string
	Containers []corev1.Container
}

type ReplacePodResponse struct {
	Pod *corev1.Pod
}

type WatchPodRequest struct {
	Name string
}

type WatchPodResponse struct {
	Pod *corev1.Pod
}

type Expression struct {
	ID   int    `json:"id"`
	Exp  string `json:"exp"`
	Resp float64
}

type CreateServiceRequest struct {
	Name     string
	Selector map[string]string
	Ports    []corev1.ServicePort
	Type     corev1.ServiceType
}

type CreateServiceResponse struct {
	Service *corev1.Service
}

type ServiceList struct {
	Services *corev1.ServiceList
}
