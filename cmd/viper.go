package cmd

import (
	"fmt"
	"github.com/MickMake/GoX32/Only"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"os"
	"strings"
)

var rootViper *viper.Viper

// initConfig reads in config file and ENV variables if set.
func initConfig(cmd *cobra.Command) error {
	var err error

	for range Only.Once {
		rootViper = viper.New()
		rootViper.AddConfigPath(Cmd.ConfigDir)
		rootViper.SetConfigFile(Cmd.ConfigFile)
		// rootViper.SetConfigName("config")

		// If a config file is found, read it in.
		err = openConfig()
		if err != nil {
			break
		}

		rootViper.SetEnvPrefix(EnvPrefix)
		rootViper.AutomaticEnv() // read in environment variables that match
		err = bindFlags(cmd, rootViper)
		if err != nil {
			break
		}
	}

	return err
}

func openConfig() error {
	var err error

	for range Only.Once {
		err = rootViper.ReadInConfig()
		if _, ok := err.(viper.UnsupportedConfigError); ok {
			break
		}

		if _, ok := err.(viper.ConfigParseError); ok {
			break
		}

		if _, ok := err.(viper.ConfigMarshalError); ok {
			break
		}

		if os.IsNotExist(err) {
			rootViper.SetDefault(flagDebug, Cmd.Debug)
			rootViper.SetDefault(flagQuiet, Cmd.Quiet)

			rootViper.SetDefault(flagX32Host, Cmd.X32Host)
			rootViper.SetDefault(flagX32Port, Cmd.X32Port)
			rootViper.SetDefault(flagX32Timeout, Cmd.X32Timeout)
			rootViper.SetDefault(flagX32Username, defaultUsername)
			rootViper.SetDefault(flagX32Password, defaultPassword)

			rootViper.SetDefault(flagMqttUsername, Cmd.MqttUsername)
			rootViper.SetDefault(flagMqttPassword, Cmd.MqttPassword)
			rootViper.SetDefault(flagMqttHost, Cmd.MqttHost)
			rootViper.SetDefault(flagMqttPort, Cmd.MqttPort)

			rootViper.SetDefault(flagGoogleSheet, Cmd.GoogleSheet)
			rootViper.SetDefault(flagGoogleSheetUpdate, Cmd.GoogleSheetUpdate)

			rootViper.SetDefault(flagGitRepo, Cmd.GitRepo)
			rootViper.SetDefault(flagGitRepoDir, Cmd.GitRepoDir)
			rootViper.SetDefault(flagGitUsername, Cmd.GitUsername)
			rootViper.SetDefault(flagGitPassword, Cmd.GitPassword)
			rootViper.SetDefault(flagGitKeyFile, Cmd.GitKeyFile)
			rootViper.SetDefault(flagGitToken, Cmd.GitToken)
			rootViper.SetDefault(flagGitDiffCmd, Cmd.GitDiffCmd)

			err = rootViper.WriteConfig()
			if err != nil {
				break
			}

			err = rootViper.ReadInConfig()
			break
		}
		if err != nil {
			break
		}

		err = rootViper.MergeInConfig()
		if err != nil {
			break
		}

		// err = viper.Unmarshal(Cmd)
	}

	return err
}

func writeConfig() error {
	var err error

	for range Only.Once {
		err = rootViper.MergeInConfig()
		if err != nil {
			break
		}

		rootViper.Set(flagDebug, Cmd.Debug)
		rootViper.Set(flagQuiet, Cmd.Quiet)

		rootViper.Set(flagX32Host, Cmd.X32Host)
		rootViper.Set(flagX32Port, Cmd.X32Port)
		rootViper.Set(flagX32Timeout, Cmd.X32Timeout)
		rootViper.Set(flagX32Username, Cmd.X32Username)
		rootViper.Set(flagX32Password, Cmd.X32Password)

		rootViper.Set(flagMqttUsername, Cmd.MqttUsername)
		rootViper.Set(flagMqttPassword, Cmd.MqttPassword)
		rootViper.Set(flagMqttHost, Cmd.MqttHost)
		rootViper.Set(flagMqttPort, Cmd.MqttPort)

		rootViper.Set(flagGoogleSheet, Cmd.GoogleSheet)
		rootViper.Set(flagGoogleSheetUpdate, Cmd.GoogleSheetUpdate)

		rootViper.Set(flagGitRepo, Cmd.GitRepo)
		rootViper.Set(flagGitRepoDir, Cmd.GitRepoDir)
		rootViper.Set(flagGitUsername, Cmd.GitUsername)
		rootViper.Set(flagGitPassword, Cmd.GitPassword)
		rootViper.Set(flagGitKeyFile, Cmd.GitKeyFile)
		rootViper.Set(flagGitToken, Cmd.GitToken)
		rootViper.Set(flagGitDiffCmd, Cmd.GitDiffCmd)

		err = rootViper.WriteConfig()
		if err != nil {
			break
		}
	}

	return err
}

