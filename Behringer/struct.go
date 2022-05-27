package Behringer

import (
	"errors"
	"fmt"
	"github.com/MickMake/GoX32/Behringer/api"
	"github.com/MickMake/GoX32/Behringer/api/output"
	"github.com/MickMake/GoX32/Only"
	"github.com/loffa/gosc"
	"log"
	"os"
	"path/filepath"
	"time"
)


type X32 struct {
	Address    string
	Prefix     string
	Client     *gosc.Client
	Error      error
	NeedLogin  bool
	Debug      bool

	OutputType output.OutputType

	Info Info

	Points api.PointsMap

	messageHandler MessageHandlerFunc

	meters    map[string]bool
	configDir string
	cache MessageMap
	cacheDir string
	cacheTimeout time.Duration
	cacheLastWrite time.Time

	autopoll api.RefreshPoints
	autopollBatchCount int
	// autopollIndex int
	// autopollTime time.Time
}

type Info struct {
	HardwareVersion    string
	Name               string
	Model              string
	FirmwareVersion    string
}
func (i *Info) String() string {
	return fmt.Sprintf("Model:\t%s\nHW Version:\t%s\nFW Version:\t%s\nName:\t%s\n",
		i.Model,
		i.HardwareVersion,
		i.FirmwareVersion,
		i.Name,
		)
}

func CreateInfo(args []any) Info {
	var ret Info
	for range Only.Once {
		switch {
			case len(args) <= 4:
				ret.FirmwareVersion = fmt.Sprintf("%s", args[3])
				fallthrough
			case len(args) <= 3:
				ret.Model = fmt.Sprintf("%s", args[2])
				fallthrough
			case len(args) <= 2:
				ret.Name = fmt.Sprintf("%s", args[1])
				fallthrough
			case len(args) <= 1:
				ret.HardwareVersion = fmt.Sprintf("%s", args[0])
		}
	}
	return ret
}


type ArgsX32 struct {
	Host      string
	Port      string
	ConfigDir    string
	CacheDir     string
	CacheTimeout time.Duration
}
func NewX32(args ArgsX32) *X32 {
// func NewX32(host string, port string, configDir string, cacheDir string, cacheTimeout time.Duration) *X32 {
	var x X32

	for range Only.Once {
		if args.Host == "" {
			x.Error = errors.New("invalid x32 host")
			break
		}
		if args.Port == "" {
			args.Port = DefaultPort
		}

		x.Address = fmt.Sprintf("%s:%s", args.Host, args.Port)
		x.Client, x.Error = gosc.NewClient(x.Address)
		if x.Error != nil {
			break
		}

		x.Error = x.Client.HandleMessageFunc("/*", x.oscMessageHandler)

		if x.Prefix == "" {
			x.Prefix = "x32"
		}

		x.Error = x.SetConfigDir(args.ConfigDir)
		if x.Error != nil {
			break
		}

		x.cache = make(MessageMap)
		x.SetCacheTimeout(args.CacheTimeout)
		x.Error = x.SetCacheDir(args.CacheDir)
		if x.Error != nil {
			break
		}
	}
	return &x
}


func (x *X32) Connect() error {
	for range Only.Once {
		x.Error = x.CacheRead()
		if x.Error != nil {
			break
		}

		fmt.Printf("Fetching console info...")
		var info *gosc.Message
		info, x.Error = x.Client.CallMessage("/info")
		if x.Error != nil {
			fmt.Printf("\nCould not get mixer info: %s\n", x.Error)
			break
		}
		fmt.Println("Done")

		x.Info = CreateInfo(info.Arguments)
		x.Prefix = x.Info.Model

		// Look for files in directory and parse all json.
		var files []string
		files, x.Error = output.DirectoryRead(x.configDir, "points_.*.json")
		if x.Error != nil {
			fmt.Printf("\nCould not get mixer info: %s\n", x.Error)
			break
		}

		fmt.Printf("Importing points...")
		x.Points, x.Error = api.ImportPoints(x.Info.Model, files...)
		if x.Error != nil {
			break
		}
		fmt.Printf("Done. %d points found.\n", len(x.Points))

		fmt.Printf("Importing refresh points...")
		x.autopoll, x.autopollBatchCount, x.Error = api.ImportRefresh(files...)
		if x.Error != nil {
			break
		}
		fmt.Printf("Done. %d points found.\n", len(x.autopoll))

		go x.XremoteSender()
	}

	return x.Error
}

