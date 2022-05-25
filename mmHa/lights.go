package mmHa

import (
	"encoding/json"
	"fmt"
	"github.com/MickMake/GoX32/Only"
)


func (m *Mqtt) PublishLightConfig(config EntityConfig) error {
// func (m *Mqtt) PublishLightConfig(id string, name string, subName string, units string, valueName string, class string) error {
	for range Only.Once {
		config.FixConfig()
		if !config.IsLight() {
			break
		}

		device := m.Device
		device.Name = JoinStrings(m.Device.Name, config.ParentId)
		device.Connections = [][]string {
			{ m.Device.Name, JoinStringsForId(m.Device.Name, config.ParentId) },
			{ JoinStringsForId(m.Device.Name, config.ParentId), JoinStringsForId(m.Device.Name, config.ParentId, config.Name) },
		}
		device.Identifiers = []string { JoinStringsForId(m.Device.Name, config.ParentId) }
		if config.StateTopic != "" {
			config.StateTopic = JoinStringsForId(m.Device.Name, config.ParentId, config.StateTopic)
			// config.StateTopic = JoinStringsForId(m.Device.Name, config.ParentName, config.Name, config.UniqueId),
		} else {
			config.StateTopic = JoinStringsForId(m.Device.Name, config.ParentId, config.Name)
			// config.StateTopic = JoinStringsForId(m.Device.Name, config.ParentName, config.Name, config.UniqueId),
		}
		uid := JoinStringsForId(m.Device.Name, config.ParentId, config.Name)

		payload := Light {
			Device:                 device,
			EnabledByDefault:       true,
			Encoding:               "utf-8",
			EntityCategory:         "",
			Icon:                   config.Icon,
			Name:                   JoinStrings(m.Device.Name, config.ParentName, config.Name),
			// ObjectId:               config.UniqueId,
			Qos:                    0,
			StateTopic:             JoinStringsForTopic(m.lightPrefix, config.StateTopic, StateTopicSuffix),
			UniqueId:               config.StateTopic,

			CommandTopic:           JoinStringsForTopic(m.lightPrefix, config.StateTopic, CmdTopicSuffix),
			StateClass:             config.StateClass,
			Retain:                 true,
			ValueTemplate:          config.ValueTemplate,
		}

		ct := JoinStringsForTopic(m.lightPrefix, uid, ConfigTopicSuffix)
		t := m.client.Publish(ct, 0, true, payload.Json())
		if !t.WaitTimeout(m.Timeout) {
			m.err = t.Error()
		}
	}

	return m.err
}

func (m *Mqtt) PublishLightValue(config EntityConfig) error {

	for range Only.Once {
		if !config.IsLight() {
			break
		}

		st := JoinStringsForId(m.Device.Name, config.ParentId, config.Name)
		cs := make(ValueMap)
		// cs["last_reset"] = m.GetLastReset(JoinStringsForId(config.ParentId, config.Name))
		cs[config.ValueName] = config.Value

		st = JoinStringsForTopic(m.lightPrefix, st, StateTopicSuffix)
		t := m.client.Publish(st, 0, true, cs.Json())
		if !t.WaitTimeout(m.Timeout) {
			m.err = t.Error()
		}
	}

	return m.err
}

func (m *Mqtt) PublishLightValues(configs []EntityConfig) error {
	for range Only.Once {
		if len(configs) == 0 {
			break
		}

		cs := make(map[string]Fields)
		topic := ""
		for i, config := range configs {
			if !config.IsBinarySensor() {
				continue
			}

			if config.StateTopic != "" {
				config.StateTopic = JoinStringsForId(m.Device.Name, config.ParentId, config.StateTopic)
			} else {
				config.StateTopic = JoinStringsForId(m.Device.Name, config.ParentId, config.Name)
				// config.StateTopic = JoinStringsForId(m.Device.Name, config.ParentName, config.Name, config.UniqueId),
			}

			if topic == "" {
				topic = JoinStringsForTopic(m.lightPrefix, config.StateTopic, StateTopicSuffix)
			}

			if config.ValueName == "" {
				config.ValueName = fmt.Sprintf("value%d", i)
			}
			if _, ok := cs[config.StateTopic]; !ok {
				cs[config.StateTopic] = make(Fields)
			}
			cs[config.StateTopic][config.ValueName] = config.Value
		}

		for _, c := range cs {
			j, _ := json.Marshal(c)
			t := m.client.Publish(topic, 0, true, string(j))
			if !t.WaitTimeout(m.Timeout) {
				m.err = t.Error()
			}
		}
	}
	return m.err
}

