package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/MickMake/GoX32/Behringer"
	"github.com/MickMake/GoX32/Behringer/api"
	"github.com/MickMake/GoX32/Only"
	"github.com/MickMake/GoX32/defaults"
	"github.com/MickMake/GoX32/mmHa"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/go-co-op/gocron"
	"github.com/spf13/cobra"
	"reflect"
	"strings"
	"time"
)


func AttachCmdMqtt(cmd *cobra.Command) *cobra.Command {
	// ******************************************************************************** //
	var cmdMqtt = &cobra.Command{
		Use:                   "mqtt",
		Aliases:               []string{""},
		Short:                 fmt.Sprintf("Connect to a HASSIO broker."),
		Long:                  fmt.Sprintf("Connect to a HASSIO broker."),
		DisableFlagParsing:    false,
		DisableFlagsInUseLine: false,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return Cmd.PreRun(cmd, args, Cmd.ProcessArgs)
		},
		RunE:                  cmdMqttFunc,
		Args:                  cobra.MinimumNArgs(1),
	}
	cmd.AddCommand(cmdMqtt)
	cmdMqtt.Example = PrintExamples(cmdMqtt, "run", "sync")


	// ******************************************************************************** //
	var cmdMqttRun = &cobra.Command{
		Use:                   "run",
		Aliases:               []string{""},
		Short:                 fmt.Sprintf("One-off sync to a HASSIO broker."),
		Long:                  fmt.Sprintf("One-off sync to a HASSIO broker."),
		DisableFlagParsing:    false,
		DisableFlagsInUseLine: false,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return Cmd.PreRun(cmd, args, Cmd.ProcessArgs, Cmd.X32Args, Cmd.MqttArgs)
		},
		RunE:                  cmdMqttRunFunc,
		Args:                  cobra.RangeArgs(0, 1),
	}
	cmdMqtt.AddCommand(cmdMqttRun)
	cmdMqttRun.Example = PrintExamples(cmdMqttRun, "")

	// ******************************************************************************** //
	var cmdMqttSync = &cobra.Command{
		Use:                   "sync",
		Aliases:               []string{""},
		Short:                 fmt.Sprintf("Sync to a HASSIO MQTT broker."),
		Long:                  fmt.Sprintf("Sync to a HASSIO MQTT broker."),
		DisableFlagParsing:    false,
		DisableFlagsInUseLine: false,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return Cmd.PreRun(cmd, args, Cmd.ProcessArgs, Cmd.X32Args, Cmd.MqttArgs)
		},
		RunE:                  cmdMqttSyncFunc,
		Args:                  cobra.RangeArgs(0, 1),
	}
	cmdMqtt.AddCommand(cmdMqttSync)
	cmdMqttSync.Example = PrintExamples(cmdMqttSync, "", "all")

	return cmdMqtt
}

func cmdMqttFunc(cmd *cobra.Command, _ []string) error {
	return cmd.Help()
}

func cmdMqttRunFunc(_ *cobra.Command, _ []string) error {
	for range Only.Once {
		Cmd.Error = Cmd.MqttCron()
		if Cmd.Error != nil {
			break
		}

		LogPrintDate("Starting ticker...\n")
		updateCounter := 0
		timer := time.NewTicker(60 * time.Second)
		for t := range timer.C {
			if updateCounter < 5 {
				updateCounter++
				LogPrintDate("Sleeping: %d\n", updateCounter)
				continue
			}

			updateCounter = 0
			LogPrintDate("Update: %s\n", t.String())
			Cmd.Error = Cmd.MqttCron()
			if Cmd.Error != nil {
				break
			}

			// ep = Cmd.X32.QueryDevice(psId)
			// if ep.IsError() {
			// 	Cmd.Error = ep.GetError()
			// 	break
			// }
			//
			// data = ep.GetData()
			// for _, r := range data.Entries {
			// 	// fmt.Printf("%s ", r.PointId)
			// 	Cmd.Error = foo.SensorPublishState(r.PointId, r.Value)
			// 	if err != nil {
			// 		break
			// 	}
			// }
			// // fmt.Println()
		}
	}

	return Cmd.Error
}

