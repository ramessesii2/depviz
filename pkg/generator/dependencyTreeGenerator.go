package generator

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/ramessesii2/depviz/pkg/artifact"
)

// ArtifactMap helps to track the dependencies of each artifact and if they are visited
type ArtifactMap struct {
	Dependencies []string
	Visited      bool
}

// GenerateDependencyMap generates a map having key as parent artifact and value as its dependencies using ArtifacMap
// This is helpful to track the dependencies of each artifact
func GenerateDependencyMap(repoDir string) (map[string]ArtifactMap, error) {
	// get the dependency text file using go mod graph
	goModGraphCmd := exec.Command("go", "mod", "graph")
	goModGraphCmd.Dir = repoDir
	modGraphOutput, err := goModGraphCmd.Output()
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
		dependencyModule := strings.Split(module, " ")[1]

		// check if the map with dependentModule as key is present
		_, ok := artifactMap[dependentModule]
		if !ok {
			artifactMap[dependentModule] = ArtifactMap{[]string{dependencyModule}, false}
		} else {
			// append the depndeencyModule to the existing list of Dependencies
			dependencies := artifactMap[dependentModule].Dependencies
			dependencies = append(dependencies, dependencyModule)
			artifactMap[dependentModule] = ArtifactMap{dependencies, false}
		}
	}

	return artifactMap, nil
}

// AppendDependencies appends the dependencies of each artifact using ArtifactMap
func AppendDependncies(artifactMap map[string]ArtifactMap, parentArtifact *artifact.Artifact) {
	// check if artifactMap key is present
	parentArtifactNameVersion := parentArtifact.Name + "@" + parentArtifact.Version
	_, ok := artifactMap[parentArtifactNameVersion]
	if !ok {
		return
	} else if artifactMap[parentArtifactNameVersion].Visited {
		return
	} else {
		artifactMap[parentArtifactNameVersion] = ArtifactMap{artifactMap[parentArtifactNameVersion].Dependencies, true}
	}

	dependencies := artifactMap[parentArtifactNameVersion].Dependencies
	for _, dependency := range dependencies {
		dependencyName := strings.Split(dependency, "@")[0]
		dependencyVersion := strings.Split(dependency, "@")[1]
		dependencyArtifact := &artifact.Artifact{
			Name:         dependencyName,
			Version:      dependencyVersion,
			Dependencies: []*artifact.Artifact{},
		}
		parentArtifact.Dependencies = append(parentArtifact.Dependencies, dependencyArtifact)
	}
}
