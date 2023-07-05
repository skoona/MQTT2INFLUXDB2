# mqttToInfluxDB

Golang program to read cerain MQTT messages for devices (sensors) using the Homie v3 discovery protocal and forward those messages to InfluxDB2; with CLI and optional GUI Interface using Fyne.io.

## Development Notes
	fyne data binding process require an address of an alternate object to be used for successful dynamic updates, and when source changes the object must be set.
	Environment variables override the default and/or save config in all cases.
	InfluxDB2 channel can be turned off via settings or environment var

## MIT License
	