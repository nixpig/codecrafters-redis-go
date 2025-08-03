package protocol

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type scenario struct {
	dataType   DataType
	dataPrefix DataPrefix
	data       any
	expected   []byte
}

func TestMarshalData(t *testing.T) {
	scenarios := map[string]scenario{
		"marshal bulk string": {
			dataType: DataBulkString,
			data:     "a bulk string",
			expected: []byte("$13\r\na bulk string\r\n"),
		},
		"marshal simple string": {
			dataType: DataSimpleString,
			data:     "a simple string",
			expected: []byte("+a simple string\r\n"),
		},
		"marshal error": {
			dataType: DataError,
			data:     "an error",
			expected: []byte("-an error\r\n"),
		},
		"marshal bulk error": {
			dataType: DataBulkError,
			data:     "a bulk error",
			expected: []byte("!12\r\na bulk error\r\n"),
		},
		"marshal big number": {
			dataType: DataBigNumber,
			data:     "123456789012345678901234567890",
			expected: []byte("(123456789012345678901234567890\r\n"),
		},
		"marshal integer": {
			dataType: DataInteger,
			data:     69,
			expected: []byte(":69\r\n"),
		},
		"marshal bool (true)": {
			dataType: DataBool,
			data:     true,
			expected: []byte("#t\r\n"),
		},
		"marshal bool (false)": {
			dataType: DataBool,
			data:     false,
			expected: []byte("#f\r\n"),
		},
		"marshal double": {
			dataType: DataDouble,
			data:     float64(69.23),
			expected: []byte(",69.23\r\n"),
		},
		"marshal array": {
			dataType: DataArray,
			data: []Message{
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
			expected: []byte(
				"*3\r\n+some simple string\r\n#t\r\n!15\r\nsome bulk error\r\n",
			),
		},
		"marshal push": {
			dataType: DataPush,
			data: []Message{
				{
					dataType: DataBulkString,
					data:     "a nested bulk string",
				},
				{
					dataType: DataInteger,
					data:     69,
				},
			},
			expected: []byte(">2\r\n$20\r\na nested bulk string\r\n:69\r\n"),
		},
		"marshal map": {
			dataType: DataMap,
			data: map[string]Message{
				"simple string": {
					dataType: DataSimpleString,
					data:     "some simple string",
				},
				"integer": {
					dataType: DataInteger,
					data:     23,
				},
			},
			expected: []byte(
				"%2\r\nsimple string+some simple string\r\ninteger:23\r\n",
			),
		},
		"marshal attributes": {
			dataType: DataAttributes,
			data: map[string]Message{
				"bulkstringkey": {
					dataType: DataBulkString,
					data:     "a bulk string",
				},
				"integerkey": {
					dataType: DataInteger,
					data:     23,
				},
			},
			expected: []byte(
				"|2\r\nbulkstringkey$13\r\na bulk string\r\nintegerkey:23\r\n",
			),
		},
		"marshal set": {
			dataType: DataSet,
			data: map[Message]struct{}{
				{
					dataType: DataBulkString,
					data:     "a bulk string",
				}: {},
				{
					dataType: DataInteger,
					data:     13,
				}: {},
			},
			expected: []byte(
				"~2\r\n$13\r\na bulk string\r\n:13\r\n",
			),
		},
	}

	for name, d := range scenarios {
		t.Run(name, func(t *testing.T) {
			testMarshalData(t, d.dataType, d.data, d.expected)
		})
	}
}

func testMarshalData(
	t *testing.T,
	dataType DataType,
	data any,
	expected []byte,
) {
	msg := &Message{dataType, data}
	b, err := Marshal(msg)
	require.NoError(t, err)
	require.Equal(t, expected, b)
}
