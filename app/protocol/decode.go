package protocol

import (
	"bytes"
	"strconv"
)

// Unmarshal parses the RESP-encoded byte data in b and returns a Command and error
func Unmarshal(b []byte) (*Message, error) {
	var err error

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

	switch commandType {
	case DataSimpleString:
		data, err = unmarshalSimpleString(b[1:])

	case DataError:
		data, err = unmarshalError(b[1:])

	case DataBulkString:
		data, err = unmarshalBulkString(b[1:])

	case DataBulkError:
		data, err = unmarshalBulkError(b[1:])

	case DataBigNumber:
		data, err = unmarshalBigNumber(b[1:])

	case DataInteger:
		data, err = unmarshalInteger(b[1:])

	case DataBool:
		data, err = unmarshalBool(b[1:])

	case DataDouble:
		data, err = unmarshalDouble(b[1:])

	}

	if err != nil {
		return nil, err
	}

	return &Message{
		dataType: commandType,
		data:     data,
	}, nil
}

func unmarshalSimpleString(b []byte) (string, error) {
	p, ok := bytes.CutSuffix(b, []byte(terminator))
	if !ok {
		return "", NewErrInvalidData("invalid simple string data")
	}

	return string(p), nil
}

func unmarshalError(b []byte) (string, error) {
	p, ok := bytes.CutSuffix(b, []byte(terminator))
	if !ok {
		return "", NewErrInvalidData("invalid error data")
	}

	return string(p), nil
}

func unmarshalBulkString(b []byte) (string, error) {
	parts := bytes.Split(b, []byte(terminator))

	if len(parts) < 2 {
		return "", NewErrInvalidData("invalid bulk string data")
	}

	return string(parts[1]), nil
}

func unmarshalBulkError(b []byte) (string, error) {
	parts := bytes.Split(b, []byte(terminator))

	if len(parts) < 2 {
		return "", NewErrInvalidData("invalid bulk error data")
	}

	return string(parts[1]), nil
}

func unmarshalBigNumber(b []byte) (string, error) {
	p, ok := bytes.CutSuffix(b, []byte(terminator))
	if !ok {
		return "", NewErrInvalidData("invalid big number data")
	}

	return string(p), nil
}

func unmarshalInteger(b []byte) (int, error) {
	p, ok := bytes.CutSuffix(b, []byte(terminator))
	if !ok {
		return 0, NewErrInvalidData("invalid integer data")
	}

	i, err := strconv.Atoi(string(p))
	if err != nil {
		return 0, NewErrInvalidData("invalid integer data: " + err.Error())
	}

	return i, nil
}

func unmarshalBool(b []byte) (bool, error) {
	p, ok := bytes.CutSuffix(b, []byte(terminator))
	if !ok {
		return false, NewErrInvalidData("invalid bool data")
	}

	if string(p) == "t" || string(p) == "T" {
		return true, nil
	}

	if string(p) == "f" || string(p) == "F" {
		return false, nil
	}

	return false, NewErrInvalidData("invalid bool data: " + string(p))
}

func unmarshalDouble(b []byte) (float64, error) {
	p, ok := bytes.CutSuffix(b, []byte(terminator))
	if !ok {
		return 0, NewErrInvalidData("invalid double data")
	}

	f, err := strconv.ParseFloat(string(p), 64)
	if err != nil {
		return 0, NewErrInvalidData("invalid double data: " + err.Error())
	}

	return f, nil
}

func unmarshalArray(b []byte) ([]Message, error) {

}
