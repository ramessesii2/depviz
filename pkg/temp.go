
// func generateDependencyMap(repoDir, parentArtifactName string) (map[*Artifact][]Artifact, error) {
// 	modGraphCmd := exec.Command("go", "mod", "graph")
// 	modGraphCmd.Dir = repoDir
// 	modGraphOutput, err := modGraphCmd.Output()
// 	if err != nil {
// 		return map[*Artifact][]Artifact{}, fmt.Errorf("failed to run 'go mod graph': %w", err)
// 	}
// 	modules := strings.Split(string(modGraphOutput), "\n")
// 	// remove the last empty line
// 	modules = modules[:len(modules)-1]

// 	// create a map of parent artifact and its dependencies
// 	artifactMap := make(map[*Artifact][]Artifact)
// 	for _, module := range modules {
// 		dependentModule := strings.Split(module, " ")[0]
// 		dependentModuleInfo := strings.Split(dependentModule, "@")
// 		dependentModuleName := dependentModuleInfo[0]
// 		dependentModuleVersion := dependentModuleInfo[1]
// 		dependentArtifact := &Artifact{
// 			Name:    dependentModuleName,
// 			Version: dependentModuleVersion,
// 		}

// 		dependencyModule := strings.Split(module, " ")[1]
// 		dependencyModuleInfo := strings.Split(dependencyModule, "@")
// 		dependencyModuleName := dependencyModuleInfo[0]
// 		dependencyModuleVersion := dependencyModuleInfo[1]
// 		dependencyArtifact := &Artifact{
// 			Name:    dependencyModuleName,
// 			Version: dependencyModuleVersion,
// 		}

// 		dependentArtifact.Dependencies = append(dependentArtifact.Dependencies, *&dependencyArtifact)
// 	}

// 	return artifactMap, nil
// }
