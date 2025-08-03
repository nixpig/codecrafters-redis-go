package protocol

type DataPrefix byte
type DataType int

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
