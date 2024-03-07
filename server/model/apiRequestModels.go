package model

type ComputeRequest struct {
	Expressions []Expression `json:"expressions"`
	CallbackURL string       `json:"callback_url"`
}

type CreateDeploymentApiRequest struct {
	// kubernete deployment request model
	Name          string `json:"name"`
	Namespace     string `json:"namespace"`
	Replicas      int    `json:"replicas"`
	AppLabel      string `json:"app_label"`
	ContainerName string `json:"container_name"`
	ImageName     string `json:"image_name"`
	Port          int    `json:"port"`
}

type UpdateDeploymentApiRequest struct {
	// kubernete deployment request model
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	Replicas  int    `json:"replicas"`
}

type CreatePodApiRequest struct {
	Name          string `json:"name"`
	ContainerName string `json:"continerNmae"`
	Port          int    `json:"port"`
	ImageName     string `json:"imageName"`
}

type CreateServiceApiRequest struct {
	Name          string
	ContainerName string
	Port          int
	TargetPort    int
	Selector      string
	Type          string
}
