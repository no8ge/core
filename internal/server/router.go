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

	r.POST("/helm/install", helm.InstallHelmChart)
	r.PUT("/helm/upgrade", helm.UpgradeHelmChart)
	r.DELETE("/helm/uninstall", helm.UninstallHelmChart)
	r.POST("/helm/rollback", helm.RollbackHelmChart)
	r.GET("/helm/releases", helm.ListHelmCharts)
	r.GET("/helm/release", helm.GetHelmChart)

	r.POST("/namespaces/:namespace/pods", k8s.CreatePod)
	r.GET("/namespaces/:namespace/pods/:name", k8s.GetPod)
	r.DELETE("/namespaces/:namespace/pods/:name", k8s.DeletePod)
	r.GET("/namespaces/:namespace/pods", k8s.ListPods)
	r.POST("/namespaces/:namespace/pods/:name/exec", k8s.ExecPod)

	r.Run(":8080")
}
