package services

import (
	// "encoding/json"

	// appsv1 "k8s.io/api/core/v1"
	// metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"bytes"
	"context"
	"io"
	"log"

	models "github.com/clusterService/server/model"
	"github.com/gin-gonic/gin"

	// clientv3 "go.etcd.io/etcd/client/v3"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	// "k8s.io/client-go/util/retry"
	// "k8s.io/apimachinery/pkg/runtime/schema"
	// "k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func CreatePod(c *gin.Context, createPodRequest *models.CreatePodRequest) (*models.CreatePodResponse, error) {

	k8Client := c.MustGet("kubeClient").(*kubernetes.Clientset)
	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: createPodRequest.Name,
		},
		Spec: corev1.PodSpec{
			Containers: createPodRequest.Containers,
		},
	}

	podsClient := k8Client.CoreV1().Pods("default")
	result, err := podsClient.Create(context.Background(), pod, metav1.CreateOptions{})

	if err != nil {
		log.Println(err)
		return nil, err
	}

	//  Create an etcd client

	// Create a client configuration
	// config := clientv3.Config{
	// 	Endpoints: []string{"10.110.189.10:30000"}, // Replace with your etcd server address
	// }

	// // Create a new client with the configuration
	// client, err := clientv3.New(config)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer client.Close()

	// // Set a key-value pair
	// ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// _, err = client.Put(ctx, "set-0-pod1", "pod_id")
	// _, err = client.Put(ctx, "set-0-pod2", "podid")
	// cancel()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// // Get the value of the key
	// ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	// resp, err := client.Get(ctx, "set-0-pod1")
	// // allkeyValues, err := client.Get("",)
	// cancel()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// if len(resp.Kvs) == 0 {
	// 	log.Fatal("key not found")
	// }
	// fmt.Printf("Value of key: %s\n", resp.Kvs[0].Value)
	return &models.CreatePodResponse{Pod: result}, nil
}

type PodList struct {
	Pods *corev1.PodList
}

func ListPods(c *gin.Context) (*PodList, error) {

	k8Client := c.MustGet("kubeClient").(*kubernetes.Clientset)
	podsClient := k8Client.CoreV1().Pods("default")
	result, err := podsClient.List(context.Background(), metav1.ListOptions{})

	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &PodList{Pods: result}, nil
}

func DeletePod(c *gin.Context, deletePodRequest *models.DeletePodRequest) (*models.DeletePodResponse, error) {
	log.Println("deletePodRequest.Name: ", deletePodRequest.Name)
	k8Client := c.MustGet("kubeClient").(*kubernetes.Clientset)
	podsClient := k8Client.CoreV1().Pods("default")
	result, err := podsClient.Get(context.Background(), deletePodRequest.Name, metav1.GetOptions{})
	if err != nil {
		log.Println(err)
		return nil, err
	}

	deletePolicy := metav1.DeletePropagationForeground
	err = podsClient.Delete(context.Background(), deletePodRequest.Name, metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	})
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &models.DeletePodResponse{Pod: result}, nil
}

func GetPod(c *gin.Context, getPodRequest *models.GetPodRequest) (*models.GetPodResponse, error) {
	log.Println("getPodRequest.Name: ", getPodRequest.Name)
	k8Client := c.MustGet("kubeClient").(*kubernetes.Clientset)
	podsClient := k8Client.CoreV1().Pods("default")
	result, err := podsClient.Get(context.Background(), getPodRequest.Name, metav1.GetOptions{})
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &models.GetPodResponse{Pod: result}, nil
}

func UpdatePod(c *gin.Context, updatePodRequest *models.UpdatePodRequest) (*models.UpdatePodResponse, error) {
	log.Println("updatePodRequest.Name: ", updatePodRequest.Name)
	k8Client := c.MustGet("kubeClient").(*kubernetes.Clientset)
	podsClient := k8Client.CoreV1().Pods("default")
	result, err := podsClient.Get(context.Background(), updatePodRequest.Name, metav1.GetOptions{})
	if err != nil {
		log.Println(err)
		return nil, err
	}

	result.Spec.Containers = updatePodRequest.Containers
	_, err = podsClient.Update(context.Background(), result, metav1.UpdateOptions{})
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &models.UpdatePodResponse{Pod: result}, nil
}

