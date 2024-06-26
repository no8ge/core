package helm

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func AddHelmRepo(c *gin.Context) {
	var req struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	helmClient, err := NewHelmClient("default")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := helmClient.AddRepo(req.Name, req.URL); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "repository added"})
}

func ListHelmRepos(c *gin.Context) {
	helmClient, err := NewHelmClient("default")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	repos, err := helmClient.ListRepos()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"repositories": repos})
}

func DeleteHelmRepo(c *gin.Context) {
	var req struct {
		Name string `json:"name"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	helmClient, err := NewHelmClient("default")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := helmClient.DeleteRepo(req.Name); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "repository deleted"})
}

func InstallHelmChart(c *gin.Context) {
	var req struct {
		Repo         string                 `json:"repo"`
		ChartName    string                 `json:"chartName"`
		ChartVersion string                 `json:"chartVersion"`
		ReleaseName  string                 `json:"releaseName"`
		Namespace    string                 `json:"namespace"`
		Values       map[string]interface{} `json:"values"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	helmClient, err := NewHelmClient(req.Namespace)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	releases, err := helmClient.InstallChart(req.Repo, req.ChartName, req.ChartVersion, req.ReleaseName, req.Values)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"releases": releases})
}
func ListHelmCharts(c *gin.Context) {
	var req struct {
		Namespace string `json:"namespace"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	helmClient, err := NewHelmClient(req.Namespace)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	releases, err := helmClient.ListReleases(req.Namespace)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"releases": releases})
}

func UpgradeHelmChart(c *gin.Context) {
	var req struct {
		Repo         string                 `json:"repo"`
		ReleaseName  string                 `json:"releaseName"`
		ChartName    string                 `json:"chartName"`
		ChartVersion string                 `json:"chartVersion"`
		Namespace    string                 `json:"namespace"`
		Values       map[string]interface{} `json:"values"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	helmClient, err := NewHelmClient(req.Namespace)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	releases, err := helmClient.UpgradeChart(req.Repo, req.ReleaseName, req.ChartName, req.ChartVersion, req.Values)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"releases": releases})

}

func RollbackHelmChart(c *gin.Context) {
	var req struct {
		ReleaseName string `json:"releaseName"`
		Revision    int    `json:"revision"`
		Namespace   string `json:"namespace"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	helmClient, err := NewHelmClient(req.Namespace)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = helmClient.RollbackChart(req.ReleaseName, req.Revision)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "chart rolled back"})
}
func UninstallHelmChart(c *gin.Context) {
	var req struct {
		ReleaseName string `json:"releaseName"`
		Namespace   string `json:"namespace"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	helmClient, err := NewHelmClient(req.Namespace)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = helmClient.UninstallChart(req.ReleaseName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "chart uninstalled"})
}