func (x *X32) SetConfigDir(basedir string) error {
	for range Only.Once {
		x.configDir = filepath.Join(basedir)
		_, x.Error = os.Stat(x.configDir)
		if os.IsExist(x.Error) {
			x.Error = nil
			break
		}

		x.Error = os.MkdirAll(x.configDir, 0700)
		if x.Error != nil {
			break
		}
	}

	return x.Error
}

func (x *X32) StartMeters(meters ...string) error {
	for range Only.Once {
		if x.meters == nil {
			x.meters = make(map[string]bool)
		}

		for _, m := range meters {
			x.meters[m] = true
		}

		x.Error = x.Emit("/meters", StringArrayToAny(meters...)...)
		if x.Error != nil {
			break
		}

		x.Error = x.Emit("/-prefs/remote/ioenable", 4089)
		if x.Error != nil {
			break
		}
	}

	return x.Error
}

func (x *X32) StopMeters(meters ...string) error {
	for range Only.Once {
		for _, m := range meters {
			delete(x.meters, m)
		}
	}

	return x.Error
}

func (x *X32) renewMeters() error {
	for range Only.Once {
		if len(x.meters) == 0 {
			break
		}

		var am []any
		for _, m := range x.meters {
			am = append(am, m)
		}

		x.Error = x.Emit("/meters", am...)
		if x.Error != nil {
			break
		}
	}

	return x.Error
}

func (x *X32) XremoteSender() {
	for range Only.Once {

		ticker := time.NewTicker(time.Second * 9)
		x.Error = x.Client.EmitMessage("/xremote")
		if x.Error != nil {
			break
		}

		for {
			select {
				case _ = <-ticker.C:
					if x.Debug { fmt.Printf("# XremoteSender()\n") }
					x.Error = x.Client.EmitMessage("/xremote")
					if x.Error != nil {
						break
					}

					if x.Debug { fmt.Printf("# renewMeters()\n") }
					x.Error = x.renewMeters()
					if x.Error != nil {
						break
					}

					if x.Debug { fmt.Printf("# renewAutopoll()\n") }
					x.Error = x.refreshAutopoll()
					if x.Error != nil {
						break
					}

					// if x.Debug { fmt.Printf("# CacheWrite()\n") }
					// x.Error = x.CacheWrite()
					// if x.Error != nil {
					// 	break
					// }
			}
		}
	}

	if x.Error != nil {
		log.Println(x.Error)
	}
}

func (x *X32) Get(address string, wait bool) *Message {
	var msg *Message

	for range Only.Once {
		if !wait {
			x.Error = x.Emit(address)
			msg.Error = x.Error
			break
		}

		msg = x.Call(address)
	}

	return msg
}

func (x *X32) Emit(address string, args ...any) error {

	for range Only.Once {
		if x.Debug {
			fmt.Printf("# Emit() - address: %v, args: %v\n", address, args)
		}

		args = FixAny(args...)

		// for i := range args {
		// 	// t := reflect.TypeOf(a)
		// 	// n := t.Name()
		// 	// fmt.Printf("TYPE: %s\n", n)
		// 	// if tv, ok := typeMapper[t]; ok {
		// 	// 	args[i] = (a)
		// 	// }
		//
		// 	switch v := args[i].(type) {
		// 		case int:
		// 			args[i] = int32(v)
		// 		case int8:
		// 			args[i] = int32(v)
		// 		case int16:
		// 			args[i] = int32(v)
		// 		case int32:
		// 			args[i] = int32(v)
		// 		case int64:
		// 			args[i] = int32(v)
		// 		case uint:
		// 			args[i] = int32(v)
		// 		case uint8:
		// 			args[i] = int32(v)
		// 		case uint16:
		// 			args[i] = int32(v)
		// 		case uint32:
		// 			args[i] = int32(v)
		// 		case uint64:
		// 			args[i] = int32(v)
		//
		// 		case float32:
		// 			args[i] = float32(v)
		// 		case float64:
		// 			args[i] = float32(v)
		//
		// 		case string:
		// 			args[i] = string(v)
		// 	}
		// }

		x.Error = x.Client.EmitMessage(address, args...)
		if x.Error != nil {
			break
		}
	}

	return x.Error
}