func PatchPod(c *gin.Context, patchPodRequest *models.PatchPodRequest) (*models.PatchPodResponse, error) {
	log.Println("patchPodRequest.Name: ", patchPodRequest.Name)
	k8Client := c.MustGet("kubeClient").(*kubernetes.Clientset)
	podsClient := k8Client.CoreV1().Pods("default")
	result, err := podsClient.Get(context.Background(), patchPodRequest.Name, metav1.GetOptions{})
	if err != nil {
		log.Println(err)
		return nil, err
	}

	result.Spec.Containers = patchPodRequest.Containers

	_, err = podsClient.Patch(context.Background(), patchPodRequest.Name, "application/strategic-merge-patch+json", []byte(`{"spec":{"containers":[{"name":"nginx","image":"nginx:1.13"}]}}`), metav1.PatchOptions{})
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &models.PatchPodResponse{Pod: result}, nil
}

// ReplacePod replaces the containers in a Kubernetes pod.
// func ReplacePod(c *gin.Context, replacePodRequest *ReplacePodRequest) (*ReplacePodResponse, error) {
// 	log.Println("replacePodRequest.Name: ", replacePodRequest.Name)

// 	if replacePodRequest == nil {
// 		return nil, errors.New("replacePodRequest is nil")
// 	}

// 	k8Client := c.MustGet("kubeClient").(*kubernetes.Clientset)
// 	podsClient := k8Client.CoreV1().Pods("default")

// 	result, err := podsClient.Get(context.Background(), replacePodRequest.Name, metav1.GetOptions{})
// 	if err != nil {
// 		log.Println(err)
// 		return nil, errors.New("failed to get pod")
// 	}

// 	if result == nil {
// 		return nil, errors.New("pod not found")
// 	}

// 	result.Spec.Containers = replacePodRequest.Containers

// 	_, err = podsClient.Replace(context.Background(), result, metav1.UpdateOptions{})
// 	if err != nil {
// 		log.Println(err)
// 		return nil, errors.New("failed to replace pod")
// 	}

// 	return &ReplacePodResponse{Pod: result}, nil
// }

func WatchPod(c *gin.Context, watchPodRequest *models.WatchPodRequest) (*models.WatchPodResponse, error) {
	log.Println("watchPodRequest.Name: ", watchPodRequest.Name)
	k8Client := c.MustGet("kubeClient").(*kubernetes.Clientset)
	podsClient := k8Client.CoreV1().Pods("default")

	watcher, err := podsClient.Watch(context.Background(), metav1.ListOptions{})
	if err != nil {
		log.Println(err)
		return nil, err
	}

	for event := range watcher.ResultChan() {
		log.Println("event.Type: ", event.Type)
		log.Println("event.Object: ", event.Object)
	}

	return &models.WatchPodResponse{}, nil
}

type GetPodLogsRequest struct {
	Name string
}

type GetPodLogsResponse struct {
	Logs string
}

func GetPodLogs(c *gin.Context, getPodLogsRequest *GetPodLogsRequest) (*GetPodLogsResponse, error) {
	log.Println("getPodLogsRequest.Name: ", getPodLogsRequest.Name)
	k8Client := c.MustGet("kubeClient").(*kubernetes.Clientset)
	podsClient := k8Client.CoreV1().Pods("default")

	req := podsClient.GetLogs(getPodLogsRequest.Name, &corev1.PodLogOptions{})
	podLogs, err := req.Stream(context.Background())
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer podLogs.Close()

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, podLogs)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	str := buf.String()
	log.Println("str: ", str)

	return &GetPodLogsResponse{Logs: str}, nil
}

type GetPodExecRequest struct {
	Name string
}

type GetPodExecResponse struct {
	Exec string
}

func GetPodExec(c *gin.Context, getPodExecRequest *GetPodExecRequest) (*GetPodExecResponse, error) {

	log.Println("getPodExecRequest.Name: ", getPodExecRequest.Name)
	k8Client := c.MustGet("kubeClient").(*kubernetes.Clientset)
	podsClient := k8Client.CoreV1().Pods("default")

	req := podsClient.GetLogs(getPodExecRequest.Name, &corev1.PodLogOptions{})
	podLogs, err := req.Stream(context.Background())
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer podLogs.Close()

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, podLogs)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	str := buf.String()
	log.Println("str: ", str)

	return &GetPodExecResponse{Exec: str}, nil
}
