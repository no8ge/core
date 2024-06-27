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

	r.POST("/namespaces/:namespace/pods", CreatePod)
	r.GET("/namespaces/:namespace/pods/:name", GetPod)
	r.GET("/namespaces/:namespace/pods", ListPods)
	r.POST("/namespaces/:namespace/pods/:name/exec", ExecPod)
	r.DELETE("/namespaces/:namespace/pods/:name", DeletePod)

	return r
}

func TestCreatePod(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	body := `{"apiVersion": "v1", "kind": "Pod", "metadata": {"name": "test-pod", "namespace": "atop-agent-test-system"}, "spec": {"containers": [{"name": "busybox", "image": "busybox", "command": ["sleep", "3600"]}]}}`
	req, _ := http.NewRequest("POST", "/namespaces/atop-agent-test-system/pods", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetPod(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/namespaces/atop-agent-test-system/pods/test-pod", nil)
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
func TestListPods(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/namespaces/atop-agent-test-system/pods", nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestExecPod(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/namespaces/atop-agent-test-system/pods/test-pod/exec?container=busybox&command=ls", nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDeletePod(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/namespaces/atop-agent-test-system/pods/test-pod", nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
