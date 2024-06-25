package k8s

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.POST("/pods", CreatePod)
	r.DELETE("/pods/:namespace/:name", DeletePod)
	r.GET("/pods/:namespace", ListPods)
	r.GET("/pods/:namespace/:pod/exec", ExecPod)
	return r
}

func TestCreatePod(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	body := `{"apiVersion": "v1", "kind": "Pod", "metadata": {"name": "test-pod", "namespace": "default"}, "spec": {"containers": [{"name": "busybox", "image": "busybox", "command": ["sleep", "3600"]}]}}`
	req, _ := http.NewRequest("POST", "/pods", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDeletePod(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/pods/default/test-pod", nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestListPods(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/pods/default", nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestExecPod(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/pods/default/test-pod/exec?container=busybox&command=ls", nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
