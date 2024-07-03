package helm

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupRouter() *gin.Engine {
	router := gin.Default()

	router.POST("/helm/repo", AddHelmRepo)
	router.GET("/helm/repos", ListHelmRepos)
	router.DELETE("/helm/repo", DeleteHelmRepo)

	router.POST("/helm/install", InstallHelmChart)
	router.PUT("/helm/upgrade", UpgradeHelmChart)
	router.DELETE("/helm/uninstall", UninstallHelmChart)
	router.POST("/helm/rollback", RollbackHelmChart)
	router.GET("/helm/releases", ListHelmCharts)
	router.GET("/helm/release", GetHelmChart)

	return router
}

func TestAddHelmRepo(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	body := `{"name": "myrepo", "url": "https://charts.bitnami.com/bitnami"}`
	req, _ := http.NewRequest("POST", "/helm/repo", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestListHelmRepos(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/helm/repos", nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDeleteHelmRepo(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	body := `{"name": "myrepo"}`
	req, _ := http.NewRequest("DELETE", "/helm/repo", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestInstallHelmChart(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	// body := `{"repo": "myrepo", "chartName": "nginx", "chartVersion":"18.1.2", "releaseName": "mynginx", "namespace": "kube-system","values":{}}`
	body := `{"repo": "oci://registry-1.docker.io/no8ge", "chartName": "core", "chartVersion":"1.0.0","releaseName": "core", "namespace": "kube-system","values":{}}`
	// body := `{"repo": "oci://172.31.34.177:30002/qingtest", "chartName": "aomaker", "chartVersion":"1.0.0","releaseName": "aomaker", "namespace": "default","values":{}}`

	req, _ := http.NewRequest("POST", "/helm/install", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUpgradeHelmChart(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	body := `{"repo": "myrepo", "chartName": "nginx", "chartVersion":"18.1.1", "releaseName": "mynginx", "namespace": "kube-system","values":{}}`
	// body := `{"repo": "oci://registry-1.docker.io/no8ge", "chartName": "core", "chartVersion":"1.0.0","releaseName": "core", "namespace": "default","values":{}}`

	req, _ := http.NewRequest("PUT", "/helm/upgrade", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRollbackHelmChart(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	body := `{"releaseName": "mynginx", "revision":1, "namespace": "kube-system"}`
	// body := `{"releaseName": "core","revision":"1", "namespace": "default"}`

	req, _ := http.NewRequest("POST", "/helm/rollback", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestListlHelmChart(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	body := `{"namespace": "kube-system"}`
	req, _ := http.NewRequest("GET", "/helm/releases", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGettHelmChart(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	body := `{"releaseName": "mynginx","namespace": "kube-system"}`
	req, _ := http.NewRequest("GET", "/helm/release", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUninstallHelmChart(t *testing.T) {
	router := setupRouter()

	body := `{"releaseName": "mynginx","namespace": "kube-system"}`
	req, _ := http.NewRequest("DELETE", "/helm/uninstall", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
