package cmd

import (
	"ci6ndex/internal"
	"fmt"

	"github.com/spf13/cobra"
)

var usersCmd = &cobra.Command{
	Use:   "users",
	Short: "ci6ndex USERS",
	Long:  `Command used to CRUD users.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("users called")
	},
}

var input string

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "ci6ndex USERS ADD",
	Long:  `Command used to add a user.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := internal.AddUsersFromFile(input)
		if err != nil {
			fmt.Println("Failed to add user: ", err)
		}
	},
}

func init() {
	addCmd.Flags().StringVarP(&input, "file", "f", "", "Input file (must be json)")

	rootCmd.AddCommand(usersCmd)
	usersCmd.AddCommand(addCmd)
}