func (x *X32) Call(address string, args ...any) *Message {
	var msg *Message

	for range Only.Once {
		msg = &Message {
			Message:    nil,
			SeenBefore: false,
			LastSeen:   time.Now(),
			Counter:    1,
			Type:       "call",
			Error:      nil,
		}
		if x.Debug {
			fmt.Printf("# Call() - msg: %v\n", msg)
		}

		args = FixAny(args...)

		msg.Message, x.Error = x.Client.CallMessage(address, args...)
		if x.Debug {
			fmt.Printf("# Call() - msg.Message: %v\n", msg.Message)
		}

		msg = x.UpdateCache(msg)
		if msg.Error != nil {
			break
		}

		msg.Point = x.Points.Resolve(address)
		if msg.Point == nil {
			msg.Error = errors.New(fmt.Sprintf("Missing Point: %v data: %v\n", address, args))
			break
		}

		x.Error = msg.Process()
		if x.Error != nil {
			break
		}
	}

	return msg
}

func (x *X32) GetTopic(msg *gosc.Message) string {
	topic := fmt.Sprintf("%s%s", x.Prefix, msg.Address)
	if x.Debug {
		fmt.Printf("# GetTopic() - Topic: %s\n", topic)
	}
	return topic
}

func (x *X32) ListEndpoints() error {
	fmt.Printf("%v", x.Points)
	return nil
}

// func (x *X32) AddToAutoPoll(address ...string) error {
// 	for range Only.Once {
// 		x.autopoll = append(x.autopoll, address...)
// 	}
// 	return x.Error
// }

func (x *X32) refreshAutopoll() error {
	for range Only.Once {
		// current := x.autopollIndex
		// size := len(x.autopoll)
		//
		// for i := 0; i < 16; i++ {
		// 	x.autopollIndex++
		// 	if x.autopollIndex > len(x.autopoll) {
		// 		x.autopollIndex = 0
		// 	}
		// 	fmt.Printf("Autopoll[%d]: %s\n", x.autopollIndex, x.autopoll[x.autopollIndex])
		// 	x.Error = x.Client.EmitMessage(x.autopoll[x.autopollIndex])
		// 	if x.Error != nil {
		// 		break
		// 	}
		// }

		batchLimit := 0
		for endpoint, v := range x.autopoll {
			if v.IsExpired() {
				// fmt.Printf("# [%d]Renewing %s Last: %s\n", batchLimit, endpoint, v.When.Format("15:04:05"))
				v.Reset()

				// if x.Debug {
				// 	fmt.Printf("# renewAutopoll() - %s\n", endpoint)
				// }
				x.Error = x.Client.EmitMessage(endpoint)
				if x.Error != nil {
					break
				}
				batchLimit++
				if batchLimit > x.autopollBatchCount {
					// fmt.Printf("# Reached batch limit.\n")
					break
				}
			}
		}

		if batchLimit > 0 {
			fmt.Printf("# Renewed %d endpoints\n", batchLimit)
		}
	}

	return x.Error
}


type MessageHandlerFunc func(msg *Message)
func (x *X32) SetMessageHandler(fn MessageHandlerFunc) error {
	for range Only.Once {
		x.messageHandler = fn
	}
	return x.Error
}

func (x *X32) oscMessageHandler(msg *gosc.Message) {
	for range Only.Once {
		// m := x.UpdateCache(&Message{ Message: msg })
		// if x.Debug {
		// 	fmt.Printf("# oscMessageHandler() - msg: %v\n", msg)
		// }
		//
		// m.Point = x.Points.Resolve(msg.Address)
		// if m.Point == nil {
		// 	x.Error = errors.New(fmt.Sprintf("Missing Point: %v data: %v\n", msg.Address, msg.Arguments))
		// 	fmt.Printf("%s", x.Error)
		// 	break
		// }
		//
		// x.Error = m.Process()
		// if x.Error != nil {
		// 	fmt.Printf("%s", x.Error)
		// 	break
		// }

		m := x.Process(msg, false)

		if x.messageHandler == nil {
			break
		}
		x.messageHandler(m)
	}
}

