package cmd

import (
	"fmt"
	"github.com/MickMake/GoX32/Behringer"
	"github.com/MickMake/GoX32/Only"
	"github.com/spf13/cobra"
)


func AttachCmdX32(cmd *cobra.Command) *cobra.Command {
	// ******************************************************************************** //
	var cmdX32 = &cobra.Command{
		Use:                   "x32",
		Aliases:               []string{},
		Short:                 fmt.Sprintf("Interact directly with a Behringer X32."),
		Long:                  fmt.Sprintf("Interact directly with a Behringer X32."),
		DisableFlagParsing:    false,
		DisableFlagsInUseLine: false,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return Cmd.PreRun(cmd, args, Cmd.ProcessArgs)
		},
		Run:                   cmdX32Func,
		Args:                  cobra.MinimumNArgs(1),
	}
	cmd.AddCommand(cmdX32)
	cmdX32.Example = PrintExamples(cmdX32, "get <endpoint>", "put <endpoint>")

	cmdX32.PersistentFlags().StringVarP(&Cmd.X32OutputType, flagX32OutputType, "o", "", fmt.Sprintf("Output type: 'json', 'raw', 'file'"))
	_ = cmdX32.PersistentFlags().MarkHidden(flagX32OutputType)

	// ******************************************************************************** //
	var cmdX32List = &cobra.Command{
		Use:                   "ls",
		Aliases:               []string{"list"},
		Short:                 fmt.Sprintf("List Behringer api endpoints/areas"),
		Long:                  fmt.Sprintf("List Behringer api endpoints/areas"),
		DisableFlagParsing:    false,
		DisableFlagsInUseLine: false,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return Cmd.PreRun(cmd, args, Cmd.ProcessArgs, Cmd.X32Args)
		},
		Run:                   cmdX32ListFunc,
		Args:                  cobra.RangeArgs(0, 1),
	}
	cmdX32.AddCommand(cmdX32List)
	cmdX32List.Example = PrintExamples(cmdX32List, "", "areas", "endpoints", "<area name>")

	// ******************************************************************************** //
	var cmdX32Login = &cobra.Command{
		Use: "login",
		// Aliases:               []string{""},
		Short:                 fmt.Sprintf("Login to Behringer"),
		Long:                  fmt.Sprintf("Login to Behringer"),
		DisableFlagParsing:    false,
		DisableFlagsInUseLine: false,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return Cmd.PreRun(cmd, args, Cmd.ProcessArgs)
		},
		Run:                   cmdX32LoginFunc,
		Args:                  cobra.MinimumNArgs(0),
	}
	cmdX32.AddCommand(cmdX32Login)
	cmdX32Login.Example = PrintExamples(cmdX32Login, "")

	// ******************************************************************************** //
	var cmdX32Get = &cobra.Command{
		Use: "get",
		// Aliases:               []string{""},
		Short:                 fmt.Sprintf("Get details from Behringer"),
		Long:                  fmt.Sprintf("Get details from Behringer"),
		DisableFlagParsing:    false,
		DisableFlagsInUseLine: false,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return Cmd.PreRun(cmd, args, Cmd.ProcessArgs, Cmd.X32Args)
		},
		Run:                   cmdX32GetFunc,
		Args:                  cobra.MinimumNArgs(1),
	}
	cmdX32.AddCommand(cmdX32Get)
	cmdX32Get.Example = PrintExamples(cmdX32Get, "[area].<endpoint>")

	// ******************************************************************************** //
	var cmdX32Raw = &cobra.Command{
		Use: "raw",
		// Aliases:               []string{""},
		Short:                 fmt.Sprintf("Raw details from Behringer"),
		Long:                  fmt.Sprintf("Raw details from Behringer"),
		DisableFlagParsing:    false,
		DisableFlagsInUseLine: false,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return Cmd.PreRun(cmd, args, Cmd.ProcessArgs, Cmd.X32Args)
		},
		Run:                   cmdX32RawFunc,
		Args:                  cobra.MinimumNArgs(1),
	}
	cmdX32.AddCommand(cmdX32Raw)
	cmdX32Raw.Example = PrintExamples(cmdX32Raw, "[area].<endpoint>")

	// ******************************************************************************** //
	var cmdX32Save = &cobra.Command{
		Use: "save",
		// Aliases:               []string{""},
		Short:                 fmt.Sprintf("Save details from Behringer"),
		Long:                  fmt.Sprintf("Save details from Behringer"),
		DisableFlagParsing:    false,
		DisableFlagsInUseLine: false,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return Cmd.PreRun(cmd, args, Cmd.ProcessArgs, Cmd.X32Args)
		},
		Run:                   cmdX32SaveFunc,
		Args:                  cobra.MinimumNArgs(1),
	}
	cmdX32.AddCommand(cmdX32Save)
	cmdX32Save.Example = PrintExamples(cmdX32Save, "[area].<endpoint>")

	// ******************************************************************************** //
	var cmdX32Put = &cobra.Command{
		Use:                   "put",
		Aliases:               []string{"write"},
		Short:                 fmt.Sprintf("Put details onto Behringer"),
		Long:                  fmt.Sprintf("Put details onto Behringer"),
		DisableFlagParsing:    false,
		DisableFlagsInUseLine: false,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return Cmd.PreRun(cmd, args, Cmd.ProcessArgs, Cmd.X32Args)
		},
		Run:                   cmdX32PutFunc,
		Args:                  cobra.RangeArgs(0, 1),
	}
	cmdX32.AddCommand(cmdX32Put)
	cmdX32Put.Example = PrintExamples(cmdX32Put, "[area].<endpoint> <value>")

	return cmdX32
}