func cmdMqttSyncFunc(_ *cobra.Command, args []string) error {

	for range Only.Once {
		// */1 * * * * /dir/command args args
		cronString := "*/5 * * * *"
		if len(args) > 0 {
			cronString = strings.Join(args[0:5], " ")
			cronString = strings.ReplaceAll(cronString, ".", "*")
		}

		Cron.Scheduler = gocron.NewScheduler(time.UTC)
		Cron.Scheduler = Cron.Scheduler.Cron(cronString)
		Cron.Scheduler = Cron.Scheduler.SingletonMode()

		Cmd.Error = Cmd.MqttCron()
		if Cmd.Error != nil {
			break
		}

		Cron.Job, Cmd.Error = Cron.Scheduler.Do(Cmd.MqttCron)
		if Cmd.Error != nil {
			break
		}

		LogPrintDate("Created job schedule using '%s'\n", cronString)
		Cron.Scheduler.StartBlocking()
		if Cmd.Error != nil {
			break
		}
	}

	return Cmd.Error
}


func (ca *CommandArgs) MqttArgs(cmd *cobra.Command, args []string) error {
	for range Only.Once {
		LogPrintDate("Connecting to MQTT HASSIO Service...\n")
		ca.Mqtt = mmHa.New(mmHa.Mqtt {
			ClientId: defaults.BinaryName,
			Username: ca.MqttUsername,
			Password: ca.MqttPassword,
			Host:     ca.MqttHost,
			Port:     ca.MqttPort,
		})
		ca.Error = ca.Mqtt.GetError()
		if ca.Error != nil {
			break
		}

		// if ca.Mqtt.EntityPrefix == "" {
		// 	ca.Mqtt.EntityPrefix = ca.X32.Prefix
		// }
		if ca.Mqtt.EntityPrefix == "" {
			ca.Mqtt.EntityPrefix = defaults.BinaryName
		}

		ca.Error = ca.Mqtt.SetDeviceConfig(
			defaults.BinaryName,
			ca.X32.Info.Model,
			defaults.BinaryName,
			ca.X32.Info.Model,
			"Behringer",
			"Studio",
			)
		if ca.Error != nil {
			break
		}

		ca.Error = ca.X32.SetMessageHandler(X32MessageHandler)
		if ca.Error != nil {
			break
		}

		ca.Error = ca.Mqtt.SetMessageHandler(MqttMessageHandler)
		if ca.Error != nil {
			break
		}

		ca.Error = ca.Mqtt.Connect()
		if ca.Error != nil {
			break
		}

		// if ca.Mqtt.PsId == 0 {
		// 	ca.Mqtt.PsId, ca.Error = ca.X32.GetPsId()
		// 	if ca.Error != nil {
		// 		break
		// 	}
		// 	LogPrintDate("Found X32 device %d\n", ca.Mqtt.PsId)
		// }
	}

	return ca.Error
}

func (ca *CommandArgs) MqttCron() error {
	for range Only.Once {
		if ca.Mqtt == nil {
			ca.Error = errors.New("mqtt not available")
			break
		}

		if ca.X32 == nil {
			ca.Error = errors.New("Behringer X32 not available")
			break
		}

		if ca.Mqtt.IsFirstRun() {
			ca.Mqtt.UnsetFirstRun()
			ca.Error = ca.X32.GetAllInfo()
			if ca.Error != nil {
				break
			}
		}

		newDay := false
		if ca.Mqtt.IsNewDay() {
			newDay = true
		}

		ca.Error = ca.Update1(newDay)
		if ca.Error != nil {
			break
		}

		time.Sleep(time.Hour * 6)

		ca.Mqtt.LastRefresh = time.Now()
	}

	if ca.Error != nil {
		LogPrintDate("Error: %s\n", ca.Error)
	}
	return ca.Error
}

