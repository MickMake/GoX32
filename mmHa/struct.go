package mmHa

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/MickMake/GoX32/Behringer/api"
	"github.com/MickMake/GoX32/Only"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"net/url"
	"time"
)


type Mqtt struct {
	ClientId      string        `json:"client_id"`
	Username      string        `json:"username"`
	Password      string        `json:"password"`
	Host          string        `json:"host"`
	Port          string        `json:"port"`
	Timeout       time.Duration `json:"timeout"`
	EntityPrefix  string        `json:"entity_prefix"`

	url           *url.URL
	client        mqtt.Client
	pubClient     mqtt.Client
	clientOptions *mqtt.ClientOptions
	onConnectHandler mqtt.OnConnectHandler
	connectionLostHandler mqtt.ConnectionLostHandler
	messageHandler mqtt.MessageHandler

	LastRefresh   time.Time `json:"-"`
	PsId          int64 `json:"-"`

	Device Device

	servicePrefix string
	sensorPrefix string
	binaryPrefix string
	lightPrefix string
	switchPrefix string
	selectPrefix string
	buttonPrefix string

	token    mqtt.Token
	firstRun bool
	err      error
}

type Fields map[string]string


func New(req Mqtt) *Mqtt {
	var ret Mqtt

	for range Only.Once {
		ret.err = ret.setUrl(req)
		if ret.err != nil {
			break
		}
		ret.firstRun = true
		ret.EntityPrefix = req.EntityPrefix

		ret.servicePrefix = "homeassistant/sensor/" + req.ClientId
		ret.binaryPrefix = "homeassistant/binary_sensor/" + req.ClientId
		ret.sensorPrefix = "homeassistant/sensor/" + req.ClientId
		ret.lightPrefix = "homeassistant/light/" + req.ClientId
		ret.switchPrefix = "homeassistant/switch/" + req.ClientId
		ret.selectPrefix = "homeassistant/select/" + req.ClientId
		ret.buttonPrefix = "homeassistant/button/" + req.ClientId
	}

	return &ret
}

func (m *Mqtt) IsFirstRun() bool {
	return m.firstRun
}

func (m *Mqtt) IsNotFirstRun() bool {
	return !m.firstRun
}

func (m *Mqtt) UnsetFirstRun() {
	m.firstRun = false
}

func (m *Mqtt) GetError() error {
	return m.err
}

func (m *Mqtt) IsError() bool {
	if m.err != nil {
		return true
	}
	return false
}

func (m *Mqtt) IsNewDay() bool {
	var yes bool
	for range Only.Once {
		last := m.LastRefresh.Format("20060102")
		now := time.Now().Format("20060102")

		if last != now {
			yes = true
			break
		}
	}
	return yes
}

func (m *Mqtt) setUrl(req Mqtt) error {

	for range Only.Once {
		// if req.Username == "" {
		// 	m.err = errors.New("username empty")
		// 	break
		// }
		m.Username = req.Username

		// if req.Password == "" {
		// 	m.err = errors.New("password empty")
		// 	break
		// }
		m.Password = req.Password

		if req.Host == "" {
			m.err = errors.New("HASSIO mqtt host not defined")
			break
		}
		m.Host = req.Host

		if req.Port == "" {
			req.Port = "1883"
		}
		m.Port = req.Port

		u := fmt.Sprintf("tcp://%s:%s@%s:%s",
			m.Username,
			m.Password,
			m.Host,
			m.Port,
			)
		m.url, m.err = url.Parse(u)
	}

	return m.err
}

func (m *Mqtt) Publish(config EntityConfig, UpdateConfig bool) error {
	for range Only.Once {
		if UpdateConfig {
			m.err = m.PublishConfig(config)
			if m.err != nil {
				break
			}
		}

		m.err = m.PublishValue(config)
		if m.err != nil {
			break
		}
	}
	return m.err
}

func (m *Mqtt) PublishConfig(config EntityConfig) error {
	for range Only.Once {
		switch {
			case config.IsSensor():
				m.err = m.PublishSensorConfig(config)

			case config.IsBinarySensor():
				m.err = m.PublishBinaryConfig(config)

			case config.IsSwitch():
				m.err = m.PublishSwitchConfig(config)

			case config.IsSelect():
				m.err = m.PublishSelectConfig(config)

			case config.IsLight():
				m.err = m.PublishLightConfig(config)

			case config.IsButton():
				m.err = m.PublishButtonConfig(config)

			default:
				m.err = m.PublishSensorConfig(config)
		}
	}
	return m.err
}

