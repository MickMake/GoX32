package Behringer

import (
	"encoding/json"
	"fmt"
	"github.com/MickMake/GoX32/Only"
)


type Auxes []Aux

type Aux struct {
	Topic string `json:"-"`
	Data  ChData `json:"data"`
}

func (a Aux) String() string {
	return fmt.Sprintf("Type:\t%T\nMQTT Topic:\t%s\n%s", a, a.Topic, a.Data)
}

func (a *Aux) Json() string {
	var ret string

	for range Only.Once {
		j, err := json.Marshal(a)
		if err != nil {
			break
		}
		ret = string(j)
	}

	return ret
}


// func (x *X32) GetAuxes() Auxes {
// 	var ret Auxes
//
// 	for range Only.Once {
// 		for c := 0; c < 8; c++ {
// 			// ret = append(ret, x.GetAux(c))
// 		}
// 	}
//
// 	return ret
// }

func (x *X32) GetAux(i int) MessageMap {
	ret := make(MessageMap)

	for range Only.Once {
		t := fmt.Sprintf("/auxin/%.2d/", i)	// channels start from index 1
		topics := map[string]string {
			"Name":        t + "config/name",
			"Colour":      t + "config/color",
			"Source":      t + "config/source",
			"Icon":        t + "config/icon",
			"Preamp Trim": t + "preamp/trim",
			"DCA Group": t + "grp/dca",
			"Mute Group": t + "grp/mute",
			// "Gate On": t + "gate/on",
			// "Compression On": t + "dyn/on",
			"EQ On": t + "eq/on",
			// "Mix On":  t + "mix/on",
			"Fader":   t + "mix/fader",
			// "Gain":    fmt.Sprintf("/headamp/%.3d/gain", i), // headamps start from index 0
			// "Phantom On": fmt.Sprintf("/headamp/%.3d/phantom", i),
		}

		for name, topic := range topics {
			msg := x.Call(topic)
			if msg.Error != nil {
				x.Error = msg.Error
				break
			}
			msg.Name = name
			ret[topic] = msg
		}
	}

	return ret
}

func (x *X32) AuxCount() []int {
	var ret []int

	for range Only.Once {
		for i := 1; i <= 8; i++ {
			ret = append(ret, i)
		}
	}

	return ret
}
