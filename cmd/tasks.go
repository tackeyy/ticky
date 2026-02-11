package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/tackeyy/ticky/internal/ticktick"

	"github.com/spf13/cobra"
)

var tasksCmd = &cobra.Command{
	Use:   "tasks",
	Short: "Task operations",
}

var tasksListCmd = &cobra.Command{
	Use:   "list",
	Short: "List tasks in a project",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := ticktick.NewClient()
		if err != nil {
			return err
		}

		projectID, _ := cmd.Flags().GetString("project")
		if projectID == "" {
			projectID, err = findInboxID(client)
			if err != nil {
				return err
			}
		}

		pd, err := client.GetProjectData(projectID)
		if err != nil {
			return fmt.Errorf("failed to list tasks: %w", err)
		}

		if outputJSON {
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(pd.Tasks)
		}

		if outputPlain {
			for _, t := range pd.Tasks {
				tags := strings.Join(t.Tags, ",")
				fmt.Printf("%s\t%s\t%s\t%s\t%s\t%s\n",
					t.ID, t.ProjectID, t.Title,
					ticktick.PriorityString(t.Priority),
					t.DueDate, tags)
			}
			return nil
		}

		if len(pd.Tasks) == 0 {
			fmt.Println("No tasks found")
			return nil
		}
		for _, t := range pd.Tasks {
			priority := ""
			if t.Priority > 0 {
				priority = fmt.Sprintf(" [%s]", ticktick.PriorityString(t.Priority))
			}
			due := ""
			if t.DueDate != "" {
				due = fmt.Sprintf(" (due: %s)", t.DueDate[:10])
			}
			tags := ""
			if len(t.Tags) > 0 {
				tags = fmt.Sprintf(" #%s", strings.Join(t.Tags, " #"))
			}
			fmt.Printf("%-24s %s%s%s%s\n", t.ID, t.Title, priority, due, tags)
		}
		return nil
	},
}

var tasksGetCmd = &cobra.Command{
	Use:   "get <task_id>",
	Short: "Get task details",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := ticktick.NewClient()
		if err != nil {
			return err
		}

		projectID, _ := cmd.Flags().GetString("project")
		if projectID == "" {
			return fmt.Errorf("--project is required for get")
		}

		task, err := client.GetTask(projectID, args[0])
		if err != nil {
			return fmt.Errorf("failed to get task: %w", err)
		}

		if outputJSON {
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(task)
		}

		if outputPlain {
			tags := strings.Join(task.Tags, ",")
			fmt.Printf("%s\t%s\t%s\t%s\t%s\t%s\n",
				task.ID, task.ProjectID, task.Title,
				ticktick.PriorityString(task.Priority),
				task.DueDate, tags)
			return nil
		}

		fmt.Printf("ID:       %s\n", task.ID)
		fmt.Printf("Project:  %s\n", task.ProjectID)
		fmt.Printf("Title:    %s\n", task.Title)
		if task.Content != "" {
			fmt.Printf("Content:  %s\n", task.Content)
		}
		fmt.Printf("Priority: %s\n", ticktick.PriorityString(task.Priority))
		if task.DueDate != "" {
			fmt.Printf("Due:      %s\n", task.DueDate)
		}
		if len(task.Tags) > 0 {
			fmt.Printf("Tags:     %s\n", strings.Join(task.Tags, ", "))
		}
		return nil
	},
}

var tasksCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new task",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := ticktick.NewClient()
		if err != nil {
			return err
		}

		title, _ := cmd.Flags().GetString("title")
		if title == "" {
			return fmt.Errorf("--title is required")
		}

		projectID, _ := cmd.Flags().GetString("project")
		content, _ := cmd.Flags().GetString("content")
		priorityStr, _ := cmd.Flags().GetString("priority")
		dueStr, _ := cmd.Flags().GetString("due")
		tagsStr, _ := cmd.Flags().GetString("tags")

		req := &ticktick.TaskCreateRequest{
			Title:     title,
			ProjectID: projectID,
			Content:   content,
		}

		if priorityStr != "" {
			p, err := ticktick.ParsePriority(priorityStr)
			if err != nil {
				return err
			}
			req.Priority = p
		}

		if dueStr != "" {
			due, err := ticktick.ParseDate(dueStr)
			if err != nil {
				return err
			}
			req.DueDate = due
		}

		if tagsStr != "" {
			req.Tags = strings.Split(tagsStr, ",")
		}

		task, err := client.CreateTask(req)
		if err != nil {
			return fmt.Errorf("failed to create task: %w", err)
		}

		if outputJSON {
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(task)
		}

		if outputPlain {
			fmt.Printf("%s\t%s\n", task.ID, task.ProjectID)
			return nil
		}

		fmt.Printf("Created task: %s (ID: %s, Project: %s)\n", task.Title, task.ID, task.ProjectID)
		return nil
	},
}

