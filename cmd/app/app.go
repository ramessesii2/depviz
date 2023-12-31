package app

import (
	"encoding/json"
	"fmt"
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
	tmpDir, err := os.MkdirTemp(currentDir, "go-cloned-repo-temp")

	if err != nil {
		log.Fatal("Failed to create temporary directory:", err)
	}
	defer os.RemoveAll(tmpDir)

	cmd := exec.Command("git", "clone", "--depth", "1", "--branch", ref, repoURL, tmpDir)
	if err := cmd.Run(); err != nil {
		errMsg := "Failed to clone repository"
		if exitError, ok := err.(*exec.ExitError); ok {
			// Check if the exit status is 0 or 1
			if exitError.ExitCode() == 1 {
				log.Fatal(errMsg + ": Command is invalid")
			}
		}
		log.Fatal(errMsg)
	}

	// Get the module name
	modInfoCmd := exec.Command("go", "list", "-m", "-json")
	modInfoCmd.Dir = tmpDir
	modInfoOutput, err := modInfoCmd.Output()
	if err != nil {
		log.Fatal("Failed to get module info:", err)
	}

	var data map[string]interface{}

	err = json.Unmarshal(modInfoOutput, &data)
	if err != nil {
		log.Fatal(err)
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
	dependencyGeneratorInitializer(depMap, rootArtifact)
	if err != nil {
		log.Fatal("Failed to generate dependency tree:", err)
	}

	// Generate the JSON file
	targetFileName := currentDir + "/dependency_tree.json"
	generator.GenerateJson(rootArtifact, targetFileName)
	fmt.Println("Dependency Tree generated successfully. Output file:", targetFileName)
}

func dependencyGeneratorInitializer(depMap map[string]generator.ArtifactMap, rootArtifact *artifact.Artifact) error {
	// Generate the dependency tree for the root artifact
	for _, dependency := range depMap[rootArtifact.Name].Dependencies {
		splitDependency := strings.Split(dependency, "@")
		var dependencyName, dependencyVersion string
		if len(splitDependency) >= 2 {
			dependencyName = splitDependency[0]
			dependencyVersion = splitDependency[1]
		} else {
			return fmt.Errorf("failed to split dependency: %s", dependency)
		}
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
	return nil
}
