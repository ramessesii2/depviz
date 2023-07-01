package generator

import (
	"encoding/json"
	"os"

	"github.com/ramessesii2/depviz/pkg/artifact"
)

func GenerateJson(rootArtifact *artifact.Artifact, outputFileName string) error {
	// Convert the dependency tree to JSON
	jsonData, err := json.MarshalIndent([]*artifact.Artifact{rootArtifact}, "", "  ")
	if err != nil {
		return err
	}

	// Write the JSON data to a file
	err = os.WriteFile(outputFileName, jsonData, 0644)
	if err != nil {
		return err
	}

	return nil
}
