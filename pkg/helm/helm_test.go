package helm_test

import (
	"testing"

	"github.com/no8ge/core/pkg/helm"
)

func TestHelmInstall(t *testing.T) {
	vals := make(map[string]interface{})
	helm.InstallChart("test", "default", "pytest", "1.0.0", vals, nil)
}

func TestHelmList(t *testing.T) {
	helm.ListChart("", nil)
}

func TestHelmUninstall(t *testing.T) {
	helm.UninstallChart("test", "default", nil)
}