var tasksUpdateCmd = &cobra.Command{
	Use:   "update <task_id>",
	Short: "Update an existing task",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := ticktick.NewClient()
		if err != nil {
			return err
		}

		projectID, _ := cmd.Flags().GetString("project")
		if projectID == "" {
			return fmt.Errorf("--project is required for update")
		}

		// Fetch existing task first
		existing, err := client.GetTask(projectID, args[0])
		if err != nil {
			return fmt.Errorf("failed to get existing task: %w", err)
		}

		req := &ticktick.TaskUpdateRequest{
			ID:        args[0],
			ProjectID: projectID,
			Title:     existing.Title,
			Content:   existing.Content,
			Tags:      existing.Tags,
		}

		p := existing.Priority
		req.Priority = &p

		if existing.DueDate != "" {
			d := existing.DueDate
			req.DueDate = &d
		}

		if cmd.Flags().Changed("title") {
			title, _ := cmd.Flags().GetString("title")
			req.Title = title
		}
		if cmd.Flags().Changed("content") {
			content, _ := cmd.Flags().GetString("content")
			req.Content = content
		}
		if cmd.Flags().Changed("priority") {
			priorityStr, _ := cmd.Flags().GetString("priority")
			pVal, err := ticktick.ParsePriority(priorityStr)
			if err != nil {
				return err
			}
			req.Priority = &pVal
		}
		if cmd.Flags().Changed("due") {
			dueStr, _ := cmd.Flags().GetString("due")
			due, err := ticktick.ParseDate(dueStr)
			if err != nil {
				return err
			}
			req.DueDate = &due
		}
		if cmd.Flags().Changed("clear-due") {
			clearDue, _ := cmd.Flags().GetBool("clear-due")
			if clearDue {
				empty := ""
				req.DueDate = &empty
			}
		}
		if cmd.Flags().Changed("tags") {
			tagsStr, _ := cmd.Flags().GetString("tags")
			if tagsStr == "" {
				req.Tags = []string{}
			} else {
				req.Tags = strings.Split(tagsStr, ",")
			}
		}
		if cmd.Flags().Changed("add-tags") {
			addTagsStr, _ := cmd.Flags().GetString("add-tags")
			for _, tag := range strings.Split(addTagsStr, ",") {
				tag = strings.TrimSpace(tag)
				if tag != "" && !containsStr(req.Tags, tag) {
					req.Tags = append(req.Tags, tag)
				}
			}
		}
		if cmd.Flags().Changed("remove-tags") {
			removeTagsStr, _ := cmd.Flags().GetString("remove-tags")
			removeTags := strings.Split(removeTagsStr, ",")
			var filtered []string
			for _, t := range req.Tags {
				if !containsStr(removeTags, t) {
					filtered = append(filtered, t)
				}
			}
			req.Tags = filtered
		}

		task, err := client.UpdateTask(req)
		if err != nil {
			return fmt.Errorf("failed to update task: %w", err)
		}

		if outputJSON {
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(task)
		}

		if outputPlain {
			fmt.Printf("%s\t%s\n", task.ID, task.ProjectID)
			return nil
		}

		fmt.Printf("Updated task: %s (ID: %s)\n", task.Title, task.ID)
		return nil
	},
}

