package main

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
	"k8s.io/helm/pkg/chartutil"
	"k8s.io/helm/pkg/engine"
	"k8s.io/helm/pkg/proto/hapi/chart"
)

func main() {
	loadedChart, err := chartutil.LoadDir("./chart")
	if err != nil {
		log.Fatal("CHART:", err)
		return
	}

	loadedValueStr, err := ioutil.ReadFile("./chart/darwin-application.yaml")
	if err != nil {
		log.Fatal("VALUES: ", err)
		return
	}

	config := &chart.Config{
		Raw: string(loadedValueStr),
	}

	overridedValues, err := chartutil.ToRenderValues(loadedChart, config, chartutil.ReleaseOptions{})
	if err != nil {
		log.Fatal("TO_RENDER_VALUES: ", err)
		return
	}

	tplEngine := engine.New()

	tmpRender, err := tplEngine.Render(loadedChart, overridedValues.AsMap())
	if err != nil {
		log.Fatal("RENDER: ", err)
		return
	}

	yamlData := make(map[interface{}]interface{})
	deplpymentYAMLStr := tmpRender["darwin-k8s-chart-values/templates/deployment.yaml"]
	err = yaml.Unmarshal([]byte(deplpymentYAMLStr), &yamlData)
	if err != nil {
		log.Fatal("YAML_UNMARSHALL: ", err)
		return
	}

	d, err := yaml.Marshal(&yamlData)
	log.Println(d)

	err = ioutil.WriteFile("./out.yaml", d, 0644)
	if err != nil {
		log.Fatal(err)
		return
	}
}
