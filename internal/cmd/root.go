/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"os"

	"github.com/csris/mdmail/internal/mdmail"
	"github.com/spf13/cobra"
)


const longDescription = `mdmail is a simple utility for writing email using Markdown.
Simply write your email in a markdown file with YAML frontmatter and
use mdmail to convert it to HTML and add it as a draft.

Example:

  $ export IMAP_SERVER=imap.gmail.com:993
  $ export IMAP_USERNAME=YOUR_EMAIL_ADDRESS
  $ export IMAP_PASSWORD=YOUR_PASSWORD

  $ mdmail message.md
`

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "mdmail [flags] message ...",
	Short: "Compose email using Markdown",
	Long: longDescription,
	Args: cobra.MinimumNArgs(1),
	Run: mdmail.CreateDraftFromMarkdown,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.mdmail.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
