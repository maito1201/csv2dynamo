package model

import (
	"fmt"
	"strings"
)

type DynamoInput []DynamoData

func (d *DynamoInput) Copy() DynamoInput {
	cp := make([]DynamoData, len(*d))
	copy(cp, *d)
	return cp
}

func (d *DynamoInput) ToJsonString(willExecute bool) string {
	data := []string{}
	for _, v := range *d {
		data = append(data, v.ToJsonValue())
	}
	if willExecute {
		return fmt.Sprintf("{%v}", strings.Join(data, ","))
	}
	return fmt.Sprintf("'{%v}'", strings.Join(data, ","))
}

type DynamoData struct {
	Key  string
	Type string
	Val  string
}

func (d *DynamoData) ToJsonValue() string {
	var val string
	switch d.Type {
	case `"S"`:
		val = fmt.Sprintf(`"%s"`, d.Val)
	case `"N"`:
		val = fmt.Sprintf(`"%s"`, d.Val)
	default:
		val = d.Val
	}
	return fmt.Sprintf(`"%s:{%s:%s}`, d.Key, d.Type, val)
}
