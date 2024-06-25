package k8s

import (
	"bytes"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/remotecommand"
)

func CreatePod(c *gin.Context) {
	var pod v1.Pod
	if err := c.ShouldBindJSON(&pod); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	podsClient := Clientset.CoreV1().Pods(pod.Namespace)
	_, err := podsClient.Create(context.TODO(), &pod, metav1.CreateOptions{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, pod)
}

func DeletePod(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")

	podsClient := Clientset.CoreV1().Pods(namespace)
	err := podsClient.Delete(context.TODO(), name, metav1.DeleteOptions{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}

func ListPods(c *gin.Context) {
	namespace := c.Param("namespace")

	podsClient := Clientset.CoreV1().Pods(namespace)
	pods, err := podsClient.List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, pods.Items)
}

func ExecPod(c *gin.Context) {
	namespace := c.Param("namespace")
	podName := c.Param("pod")
	containerName := c.Query("container")
	command := c.QueryArray("command")

	req := Clientset.CoreV1().RESTClient().
		Post().
		Resource("pods").
		Name(podName).
		Namespace(namespace).
		SubResource("exec").
		Param("container", containerName).
		Param("stderr", "true").
		Param("stdout", "true").
		Param("tty", "false").
		Param("stdin", "false")

	for _, cmd := range command {
		req.Param("command", cmd)
	}

	exec, err := remotecommand.NewSPDYExecutor(Config, "POST", req.URL())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var stdout, stderr bytes.Buffer

	streamOptions := remotecommand.StreamOptions{
		Stdout: &stdout,
		Stderr: &stderr,
		Tty:    false,
	}

	err = exec.StreamWithContext(context.TODO(), streamOptions)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"stdout": stdout.String(),
		"stderr": stderr.String(),
	})
}
