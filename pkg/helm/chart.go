package helm

import (
	"fmt"
	"log"
	"os"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/release"
)

func MergeMap(mObj ...map[string]string) map[string]string {
	newObj := map[string]string{}
	for _, m := range mObj {
		for k, v := range m {
			newObj[k] = v
		}
	}
	return newObj
}

// install chart by helm
func InstallChart(name string, namespace string, chart string, version string, vals map[string]interface{}, args []string) (*release.Release, error) {
	settings := cli.New()
	actionConfig := new(action.Configuration)

	if err := actionConfig.Init(
		settings.RESTClientGetter(),
		namespace,
		os.Getenv("HELM_DRIVER"),
		log.Printf,
	); err != nil {
		log.Printf("%+v", err)
		os.Exit(1)
	}
	client := action.NewInstall(actionConfig)
	client.ReleaseName = name
	client.Namespace = namespace

	// load Chart
	cp, err := client.ChartPathOptions.LocateChart(fmt.Sprintf("%s/%s-%s.tgz", "https://no8ge.github.io/chartrepo", chart, version), settings)
	if err != nil {
		fmt.Println(err)
	}

	chartRequested, err := loader.Load(cp)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	chartRequested.Values = vals

	rel, err := client.Run(chartRequested, nil)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	fmt.Println("Chart installed successfully!")
	return rel, nil
}

// uninstall chart by helm
func UninstallChart(name string, namespace string, args []string) (*release.UninstallReleaseResponse, error) {
	settings := cli.New()
	actionConfig := new(action.Configuration)

	if err := actionConfig.Init(
		settings.RESTClientGetter(),
		namespace,
		os.Getenv("HELM_DRIVER"),
		log.Printf,
	); err != nil {
		log.Printf("%+v", err)
		os.Exit(1)
	}
	client := action.NewUninstall(actionConfig)

	result, err := client.Run(name)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	fmt.Println("Chart uninstalled successfully!")
	return result, nil
}

// get release list by helm, print table
func ListChart(namespace string, args []string) ([]*release.Release, error) {
	settings := cli.New()
	actionConfig := new(action.Configuration)

	if err := actionConfig.Init(
		settings.RESTClientGetter(),
		namespace,
		os.Getenv("HELM_DRIVER"),
		log.Printf,
	); err != nil {
		log.Printf("%+v", err)
		os.Exit(1)
	}

	client := action.NewList(actionConfig)
	rels, err := client.Run()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	fmt.Println("Chart list successfully!")
	return rels, nil
}
