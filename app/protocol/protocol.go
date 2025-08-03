package protocol

import (
	"bytes"
	"fmt"
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

// Unmarshal parses the RESP-encoded byte data in b and returns a Command and error
func Unmarshal(b []byte) (*Message, error) {
	if len(b) == 0 {
		return nil, fmt.Errorf("no data")
	}

	commandType, ok := dataMap[DataPrefix(b[0])]
	if !ok {
		return nil, fmt.Errorf("invalid prefix")
	}

	var data any

	switch commandType {
	case DataBulkString:
		data = ""

	case DataSimpleString:
		data = ""

	case DataArray:
		data = []Message{}
	}

	return &Message{
		dataType: commandType,
		data:     data,
	}, nil
}

// Marshal encodes Command to RESP-encoded byte slice
func Marshal(cmd *Message) ([]byte, error) {
	switch data := (cmd.data).(type) {
	case string:
		switch cmd.dataType {
		case DataBulkString:
			return []byte(
				"$" + strconv.Itoa(len(data)) + "\r\n" + data + "\r\n",
			), nil

		case DataSimpleString:
			return []byte("+" + data + "\r\n"), nil

		case DataError:
			return []byte("-" + data + "\r\n"), nil

		case DataBulkError:
			return []byte(
				"!" + strconv.Itoa(len(data)) + "\r\n" + data + "\r\n",
			), nil

		case DataBigNumber:
			return []byte("(" + data + "\r\n"), nil
		}

	case []Message:
		switch cmd.dataType {
		case DataArray:
			var buf bytes.Buffer

			buf.WriteString("*" + strconv.Itoa(len(data)) + "\r\n")

			for _, item := range data {
				b, err := Marshal(&item)
				if err != nil {
					return nil, fmt.Errorf("invalid data")
				}
				buf.Write(b)
			}

			return buf.Bytes(), nil

		case DataPush:
			data, ok := (cmd.data).([]Message)
			if !ok {
				return nil, fmt.Errorf("invalid data")
			}

			var buf bytes.Buffer

			buf.WriteString(">" + strconv.Itoa(len(data)) + "\r\n")

			for _, item := range data {
				d, err := Marshal(&item)
				if err != nil {
					return nil, fmt.Errorf("invalid data")
				}

				buf.Write(d)
			}

			return buf.Bytes(), nil
		}

	case int:
		switch cmd.dataType {
		case DataInteger:
			return []byte(":" + strconv.Itoa(data) + "\r\n"), nil
		}

	case bool:
		switch cmd.dataType {
		case DataBool:
			var v string
			if data {
				v = "t"
			} else {
				v = "f"
			}

			return []byte("#" + v + "\r\n"), nil
		}

	case float64:
		switch cmd.dataType {
		case DataDouble:
			return []byte(
				"," + strconv.FormatFloat(data, 'G', -1, 64) + "\r\n",
			), nil
		}

	case map[string]Message:
		switch cmd.dataType {

		case DataMap:
			var buf bytes.Buffer

			buf.WriteString("%" + strconv.Itoa(len(data)) + "\r\n")

			for k, v := range data {
				buf.WriteString(k)
				d, err := Marshal(&v)
				if err != nil {
					return nil, fmt.Errorf("invalid data")
				}

				buf.Write(d)
			}

			return buf.Bytes(), nil

		case DataAttributes:
			var buf bytes.Buffer

			buf.WriteString("|" + strconv.Itoa(len(data)) + "\r\n")

			for k, v := range data {
				buf.WriteString(k)
				d, err := Marshal(&v)
				if err != nil {
					return nil, fmt.Errorf("invalid data")
				}
				buf.Write(d)
			}

			return buf.Bytes(), nil
		}

	case map[string]struct{}:
		switch cmd.dataType {
		case DataSet:
			var buf bytes.Buffer

			buf.WriteString("~" + strconv.Itoa(len(data)) + "\r\n")

			for k := range data {
				buf.WriteString(k)
			}

			return buf.Bytes(), nil
		}
	}

	return nil, fmt.Errorf("invalid command")
}
