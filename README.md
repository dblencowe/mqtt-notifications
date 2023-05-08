# MQTT Notifications
Receive messages from an MQTT broker right on your desktop

## Building
To build the application, you wil require at least Go 1.20 installed on your machine
```bash
$ go version
go version go1.20.3 linux/amd64
```

Build
```bash
make build
```

## Configuration
The application is configured through Environment Variables provided
at runtime

| Variable                 | Optional | Example               | Description                     |
|--------------------------|----------|-----------------------|---------------------------------|
| MQTT_BROKER_URL          |          | localhost             | The hostname of the MQTT server |
| MQTT_BROKER_PORT         |          | 1883                  | The port of the MQTT Server     |
| MQTT_NOTIFICATIONS_TOPIC |          | notifications/publish | The topic to subscribe to       |
| MQTT_BROKER_USERNAME     | optional | admin                 | Broker username                 |
| MQTT_BROKER_PASSWORD     | optional | password              | Broker password                 |

## Installation
### Basic user install
Run `make install` to install the binary to your $HOME directory and then add it to your
startup applications

### Running at login
I'm currently still figuring out the best methods for running on user login for each
major operating system. Until then I'm using the following Bash snippet for my Linux
Desktop:

```bash
BIN_PATH="$HOME/.mqtt-notifications/bin/mqtt-notifications"
if ! pidof -x "$BIN_PATH" 1> /dev/null
then
    MQTT_BROKER_URL=localhost \
    MQTT_BROKER_PORT=1883 \
    MQTT_NOTIFICATIONS_TOPIC=notifications/publish \
    $BIN_PATH > /var/log/mqtt-notifications.log 2>&1 &
fi
```