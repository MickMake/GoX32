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
		fmt.Println("Done")

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

					if x.Debug { fmt.Printf("# CacheWrite()\n") }
					x.Error = x.CacheWrite()
					if x.Error != nil {
						break
					}
			}
		}
	}

	if x.Error != nil {
		log.Println(x.Error)
	}
}

// func (x *X32) Process(point string, value ...any) Message {		// (*api.Point, api.UnitValueMap, error) {
// 	// var ret *api.Point
// 	// values := make(api.UnitValueMap)
// 	// var err error
// 	var msg Message
//
// 	for range Only.Once {
// 		msg.Point = x.Points.Resolve(point)
// 		if msg.Point == nil {
// 			msg.Error = errors.New(fmt.Sprintf("Missing Point: %v data: %v\n", point, value))
// 			break
// 		}
//
// 		// gv := msg.Point.Convert.GetValues(value...)
// 		// keys := make([]string, 0, len(gv))
// 		// for k := range gv {
// 		// 	keys = append(keys, k)
// 		// }
// 		// sort.Strings(keys)
// 		//
// 		// for _, k := range keys {
// 		// 	v2 := gv[k]
// 		// 	vf, _ := strconv.ParseFloat(v2, 64)
// 		// 	vi, _ := strconv.ParseInt(v2, 10, 64)
// 		// 	vb = fmt.Sprintf("%t", )
// 		//
// 		// 	values[k] = api.UnitValue {
// 		// 		Unit:        ret.Unit,
// 		// 		ValueString: v2,
// 		// 		ValueFloat:  vf,
// 		// 		ValueInt:    vi,
// 		// 		ValueBool:   vb,
// 		// 	}
// 		// }
// 	}
//
// 	return msg
// 	// return ret, values, err
// }

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
		x.Error = x.Client.EmitMessage(address, args...)
		break
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

// func (x *X32) CallProcess(address string, args ...any) *Message {	// (*api.Point, api.UnitValueMap, error) {
// 	// var ret *api.Point
// 	// values := make(api.UnitValueMap)
// 	// var err error
// 	var msg *Message
//
// 	for range Only.Once {
// 		msg = x.Call(address, args...)
// 		if msg.Error != nil {
// 			x.Error = msg.Error
// 			break
// 		}
//
// 		x.Error = msg.Process()
// 		if x.Error != nil {
// 			break
// 		}
// 	}
//
// 	return msg
// 	// return ret, values, err
// }

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

type MessageHandlerFunc func(msg *Message)
func (x *X32) SetMessageHandler(fn MessageHandlerFunc) error {
	for range Only.Once {
		x.messageHandler = fn
	}
	return x.Error
}

func (x *X32) oscMessageHandler(msg *gosc.Message) {
	for range Only.Once {
		m := x.UpdateCache(&Message{ Message: msg })
		if x.Debug {
			fmt.Printf("# oscMessageHandler() - msg: %v\n", msg)
		}

		m.Point = x.Points.Resolve(msg.Address)
		if m.Point == nil {
			x.Error = errors.New(fmt.Sprintf("Missing Point: %v data: %v\n", msg.Address, msg.Arguments))
			fmt.Printf("%s", x.Error)
			break
		}

		x.Error = m.Process()
		if x.Error != nil {
			fmt.Printf("%s", x.Error)
			break
		}

		if x.messageHandler == nil {
			break
		}
		x.messageHandler(m)
	}
}

// func (x *X32) MultipleMessageHandler(msg *gosc.Message) {
// 	for range Only.Once {
// 		m := x.UpdateCache(msg)
// 		if x.Debug {
// 			fmt.Printf("# oscMessageHandler() - msg: %v\n", msg)
// 		}
//
// 		m.Point, m.UnitValueMap, m.Error = x.Process(msg.Address, msg.Arguments...)
// 		if m.Error != nil {
// 			break
// 		}
// 	}
// }

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
