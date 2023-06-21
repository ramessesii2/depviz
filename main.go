package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
)

type Artifact struct {
	Name         string      `json:"name"`
	Version      string      `json:"version"`
	Dependencies []*Artifact `json:"dependencies"`
}

func main() {
	repoURL := "https://github.com/ramessesii2/sprayproxy" // Replace with the actual GitHub repository URL
	branchOrTag := "main"                                  // Replace with the desired branch or tag

	// Clone the repository to a temporary directory
	tmpDir, err := ioutil.TempDir("/home/ramesses/RAMESSESII2/depviz", "go-cloned-repo-temp")
	if err != nil {
		log.Fatal("Failed to create temporary directory:", err)
	}
	defer os.RemoveAll(tmpDir)

	cmd := exec.Command("git", "clone", "--depth", "1", "--branch", branchOrTag, repoURL, tmpDir)
	if err := cmd.Run(); err != nil {
		log.Fatal("Failed to clone repository:", err)
	}

	// Get the module name and version
	modInfoCmd := exec.Command("go", "list", "-m", "-json")
	modInfoCmd.Dir = tmpDir
	modInfoOutput, err := modInfoCmd.Output()
	if err != nil {
		log.Fatal("Failed to get module info:", err)
	} else {
		fmt.Println("modInfoOutput:", string(modInfoOutput))
	}

	var data map[string]interface{}

	err = json.Unmarshal(modInfoOutput, &data)
	if err != nil {
		fmt.Println(err)
	}

	rootArtifact := &Artifact{
		Name:         data["Path"].(string),
		Dependencies: []*Artifact{},
	}

	// Generate the dependency tree
	// err = generateDependencyTree(tmpDir, rootArtifact)
	// if err != nil {
	// 	log.Fatal("Failed to generate dependency tree:", err)
	// }

	depMap, err := generateDependencyMap(tmpDir)
	if err != nil {
		log.Fatal("Failed to generate dependency tree:", err)
	}
	for _, dependency := range depMap[rootArtifact.Name].Dependencies {
		dependencyName := strings.Split(dependency, "@")[0]
		dependencyVersion := strings.Split(dependency, "@")[1]
		dependencyArtifact := &Artifact{
			Name:         dependencyName,
			Version:      dependencyVersion,
			Dependencies: []*Artifact{},
		}
		AppendDependncies(depMap, dependencyArtifact)
		rootArtifact.Dependencies = append(rootArtifact.Dependencies, dependencyArtifact)
	}
	// Convert the dependency tree to JSON
	jsonData, err := json.MarshalIndent([]*Artifact{rootArtifact}, "", "  ")
	if err != nil {
		log.Fatal("Failed to marshal JSON:", err)
	}

	// Write the JSON data to a file
	outputFile := "dependency_tree.json" // Replace with the desired output file name
	err = ioutil.WriteFile(outputFile, jsonData, 0644)
	if err != nil {
		log.Fatal("Failed to write JSON file:", err)
	}

	fmt.Println("Dependency tree generated successfully. Output file:", outputFile)
}

// func generateDependencyTree(dir string, parentArtifact *Artifact) error {
// 	modGraphCmd := exec.Command("go", "mod", "graph")
// 	modGraphCmd.Dir = dir
// 	modGraphOutput, err := modGraphCmd.Output()
// 	if err != nil {
// 		return fmt.Errorf("failed to run 'go mod graph': %w", err)
// 	}

// 	modules := strings.Split(string(modGraphOutput), "\n")
// 	// remove the last empty line
// 	modules = modules[:len(modules)-1]
// 	for _, module := range modules {
// 		dependency := strings.Split(module, " ")
// 		dependencyInfo := strings.Split(dependency[1], "@")

// 		dependencyName := dependencyInfo[0]
// 		dependencyVersion := dependencyInfo[1]

// 		dependencyArtifact := &Artifact{
// 			Name:    dependencyName,
// 			Version: dependencyVersion,
// 		}

