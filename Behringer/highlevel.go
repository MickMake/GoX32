package Behringer

import (
	"encoding/json"
	"fmt"
	"github.com/MickMake/GoX32/Only"
)


type Channel struct {
	Topic       string  `json:"-"`

	Data struct {
		Name     string  `json:"name"`
		Colour   string  `json:"colour"`
		Icon     string  `json:"icon"`
		Mute     bool    `json:"mute"`
		Solo     bool    `json:"solo"`
		Source   string  `json:"source"`
		Gain     float64 `json:"gain"`
		Trim     float64 `json:"trim"`
		Phantom  bool    `json:"phantom"`
		Phantom2 bool    `json:"phantom2"`
		Selected bool    `json:"selected"`

		Fader float64 `json:"fader"`
		Level float64 `json:"level"`
	}
}
type Channels []Channel	// We're keeping a 1:1 mapping of array -> channel numbers.

func (c *Channel) Json() string {
	var ret string

	for range Only.Once {
		j, err := json.Marshal(c)
		if err != nil {
			break
		}
		ret = string(j)
	}

	return ret
}

func (x *X32) GetChannel(i int) Channel {
	var ret Channel

	for range Only.Once {
		// channels start from index 1
		ret.Topic = fmt.Sprintf("/ch/%.2d/", i + 1)

		msg := x.Call(ret.Topic + "config/name")
		if msg.Error != nil {
			break
		}
		ret.Data.Name = msg.GetStringValue()

		msg = x.Call(ret.Topic + "config/color")
		if msg.Error != nil {
			break
		}
		ret.Data.Colour = msg.GetStringValue()

		msg = x.Call(ret.Topic + "config/source")
		if msg.Error != nil {
			break
		}
		ret.Data.Source = msg.GetStringValue()

		msg = x.Call(ret.Topic + "config/icon")
		if msg.Error != nil {
			break
		}
		ret.Data.Icon = msg.GetStringValue()

		msg = x.Call(ret.Topic + "preamp/trim")
		if msg.Error != nil {
			break
		}
		ret.Data.Trim = msg.GetFloatValue()

		// msg = x.Call(ret.Topic + "pramp/hpon")
		// if msg.Error != nil {
		// 	break
		// }
		// ret.Phantom2 = msg.GetBoolValue()

		msg = x.Call(ret.Topic + "mix/on")
		if msg.Error != nil {
			break
		}
		ret.Data.Mute = !msg.GetBoolValue()		// Mix ON == Mute OFF

		msg = x.Call(ret.Topic + "mix/fader")
		if msg.Error != nil {
			break
		}
		ret.Data.Fader = msg.GetFloatValue()

		msg = x.Call(ret.Topic + "mix/01/level")
		if msg.Error != nil {
			break
		}
		ret.Data.Level = msg.GetFloatValue()


		// headamps start from index 0
		// base = fmt.Sprintf("/headamp/%.3d/", i)

		msg = x.Call(fmt.Sprintf("/headamp/%.3d/gain", i))
		if msg.Error != nil {
			break
		}
		ret.Data.Gain = msg.GetFloatValue()

		msg = x.Call(fmt.Sprintf("/headamp/%.3d/phantom", i))
		if msg.Error != nil {
			break
		}
		ret.Data.Mute = msg.GetBoolValue()
	}

	return ret
}

func (x *X32) GetChannels() Channels {
	var ret Channels

	for range Only.Once {
		for c := 0; c < 32; c++ {
			ret = append(ret, x.GetChannel(c))
		}
	}

	return ret
}


type Scene struct {
	Topic string `json:"-"`

	Data struct {
		HasData bool   `json:"has_data"`
		Name    string `json:"name"`
		Notes   string `json:"notes"`
		Safes   string `json:"safes"`
	}
}
type Scenes []Scene

