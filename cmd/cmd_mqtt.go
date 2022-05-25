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
	"github.com/loffa/gosc"
	"github.com/spf13/cobra"
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

		time.Sleep(time.Hour * 24)

		ca.Mqtt.LastRefresh = time.Now()
	}

	if ca.Error != nil {
		LogPrintDate("Error: %s\n", ca.Error)
	}
	return ca.Error
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

		data := make(map[string]string)

		Cmd.Error = json.Unmarshal(message.Payload(), &data)
		if Cmd.Error != nil {
			LogPrintDate("MQTT: Invalid message payload: %s\n", Cmd.Error)
			break
		}

		// values := make([]any, 0, len(res))
		// for _, p := range res {
		// 	switch p.Type {
		// 		case reflect.TypeOf(float32(0)).String():
		// 			values = append(values, float32(p.Value.(float64)))
		// 		case reflect.TypeOf("").String():
		// 			values = append(values, p.Value.(string))
		// 	}
		// }

		for address, value := range data {
			// address := "/" + strings.Join(parts, "/")
			LogPrintDate("Update: %v -> %v\n", address, value)
			Cmd.Error = Cmd.X32.Set(address, value)
			if Cmd.Error != nil {
				LogPrintDate("Could not send OSC message: %s\n", Cmd.Error)
				break
			}

			m := Cmd.X32.Process(&gosc.Message {
				Address:   address,
				Arguments: []any{value},
			})
			X32MessageHandler(m)
		}
	}
}

func X32MessageHandler(msg *Behringer.Message) {
	for range Only.Once {
		// Single value.
		if len(msg.UnitValueMap) == 1 {
			LogPrintDate("# Single Point:\n\t%s\n\tUnitValue: %s (%v)\n", msg.Point, msg.UnitValueMap, msg.Arguments[0])

			haType := "sensor"
			if msg.IsSwitch() {
				haType = "binary" // Will become a sensor and a toggle.
			}
			if msg.IsMomentary() {
				haType = "binary"	// Will become a sensor and a toggle.
			}
			// if msg.IsMap() {
			// 	haType = "select"
			// }

			ec := mmHa.EntityConfig {
				Name:        msg.Point.Name,
				SubName:     "",
				ParentId:    msg.Point.ParentId,
				ParentName:  msg.Point.ParentId,
				UniqueId:    api.CleanString(msg.Point.Id),
				Units:       msg.Point.Unit, // msg.GetType(),
				ValueName:   msg.Point.Id,
				DeviceClass: "",
				StateClass:  "measurement",
				Value:       msg.UnitValueMap.GetFirst().ValueString,
				HaType:      haType,

				// Icon:                   "",
				// ValueTemplate:          "",
				// LastReset:              "",
				// LastResetValueTemplate: "",
			}

			seen := msg.SeenBefore
			Cmd.Error = Cmd.Mqtt.Publish(ec, !seen)
			if Cmd.Error != nil {
				LogPrintDate("MQTT: Could not publish: %s\n", Cmd.Error)
				break
			}

			if msg.IsSwitch() {
				ec.HaType = "switch"
				Cmd.Error = Cmd.Mqtt.Publish(ec, !seen)
				if Cmd.Error != nil {
					LogPrintDate("MQTT: Could not publish: %s\n", Cmd.Error)
					break
				}
				break
			}

			if msg.IsMomentary() {
				ec.HaType = "button"
				ec.Value = "ON"
				Cmd.Error = Cmd.Mqtt.Publish(ec, !seen)
				if Cmd.Error != nil {
					LogPrintDate("MQTT: Could not publish: %s\n", Cmd.Error)
					break
				}
				break
			}

			if msg.IsIndex() {
				ec.HaType = "select"
				// ec.Value = "ON"
				ec.Options = msg.GetIndexOptions()
				Cmd.Error = Cmd.Mqtt.Publish(ec, !seen)
				if Cmd.Error != nil {
					LogPrintDate("MQTT: Could not publish: %s\n", Cmd.Error)
					break
				}
				break
			}

			break
		}


		// Multiple values.
		LogPrintDate("# Multiple Point:\n\t%s\n\tUnitValue: %s\n", msg.Point, msg.UnitValueMap)

		var entities []mmHa.EntityConfig
		for i, u := range msg.UnitValueMap {
			id := api.JoinStringsForId(i)

			haType := "sensor"
			if msg.Point.IsSwitch() {
				haType = "binary"
			}

			ec := mmHa.EntityConfig {
				Name:        fmt.Sprintf("%s %s", msg.Point.Name, i),
				SubName:     "",
				ParentId:    msg.Point.ParentId,
				ParentName:  msg.Point.ParentId,
				UniqueId:    api.CleanString(fmt.Sprintf("%s_%s", msg.Point.Id, id)),
				Units:       u.Unit,
				ValueName:   id,
				DeviceClass: "",
				StateClass:  "measurement",
				Value:       u.ValueString,
				HaType:      haType,

				StateTopic:    msg.Point.Name,
				ValueTemplate: fmt.Sprintf("{{ value_json.%s }}", id),

				// Icon:                   "",
				// LastReset:              "",
				// LastResetValueTemplate: "",
			}
			entities = append(entities, ec)
		}

		if !msg.SeenBefore {
			Cmd.Error = Cmd.Mqtt.PublishConfigs(entities)
			if Cmd.Error != nil {
				LogPrintDate("MQTT: Could not publish: %s\n", Cmd.Error)
				break
			}
		}

		Cmd.Error = Cmd.Mqtt.PublishSensorValues(entities)
		if Cmd.Error != nil {
			LogPrintDate("MQTT: Could not publish: %s\n", Cmd.Error)
			break
		}
	}
}