func (ca *CommandArgs) Update1(newDay bool) error {
	for range Only.Once {
		// var pm api.PointsMap
		// pm, ca.Error = api.ImportPoints("points.json", ca.X32.Info.Model)
		// fmt.Printf("\n%v\n", pm)

		// /-show/showfile/scene/000/name
		ca.Error = ca.X32.Emit("/-show/showfile/show/name")
		ca.Error = ca.X32.Emit("/-prefs/lamp")
		ca.Error = ca.X32.Emit("/-prefs/lampon")
		ca.Error = ca.X32.Emit("/-prefs/lcdbright")
		ca.Error = ca.X32.Emit("/-prefs/lcdcont")
		ca.Error = ca.X32.Emit("/-prefs/rec_control")
		ca.Error = ca.X32.Emit("/-prefs/lamp")
		ca.Error = ca.X32.Emit("/-prefs/sceneadvance")
		ca.Error = ca.X32.Emit("/-prefs/selfollowsbank")
		ca.Error = ca.X32.Emit("/-prefs/show_control")
		ca.Error = ca.X32.Emit("/-prefs/style")
		ca.Error = ca.X32.Emit("/-prefs/viewrtn")
		ca.Error = ca.X32.Emit("/-prefs/??????")
		// ca.Error = ca.X32.Call("foo")

		// ca.Error = ca.X32.Call("/-show/showfile/show/buses")
		// fmt.Printf("%v\n", foo)
		// ca.Error = ca.X32.Call("/-show/showfile/show/chan16")
		// fmt.Printf("%v\n", foo)
		// ca.Error = ca.X32.Call("/-show/showfile/show/chan32")
		// fmt.Printf("%v\n", foo)
		// ca.Error = ca.X32.Call("/-show/showfile/show/console")
		// fmt.Printf("%v\n", foo)
		// ca.Error = ca.X32.Call("/-show/showfile/show/effects")
		// fmt.Printf("%v\n", foo)
		// ca.Error = ca.X32.Call("/-show/showfile/show/inputs")
		// fmt.Printf("%v\n", foo)
		// // ca.Error = ca.X32.Call("/-show/showfile/show/lrmtxdce")
		// // fmt.Printf("%v\n", foo)
		// ca.Error = ca.X32.Call("/-show/showfile/show/mxbuses")
		// fmt.Printf("%v\n", foo)
		// ca.Error = ca.X32.Call("/-show/showfile/show/mxsends")
		// fmt.Printf("%v\n", foo)
		// ca.Error = ca.X32.Call("/-show/showfile/show/return")
		// fmt.Printf("%v\n", foo)

		foo2 := ca.X32.GetScene(0)
		fmt.Printf("%v\n", foo2)
		foo2 = ca.X32.GetScene(1)
		fmt.Printf("%v\n", foo2)
		foo2 = ca.X32.GetScene(2)
		fmt.Printf("%v\n", foo2)

		time.Sleep(time.Hour * 24)

		ca.PublishChannels()



	// 	// Also getPowerStatistics, getHouseholdStoragePsReport, getPsList, getUpTimePoint,
	// 	ep := Cmd.X32.QueryDevice(Cmd.Mqtt.PsId)
	// 	if ep.IsError() {
	// 		Cmd.Error = ep.GetError()
	// 		break
	// 	}
	// 	data := ep.GetData()
	//
	// 	if newDay {
	// 		LogPrintDate("New day: Configuring %d entries in HASSIO.\n", len(data.Entries))
	// 		for _, o := range data.Order {
	// 			r := data.Entries[o]
	//
	// 			fmt.Printf("C")
	// 			re := mmHa.EntityConfig {
	// 				Name:        r.Point.Id, // PointName,
	// 				SubName:     "",
	// 				ParentId:    r.EndPoint,
	// 				ParentName:  "",
	// 				UniqueId:    r.Point.Id,
	// 				FullId:      r.Point.FullId,
	// 				Units:       r.Point.Unit,
	// 				ValueName:   r.Point.Name,
	// 				DeviceClass: "",
	// 				StateClass:  r.Point.Type,
	// 				Value:       r.Value,
	//
	// 				// Icon:                   "",
	// 				// ValueTemplate:          "",
	// 				// LastReset:              "",
	// 				// LastResetValueTemplate: "",
	// 			}
	//
	// 			// if re.LastResetValueTemplate != "" {
	// 			// 	fmt.Printf("HEY\n")
	// 			// }
	//
	// 			Cmd.Error = Cmd.Mqtt.BinarySensorPublishConfig(re)
	// 			if Cmd.Error != nil {
	// 				break
	// 			}
	//
	// 			Cmd.Error = Cmd.Mqtt.SensorPublishConfig(re)
	// 			if Cmd.Error != nil {
	// 				break
	// 			}
	// 		}
	// 		fmt.Println()
	// 	}
	//
	// 	LogPrintDate("Updating %d entries to HASSIO.\n", len(data.Entries))
	// 	for _, o := range data.Order {
	// 		r := data.Entries[o]
	//
	// 		fmt.Printf("U")
	// 		re := mmHa.EntityConfig {
	// 			Name:        r.Point.Id, // PointName,
	// 			SubName:     "",
	// 			ParentId:    r.EndPoint,
	// 			ParentName:  "",
	// 			UniqueId:    r.Point.Id,
	// 			// UniqueId:    r.Id,
	// 			FullId: r.Point.FullId,
	// 			// FullName:    r.Point.Name,
	// 			Units:       r.Point.Unit,
	// 			ValueName:   r.Point.Name,
	// 			// ValueName:   r.Id,
	// 			DeviceClass: "",
	// 			StateClass:  r.Point.Type,
	// 			Value:       r.Value,
	// 		}
	//
	// 		Cmd.Error = Cmd.Mqtt.BinarySensorPublishValue(re)
	// 		if Cmd.Error != nil {
	// 			break
	// 		}
	//
	// 		Cmd.Error = Cmd.Mqtt.SensorPublishValue(re)
	// 		if Cmd.Error != nil {
	// 			break
	// 		}
	// 	}
	// 	fmt.Println()
	}

	if Cmd.Error != nil {
		LogPrintDate("Error: %s\n", Cmd.Error)
	}
	return Cmd.Error
}

