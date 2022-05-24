package Behringer

import (
	"encoding/json"
	"fmt"
	"github.com/MickMake/GoX32/Only"
)


type Scenes []Scene

type Scene struct {
	Topic string `json:"-"`

	Data struct {
		HasData bool   `json:"has_data"`
		Name    string `json:"name"`
		Notes   string `json:"notes"`
		Safes   string `json:"safes"`
	}
}

func (s Scene) String() string {
	return fmt.Sprintf("Type:\t%T\nMQTT Topic:\t%s\nName:\t%s\nHas Data:\t%t\nSafes:\t%s\nNotes:\t%s",
		s,
		s.Topic,
		s.Data.Name,
		s.Data.HasData,
		s.Data.Safes,
		s.Data.Notes,
	)
}

func (s *Scene) Json() string {
	var ret string

	for range Only.Once {
		j, err := json.Marshal(s)
		if err != nil {
			break
		}
		ret = string(j)
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

func (x *X32) GetScene(i int) Scene {
	var ret Scene

	for range Only.Once {
		// channels start from index 1
		ret.Topic = fmt.Sprintf("/-show/showfile/scene/%.3d/", i)

		msg := x.Call(ret.Topic + "name")
		if msg.Error != nil {
			break
		}
		ret.Data.Name = msg.GetValueString()

		msg = x.Call(ret.Topic + "hasdata")
		if msg.Error != nil {
			break
		}
		ret.Data.HasData = msg.GetValueBool()

		msg = x.Call(ret.Topic + "notes")
		if msg.Error != nil {
			break
		}
		ret.Data.Notes = msg.GetValueString()

		msg = x.Call(ret.Topic + "safes")
		if msg.Error != nil {
			break
		}
		ret.Data.Safes = msg.GetValueString()
	}

	return ret
}
