package mmHa

import (
	"encoding/json"
	"fmt"
	"github.com/MickMake/GoX32/Only"
)


func (m *Mqtt) SwitchPublishConfig(config EntityConfig) error {

	for range Only.Once {
		config.FixConfig()
		if !config.IsSwitch() {
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

		payload := Switch {
			Device:                 device,
			Name:                   JoinStrings(m.Device.Name, config.ParentName, config.Name),
			StateTopic:             JoinStringsForTopic(m.binarySensorPrefix, config.StateTopic, "state"),
			// StateClass:             config.StateClass,
			UniqueId:               config.StateTopic,
			// UnitOfMeasurement:      config.Units,
			// DeviceClass:            config.DeviceClass,
			Qos:                    0,
			Retain:                 true,
			// ExpireAfter:            0,
			// Encoding:               "utf-8",
			// EnabledByDefault:       true,
			// LastResetValueTemplate: config.LastResetValueTemplate,
			// LastReset:              config.LastReset,
			ValueTemplate:          config.ValueTemplate,
			Icon:                   config.Icon,
			PayloadOn:              "ON",
			PayloadOff:             "OFF",

			StateOn:                "ON",
			StateOff:               "OFF",
			CommandTopic:           JoinStringsForTopic(m.binarySensorPrefix, config.StateTopic, "cmd"),

			// Availability:           &Availability {
			// 	PayloadAvailable:    "",
			// 	PayloadNotAvailable: "",
			// 	Topic:               "",
			// 	ValueTemplate:       "",
			// },
			// AvailabilityMode:       "",
			// AvailabilityTemplate:   "",
			// AvailabilityTopic:      "",
			// JSONAttributesTemplate: "",
			// JSONAttributesTopic:    "",
			// Optimistic:             false,
			// PayloadAvailable:       "",
			// PayloadNotAvailable:    "",
		}

		ct := JoinStringsForTopic(m.switchPrefix, uid, "config")
		t := m.client.Publish(ct, 0, true, payload.Json())
		if !t.WaitTimeout(m.Timeout) {
			m.err = t.Error()
		}
	}

	return m.err
}

func (m *Mqtt) SwitchPublishValue(config EntityConfig) error {

	for range Only.Once {
		if !config.IsSwitch() {
			break
		}

		st := JoinStringsForId(m.Device.Name, config.ParentId, config.Name)
		cs := make(ValueMap)
		cs["last_reset"] = m.GetLastReset(JoinStringsForId(config.ParentId, config.Name))
		cs[config.ValueName] = config.Value

		// payload := MqttState {
		// 	LastReset: m.GetLastReset(config.UniqueId),
		// 	Value:     config.Value,
		// }

		st = JoinStringsForTopic(m.binarySensorPrefix, st, "state")
		t := m.client.Publish(st, 0, true, cs.Json())
		if !t.WaitTimeout(m.Timeout) {
			m.err = t.Error()
		}
	}

	return m.err
}

func (m *Mqtt) SwitchPublishValues(configs []EntityConfig) error {
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
				topic = JoinStringsForTopic(m.switchPrefix, config.StateTopic, "state")
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


type Switch struct {
	AvailabilityTopic      string `json:"avty_t,omitempty"`
	CommandTopic           string `json:"cmd_t"`
	Device                 Device `json:"dev,omitempty"`
	Icon                   string `json:"ic,omitempty"`
	JSONAttributesTemplate string `json:"json_attr_tpl,omitempty"`
	JSONAttributesTopic    string `json:"json_attr_t,omitempty"`
	Name                   string `json:"name,omitempty"`
	Optimistic             bool   `json:"opt,omitempty"`
	PayloadAvailable       string `json:"pl_avail,omitempty"`
	PayloadNotAvailable    string `json:"pl_not_avail,omitempty"`
	PayloadOff             string `json:"pl_off,omitempty"`
	PayloadOn              string `json:"pl_on,omitempty"`
	Qos                    int    `json:"qos,omitempty"`
	Retain                 bool   `json:"ret,omitempty"`
	StateOff               string `json:"stat_off,omitempty"`
	StateOn                string `json:"stat_on,omitempty"`
	StateTopic             string `json:"stat_t,omitempty"`
	UniqueId               string `json:"uniq_id,omitempty"`
	ValueTemplate          string `json:"val_tpl,omitempty"`

	// CommandFunc func(mqtt.Message, mqtt.Client) `json:"-"`
	// StateFunc   func() string                   `json:"-"`
	//
	// UpdateInterval  float64 `json:"-"`
	// ForceUpdateMQTT bool    `json:"-"`
	//
	// messageHandler mqtt.MessageHandler
}

func (c *Switch) Json() string {
	j, _ := json.Marshal(*c)
	return string(j)
}
