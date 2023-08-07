package domain

import (
	"encoding/json"
	"fmt"
	"github.com/skoona/mqttToInfluxDB/internal/commons"
	"strings"
	"time"
)

type baseMessage struct {
	topic     string
	network   string
	device    string
	node      string
	property  string
	value     string
	timestamp time.Time
	ActualI   int     `json:"range"`
	Average   int     `json:"average"`
	Mapped    int     `json:"mapped"`
	Status    string  `json:"status"`
	RawStatus int     `json:"raw_status"`
	Signal    float32 `json:"signal"`
	AmbientF  float32 `json:"ambient"`
	Movement  string  `json:"movement"`
}

func NewStreamMessage(topic, value string) (*baseMessage, error) {
	base := &baseMessage{}
	parts := strings.Split(topic, "/")

	if strings.Contains(parts[3], commons.GarageProperty) {
		err := json.Unmarshal([]byte(value), base)
		if err != nil {
			fmt.Println("JSON Parse Error: ", err.Error())
			return base, err
		}
	}
	base.topic = topic
	base.value = value
	base.network = parts[0]
	base.device = parts[1]
	base.node = parts[2]
	base.property = parts[3]
	base.timestamp = time.Now()

	return base, nil
}

func (s *baseMessage) IsGarageDoor() bool {
	return s.property == commons.GarageProperty
}
func (s *baseMessage) Topic() string {
	return s.topic
}
func (s *baseMessage) Network() string {
	return s.network
}
func (s *baseMessage) Device() string {
	return s.device
}
func (s *baseMessage) Node() string {
	return s.node
}
func (s *baseMessage) Property() string {
	return s.property
}
func (s *baseMessage) Value() string {
	return s.value
}
func (s *baseMessage) Timestamp() string {
	return s.timestamp.Format(time.RFC3339)
}
func (s *baseMessage) Ambient() float32 {
	return s.AmbientF
}
func (s *baseMessage) Actual() int {
	return s.ActualI
}
func (s *baseMessage) Position() int {
	return s.Mapped
}
func (s *baseMessage) SignalStrength() float32 {
	return s.Signal
}
func (s *baseMessage) State() string {
	return s.Movement
}
