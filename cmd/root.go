package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use: "ek",
	Short: "k8s终端工具",
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
	},
}

func Execute() {
	if err := rootCmd.Execute(); err !=nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}