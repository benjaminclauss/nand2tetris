package command

import "github.com/spf13/cobra"

var RootCmd = &cobra.Command{Use: "nand2tetris"}

func init() {
	RootCmd.AddCommand(assemblerCmd)
	RootCmd.AddCommand(vmTranslatorCommand)
}
