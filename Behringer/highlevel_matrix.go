package Behringer

import (
	"encoding/json"
	"fmt"
	"github.com/MickMake/GoX32/Only"
)


type Matrices []Matrix

type Matrix struct {
	Topic string `json:"-"`
	Data  ChData `json:"data"`
}

func (m Matrix) String() string {
	return fmt.Sprintf("Type:\t%T\nMQTT Topic:\t%s\n%s", m, m.Topic, m.Data)
}

func (m *Matrix) Json() string {
	var ret string

	for range Only.Once {
		j, err := json.Marshal(m)
		if err != nil {
			break
		}
		ret = string(j)
	}

	return ret
}


// func (x *X32) GetMatrices() Matrices {
// 	var ret Matrices
//
// 	for range Only.Once {
// 		for c := 0; c < 6; c++ {
// 			// ret = append(ret, x.GetMatrix(c))
// 		}
// 	}
//
// 	return ret
// }

func (x *X32) GetMatrix(i int) MessageMap {
	ret := make(MessageMap)

	for range Only.Once {
		t := fmt.Sprintf("/mtx/%.2d/", i)	// channels start from index 1
		topics := map[string]string {
			"Name":        t + "config/name",
			"Colour":      t + "config/color",
			// "Source":      t + "config/source",
			"Icon":        t + "config/icon",
			// "Preamp Trim": t + "preamp/trim",
			// "": topic + "preamp/hpon",
			// "": topic + "mix/01/level",
			"Mix On":  t + "mix/on",
			"Fader":   t + "mix/fader",
			"Gain":    fmt.Sprintf("/headamp/%.3d/gain", i), // headamps start from index 0
			"Phantom": fmt.Sprintf("/headamp/%.3d/phantom", i),
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

func (x *X32) MatrixCount() []int {
	var ret []int

	for range Only.Once {
		for i := 1; i <= 6; i++ {
			ret = append(ret, i)
		}
	}

	return ret
}
