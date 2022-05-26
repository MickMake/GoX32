package Behringer

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/MickMake/GoX32/Only"
)


type ChData struct {
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
	Fader    float64 `json:"fader"`
}

func (c *ChData) Json() string {
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

func (c ChData) String() string {
	var ret string

	ret += fmt.Sprintf("Name:\t%s\n", c.Name)
	ret += fmt.Sprintf("Colour:\t%s\n", c.Colour)
	ret += fmt.Sprintf("Icon:\t%s\n", c.Icon)
	ret += fmt.Sprintf("Source:\t%s\n", c.Source)
	ret += fmt.Sprintf("Gain:\t%f\n", c.Gain)
	ret += fmt.Sprintf("Mute:\t%t\n", c.Mute)
	ret += fmt.Sprintf("Fader:\t%f\n", c.Fader)
	ret += fmt.Sprintf("Solo:\t%t\n", c.Solo)
	ret += fmt.Sprintf("Trim:\t%f\n", c.Trim)
	ret += fmt.Sprintf("Phantom:\t%t\n", c.Phantom)
	ret += fmt.Sprintf("Selected:\t%t\n", c.Selected)

	return ret
}


func (x *X32) GetAllInfo() error {
	for range Only.Once {
		// return nil

		x.Error = x.EmitStatus()
		if x.Error != nil {
			break
		}

		x.Error = x.EmitInfo()
		if x.Error != nil {
			break
		}

		x.Error = x.EmitXinfo()
		if x.Error != nil {
			break
		}

		// x.Error = x.EmitNode()
		// if x.Error != nil {
		// 	break
		// }

		// x.Error = x.EmitId()
		// if x.Error != nil {
		// 	break
		// }

		// x.GetDeskName()
	}

	return x.Error
}

const CmdStatus = "/status"
func (x *X32) EmitStatus() error  { return x.Emit(CmdStatus) }
func (x *X32) GetStatus() *Message { return x.Call(CmdStatus) }

const CmdInfo = "/info"
func (x *X32) EmitInfo() error  { return x.Emit(CmdInfo) }
func (x *X32) GetInfo() *Message { return x.Call(CmdInfo) }

const CmdXinfo = "/xinfo"
func (x *X32) EmitXinfo() error  { return x.Emit(CmdXinfo) }
func (x *X32) GetXinfo() *Message { return x.Call(CmdXinfo) }

const CmdShowDump = "/showdump"
func (x *X32) EmitShowDump() error  { return x.Emit(CmdShowDump) }
func (x *X32) GetShowDump() *Message { return x.Call(CmdShowDump) }

const CmdNode = "/node"
func (x *X32) EmitNode() error  { return x.Emit(CmdNode) }
func (x *X32) GetNode() *Message { return x.Call(CmdNode) }

const CmdId = "/-prefs/??????"
func (x *X32) EmitId() error  { return x.Emit(CmdId) }
func (x *X32) GetId() *Message { return x.Call(CmdId) }

func (x *X32) GetDeskName() *Message { return x.Call("/-prefs/name") }


func (x *X32) Set(address string, value any) error {

	for range Only.Once {
		if x.Debug {
			fmt.Printf("# Set() - address: %v, args: %v\n", address, value)
		}

		point := x.Points.Resolve(address)
		if point == nil {
			x.Error = errors.New(fmt.Sprintf("point '%s' not found", address))
			break
		}
		value = point.Convert.SetValue(value)

		x.Error = x.Emit(address, value)
		if x.Error != nil {
			break
		}
	}

	return x.Error
}
