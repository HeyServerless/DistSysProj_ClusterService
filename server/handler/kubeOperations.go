package handler

import (
	// "encoding/json"

	"context"
	"fmt"
	"net/http"
	"time"

	// appsv1 "k8s.io/api/core/v1"
	// metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"log"

	services "github.com/clusterService/server/services"

	models "github.com/clusterService/server/model"
	"github.com/gin-gonic/gin"
	clientv3 "go.etcd.io/etcd/client/v3"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

const ETCDIP = "34.198.214.90:30000"

func int32Ptr(i int32) *int32 { return &i }

// Deployment apis start here

func CreateDeployment(c *gin.Context) {

	deployRequest := models.CreateDeploymentApiRequest{}

	if err := c.ShouldBindJSON(&deployRequest); err != nil {
		log.Println("Error: ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	CreateDeploymentRequest := &models.CreateDeploymentRequest{
		Name:     deployRequest.Name,
		Replicas: int32(deployRequest.Replicas),
		Selector: map[string]string{
			"app": deployRequest.Name,
		},
		Labels: map[string]string{
			"app": deployRequest.AppLabel,
		},
		Containers: []corev1.Container{
			{
				Name:  deployRequest.ContainerName,
				Image: deployRequest.ImageName,
				Ports: []corev1.ContainerPort{
					{
						ContainerPort: int32(deployRequest.Port),
					},
				},
			},
		},
	}
	CreateDeploymentResponse := &models.CreateDeploymentResponse{}
	CreateDeploymentResponse, err := services.CreateDeployment(c, CreateDeploymentRequest)

	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"error":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   CreateDeploymentResponse,
	})
}

func GetAllDeployments(c *gin.Context) {

	GetDeploymentsResponse, err := services.GetDeployments(c)

	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"error":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   GetDeploymentsResponse,
	})
}

func GetDeployment(c *gin.Context) {

	deploymentName := c.Param("deploymentName")
	GetDeploymentResponse, err := services.GetDeployment(c, deploymentName)

	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"error":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   GetDeploymentResponse,
	})
}

func DeleteDeployment(c *gin.Context) {

	deploymentName := c.Param("deploymentName")
	if deploymentName == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"error":  "deployment name is required",
		})
		return
	}
	DeleteDeploymentResponse, err := services.DeleteDeployment(c, deploymentName)

	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"error":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   DeleteDeploymentResponse,
	})
}

func UpdateDeployment(c *gin.Context) {
	updateDeployementApiRequest := models.UpdateDeploymentApiRequest{}

	if err := c.ShouldBindJSON(&updateDeployementApiRequest); err != nil {
		log.Println("Error: ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Println(updateDeployementApiRequest)
	UpdateDeploymentRequest := &models.UpdateDeploymentRequest{
		Name:      updateDeployementApiRequest.Name,
		Namespace: updateDeployementApiRequest.Namespace,
		Replicas:  int32(updateDeployementApiRequest.Replicas),
	}
	UpdateDeploymentResponse, err := services.UpdateDeployment(c, UpdateDeploymentRequest)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"error":  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   UpdateDeploymentResponse,
	})
}

// Deployment apis end here

// Pod apis start here

func CreatePod(c *gin.Context) {
	//create pod functionality
	podApiRequest := models.CreatePodApiRequest{}
	fmt.Println("create pod=>")

	client, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{ETCDIP},
		DialTimeout: 5 * time.Second,
	})
	// Create a new client with the configuration
	// client, err := clientv3.New(config)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	// Set a key-value pair
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	_, err = client.Put(ctx, "set-0-pod1", "pod_id")
	_, err = client.Put(ctx, "set-0-pod2", "podid")
	cancel()
	if err != nil {
		log.Fatal(err)
	}

	// Get the value of the key
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	resp, err := client.Get(ctx, "set-0-pod1")
	// allkeyValues, err := client.Get("",)
	cancel()
	if err != nil {
		log.Fatal(err)
	}
	if len(resp.Kvs) == 0 {
		log.Fatal("key not found")
	}
	fmt.Printf("Value of key: %s\n", resp.Kvs[0].Value)

	if err := c.ShouldBindJSON(&podApiRequest); err != nil {
		log.Println("Error: ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	CreatePodRequest := &models.CreatePodRequest{
		Name: podApiRequest.Name,
		Containers: []corev1.Container{
			{
				Name:  podApiRequest.ContainerName,
				Image: podApiRequest.ImageName,
				Ports: []corev1.ContainerPort{
					{
						ContainerPort: int32(podApiRequest.Port),
					},
				},
			},
		},
	}
	CreatePodResponse, err := services.CreatePod(c, CreatePodRequest)

	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"error":  err.Error(),
		})
		return
	}

	// register the pod in etcd

	// etcdClient := etcdv1.GetEtcdClient()
	// etcdClient.Put(context.Background(), "/pods/"+podApiRequest.Name, podApiRequest.Name)

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   CreatePodResponse,
	})
}

//
func GetAllPods(c *gin.Context) {

	GetPodsResponse, err := services.ListPods(c)

	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"error":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   GetPodsResponse,
	})
}

func GetPod(c *gin.Context) {

	podName := c.Param("podName")
	PodRequest := &models.GetPodRequest{
		Name: podName,
	}
	GetPodResponse, err := services.GetPod(c, PodRequest)

	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"error":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   GetPodResponse,
	})
}

func DeletePod(c *gin.Context) {

	podName := c.Param("podName")
	if podName == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"error":  "pod name is required",
		})
	}
	deleteRequest := &models.DeletePodRequest{
		Name: podName,
	}
	DeletePodResponse, err := services.DeletePod(c, deleteRequest)

	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"error":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   DeletePodResponse,
	})
}

// Pod apis end here

// Service apis start here

func CreateService(c *gin.Context) {

	createServiceApiRequest := models.CreateServiceApiRequest{}

	if err := c.ShouldBindJSON(&createServiceApiRequest); err != nil {
		log.Println("Error: ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	CreateServiceRequest := &models.CreateServiceRequest{
		Name: createServiceApiRequest.Name,
		Ports: []corev1.ServicePort{
			{
				Port:       int32(createServiceApiRequest.Port),
				TargetPort: intstr.FromInt(createServiceApiRequest.TargetPort),
			},
		},
		Selector: map[string]string{
			"app": createServiceApiRequest.Selector,
		},
		Type: corev1.ServiceType(createServiceApiRequest.Type),
	}
	CreateServiceResponse := &models.CreateServiceResponse{}
	CreateServiceResponse, err := services.CreateService(c, CreateServiceRequest)

	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"error":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   CreateServiceResponse,
	})
}

func GetAllServices(c *gin.Context) {

	GetServicesResponse, err := services.GetAllServices(c)

	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"error":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   GetServicesResponse,
	})
}

func GetService(c *gin.Context) {

	serviceName := c.Param("serviceName")
	GetServiceResponse, err := services.GetService(c, serviceName)

	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"error":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   GetServiceResponse,
	})
}

func DeleteService(c *gin.Context) {

	serviceName := c.Param("serviceName")
	err := services.DeleteService(c, serviceName)

	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"error":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   serviceName + " Service Deleted Successfully",
	})
}

// Service apis end here