func (m *Mqtt) PublishValue(config EntityConfig) error {
	for range Only.Once {
		switch {
			case config.IsSensor():
				m.err = m.PublishSensorValue(config)

			case config.IsBinarySensor():
				m.err = m.PublishBinaryValue(config)

			case config.IsSwitch():
				m.err = m.PublishSwitchValue(config)

			case config.IsSelect():
				m.err = m.PublishSelectValue(config)

			case config.IsLight():
				m.err = m.PublishLightValue(config)

			case config.IsButton():
				m.err = m.PublishButtonValue(config)

			default:
				m.err = m.PublishSensorValue(config)
		}
	}
	return m.err
}

func (m *Mqtt) PublishValues(config []EntityConfig) error {
	for range Only.Once {
		m.err = m.PublishSensorValues(config)
		if m.err != nil {
			break
		}

		m.err = m.PublishBinaryValues(config)
		if m.err != nil {
			break
		}

		m.err = m.PublishSwitchValues(config)
		if m.err != nil {
			break
		}

		m.err = m.PublishSensorValues(config)
		if m.err != nil {
			break
		}

		m.err = m.PublishLightValues(config)
		if m.err != nil {
			break
		}

		m.err = m.PublishButtonValues(config)
		if m.err != nil {
			break
		}

		// m.err = m.PublishSelectValues(config)
		// if m.err != nil {
		// 	break
		// }
	}
	return m.err
}

func (m *Mqtt) PublishConfigs(configs []EntityConfig) error {
	for range Only.Once {
		for _, config := range configs {
			m.err = m.PublishConfig(config)
			if m.err != nil {
				break
			}
		}
	}
	return m.err
}

func (m *Mqtt) SetDeviceConfig(swname string, id string, name string, model string, vendor string, area string) error {
	for range Only.Once {
		id = JoinStringsForId(m.EntityPrefix, id)

		m.Device = Device {
			Connections:  [][]string{
				{swname, id},
			},
			Identifiers:  []string{id},
			Manufacturer: vendor,
			Model:        model,
			Name:         name,
			SwVersion:    swname + " https://github.com/MickMake/" + swname,
			ViaDevice:    swname,
			SuggestedArea: area,
		}
	}
	return m.err
}

func (m *Mqtt) GetSensorStateTopic(config EntityConfig) string {
	st := JoinStringsForId(m.Device.Name, config.ParentId, config.Name)
	st = JoinStringsForTopic(m.sensorPrefix, st, "state")		// m.GetSensorStateTopic(name, config.SubName),m.EntityPrefix, m.Device.FullName, config.SubName
	return st
}


type MqttState struct {
	LastReset string `json:"last_reset,omitempty"`
	Value string `json:"value"`
}
func (mq *MqttState) Json() string {
	var ret string
	for range Only.Once {
		j, err := json.Marshal(*mq)
		if err != nil {
			ret = fmt.Sprintf("{ \"error\": \"%s\"", err)
			break
		}
		ret = string(j)
	}
	return ret
}


type ValueMap map[string]string
func (vm *ValueMap) Json() string {
	var ret string
	for range Only.Once {
		j, err := json.Marshal(*vm)
		if err != nil {
			ret = fmt.Sprintf("{ \"error\": \"%s\"", err)
			break
		}
		ret = string(j)
	}
	return ret
}


type Availability struct {
	PayloadAvailable    string `json:"payload_available,omitempty" required:"false"`
	PayloadNotAvailable string `json:"payload_not_available,omitempty" required:"false"`
	Topic               string `json:"topic,omitempty" required:"true"`
	ValueTemplate       string `json:"value_template,omitempty" required:"false"`
}
type SensorState string


func (m *Mqtt) GetLastReset(pointType string) string {
	var ret string

	for range Only.Once {
		pt := api.GetDevicePoint(pointType)
		if !pt.Valid {
			break
		}
		if pt.Type == "" {
			break
		}
		ret = pt.WhenReset()
	}

	return ret
}


type EntityConfig struct {
	// Type       string
	Name        string
	SubName     string

	ParentId    string
	ParentName  string

	UniqueId    string
	// FullId      string
	Units       string
	ValueName   string
	DeviceClass string
	StateClass  string
	StateTopic  string
	Icon        string

	Value         string
	ValueTemplate string
	Options       []string

	LastReset              string
	LastResetValueTemplate string

	HaType string
}

