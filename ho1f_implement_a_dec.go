package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/hashicorp/hcl/v2/hcldec"
)

type Pipeline struct {
	Stages []Stage `hcl:"stage,attr"`
}

type Stage struct {
	Name  string `hcl:"name,attr"`
	Tasks []Task `hcl:"task,attr"`
}

type Task struct {
	Name string `hcl:"name,attr"`
	Cmd  string `hcl:"cmd,attr"`
}

func parsePipelineConfig(configPath string) (*Pipeline, error) {
	configBytes, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var pipeline Pipeline
	err = hcldec.Decode(&pipeline, bytes.NewBuffer(configBytes), nil)
	if err != nil {
		return nil, err
	}

	return &pipeline, nil
}

func main() {
	configPath := "pipeline.config"
	pipeline, err := parsePipelineConfig(configPath)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Pipeline:\n")
	fmt.Printf("  Stages:\n")
	for _, stage := range pipeline.Stages {
		fmt.Printf("  - %s:\n", stage.Name)
		fmt.Printf("    Tasks:\n")
		for _, task := range stage.Tasks {
			fmt.Printf("    - %s: %s\n", task.Name, task.Cmd)
		}
	}
}

func testParsePipelineConfig() {
	config := `
stage "build" {
  task "compile" {
    cmd = "go build main.go"
  }
  task "test" {
    cmd = "go test -v"
  }
}

stage "deploy" {
  task "deploy" {
    cmd = "kubectl deploy -f deployment.yaml"
  }
}
`

	configPath := filepath.Join(os.TempDir(), "pipeline.config")
	err := ioutil.WriteFile(configPath, []byte(config), 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(configPath)

	pipeline, err := parsePipelineConfig(configPath)
	if err != nil {
		log.Fatal(err)
	}

	if len(pipeline.Stages) != 2 {
		log.Fatal("expected 2 stages, got:", len(pipeline.Stages))
	}

	buildStage := pipeline.Stages[0]
	if buildStage.Name != "build" {
		log.Fatal("expected stage name 'build', got:", buildStage.Name)
	}

	if len(buildStage.Tasks) != 2 {
		log.Fatal("expected 2 tasks in build stage, got:", len(buildStage.Tasks))
	}

	deployStage := pipeline.Stages[1]
	if deployStage.Name != "deploy" {
		log.Fatal("expected stage name 'deploy', got:", deployStage.Name)
	}

	if len(deployStage.Tasks) != 1 {
		log.Fatal("expected 1 task in deploy stage, got:", len(deployStage.Tasks))
	}
}

func init() {
	testParsePipelineConfig()
}