func (x *X32) Process(msg *gosc.Message, toX32 bool) *Message {
	var m *Message

	for range Only.Once {
		m = x.UpdateCache(&Message{ Message: msg })
		if x.Debug {
			fmt.Printf("# oscMessageHandler() - msg: %v\n", msg)
		}

		m.Point = x.Points.Resolve(msg.Address)
		if m.Point == nil {
			x.Error = errors.New(fmt.Sprintf("Missing Point: %v data: %v\n", msg.Address, msg.Arguments))
			fmt.Printf("%s", x.Error)
			break
		}

		if toX32 {
			x.Error = m.Process()
			if x.Error != nil {
				break
			}
			break
		}

		x.Error = m.Process()
		if x.Error != nil {
			break
		}
	}

	return m
}

func (x *X32) Output(endpoint api.EndPoint, table *output.Table, graphFilter string) error {
	for range Only.Once {
		switch {
			case x.OutputType.IsNone():

			case x.OutputType.IsHuman():
				if table == nil {
					break
				}
				table.Print()

			case x.OutputType.IsFile():
				if table == nil {
					break
				}
				x.Error = table.WriteCsvFile()

			case x.OutputType.IsRaw():
				fmt.Println(endpoint.GetJsonData(true))

			case x.OutputType.IsJson():
				fmt.Println(endpoint.GetJsonData(false))

			case x.OutputType.IsGraph():
				if table == nil {
					break
				}
				x.Error = table.SetGraphFromJson(output.Json(graphFilter))
				if x.Error != nil {
					break
				}
				x.Error = table.CreateGraph()
				if x.Error != nil {
					break
				}

			default:
		}
	}

	return x.Error
}

func (x *X32) OutputTable(table *output.Table) error {
	for range Only.Once {
		switch {
			case x.OutputType.IsNone():

			case x.OutputType.IsHuman():
				if table == nil {
					break
				}
				table.Print()

			case x.OutputType.IsFile():
				if table == nil {
					break
				}
				x.Error = table.WriteCsvFile()

			default:
		}
	}

	return x.Error
}


// var typeMapper = map[reflect.Type]reflect.Type {
// 	reflect.TypeOf(int(0)):   int32(0),
// 	reflect.TypeOf(int8(0)):   int32(0),
// 	reflect.TypeOf(int16(0)):   int32(0),
// 	reflect.TypeOf(int32(0)):   int32(0),
// 	reflect.TypeOf(int64(0)):   int32(0),
// 	reflect.TypeOf(uint(0)):   int32(0),
// 	reflect.TypeOf(uint8(0)):   int32(0),
// 	reflect.TypeOf(uint16(0)):   int32(0),
// 	reflect.TypeOf(uint32(0)):   int32(0),
// 	reflect.TypeOf(uint64(0)):   int32(0),
// 	reflect.TypeOf(byte(0)):   int32(0),
//
// 	reflect.TypeOf(float32(0)): float32(0),
// 	reflect.TypeOf(float64(0)): float32(0),
//
// 	reflect.TypeOf(""):         "",
// }

func AnyToStringArray(array ...any) []string {
	var ret []string
	for _, a := range array {
		ret = append(ret, fmt.Sprintf("%v", a))
	}
	return ret
}

func StringArrayToAny(array ...string) []any {
	var ret []any
	for _, a := range array {
		ret = append(ret, a)
	}
	return ret
}

func FixAny(args ...any) []any {
	for i := range args {
		// t := reflect.TypeOf(a)
		// n := t.Name()
		// fmt.Printf("TYPE: %s\n", n)
		// if tv, ok := typeMapper[t]; ok {
		// 	args[i] = (a)
		// }

		switch v := args[i].(type) {
			case int:
				args[i] = int32(v)
			case int8:
				args[i] = int32(v)
			case int16:
				args[i] = int32(v)
			case int32:
				args[i] = int32(v)
			case int64:
				args[i] = int32(v)
			case uint:
				args[i] = int32(v)
			case uint8:
				args[i] = int32(v)
			case uint16:
				args[i] = int32(v)
			case uint32:
				args[i] = int32(v)
			case uint64:
				args[i] = int32(v)

			case float32:
				args[i] = float32(v)
			case float64:
				args[i] = float32(v)

			case string:
				args[i] = string(v)
		}
	}
	return args
}

