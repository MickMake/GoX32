package mmHa

import (
	"encoding/json"
	"fmt"
	"github.com/MickMake/GoX32/Only"
)


func (m *Mqtt) BinarySensorPublishConfig(config EntityConfig) error {

	for range Only.Once {
		config.FixConfig()
		if !config.IsBinarySensor() {
			break
		}

		// LastReset := m.GetLastReset(config.UniqueId)
		// LastResetValueTemplate := ""
		// if LastReset != "" {
		// 	LastResetValueTemplate = "{{ value_json.last_reset | as_datetime() }}"
		// }

		device := m.Device
		device.Name = JoinStrings(m.Device.Name, config.ParentId)
		device.Connections = [][]string {
			{ m.Device.Name, JoinStringsForId(m.Device.Name, config.ParentId) },
			{ JoinStringsForId(m.Device.Name, config.ParentId), JoinStringsForId(m.Device.Name, config.ParentId, config.Name) },
		}
		device.Identifiers = []string{ JoinStringsForId(m.Device.Name, config.ParentId) }
		if config.StateTopic != "" {
			config.StateTopic = JoinStringsForId(m.Device.Name, config.ParentId, config.StateTopic)
			// config.StateTopic = JoinStringsForId(m.Device.Name, config.ParentName, config.Name, config.UniqueId),
		} else {
			config.StateTopic = JoinStringsForId(m.Device.Name, config.ParentId, config.Name)
			// config.StateTopic = JoinStringsForId(m.Device.Name, config.ParentName, config.Name, config.UniqueId),
		}
		uid := JoinStringsForId(m.Device.Name, config.ParentId, config.Name)

		payload := BinarySensor {
			Device:                 device,
			Name:                   JoinStrings(m.Device.Name, config.ParentName, config.Name),
			StateTopic:             JoinStringsForTopic(m.binarySensorPrefix, config.StateTopic, "state"),
			StateClass:             config.StateClass,
			UniqueId:               uid,
			UnitOfMeasurement:      config.Units,
			DeviceClass:            config.DeviceClass,
			Qos:                    0,
			ForceUpdate:            true,
			ExpireAfter:            0,
			Encoding:               "utf-8",
			EnabledByDefault:       true,
			LastResetValueTemplate: config.LastResetValueTemplate,
			// LastReset:              config.LastReset,
			ValueTemplate:          config.ValueTemplate,
			Icon:                   config.Icon,
			PayloadOn:              "ON",
			PayloadOff:             "OFF",

			// Availability:           &Availability {
			// 	PayloadAvailable:    "",
			// 	PayloadNotAvailable: "",
			// 	Topic:               "",
			// 	ValueTemplate:       "",
			// },
			// AvailabilityMode:       "",
			// AvailabilityTemplate:   "",
			// AvailabilityTopic:      "",
			// EntityCategory:         "",
			// JsonAttributesTemplate: "",
			// JsonAttributesTopic:    "",
			// ObjectId:               "",
			// PayloadAvailable:       "",
			// PayloadNotAvailable:    "",
			// OffDelay:               0,
		}

		ct := JoinStringsForTopic(m.binarySensorPrefix, uid, "config")
		t := m.client.Publish(ct, 0, true, payload.Json())
		if !t.WaitTimeout(m.Timeout) {
			m.err = t.Error()
		}
	}

	return m.err
}

func (m *Mqtt) BinarySensorPublishValue(config EntityConfig) error {

	for range Only.Once {
		if !config.IsBinarySensor() {
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

func (m *Mqtt) BinarySensorPublishValues(configs []EntityConfig) error {
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
				topic = JoinStringsForTopic(m.binarySensorPrefix, config.StateTopic, "state")
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
			// fmt.Printf("%s (%s) -> %s\n", topic, n, string(j))

			// log.Fatalln("NEED TO FIX THIS")
			// m.err = m.PublishValue(n, topic, string(j))
			t := m.client.Publish(topic, 0, true, string(j))
			if !t.WaitTimeout(m.Timeout) {
				m.err = t.Error()
			}
		}
	}
	return m.err
}


type BinarySensor struct {
	Availability           *Availability `json:"availability,omitempty" required:"false"`
	AvailabilityMode       string       `json:"availability_mode,omitempty" required:"false"`
	AvailabilityTemplate   string       `json:"availability_template,omitempty" required:"false"`
	AvailabilityTopic      string       `json:"availability_topic,omitempty" required:"false"`
	Device                 Device       `json:"device,omitempty" required:"false"`
	DeviceClass            string       `json:"device_class,omitempty" required:"false"`
	EnabledByDefault       bool         `json:"enabled_by_default,omitempty" required:"false"`
	Encoding               string       `json:"encoding,omitempty" required:"false"`
	EntityCategory         string       `json:"entity_category,omitempty" required:"false"`
	ExpireAfter            int          `json:"expire_after,omitempty" required:"false"`
	ForceUpdate            bool         `json:"force_update,omitempty" required:"false"`
	Icon                   string       `json:"icon,omitempty" required:"false"`
	JsonAttributesTemplate string       `json:"json_attributes_template,omitempty" required:"false"`
	JsonAttributesTopic    string       `json:"json_attributes_topic,omitempty" required:"false"`
	LastResetValueTemplate string       `json:"last_reset_value_template,omitempty" required:"false"`
	Name                   string       `json:"name,omitempty" required:"false"`
	ObjectId               string       `json:"object_id,omitempty" required:"false"`
	PayloadAvailable       string       `json:"payload_available,omitempty" required:"false"`
	PayloadNotAvailable    string       `json:"payload_not_available,omitempty" required:"false"`
	Qos                    int          `json:"qos,omitempty" required:"false"`
	StateClass             string       `json:"state_class,omitempty" required:"false"`
	StateTopic             string       `json:"state_topic" required:"true"`
	UniqueId               string       `json:"unique_id,omitempty" required:"false"`
	UnitOfMeasurement      string       `json:"unit_of_measurement,omitempty" required:"false"`
	ValueTemplate          string       `json:"value_template,omitempty" required:"false"`

	OffDelay               int          `json:"off_delay,omitempty" required:"false"`
	PayloadOff             string       `json:"pl_off,omitempty" required:"false"`
	PayloadOn              string       `json:"pl_on,omitempty" required:"false"`
}
func (c *BinarySensor) Json() string {
	j, _ := json.Marshal(*c)
	return string(j)
}
