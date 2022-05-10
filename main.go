package main

import (
	"fmt"
	"github.com/MickMake/GoX32/Only"
	"github.com/MickMake/GoX32/cmd"
	"log"
	"os"
	"os/signal"
	"syscall"
)


// func configure() {
// 	viper.SetConfigName("x32-mqtt")
// 	viper.SetConfigType("yaml")
// 	viper.AddConfigPath(".")
// 	viper.AddConfigPath("/etc/x32-mqtt")
// 	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
// 	viper.AutomaticEnv()
//
// 	// OSC Configuration
// 	viper.SetDefault("osc.host", "localhost")
// 	viper.SetDefault("osc.port", 1234)
//
// 	// MQTT Configuration
// 	viper.SetDefault("mqtt.broker", "localhost:1883")
// 	viper.SetDefault("mqtt.username", "")
// 	viper.SetDefault("mqtt.password", "")
// 	viper.SetDefault("mqtt.client_id", "osc-mqtt")
// 	viper.SetDefault("mqtt.keep_alive", "300s")
// 	viper.SetDefault("topic_prefix", "sound")
//
// 	if err := viper.ReadInConfig(); err != nil {
// 		log.Println(err)
// 		_ = viper.SafeWriteConfig()
// 		log.Println("Running with default config...")
// 	}
// }

func main() {
	var err error

	for range Only.Once {
		err = cmd.Execute()
		if err != nil {
			break
		}

	}

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
	}

	// configure()
	log.Println("Starting x32-mqtt")

	// setupMQTTClient()
	// setupOSCClient()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	<-sigChan
	log.Println("Exiting...")
}