func (ca *CommandArgs) Update1(newDay bool) error {
	for range Only.Once {

		// var pm api.PointsMap
		// pm, ca.Error = api.ImportPoints("points.json", ca.X32.Info.Model)
		// fmt.Printf("\n%v\n", pm)
		// ca.Error = ca.X32.AddMeters("/meters/11")
		//
		// ca.Error = ca.X32.GetAllInfo()
		//
		// sm := gosc.Message {
		// 	Address:   "/node",
		// 	Arguments: []any{"ch/01/config"},
		// }
		// fmt.Printf("TYPE: %s\n", sm.GetType())
		// m := ca.X32.Call("/-prefs/viewrtn")
		// fmt.Printf("FOO:\n%v\n", m)
		//
		// m2, e := ca.X32.Client.SendAndReceiveMessage(&sm)
		// fmt.Printf("FOO:\n%v\n%s", m2, e)
		//
		// m = ca.X32.Call("/node", "ch/01/config")
		// fmt.Printf("FOO:\n%v\n", m)
		//
		// x.Error = x.Emit("/formatsubscribe", "hidden/states", "/-stat/tape/state", "/-usb/path", "/-usb/title", "/-stat/tape/etime", "/-stat/tape/rtime", "/-stat/aes50/state", "/-stat/aes50/A", "/-stat/aes50/B", "/-show/prepos/current", "/-stat/usbmounted", "/-usb/dir/dirpos", "/-usb/dir/maxpos", "/-stat/xcardtype", "/-stat/xcardsync", "/-stat/rtasource", "/-stat/talk/A", "/-stat/talk/B", "/-stat/osc/on", "/-stat/keysolo", "/-stat/urec/state", "/-stat/urec/etime", "/-stat/urec/rtime", 0, 0, 4)
		// x.Error = x.Emit("/formatsubscribe", "hidden/solo", "/-stat/solosw/**", 1, 80, 20)
		// x.Error = x.Emit("/formatsubscribe", "hidden/prefs", "/-prefs/clockrate", "/-prefs/clocksource", "/-prefs/scene_advance", "/-prefs/safe_masterlevels", "/-prefs/clockmode", "/-prefs/show_control", "/-prefs/haflags", "/-prefs/hardmute", "/-prefs/dcamute", "/-prefs/invertmutes", "/-prefs/remote/ioenable", "/-prefs/rta/source", "/-prefs/rta/pos", "/-prefs/rta/det", 0, 0, 10)
		// x.Error = x.Emit("/batchsubscribe", "meters/6", "/meters/6", 0, 0, 2)
		// x.Error = x.Emit("/batchsubscribe", "meters/9", "/meters/9", 0, 0, 2)
		// x.Error = x.Emit("/batchsubscribe", "meters/14", "/meters/14", 0, 0, 2)
		// x.Error = x.Emit("/batchsubscribe", "meters/10", "/meters/10", 0, 0, 2)
		// x.Error = x.Emit("/meters", "/meters/11")	//, 0, 0, 2)
		// x.Error = x.Emit("/-prefs/remote/ioenable", 4089)	//, 0, 0, 2)
		// x.Error = x.AddMeters("/meters/11")
		//
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
		//
		// m := x.Points["/meters/2"]
		// foo := m.Convert.Blob.Get(data)
		// fmt.Printf("FOO: %s\n", foo)
		//
		// hey1 := x.Call("/status")
		// fmt.Printf("%v\n", hey1)
		//
		// fmt.Println("")
		//
		//
		// hey := x.Emit("/meters/0", "")
		// fmt.Printf("%v\n", hey)
		// fmt.Println("")
		//
		// foo2 := ca.X32.GetScene(0)
		// fmt.Printf("%s\n", foo2)
		// foo2 = ca.X32.GetScene(1)
		// fmt.Printf("%s\n", foo2)
		// foo2 = ca.X32.GetScene(2)
		// fmt.Printf("%s\n", foo2)
		//
		// fmt.Println(ca.X32.Points.String())
		//
		// ca.Error = ca.X32.StartMeters("/meters/11")
		// time.Sleep(time.Second * 60)
		// ca.Error = ca.X32.StopMeters("/meters/11")
		//
		// fmt.Println("HEY1")
		// time.Sleep(time.Second * 5)
		// fmt.Println("HEY2")
		// ca.X32.Emit("/showdump")
		// time.Sleep(time.Second * 5)
		// fmt.Println("HEY3")
		// ca.X32.Emit("/showdump")
		// time.Sleep(time.Second * 5)
		// fmt.Println("HEY4")
		// ca.X32.Emit("/showdump")
		//
		// ca.PublishChannel(7)
		//
		// time.Sleep(time.Second * 5)
		// ca.PublishChannels()
		// time.Sleep(time.Second * 5)
		// ca.PublishMatrices()
		// time.Sleep(time.Second * 5)
		// ca.PublishBusses()
		// time.Sleep(time.Second * 5)
		// ca.PublishAuxes()
		//
		// ca.GetInitial()
		// ca.Error = ca.X32.Emit("/-show/showfile/show/name")
		// ca.Error = ca.X32.Emit("/-prefs/lamp")
		// ca.Error = ca.X32.Emit("/-prefs/lampon")
		// ca.Error = ca.X32.Emit("/-prefs/lcdbright")
		// ca.Error = ca.X32.Emit("/-prefs/lcdcont")
		// ca.Error = ca.X32.Emit("/-prefs/rec_control")
		// ca.Error = ca.X32.Emit("/-prefs/lamp")
		// ca.Error = ca.X32.Emit("/-prefs/sceneadvance")
		// ca.Error = ca.X32.Emit("/-prefs/selfollowsbank")
		// ca.Error = ca.X32.Emit("/-prefs/show_control")
		// ca.Error = ca.X32.Emit("/-prefs/style")
		// ca.Error = ca.X32.Emit("/-prefs/viewrtn")
		// ca.Error = ca.X32.Emit("/-prefs/??????")
		// ca.Error = ca.X32.Emit("/-stat/geqpos")
		// ca.Error = ca.X32.Emit("/-stat/rtageqpost")
		//
		// ca.Error = ca.X32.Emit("/-prefs/remote/ioenable", int32(4089))	//, 0, 0, 2)
		//
		// ca.Error = ca.X32.Emit("/formatsubscribe", "hidden/solo", "/-stat/solosw/**")	//, int32(1), int32(80), int32(20))
		//
		// ca.Error = ca.X32.Set("/ch/01/mix/on", api.Off)
		// ca.Error = ca.X32.Set("/-stat/solosw/01", api.Off)
		// time.Sleep(time.Second * 1)
		// ca.Error = ca.X32.Set("/ch/01/mix/on", api.On)
		// ca.Error = ca.X32.Set("/-stat/solosw/01", api.On)
		// time.Sleep(time.Second * 1)
		// ca.Error = ca.X32.Set("/ch/01/mix/on", api.Off)
		// ca.Error = ca.X32.Set("/-stat/solosw/01", api.Off)
		// time.Sleep(time.Second * 1)
		// ca.Error = ca.X32.Set("/ch/01/mix/on", api.On)
		// ca.Error = ca.X32.Set("/-stat/solosw/01", api.On)
		// time.Sleep(time.Second * 1)
		//
		// ca.Error = ca.X32.Set("/-action/clearsolo", api.On)
		//
		// p1 := ca.X32.Points.Resolve("/ch/01/mix/on")
		// p2 := ca.X32.Points.Resolve("/-stat/solosw/01")
		// ca.Error = ca.X32.Emit(p1.EndPoint, p1.Convert.SetValue(api.On))
		// ca.Error = ca.X32.Emit(p2.EndPoint, p2.Convert.SetValue(api.On))
		// time.Sleep(time.Second * 2)
		// ca.Error = ca.X32.Emit(p1.EndPoint, p1.Convert.SetValue(api.Off))
		// ca.Error = ca.X32.Emit(p2.EndPoint, p2.Convert.SetValue(api.Off))
		// time.Sleep(time.Second * 2)
		//
		// ca.Error = ca.X32.Emit("/formatsubscribe",
		// 	"hidden/names",
		// 	"/ch/**/config/name",
		// 	// "/bus/*/config/name",
		// 	// "/ch/*/config/name",
		// 	// "/dca/*/config/name",
		// 	// "/fxrtn/*/config/name",
		// 	// "/main/m/config/name",
		// 	// "/main/st/config/name",
		// 	// "/mtx/*/config/name",
		// 	int32(1), int32(8), int32(4),
		// )
		//
		// ca.Error = ca.X32.Emit("/-stat/screen/screen", int32(1))
		// time.Sleep(time.Second * 2)
		// ca.Error = ca.X32.Emit("/-stat/screen/screen", int32(2))
		// time.Sleep(time.Second * 2)
		// ca.Error = ca.X32.Emit("/-stat/screen/screen", int32(0))
		// time.Sleep(time.Second * 2)
		//
		// ca.Error = ca.X32.Emit("/-stat/screen/CHAN/page", int32(0))
		// time.Sleep(time.Second * 2)
		// ca.Error = ca.X32.Emit("/-stat/screen/CHAN/page", int32(1))
		// time.Sleep(time.Second * 2)
		// ca.Error = ca.X32.Emit("/-stat/screen/CHAN/page", int32(2))
		// time.Sleep(time.Second * 2)
		// ca.Error = ca.X32.Emit("/-stat/screen/CHAN/page", int32(3))
		// time.Sleep(time.Second * 2)
		// ca.Error = ca.X32.Emit("/-stat/screen/CHAN/page", int32(4))
		// time.Sleep(time.Second * 2)
		// ca.Error = ca.X32.Emit("/-stat/screen/CHAN/page", int32(5))
		// time.Sleep(time.Second * 2)
		// ca.Error = ca.X32.Emit("/-stat/screen/CHAN/page", int32(6))
		// time.Sleep(time.Second * 2)

		time.Sleep(time.Hour * 24)
	}

	if Cmd.Error != nil {
		LogPrintDate("Error: %s\n", Cmd.Error)
	}
	return Cmd.Error
}


