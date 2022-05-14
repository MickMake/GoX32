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

	OutputType output.OutputType

	Info Info

	Points api.PointsMap

	messageHandler MessageHandlerFunc

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

		x.Points, x.Error = api.ImportPoints(x.Info.Model, files...)
		if x.Error != nil {
			break
		}

		go x.XremoteSender()

		// fmt.Println("Got info from the OSC server:", info.Arguments)
		// time.Sleep(time.Second * 30)
		// sg.Areas = make(api.Areas)
		// sg.Areas[api.GetArea(AliSmsService.Area{})] = api.AreaStruct(AliSmsService.Init(sg.ApiRoot))
		// sg.Areas[api.GetArea(AppService.Area{})] = api.AreaStruct(AppService.Init(sg.ApiRoot))
		// sg.Areas[api.GetArea(MttvScreenService.Area{})] = api.AreaStruct(MttvScreenService.Init(sg.ApiRoot))
		// sg.Areas[api.GetArea(PowerPointService.Area{})] = api.AreaStruct(PowerPointService.Init(sg.ApiRoot))
		// sg.Areas[api.GetArea(WebAppService.Area{})] = api.AreaStruct(WebAppService.Init(sg.ApiRoot))
		// sg.Areas[api.GetArea(WebIscmAppService.Area{})] = api.AreaStruct(WebIscmAppService.Init(sg.ApiRoot))
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
					fmt.Printf("# XremoteSender()\n")
					x.Error = x.Client.EmitMessage("/xremote")
					if x.Error != nil {
						break
					}
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

func (x *X32) Send(msg string) error {
	x.Error = x.Client.EmitMessage(msg)
	return x.Error
}
func (x *X32) GetStatus() error {
	return x.Send("/status")
}
func (x *X32) GetInfo() error {
	return x.Send("/info")
}
func (x *X32) GetXinfo() error {
	return x.Send("/xinfo")
}
func (x *X32) GetShowDump() error {
	return x.Send("/showdump")
}

func (x *X32) Call(address string) *Message {
	var msg Message
	for range Only.Once {
		var m *gosc.Message
		m, msg.Error = x.Client.CallMessage(address)
		if msg.Error != nil {
			x.Error = msg.Error
			break
		}
		msg = Message {
			Message:    *m,
			SeenBefore: false,
			LastSeen:   time.Now(),
			Counter:    1,
			Type:       "call",
			Error:      nil,
		}
	}
	return &msg
}
func (x *X32) CallStatus() *Message {
	return x.Call("/status")
}
func (x *X32) CallInfo() *Message {
	return x.Call("/info")
}
func (x *X32) CallXinfo() *Message {
	return x.Call("/xinfo")
}
func (x *X32) CallShowDump() *Message {
	return x.Call("/showdump")
}
func (x *X32) CallDeskName() *Message {
	return x.Call("/-prefs/name")
}

func (x *X32) GetTopic(msg *gosc.Message) string {
	topic := fmt.Sprintf("%s%s", x.Prefix, msg.Address)
	fmt.Printf("# GetTopic() - Topic: %s\n", topic)
	return topic
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
		m := x.UpdateCache(msg)
		topic := x.GetTopic(msg)
		fmt.Printf("MSG: %s\n%v\n", topic, msg)
		if x.messageHandler == nil {
			break
		}
		x.messageHandler(m)
	}
}

