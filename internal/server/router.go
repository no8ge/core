package server

import (
	"github.com/no8ge/core/pkg/helm"
	"github.com/no8ge/core/pkg/k8s"

	"github.com/gin-gonic/gin"
)

func Run() {
	r := gin.Default()

	r.POST("/helm/repo", helm.AddHelmRepo)
	r.GET("/helm/repos", helm.ListHelmRepos)
	r.DELETE("/helm/repo", helm.DeleteHelmRepo)

	r.POST("/install", helm.InstallHelmChart)
	r.PUT("/upgrade", helm.UpgradeHelmChart)
	r.DELETE("/uninstall", helm.UninstallHelmChart)
	r.POST("/rollback", helm.RollbackHelmChart)
	r.GET("/releases", helm.ListHelmCharts)

	r.POST("/pods", k8s.CreatePod)
	r.DELETE("/pods/:namespace/:name", k8s.DeletePod)
	r.GET("/pods/:namespace", k8s.ListPods)
	r.GET("/pods/:namespace/:pod/exec", k8s.ExecPod)

	r.Run(":8080")
}
