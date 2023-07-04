package entities

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type baseMessage struct {
	topic     string    `json:"-,omitempty"`
	network   string    `json:"-,omitempty"`
	device    string    `json:"-,omitempty"`
	node      string    `json:"-,omitempty"`
	property  string    `json:"-,omitempty"`
	value     string    `json:"-,omitempty"`
	timestamp time.Time `json:"-,omitempty"`
	actual    int       `json:"range"`
	average   int       `json:"average"`
	mapped    int       `json:"mapped"`
	status    string    `json:"status"`
	rawStatus int       `json:"raw_status"`
	signal    float32   `json:"signal"`
	ambient   float32   `json:"ambient"`
	movement  string    `json:"movement"`
}

func NewStreamMessage(topic, value string) (*baseMessage, error) {
	base := &baseMessage{}
	parts := strings.Split(topic, "/")

	if strings.Contains(parts[3], GarageProperty) {
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
	return s.property == GarageProperty
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
	return s.ambient
}
func (s *baseMessage) Actual() int {
	return s.actual
}
func (s *baseMessage) Position() int {
	return s.mapped
}
func (s *baseMessage) SignalStrength() float32 {
	return s.signal
}
func (s *baseMessage) State() string {
	return s.movement
}
