# GoX32 - Behringer X32 for HASSIO

This application provides a link between the [Behringer X32](https://www.behringer.com/behringer/product?modelCode=P0AWN) Sound-mixer and HASSIO.
All X32 messages are transmitted to [HomeAssistant](https://www.home-assistant.io/) via MQTT auto-discovery.
A JSON config file allows changes to be made to X32 endpoints, (icons, types, maps, etc)


## What is it?

This GoLang package allows you to sync your [Behringer X32](https://www.behringer.com/behringer/product?modelCode=P0AWN) mixer with your [HomeAssistant](https://www.home-assistant.io/) instance.
It has a complete implementation of all X32 endpoints. However, there are a few dodgy things.

![alt text](https://github.com/MickMake/GoX32/blob/master/docs/X32RACK.png?raw=true)

![alt text](https://github.com/MickMake/GoX32/blob/master/docs/X32.png?raw=true)

I'm currently using it in my [HomeAssistant](https://www.home-assistant.io/) instance.


## What state is it in?

This is currently usable for my needs, (seeing all data in [HomeAssistant](https://www.home-assistant.io/)), but there's quite a few things that need to be completed.

![alt text](https://github.com/MickMake/GoX32/blob/master/docs/ScreenShot.png?raw=true)

I've implemented most of the features I've wanted to, but the current roadmap is:

1. Make this a HASSIO add-on.
2. Provide mixer icons and images for HASSIO.
3. Add snazzy HASSIO cards to enable quick adding of mixer components. 

Note: Until I create a HASSIO add-on, you will have to run this manually.


## Using GoX32:

### Initial config.

The following config options are required. Make sure you set these correctly.
```
$ ./bin/GoSungrow config write \
--mqtt-host="<YOUR HASSIO IP ADDRESS>" \
--mqtt-password="<YOUR HASSIO PASSWORD>" \
--mqtt-port="1883" \
--mqtt-user="<YOUR HASSIO USER>" \
--x32-host="<YOUR X32 IP ADDRESS>"

Using config file '/Users/mick/.GoSungrow/config.json'
```


## Flags & environment variables.
```
Using environment variables instad of flags.
+-----------------+------------+-------------------+--------------------------------+--------------------------------+
|      FLAG       | SHORT FLAG |    ENVIRONMENT    |          DESCRIPTION           |            DEFAULT             |
+-----------------+------------+-------------------+--------------------------------+--------------------------------+
| --x32-host      |            | X32_X32_HOST      | Behringer X32: Host / IP       |                                |
|                 |            |                   | address.                       |                                |
| --x32-port      |            | X32_X32_PORT      | Behringer X32: Port.           |                          10023 |
| --x32-user      | -u         | X32_X32_USER      | Behringer X32: Username.       |                                |
| --x32-password  | -p         | X32_X32_PASSWORD  | Behringer X32: Password.       |                                |
| --x32-timeout   |            | X32_X32_TIMEOUT   | Behringer X32: Timeout.        | 30s                            |
| --mqtt-user     |            | X32_MQTT_USER     | HASSIO: mqtt username.         |                                |
| --mqtt-password |            | X32_MQTT_PASSWORD | HASSIO: mqtt password.         |                                |
| --mqtt-host     |            | X32_MQTT_HOST     | HASSIO: mqtt host.             |                                |
| --mqtt-port     |            | X32_MQTT_PORT     | HASSIO: mqtt port.             |                                |
| --config        |            | X32_CONFIG        | GoX32: config file.            | /Users/mick/.GoX32/config.json |
| --debug         |            | X32_DEBUG         | GoX32: Debug mode.             | false                          |
| --quiet         | -q         | X32_QUIET         | GoX32: Silence all messages.   | false                          |
+-----------------+------------+-------------------+--------------------------------+--------------------------------+
```