// func (sg *Behringer) AppendUrl(endpoint string) api.EndPointUrl {
// 	return sg.ApiRoot.AppendUrl(endpoint)
// }
//
// func (sg *Behringer) GetEndpoint(ae string) api.EndPoint {
// 	var ep api.EndPoint
// 	for range Only.Once {
// 		area, endpoint := sg.SplitEndPoint(ae)
// 		if sg.Error != nil {
// 			break
// 		}
//
// 		ep = sg.Areas.GetEndPoint(api.AreaName(area), api.EndPointName(endpoint))
// 		if ep == nil {
// 			sg.Error = errors.New("EndPoint not found")
// 			break
// 		}
//
// 		if ep.IsDisabled() {
// 			sg.Error = errors.New("API EndPoint is not implemented")
// 			break
// 		}
//
// 		if sg.Auth.Token() != "" {
// 			ep = ep.SetRequest(api.RequestCommon{
// 				Appkey:    sg.GetAppKey(), // sg.Auth.RequestCommon.Appkey
// 				Lang:      "_en_US",
// 				SysCode:   "200",
// 				Token:     sg.GetToken(),
// 				UserID:    sg.GetUserId(),
// 				ValidFlag: "1,3",
// 			})
// 		}
// 	}
// 	return ep
// }
//
// func (sg *Behringer) GetByJson(endpoint string, request string) api.EndPoint {
// 	var ret api.EndPoint
// 	for range Only.Once {
// 		ret = sg.GetEndpoint(endpoint)
// 		if sg.Error != nil {
// 			break
// 		}
// 		if ret.IsError() {
// 			sg.Error = ret.GetError()
// 			break
// 		}
//
// 		if request != "" {
// 			ret = ret.SetRequestByJson(output.Json(request))
// 			if ret.IsError() {
// 				fmt.Println(ret.Help())
// 				sg.Error = ret.GetError()
// 				break
// 			}
// 		}
//
// 		ret = ret.Call()
// 		if ret.IsError() {
// 			fmt.Println(ret.Help())
// 			sg.Error = ret.GetError()
// 			break
// 		}
//
// 		switch {
// 		case sg.OutputType.IsNone():
//
// 		case sg.OutputType.IsFile():
// 			sg.Error = ret.WriteDataFile()
//
// 		case sg.OutputType.IsRaw():
// 			fmt.Println(ret.GetJsonData(true))
//
// 		case sg.OutputType.IsJson():
// 			fmt.Println(ret.GetJsonData(false))
//
// 		default:
// 		}
// 	}
// 	return ret
// }
//
// func (sg *Behringer) GetByStruct(endpoint string, request interface{}, cache time.Duration) api.EndPoint {
// 	var ret api.EndPoint
// 	for range Only.Once {
// 		ret = sg.GetEndpoint(endpoint)
// 		if sg.Error != nil {
// 			break
// 		}
// 		if ret.IsError() {
// 			sg.Error = ret.GetError()
// 			break
// 		}
//
// 		if request != nil {
// 			ret = ret.SetRequest(request)
// 			if ret.IsError() {
// 				sg.Error = ret.GetError()
// 				break
// 			}
// 		}
//
// 		ret = ret.SetCacheTimeout(cache)
// 		// if ret.CheckCache() {
// 		// 	ret = ret.ReadCache()
// 		// 	if !ret.IsError() {
// 		// 		break
// 		// 	}
// 		// }
//
// 		ret = ret.Call()
// 		if ret.IsError() {
// 			sg.Error = ret.GetError()
// 			break
// 		}
//
// 		// sg.Error = ret.WriteCache()
// 		// if sg.Error != nil {
// 		// 	break
// 		// }
// 	}
//
// 	return ret
// }
//
// func (sg *Behringer) SplitEndPoint(ae string) (string, string) {
// 	var area string
// 	var endpoint string
//
// 	for range Only.Once {
// 		s := strings.Split(ae, ".")
// 		switch len(s) {
// 		case 0:
// 			sg.Error = errors.New("empty endpoint")
//
// 		case 1:
// 			area = "AppService"
// 			endpoint = s[0]
//
// 		case 2:
// 			area = s[0]
// 			endpoint = s[1]
//
// 		default:
// 			sg.Error = errors.New("too many delimeters defined, (only one '.' allowed)")
// 		}
// 	}
//
// 	return area, endpoint
// }
//
// func (sg *Behringer) ListEndpoints(area string) error {
// 	return sg.Areas.ListEndpoints(area)
// }
//
// func (sg *Behringer) ListAreas() {
// 	sg.Areas.ListAreas()
// }
//
// func (sg *Behringer) AreaExists(area string) bool {
// 	return sg.Areas.Exists(area)
// }
//
// func (sg *Behringer) AreaNotExists(area string) bool {
// 	return sg.Areas.NotExists(area)
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

// func (sg *Behringer) Login(auth login.BehringerAuth) error {
// 	for range Only.Once {
// 		a := sg.GetEndpoint(AppService.GetAreaName() + ".login")
// 		sg.Auth = login.Assert(a)
//
// 		sg.Error = sg.Auth.Login(&auth)
// 		if sg.Error != nil {
// 			break
// 		}
// 	}
//
// 	return sg.Error
// }
//
// func (sg *Behringer) GetToken() string {
// 	return sg.Auth.Token()
// }
//
// func (sg *Behringer) GetUserId() string {
// 	return sg.Auth.UserId()
// }
//
// func (sg *Behringer) GetAppKey() string {
// 	return sg.Auth.AppKey()
// }
//
// func (sg *Behringer) GetLastLogin() string {
// 	return sg.Auth.LastLogin().Format(login.DateTimeFormat)
// }
//
// func (sg *Behringer) GetUserName() string {
// 	return sg.Auth.UserName()
// }
//
// func (sg *Behringer) GetUserEmail() string {
// 	return sg.Auth.Email()
// }
//
// func (sg *Behringer) HasTokenChanged() bool {
// 	return sg.Auth.HasTokenChanged()
// }