func readConfig() error {
	var err error

	for range Only.Once {
		err = rootViper.ReadInConfig()
		if err != nil {
			break
		}

		_, _ = fmt.Fprintln(os.Stderr, "Config file settings:")

		_, _ = fmt.Fprintf(os.Stderr, "Behringer X32 Host:	%v\n", rootViper.Get(flagX32Host))
		_, _ = fmt.Fprintf(os.Stderr, "Behringer X32 Port:	%v\n", rootViper.Get(flagX32Port))
		_, _ = fmt.Fprintf(os.Stderr, "Behringer X32 User:	%v\n", rootViper.Get(flagX32Username))
		_, _ = fmt.Fprintf(os.Stderr, "Behringer X32 Password:	%v\n", rootViper.Get(flagX32Password))
		_, _ = fmt.Fprintf(os.Stderr, "Behringer X32 Timeout:	%v\n", rootViper.Get(flagX32Timeout))
		_, _ = fmt.Fprintln(os.Stderr)

		_, _ = fmt.Fprintf(os.Stderr, "HASSIO mqtt Username:		%v\n", rootViper.Get(flagMqttUsername))
		_, _ = fmt.Fprintf(os.Stderr, "HASSIO mqtt Password:		%v\n", rootViper.Get(flagMqttPassword))
		_, _ = fmt.Fprintf(os.Stderr, "HASSIO mqtt Host:		%v\n", rootViper.Get(flagMqttHost))
		_, _ = fmt.Fprintf(os.Stderr, "HASSIO mqtt Port:		%v\n", rootViper.Get(flagMqttPort))
		_, _ = fmt.Fprintln(os.Stderr)

		_, _ = fmt.Fprintf(os.Stderr, "Git Repo URL:		%v\n", rootViper.Get(flagGitRepo))
		_, _ = fmt.Fprintf(os.Stderr, "Git Repo Dir:		%v\n", rootViper.Get(flagGitRepoDir))
		_, _ = fmt.Fprintf(os.Stderr, "Git Repo User:	%v\n", rootViper.Get(flagGitUsername))
		_, _ = fmt.Fprintf(os.Stderr, "Git Repo ApiPassword:	%v\n", rootViper.Get(flagGitPassword))
		_, _ = fmt.Fprintf(os.Stderr, "Git SSH keyfile:	%v\n", rootViper.Get(flagGitKeyFile))
		_, _ = fmt.Fprintf(os.Stderr, "Git Auth:	%v\n", rootViper.Get(flagGitToken))
		_, _ = fmt.Fprintf(os.Stderr, "Git Diff Command:	%v\n", rootViper.Get(flagGitDiffCmd))
		_, _ = fmt.Fprintln(os.Stderr)

		_, _ = fmt.Fprintf(os.Stderr, "Debug:		%v\n", rootViper.Get(flagDebug))
		_, _ = fmt.Fprintf(os.Stderr, "Quiet:		%v\n", rootViper.Get(flagQuiet))
	}

	return err
}

func bindFlags(cmd *cobra.Command, v *viper.Viper) error {
	var err error

	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		// Environment variables can't have dashes in them, so bind them to their equivalent
		// keys with underscores, e.g. --favorite-color to STING_FAVORITE_COLOR
		if strings.Contains(f.Name, "-") {
			envVarSuffix := strings.ToUpper(strings.ReplaceAll(f.Name, "-", "_"))
			err = v.BindEnv(f.Name, fmt.Sprintf("%s_%s", EnvPrefix, envVarSuffix))
		}

		// Apply the viper config value to the flag when the flag is not set and viper has a value
		if !f.Changed && v.IsSet(f.Name) {
			val := v.Get(f.Name)
			err = cmd.Flags().Set(f.Name, fmt.Sprintf("%v", val))
		}
	})

	return err
}