// func (m *Mqtt) PublishLight(subtopic string, payload interface{}) error {
// 	for range Only.Once {
// 		t := m.client.Publish(JoinStringsForTopic(m.lightPrefix, subtopic), 0, true, payload)
// 		if !t.WaitTimeout(m.Timeout) {
// 			m.err = t.Error()
// 		}
// 	}
// 	return m.err
// }
//
// func (m *Mqtt) PublishLightState(topic string, payload interface{}) error {
// 	for range Only.Once {
// 		topic = JoinStringsForId(m.EntityPrefix, m.Device.Name, topic)
// 		t := m.client.Publish(JoinStringsForTopic(m.lightPrefix, topic, "state"), 0, true, payload)
// 		if !t.WaitTimeout(m.Timeout) {
// 			m.err = t.Error()
// 		}
// 	}
// 	return m.err
// }


type Light struct {
	// Common fields.
	Availability             *Availability `json:"availability,omitempty" required:"false"`
	AvailabilityMode         string       `json:"availability_mode,omitempty" required:"false"`
	AvailabilityTemplate     string       `json:"availability_template,omitempty" required:"false"`
	AvailabilityTopic        string       `json:"availability_topic,omitempty" required:"false"`
	Device                   Device       `json:"device,omitempty" required:"false"`
	EnabledByDefault         bool         `json:"enabled_by_default,omitempty" required:"false"`
	Encoding                 string       `json:"encoding,omitempty" required:"false"`
	EntityCategory           string       `json:"entity_category,omitempty" required:"false"`
	Icon                     string       `json:"icon,omitempty" required:"false"`
	JsonAttributesTemplate   string       `json:"json_attributes_template,omitempty" required:"false"`
	JsonAttributesTopic      string       `json:"json_attributes_topic,omitempty" required:"false"`
	Name                     string       `json:"name,omitempty" required:"false"`
	ObjectId                 string       `json:"object_id,omitempty" required:"false"`
	PayloadAvailable         string       `json:"payload_available,omitempty" required:"false"`
	PayloadNotAvailable      string       `json:"payload_not_available,omitempty" required:"false"`
	Qos                      int          `json:"qos,omitempty" required:"false"`
	StateTopic               string       `json:"state_topic" required:"true"`
	UniqueId                 string       `json:"unique_id,omitempty" required:"false"`
	ValueTemplate            string       `json:"value_template,omitempty" required:"false"`

	// Less common fields.
	//  CommandTemplate          string       `json:"command_template,omitempty" required:"false"`
	CommandTopic             string       `json:"command_topic,omitempty" required:"true"`
	//  DeviceClass              string       `json:"device_class,omitempty" required:"false"`
	//	ExpireAfter              int          `json:"expire_after,omitempty" required:"false"`
	//	ForceUpdate              bool         `json:"force_update,omitempty" required:"false"`
	//	LastResetValueTemplate string       `json:"last_reset_value_template,omitempty" required:"false"`
	//	OffDelay               int          `json:"off_delay,omitempty" required:"false"`
	//  Options                  []string     `json:"options,omitempty" required:"true"`
	Optimistic               bool         `json:"optimistic,omitempty" required:"false"`
	PayloadOff               string       `json:"payload_off,omitempty" required:"false"`
	PayloadOn                string       `json:"payload_on,omitempty" required:"false"`
	Retain                   bool         `json:"retain,omitempty" required:"false"`
	StateClass               string       `json:"state_class,omitempty" required:"false"`
	//	StateOff                 string       `json:"state_off,omitempty" required:"false"`
	//	StateOn                  string       `json:"state_on,omitempty" required:"false"`
	//	UnitOfMeasurement        string       `json:"unit_of_measurement,omitempty" required:"false"`

	// Unique fields.
	BrightnessCommandTopic   string   `json:"brightness_command_topic,omitempty"`
	BrightnessCommandTemplate   string   `json:"brightness_command_template,omitempty"`
	BrightnessScale          uint8    `json:"brightness_scale,omitempty"`
	BrightnessStateTopic     string   `json:"brightness_state_topic,omitempty"`
	BrightnessValueTemplate  string   `json:"brightness_value_template,omitempty"`
	ColorModeStateTopic      string   `json:"color_mode_state_topic,omitempty"`
	ColorModeValueTemplate   string   `json:"color_mode_value_template,omitempty"`
	ColorTempCommandTemplate string   `json:"color_temp_command_template,omitempty"`
	ColorTempCommandTopic    string   `json:"color_temp_command_topic,omitempty"`
	ColorTempStateTopic      string   `json:"color_temp_state_topic,omitempty"`
	ColorTempValueTemplate   string   `json:"color_temp_value_template,omitempty"`
	EffectCommandTopic       string   `json:"effect_command_topic,omitempty"`
	EffectCommandTemplate    string   `json:"effect_command_template,omitempty"`
	EffectList               []string `json:"effect_list,omitempty"`
	EffectStateTopic         string   `json:"effect_state_topic,omitempty"`
	EffectValueTemplate      string   `json:"effect_value_template,omitempty"`
	HsCommandTopic           string   `json:"hs_command_topic,omitempty"`
	HsStateTopic             string   `json:"hs_state_topic,omitempty"`
	HsValueTemplate          string   `json:"hs_value_template,omitempty"`
	MaxMireds                int      `json:"max_mireds,omitempty"`
	MinMireds                int      `json:"min_mireds,omitempty"`
	OnCommandType            string   `json:"on_command_type,omitempty"`
	RgbCommandTemplate       string   `json:"rgb_command_template,omitempty"`
	RgbCommandTopic          string   `json:"rgb_command_topic,omitempty"`
	RgbStateTopic            string   `json:"rgb_state_topic,omitempty"`
	RgbValueTemplate         string   `json:"rgb_value_template,omitempty"`
	RgbwCommandTemplate      string   `json:"rgbw_command_template,omitempty"`
	RgbwCommandTopic         string   `json:"rgbw_command_topic,omitempty"`
	RgbwStateTopic           string   `json:"rgbw_state_topic,omitempty"`
	RgbwValueTemplate        string   `json:"rgbw_value_template,omitempty"`
	RgbwwCommandTemplate     string   `json:"rgbww_command_template,omitempty"`
	RgbwwCommandTopic        string   `json:"rgbww_command_topic,omitempty"`
	RgbwwStateTopic          string   `json:"rgbww_state_topic,omitempty"`
	RgbwwValueTemplate       string   `json:"rgbww_value_template,omitempty"`
	Schema                   string   `json:"schema,omitempty"`
	StateValueTemplate       string   `json:"state_value_template,omitempty"`
	WhiteCommandTopic        string   `json:"white_command_topic,omitempty"`
	WhiteScale               int      `json:"white_scale,omitempty"`
	XyCommandTopic           string   `json:"xy_command_topic,omitempty"`
	XyStateTopic             string   `json:"xy_state_topic,omitempty"`
	XyValueTemplate          string   `json:"xy_value_template,omitempty"`

}
func (c *Light) Json() string {
	j, _ := json.Marshal(*c)
	return string(j)
}


