package beanstalk

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"github.com/ohsu-comp-bio/funnel/config"
	"os"
	"text/template"
)

var dockerfileTpl = `
FROM {{.Image}}

ADD ./config.yaml /opt/config.yaml

ENTRYPOINT ["funnel", "server", "run", "--config", "/opt/config.yaml"]
`

func createBundle(zipPath string, image string, confPath string) error {
	// create zip file
	zFile, err := os.Create(zipPath)
	if err != nil {
		return err
	}
	w := zip.NewWriter(zFile)

	// read / create config file
	cFile, err := w.Create("config.yaml")
	if err != nil {
		return err
	}
	conf := config.DefaultConfig()
	err = config.ParseFile(confPath, &conf)
	if err != nil {
		return fmt.Errorf("error reading config file: %v", err)
	}
	conf.Server.HTTPPort = "8000"
	conf.Server.RPCPort = "9090"
	conf.Server.Logger.OutputFile = "/var/log/funnel/funnel.log"
	cBinary := conf.ToYaml()
	_, err = cFile.Write(cBinary)
	if err != nil {
		return err
	}

	// create dockerfile
	dFile, err := w.Create("Dockerfile")
	if err != nil {
		return err
	}

	tpl, err := template.New("Dockerfile").Parse(dockerfileTpl)
	tpl.Execute(dFile, map[string]interface{}{
		"Image": image,
	})
	if err != nil {
		return err
	}

	// create Dockerrun.aws.json file
	jFile, err := w.Create("Dockerrun.aws.json")
	if err != nil {
		return err
	}
	jContent := map[string]interface{}{
		"AWSEBDockerrunVersion": "1",
		"Ports": []map[string]string{
			{
				"ContainerPort": conf.Server.HTTPPort,
			},
		},
		"Logging": "/var/log/funnel",
	}
	jBinary, err := json.Marshal(jContent)
	if err != nil {
		return fmt.Errorf("error marshalling json for Dockerrun.aws.json: %v", err)
	}
	_, err = jFile.Write(jBinary)
	if err != nil {
		return err
	}

	// close zip writer
	err = w.Close()
	if err != nil {
		return err
	}

	return nil
}