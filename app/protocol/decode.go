package protocol

import "bytes"

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