// {
//	"brightness": true,
//	"cmd_t": "homeassistant/light/cbus_20/set",
//	"device": {
//		"connections": [
//			[
//				"cbus_group_address",
//				"20"
//			]
//		],
//		"identifiers": [
//			"cbus_light_20"
//		],
//		"manufacturer": "Clipsal",
//		"model": "C-Bus Lighting Application",
//		"name": "C-Bus Light 020",
//		"sw_version": "cmqttd https://github.com/micolous/cbus",
//		"via_device": "cmqttd"
//	},
//	"name": "C-Bus Light 020",
//	"schema": "json",
//	"stat_t": "homeassistant/light/cbus_20/state",
//	"unique_id": "cbus_light_20"
// }
//
// type LightConfig struct {
// 	Name        string      `json:"name"`
// 	UniqueId    string      `json:"unique_id"`
// 	CmdT        string      `json:"cmd_t"`
// 	StatT       string      `json:"stat_t"`
// 	Schema      string      `json:"schema"`
// 	Brightness  bool        `json:"brightness"`
// 	LightDevice LightDevice `json:"device"`
// }
// type LightDevice struct {
// 	Identifiers  []string   `json:"identifiers"`
// 	Connections  [][]string `json:"connections"`
// 	SwVersion    string     `json:"sw_version"`
// 	Name         string     `json:"name"`
// 	Manufacturer string     `json:"manufacturer"`
// 	Model        string     `json:"model"`
// 	ViaDevice    string     `json:"via_device"`
// }

// {
//	"brightness": 255,
//	"cbus_source_addr": 7,
//	"state": "ON",
//	"transition": 0
// }

// type LightState struct {
// 	State          string `json:"state"`
// 	Brightness     int    `json:"brightness"`
// 	Transition     int    `json:"transition"`
// 	CbusSourceAddr int    `json:"cbus_source_addr"`
// }