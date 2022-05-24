package mmHa

import (
	"errors"
	"fmt"
	"github.com/MickMake/GoX32/Only"
	"github.com/MickMake/GoX32/defaults"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"log"
	"time"
)


func (m *Mqtt) SetAuth(username string, password string) error {

	for range Only.Once {
		if username == "" {
			m.err = errors.New("username empty")
			break
		}
		m.Username = username

		if password == "" {
			m.err = errors.New("password empty")
			break
		}
		m.Password = password
	}

	return m.err
}

func (m *Mqtt) Connect() error {
	for range Only.Once {
		m.err = m.createClientOptions()
		if m.err != nil {
			break
		}

		m.client = mqtt.NewClient(m.clientOptions)
		for m.token = m.client.Connect(); m.token.Wait() && m.token.Error() != nil; m.token = m.client.Connect() {
			log.Println("MQTT: error connecting:", m.token.Error())
			time.Sleep(time.Second * 5)
		}
		if m.err = m.token.Error(); m.err != nil {
			break
		}
		if m.ClientId == "" {
			m.ClientId = defaults.BinaryName
		}

		device := Config {
			Entry:      m.servicePrefix,
			Name:       m.ClientId,
			UniqueId:   m.ClientId, 	// + "_Service",
			StateTopic:   "~/state",
			DeviceConfig: DeviceConfig {
				Identifiers:  []string{defaults.BinaryName},
				SwVersion:    fmt.Sprintf("%s https://%s", defaults.BinaryName, defaults.Repo),
				Name:         m.ClientId + " Service",
				Manufacturer: "MickMake",
				Model:        defaults.BinaryName,
			},
		}

		m.err = m.RawPublish(JoinStringsForTopic(m.servicePrefix, "config"), 0, true, device.Json())
		if m.err != nil {
			break
		}

		m.err = m.RawPublish(JoinStringsForTopic(m.servicePrefix, "state"), 0, true, "ON")
		if m.err != nil {
			break
		}

	}

	return m.err
}

func (m *Mqtt) Disconnect() error {
	for range Only.Once {
		m.client.Disconnect(250)
		time.Sleep(1 * time.Second)
	}
	return m.err
}

func (m *Mqtt) createClientOptions() error {
	for range Only.Once {
		m.clientOptions = mqtt.NewClientOptions()
		m.clientOptions.AddBroker(fmt.Sprintf("tcp://%s", m.url.Host))
		m.clientOptions.SetUsername(m.url.User.Username())
		password, _ := m.url.User.Password()
		m.clientOptions.SetPassword(password)
		m.clientOptions.SetClientID(m.ClientId)

		m.clientOptions.WillTopic = JoinStringsForTopic(m.servicePrefix, "state")
		m.clientOptions.WillPayload = []byte("OFF")
		m.clientOptions.WillQos = 0
		m.clientOptions.WillRetained = true
		m.clientOptions.WillEnabled = true

		// mqttOptions := mqtt.NewClientOptions()
		// m.clientOptions.SetClientID(viper.GetString("mqtt.client_id"))
		// m.clientOptions.SetUsername(viper.GetString("mqtt.username"))
		// m.clientOptions.SetPassword(viper.GetString("mqtt.password"))
		m.clientOptions.SetMaxReconnectInterval(time.Second * 5)
		m.clientOptions.SetConnectTimeout(time.Second)
		m.clientOptions.SetCleanSession(true)
		m.clientOptions.SetAutoReconnect(true)
		m.clientOptions.SetOrderMatters(true)
		m.clientOptions.SetKeepAlive(time.Second * 300)
		// m.clientOptions.AddBroker(viper.GetString("mqtt.broker"))

		m.err = m.SetOnConnectHandler(m.onConnectHandler)
		if m.err != nil {
			break
		}

		m.err = m.SetConnectionLostHandler(m.connectionLostHandler)
		if m.err != nil {
			break
		}
	}
	return m.err
}