func cmdX32Func(cmd *cobra.Command, args []string) {
	for range Only.Once {
		if len(args) == 0 {
			Cmd.Error = cmd.Help()
			break
		}
	}
}

func cmdX32ListFunc(cmd *cobra.Command, args []string) {
	// for range Only.Once {
	// 	switch {
	// 	case len(args) == 0:
	// 		fmt.Println("Unknown sub-command.")
	// 		_ = cmd.Help()
	//
	// 	case args[0] == "endpoints":
	// 		Cmd.Error = Cmd.X32.ListEndpoints("")
	//
	// 	case args[0] == "areas":
	// 		Cmd.X32.ListAreas()
	//
	// 	default:
	// 		Cmd.Error = Cmd.X32.ListEndpoints(args[0])
	// 	}
	// }
}

func cmdX32LoginFunc(_ *cobra.Command, _ []string) {
	// for range Only.Once {
	// 	Cmd.Error = Cmd.X32.Login(login.X32Auth{
	// 		AppKey:       Cmd.X32AppKey,
	// 		UserAccount:  Cmd.X32Username,
	// 		UserPassword: Cmd.X32Password,
	// 		TokenFile:    Cmd.X32TokenFile,
	// 		Force:        true,
	// 	})
	// 	if Cmd.Error != nil {
	// 		break
	// 	}
	//
	// 	Cmd.X32.Auth.Print()
	//
	// 	if Cmd.X32.HasTokenChanged() {
	// 		Cmd.X32LastLogin = Cmd.X32.GetLastLogin()
	// 		Cmd.X32Token = Cmd.X32.GetToken()
	// 		Cmd.Error = writeConfig()
	// 	}
	// }
}

func cmdX32GetFunc(_ *cobra.Command, args []string) {
	// for range Only.Once {
	// 	Cmd.X32.OutputType.SetJson()
	//
	// 	args = fillArray(2, args)
	// 	if args[0] == "all" {
	// 		Cmd.Error = Cmd.X32.AllCritical()
	// 		break
	// 	}
	//
	// 	ep := Cmd.X32.GetByJson(args[0], args[1])
	// 	if Cmd.X32.Error != nil {
	// 		Cmd.Error = Cmd.X32.Error
	// 		break
	// 	}
	// 	if Cmd.Error != nil {
	// 		break
	// 	}
	//
	// 	Cmd.Error = ep.GetError()
	// 	if Cmd.Error != nil {
	// 		break
	// 	}
	// }
}

func cmdX32RawFunc(_ *cobra.Command, args []string) {
	// for range Only.Once {
	// 	Cmd.X32.OutputType.SetRaw()
	//
	// 	args = fillArray(2, args)
	// 	if args[0] == "all" {
	// 		Cmd.Error = Cmd.X32.AllCritical()
	// 		break
	// 	}
	//
	// 	ep := Cmd.X32.GetByJson(args[0], args[1])
	// 	if Cmd.X32.Error != nil {
	// 		Cmd.Error = Cmd.X32.Error
	// 		break
	// 	}
	// 	if Cmd.Error != nil {
	// 		break
	// 	}
	//
	// 	Cmd.Error = ep.GetError()
	// 	if Cmd.Error != nil {
	// 		break
	// 	}
	// }
}

