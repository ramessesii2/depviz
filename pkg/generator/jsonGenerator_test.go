package generator

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ramessesii2/depviz/pkg/artifact"
)

func TestGenerateJson(t *testing.T) {

	tests := []struct {
		name         string
		rootArtifact *artifact.Artifact
	}{
		{
			name: "test1",
			rootArtifact: &artifact.Artifact{
				Name:    "pkg1",
				Version: "v1.0.0",
				Dependencies: []*artifact.Artifact{
					{
						Name:    "pkg2",
						Version: "v1.0.0",
						Dependencies: []*artifact.Artifact{
							{
								Name:         "pkg3",
								Version:      "v1.0.0",
								Dependencies: []*artifact.Artifact{},
							},
							{
								Name:    "pkg8",
								Version: "v1.0.0",
								Dependencies: []*artifact.Artifact{
									{
										Name:    "pkg9",
										Version: "v1.0.0",
										Dependencies: []*artifact.Artifact{
											{
												Name:         "pkg10",
												Version:      "v1.0.0",
												Dependencies: []*artifact.Artifact{},
											},
										},
									},
								},
							},
							{
								Name:         "pkg13",
								Version:      "v1.0.0",
								Dependencies: []*artifact.Artifact{},
							},
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a temporary dir for testing
			tmpDir, err := ioutil.TempDir("", "test-output")
			if err != nil {
				t.Fatal("Failed to create temporary directory: ", err)
			}
			// defer os.RemoveAll(tmpDir)

			targetFileName := filepath.Join(tmpDir, "dependency_tree.json")
			err = GenerateJson(tt.rootArtifact, targetFileName)
			if err != nil {
				t.Errorf("Failed to generate JSON: %v", err)
			}

			// Check if the JSON file was generated correctly
			_, err = os.Stat(targetFileName)

			if os.IsNotExist(err) {
				t.Error("JSON file doesn't exist")
			}

			// Read the content of teh JSON file
			generatedJSON, err := ioutil.ReadFile(targetFileName)
			if err != nil {
				t.Fatal("Failed to read JSON file:", err)
			}

			expectedJSON := `[
				{
					"name": "pkg1",
					"version": "v1.0.0",
					"dependencies": [
					{
						"name": "pkg2",
						"version": "v1.0.0",
						"dependencies": [
						{
							"name": "pkg3",
							"version": "v1.0.0",
							"dependencies": []
						},
						{
							"name": "pkg8",
							"version": "v1.0.0",
							"dependencies": [
							{
								"name": "pkg9",
								"version": "v1.0.0",
								"dependencies": [
								{
									"name": "pkg10",
									"version": "v1.0.0",
									"dependencies": []
								}
								]
							}
							]
						},
						{
							"name": "pkg13",
							"version": "v1.0.0",
							"dependencies": []
						}
						]
					}
					]
				}
			]`

			// Remove all the whitespaces from the expected JSON
			expectedJSONStr := removeWhitespace(expectedJSON)
			generatedJSONStr := removeWhitespace(string(generatedJSON))

			if expectedJSONStr != generatedJSONStr {
				t.Errorf("Generated JSON does not match expected. \nExpected: %s \nGenerated: %s", expectedJSONStr, generatedJSONStr)
			}

		})
	}
}

func removeWhitespace(s string) string {
	return strings.Join(strings.Fields(s), "")
}
