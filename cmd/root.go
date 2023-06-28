package cmd

import (
	run "gitget/pkg"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use: `gitget <orga>/<repo> <git-ref> <path> [flags]

  <orga> is the name of the organization
  <repo> is the name of the repository
  <git-ref> is the git reference, it can be a tag, branch, sha1 or a pull request reference
  <path> is the file path relative to the repository root`,
	Short:   "gitget downloads a git file from a repository",
	Version: "0.0.1",
	Args:    cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		cobra.CheckErr(run.GitGet(args[0], args[1], args[2]))
	},
	DisableFlagsInUseLine: true,
}

// Add flags to change the output path
func init() {
	rootCmd.PersistentFlags().StringVarP(&run.OutPutPath, "output", "o", ".", "output path, can be overwritten using GITHUB_SERVER_URL environment variable")
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}
