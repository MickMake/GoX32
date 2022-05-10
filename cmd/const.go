package cmd

import "time"


//goland:noinspection SpellCheckingInspection
const (
	EnvPrefix         = "X32"
	defaultConfigFile = "config.json"
	// defaultTokenFile  = "AuthToken.json"

	flagConfigFile = "config"
	flagDebug      = "debug"
	flagQuiet      = "quiet"

	flagX32Host       = "x32-host"
	flagX32Port       = "x32-port"
	flagX32Timeout    = "x32-timeout"
	flagX32Username   = "x32-user"
	flagX32Password   = "x32-password"
	flagX32OutputType = "out"

	flagMqttUsername   = "mqtt-user"
	flagMqttPassword   = "mqtt-password"
	flagMqttHost       = "mqtt-host"
	flagMqttPort       = "mqtt-port"

	flagGoogleSheet       = "google-sheet"
	flagGoogleSheetUpdate = "update"

	flagGitUsername = "git-username"
	flagGitPassword = "git-password"
	flagGitKeyFile  = "git-sshkey"
	flagGitToken    = "git-token"
	flagGitRepo     = "git-repo"
	flagGitRepoDir  = "git-dir"
	flagGitDiffCmd  = "diff-cmd"

	defaultHost      = ""
	defaultPort      = "10023"
	defaultUsername  = ""
	defaultPassword  = ""

	defaultTimeout = time.Duration(time.Second * 30)
)
