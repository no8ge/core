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
	body := `{"repoName": "myrepo", "chartName": "nginx", "releaseName": "mynginx", "namespace": "kube-system","values":{}}`
	req, _ := http.NewRequest("POST", "/helm/install", strings.NewReader(body))
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

func TestUninstallHelmChart(t *testing.T) {
	router := setupRouter()

	body := `{"releaseName": "mynginx","namespace": "kube-system"}`
	req, _ := http.NewRequest("DELETE", "/helm/uninstall", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