var SensorLabels = Labels{"sensor", "int", "int32", "int64", "float", "float32", "float64", "string"}
func (config *EntityConfig) IsSensor() bool {
	var ok bool

	for range Only.Once {
		if config.IsBinarySensor() {
			break
		}
		if config.IsSwitch() {
			break
		}
		if config.IsLight() {
			break
		}
		if config.IsButton() {
			break
		}
		if config.IsSelect() {
			break
		}

		if SensorLabels.ValueExists(config.Units) {
			ok = true
			break
		}

		ok = true
	}

	return ok
}

var BinaryLabels = Labels{"binary", "state"}
func (config *EntityConfig) IsBinarySensor() bool {
	var ok bool

	for range Only.Once {
		if BinaryLabels.ValueExists(config.HaType) {
			ok = true
			break
		}
	}

	return ok
}

var SwitchLabels = Labels{"switch", "toggle"}
func (config *EntityConfig) IsSwitch() bool {
	var ok bool

	for range Only.Once {
		if SwitchLabels.ValueExists(config.HaType) {
			ok = true
			break
		}
	}

	return ok
}

var SelectLabels = Labels{"select", "array", "map"}
func (config *EntityConfig) IsSelect() bool {
	var ok bool

	for range Only.Once {
		if SelectLabels.ValueExists(config.HaType) {
			ok = true
			break
		}
	}

	return ok
}

var LightLabels = Labels{"light"}
func (config *EntityConfig) IsLight() bool {
	var ok bool

	for range Only.Once {
		if LightLabels.ValueExists(config.HaType) {
			ok = true
			break
		}
	}

	return ok
}

var ButtonLabels = Labels{"button"}
func (config *EntityConfig) IsButton() bool {
	var ok bool

	for range Only.Once {
		if ButtonLabels.ValueExists(config.HaType) {
			ok = true
			break
		}
	}

	return ok
}

