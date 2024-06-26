package helm

import (
	"errors"
	"fmt"
	"log"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/getter"
	"helm.sh/helm/v3/pkg/release"
	"helm.sh/helm/v3/pkg/repo"
)

type HelmClient struct {
	Config   *action.Configuration
	Settings *cli.EnvSettings
}

func NewHelmClient(namespace string) (*HelmClient, error) {
	settings := cli.New()
	config := new(action.Configuration)
	if err := config.Init(settings.RESTClientGetter(), namespace, "configmap", debug); err != nil {
		return nil, err
	}
	return &HelmClient{Config: config, Settings: settings}, nil
}

func debug(format string, v ...interface{}) {
	log.Printf(format, v...)
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

func (hc *HelmClient) InstallChart(repoName, chartName, releaseName, namespace string) (*release.Release, error) {
	install := action.NewInstall(hc.Config)
	install.ReleaseName = releaseName
	install.Namespace = namespace

	fullName := fmt.Sprintf("%s/%s", repoName, chartName)
	chartPath, err := install.LocateChart(fullName, hc.Settings)
	if err != nil {
		return nil, err
	}

	chart, err := loader.Load(chartPath)
	if err != nil {
		return nil, err
	}

	rel, err := install.Run(chart, nil)

	return rel, err
}

func (hc *HelmClient) UpgradeChart(releaseName, chartName, namespace string) error {
	upgrade := action.NewUpgrade(hc.Config)
	upgrade.Namespace = namespace

	chartPath, err := upgrade.ChartPathOptions.LocateChart(chartName, hc.Settings)
	if err != nil {
		return err
	}

	chart, err := loader.Load(chartPath)
	if err != nil {
		return err
	}

	_, err = upgrade.Run(releaseName, chart, nil)
	return err
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
