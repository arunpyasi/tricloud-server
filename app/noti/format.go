package noti

import "encoding/json"

type MessageFormat struct {
	Type string
	Data interface{}
}

func NewFormattedMessage(msgtype string, data interface{}) string {
	msg := &MessageFormat{
		Type: msgtype,
		Data: data,
	}

	bytes, err := json.Marshal(msg)
	if err != nil {
		return ""
	}
	return string(bytes)
}
