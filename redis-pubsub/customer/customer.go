package Customer

import "encoding/json"

const (
	ChannelCreation = "customer.created"
)

type Customer struct {
	Name string `json:"name"`
	Age  int64  `json:"age"`
}

func (c Customer) Json() (string, error) {
	jsonBytes, err := json.Marshal(c)

	return string(jsonBytes), err
}