var tasksCompleteCmd = &cobra.Command{
	Use:   "complete <task_id>",
	Short: "Mark a task as complete",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := ticktick.NewClient()
		if err != nil {
			return err
		}

		projectID, _ := cmd.Flags().GetString("project")
		if projectID == "" {
			return fmt.Errorf("--project is required for complete")
		}

		if err := client.CompleteTask(projectID, args[0]); err != nil {
			return fmt.Errorf("failed to complete task: %w", err)
		}

		if outputJSON {
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(map[string]string{
				"status":  "completed",
				"task_id": args[0],
			})
		}

		if outputPlain {
			fmt.Printf("%s\tcompleted\n", args[0])
			return nil
		}

		fmt.Printf("Task %s completed\n", args[0])
		return nil
	},
}

var tasksDeleteCmd = &cobra.Command{
	Use:   "delete <task_id>",
	Short: "Delete a task",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := ticktick.NewClient()
		if err != nil {
			return err
		}

		projectID, _ := cmd.Flags().GetString("project")
		if projectID == "" {
			return fmt.Errorf("--project is required for delete")
		}

		if err := client.DeleteTask(projectID, args[0]); err != nil {
			return fmt.Errorf("failed to delete task: %w", err)
		}

		if outputJSON {
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(map[string]string{
				"status":  "deleted",
				"task_id": args[0],
			})
		}

		if outputPlain {
			fmt.Printf("%s\tdeleted\n", args[0])
			return nil
		}

		fmt.Printf("Task %s deleted\n", args[0])
		return nil
	},
}

func init() {
	tasksListCmd.Flags().String("project", "", "Project ID (default: Inbox)")

	tasksGetCmd.Flags().String("project", "", "Project ID (required)")

	tasksCreateCmd.Flags().String("title", "", "Task title (required)")
	tasksCreateCmd.Flags().String("project", "", "Project ID (default: Inbox)")
	tasksCreateCmd.Flags().String("content", "", "Task content/description")
	tasksCreateCmd.Flags().String("priority", "", "Priority: none, low, medium, high")
	tasksCreateCmd.Flags().String("due", "", "Due date: today, tomorrow, +3d, YYYY-MM-DD")
	tasksCreateCmd.Flags().String("tags", "", "Comma-separated tags")

	tasksUpdateCmd.Flags().String("project", "", "Project ID (required)")
	tasksUpdateCmd.Flags().String("title", "", "New title")
	tasksUpdateCmd.Flags().String("content", "", "New content")
	tasksUpdateCmd.Flags().String("priority", "", "Priority: none, low, medium, high")
	tasksUpdateCmd.Flags().String("due", "", "Due date: today, tomorrow, +3d, YYYY-MM-DD")
	tasksUpdateCmd.Flags().Bool("clear-due", false, "Clear the due date")
	tasksUpdateCmd.Flags().String("tags", "", "Replace all tags (comma-separated)")
	tasksUpdateCmd.Flags().String("add-tags", "", "Add tags (comma-separated)")
	tasksUpdateCmd.Flags().String("remove-tags", "", "Remove tags (comma-separated)")

	tasksCompleteCmd.Flags().String("project", "", "Project ID (required)")

	tasksDeleteCmd.Flags().String("project", "", "Project ID (required)")

	tasksCmd.AddCommand(tasksListCmd)
	tasksCmd.AddCommand(tasksGetCmd)
	tasksCmd.AddCommand(tasksCreateCmd)
	tasksCmd.AddCommand(tasksUpdateCmd)
	tasksCmd.AddCommand(tasksCompleteCmd)
	tasksCmd.AddCommand(tasksDeleteCmd)
	rootCmd.AddCommand(tasksCmd)
}

// findInboxID finds the Inbox project ID.
func findInboxID(client *ticktick.Client) (string, error) {
	projects, err := client.GetProjects()
	if err != nil {
		return "", fmt.Errorf("failed to list projects: %w", err)
	}
	for _, p := range projects {
		if p.Kind == "INBOX" || p.Name == "Inbox" {
			return p.ID, nil
		}
	}
	if len(projects) > 0 {
		return projects[0].ID, nil
	}
	return "", fmt.Errorf("no projects found")
}

func containsStr(slice []string, s string) bool {
	for _, v := range slice {
		if strings.TrimSpace(v) == strings.TrimSpace(s) {
			return true
		}
	}
	return false
}