func (ca *CommandArgs) GetInitial() {
	for range Only.Once {
		var entities []mmHa.EntityConfig

		for _, c := range ca.X32.ChannelCount() {
			for _, r := range ca.X32.GetChannel(c) {
				haType := "sensor"
				if r.Point.IsSwitch() {
					haType = "binary"
				}

				entity := mmHa.EntityConfig {
					Name:        r.Point.Name,
					SubName:     "",
					ParentId:    r.Point.ParentId,
					ParentName:  r.Point.ParentId,
					UniqueId:    r.Point.Id,
					Units:       r.Point.Unit,
					ValueName:   r.Point.Id,
					DeviceClass: "",
					StateClass:  "measurement",
					Value:       r.GetValueString(),
					HaType:      haType,

					StateTopic:    r.Point.Name,
					// ValueTemplate: fmt.Sprintf("{{ value_json.%s }}", r.Point.Id),

					// Icon:                   "",
					// LastReset:              "",
					// LastResetValueTemplate: "",
				}
				entities = append(entities, entity)
			}
		}

		for _, c := range ca.X32.BusCount() {
			for _, r := range ca.X32.GetBus(c) {
				haType := "sensor"
				if r.Point.IsSwitch() {
					haType = "binary"
				}

				entity := mmHa.EntityConfig {
					Name:        r.Point.Name,
					SubName:     "",
					ParentId:    r.Point.ParentId,
					ParentName:  r.Point.ParentId,
					UniqueId:    r.Point.Id,
					Units:       r.Point.Unit,
					ValueName:   r.Point.Id,
					DeviceClass: "",
					StateClass:  "measurement",
					Value:       r.GetValueString(),
					HaType:      haType,

					StateTopic:    r.Point.Name,
					// ValueTemplate: fmt.Sprintf("{{ value_json.%s }}", r.Point.Id),

					// Icon:                   "",
					// LastReset:              "",
					// LastResetValueTemplate: "",
				}
				entities = append(entities, entity)
			}
		}

		for _, c := range ca.X32.MatrixCount() {
			for _, r := range ca.X32.GetMatrix(c) {
				haType := "sensor"
				if r.Point.IsSwitch() {
					haType = "binary"
				}

				entity := mmHa.EntityConfig {
					Name:        r.Point.Name,
					SubName:     "",
					ParentId:    r.Point.ParentId,
					ParentName:  r.Point.ParentId,
					UniqueId:    r.Point.Id,
					Units:       r.Point.Unit,
					ValueName:   r.Point.Id,
					DeviceClass: "",
					StateClass:  "measurement",
					Value:       r.GetValueString(),
					HaType:      haType,

					StateTopic:    r.Point.Name,
					// ValueTemplate: fmt.Sprintf("{{ value_json.%s }}", r.Point.Id),

					// Icon:                   "",
					// LastReset:              "",
					// LastResetValueTemplate: "",
				}
				entities = append(entities, entity)
			}
		}

		for _, c := range ca.X32.AuxCount() {
			for _, r := range ca.X32.GetAux(c) {
				haType := "sensor"
				if r.Point.IsSwitch() {
					haType = "binary"
				}

				entity := mmHa.EntityConfig {
					Name:        r.Point.Name,
					SubName:     "",
					ParentId:    r.Point.ParentId,
					ParentName:  r.Point.ParentId,
					UniqueId:    r.Point.Id,
					Units:       r.Point.Unit,
					ValueName:   r.Point.Id,
					DeviceClass: "",
					StateClass:  "measurement",
					Value:       r.GetValueString(),
					HaType:      haType,

					StateTopic:    r.Point.Name,
					// ValueTemplate: fmt.Sprintf("{{ value_json.%s }}", r.Point.Id),

					// Icon:                   "",
					// LastReset:              "",
					// LastResetValueTemplate: "",
				}
				entities = append(entities, entity)
			}
		}

		for _, entity := range entities {
			Cmd.Error = Cmd.Mqtt.Publish(entity, true)
			if Cmd.Error != nil {
				LogPrintDate("MQTT: Could not publish: %s\n", Cmd.Error)
				break
			}

			time.Sleep(time.Millisecond * 10)
		}
	}
}


