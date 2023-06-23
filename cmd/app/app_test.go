package app

import (
	"testing"

	"github.com/ramessesii2/depviz/pkg/artifact"
	"github.com/ramessesii2/depviz/pkg/generator"
)

func TestDependencyGeneratorInitializer(t *testing.T) {
	type args struct {
		depMap       map[string]generator.ArtifactMap
		rootArtifact *artifact.Artifact
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Valid dependencies",
			args: args{
				depMap: map[string]generator.ArtifactMap{
					"pkg1": {
						Dependencies: []string{"pkg2@v1", "pkg3@v2", "pkg4@v3"},
						Visited:      false,
					},
					"pkg2@v1": {
						Dependencies: []string{"pkg9@v9", "pkg8@v10", "pkg7@v11"},
						Visited:      false,
					},
					"pkg4@v3": {
						Dependencies: []string{"pkg7@v3", "pkg8@v4", "pkg9@v5"},
						Visited:      false,
					},
				},
				rootArtifact: &artifact.Artifact{
					Name:    "pkg1",
					Version: "v1.0.0",
				},
			},
			wantErr: false,
		},
		{
			name: "Invalid dependencies due to missing version",
			args: args{
				depMap: map[string]generator.ArtifactMap{
					"pkg1": {
						Dependencies: []string{"pkg2", "pkg3", "pkg4@v3"},
						Visited:      false,
					},
					"pkg2": {
						Dependencies: []string{"pkg9@v9", "pkg8@v10", "pkg7@v11"},
						Visited:      false,
					},
					"pkg4@v3": {
						Dependencies: []string{"pkg7@v3", "pkg8@v4", "pkg9@v5"},
						Visited:      false,
					},
				},
				rootArtifact: &artifact.Artifact{
					Name:    "pkg1",
					Version: "v1.0.0",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := dependencyGeneratorInitializer(tt.args.depMap, tt.args.rootArtifact)
			if err != nil && !tt.wantErr {
				t.Errorf("dependencyGeneratorInitializer() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
