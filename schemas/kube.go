package schemas

type Command struct {
	PodName       string `json:"podname"`
	Namespace     string `json:"namespace"`
	ContainerName string `json:"containername"`
	EndPoint      string `json:"endpoint"`
	FunctionName  string `json:"functionname"`
	Route         string `json:"route"`
	FunctionBody  string `json:"functionbody"`
}

type PodInfo struct {
	Name       string          `json:"name"`
	Containers []ContainerInfo `json:"containers"`
	Status     string          `json:"status"`
}

type ContainerInfo struct {
	Name  string `json:"name"`
	Image string `json:"image"`
}

type PodCreateSchema struct {
	PodName       string `json:"podname"`
	NameSpace     string `json:"namespace"`
	ContainerName string `json:"containername"`
	Image         string `json:"image"`
}