func (ca *CommandArgs) PublishChannels() {
	for range Only.Once {
		for _, c := range ca.X32.ChannelCount() {
			ca.PublishChannel(c)
		}
	}
}

func (ca *CommandArgs) PublishChannel(id int) {
	for range Only.Once {
		LogPrintDate("PublishChannel(%v)\n", id)

		c := ca.X32.GetChannel(id)
		name := fmt.Sprintf("Virtual Channel %.2d", id + 1)
		topic := api.JoinStringsForId(name)
		var entities []mmHa.EntityConfig
		for _, r := range c {

			haType := "sensor"
			if r.Point.IsSwitch() {
				haType = "binary"
			}

			ec := mmHa.EntityConfig {
				Name:        name + " " + r.Name,
				SubName:     "",
				ParentId:    r.Point.ParentId,
				ParentName:  r.Point.ParentId,
				UniqueId:    mmHa.JoinStringsForId(topic, r.Point.Id),
				Units:       r.Point.Unit,
				ValueName:   r.Point.Id,
				DeviceClass: "",
				StateClass:  "measurement",
				Value:       r.GetValueString(),
				HaType:      haType,

				StateTopic:    topic,
				ValueTemplate: fmt.Sprintf("{{ value_json.%s }}", r.Point.Id),

				// Icon:                   "",
				// LastReset:              "",
				// LastResetValueTemplate: "",
			}

			entities = append(entities, ec)
		}

		Cmd.Error = Cmd.Mqtt.PublishConfigs(entities)
		if Cmd.Error != nil {
			LogPrintDate("MQTT: Could not publish: %s\n", Cmd.Error)
			break
		}

		Cmd.Error = Cmd.Mqtt.PublishValues(entities)
		if Cmd.Error != nil {
			LogPrintDate("MQTT: Could not publish: %s\n", Cmd.Error)
			break
		}
	}
}


