package protocol

import (
	"bytes"
	"strconv"
)

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