func (ca *CommandArgs) PublishChannel(channel int) {
	for range Only.Once {
		LogPrintDate("PublishChannel(%v)\n", channel)

		c := ca.X32.GetChannel(channel)
		fmt.Printf("\n%s\n", c.Json())
		dm := api.NewDataMap()
		dm.StructToPoints("ch0", c.Data)
		fmt.Printf("dm:\n%v\n", dm)
		fmt.Println("")
		// name := mmHa.JoinStringsForId(c.Topic)

		var entities []mmHa.EntityConfig
		for _, o := range dm.Order {
			r := dm.Entries[o]

			entities = append(entities, mmHa.EntityConfig {
				Name:        mmHa.JoinStrings(c.Topic, r.Point.Id),	// r.Point.Id,
				SubName:     "",
				ParentId:    ca.X32.Info.Model,	// r.EndPoint,
				ParentName:  ca.X32.Info.Model,	// "",
				UniqueId:    mmHa.JoinStringsForId(c.Topic, r.Point.Id),	// r.Point.Id,
				// FullId:      mmHa.JoinStringsForId(c.Topic, r.Point.Id),	// r.Point.FullId,
				Units:       r.Point.Unit,
				ValueName:   r.Point.Name,
				DeviceClass: "",
				StateClass:  "measurement",	// r.Point.Type,
				StateTopic:  mmHa.JoinStringsForId(c.Topic),
				Value:       r.Value,
				ValueTemplate:          fmt.Sprintf("{{ value_json.%s }}", r.Point.Name),

				// ParentId:    Cmd.X32.Info.Model,
				// ParentName:  Cmd.X32.Info.Model,
				// UniqueId:    fmt.Sprintf("%s %d", msg.Address, i),
				// FullId:      fmt.Sprintf("%s %d", msg.Address, i),
				// Units:       msg.GetType(),
				// ValueName:   fmt.Sprintf("%v", v),
				// DeviceClass: "",
				// StateClass:  "r.Point.Type",
				// StateTopic:  name,
				// Value:       fmt.Sprintf("%v", v),
				// ValueTemplate:          fmt.Sprintf("{{ value_json.value%d }}", i),

				// Icon:                   "",
				// LastReset:              "",
				// LastResetValueTemplate: "",
			})
			fmt.Println()

			// if !msg.SeenBefore {
			// 	Cmd.Error = Cmd.Mqtt.PublishConfig(ec)
			// 	if Cmd.Error != nil {
			// 		LogPrintDate("MQTT: Could not publish: %s\n", Cmd.Error)
			// 		break
			// 	}
			// }
			//
			// Cmd.Error = Cmd.Mqtt.PublishValue(ec)
			// if Cmd.Error != nil {
			// 	LogPrintDate("MQTT: Could not publish: %s\n", Cmd.Error)
			// 	break
			// }
		}

		Cmd.Error = Cmd.Mqtt.PublishConfigs(entities)
		if Cmd.Error != nil {
			LogPrintDate("MQTT: Could not publish: %s\n", Cmd.Error)
			break
		}

		Cmd.Error = Cmd.Mqtt.SensorPublishValues(entities)
		if Cmd.Error != nil {
			LogPrintDate("MQTT: Could not publish: %s\n", Cmd.Error)
			break
		}
	}
}

