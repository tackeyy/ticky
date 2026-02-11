package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/tackeyy/ticky/internal/ticktick"

	"github.com/spf13/cobra"
)

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Authentication commands",
}

var authLoginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to TickTick via OAuth",
	RunE: func(cmd *cobra.Command, args []string) error {
		token, err := ticktick.OAuthLogin()
		if err != nil {
			return fmt.Errorf("login failed: %w", err)
		}

		if err := ticktick.SaveToken(token); err != nil {
			return fmt.Errorf("failed to save token: %w", err)
		}

		fmt.Println("Successfully authenticated!")
		fmt.Printf("Token saved to %s\n", ticktick.TokenPath())
		return nil
	},
}

var authStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show authentication status",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := ticktick.NewClient()
		if err != nil {
			if outputJSON {
				enc := json.NewEncoder(os.Stdout)
				enc.SetIndent("", "  ")
				return enc.Encode(map[string]string{"status": "not authenticated", "error": err.Error()})
			}
			fmt.Println("Not authenticated")
			return nil
		}

		user, err := client.GetUser()
		if err != nil {
			return fmt.Errorf("failed to get user info: %w", err)
		}

		if outputJSON {
			out := map[string]string{
				"status":   "authenticated",
				"id":       user.ID,
				"username": user.Username,
				"name":     user.Name,
			}
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(out)
		}

		if outputPlain {
			fmt.Printf("%s\t%s\t%s\n", user.ID, user.Username, user.Name)
			return nil
		}

		fmt.Printf("Authenticated as: %s (%s)\n", user.Name, user.Username)
		fmt.Printf("User ID: %s\n", user.ID)
		return nil
	},
}

var authLogoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Remove saved authentication token",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := ticktick.DeleteToken(); err != nil {
			return fmt.Errorf("logout failed: %w", err)
		}
		fmt.Println("Logged out successfully")
		return nil
	},
}

func init() {
	authCmd.AddCommand(authLoginCmd)
	authCmd.AddCommand(authStatusCmd)
	authCmd.AddCommand(authLogoutCmd)
	rootCmd.AddCommand(authCmd)
}
