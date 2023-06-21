package generator

import (
	"testing"

	"github.com/ramessesii2/depviz/pkg/artifact"
)

// test for AppendDependencies()
func TestAppendDependencies(t *testing.T) {
	// create a map of parent artifact and its dependencies to track the dependencies
	// of each artifact
	artifactMap := make(map[string]ArtifactMap)
	artifactMap["pkg1@v1.0.0"] = ArtifactMap{[]string{"pkg2@v1.0.0"}, false}
	artifactMap["pkg2@v1.0.0"] = ArtifactMap{[]string{"pkg3@v1.0.0", "pkg8@v1.0.0", "pkg13@v1.0.0"}, false}
	artifactMap["pkg3@v1.0.0"] = ArtifactMap{[]string{}, false}
	artifactMap["pkg8@v1.0.0"] = ArtifactMap{[]string{"pkg9@v1.0.0"}, false}
	artifactMap["pkg9@v1.0.0"] = ArtifactMap{[]string{"pkg10@v1.0.0"}, false}
	artifactMap["pkg10@v1.0.0"] = ArtifactMap{[]string{}, false}
	artifactMap["pkg13@v1.0.0"] = ArtifactMap{[]string{}, false}

	parentArtifact := &artifact.Artifact{
		Name:         "pkg1",
		Version:      "v1.0.0",
		Dependencies: []*artifact.Artifact{},
	}

	AppendDependncies(artifactMap, parentArtifact)

	// check if the parentArtifact has the correct dependencies
	if len(parentArtifact.Dependencies) != 1 {
		t.Errorf("parentArtifact has incorrect number of dependencies, got: %d, want: %d", len(parentArtifact.Dependencies), 1)
	}
	if parentArtifact.Dependencies[0].Name != "pkg2" {
		t.Errorf("parentArtifact has incorrect dependency, got: %s, want: %s", parentArtifact.Dependencies[0].Name, "pkg2")
	}
	if parentArtifact.Dependencies[0].Version != "v1.0.0" {
		t.Errorf("parentArtifact has incorrect dependency, got: %s, want: %s", parentArtifact.Dependencies[0].Version, "v1.0.0")
	}
	if len(parentArtifact.Dependencies[0].Dependencies) != 3 {
		t.Errorf("parentArtifact has incorrect number of dependencies, got: %d, want :%d", len(parentArtifact.Dependencies[0].Dependencies), 3)
	}
}

func TestGenerateDependencyMap(t *testing.T) {
	modules := []string{"github.com/ramessesii2/depviz/pkg/generator github.com/ramessesii2/depviz/pkg/artifact", "github.com/ramessesii2/depviz/pkg/generator github.com/ramessesii2/depviz/pkg/util"}
	artifactMap := GenerateDependencyMap(modules)

	// check if the map is correctly generated
	if len(artifactMap) != 1 {
		t.Errorf("artifactMap has incorrect number of dependencies, got: %d, want: %d", len(artifactMap), 1)
	}
	if len(artifactMap["github.com/ramessesii2/depviz/pkg/generator"].Dependencies) != 2 {
		t.Errorf("artifactMap has incorrect number of dependencies, got: %d, want: %d", len(artifactMap["github.com/ramessesii2/depviz/pkg/generator"].Dependencies), 2)
	}
	if artifactMap["github.com/ramessesii2/depviz/pkg/generator"].Dependencies[0] != "github.com/ramessesii2/depviz/pkg/artifact" {
		t.Errorf("artifactMap has incorrect dependency, got: %s, want: %s", artifactMap["github.com/ramessesii2/depviz/pkg/generator"].Dependencies[0], "github.com/ramessesii2/depviz/pkg/artifact")
	}
	if artifactMap["github.com/ramessesii2/depviz/pkg/generator"].Dependencies[1] != "github.com/ramessesii2/depviz/pkg/util" {
		t.Errorf("artifactMap has incorrect dependency, got: %s, want: %s", artifactMap["github.com/ramessesii2/depviz/pkg/generator"].Dependencies[1], "github.com/ramessesii2/depviz/pkg/util")
	}
	if artifactMap["github.com/ramessesii2/depviz/pkg/generator"].Visited != false {
		t.Errorf("artifactMap has incorrect visited status, got: %t, want: %t", artifactMap["github.com/ramessesii2/depviz/pkg/generator"].Visited, false)
	}

}
