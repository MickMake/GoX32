package cmd


// func AttachCmdDataTemplatePoints(cmd *cobra.Command) *cobra.Command {
// 	// ********************************************************************************
// 	var c = &cobra.Command{
// 		Use:                   "template-points <template_id>",
// 		Aliases:               []string{""},
// 		Short:                 fmt.Sprintf("List data points used in report template."),
// 		Long:                  fmt.Sprintf("List data points used in report template."),
// 		DisableFlagParsing:    false,
// 		DisableFlagsInUseLine: false,
// 		PreRunE:               Cmd.X32Args,
// 		RunE:                  func(cmd *cobra.Command, args []string) error {
// 			_ = SwitchOutput(cmd)
// 			return Cmd.X32.GetTemplatePoints(args[0])
// 		},
// 		Args:                  cobra.ExactArgs(1),
// 	}
// 	cmd.AddCommand(c)
// 	c.Example = PrintExamples(c, "8042", "8040")
//
// 	return cmd
// }
//
// func AttachCmdDataTemplate(cmd *cobra.Command) *cobra.Command {
// 	// ********************************************************************************
// 	var c = &cobra.Command{
// 		Use:                   "template <template_id> <date> [filter]",
// 		Short:                 fmt.Sprintf("Get data from report template."),
// 		Long:                  fmt.Sprintf("Get data from report template."),
// 		DisableFlagParsing:    false,
// 		DisableFlagsInUseLine: false,
// 		PreRunE:               Cmd.X32Args,
// 		RunE:                  func(cmd *cobra.Command, args []string) error {
// 			_ = SwitchOutput(cmd)
// 			args = fillArray(3, args)
// 			return Cmd.X32.GetTemplateData(args[0], args[1], args[2])
// 		},
// 		Args:                  cobra.RangeArgs(2, 3),
// 	}
// 	cmd.AddCommand(c)
// 	c.Example = PrintExamples(c, "8042 20220212", "8042 20220212 '{\"search_string\":\"p83106\",\"min_left_axis\":-6000,\"max_left_axis\":12000}'")
//
// 	return cmd
// }
//
// func AttachCmdDataStats(cmd *cobra.Command) *cobra.Command {
// 	// ********************************************************************************
// 	var c = &cobra.Command{
// 		Use:                   "stats",
// 		Aliases:               []string{""},
// 		Short:                 fmt.Sprintf("Get current inverter stats, (last 5 minutes)."),
// 		Long:                  fmt.Sprintf("Get current inverter stats, (last 5 minutes)."),
// 		DisableFlagParsing:    false,
// 		DisableFlagsInUseLine: false,
// 		PreRunE:               Cmd.X32Args,
// 		RunE:                  func(cmd *cobra.Command, args []string) error {
// 			_ = SwitchOutput(cmd)
// 			return Cmd.X32.PrintCurrentStats()
// 		},
// 		Args:                  cobra.ExactArgs(0),
// 	}
// 	cmd.AddCommand(c)
// 	c.Example = PrintExamples(c, "")
//
// 	return cmd
// }
//
// func AttachCmdDataPoints(cmd *cobra.Command) *cobra.Command {
// 	// ********************************************************************************
// 	var cmdDataPoints = &cobra.Command{
// 		Use:                   "points <date> <device_id.point_id> ...",
// 		Aliases:               []string{""},
// 		Short:                 fmt.Sprintf("Get points data for a specific date."),
// 		Long:                  fmt.Sprintf("Get points data for a specific date."),
// 		DisableFlagParsing:    false,
// 		DisableFlagsInUseLine: false,
// 		PreRunE:               Cmd.X32Args,
// 		RunE:                  func(cmd *cobra.Command, args []string) error {
// 			_ = SwitchOutput(cmd)
// 			return Cmd.X32.GetPointData(args[0], api.CreatePoints(args[1:]))
// 		},
// 		Args:                  cobra.MinimumNArgs(2),
// 	}
// 	cmd.AddCommand(cmdDataPoints)
// 	cmdDataPoints.Example = PrintExamples(cmdDataPoints, "20220202 1129147.p13019 1129147.p83106")
//
// 	return cmd
// }
//
// func AttachCmdDataTemplates(cmd *cobra.Command) *cobra.Command {
// 	// ********************************************************************************
// 	var c = &cobra.Command{
// 		Use:                   "templates",
// 		Aliases:               []string{""},
// 		Short:                 fmt.Sprintf("Get all defined templates."),
// 		Long:                  fmt.Sprintf("Get all defined templates."),
// 		DisableFlagParsing:    false,
// 		DisableFlagsInUseLine: false,
// 		PreRunE:               Cmd.X32Args,
// 		RunE:                  func(cmd *cobra.Command, args []string) error {
// 			_ = SwitchOutput(cmd)
// 			return Cmd.X32.GetTemplates()
// 		},
// 		Args:                  cobra.ExactArgs(0),
// 	}
// 	cmd.AddCommand(c)
// 	c.Example = PrintExamples(c, "")
//
// 	return cmd
// }
//
// func AttachCmdDataMqtt(cmd *cobra.Command) *cobra.Command {
// 	// ********************************************************************************
// 	var c = &cobra.Command{
// 		Use:                   "mqtt-server",
// 		Aliases:               []string{"mqtt"},
// 		Short:                 fmt.Sprintf("Get Behringer MQTT service login details."),
// 		Long:                  fmt.Sprintf("Get Behringer MQTT service login details."),
// 		DisableFlagParsing:    false,
// 		DisableFlagsInUseLine: false,
// 		PreRunE:               Cmd.X32Args,
// 		RunE:                  func(cmd *cobra.Command, args []string) error {
// 			_ = SwitchOutput(cmd)
// 			args = fillArray(1, args)
// 			return Cmd.X32.GetIsolarcloudMqtt(args[0])
// 		},
// 		Args:                  cobra.RangeArgs(0, 1),
// 	}
// 	cmd.AddCommand(c)
// 	c.Example = PrintExamples(c, "")
//
// 	return cmd
// }
//
// func AttachCmdDataRealTime(cmd *cobra.Command) *cobra.Command {
// 	// ********************************************************************************
// 	var c = &cobra.Command{
// 		Use:                   "real-time",
// 		Aliases:               []string{"realtime"},
// 		Short:                 fmt.Sprintf("Get Behringer real-time data."),
// 		Long:                  fmt.Sprintf("Get Behringer real-time data."),
// 		DisableFlagParsing:    false,
// 		DisableFlagsInUseLine: false,
// 		PreRunE:               Cmd.X32Args,
// 		RunE:                  func(cmd *cobra.Command, args []string) error {
// 			_ = SwitchOutput(cmd)
// 			args = fillArray(1, args)
// 			return Cmd.X32.GetRealTimeData(args[0])
// 		},
// 		Args:                  cobra.RangeArgs(0, 1),
// 	}
// 	cmd.AddCommand(c)
// 	c.Example = PrintExamples(c, "")
//
// 	return cmd
// }
//
// func AttachCmdDataPsDetails(cmd *cobra.Command) *cobra.Command {
// 	// ********************************************************************************
// 	var c = &cobra.Command{
// 		Use:                   "psdetails",
// 		Aliases:               []string{"ps-details"},
// 		Short:                 fmt.Sprintf("Get Behringer ps details."),
// 		Long:                  fmt.Sprintf("Get Behringer ps details."),
// 		DisableFlagParsing:    false,
// 		DisableFlagsInUseLine: false,
// 		PreRunE:               Cmd.X32Args,
// 		RunE:                  func(cmd *cobra.Command, args []string) error {
// 			_ = SwitchOutput(cmd)
// 			args = fillArray(1, args)
// 			return Cmd.X32.GetPsDetails(args[0])
// 		},
// 		Args:                  cobra.RangeArgs(0, 1),
// 	}
// 	cmd.AddCommand(c)
// 	c.Example = PrintExamples(c, "")
//
// 	return cmd
// }
//
// func AttachCmdDataPointNames(cmd *cobra.Command) *cobra.Command {
// 	// ********************************************************************************
// 	var c = &cobra.Command{
// 		Use:                   "point-names",
// 		Aliases:               []string{"names"},
// 		Short:                 fmt.Sprintf("Get Behringer point names."),
// 		Long:                  fmt.Sprintf("Get Behringer point names."),
// 		DisableFlagParsing:    false,
// 		DisableFlagsInUseLine: false,
// 		PreRunE:               Cmd.X32Args,
// 		RunE:                  func(cmd *cobra.Command, args []string) error {
// 			_ = SwitchOutput(cmd)
// 			return Cmd.X32.GetPointNames(args...)
// 		},
// 		Args:                  cobra.MinimumNArgs(0),
// 	}
// 	cmd.AddCommand(c)
// 	c.Example = PrintExamples(c, "")
//
// 	return cmd
// }
//
// // AttachCmdDataSearchPointInfo
// func AttachCmdDataSearchPointInfo(cmd *cobra.Command) *cobra.Command {
// 	// ********************************************************************************
// 	var c = &cobra.Command{
// 		Use:                   "search-point-names",
// 		Aliases:               []string{"names"},
// 		Short:                 fmt.Sprintf("Get Behringer search point names."),
// 		Long:                  fmt.Sprintf("Get Behringer search point names."),
// 		DisableFlagParsing:    false,
// 		DisableFlagsInUseLine: false,
// 		PreRunE:               Cmd.X32Args,
// 		RunE:                  func(cmd *cobra.Command, args []string) error {
// 			_ = SwitchOutput(cmd)
// 			return Cmd.X32.SearchPointNames(args...)
// 		},
// 		Args:                  cobra.MinimumNArgs(0),
// 	}
// 	cmd.AddCommand(c)
// 	c.Example = PrintExamples(c, "")
//
// 	return cmd
// }
