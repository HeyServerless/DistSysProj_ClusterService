package server

import (
	"context"
	"flag"
	"fmt"
	"log"
	"path/filepath"
	"strconv"
	"time"

	// EKS requires the AWS SDK to be imported

	"github.com/gin-gonic/gin"
	clientv3 "go.etcd.io/etcd/client/v3"

	// "honnef.co/go/tools/conf
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"

	// pb "github.com/example/mypackage"

	"github.com/clusterService/server/handler"
	models "github.com/clusterService/server/model"
	services "github.com/clusterService/server/services"
	utils "github.com/clusterService/server/utils"
)

// Server is srv struct that holds srv Kubernetes client
type Server struct {
	KubeClient *kubernetes.Clientset
}

const ETCDIP = "34.198.214.90:30000"

func ApiMiddleware(cli *kubernetes.Clientset) gin.HandlerFunc {
	// do something with the request
	return func(c *gin.Context) {
		// do something with the request

		c.Set("kubeClient", cli)
		c.Next()
	}
}

func (srv *Server) Initialize() {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err)
	}
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}
	srv = &Server{KubeClient: client}
	ctx := context.Background()

	//listing the namespaces
	namespaces, _ := client.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
	for _, namespace := range namespaces.Items {
		fmt.Println(namespace.Name)
	}
	// initPods := initSetOfPodCreation(client)
	initPods := initPodCreation(client)
	fmt.Println(initPods)

	fmt.Println("=================================starting server=================================")
	r := gin.Default()
	r.Use(ApiMiddleware(client))

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	/** temop routes start */
	r.POST("/CreateDeployment", handler.CreateDeployment)
	r.POST("/UpdateDeployment", handler.UpdateDeployment)
	r.GET("/getDeployment", handler.GetDeployment)
	r.GET("/getAllDeployments", handler.GetAllDeployments)
	r.DELETE("/deleteDeployment", handler.DeleteDeployment)

	r.POST("/CreateService", handler.CreateService)
	r.GET("/getService", handler.GetService)
	r.GET("/getAllServices", handler.GetAllServices)
	r.DELETE("/deleteService", handler.DeleteService)

	r.POST("/CreatePod", handler.CreatePod)
	r.GET("/getPod", handler.GetPod)
	r.GET("/getAllPods", handler.GetAllPods)
	r.DELETE("/deletePod", handler.DeletePod)

	r.Run(":8081")
}

func initPodCreation(cli *kubernetes.Clientset) string {
	c := &gin.Context{}

	c.Set("kubeClient", cli)
	fmt.Println("Creating Pods")
	for _, image := range utils.ImagesList {
		fmt.Println("Creating Deployment for %s", image)
		deployment, err := services.CreateDeployment(c, &models.CreateDeploymentRequest{
			Name:     image.FuncName,
			Replicas: 4,
			Selector: map[string]string{
				"app": image.FuncName,
			},
			Labels: map[string]string{
				"app": image.FuncName,
			},
			Containers: []corev1.Container{
				{
					Name:  image.FuncName + "-container",
					Image: image.ImageName,
					Ports: []corev1.ContainerPort{
						{
							ContainerPort: image.Port,
						},
					},
				},
			},
		})
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(deployment)
		fmt.Println("=================================")
		fmt.Println("deployment done")
		// strings.ReplaceAll(strings.ReplaceAll(image, ":", "-"), "/", "-")
		service, err := services.CreateService(c, &models.CreateServiceRequest{
			Name: image.FuncName,
			Selector: map[string]string{
				"app": image.FuncName,
			},
			Ports: []corev1.ServicePort{
				{
					Port:     image.Port,
					NodePort: image.NodePort,
				},
			},
			Type: "NodePort",
		})
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(service)

	}
	return "init Pods Created Successfully"
}

