package services

import (
	// "encoding/json"

	// appsv1 "k8s.io/api/core/v1"
	// metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"context"
	"log"
	"strings"
	"time"

	models "github.com/clusterService/server/model"
	"github.com/gin-gonic/gin"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	// "k8s.io/client-go/util/retry"
	// "k8s.io/apimachinery/pkg/runtime/schema"
	// "k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func CreateService(c *gin.Context, createServiceRequest *models.CreateServiceRequest) (*models.CreateServiceResponse, error) {
	log.Println("================================Creating service================================")
	k8Client := c.MustGet("kubeClient").(*kubernetes.Clientset)

	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: createServiceRequest.Name,
		},
		Spec: corev1.ServiceSpec{
			Selector: createServiceRequest.Selector,
			Ports:    createServiceRequest.Ports,
			Type:     createServiceRequest.Type,
		},
	}
	log.Println("service spec: ", service.Spec)
	log.Println("service spec cluster ip: ", service.Spec.ClusterIP)
	log.Println("service spec ports: ", service.Spec.Ports)
	log.Println("service spec selector: ", service.Spec.Selector)
	log.Println("service spec type: ", service.Spec.Type)
	log.Println("service spec status: ", service.Status)

	result, err := k8Client.CoreV1().Services("default").Create(context.Background(), service, metav1.CreateOptions{})
	if err != nil {
		log.Println(err)
		log.Println("Error creating service: %v", err.Error)
		if strings.Contains(err.Error(), "already exists") {
			log.Println("================================Service already exists================================")
			return &models.CreateServiceResponse{
				Service: result,
			}, nil
		} else {
			log.Println("================================Error creating service================================")
			return nil, err
		}

	}

	// get the status of service

	for {
		serviceStatus, err := k8Client.CoreV1().Services("default").Get(context.Background(), createServiceRequest.Name, metav1.GetOptions{})
		if err != nil {
			log.Fatal(err)
		}
		status := serviceStatus.Status
		if len(service.Spec.Ports) > 0 {
			// service has one or more ports
			log.Println("Service has ports")
			for _, port := range service.Spec.Ports {
				log.Println("port: ", port)

				log.Println("port status: ", status)
				return &models.CreateServiceResponse{
					Service: result,
				}, nil

			}

		}
		log.Println("Service is not ready")
		sleepTime := 5
		log.Printf("Sleeping for %d seconds", sleepTime)
		time.Sleep(time.Duration(sleepTime) * time.Second)
	}

	// list end points of service
	endpoints, err := k8Client.CoreV1().Endpoints("default").List(context.Background(), metav1.ListOptions{})
	if err != nil {
		log.Fatal(err)
	}
	log.Println("endpoints: ", endpoints)
	for _, endpoint := range endpoints.Items {
		log.Println("endpoint: ", endpoint)
	}

	return &models.CreateServiceResponse{
		Service: result,
	}, nil

}

func GetAllServices(c *gin.Context) (*models.ServiceList, error) {

	k8Client := c.MustGet("kubeClient").(*kubernetes.Clientset)

	services, err := k8Client.CoreV1().Services("default").List(context.Background(), metav1.ListOptions{})
	if err != nil {
		log.Fatal(err)
	}

	return &models.ServiceList{
		Services: services,
	}, nil

}

func GetService(c *gin.Context, serviceName string) (*corev1.Service, error) {

	k8Client := c.MustGet("kubeClient").(*kubernetes.Clientset)

	service, err := k8Client.CoreV1().Services("default").Get(context.Background(), serviceName, metav1.GetOptions{})
	if err != nil {
		log.Fatal(err)
	}

	return service, nil

}

func DeleteService(c *gin.Context, serviceName string) error {

	k8Client := c.MustGet("kubeClient").(*kubernetes.Clientset)

	err := k8Client.CoreV1().Services("default").Delete(context.Background(), serviceName, metav1.DeleteOptions{})
	if err != nil {
		log.Fatal(err)
	}

	return nil

}

func UpdateService(c *gin.Context, serviceName string, service *corev1.Service) (*corev1.Service, error) {

	k8Client := c.MustGet("kubeClient").(*kubernetes.Clientset)

	service, err := k8Client.CoreV1().Services("default").Update(context.Background(), service, metav1.UpdateOptions{})
	if err != nil {
		log.Fatal(err)
	}

	return service, nil

}

// func patchService(c *gin.Context, serviceName string, service *corev1.Service) (*corev1.Service, error) {
// 	k8Client := c.MustGet("kubeClient").(*kubernetes.Clientset)
// 	servicesClient := k8Client.CoreV1().Services("default")

// 	// Get the current service object
// 	currentService, err := servicesClient.Get(context.Background(), serviceName, metav1.GetOptions{})
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Update the current service object with the new fields
// 	currentService.Spec.Ports = service.Spec.Ports

// 	// Patch the service with the updated object
// 	patchedService, err := servicesClient.Patch(context.Background(), serviceName, types.MergePatchType, []byte{}, metav1.PatchOptions{}, currentervice, "application/json-patch+json")
// 	if err != nil {
// 		return nil, err
// 	}

// 	return patchedService, nil
// }