func cmdX32SaveFunc(_ *cobra.Command, args []string) {
	// for range Only.Once {
	// 	Cmd.X32.OutputType.SetFile()
	//
	// 	args = fillArray(2, args)
	// 	if args[0] == "all" {
	// 		Cmd.Error = Cmd.X32.AllCritical()
	// 		break
	// 	}
	//
	// 	ep := Cmd.X32.GetByJson(args[0], args[1])
	// 	if Cmd.X32.Error != nil {
	// 		Cmd.Error = Cmd.X32.Error
	// 		break
	// 	}
	// 	if Cmd.Error != nil {
	// 		break
	// 	}
	//
	// 	Cmd.Error = ep.GetError()
	// 	if Cmd.Error != nil {
	// 		break
	// 	}
	// }
}

func cmdX32PutFunc(_ *cobra.Command, _ []string) {
	for range Only.Once {
		fmt.Println("Not yet implemented.")
		// args = fillArray(1, args)
		// Cmd.Error = X32.Init()
		// if Cmd.Error != nil {
		// 	break
		// }
	}
}


func SwitchOutput(cmd *cobra.Command) error {
	var err error
	for range Only.Once {
		foo := cmd.Parent()
		switch foo.Use {
			case "get":
				Cmd.X32.OutputType.SetHuman()
			case "raw":
				Cmd.X32.OutputType.SetRaw()
			case "save":
				Cmd.X32.OutputType.SetFile()
			case "graph":
				Cmd.X32.OutputType.SetGraph()
			default:
				Cmd.X32.OutputType.SetHuman()
		}
	}

	return err
}


func (ca *CommandArgs) X32Args(cmd *cobra.Command, args []string) error {
	for range Only.Once {
		ca.X32 = Behringer.NewX32(Behringer.ArgsX32 {
			Host:         ca.X32Host,
			Port:         ca.X32Port,
			ConfigDir:    ".",	// ca.ConfigDir,	// @TODO - DEBUG
			CacheDir:     ca.CacheDir,
			CacheTimeout: ca.X32Timeout,
		})
		if ca.X32.Error != nil {
			ca.Error = ca.X32.Error
			break
		}

		switch ca.X32OutputType {
			case "json":
				ca.X32.OutputType.SetJson()
			case "raw":
				ca.X32.OutputType.SetRaw()
			case "file":
				ca.X32.OutputType.SetFile()
			default:
				ca.X32.OutputType.SetJson()
		}

		LogPrintDate("Connecting to X32...\n")
		ca.Error = ca.X32.Connect()
		if ca.Error != nil {
			break
		}
		LogPrintDate("Found X32 device.\n%v\n", ca.X32.Info.String())

		// var id int64
		// id, ca.Error = ca.X32.GetPsId()
		// if ca.Error != nil {
		// 	break
		// }
		//
		// var model string
		// model, ca.Error = ca.X32.GetPsModel()
		// if ca.Error != nil {
		// 	break
		// }
		//
		// var serial string
		// serial, ca.Error = ca.X32.GetPsSerial()
		// if ca.Error != nil {
		// 	break
		// }
		// LogPrintDate("Found X32 device %s id:%d serial:%s\n", model, id, serial)
		//
		// ca.Error = ca.X32.Login(login.X32Auth{
		// 	AppKey:       ca.ApiAppKey,
		// 	UserAccount:  ca.ApiUsername,
		// 	UserPassword: ca.ApiPassword,
		// 	TokenFile:    ca.ApiTokenFile,
		// 	Force:        false,
		// })
		// if ca.Error != nil {
		// 	break
		// }
		//
		// if ca.Debug {
		// 	ca.X32.Auth.Print()
		// }
		//
		// if ca.X32.HasTokenChanged() {
		// 	ca.ApiLastLogin = ca.X32.GetLastLogin()
		// 	ca.ApiToken = ca.X32.GetToken()
		// 	ca.Error = writeConfig()
		// }

		ca.Valid = true
	}

	return ca.Error
}
