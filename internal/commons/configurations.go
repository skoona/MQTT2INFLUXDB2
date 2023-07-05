package commons

import (
	"context"
	"fyne.io/fyne/v2"
	"strings"

	"os"
)

// configuration keys
const (
	// basic runtime
	ApplicationName  = "applicationName"
	ApplicationId    = "applicationId"
	ApplicationTitle = "applicationTitle"
	CompanyName      = "companyName"
	InfluxHostUri    = "influxHostUri"
	InfluxBucket     = "influxBucket"
	InfluxOrg        = "influxOrg"
	InfluxToken      = "influxToken"
	MqttHostUri      = "mqttHostUri"
	MqttUser         = "mqttUser"
	MqttPass         = "mqttPass"
	TestMode         = "testMode"
	DebugMode        = "debugMode"

	// data modeling
	GarageType             = "garage"
	SensorType             = "sensor"
	GarageProperty         = "Details"
	ActualProperty         = "Actual"
	AmbientProperty        = "Ambient"
	PositionProperty       = "Position"
	SignalStrengthProperty = "SignalStrength"
	StateProperty          = "State"

	// context keys
	TestModeKey      = 7
	DebugModeKey     = 8
	SknAppIDKey      = 9
	InfluxHostUriKey = 10
	InfluxBucketKey  = 11
	InfluxOrgKey     = 12
	InfluxTokenKey   = 13
	MqttHostUriKey   = 14
	MqttUserKey      = 15
	MqttPassKey      = 16
)

var (
	fyne_instance fyne.App
	appSettings   map[string]string
)

/*
	  ConfigFromCtx extract the app config map[string]string from the context

		func ConfigFromCtx(ctx context.Context) *map[string]string {
			v := ctx.Value(CfgKey)
			t := reflect.ValueOf(v)
			i := t.Interface()
			return i.(*map[string]string)
		}
*/
func IsDebugMode() bool {
	return (appSettings[DebugMode] == "true")
}
func IsTestMode() bool {
	return (appSettings[TestMode] == "true")
}
func GetCompanyName() string {
	return appSettings[CompanyName]
}
func GetApplicationName() string {
	return appSettings[ApplicationName]
}
func GetApplicationTitle() string {
	return appSettings[ApplicationTitle]
}

func GetInfluxHostUri() string {
	return appSettings[InfluxHostUri]
}
func GetInfluxBucket() string {
	return appSettings[InfluxBucket]
}
func GetInfluxOrg() string {
	return appSettings[InfluxOrg]
}
func GetInfluxToken() string {
	return appSettings[InfluxToken]
}

func GetMqttHostUri() string {
	return appSettings[MqttHostUri]
}
func GetMqttUser() string {
	return appSettings[MqttUser]
}
func GetMqttPass() string {
	return appSettings[MqttPass]
}

// AppSettings collects env params and same preferences with env value priority
func AppSettings(ctx context.Context, a fyne.App) map[string]string {
	fyne_instance = a
	appID := ctx.Value(SknAppIDKey).(string)
	if appID == "" {
		appID = "net.skoona.mq2influx"
	}

	cfg := map[string]string{
		ApplicationName:  "mqttToInfluxDB",
		ApplicationId:    appID,
		ApplicationTitle: "Homie v3/MQTT to InfluxDB2",
		CompanyName:      "Skoona Development",
		InfluxHostUri:    "http://10.100.1.17:8086",
		InfluxBucket:     "SknSensors",
		InfluxOrg:        "skoona.net",
		InfluxToken:      "1ac1a18911f80510ee8c1de0d5a0d132dbec9e31cc1b9dc422f55d8c612d5498",
		MqttHostUri:      "tcp://10.100.1.16:1883",
		MqttUser:         "openhabian",
		MqttPass:         "Apache.Tomcat.8",
		TestMode:         "true",
		DebugMode:        "true",
	}
	appSettings = cfg

	//for key, value := range cfg {
	//	cfg[key] = fyne_instance.Preferences().StringWithFallback(key, value)
	//}

	value := os.Getenv("ENABLE_TEST_MODE")
	if value != "" {
		cfg[TestMode] = value
	}
	value = os.Getenv("ENABLE_DEBUG_MODE")
	if value != "" {
		cfg[DebugMode] = value
	}

	value = os.Getenv("INFLUXDB_URI")
	if value != "" {
		if strings.Contains(value, ":") {
			cfg[InfluxHostUri] = value
		} else {
			cfg[InfluxHostUri] = value + ":8086"
		}
	}
	value = os.Getenv("INFLUXDB_BUCKET")
	if value != "" {
		cfg[InfluxBucket] = value
	}
	value = os.Getenv("INFLUXDB_TOKEN")
	if value != "" {
		cfg[InfluxToken] = value
	}
	value = os.Getenv("INFLUXDB_ORG")
	if value != "" {
		cfg[InfluxOrg] = value
	}

	value = os.Getenv("MQTT_URI")
	if value != "" {
		if strings.Contains(value, ":") {
			cfg[MqttHostUri] = value
		} else {
			cfg[MqttHostUri] = value + ":1883"
		}
	}
	value = os.Getenv("MQTT_USER")
	if value != "" {
		cfg[MqttUser] = value
	}
	value = os.Getenv("MQTT_PASS")
	if value != "" {
		cfg[MqttPass] = value
	}

	for key, value := range cfg {
		fyne_instance.Preferences().SetString(key, value)
	}

	return cfg
}