// 		parentArtifact.Dependencies = append(parentArtifact.Dependencies, dependencyArtifact)

// 		// Recursively generate the dependency tree for the current dependency
// 		dependencyDir := filepath.Join(dir, "vendor", dependencyName)
// 		if dependencyDir == "github.com/redhat-appstudio/sprayproxy" {
// 			continue
// 		}
// 		if _, err := os.Stat(dependencyDir); os.IsNotExist(err) {
// 			// If the dependency is not in the vendor directory, recursively generate its dependency tree
// 			err = generateDependencyTree(dependencyDir, dependencyArtifact)
// 			if err != nil {
// 				return err
// 			}
// 		}
// 	}

// 	return nil
// }

type ArtifactMap struct {
	Dependencies []string
	// visited with default value false
	Visited bool
}

func generateDependencyMap(repoDir string) (map[string]ArtifactMap, error) {
	modGraphCmd := exec.Command("go", "mod", "graph")
	modGraphCmd.Dir = repoDir
	modGraphOutput, err := modGraphCmd.Output()
	if err != nil {
		return map[string]ArtifactMap{}, fmt.Errorf("failed to run 'go mod graph': %w", err)
	}
	modules := strings.Split(string(modGraphOutput), "\n")
	// remove the last empty line
	modules = modules[:len(modules)-1]

	// create a map of parent artifact and its dependencies to track the dependencies
	// of each artifact
	artifactMap := make(map[string]ArtifactMap)

	for _, module := range modules {
		dependentModule := strings.Split(module, " ")[0]
		// dependentModuleInfo := strings.Split(dependentModule, "@")
		// dependentModuleName := dependentModuleInfo[0]
		// dependentModuleVersion := dependentModuleInfo[1]
		// dependentArtifact := &Artifact{
		// 	Name:    dependentModuleName,
		// 	Version: dependentModuleVersion,
		// }

		dependencyModule := strings.Split(module, " ")[1]
		// dependencyModuleInfo := strings.Split(dependencyModule, "@")
		// dependencyModuleName := dependencyModuleInfo[0]
		// dependencyModuleVersion := dependencyModuleInfo[1]
		// dependencyArtifact := &Artifact{
		// 	Name:    dependencyModuleName,
		// 	Version: dependencyModuleVersion,
		// }

		// check if dependentModule is present in the map
		_, ok := artifactMap[dependentModule]
		if !ok {
			artifactMap[dependentModule] = ArtifactMap{[]string{dependencyModule}, false}
		} else {
			// append teh depndeencyModule to the existing list of Dependencies
			dependencies := artifactMap[dependentModule].Dependencies
			dependencies = append(dependencies, dependencyModule)
			artifactMap[dependentModule] = ArtifactMap{dependencies, false}
		}

		// dependentArtifact.Dependencies = append(dependentArtifact.Dependencies, dependencyArtifact)
	}

	return artifactMap, nil
}

func AppendDependncies(artifactMap map[string]ArtifactMap, artifact *Artifact) {
	// check if artifactMap key is present
	// if not present, return
	_, ok := artifactMap[artifact.Name+"@"+artifact.Version]
	if !ok {
		return
	} else if artifactMap[artifact.Name+"@"+artifact.Version].Visited {
		return
	} else {
		artifactMap[artifact.Name+"@"+artifact.Version] = ArtifactMap{artifactMap[artifact.Name+"@"+artifact.Version].Dependencies, true}
	}

	dependencies := artifactMap[artifact.Name+"@"+artifact.Version].Dependencies
	for _, dependency := range dependencies {
		dependencyName := strings.Split(dependency, "@")[0]
		dependencyVersion := strings.Split(dependency, "@")[1]
		dependencyArtifact := &Artifact{
			Name:         dependencyName,
			Version:      dependencyVersion,
			Dependencies: []*Artifact{},
		}
		artifact.Dependencies = append(artifact.Dependencies, dependencyArtifact)
	}
}
