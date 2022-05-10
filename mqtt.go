package main

//
// var (
// 	mq mqtt.Client
// )
//
// type dataField struct {
// 	Type  string `json:"type"`
// 	Value any    `json:"value"`
// }
// type mqttPayload []*dataField
//
// func setupMQTTClient() {
// 	mqttOptions := mqtt.NewClientOptions()
// 	mqttOptions.SetClientID(viper.GetString("mqtt.client_id"))
// 	mqttOptions.SetUsername(viper.GetString("mqtt.username"))
// 	mqttOptions.SetPassword(viper.GetString("mqtt.password"))
// 	mqttOptions.SetMaxReconnectInterval(time.Second * 5)
// 	mqttOptions.SetConnectTimeout(time.Second)
// 	mqttOptions.SetCleanSession(viper.GetBool("mqtt.clean_session"))
// 	mqttOptions.SetAutoReconnect(true)
// 	mqttOptions.SetOnConnectHandler(connectHandler)
// 	mqttOptions.SetConnectionLostHandler(connectionLostHandler)
// 	mqttOptions.SetOrderMatters(true)
// 	mqttOptions.SetKeepAlive(viper.GetDuration("mqtt.keep_alive"))
// 	mqttOptions.AddBroker(viper.GetString("mqtt.broker"))
//
// 	mq = mqtt.NewClient(mqttOptions)
// 	for token := mq.Connect(); token.Wait() && token.Error() != nil; token = mq.Connect() {
// 		log.Println("MQTT: error connecting:", token.Error())
// 		time.Sleep(time.Second * 5)
// 	}
// }
//
// func connectionLostHandler(_ mqtt.Client, err error) {
// 	log.Println("MQTT: Connection lost:", err)
// }
//
// func connectHandler(_ mqtt.Client) {
// 	log.Println("MQTT: Connected!")
// 	prefix := viper.GetString("mqtt.topic_prefix")
// 	topic := fmt.Sprintf("%s%s", prefix, "/set/#")
// 	fmt.Printf("# connectHandler() - Topic: %s\n", topic)
// 	mq.Subscribe(topic, 0, onMessage)
// }
//
// func onMessage(_ mqtt.Client, message mqtt.Message) {
// 	parts := strings.Split(message.Topic(), "/")
// 	prefixParts := strings.Split(viper.GetString("mqtt.topic_prefix"), "/")
// 	fmt.Printf("# onMessage() - Topic: %v\n", prefixParts)
//
// 	parts = parts[len(prefixParts)+1:]
// 	res := make(mqttPayload, 0, 2)
// 	err := json.Unmarshal(message.Payload(), &res)
// 	if err != nil {
// 		log.Println("MQTT: Invalid message payload:", err)
// 		return
// 	}
//
// 	values := make([]any, 0, len(res))
// 	for _, p := range res {
// 		switch p.Type {
// 			case reflect.TypeOf(float32(0)).String():
// 				values = append(values, float32(p.Value.(float64)))
// 			case reflect.TypeOf("").String():
// 				values = append(values, p.Value.(string))
// 		}
// 	}
//
// 	address := "/" + strings.Join(parts, "/")
// 	log.Println(address)
// 	err = cli.EmitMessage(address, values...)
// 	if err != nil {
// 		log.Println("Could not send OSC message:", err)
// 		return
// 	}
// }
