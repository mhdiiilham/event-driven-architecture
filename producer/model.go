package main

import "encoding/json"

type Decible struct {
	Value int
}

type Event struct {
	CorrelationID string
	ContentType   string
	Body          any
}

func (e *Event) ToBytes() ([]byte, error) {
	return json.Marshal(e)
}

func (e Event) String() string {
	b, _ := e.ToBytes()
	return string(b)
}
