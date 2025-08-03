package protocol

type Messager interface {
	Type() DataType
	Data() any
}

type Message struct {
	dataType DataType
	data     any
}

func (c *Message) Type() DataType {
	return c.dataType
}

func (c *Message) Data() any {
	return c.data
}
