package main

import (
	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{Use: "dockcli"}
	rootCmd.AddCommand(listCmd, startCmd, stopCmd, statsCmd)
	rootCmd.Execute()
}
