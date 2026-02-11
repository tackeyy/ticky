package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/tackeyy/ticky/internal/ticktick"

	"github.com/spf13/cobra"
)

var projectsCmd = &cobra.Command{
	Use:   "projects",
	Short: "Project operations",
}

var projectsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all projects",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := ticktick.NewClient()
		if err != nil {
			return err
		}

		projects, err := client.GetProjects()
		if err != nil {
			return fmt.Errorf("failed to list projects: %w", err)
		}

		if outputJSON {
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(projects)
		}

		if outputPlain {
			for _, p := range projects {
				fmt.Printf("%s\t%s\n", p.ID, p.Name)
			}
			return nil
		}

		for _, p := range projects {
			fmt.Printf("%-24s %s\n", p.ID, p.Name)
		}
		return nil
	},
}

var projectsGetCmd = &cobra.Command{
	Use:   "get <project_id>",
	Short: "Get project details",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := ticktick.NewClient()
		if err != nil {
			return err
		}

		project, err := client.GetProject(args[0])
		if err != nil {
			return fmt.Errorf("failed to get project: %w", err)
		}

		if outputJSON {
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(project)
		}

		if outputPlain {
			fmt.Printf("%s\t%s\n", project.ID, project.Name)
			return nil
		}

		fmt.Printf("ID:   %s\n", project.ID)
		fmt.Printf("Name: %s\n", project.Name)
		if project.Color != "" {
			fmt.Printf("Color: %s\n", project.Color)
		}
		return nil
	},
}

func init() {
	projectsCmd.AddCommand(projectsListCmd)
	projectsCmd.AddCommand(projectsGetCmd)
	rootCmd.AddCommand(projectsCmd)
}
