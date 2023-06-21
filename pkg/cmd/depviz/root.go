package depviz

import (
	"github.com/ramessesii2/depviz/cmd/app"
	"github.com/spf13/cobra"
)

func Root() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "depviz",
		Short: "depviz is a tool for generating dependency tree for a Go repo",
		Long:  "depviz is a tool for generating dependency tree for a Go repo",
		Run: func(cmd *cobra.Command, args []string) {
			repo, _ := cmd.Flags().GetString("repo")
			ref, _ := cmd.Flags().GetString("ref")

			app.App(repo, ref)
		},
	}

	cmd.Flags().StringP("repo", "r", "", "Fully qualified domain name of GitHub repository for e.g.- `https://www.github.com/ramessesii2/depviz`")
	cmd.Flags().StringP("ref", "", "main", "Branch or tag for e.g.- `main`")

	cmd.MarkFlagRequired("repo")

	return cmd
}