func (x *X32) GetScene(i int) Scene {
	var ret Scene

	for range Only.Once {
		// channels start from index 1
		ret.Topic = fmt.Sprintf("/-show/showfile/scene/%.3d/", i)

		msg := x.Call(ret.Topic + "name")
		if msg.Error != nil {
			break
		}
		ret.Data.Name = msg.GetStringValue()

		msg = x.Call(ret.Topic + "hasdata")
		if msg.Error != nil {
			break
		}
		ret.Data.HasData = msg.GetBoolValue()

		msg = x.Call(ret.Topic + "notes")
		if msg.Error != nil {
			break
		}
		ret.Data.Notes = msg.GetStringValue()

		msg = x.Call(ret.Topic + "safes")
		if msg.Error != nil {
			break
		}
		ret.Data.Safes = msg.GetStringValue()
	}

	return ret
}

func (x *X32) GetScenes() Scenes {
	var ret Scenes

	for range Only.Once {
		for c := 0; c < 32; c++ {
			ret = append(ret, x.GetScene(c))
		}
	}

	return ret
}


// func (x *X32) GetPointNamesFromTemplate(template string) api.TemplatePoints {
// 	var ret api.TemplatePoints
//
// 	// for range Only.Once {
// 	// 	if template == "" {
// 	// 		sg.Error = errors.New("no template defined")
// 	// 		break
// 	// 	}
// 	//
// 	// 	ep := sg.GetByStruct(
// 	// 		"WebAppService.queryUserCurveTemplateData",
// 	// 		queryUserCurveTemplateData.RequestData{TemplateID: template},
// 	// 		time.Hour,
// 	// 	)
// 	// 	if sg.Error != nil {
// 	// 		break
// 	// 	}
// 	//
// 	// 	data := queryUserCurveTemplateData.AssertResultData(ep)
// 	// 	for dn, dr := range data.PointsData.Devices {
// 	// 		for _, pr := range dr.Points {
// 	// 			ret = append(ret, api.TemplatePoint{
// 	// 				PsKey:       dn,
// 	// 				PointId:     "p" + pr.PointID,
// 	// 				Description: pr.PointName,
// 	// 				Unit:        pr.Unit,
// 	// 			})
// 	// 		}
// 	// 	}
// 	// }
//
// 	return ret
// }
//
// func (x *X32) GetTemplateData(template string, date string, filter string) error {
// 	// for range Only.Once {
// 	// 	if template == "" {
// 	// 		template = "8042"
// 	// 	}
// 	//
// 	// 	if date == "" {
// 	// 		date = api.NewDateTime("").String()
// 	// 	}
// 	// 	when := api.NewDateTime(date)
// 	//
// 	// 	var psId int64
// 	// 	psId, sg.Error = sg.GetPsId()
// 	// 	if sg.Error != nil {
// 	// 		break
// 	// 	}
// 	//
// 	// 	pointNames := sg.GetPointNamesFromTemplate(template)
// 	// 	if sg.Error != nil {
// 	// 		break
// 	// 	}
// 	//
// 	// 	ep := sg.GetByStruct(
// 	// 		"AppService.queryMutiPointDataList",
// 	// 		queryMutiPointDataList.RequestData{
// 	// 			PsID:           psId,
// 	// 			PsKey:          pointNames.PrintKeys(),
// 	// 			Points:         pointNames.PrintPoints(),
// 	// 			MinuteInterval: "5",
// 	// 			StartTimeStamp: when.GetDayStartTimestamp(),
// 	// 			EndTimeStamp:   when.GetDayEndTimestamp(),
// 	// 		},
// 	// 		DefaultCacheTimeout,
// 	// 	)
// 	// 	if sg.Error != nil {
// 	// 		break
// 	// 	}
// 	//
// 	// 	// data := queryMutiPointDataList.AssertResultData(ep)
// 	// 	data := queryMutiPointDataList.Assert(ep)
// 	// 	table := data.GetDataTable(pointNames)
// 	// 	if table.Error != nil {
// 	// 		sg.Error = table.Error
// 	// 		break
// 	// 	}
// 	//
// 	// 	fn := data.SetFilenamePrefix("%s-%s", when.String(), template)
// 	// 	sg.Error = table.SetFilePrefix(fn)
// 	// 	if sg.Error != nil {
// 	// 		break
// 	// 	}
// 	//
// 	// 	sg.Error = sg.Output(ep, &table, filter)
// 	// 	if sg.Error != nil {
// 	// 		break
// 	// 	}
// 	// }
//
// 	return x.Error
// }
//
// func (x *X32) GetTemplatePoints(template string) error {
// 	for range Only.Once {
// 		if template == "" {
// 			template = "8042"
// 		}
//
// 		table := output.NewTable()
// 		x.Error = table.SetHeader(
// 			"PointStruct Id",
// 			"Description",
// 			"Unit",
// 		)
// 		if x.Error != nil {
// 			break
// 		}
//
// 		ss := x.GetPointNamesFromTemplate(template)
// 		for _, s := range ss {
// 			x.Error = table.AddRow(
// 				api.NameDevicePoint(s.PsKey, s.PointId),
// 				s.Description,
// 				s.Unit,
// 			)
// 			if x.Error != nil {
// 				break
// 			}
// 		}
// 		if x.Error != nil {
// 			break
// 		}
//
// 		table.Print()
// 	}
//
// 	return x.Error
// }
//
// func (x *X32) AllCritical() error {
// 	var ep api.EndPoint
// 	// for range Only.Once {
// 	// 	ep = sg.GetByJson("AppService.powerDevicePointList", "")
// 	// 	if ep.IsError() {
// 	// 		break
// 	// 	}
// 	//
// 	// 	ep = sg.GetByJson("AppService.getPsList", "")
// 	// 	if ep.IsError() {
// 	// 		break
// 	// 	}
// 	//
// 	// 	_getPsList := getPsList.AssertResultData(ep)
// 	// 	psId := _getPsList.GetPsId()
// 	//
// 	// 	ep = sg.GetByJson("AppService.queryDeviceList", fmt.Sprintf(`{"ps_id":"%d"}`, psId))
// 	// 	if ep.IsError() {
// 	// 		break
// 	// 	}
// 	//
// 	// 	ep = sg.GetByJson("AppService.queryDeviceListForApp", fmt.Sprintf(`{"ps_id":"%d"}`, psId))
// 	// 	if ep.IsError() {
// 	// 		break
// 	// 	}
// 	//
// 	// 	ep = sg.GetByJson("WebAppService.showPSView", fmt.Sprintf(`{"ps_id":"%d"}`, psId))
// 	// 	if ep.IsError() {
// 	// 		break
// 	// 	}
// 	//
// 	// 	// ep = sg.GetByJson("AppService.findPsType", fmt.Sprintf(`{"ps_id":"%d"}`, psId))
// 	// 	// if ep.IsError() {
// 	// 	// 	break
// 	// 	// }
// 	//
// 	// 	ep = sg.GetByJson("AppService.getPowerStatistics", fmt.Sprintf(`{"ps_id":"%d"}`, psId))
// 	// 	if ep.IsError() {
// 	// 		break
// 	// 	}
// 	//
// 	// 	ep = sg.GetByJson("AppService.getPsDetail", fmt.Sprintf(`{"ps_id":"%d"}`, psId))
// 	// 	if ep.IsError() {
// 	// 		break
// 	// 	}
// 	//
// 	// 	ep = sg.GetByJson("AppService.getPsDetailWithPsType", fmt.Sprintf(`{"ps_id":"%d"}`, psId))
// 	// 	if ep.IsError() {
// 	// 		break
// 	// 	}
// 	//
// 	// 	ep = sg.GetByJson("AppService.getPsHealthState", fmt.Sprintf(`{"ps_id":"%d"}`, psId))
// 	// 	if ep.IsError() {
// 	// 		break
// 	// 	}
// 	//
// 	// 	ep = sg.GetByJson("AppService.getPsListStaticData", fmt.Sprintf(`{"ps_id":"%d"}`, psId))
// 	// 	if ep.IsError() {
// 	// 		break
// 	// 	}
// 	//
// 	// 	ep = sg.GetByJson("AppService.getPsWeatherList", fmt.Sprintf(`{"ps_id":"%d"}`, psId))
// 	// 	if ep.IsError() {
// 	// 		break
// 	// 	}
// 	//
// 	// 	// ep = sg.GetByJson("AppService.queryAllPsIdAndName", fmt.Sprintf(`{"ps_id":"%d"}`, psId))
// 	// 	// if ep.IsError() {
// 	// 	// 	break
// 	// 	// }
// 	//
// 	// 	// ep = sg.GetByJson("AppService.queryDeviceListByUserId", fmt.Sprintf(`{"ps_id":"%d"}`, psId))
// 	// 	// if ep.IsError() {
// 	// 	// 	break
// 	// 	// }
// 	//
// 	// 	ep = sg.GetByJson("AppService.queryDeviceListForApp", fmt.Sprintf(`{"ps_id":"%d"}`, psId))
// 	// 	if ep.IsError() {
// 	// 		break
// 	// 	}
// 	//
// 	// }
//
// 	x.Error = ep.GetError()
// 	return x.Error
// }
//
// func (x *X32) PrintCurrentStats() error {
// 	// var ep api.EndPoint
// 	// for range Only.Once {
// 	// 	ep = sg.GetByStruct("AppService.getPsList", nil, DefaultCacheTimeout)
// 	// 	if ep.IsError() {
// 	// 		break
// 	// 	}
// 	// 	_getPsList := getPsList.Assert(ep)
// 	// 	psId := _getPsList.GetPsId()
// 	// 	table := _getPsList.GetDataTable()
// 	// 	if table.Error != nil {
// 	// 		sg.Error = table.Error
// 	// 		break
// 	// 	}
// 	//
// 	// 	sg.Error = sg.Output(_getPsList, &table, "")
// 	// 	if sg.Error != nil {
// 	// 		break
// 	// 	}
// 	//
// 	// 	ep = sg.GetByStruct(
// 	// 		"AppService.queryDeviceList",
// 	// 		queryDeviceList.RequestData{PsId: strconv.FormatInt(psId, 10)},
// 	// 		time.Second*60,
// 	// 	)
// 	// 	if sg.Error != nil {
// 	// 		break
// 	// 	}
// 	//
// 	// 	ep2 := queryDeviceList.Assert(ep)
// 	// 	table = ep2.GetDataTable()
// 	// 	if table.Error != nil {
// 	// 		sg.Error = table.Error
// 	// 		break
// 	// 	}
// 	//
// 	// 	sg.Error = sg.Output(ep2, &table, "")
// 	// 	if sg.Error != nil {
// 	// 		break
// 	// 	}
// 	// }
//
// 	return x.Error
// }
//
// // func (sg *X32) QueryDevice(psId int64) queryDeviceList.EndPoint {
// // 	var ret queryDeviceList.EndPoint
// // 	// for range Only.Once {
// // 	// 	if psId == 0 {
// // 	// 		psId, sg.Error = sg.GetPsId()
// // 	// 		if sg.Error != nil {
// // 	// 			break
// // 	// 		}
// // 	// 	}
// // 	//
// // 	// 	// ep = sg.GetByJson("AppService.queryDeviceList", fmt.Sprintf(`{"ps_id":"%d"}`, psId))
// // 	// 	ep := sg.GetByStruct(
// // 	// 		"AppService.queryDeviceList",
// // 	// 		queryDeviceList.RequestData{PsId: strconv.FormatInt(psId, 10)},
// // 	// 		time.Second*60,
// // 	// 	)
// // 	// 	// if sg.Error != nil {
// // 	// 	// 	break
// // 	// 	// }
// // 	//
// // 	// 	ret = queryDeviceList.Assert(ep)
// // 	// }
// //
// // 	return ret
// // }
// //
// // func (sg *X32) QueryPs(psId int64) getPsList.EndPoint {
// // 	var ret getPsList.EndPoint
// // 	// for range Only.Once {
// // 	// 	if psId == 0 {
// // 	// 		psId, sg.Error = sg.GetPsId()
// // 	// 		if sg.Error != nil {
// // 	// 			break
// // 	// 		}
// // 	// 	}
// // 	//
// // 	// 	// ep = sg.GetByJson("AppService.queryDeviceList", fmt.Sprintf(`{"ps_id":"%d"}`, psId))
// // 	// 	ep := sg.GetByStruct(
// // 	// 		"AppService.getPsList",
// // 	// 		getPsList.RequestData{},
// // 	// 		time.Second*60,
// // 	// 	)
// // 	// 	// if sg.Error != nil {
// // 	// 	// 	break
// // 	// 	// }
// // 	//
// // 	// 	ret = getPsList.Assert(ep)
// // 	// }
// //
// // 	return ret
// // }
//
// func (x *X32) GetPointNames(devices ...string) error {
// 	// for range Only.Once {
// 	// 	if len(devices) == 0 {
// 	// 		devices = getPowerDevicePointNames.DeviceTypes
// 	// 	}
// 	// 	for _, dt := range devices {
// 	// 		ep := sg.GetByStruct(
// 	// 			"AppService.getPowerDevicePointNames",
// 	// 			getPowerDevicePointNames.RequestData{DeviceType: dt},
// 	// 			DefaultCacheTimeout,
// 	// 		)
// 	// 		if sg.Error != nil {
// 	// 			break
// 	// 		}
// 	//
// 	// 		ep2 := getPowerDevicePointNames.Assert(ep)
// 	// 		table := ep2.GetDataTable()
// 	// 		if table.Error != nil {
// 	// 			sg.Error = table.Error
// 	// 			break
// 	// 		}
// 	//
// 	// 		sg.Error = sg.Output(ep2, &table, "")
// 	// 		if sg.Error != nil {
// 	// 			break
// 	// 		}
// 	// 	}
// 	// }
//
// 	return x.Error
// }
//
// func (x *X32) GetTemplates() error {
// 	// for range Only.Once {
// 	// 	ep := sg.GetByStruct(
// 	// 		"AppService.getTemplateList",
// 	// 		getTemplateList.RequestData{},
// 	// 		DefaultCacheTimeout,
// 	// 	)
// 	// 	if sg.Error != nil {
// 	// 		break
// 	// 	}
// 	//
// 	// 	ep2 := getTemplateList.Assert(ep)
// 	// 	table := ep2.GetDataTable()
// 	// 	if table.Error != nil {
// 	// 		sg.Error = table.Error
// 	// 		break
// 	// 	}
// 	//
// 	// 	sg.Error = sg.Output(ep2, &table, "")
// 	// 	if sg.Error != nil {
// 	// 		break
// 	// 	}
// 	// }
//
// 	return x.Error
// }
//
// func (x *X32) GetIsolarcloudMqtt(appKey string) error {
// 	// for range Only.Once {
// 	// 	if appKey == "" {
// 	// 		appKey = sg.GetAppKey()
// 	// 	}
// 	//
// 	// 	ep := sg.GetByStruct(
// 	// 		"WebAppService.getMqttConfigInfoByAppkey",
// 	// 		getMqttConfigInfoByAppkey.RequestData{AppKey: appKey},
// 	// 		DefaultCacheTimeout,
// 	// 	)
// 	// 	if sg.Error != nil {
// 	// 		break
// 	// 	}
// 	//
// 	// 	ep2 := getMqttConfigInfoByAppkey.Assert(ep)
// 	// 	table := ep2.GetDataTable()
// 	// 	if table.Error != nil {
// 	// 		sg.Error = table.Error
// 	// 		break
// 	// 	}
// 	//
// 	// 	sg.Error = sg.Output(ep2, &table, "")
// 	// 	if sg.Error != nil {
// 	// 		break
// 	// 	}
// 	// }
//
// 	return x.Error
// }
//
// func (x *X32) GetRealTimeData(psKey string) error {
// 	// for range Only.Once {
// 	// 	if psKey == "" {
// 	// 		var psKeys []string
// 	// 		psKeys, sg.Error = sg.GetPsKeys()
// 	// 		if sg.Error != nil {
// 	// 			break
// 	// 		}
// 	// 		fmt.Printf("%v\n", psKeys)
// 	// 		psKey = strings.Join(psKeys, ",")
// 	// 	}
// 	//
// 	// 	ep := sg.GetByStruct(
// 	// 		"AppService.queryDeviceRealTimeDataByPsKeys",
// 	// 		queryDeviceRealTimeDataByPsKeys.RequestData{PsKeyList: psKey},
// 	// 		DefaultCacheTimeout,
// 	// 	)
// 	// 	if sg.Error != nil {
// 	// 		break
// 	// 	}
// 	//
// 	// 	ep2 := queryDeviceRealTimeDataByPsKeys.Assert(ep)
// 	// 	table := ep2.GetDataTable()
// 	// 	if table.Error != nil {
// 	// 		sg.Error = table.Error
// 	// 		break
// 	// 	}
// 	//
// 	// 	sg.Error = sg.Output(ep2, nil, "")
// 	// 	if sg.Error != nil {
// 	// 		break
// 	// 	}
// 	// }
//
// 	return x.Error
// }
//
// func (x *X32) GetPsDetails(psid string) error {
// 	// for range Only.Once {
// 	// 	var psId int64
// 	// 	if psid == "" {
// 	// 		psId, sg.Error = sg.GetPsId()
// 	// 	} else {
// 	// 		psId, sg.Error = strconv.ParseInt(psid, 10, 64)
// 	// 	}
// 	// 	if sg.Error != nil {
// 	// 		break
// 	// 	}
// 	//
// 	// 	ep := sg.GetByStruct(
// 	// 		"AppService.getPsDetailWithPsType",
// 	// 		getPsDetailWithPsType.RequestData{PsId: strconv.FormatInt(psId, 10)},
// 	// 		DefaultCacheTimeout)
// 	// 	if ep.IsError() {
// 	// 		sg.Error = ep.GetError()
// 	// 		break
// 	// 	}
// 	//
// 	// 	ep2 := getPsDetailWithPsType.Assert(ep)
// 	// 	table := ep2.GetDataTable()
// 	// 	if table.Error != nil {
// 	// 		sg.Error = table.Error
// 	// 		break
// 	// 	}
// 	//
// 	// 	sg.Error = sg.Output(ep2, &table, "")
// 	// 	if sg.Error != nil {
// 	// 		break
// 	// 	}
// 	// }
//
// 	return x.Error
// }
//
// func (x *X32) GetPointData(date string, pointNames api.TemplatePoints) error {
// 	// for range Only.Once {
// 	// 	if len(pointNames) == 0 {
// 	// 		sg.Error = errors.New("no points defined")
// 	// 		break
// 	// 	}
// 	//
// 	// 	if date == "" {
// 	// 		date = api.NewDateTime("").String()
// 	// 	}
// 	// 	when := api.NewDateTime(date)
// 	//
// 	// 	var psId int64
// 	// 	psId, sg.Error = sg.GetPsId()
// 	// 	if sg.Error != nil {
// 	// 		break
// 	// 	}
// 	//
// 	// 	ep := sg.GetByStruct(
// 	// 		"AppService.queryMutiPointDataList",
// 	// 		queryMutiPointDataList.RequestData{
// 	// 			PsID:           psId,
// 	// 			PsKey:          pointNames.PrintKeys(),
// 	// 			Points:         pointNames.PrintPoints(),
// 	// 			MinuteInterval: "5",
// 	// 			StartTimeStamp: when.GetDayStartTimestamp(),
// 	// 			EndTimeStamp:   when.GetDayEndTimestamp(),
// 	// 		},
// 	// 		DefaultCacheTimeout,
// 	// 	)
// 	// 	if sg.Error != nil {
// 	// 		break
// 	// 	}
// 	//
// 	// 	ep2 := queryMutiPointDataList.Assert(ep)
// 	// 	table := ep2.GetDataTable(pointNames)
// 	// 	if table.Error != nil {
// 	// 		sg.Error = table.Error
// 	// 		break
// 	// 	}
// 	//
// 	// 	sg.Error = sg.Output(ep2, &table, "")
// 	// 	if sg.Error != nil {
// 	// 		break
// 	// 	}
// 	// }
//
// 	return x.Error
// }
//
// func (x *X32) SearchPointNames(pns ...string) error {
// 	// for range Only.Once {
// 	// 	table := output.NewTable()
// 	// 	sg.Error = table.SetTitle("")
// 	// 	if sg.Error != nil {
// 	// 		break
// 	// 	}
// 	// 	_ = table.SetHeader(
// 	// 		"DeviceType",
// 	// 		"Id",
// 	// 		"Period",
// 	// 		"Point Id",
// 	// 		"Point Name",
// 	// 		"Show Point Name",
// 	// 		"Translation Id",
// 	// 	)
// 	//
// 	// 	if len(pns) == 0 {
// 	// 		fmt.Println("Searching up to id 1000 within getPowerDevicePointInfo")
// 	// 		for pni := 0; pni < 1000; pni++ {
// 	// 			PrintPause(pni, 20)
// 	//
// 	// 			ep := sg.GetByStruct(
// 	// 				"AppService.getPowerDevicePointInfo",
// 	// 				getPowerDevicePointInfo.RequestData{Id: strconv.FormatInt(int64(pni), 10)},
// 	// 				DefaultCacheTimeout,
// 	// 			)
// 	// 			if sg.Error != nil {
// 	// 				break
// 	// 			}
// 	//
// 	// 			ep2 := getPowerDevicePointInfo.Assert(ep)
// 	// 			table = ep2.AddDataTable(table)
// 	// 			if table.Error != nil {
// 	// 				sg.Error = table.Error
// 	// 				break
// 	// 			}
// 	// 		}
// 	// 		fmt.Println("")
// 	// 	} else {
// 	// 		fmt.Printf("Searching for %v within getPowerDevicePointInfo\n", pns)
// 	// 		for _, pn := range pns {
// 	// 			ep := sg.GetByStruct(
// 	// 				"AppService.getPowerDevicePointInfo",
// 	// 				getPowerDevicePointInfo.RequestData{Id: pn},
// 	// 				DefaultCacheTimeout,
// 	// 			)
// 	// 			if sg.Error != nil {
// 	// 				break
// 	// 			}
// 	//
// 	// 			ep2 := getPowerDevicePointInfo.Assert(ep)
// 	// 			table := ep2.GetDataTable()
// 	// 			if table.Error != nil {
// 	// 				sg.Error = table.Error
// 	// 				break
// 	// 			}
// 	// 		}
// 	// 		fmt.Println("")
// 	// 	}
// 	//
// 	// 	sg.Error = sg.OutputTable(&table)
// 	// 	if sg.Error != nil {
// 	// 		break
// 	// 	}
// 	// }
//
// 	return x.Error
// }
//
// func PrintPause(index int, max int) {
// 	for range Only.Once {
// 		if index == 0 {
// 			fmt.Printf("\n%.3d ", index)
// 			break
// 		}
//
// 		m := math.Mod(float64(index), float64(max))
// 		if m == 0 {
// 			fmt.Printf("PAUSE")
// 			time.Sleep(time.Millisecond * 500)
// 			// fmt.Printf("\r%s%.3d ", strings.Repeat(" ", 4), pni)
// 			fmt.Printf("\r%.3d ", index)
// 		} else {
// 			time.Sleep(time.Millisecond * 100)
// 			fmt.Printf(".")
// 		}
// 	}
// }
//
// func (x *X32) GetPointName(pn string) error {
// 	// for range Only.Once {
// 	// 	ep := sg.GetByStruct(
// 	// 		"AppService.getPowerDevicePointInfo",
// 	// 		getPowerDevicePointInfo.RequestData{Id: pn},
// 	// 		DefaultCacheTimeout,
// 	// 	)
// 	// 	if sg.Error != nil {
// 	// 		break
// 	// 	}
// 	//
// 	// 	ep2 := getPowerDevicePointInfo.Assert(ep)
// 	// 	table := ep2.GetDataTable()
// 	// 	if table.Error != nil {
// 	// 		sg.Error = table.Error
// 	// 		break
// 	// 	}
// 	//
// 	// 	// table2 := ep2.GetData()
// 	// 	// fmt.Printf("%v\n", table2)
// 	//
// 	// 	sg.Error = sg.Output(ep2, &table, "")
// 	// 	if sg.Error != nil {
// 	// 		break
// 	// 	}
// 	// }
//
// 	return x.Error
// }
//
// func (x *X32) GetPsId() (int64, error) {
// 	var ret int64
//
// 	// for range Only.Once {
// 	//
// 	// 	ep := sg.GetByStruct("AppService.getPsList", nil, DefaultCacheTimeout)
// 	// 	if ep.IsError() {
// 	// 		sg.Error = ep.GetError()
// 	// 		break
// 	// 	}
// 	//
// 	// 	_getPsList := getPsList.AssertResultData(ep)
// 	// 	ret = _getPsList.GetPsId()
// 	// }
//
// 	return ret, x.Error
// }
//
// func (x *X32) GetPsName() (string, error) {
// 	var ret string
//
// 	// for range Only.Once {
// 	//
// 	// 	ep := sg.GetByStruct("AppService.getPsList", nil, DefaultCacheTimeout)
// 	// 	if ep.IsError() {
// 	// 		sg.Error = ep.GetError()
// 	// 		break
// 	// 	}
// 	//
// 	// 	_getPsList := getPsList.AssertResultData(ep)
// 	// 	ret = _getPsList.GetPsName()
// 	// }
//
// 	return ret, x.Error
// }
//
// func (x *X32) GetPsModel() (string, error) {
// 	var ret string
//
// 	// for range Only.Once {
// 	// 	var psId int64
// 	// 	psId, sg.Error = sg.GetPsId()
// 	// 	if sg.Error != nil {
// 	// 		break
// 	// 	}
// 	//
// 	// 	ep := sg.GetByStruct(
// 	// 		"AppService.getPsDetailWithPsType",
// 	// 		getPsDetailWithPsType.RequestData{PsId: strconv.FormatInt(psId, 10)},
// 	// 		DefaultCacheTimeout)
// 	// 	if ep.IsError() {
// 	// 		sg.Error = ep.GetError()
// 	// 		break
// 	// 	}
// 	//
// 	// 	ep2 := getPsDetailWithPsType.Assert(ep)
// 	// 	ret = ep2.GetDeviceName()
// 	// }
//
// 	return ret, x.Error
// }
//
// func (x *X32) GetPsSerial() (string, error) {
// 	var ret string
//
// 	// for range Only.Once {
// 	// 	var psId int64
// 	// 	psId, sg.Error = sg.GetPsId()
// 	// 	if sg.Error != nil {
// 	// 		break
// 	// 	}
// 	//
// 	// 	ep := sg.GetByStruct(
// 	// 		"AppService.getPsDetailWithPsType",
// 	// 		getPsDetailWithPsType.RequestData{PsId: strconv.FormatInt(psId, 10)},
// 	// 		DefaultCacheTimeout)
// 	// 	if ep.IsError() {
// 	// 		sg.Error = ep.GetError()
// 	// 		break
// 	// 	}
// 	//
// 	// 	ep2 := getPsDetailWithPsType.Assert(ep)
// 	// 	ret = ep2.GetDeviceSerial()
// 	// }
//
// 	return ret, x.Error
// }
//
// func (x *X32) GetPsKeys() ([]string, error) {
// 	var ret []string
//
// 	// for range Only.Once {
// 	// 	var psId int64
// 	// 	psId, sg.Error = sg.GetPsId()
// 	// 	if sg.Error != nil {
// 	// 		break
// 	// 	}
// 	//
// 	// 	ep := sg.GetByStruct(
// 	// 		"AppService.getPsDetailWithPsType",
// 	// 		getPsDetailWithPsType.RequestData{PsId: strconv.FormatInt(psId, 10)},
// 	// 		DefaultCacheTimeout)
// 	// 	if ep.IsError() {
// 	// 		sg.Error = ep.GetError()
// 	// 		break
// 	// 	}
// 	//
// 	// 	ep2 := getPsDetailWithPsType.Assert(ep)
// 	// 	ret = ep2.GetPsKeys()
// 	// }
//
// 	return ret, x.Error
// }