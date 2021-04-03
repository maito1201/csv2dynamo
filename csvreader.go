package csv2dynamo

import (
	"encoding/csv"
	"errors"
	"fmt"
	"github.com/maito1201/clearout"
	"os"
	"regexp"
	"strings"
)

var (
	Input               []DynamoInput
	InputTemplate       DynamoInput
	errUnexpectedHeader = errors.New("unexpected header format")
)

func readCSV(path string) ([]DynamoInput, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	scanner := csv.NewReader(f)
	records, err := scanner.ReadAll()
	if err != nil {
		return nil, err
	}
	if err := readHeader(records[0]); err != nil {
		return nil, err
	}
	cout := clearout.Output{Prefix: "read and compile csv\n"}
	for i := 1; i < len(records); i++ {
		cout.Printf("progress: %d/%d\n", i, len(records)-1)
		if i == len(records)-1 {
			cout.Println("complete!")
		}
		cout.Render()
		readValue(records[i])
	}
	return Input, nil
}

func readHeader(ss []string) error {
	for _, v := range ss {
		ss := strings.Split(v, " ")
		if len(ss) != 2 {
			return fmt.Errorf("%v %w", ss, errUnexpectedHeader)
		}
		k := fmt.Sprintf(`%s"`, ss[0])
		InputTemplate = append(InputTemplate, DynamoData{Key: k, Type: normalizeType(ss[1])})
	}
	return nil
}

func readValue(ss []string) {
	data := InputTemplate.Copy()
	for i := 0; i < len(ss); i++ {
		data[i].Val = ss[i]
	}
	Input = append(Input, data)
}

var typeRegex = regexp.MustCompile("[A-Z]+")

func normalizeType(s string) string {
	return fmt.Sprintf(`"%s"`, typeRegex.FindString(s))
}
