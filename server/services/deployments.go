package services

import (
	// "encoding/json"

	// appsv1 "k8s.io/api/core/v1"
	// metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"context"
	"log"

	"strings"

	models "github.com/clusterService/server/model"
	"github.com/gin-gonic/gin"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	// "k8s.io/client-go/util/retry"
	// "k8s.io/apimachinery/pkg/runtime/schema"
	// "k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func int32Ptr(i int32) *int32 { return &i }

func CreateDeployment(c *gin.Context, createDeploymentRequest *models.CreateDeploymentRequest) (*models.CreateDeploymentResponse, error) {

	k8Client := c.MustGet("kubeClient").(*kubernetes.Clientset)
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: createDeploymentRequest.Name,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(createDeploymentRequest.Replicas),
			Selector: &metav1.LabelSelector{
				MatchLabels: createDeploymentRequest.Selector,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: createDeploymentRequest.Labels,
				},
				Spec: corev1.PodSpec{
					Containers: createDeploymentRequest.Containers,
				},
			},
		},
	}
	result, err := k8Client.AppsV1().Deployments("default").Create(context.Background(), deployment, metav1.CreateOptions{})
	if err != nil {
		log.Println(err)
		if strings.Contains(err.Error(), "already exists") {
			log.Println("Deployment already exists")
			// delete the deployment
			// deletePolicy := metav1.DeletePropagationForeground
			// if err := k8Client.AppsV1().Deployments("default").Delete(context.Background(), "http1", metav1.DeleteOptions{
			// 	PropagationPolicy: &deletePolicy,
			// }); err != nil {
			// 	log.Fatal(err)
			// }
		}
	}

	// get the status of deployment
	for {
		deploymentStatus, err := k8Client.AppsV1().Deployments("default").Get(context.Background(), createDeploymentRequest.Name, metav1.GetOptions{})
		if err != nil {
			log.Fatal(err)
		}
		if deploymentStatus.Status.ReadyReplicas == 4 {

			break
		}
	}

	return &models.CreateDeploymentResponse{Deployment: result}, nil

}

func UpdateDeployment(c *gin.Context, updateDeploymentRequest *models.UpdateDeploymentRequest) (*models.UpdateDeploymentResponse, error) {

	k8Client := c.MustGet("kubeClient").(*kubernetes.Clientset)
	deployment, err := k8Client.AppsV1().Deployments(updateDeploymentRequest.Namespace).Get(context.Background(), updateDeploymentRequest.Name, metav1.GetOptions{})
	if err != nil {
		log.Fatal(err)
	}
	deployment.Spec.Replicas = int32Ptr(updateDeploymentRequest.Replicas)
	result, err := k8Client.AppsV1().Deployments(updateDeploymentRequest.Namespace).Update(context.Background(), deployment, metav1.UpdateOptions{})
	if err != nil {
		log.Fatal(err)
	}

	return &models.UpdateDeploymentResponse{Deployment: result}, nil
}

func GetDeployments(c *gin.Context) (*models.DeploymentList, error) {

	k8Client := c.MustGet("kubeClient").(*kubernetes.Clientset)
	deployments, err := k8Client.AppsV1().Deployments("default").List(context.Background(), metav1.ListOptions{})
	if err != nil {
		log.Fatal(err)
	}
	for _, deployment := range deployments.Items {
		log.Println(deployment.Name)
	}

	return &models.DeploymentList{Deployments: deployments}, nil
}

func GetDeployment(c *gin.Context, name string) (*appsv1.Deployment, error) {

	k8Client := c.MustGet("kubeClient").(*kubernetes.Clientset)
	deployment, err := k8Client.AppsV1().Deployments("default").Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		log.Fatal(err)
	}
	return deployment, nil
}

func DeleteDeployment(c *gin.Context, name string) (string, error) {

	k8Client := c.MustGet("kubeClient").(*kubernetes.Clientset)
	deletePolicy := metav1.DeletePropagationForeground
	if err := k8Client.AppsV1().Deployments("default").Delete(context.Background(), name, metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	}); err != nil {
		log.Fatal(err)
	}
	return "Deployment deleted", nil

}

// func GetDeployment(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
// 	vars := mux.Vars(r)

// 	title := vars["title"]
// 	project := getProjectOr404(db, title, w, r)
// 	if project == nil {
// 		return
// 	}
// 	respondJSON(w, http.StatusOK, project)
// }

// func UpdateDeployment(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
// 	vars := mux.Vars(r)

// 	title := vars["title"]
// 	project := getProjectOr404(db, title, w, r)
// 	if project == nil {
// 		return
// 	}

// 	decoder := json.NewDecoder(r.Body)
// 	if err := decoder.Decode(&project); err != nil {
// 		respondError(w, http.StatusBadRequest, err.Error())
// 		return
// 	}
// 	defer r.Body.Close()

// 	if err := db.Save(&project).Error; err != nil {
// 		respondError(w, http.StatusInternalServerError, err.Error())
// 		return
// 	}
// 	respondJSON(w, http.StatusOK, project)
// }

// func DeleteDeployment(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
// 	vars := mux.Vars(r)

// 	title := vars["title"]
// 	project := getProjectOr404(db, title, w, r)
// 	if project == nil {
// 		return
// 	}
// 	if err := db.Delete(&project).Error; err != nil {
// 		respondError(w, http.StatusInternalServerError, err.Error())
// 		return
// 	}
// 	respondJSON(w, http.StatusNoContent, nil)
// }