func (m *Mqtt) SetOnConnectHandler(fn mqtt.OnConnectHandler) error {
	for range Only.Once {
		m.onConnectHandler = fn
		if m.onConnectHandler == nil {
			m.onConnectHandler = m.handlerOnConnect
		}
		m.clientOptions.SetOnConnectHandler(m.onConnectHandler)
	}
	return m.err
}
func (m *Mqtt) handlerOnConnect(_ mqtt.Client) {
	log.Println("MQTT: Connected!")
	if m.EntityPrefix == "" {
		m.EntityPrefix = "x32"
	}
	topic := fmt.Sprintf("%s/set/#", m.EntityPrefix)
	fmt.Printf("# connectHandler() - Topic: %s\n", topic)
	m.client.Subscribe(topic, 0, m.onMessage)
}

func (m *Mqtt) SetConnectionLostHandler(fn mqtt.ConnectionLostHandler) error {
	for range Only.Once {
		m.connectionLostHandler = fn
		if m.connectionLostHandler == nil {
			m.connectionLostHandler = m.handlerConnectionLost
		}
		m.clientOptions.SetConnectionLostHandler(m.connectionLostHandler)
	}
	return m.err
}
func (m *Mqtt) handlerConnectionLost(_ mqtt.Client, err error) {
	fmt.Println("MQTT: Connection lost:", err)
}


func (m *Mqtt) SetMessageHandler(fn mqtt.MessageHandler) error {
	for range Only.Once {
		m.messageHandler = fn
		if m.messageHandler == nil {
			m.messageHandler = handlerMessage
		}
	}
	return m.err
}
func handlerMessage(_ mqtt.Client, message mqtt.Message) {
	fmt.Printf("# handlerMessage()\n")
	fmt.Printf("\t- %s\n", message.Topic())
	fmt.Printf("\t- %d\n", message.MessageID())
	fmt.Printf("\t- %v\n", message.Qos())
	fmt.Printf("\t- %v\n", message.Duplicate())
	fmt.Printf("\t- %v\n", message.Retained())
	fmt.Printf("\t- %s\n", string(message.Payload()))
}
func (m *Mqtt) onMessage(client mqtt.Client, message mqtt.Message) {
	for range Only.Once {
		if m.messageHandler == nil {
			break
		}
		m.messageHandler(client, message)

		// prefixParts := strings.Split(m.EntityPrefix, "/")
		// fmt.Printf("# onMessage() - Topic: %v\n", prefixParts)
		//
		// parts := strings.Split(message.Topic(), "/")
		// parts = parts[len(prefixParts)+1:]
		// res := make(mqttPayload, 0, 2)
		// m.err = json.Unmarshal(message.Payload(), &res)
		// if m.err != nil {
		// 	log.Println("MQTT: Invalid message payload:", m.err)
		// 	break
		// }
		//
		// values := make([]any, 0, len(res))
		// for _, p := range res {
		// 	switch p.Type {
		// 		case reflect.TypeOf(float32(0)).String():
		// 			values = append(values, float32(p.Value.(float64)))
		// 		case reflect.TypeOf("").String():
		// 			values = append(values, p.Value.(string))
		// 	}
		// }
		//
		// address := "/" + strings.Join(parts, "/")
		// log.Println(address)
		// m.err = osc.EmitMessage(address, values...)
		// if m.err != nil {
		// 	log.Println("Could not send OSC message:", m.err)
		// 	break
		// }
	}
}


func (m *Mqtt) RawSubscribe(topic string) error {
	for range Only.Once {
		t := m.client.Subscribe(topic, 0, func(client mqtt.Client, msg mqtt.Message) {
			fmt.Printf("* [%s] %s\n", msg.Topic(), string(msg.Payload()))
		})
		if !t.WaitTimeout(m.Timeout) {
			m.err = t.Error()
			// m.err = errors.New("mqtt subscribe timeout")
		}
	}
	return m.err
}

func (m *Mqtt) RawPublish(topic string, qos byte, retained bool, payload interface{}) error {
	for range Only.Once {
		t := m.client.Publish(topic, qos, retained, payload)
		if !t.WaitTimeout(m.Timeout) {
			m.err = t.Error()
			// m.err = errors.New("mqtt publish timeout")
		}
	}
	return m.err
}
