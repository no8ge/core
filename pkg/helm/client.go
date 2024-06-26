package helm

import (
	"errors"
	"fmt"
	"log"
	"os"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/getter"
	"helm.sh/helm/v3/pkg/registry"
	"helm.sh/helm/v3/pkg/release"
	"helm.sh/helm/v3/pkg/repo"
)

type HelmClient struct {
	Config   *action.Configuration
	Settings *cli.EnvSettings
}

func NewHelmClient(namespace string) (*HelmClient, error) {
	settings := cli.New()
	settings.SetNamespace(namespace)

	registryClient, err := registry.NewClient(
		registry.ClientOptDebug(true),
		registry.ClientOptEnableCache(true),
		registry.ClientOptWriter(os.Stdout),
	)

	if err != nil {
		return nil, err
	}

	config := new(action.Configuration)
	if err := config.Init(settings.RESTClientGetter(), settings.Namespace(), "secret", log.Printf); err != nil {
		return nil, err
	}
	config.RegistryClient = registryClient

	return &HelmClient{Config: config, Settings: settings}, nil
}

func (hc *HelmClient) AddRepo(name, url string) error {
	repoEntry := repo.Entry{
		Name: name,
		URL:  url,
	}

	chartRepo, err := repo.NewChartRepository(&repoEntry, getter.All(hc.Settings))
	if err != nil {
		return err
	}

	_, err = chartRepo.DownloadIndexFile()
	if err != nil {
		return err
	}

	repoFile := hc.Settings.RepositoryConfig
	f, err := repo.LoadFile(repoFile)
	if err != nil {
		return err
	}

	f.Update(&repoEntry)

	return f.WriteFile(repoFile, 0644)
}

func (hc *HelmClient) ListRepos() ([]*repo.Entry, error) {
	repoFile := hc.Settings.RepositoryConfig
	f, err := repo.LoadFile(repoFile)
	if err != nil {
		return nil, err
	}

	return f.Repositories, nil
}

func (hc *HelmClient) DeleteRepo(name string) error {
	repoFile := hc.Settings.RepositoryConfig
	f, err := repo.LoadFile(repoFile)
	if err != nil {
		return err
	}

	if !f.Remove(name) {
		return errors.New("repository not found")
	}

	return f.WriteFile(repoFile, 0644)
}

func (hc *HelmClient) InstallChart(repo, chartName, chartVersion, releaseName string, values map[string]interface{}) (*release.Release, error) {
	install := action.NewInstall(hc.Config)
	install.ReleaseName = releaseName
	install.Namespace = hc.Settings.Namespace()
	install.Version = chartVersion

	chartPathRef := fmt.Sprintf("%s/%s", repo, chartName)
	chartPath, err := install.ChartPathOptions.LocateChart(chartPathRef, hc.Settings)
	if err != nil {
		return nil, err
	}

	chart, err := loader.Load(chartPath)
	if err != nil {
		return nil, err
	}

	rel, err := install.Run(chart, values)
	return rel, err
}

func (hc *HelmClient) UpgradeChart(repo, releaseName, chartName, chartVersion string, values map[string]interface{}) (*release.Release, error) {
	upgrade := action.NewUpgrade(hc.Config)
	upgrade.Namespace = hc.Settings.Namespace()
	upgrade.Version = chartVersion

	fullName := fmt.Sprintf("%s/%s", repo, chartName)
	chartPath, err := upgrade.ChartPathOptions.LocateChart(fullName, hc.Settings)
	if err != nil {
		return nil, err
	}

	chart, err := loader.Load(chartPath)
	if err != nil {
		return nil, err
	}

	rel, err := upgrade.Run(releaseName, chart, values)
	return rel, err
}

func (hc *HelmClient) RollbackChart(releaseName string, revision int) error {
	rollback := action.NewRollback(hc.Config)
	rollback.Version = revision

	err := rollback.Run(releaseName)
	if err != nil {
		return err
	}

	return nil
}

func (hc *HelmClient) UninstallChart(releaseName string) error {
	uninstall := action.NewUninstall(hc.Config)

	_, err := uninstall.Run(releaseName)
	if err != nil {
		return err
	}

	return nil
}

func (hc *HelmClient) ListReleases(namespace string) ([]*release.Release, error) {
	list := action.NewList(hc.Config)

	releases, err := list.Run()
	if err != nil {
		return nil, err
	}

	return releases, nil
}