func (config *EntityConfig) FixConfig() {

	for range Only.Once {
		// mdi:power-socket-au
		// mdi:solar-power
		// mdi:home-lightning-bolt-outline
		// mdi:transmission-tower
		// mdi:transmission-tower-export
		// mdi:transmission-tower-import
		// mdi:transmission-tower-off
		// mdi:home-battery-outline
		// mdi:lightning-bolt
		// mdi:check-circle-outline | mdi:arrow-right-bold

		switch config.Units {
			case "MW":
				fallthrough
			case "kW":
				fallthrough
			case "W":
				config.DeviceClass = SetDefault(config.DeviceClass, "power")
				config.Icon = SetDefault(config.Icon, "mdi:lightning-bolt")
				// config.ValueTemplate = SetDefault(config.ValueTemplate, "{{ value_json.value | float }}")
				config.ValueTemplate = SetDefault(config.ValueTemplate, fmt.Sprintf("{{ value_json.%s | float }}", config.ValueName))
				// - Used with merged values.

			case "MWh":
				fallthrough
			case "kWh":
				fallthrough
			case "Wh":
				config.DeviceClass = SetDefault(config.DeviceClass, "energy")
				config.Icon = SetDefault(config.Icon, "mdi:lightning-bolt")
				config.ValueTemplate = SetDefault(config.ValueTemplate, fmt.Sprintf("{{ value_json.%s | float }}", config.ValueName))

			case "kvar":
				config.DeviceClass = SetDefault(config.DeviceClass, "reactive_power")
				config.Icon = SetDefault(config.Icon, "mdi:lightning-bolt")
				config.ValueTemplate = SetDefault(config.ValueTemplate, fmt.Sprintf("{{ value_json.%s | float }}", config.ValueName))

			case "Hz":
				config.DeviceClass = SetDefault(config.DeviceClass, "frequency")
				config.Icon = SetDefault(config.Icon, "mdi:sine-wave")
				config.ValueTemplate = SetDefault(config.ValueTemplate, fmt.Sprintf("{{ value_json.%s | float }}", config.ValueName))

			case "V":
				config.DeviceClass = SetDefault(config.DeviceClass, "voltage")
				config.Icon = SetDefault(config.Icon, "mdi:current-dc")
				config.ValueTemplate = SetDefault(config.ValueTemplate, fmt.Sprintf("{{ value_json.%s | float }}", config.ValueName))

			case "A":
				config.DeviceClass = SetDefault(config.DeviceClass, "current")
				config.Icon = SetDefault(config.Icon, "mdi:current-ac")
				config.ValueTemplate = SetDefault(config.ValueTemplate, fmt.Sprintf("{{ value_json.%s | float }}", config.ValueName))

			case "??C":
				fallthrough
			case "C":
				fallthrough
			case "???":
				config.DeviceClass = SetDefault(config.DeviceClass, "temperature")
				config.Units = "??C"
				config.Icon = SetDefault(config.Icon, "mdi:thermometer")
				config.ValueTemplate = SetDefault(config.ValueTemplate, fmt.Sprintf("{{ value_json.%s | float }}", config.ValueName))

			case "battery":
				config.DeviceClass = SetDefault(config.DeviceClass, "battery")
				config.Icon = SetDefault(config.Icon, "mdi:home-battery-outline")
				config.ValueTemplate = SetDefault(config.ValueTemplate, fmt.Sprintf("{{ value_json.%s | float }}", config.ValueName))

			case "%":
				config.DeviceClass = SetDefault(config.DeviceClass, "percent")
				config.Icon = SetDefault(config.Icon, "")
				config.ValueTemplate = SetDefault(config.ValueTemplate, fmt.Sprintf("{{ value_json.%s | float }}", config.ValueName))

			case "dB":
				config.DeviceClass = SetDefault(config.DeviceClass, "signal_strength")
				config.Icon = SetDefault(config.Icon, "mdi:volume-high")
				config.ValueTemplate = SetDefault(config.ValueTemplate, fmt.Sprintf("{{ value_json.%s | float }}", config.ValueName))

			case "mS":
				config.DeviceClass = SetDefault(config.DeviceClass, "")	// "time")
				config.Icon = SetDefault(config.Icon, "mdi:clock")
				config.ValueTemplate = SetDefault(config.ValueTemplate, fmt.Sprintf("{{ value_json.%s | float }}", config.ValueName))

			default:
				config.DeviceClass = SetDefault(config.DeviceClass, "")
				config.Icon = SetDefault(config.Icon, "")
				config.ValueTemplate = SetDefault(config.ValueTemplate, fmt.Sprintf("{{ value_json.%s }}", config.ValueName))
		}

		switch {
			case BinaryLabels.ValueExists(config.HaType):
				config.DeviceClass = SetDefault(config.DeviceClass, "")
				config.Icon = SetDefault(config.Icon, "mdi:check-circle-outline")
				config.ValueTemplate = SetDefault(config.ValueTemplate, fmt.Sprintf("{{ value_json.%s }}", config.ValueName))

			case SwitchLabels.ValueExists(config.HaType):
				config.DeviceClass = SetDefault(config.DeviceClass, "")
				config.Icon = SetDefault(config.Icon, "mdi:toggle-switch-off-outline")
				config.ValueTemplate = SetDefault(config.ValueTemplate, fmt.Sprintf("{{ value_json.%s }}", config.ValueName))

			case SelectLabels.ValueExists(config.HaType):
				config.DeviceClass = SetDefault(config.DeviceClass, "")
				config.Icon = SetDefault(config.Icon, "mdi:format-list-group")
				config.ValueTemplate = SetDefault(config.ValueTemplate, fmt.Sprintf("{{ value_json.%s }}", config.ValueName))

			case LightLabels.ValueExists(config.HaType):
				config.DeviceClass = SetDefault(config.DeviceClass, "")
				config.Icon = SetDefault(config.Icon, "mdi:lightbulb")
				config.ValueTemplate = SetDefault(config.ValueTemplate, fmt.Sprintf("{{ value_json.%s }}", config.ValueName))

			case ButtonLabels.ValueExists(config.HaType):
				config.DeviceClass = SetDefault(config.DeviceClass, "")
				config.Icon = SetDefault(config.Icon, "mdi:gesture-tap-button")
				config.ValueTemplate = SetDefault(config.ValueTemplate, fmt.Sprintf("{{ value_json.%s }}", config.ValueName))
		}

		// if config.ValueTemplate == "" {
		// 	config.ValueTemplate = fmt.Sprintf("{{ value_json.%s }}", config.UniqueId)
		// }

		if config.StateClass == "instant" {
			config.StateClass = "measurement"
			break
		}

		if config.StateClass == "" {
			config.StateClass = "measurement"
			break
		}

		if config.LastReset != "" {
			break
		}

		pt := api.GetDevicePoint(config.UniqueId)
		if !pt.Valid {
			break
		}

		config.LastReset = pt.WhenReset()
		config.LastResetValueTemplate = SetDefault(config.LastResetValueTemplate, "{{ value_json.last_reset | as_datetime() }}")
		// config.LastResetValueTemplate = SetDefault(config.LastResetValueTemplate, "{{ value_json.last_reset | int | timestamp_local | as_datetime }}")

		if config.LastReset == "" {
			config.StateClass = "measurement"
			break
		}
		config.StateClass = "total"
	}
}

func SetDefault(value string, def string) string {
	if value == "" {
		value = def
	}
	return value
}
