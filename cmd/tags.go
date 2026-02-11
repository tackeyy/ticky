package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"

	"github.com/tackeyy/ticky/internal/ticktick"

	"github.com/spf13/cobra"
)

var tagsCmd = &cobra.Command{
	Use:   "tags",
	Short: "Tag operations",
}

var tagsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all tags (aggregated from tasks across all projects)",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := ticktick.NewClient()
		if err != nil {
			return err
		}

		projects, err := client.GetProjects()
		if err != nil {
			return fmt.Errorf("failed to list projects: %w", err)
		}

		tagCount := make(map[string]int)
		for _, p := range projects {
			pd, err := client.GetProjectData(p.ID)
			if err != nil {
				continue
			}
			for _, t := range pd.Tasks {
				for _, tag := range t.Tags {
					tagCount[tag]++
				}
			}
		}

		type tagInfo struct {
			Name  string `json:"name"`
			Count int    `json:"count"`
		}

		var tags []tagInfo
		for name, count := range tagCount {
			tags = append(tags, tagInfo{Name: name, Count: count})
		}
		sort.Slice(tags, func(i, j int) bool {
			return tags[i].Count > tags[j].Count
		})

		if outputJSON {
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(tags)
		}

		if outputPlain {
			for _, t := range tags {
				fmt.Printf("%s\t%d\n", t.Name, t.Count)
			}
			return nil
		}

		if len(tags) == 0 {
			fmt.Println("No tags found")
			return nil
		}
		for _, t := range tags {
			fmt.Printf("#%-20s (%d tasks)\n", t.Name, t.Count)
		}
		return nil
	},
}

func init() {
	tagsCmd.AddCommand(tagsListCmd)
	rootCmd.AddCommand(tagsCmd)
}
