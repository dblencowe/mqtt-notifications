.PHONY: .local-mqtt

local-mqtt:
	docker run -d -p 1883:1883 -p 9001:9001 --name mqtt eclipse-mosquitto mosquitto -c /mosquitto-no-auth.conf

run-local:
	MQTT_BROKER_URL=localhost \
	MQTT_BROKER_PORT=1883 \
	MQTT_NOTIFICATIONS_TOPIC=notifications/publish \
	go run -tags nodbus ./cmd/main.go

build:
	go build -o ./bin/mqtt-notifications ./cmd/main.go