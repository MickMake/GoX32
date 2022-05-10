package cmd

import (
	"fmt"
	"github.com/MickMake/GoX32/Only"
	"github.com/MickMake/GoX32/defaults"
	"github.com/spf13/cobra"
)


func AttachCmdVersion(cmd *cobra.Command) *cobra.Command {
	// ******************************************************************************** //
	var cmdVersion = &cobra.Command{
		Use:                   "version",
		Short:                 "Version info.",
		Long:                  "Version info.",
		DisableFlagParsing:    false,
		DisableFlagsInUseLine: false,
		PreRunE:               Cmd.ProcessArgs,
		Run:                   cmdVersionFunc,
		Args:                  cobra.RangeArgs(0, 0),
	}
	cmd.AddCommand(cmdVersion)
	cmdVersion.Example = PrintExamples(cmdVersion, "")

	return cmdVersion
}


// ******************************************************************************** //
//goland:noinspection GoUnusedParameter
func cmdVersionFunc(cmd *cobra.Command, args []string) {
	for range Only.Once {
		fmt.Printf("%s v%s\n", defaults.BinaryName, defaults.Version)
	}
}
