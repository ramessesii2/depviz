package generator

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/ramessesii2/depviz/pkg/artifact"
)

func GenerateJson(rootArtifact *artifact.Artifact, outFileDir string) {
	// Convert the dependency tree to JSON
	jsonData, err := json.MarshalIndent([]*artifact.Artifact{rootArtifact}, "", "  ")
	if err != nil {
		log.Fatal("Failed to marshal JSON:", err)
	}

	// Write the JSON data to a file
	if !strings.HasSuffix(outFileDir, "/") {
		outFileDir = outFileDir + "/"
	}
	outputFile := outFileDir + "dependency_tree.json"
	err = ioutil.WriteFile(outputFile, jsonData, 0644)
	if err != nil {
		log.Fatal("Failed to write JSON file:", err)
	}

	fmt.Println("Dependency Tree generated successfully. Output file:", outputFile)
}