func (ca *CommandArgs) PublishChannels() {
	for range Only.Once {
		for c := 0; c <= 32; c++ {
			ca.PublishChannel(c)
		}
	}
}


// Publish example:
// topic: GoX32/set/ch/01/mix/fader
// message: [{"type":"string","value":"0.8"}]

type dataField struct {
	Type  string `json:"type"`
	Value any    `json:"value"`
}
type mqttPayload []*dataField

func MqttMessageHandler(_ mqtt.Client, message mqtt.Message) {
	for range Only.Once {
		prefixParts := strings.Split(Cmd.Mqtt.EntityPrefix, "/")
		LogPrintDate("MqttMessageHandler() - Topic: %v\n", prefixParts)

		parts := strings.Split(message.Topic(), "/")
		parts = parts[len(prefixParts)+1:]
		res := make(mqttPayload, 0, 2)
		Cmd.Error = json.Unmarshal(message.Payload(), &res)
		if Cmd.Error != nil {
			LogPrintDate("MQTT: Invalid message payload: %s\n", Cmd.Error)
			break
		}

		values := make([]any, 0, len(res))
		for _, p := range res {
			switch p.Type {
				case reflect.TypeOf(float32(0)).String():
					values = append(values, float32(p.Value.(float64)))
				case reflect.TypeOf("").String():
					values = append(values, p.Value.(string))
			}
		}

		address := "/" + strings.Join(parts, "/")
		LogPrintDate("Address: %v\n", address)
		Cmd.Error = Cmd.X32.Client.EmitMessage(address, values...)
		if Cmd.Error != nil {
			LogPrintDate("Could not send OSC message: %s\n", Cmd.Error)
			break
		}
	}
}

