# mqttToInfluxDB
<a href="https://homieiot.github.io/">
  <img src="https://homieiot.github.io/img/works-with-homie.png" alt="works with MQTT Homie">
</a>

Golang program to read certain MQTT messages for devices (sensors) that use the Homie v3 discovery protocol and forward those messages to InfluxDB2; with CLI and optional GUI Interface using Fyne.io.

## Two version of this program are available.
    CLI: Command Line Interface just does the work of forwarding messages.
    GUI: Does the work of forwards and displays the live status of the Homie devices discoveryed.

### Environment Variables used for runtime configuration
    $ export INFLUXDB_ORG="skoona.net"
    $ export INFLUXDB_BUCKET="SknSensors"
    $ export INFLUXDB_TOKEN="reallylonghexstringcalledfrominfluxdbtoken"
    $ export INFLUXDB_URI="http://ipOfInfluxDB:8086"
    $ export MQTT_URI="tcp://ipOfMqtt:1883"
    $ export MQTT_USER="username"
    $ export MQTT_PASS="userPassword"
    $ 
    $ export ENABLE_TEST_MODE="true"
    $ export ENABLE_DEBUG_MODE="true"
    $ export ENABLE_INFLUXDB_MODE="true"
Can also be configured thru the `settings` menu option.

### MQTT Subscriptions
    `.internal/adapters/repository/homieprovider.go` contains this list of scriptions:
* 	"+/+/+/humidity"
*  	"+/+/+/temperature"
*  	"+/+/+/motion"    
*   "+/+/+/occupancy" 
*  	"+/+/+/Position"  
*  	"+/+/+/State"   
*  	"+/+/+/Details"   
*  	"+/+/+/message"    
*   "+/+/+/name"      
*  	"+/+/+/heartbeat" 

The focus is to collect environmental reading from sensors on the Homie network.


## Development Notes
    Use build tags to generate gui or cli versions; --tags [gui|cli]
	Environment variables override the defaults and/or saved config in all cases.
	InfluxDB2 channel can be turned off via settings or environment var
    Two icons on status bar relate to (folder)message provider and (storage) message consumer being enabled
    The Main menu and or SysTray menu have the Quit option.  Closing the main window simply hides it.
    Using the CLI version with InfluxDB turned off, doesn't make much sense but can be done.
    There are five main components driving the whole app.
    * ViewProvider    Fyne GUI page manager
    * StreamService   Provide structure and control data flow usecase
    * StreamStorage   transforms MQTT messages into Entities
    * StreamConsumer  sends MQTT Messages to InfluxDB2
    * StreamProvider  reads MQTT messages

### File Tree
Modeled after [Hexagonal Architecture Pattern](https://en.wikipedia.org/wiki/Hexagonal_architecture_(software))
```text
├── LICENSE
├── Makefile
├── README.md
├── cmd
│   └── cli
│       └── main.go
├── go.mod
├── go.sum
├── internal
│   ├── adapters
│   │   ├── handler
│   │   │   └── ui
│   │   │       ├── menus.go
│   │   │       └── viewhandler.go
│   │   └── repository
│   │       ├── homieprovider.go
│   │       ├── influxconsumer.go
│   │       └── storagerepository.go
│   ├── commons
│   │   ├── configurations.go
│   │   ├── imageManager.go
│   │   ├── resources
│   │   │   ├── garage-closed.svg
│   │   │   ├── garage-open.svg
│   │   │   └── sensorsOn-mbo-24px.svg
│   │   └── svgImages.go
│   └── core
│       ├── domain
│       │   ├── baseMessage.go
│       │   └── devices.go
│       ├── ports
│       │   ├── storagerepository.go
│       │   ├── streamMessage.go
│       │   ├── streamService.go
│       │   ├── streamconsumer.go
│       │   └── streamprovider.go
│       └── services
│           └── streamService.go
├── main.go
└── skoona.png
```


## MIT License
The application is available as open source under the terms of the [MIT License](http://opensource.org/licenses/MIT).	