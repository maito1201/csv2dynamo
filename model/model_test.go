package model

import (
	"github.com/google/go-cmp/cmp"
	"reflect"
	"testing"
)

func TestCopy(t *testing.T) {
	d := DynamoInput{
		DynamoData{
			Key:  "SKey",
			Type: "S",
			Val:  "",
		},
	}
	cp := d.Copy()
	if diff := cmp.Diff(d, cp); diff != "" {
		t.Fatalf("copied obj mismatch (-data +copy):\n%s", diff)
	}
	if reflect.ValueOf(d).Pointer() == reflect.ValueOf(cp).Pointer() {
		t.Fatalf("Copy is not shallow copy")
	}
}

func Test_DynamoInput_ToJsonString(t *testing.T) {
	tests := []struct {
		name    string
		in      DynamoInput
		execute bool
		want    string
	}{
		{
			name: `single value "S"`,
			in: DynamoInput{
				DynamoData{
					Key:  `key_S"`,
					Type: `"S"`,
					Val:  "val_S",
				},
			},
			execute: false,
			want:    `'{"key_S":{"S":"val_S"}}'`,
		},
		{
			name: `single value "N"`,
			in: DynamoInput{
				DynamoData{
					Key:  `key_N"`,
					Type: `"N"`,
					Val:  "1",
				},
			},
			execute: false,
			want:    `'{"key_N":{"N":"1"}}'`,
		},
		{
			name: `single value "BOOL"`,
			in: DynamoInput{
				DynamoData{
					Key:  `key_BOOL"`,
					Type: `"BOOL"`,
					Val:  "true",
				},
			},
			execute: false,
			want:    `'{"key_BOOL":{"BOOL":true}}'`,
		},
		{
			name: `multiple value`,
			in: DynamoInput{
				DynamoData{
					Key:  `key_S"`,
					Type: `"S"`,
					Val:  "val_S",
				},
				DynamoData{
					Key:  `key_N"`,
					Type: `"N"`,
					Val:  "1",
				},
				DynamoData{
					Key:  `key_BOOL"`,
					Type: `"BOOL"`,
					Val:  "true",
				},
			},
			execute: false,
			want:    `'{"key_S":{"S":"val_S"},"key_N":{"N":"1"},"key_BOOL":{"BOOL":true}}'`,
		},
		{
			name: `output for dirext execute`,
			in: DynamoInput{
				DynamoData{
					Key:  `key_S"`,
					Type: `"S"`,
					Val:  "val_S",
				},
				DynamoData{
					Key:  `key_N"`,
					Type: `"N"`,
					Val:  "1",
				},
				DynamoData{
					Key:  `key_BOOL"`,
					Type: `"BOOL"`,
					Val:  "true",
				},
			},
			execute: true,
			want:    `{"key_S":{"S":"val_S"},"key_N":{"N":"1"},"key_BOOL":{"BOOL":true}}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.in.ToJsonString(tt.execute)
			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("json string mismatch (-got +want):\n%s", diff)
			}
		})
	}
}
