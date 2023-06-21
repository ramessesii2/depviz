package dependecytree

import "github.com/ramessesii2/depviz/pkg/artifact"

// GenerateDependencyTree generates a dependency tree from a given package
// It uses the `go mod graph` command to get the dependency graph
func GenerateDependencyTree(pkg string, parentArtifact *artifact.Artifact) error{
}

