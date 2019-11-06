package utils

import (
	"github.com/savaki/jq"
	"strconv"
)

func GetJsonStringUnquoteAttribute(exp string, data string) string {
	op, _ := jq.Parse(exp)
	value, _ := op.Apply([]byte(data))
	unquoteValue, _ := strconv.Unquote(string(value))

	return unquoteValue
}

func GetJsonAttribute(exp string, data string) []byte {
	op, _ := jq.Parse(exp)
	value, _ := op.Apply([]byte(data))

	return value
}

func GetJsonStringAttribute(exp string, data string) string {
	op, _ := jq.Parse(exp)
	value, _ := op.Apply([]byte(data))

	return TrimQuotes(string(value))
}