func initSetOfPodCreation(cli *kubernetes.Clientset) string {
	c := &gin.Context{}

	c.Set("kubeClient", cli)
	fmt.Println("Creating Pods")
	for _, image := range utils.PbftImagesList {

		fmt.Println("Creating set of pods and registering them with etcd for %s", image)
		fmt.Println("Creating:", image.ImageName+"-worker:latest")

		// create worker pod for each operation image
		createPodRequest := models.CreatePodRequest{
			Name: image.FuncName + "-worker",
			Containers: []corev1.Container{
				{
					Name:  image.FuncName + "-worker",
					Image: image.ImageName + "-worker:latest",
					Ports: []corev1.ContainerPort{
						{
							ContainerPort: image.Port,
						},
					},
					ImagePullPolicy: corev1.PullNever,
				},
			},
		}
		createWorkerPodResponse, er := services.CreatePod(c, &createPodRequest)
		if er != nil {
			fmt.Println(er)
		}
		fmt.Println("pod creation done")
		fmt.Println(createWorkerPodResponse)
		fmt.Println("================================================")
		// create a service which exposes the worker pod via nodeport
		createServiceRequest := models.CreateServiceRequest{
			Name: image.FuncName + "-worker",
			Selector: map[string]string{
				"app": image.FuncName + "-worker",
			},
			Ports: []corev1.ServicePort{
				{
					Port:     image.Port,
					Protocol: corev1.ProtocolTCP,
					NodePort: image.NodePort,
				},
			},
			Type: "NodePort",
		}
		createServiceResponse, er := services.CreateService(c, &createServiceRequest)
		if er != nil {
			fmt.Println(er)
		}
		fmt.Println("================================================================service creation done================================================================")
		fmt.Println(createServiceResponse)
		//  Create an etcd client
		etcdClient, err := clientv3.New(clientv3.Config{
			Endpoints:   []string{ETCDIP},
			DialTimeout: 5 * time.Second,
		})
		if err != nil {
			log.Fatal(err)
		}
		defer etcdClient.Close()

		// Register the pod with etcd
		key := fmt.Sprintf("pods/%s", "set-"+strconv.Itoa(int(image.Id))+"_"+image.FuncName+"-worker")
		value := fmt.Sprintf("%s:%s", createWorkerPodResponse.Pod.Status.PodIP, ":"+strconv.Itoa(int(image.Port)))
		_, err = etcdClient.Put(context.Background(), key, value)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Registered pod with etcd: %s -> %s\n", key, value)

		// for loop of to iterate over the 4 pods an deploye the replicas and register them
		for i := 0; i < 5; i++ {
			fmt.Println(i)
			// create replica pod for each id
			createPodRequest := models.CreatePodRequest{
				Name: image.FuncName + "-replica-" + strconv.Itoa(i),
				Containers: []corev1.Container{
					{
						Name:  image.FuncName + "-replica-" + strconv.Itoa(i),
						Image: image.ImageName + "-replica:latest",
						Ports: []corev1.ContainerPort{
							{
								ContainerPort: image.Port,
							},
						},
						ImagePullPolicy: corev1.PullNever,
					},
				},
			}
			createReplicaPodResponse, er := services.CreatePod(c, &createPodRequest)
			if er != nil {
				fmt.Println(er)
			}
			fmt.Println("pod creation done")

			fmt.Println(createReplicaPodResponse)
			fmt.Println("================================================")
			//  Create an etcd client
			etcdClient, err := clientv3.New(clientv3.Config{
				Endpoints:   []string{ETCDIP},
				DialTimeout: 5 * time.Second,
			})
			if err != nil {
				log.Fatal(err)
			}
			defer etcdClient.Close()
			// Register the pod with etcd
			key := fmt.Sprintf("pods/%s", "set-"+strconv.Itoa(int(image.Id))+"_"+image.FuncName+"-replica-"+strconv.Itoa(i))
			value := fmt.Sprintf("%s:%s", createReplicaPodResponse.Pod.Status.PodIP, ":"+strconv.Itoa(int(image.Port)))
			_, err = etcdClient.Put(context.Background(), key, value)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("Registered pod with etcd: %s -> %s\n", key, value)
		}

	}

	// strings.ReplaceAll(strings.ReplaceAll(image, ":", "-"), "/", "-")
	// service, err := services.CreateService(c, &models.CreateServiceRequest{
	// 	Name: image.FuncName,
	// 	Selector: map[string]string{
	// 		"app": image.FuncName,
	// 	},
	// 	Ports: []corev1.ServicePort{
	// 		{
	// 			Port:     image.Port,
	// 			NodePort: image.NodePort,
	// 		},
	// 	},
	// 	Type: "NodePort",
	// })
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println(service)
	return "init Pods Created Successfully"
}