func (ca *CommandArgs) PublishMatrices() {
	for range Only.Once {
		for _, c := range ca.X32.MatrixCount() {
			ca.PublishMatrix(c)
		}
	}
}

func (ca *CommandArgs) PublishMatrix(id int) {
	for range Only.Once {
		LogPrintDate("PublishMatrix(%v)\n", id)

		c := ca.X32.GetMatrix(id)
		name := fmt.Sprintf("Virtual Matrix %.2d", id + 1)
		topic := api.JoinStringsForId(name)
		var entities []mmHa.EntityConfig
		for _, r := range c {

			haType := "sensor"
			if r.Point.IsSwitch() {
				haType = "binary"
			}

			ec := mmHa.EntityConfig {
				Name:        name + " " + r.Name,
				SubName:     "",
				ParentId:    r.Point.ParentId,
				ParentName:  r.Point.ParentId,
				UniqueId:    mmHa.JoinStringsForId(topic, r.Point.Id),
				Units:       r.Point.Unit,
				ValueName:   r.Point.Id,
				DeviceClass: "",
				StateClass:  "measurement",
				Value:       r.GetValueString(),
				HaType:      haType,

				StateTopic:    topic,
				ValueTemplate: fmt.Sprintf("{{ value_json.%s }}", r.Point.Id),

				// Icon:                   "",
				// LastReset:              "",
				// LastResetValueTemplate: "",
			}

			entities = append(entities, ec)
		}

		Cmd.Error = Cmd.Mqtt.PublishConfigs(entities)
		if Cmd.Error != nil {
			LogPrintDate("MQTT: Could not publish: %s\n", Cmd.Error)
			break
		}

		Cmd.Error = Cmd.Mqtt.PublishSensorValues(entities)
		if Cmd.Error != nil {
			LogPrintDate("MQTT: Could not publish: %s\n", Cmd.Error)
			break
		}
	}
}