func X32MessageHandler(msg *Behringer.Message) {
	for range Only.Once {
		LogPrintDate("X32MessageHandler() - %v\t", msg.Address)
		// a := fmt.Sprintf("%s", msg.Address)
		// a = strings.TrimPrefix(a, "-")
		// a = strings.TrimSuffix(a, "-")
		// a = strings.TrimPrefix(a, "_")
		// a = strings.TrimSuffix(a, "_")
		// name := mmHa.JoinStringsForId(a)

		p := Cmd.X32.Points.Resolve(msg.Address)
		if p == nil {
			fmt.Printf("- Missing Point: %v\n", p)
			break
		}

		// p.CorrectUnit(msg.GetType())	// @TODO - Not really required since we're now referencing all endpoints.

		if len(msg.Arguments) == 1 {
			value := fmt.Sprintf("%v", msg.Arguments[0])

			// if msg.Address == "/ch/01/mix/fader" {
			// 	foo1 := api.ToLinDb(value, "0", "1")
			// 	foo2 := api.ToLogDb(value)
			// 	fmt.Printf("foo1: %s\tfoo2: %s\n", foo1, foo2)
			// }

			value = p.Convert.GetString(value)
			fmt.Printf("- Value: %s %s\t (%v)\n", value, p.Unit, msg.Arguments[0])

			ec := mmHa.EntityConfig {
				Name:        p.Name,
				SubName:     "",
				ParentId:    p.ParentId,
				ParentName:  p.ParentId,
				UniqueId:    api.CleanString(p.Id),
				Units:       p.Unit,	// msg.GetType(),
				ValueName:   "",
				DeviceClass: "",
				StateClass:  "measurement",
				Value:       value,

				// Icon:                   "",
				// ValueTemplate:          "",
				// LastReset:              "",
				// LastResetValueTemplate: "",
			}

			if !msg.SeenBefore {
				Cmd.Error = Cmd.Mqtt.PublishConfig(ec)
				if Cmd.Error != nil {
					LogPrintDate("MQTT: Could not publish: %s\n", Cmd.Error)
					break
				}
			}

			Cmd.Error = Cmd.Mqtt.PublishValue(ec)
			if Cmd.Error != nil {
				LogPrintDate("MQTT: Could not publish: %s\n", Cmd.Error)
				break
			}
			break
		}

		var entities []mmHa.EntityConfig
		var values []string
		for i := range msg.Arguments {
			value := fmt.Sprintf("%v", msg.Arguments[i])
			values = append(values, value)
			p.Unit = ""		// @TODO - Fix this up!

			ec := mmHa.EntityConfig {
				Name:        fmt.Sprintf("%s %d", p.Name, i),
				SubName:     "",
				ParentId:    p.ParentId,
				ParentName:  p.ParentId,
				UniqueId:    api.CleanString(fmt.Sprintf("%s_%d", p.Id, i)),
				Units:       p.Unit,	// msg.GetType(),
				ValueName:   "",
				DeviceClass: "",
				StateClass:  "measurement",
				Value:       value,

				StateTopic:  p.Name,
				ValueTemplate:          fmt.Sprintf("{{ value_json.value%d }}", i),

				// Icon:                   "",
				// LastReset:              "",
				// LastResetValueTemplate: "",
			}
			entities = append(entities, ec)

			// if !msg.SeenBefore {
			// 	Cmd.Error = Cmd.Mqtt.PublishConfig(ec)
			// 	if Cmd.Error != nil {
			// 		LogPrintDate("MQTT: Could not publish: %s\n", Cmd.Error)
			// 		break
			// 	}
			// }
			//
			// Cmd.Error = Cmd.Mqtt.PublishValue(ec)
			// if Cmd.Error != nil {
			// 	LogPrintDate("MQTT: Could not publish: %s\n", Cmd.Error)
			// 	break
			// }
		}
		fmt.Printf("- Values: %s %s\n", strings.Join(values, ", "), p.Unit)

		if !msg.SeenBefore {
			Cmd.Error = Cmd.Mqtt.PublishConfigs(entities)
			if Cmd.Error != nil {
				LogPrintDate("MQTT: Could not publish: %s\n", Cmd.Error)
				break
			}
		}

		Cmd.Error = Cmd.Mqtt.SensorPublishValues(entities)
		if Cmd.Error != nil {
			LogPrintDate("MQTT: Could not publish: %s\n", Cmd.Error)
			break
		}
	}
}
