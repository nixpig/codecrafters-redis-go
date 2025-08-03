package protocol

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type unmarshalScenario struct {
	data             []byte
	expectedDataType DataType
	expectedData     any
}

func TestUnmarshalData(t *testing.T) {
	scenarios := map[string]unmarshalScenario{
		"unmarshal bulk string": {
			data:             []byte("$13\r\na bulk string\r\n"),
			expectedDataType: DataBulkString,
			expectedData:     "a bulk string",
		},
		"unmarshal simple string": {
			data:             []byte("+a simple string\r\n"),
			expectedDataType: DataSimpleString,
			expectedData:     "a simple string",
		},
		"unmarshal error": {
			data:             []byte("-an error\r\n"),
			expectedDataType: DataError,
			expectedData:     "an error",
		},
		"unmarshal bulk error": {
			data:             []byte("!12\r\na bulk error\r\n"),
			expectedDataType: DataBulkError,
			expectedData:     "a bulk error",
		},
		"unmarshal big number": {
			data:             []byte("(123456789012345678901234567890\r\n"),
			expectedDataType: DataBigNumber,
			expectedData:     "123456789012345678901234567890",
		},
		"unmarshal integer": {
			data:             []byte(":69\r\n"),
			expectedDataType: DataInteger,
			expectedData:     69,
		},
		"unmarshal bool (true)": {
			data:             []byte("#t\r\n"),
			expectedDataType: DataBool,
			expectedData:     true,
		},
		"unmarshal bool (false)": {
			data:             []byte("#f\r\n"),
			expectedDataType: DataBool,
			expectedData:     false,
		},
		"unmarshal double": {
			data:             []byte(",69.23\r\n"),
			expectedDataType: DataDouble,
			expectedData:     float64(69.23),
		},
		"unmarshal array": {
			data: []byte(
				"*3\r\n+some simple string\r\n#t\r\n!15\r\nsome bulk error\r\n",
			),
			expectedDataType: DataArray,
			expectedData: []Message{
				{
					dataType: DataSimpleString,
					data:     "some simple string",
				},
				{
					dataType: DataBool,
					data:     true,
				},
				{
					dataType: DataBulkError,
					data:     "some bulk error",
				},
			},
		},
		// "unmarshal push": {
		// 	data: []byte(
		// 		">2\r\n$20\r\na nested bulk string\r\n:69\r\n",
		// 	),
		// 	expectedDataType: DataPush,
		// 	expectedData: []Message{
		// 		{
		// 			dataType: DataBulkString,
		// 			data:     "a nested bulk string",
		// 		},
		// 		{
		// 			dataType: DataInteger,
		// 			data:     69,
		// 		},
		// 	},
		// },
		// "unmarshal map": {
		// 	data: []byte(
		// 		"%2\r\nsimple string+some simple string\r\ninteger:23\r\n",
		// 	),
		// 	expectedDataType: DataMap,
		// 	expectedData: map[string]Message{
		// 		"simple string": {
		// 			dataType: DataSimpleString,
		// 			data:     "some simple string",
		// 		},
		// 		"integer": {
		// 			dataType: DataInteger,
		// 			data:     23,
		// 		},
		// 	},
		// },
		// "unmarshal attributes": {
		// 	data: []byte(
		// 		"|2\r\nbulkstringkey$13\r\na bulk string\r\nintegerkey:23\r\n",
		// 	),
		// 	expectedDataType: DataAttributes,
		// 	expectedData: map[string]Message{
		// 		"bulkstringkey": {
		// 			dataType: DataBulkString,
		// 			data:     "a bulk string",
		// 		},
		// 		"integerkey": {
		// 			dataType: DataInteger,
		// 			data:     23,
		// 		},
		// 	},
		// },
		// "unmarshal set": {
		// 	data: []byte(
		// 		"~2\r\n$13\r\na bulk string\r\n:13\r\n",
		// 	),
		// 	expectedDataType: DataSet,
		// 	expectedData: map[Message]struct{}{
		// 		{
		// 			dataType: DataBulkString,
		// 			data:     "a bulk string",
		// 		}: {},
		// 		{
		// 			dataType: DataInteger,
		// 			data:     13,
		// 		}: {},
		// 	},
		// },
		// "unmarshal empty": {
		// 	data:             []byte(nil),
		// 	expectedDataType: DataType(-1),
		// 	expectedData:     nil,
		// },
	}

	for name, d := range scenarios {
		t.Run(name, func(t *testing.T) {
			testUnmarshalData(t, d.data, d.expectedDataType, d.expectedData)
		})
	}
}

func testUnmarshalData(
	t *testing.T,
	data []byte,
	expectedDataType DataType,
	expectedData any,
) {
	msg, err := Unmarshal(data)
	require.NoError(t, err)
	require.Equal(t, expectedDataType, msg.Type())
	require.Equal(t, expectedData, msg.Data())
}