func (ca *CommandArgs) PublishBusses() {
	for range Only.Once {
		for _, c := range ca.X32.BusCount() {
			ca.PublishBus(c)
		}
	}
}

func (ca *CommandArgs) PublishBus(id int) {
	for range Only.Once {
		LogPrintDate("PublishBus(%v)\n", id)

		c := ca.X32.GetBus(id)
		name := fmt.Sprintf("Virtual Bus %.2d", id + 1)
		topic := api.JoinStringsForId(name)
		var entities []mmHa.EntityConfig
		for _, r := range c {

			haType := "sensor"
			if r.Point.IsSwitch() {
				haType = "binary"
			}

			ec := mmHa.EntityConfig {
				Name:        name + " " + r.Name,
				SubName:     "",
				ParentId:    r.Point.ParentId,
				ParentName:  r.Point.ParentId,
				UniqueId:    mmHa.JoinStringsForId(topic, r.Point.Id),
				Units:       r.Point.Unit,
				ValueName:   r.Point.Id,
				DeviceClass: "",
				StateClass:  "measurement",
				Value:       r.GetValueString(),
				HaType:      haType,

				StateTopic:    topic,
				ValueTemplate: fmt.Sprintf("{{ value_json.%s }}", r.Point.Id),

				// Icon:                   "",
				// LastReset:              "",
				// LastResetValueTemplate: "",
			}

			entities = append(entities, ec)
		}

		Cmd.Error = Cmd.Mqtt.PublishConfigs(entities)
		if Cmd.Error != nil {
			LogPrintDate("MQTT: Could not publish: %s\n", Cmd.Error)
			break
		}

		Cmd.Error = Cmd.Mqtt.PublishSensorValues(entities)
		if Cmd.Error != nil {
			LogPrintDate("MQTT: Could not publish: %s\n", Cmd.Error)
			break
		}
	}
}


