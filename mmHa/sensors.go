package mmHa

import (
	"encoding/json"
	"fmt"
	"github.com/MickMake/GoX32/Only"
)


func (m *Mqtt) SensorPublishConfig(config EntityConfig) error {

	for range Only.Once {
		config.FixConfig()
		if !config.IsSensor() {
			break
		}

		// if config.LastReset == "" {
		// 	config.LastReset = m.GetLastReset(config.UniqueId)
		// }
		// if config.LastReset != "" {
		// 	if config.LastResetValueTemplate == "" {
		// 		config.LastResetValueTemplate = "{{ value_json.last_reset | as_datetime() }}"
		// 	}
		// 	// LastResetValueTemplate = "{{ value_json.last_reset | int | timestamp_local | as_datetime }}"
		// }

		// switch config.Units {
		// 	case "MW":
		// 		fallthrough
		// 	case "kW":
		// 		fallthrough
		// 	case "W":
		// 		config.DeviceClass = "power"
		//
		// 	case "MWh":
		// 		fallthrough
		// 	case "kWh":
		// 		fallthrough
		// 	case "Wh":
		// 		config.DeviceClass = "energy"
		//
		// 	case "kvar":
		// 		config.DeviceClass = "reactive_power"
		//
		// 	case "Hz":
		// 		config.DeviceClass = "frequency"
		//
		// 	case "V":
		// 		config.DeviceClass = "voltage"
		//
		// 	case "A":
		// 		config.DeviceClass = "current"
		//
		// 	case "℃":
		// 		config.DeviceClass = "temperature"
		// 		// point.Unit = "C"
		//
		// 	case "C":
		// 		config.DeviceClass = "temperature"
		// 		config.Units = "℃"
		//
		// 	case "%":
		// 		config.DeviceClass = "battery"
		//
		// 	default:
		// 		ValueTemplate = "{{ value_json.value }}"
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
		} else {
			config.StateTopic = JoinStringsForId(m.Device.Name, config.ParentId, config.Name)
			// config.StateTopic = JoinStringsForId(m.Device.Name, config.ParentName, config.Name, config.UniqueId),
		}
		uid := JoinStringsForId(m.Device.Name, config.ParentId, config.Name)

		payload := Sensor {
			Device:                 device,
			Name:                   JoinStrings(m.Device.Name, config.ParentName, config.Name),
			StateTopic:             JoinStringsForTopic(m.sensorPrefix, config.StateTopic, "state"),
			// StateTopic:             m.GetSensorStateTopic(name, config.SubName),m.EntityPrefix, m.Device.FullName, config.SubName
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
			LastReset:              config.LastReset,
			ValueTemplate:          config.ValueTemplate,
			Icon:                   config.Icon,

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
		}

		ct := JoinStringsForTopic(m.sensorPrefix, uid, "config")
		t := m.client.Publish(ct, 0, true, payload.Json())
		if !t.WaitTimeout(m.Timeout) {
			m.err = t.Error()
		}
	}

	return m.err
}

func (m *Mqtt) SensorPublishValue(config EntityConfig) error {

	for range Only.Once {
		if !config.IsSensor() {
			break
		}

		st := JoinStringsForId(m.Device.Name, config.ParentId, config.Name)
		payload := MqttState {
			LastReset: m.GetLastReset(JoinStringsForId(config.ParentId, config.Name)),
			Value:     config.Value,
		}
		st = JoinStringsForTopic(m.sensorPrefix, st, "state")
		t := m.client.Publish(st, 0, true, payload.Json())
		if !t.WaitTimeout(m.Timeout) {
			m.err = t.Error()
		}
	}

	return m.err
}

func (m *Mqtt) SensorPublishValues(configs []EntityConfig) error {
	for range Only.Once {
		if len(configs) == 0 {
			break
		}

		cs := make(map[string]Fields)
		topic := ""
		for i, config := range configs {
			if config.StateTopic != "" {
				config.StateTopic = JoinStringsForId(m.Device.Name, config.ParentId, config.StateTopic)
			} else {
				config.StateTopic = JoinStringsForId(m.Device.Name, config.ParentId, config.Name)
				// config.StateTopic = JoinStringsForId(m.Device.Name, config.ParentName, config.Name, config.UniqueId),
			}

			if topic == "" {
				topic = JoinStringsForTopic(m.sensorPrefix, config.StateTopic, "state")
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


type Sensor struct {
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
	LastReset              string       `json:"last_reset,omitempty" required:"false"`

	// StateFunc func() string `json:"-"`
	//
	// UpdateInterval  float64 `json:"-"`
	// ForceUpdateMQTT bool    `json:"-"`
}
func (c *Sensor) Json() string {
	j, _ := json.Marshal(*c)
	return string(j)
}
