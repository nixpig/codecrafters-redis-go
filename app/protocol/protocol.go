package protocol

import (
	"bytes"
	"strconv"
)

type DataPrefix byte
type DataType int

type Message struct {
	dataType DataType
	data     any
}

type Messager interface {
	Type() DataType
	Data() any
}

const (
	terminator = "\r\n"
)

const (
	PrefixBulkString   DataPrefix = '$'
	PrefixSimpleString DataPrefix = '+'
	PrefixBulkError    DataPrefix = '!'
	PrefixError        DataPrefix = '-'
	PrefixInteger      DataPrefix = ':'
	PrefixArray        DataPrefix = '*'
	PrefixBool         DataPrefix = '#'
	PrefixDouble       DataPrefix = ','
	PrefixBigNumber    DataPrefix = '('
	PrefixMap          DataPrefix = '%'
	PrefixAttributes   DataPrefix = '|'
	PrefixSet          DataPrefix = '~'
	PrefixPush         DataPrefix = '>'
	// PrefixVerbatimString DataPrefix = '='
)

const (
	DataBulkString DataType = iota
	DataSimpleString
	DataArray
	DataError
	DataBulkError
	DataInteger
	DataNull
	DataBool
	DataDouble
	DataBigNumber
	DataMap
	DataAttributes
	DataSet
	DataPush

	// VerbatimString
	// NullBulkString
	// NullArray
	// Null
	// PositiveInf
	// NegativeInf
	// NaN
)

var dataMap = map[DataPrefix]DataType{
	PrefixBulkString:   DataBulkString,
	PrefixSimpleString: DataSimpleString,
	PrefixArray:        DataArray,
	PrefixBulkError:    DataBulkError,
	PrefixError:        DataError,
	PrefixInteger:      DataInteger,
	PrefixBool:         DataBool,
	PrefixDouble:       DataDouble,
	PrefixBigNumber:    DataBigNumber,
	PrefixMap:          DataMap,
	PrefixAttributes:   DataAttributes,
	PrefixSet:          DataSet,
	PrefixPush:         DataPush,
}

func (c *Message) Type() DataType {
	return c.dataType
}

func (c *Message) Data() any {
	return c.data
}

// Unmarshal parses the RESP-encoded byte data in b and returns a Command and error
func Unmarshal(b []byte) (*Message, error) {
	if len(b) == 0 {
		return nil, NewErrInvalidData("zero length data")
	}

	if !bytes.HasSuffix(b, []byte(terminator)) {
		return nil, NewErrInvalidData("invalid ending")
	}

	prefix := DataPrefix(b[0])

	commandType, ok := dataMap[prefix]
	if !ok {
		return nil, NewErrInvalidCommand("invalid command: " + string(b[0]))
	}

	var data any
	// offset := 1

	switch commandType {
	case DataBulkString:
		data = ""

	case DataSimpleString:
		data = string(b[1:])

	// TODO: get this working
	case DataArray:
		// parts := bytes.SplitN(b[offset:], []byte(terminator), 3)
		// if len(parts) < 3 {
		// 	return nil, NewErrInvalidData("not enough parts")
		// }
		//
		// n, _ := strconv.Atoi(string(parts[0]))
		// data := make([]Message, n)

	}

	return &Message{
		dataType: commandType,
		data:     data,
	}, nil
}

// Marshal encodes Command to RESP-encoded byte slice
func Marshal(cmd *Message) ([]byte, error) {
	var buf bytes.Buffer

	switch data := (cmd.data).(type) {
	case string:
		switch cmd.dataType {
		case DataBulkString:
			buf.WriteByte(byte(PrefixBulkString))
			buf.WriteString(strconv.Itoa(len(data)) + terminator + data + terminator)

		case DataSimpleString:
			buf.WriteByte(byte(PrefixSimpleString))
			buf.WriteString(data + terminator)

		case DataError:
			buf.WriteByte(byte(PrefixError))
			buf.WriteString(data + terminator)

		case DataBulkError:
			buf.WriteByte(byte(PrefixBulkError))
			buf.WriteString(strconv.Itoa(len(data)) + terminator + data + terminator)

		case DataBigNumber:
			buf.WriteByte(byte(PrefixBigNumber))
			buf.WriteString(data + terminator)
		}

	case []Message:
		switch cmd.dataType {
		case DataArray:
			buf.WriteByte(byte(PrefixArray))
			buf.WriteString(strconv.Itoa(len(data)))
			buf.WriteString(terminator)

			for _, item := range data {
				b, err := Marshal(&item)
				if err != nil {
					return nil, NewErrInvalidData("invalid data")
				}
				buf.Write(b)
			}

		case DataPush:
			buf.WriteByte(byte(PrefixPush))
			buf.WriteString(strconv.Itoa(len(data)) + terminator)

			for _, item := range data {
				d, err := Marshal(&item)
				if err != nil {
					return nil, NewErrInvalidData("invalid data")
				}

				buf.Write(d)
			}
		}

	case int:
		switch cmd.dataType {
		case DataInteger:
			buf.WriteByte(byte(PrefixInteger))
			buf.WriteString(strconv.Itoa(data) + terminator)
		}

	case bool:
		switch cmd.dataType {
		case DataBool:
			buf.WriteByte(byte(PrefixBool))
			var v string
			if data {
				v = "t"
			} else {
				v = "f"
			}

			buf.WriteString(v + terminator)
		}

	case float64:
		switch cmd.dataType {
		case DataDouble:
			buf.WriteByte(byte(PrefixDouble))
			buf.WriteString(strconv.FormatFloat(data, 'G', -1, 64) + terminator)
		}

	case map[string]Message:
		switch cmd.dataType {

		case DataMap:
			buf.WriteByte(byte(PrefixMap))
			buf.WriteString(strconv.Itoa(len(data)) + terminator)

			for k, v := range data {
				buf.WriteString(k)
				d, err := Marshal(&v)
				if err != nil {
					return nil, NewErrInvalidData("invalid data")
				}

				buf.Write(d)
			}

		case DataAttributes:
			buf.WriteByte(byte(PrefixAttributes))
			buf.WriteString(strconv.Itoa(len(data)) + terminator)

			for k, v := range data {
				buf.WriteString(k)
				d, err := Marshal(&v)
				if err != nil {
					return nil, NewErrInvalidData("invalid data")
				}
				buf.Write(d)
			}
		}

	case map[Message]struct{}:
		switch cmd.dataType {
		case DataSet:
			buf.WriteByte(byte(PrefixSet))
			buf.WriteString(strconv.Itoa(len(data)) + terminator)

			for k := range data {
				d, err := Marshal(&k)
				if err != nil {
					return nil, NewErrInvalidData("invalid data")
				}
				buf.Write(d)
			}
		}
	}

	return buf.Bytes(), nil
}