func (ca *CommandArgs) PublishAuxes() {
	for range Only.Once {
		for _, c := range ca.X32.AuxCount() {
			ca.PublishAux(c)
		}
	}
}

func (ca *CommandArgs) PublishAux(id int) {
	for range Only.Once {
		LogPrintDate("PublishAux(%v)\n", id)

		c := ca.X32.GetAux(id)
		name := fmt.Sprintf("Virtual Aux %.2d", id + 1)
		topic := api.JoinStringsForId(name)
		var entities []mmHa.EntityConfig
		for _, r := range c {

			haType := "sensor"
			if r.Point.IsSwitch() {
				haType = "binary"
			}

			ec := mmHa.EntityConfig {
				Name:        name + " " + r.Name,
				SubName:     "",
				ParentId:    r.Point.ParentId,
				ParentName:  r.Point.ParentId,
				UniqueId:    mmHa.JoinStringsForId(topic, r.Point.Id),
				Units:       r.Point.Unit,
				ValueName:   r.Point.Id,
				DeviceClass: "",
				StateClass:  "measurement",
				Value:       r.GetValueString(),
				HaType:      haType,

				StateTopic:    topic,
				ValueTemplate: fmt.Sprintf("{{ value_json.%s }}", r.Point.Id),

				// Icon:                   "",
				// LastReset:              "",
				// LastResetValueTemplate: "",
			}

			entities = append(entities, ec)
		}

		Cmd.Error = Cmd.Mqtt.PublishConfigs(entities)
		if Cmd.Error != nil {
			LogPrintDate("MQTT: Could not publish: %s\n", Cmd.Error)
			break
		}

		Cmd.Error = Cmd.Mqtt.PublishSensorValues(entities)
		if Cmd.Error != nil {
			LogPrintDate("MQTT: Could not publish: %s\n", Cmd.Error)
			break
		}
	}
}
