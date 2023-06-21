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
