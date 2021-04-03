package csv2dynamo

import (
	"errors"
	"github.com/google/go-cmp/cmp"
	"io/fs"
	"testing"
)

func TestReadCSV(t *testing.T) {
	tests := []struct {
		name    string
		in      string
		want    []DynamoInput
		wantErr error
	}{
		{
			name: "success",
			in:   "./testdata/sample.csv",
			want: []DynamoInput{
				DynamoInput{
					DynamoData{Key: `s_value"`, Type: `"S"`, Val: "sample1"},
					DynamoData{Key: `n_value"`, Type: `"N"`, Val: "1"},
					DynamoData{Key: `bool_value"`, Type: `"BOOL"`, Val: "true"},
				},
				DynamoInput{
					DynamoData{Key: `s_value"`, Type: `"S"`, Val: "sample2"},
					DynamoData{Key: `n_value"`, Type: `"N"`, Val: "2"},
					DynamoData{Key: `bool_value"`, Type: `"BOOL"`, Val: "false"},
				},
			},
			wantErr: nil,
		},
		{
			name:    "invalid header",
			in:      "./testdata/invalid_header.csv",
			want:    nil,
			wantErr: errUnexpectedHeader,
		},
		{
			name:    "file missing",
			in:      "./testdata/not_exist.csv",
			want:    nil,
			wantErr: &fs.PathError{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := readCSV(tt.in)
			werr := tt.wantErr
			if err != nil && !errors.As(err, &werr) {
				t.Fatalf("unexpected error: %v %T, want %v", err, err, tt.wantErr)
			}

			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("read value is mismatch (-got +want):\n%s", diff)
			}
		})
	}
}
