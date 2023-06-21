package app

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/ramessesii2/depviz/pkg/artifact"
	"github.com/ramessesii2/depviz/pkg/generator"
	"github.com/ramessesii2/depviz/pkg/util"
)

func App(repoURL, ref string) {

	// Clone the repository to a temporary directory
	currentDir, err := util.FindRootDir()
	if err != nil {
		log.Fatal("Failed to get current directory:", err)
	}
	tmpDir, err := ioutil.TempDir(currentDir, "go-cloned-repo-temp")
	if err != nil {
		log.Fatal("Failed to create temporary directory:", err)
	}
	defer os.RemoveAll(tmpDir)

	cmd := exec.Command("git", "clone", "--depth", "1", "--branch", ref, repoURL, tmpDir)
	if err := cmd.Run(); err != nil {
		log.Fatal("Failed to clone repository:", err)
	}

	// Get the module name
	modInfoCmd := exec.Command("go", "list", "-m", "-json")
	modInfoCmd.Dir = tmpDir
	modInfoOutput, err := modInfoCmd.Output()
	if err != nil {
		log.Fatal("Failed to get module info:", err)
	}
	// else {
	// 	fmt.Println("modInfoOutput:", string(modInfoOutput))
	// }

	var data map[string]interface{}

	err = json.Unmarshal(modInfoOutput, &data)
	if err != nil {
		fmt.Println(err)
	}

	var version string

	if isTag := util.IsRefTag(tmpDir, ref); err != nil {
		log.Fatal("Failed to check if ref is a tag:", err)
	} else if isTag {
		version = ref
	} else if isBranch := util.IsRefBranch(tmpDir, ref); err != nil {
		log.Fatal("Failed to check if ref is a branch:", err)
	} else if isBranch {
		if ref == "main" {
			version = "latest"
		}
	} else {
		version = ""
	}

	rootArtifact := &artifact.Artifact{
		Name:         data["Path"].(string),
		Version:      version,
		Dependencies: []*artifact.Artifact{},
	}

	// Get all the Dependencies modules in string format
	modules, err := generator.GetDependencies(tmpDir)
	if err != nil {
		log.Fatal("Failed to get dependencies:", err)
	}
	// Generate the dependency Map
	depMap := generator.GenerateDependencyMap(modules)
	if err != nil {
		log.Fatal("Failed to generate dependency tree:", err)
	}

	// Generate the dependency tree for the root artifact
	for _, dependency := range depMap[rootArtifact.Name].Dependencies {
		dependencyName := strings.Split(dependency, "@")[0]
		dependencyVersion := strings.Split(dependency, "@")[1]
		dependencyArtifact := &artifact.Artifact{
			Name:         dependencyName,
			Version:      dependencyVersion,
			Dependencies: []*artifact.Artifact{},
		}
		// Generate the dependency tree for the dependency artifacts
		generator.AppendDependncies(depMap, dependencyArtifact)
		// Append the dependency artifact to the root artifact
		rootArtifact.Dependencies = append(rootArtifact.Dependencies, dependencyArtifact)
	}

	// Generate the JSON file
	targetFileName := currentDir + "/dependency_tree.json"
	generator.GenerateJson(rootArtifact, targetFileName)
	fmt.Println("Dependency Tree generated successfully. Output file:", targetFileName)
}
