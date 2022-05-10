package cmd

import (
	"fmt"
	"github.com/MickMake/GoX32/Only"
	"github.com/MickMake/GoX32/defaults"
	"github.com/spf13/cobra"
)

func AttachRootCmd(cmd *cobra.Command) *cobra.Command {
	// ******************************************************************************** //
	var rootCmd = &cobra.Command{
		Use:              defaults.BinaryName,
		Short:            fmt.Sprintf("%s - Manage an Behringer X32  instance", defaults.BinaryName),
		Long:             fmt.Sprintf("%s - Manage an Behringer X32  instance", defaults.BinaryName),
		Run:              gbRootFunc,
		TraverseChildren: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// You can bind cobra and viper in a few locations, but PersistencePreRunE on the root command works well
			return initConfig(cmd)
		},
	}
	if cmd != nil {
		cmd.AddCommand(rootCmd)
	}
	rootCmd.Example = PrintExamples(rootCmd, "")

	rootCmd.SetHelpTemplate(DefaultHelpTemplate)
	rootCmd.SetUsageTemplate(DefaultUsageTemplate)
	rootCmd.SetVersionTemplate(DefaultVersionTemplate)

	rootCmd.PersistentFlags().StringVarP(&Cmd.X32Host, flagX32Host, "", defaultHost, fmt.Sprintf("Behringer X32: Host / IP address."))
	rootViper.SetDefault(flagX32Host, defaultHost)
	rootCmd.PersistentFlags().StringVarP(&Cmd.X32Port, flagX32Port, "", defaultPort, fmt.Sprintf("Behringer X32: Port."))
	rootViper.SetDefault(flagX32Port, defaultPort)
	rootCmd.PersistentFlags().StringVarP(&Cmd.X32Username, flagX32Username, "u", "", fmt.Sprintf("Behringer X32: Username."))
	rootViper.SetDefault(flagX32Username, "")
	rootCmd.PersistentFlags().StringVarP(&Cmd.X32Password, flagX32Password, "p", "", fmt.Sprintf("Behringer X32: Password."))
	rootViper.SetDefault(flagX32Password, "")
	rootCmd.PersistentFlags().DurationVarP(&Cmd.X32Timeout, flagX32Timeout, "", defaultTimeout, fmt.Sprintf("Behringer X32: Timeout."))
	rootViper.SetDefault(flagX32Timeout, defaultTimeout)

	rootCmd.PersistentFlags().StringVarP(&Cmd.MqttUsername, flagMqttUsername, "", "", fmt.Sprintf("HASSIO: mqtt username."))
	rootViper.SetDefault(flagMqttUsername, "")
	rootCmd.PersistentFlags().StringVarP(&Cmd.MqttPassword, flagMqttPassword, "", "", fmt.Sprintf("HASSIO: mqtt password."))
	rootViper.SetDefault(flagMqttPassword, "")
	rootCmd.PersistentFlags().StringVarP(&Cmd.MqttHost, flagMqttHost, "", "", fmt.Sprintf("HASSIO: mqtt host."))
	rootViper.SetDefault(flagMqttHost, "")
	rootCmd.PersistentFlags().StringVarP(&Cmd.MqttPort, flagMqttPort, "", "", fmt.Sprintf("HASSIO: mqtt port."))
	rootViper.SetDefault(flagMqttPort, "")

	// rootCmd.PersistentFlags().StringVarP(&Cmd.GoogleSheet, flagGoogleSheet, "", "", fmt.Sprintf("Google: Sheet URL for updates."))
	// rootViper.SetDefault(flagGoogleSheet, "")
	// rootCmd.PersistentFlags().BoolVarP(&Cmd.GoogleSheetUpdate, flagGoogleSheetUpdate, "", false, fmt.Sprintf("Update Google sheets."))
	// rootViper.SetDefault(flagGoogleSheetUpdate, false)
	// _ = rootCmd.PersistentFlags().MarkHidden(flagGoogleSheetUpdate)
	//
	// rootCmd.PersistentFlags().StringVarP(&Cmd.GitRepo, flagGitRepo, "", "", fmt.Sprintf("Git: Repo url for updates."))
	// rootViper.SetDefault(flagGitRepo, "")
	// rootCmd.PersistentFlags().StringVarP(&Cmd.GitRepoDir, flagGitRepoDir, "", "", fmt.Sprintf("Git: Local repo directory."))
	// rootViper.SetDefault(flagGitRepoDir, "")
	// rootCmd.PersistentFlags().StringVarP(&Cmd.GitUsername, flagGitUsername, "", "", fmt.Sprintf("Git: Repo username."))
	// rootViper.SetDefault(flagGitUsername, "")
	// rootCmd.PersistentFlags().StringVarP(&Cmd.GitPassword, flagGitPassword, "", "", fmt.Sprintf("Git: Repo password."))
	// rootViper.SetDefault(flagGitPassword, "")
	// rootCmd.PersistentFlags().StringVarP(&Cmd.GitKeyFile, flagGitKeyFile, "", "", fmt.Sprintf("Git: Repo SSH keyfile."))
	// rootViper.SetDefault(flagGitKeyFile, "")
	// rootCmd.PersistentFlags().StringVarP(&Cmd.GitToken, flagGitToken, "", "", fmt.Sprintf("Git: Repo token string."))
	// rootViper.SetDefault(flagGitToken, "")
	// rootCmd.PersistentFlags().StringVarP(&Cmd.GitDiffCmd, flagGitDiffCmd, "", "tkdiff", fmt.Sprintf("Git: Command for diffs."))
	// rootViper.SetDefault(flagGitDiffCmd, "tkdiff")

	rootCmd.PersistentFlags().StringVar(&Cmd.ConfigFile, flagConfigFile, Cmd.ConfigFile, fmt.Sprintf("%s: config file.", defaults.BinaryName))
	// _ = rootCmd.PersistentFlags().MarkHidden(flagConfigFile)
	rootCmd.PersistentFlags().BoolVarP(&Cmd.Debug, flagDebug, "", false, fmt.Sprintf("%s: Debug mode.", defaults.BinaryName))
	rootViper.SetDefault(flagDebug, false)
	rootCmd.PersistentFlags().BoolVarP(&Cmd.Quiet, flagQuiet, "q", false, fmt.Sprintf("%s: Silence all messages.", defaults.BinaryName))
	rootViper.SetDefault(flagQuiet, false)

	rootCmd.PersistentFlags().SortFlags = false
	rootCmd.Flags().SortFlags = false

	return rootCmd
}

func gbRootFunc(cmd *cobra.Command, args []string) {
	for range Only.Once {
		if len(args) == 0 {
			_ = cmd.Help()
			break
		}
	}
}
