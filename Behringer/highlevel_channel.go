package Behringer

import (
	"encoding/json"
	"fmt"
	"github.com/MickMake/GoX32/Only"
)


type Channels []Channel	// We're keeping a 1:1 mapping of array -> channel numbers.

type Channel struct {
	Topic string `json:"-"`
	Data  ChData `json:"data"`
}

func (c Channel) String() string {
	return fmt.Sprintf("Type:\t%T\nMQTT Topic:\t%s\n%s", c, c.Topic, c.Data)
}

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


// func (x *X32) GetChannels() []MessageMap {
// 	var ret []MessageMap
//
// 	for range Only.Once {
// 		for c := 0; c < 32; c++ {
// 			// ret = append(ret, x.GetChannel(c))
// 		}
// 	}
//
// 	return ret
// }

func (c *Channel) GetInfo(i int) map[string]string {
	t := fmt.Sprintf("/ch/%.2d/", i)	// channels start from index 1
	ret := map[string]string {
		"Name":        t + "config/name",
		"Colour":      t + "config/color",
		"Source":      t + "config/source",
		"Icon":        t + "config/icon",
		"Preamp Trim": t + "preamp/trim",
		"DCA Group": t + "grp/dca",
		"Mute Group": t + "grp/mute",
		"Gate On": t + "gate/on",
		"Compression On": t + "dyn/on",
		"EQ On": t + "eq/on",
		// "": topic + "preamp/hpon",
		// "": topic + "mix/01/level",
		"Mix On":  t + "mix/on",
		"Fader":   t + "mix/fader",
		"Gain":    fmt.Sprintf("/headamp/%.3d/gain", i - 1), // headamps start from index 0
		"Phantom On": fmt.Sprintf("/headamp/%.3d/phantom", i - 1),
	}

	return ret
}

func (x *X32) GetChannel(i int) MessageMap {
	ret := make(MessageMap)

	for range Only.Once {
		channel := Channel{}
		topics := channel.GetInfo(i)

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

func (x *X32) ChannelCount() []int {
	var ret []int

	for range Only.Once {
		for i := 1; i <= 32; i++ {
			ret = append(ret, i)
		}
	}

	return ret
}
