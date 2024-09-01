package command

import "github.com/spf13/cobra"

func NewRootCommand() *cobra.Command {
	cmd := &cobra.Command{Use: "nand2tetris"}
	cmd.AddCommand(assemblerCmd)
	cmd.AddCommand(vmTranslatorCommand)

	return cmd
}
