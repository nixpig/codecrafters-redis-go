package protocol

import (
	"strconv"
	"strings"
)

type RESP interface {
	String() string
}

// ---

type BulkString struct {
	data string
}

func (b *BulkString) String() string {
	return "$" + strconv.Itoa(len(b.data)) + "\r\n" + b.data + "\r\n"
}

func NewBulkString(data string) *BulkString {
	return &BulkString{data}
}

// ---

func NewNullBulkString() string {
	return "$-1\r\n"
}

func NewSimpleString(data string) string {
	return "+" + data + "\r\n"
}

func NewError(prefix, msg string) string {
	return "-" + prefix + " " + msg + "\r\n"
}

func NewBulkError(prefix, msg string) string {
	return "!" + strconv.Itoa(
		len(prefix)+len(msg)+1,
	) + "\r\n" + prefix + " " + msg + "\r\n"
}

func NewInteger(num int) string {
	return ":" + strconv.Itoa(num) + "\r\n"
}

func NewArray(items []RESP) string {
	var sb strings.Builder

	sb.WriteString("*" + strconv.Itoa(len(items)) + "\r\n")

	for _, item := range items {
		sb.WriteString(item.String())
	}

	return sb.String()
}

func NewNullArray() string {
	return "*-1\r\n"
}

func NewNull() string {
	return "_\r\n"
}

func NewBool(val bool) string {
	var v string
	if val {
		v = "t"
	} else {
		v = "f"
	}

	return "#" + v + "\r\n"
}

func NewDouble(val float64) string {
	return "," + strconv.FormatFloat(val, 'G', -1, 64) + "\r\n"
}

func NewPositiveInf() string {
	return ",inf\r\n"
}

func NewNegativeInf() string {
	return ",-inf\r\n"
}

func NewNaN() string {
	return ",nan\r\n"
}

func NewBigNumber(val string) string {
	return "(" + val + "\r\n"
}

func NewVerbatimString(enc, val string) string {
	return "=" + strconv.Itoa(
		len(enc)+len(val)+1,
	) + "\r\n" + enc + ":" + val + "\r\n"
}

func NewMap(items map[string]RESP) string {
	var sb strings.Builder

	sb.WriteString("%" + strconv.Itoa(len(items)) + "\r\n")

	for k, v := range items {
		sb.WriteString(k + v.String())
	}

	return sb.String()
}

func NewAttributes(items map[string]RESP) string {
	var sb strings.Builder

	sb.WriteString("|" + strconv.Itoa(len(items)) + "\r\n")

	for k, v := range items {
		sb.WriteString(k + v.String())
	}

	return sb.String()
}

func NewSet(items map[RESP]struct{}) string {
	var sb strings.Builder

	sb.WriteString("~" + strconv.Itoa(len(items)) + "\r\n")

	for k := range items {
		sb.WriteString(k.String())
	}

	return sb.String()
}

func NewPush(items ...RESP) string {
	var sb strings.Builder

	sb.WriteString(">" + strconv.Itoa(len(items)) + "\r\n")

	for _, item := range items {
		sb.WriteString(item.String())
	}

	return sb.String()
}
