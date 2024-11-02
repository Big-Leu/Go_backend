package controllers

import (
	"bytes"
	"context"
	"fmt"
	"kubequntumblock/internal/initializer"
	"kubequntumblock/schemas"
	"net/http"
	"os"
	"strings"
	"github.com/gin-gonic/gin"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/remotecommand"
	"k8s.io/client-go/util/retry"
	
)

type KubeService interface {

	GetPods(c *gin.Context)
	CreatePods(c *gin.Context)
	ExecCommandInPod(c *gin.Context)
	PatchPod(c *gin.Context)
	DeletePod(c *gin.Context)

}
type LogStreamer struct {
	b bytes.Buffer
}

func (l *LogStreamer) String() string {
	return l.b.String()
}


func (l *LogStreamer) Write(p []byte) (n int, err error) {
	a := strings.TrimSpace(string(p))
	l.b.WriteString(a)
	return len(p), nil
}

func CreatePods(c *gin.Context) {
	var pod schemas.PodCreateSchema
	if err := c.ShouldBindJSON(&pod); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fmt.Print(pod)

	podDefintion := &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: pod.PodName,
			Namespace:    initializer.K.Namespace,
		},
		Spec: v1.PodSpec{
			Containers: []v1.Container{
				{
					Name:  pod.ContainerName,
					Image: pod.Image,
				},
			},
		},
	}
	newPods, err := initializer.K.Pod.Create(context.Background(), podDefintion, metav1.CreateOptions{})
	
    var createdPod []schemas.PodInfo
	var containers []schemas.ContainerInfo
	for _, container := range newPods.Spec.Containers {
		containers = append(containers, schemas.ContainerInfo{
			Name:  container.Name,
			Image: container.Image,
		})
    }
	fmt.Printf(newPods.Name)
	createdPod = append(createdPod, schemas.PodInfo{
		Name:       newPods.Name,
		Containers: containers,
		Status:      string(newPods.Status.Phase),
	})
	
	if err != nil {
		panic(err.Error())
	}
	c.IndentedJSON(http.StatusAccepted, createdPod)

}

func PatchPod(c *gin.Context){
	podName := c.Query("podname")
    var PodName string
	retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {

		currentPod, updateErr := initializer.K.Pod.Get(context.TODO(), podName, metav1.GetOptions{})
		if updateErr != nil {
			panic(updateErr.Error())
		}

		currentPod.Spec.Containers[0].Image = "nginx:1.25.4"
		updatedPod, updateErr := initializer.K.Pod.Update(context.TODO(), currentPod, metav1.UpdateOptions{})
		fmt.Printf("Updated pod: %s", updatedPod.Name)
        PodName = updatedPod.Name
		return updateErr
	})
	if retryErr != nil {
		c.IndentedJSON(http.StatusNotFound,gin.H{"message":"Error while processing"})
		panic(retryErr.Error())
	}
	c.IndentedJSON(http.StatusOK, gin.H{"PodName Updated":PodName})
}
func DeletePod(c *gin.Context){
	podName := c.Query("podname")
	deleteErr := initializer.K.Pod.Delete(context.TODO(),podName, metav1.DeleteOptions{})
	if deleteErr != nil {
		c.IndentedJSON(http.StatusNotFound,gin.H{"message":"Error while processing"})
		panic(deleteErr.Error())
	}
	c.IndentedJSON(http.StatusOK,gin.H{"PodName Deleted":podName})
}
func ExecCommandInPod(c *gin.Context) {
	var pod schemas.Command
	if err := c.ShouldBindJSON(&pod); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	pythonCommand := fmt.Sprintf(`python script_to.py --endpoint_type %s --function_name %s --route %s --function_file '%s'`, pod.EndPoint,
		pod.FunctionName, pod.Route, pod.FunctionBody)

	fmt.Print(pythonCommand)
	command := []string{"/bin/sh", "-c", pythonCommand}

	req := initializer.K.Client.CoreV1().RESTClient().
		Post().
		Resource("pods").
		Name(pod.PodName).
		Namespace(pod.Namespace).
		SubResource("exec").
		Param("container", pod.ContainerName).
		Param("stdout", "true").
		Param("stderr", "true").
		Param("tty", "false")

	// Add command parameters correctly
	for _, cmd := range command {
		req.Param("command", cmd)
	}
	l := &LogStreamer{}
	Executor, err := remotecommand.NewSPDYExecutor(initializer.K.Config, http.MethodPost, req.URL())

	Executor.StreamWithContext(context.Background(), remotecommand.StreamOptions{
		Stdin:  os.Stdin,
		Stdout: l,
		Stderr: nil,
		Tty:    true,
	})

	if err != nil {
		c.String(http.StatusInternalServerError, "Error executing command: %s", err.Error())
		return
	}

	c.String(http.StatusOK, "Command executed successfully")
}


func GetPods(c *gin.Context) {

	pods, err := initializer.K.Pod.List(context.Background(), metav1.ListOptions{})

	if err != nil {
		panic(err.Error())
	}

	var podDetails []schemas.PodInfo
	for _, pod := range pods.Items {
		var containers []schemas.ContainerInfo

		for _, container := range pod.Spec.Containers {
			containers = append(containers, schemas.ContainerInfo{
				Name:  container.Name,
				Image: container.Image,
			})
		}

		podDetails = append(podDetails, schemas.PodInfo{
			Name:       pod.Name,
			Containers: containers,
			Status:     string(pod.Status.Phase),
		})
	}

	if err != nil {
		panic(err.Error())
	}

	c.IndentedJSON(http.StatusOK, podDetails)
}