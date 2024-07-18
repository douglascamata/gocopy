package gocopy

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testDataFetcher struct {
	filename string
}

func (f testDataFetcher) Fetch() ([]byte, error) {
	file, err := os.ReadFile(filepath.Join("testdata", f.filename))
	return file, err
}

func mustReadTestData(filename string) []byte {
	file, err := os.ReadFile(filepath.Join("testdata", filename))
	if err != nil {
		panic(err)
	}
	return bytes.TrimSpace(file)
}

func TestCopyFunction(t *testing.T) {
	type args struct {
		fetcher sourceFetcher
		fnName  string
	}
	tests := []struct {
		name         string
		args         args
		wantErr      error
		wantFunction []byte
	}{
		{
			name: "Copy existing function",
			args: args{
				fetcher: testDataFetcher{filename: "decode.go.test"},
				fnName:  "Unmarshal",
			},
			wantErr:      nil,
			wantFunction: mustReadTestData("decode_unmarshal.go.test"),
		},
		{
			name: "Copy non-existing function",
			args: args{
				fetcher: testDataFetcher{filename: "decode.go.test"},
				fnName:  "Foobar",
			},
			wantErr: ErrFunctionNotFound{functionName: "Foobar"},
		},
		{
			name: "Copy from invalid syntax file",
			args: args{
				fetcher: testDataFetcher{filename: "bad_decode.go.test"},
				fnName:  "Foobar",
			},
			wantErr: ErrInvalidGoSource{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out, err := CopyFunction(tt.args.fetcher, tt.args.fnName)
			if tt.wantErr != nil {
				assert.ErrorContains(t, err, tt.wantErr.Error())
				return
			}
			assert.NoError(t, err)
			assert.EqualValues(t, string(tt.wantFunction), out)
		})
	}
}

func TestCopyTypeFunction(t *testing.T) {
	type args struct {
		fetcher        sourceFetcher
		typeName       string
		includeMethods bool
	}
	tests := []struct {
		name     string
		args     args
		wantErr  error
		wantType []byte
	}{
		{
			name: "Copy existing type",
			args: args{
				fetcher:        testDataFetcher{filename: "decode.go.test"},
				typeName:       "Number",
				includeMethods: true,
			},
			wantErr:  nil,
			wantType: mustReadTestData("decode_number.go.test"),
		},
		{
			name: "Copy non-existing type",
			args: args{
				fetcher:        testDataFetcher{filename: "decode.go.test"},
				typeName:       "Foobar",
				includeMethods: true,
			},
			wantErr: ErrTypeNotFound{typeName: "Foobar"},
		},
		{
			name: "Copy from invalid syntax file",
			args: args{
				fetcher:        testDataFetcher{filename: "bad_decode.go.test"},
				typeName:       "Foobar",
				includeMethods: true,
			},
			wantErr: ErrInvalidGoSource{},
		},
		{
			name: "Parses grouped types correctly",
			args: args{
				fetcher:        testDataFetcher{filename: "grouped_types.go.test"},
				typeName:       "A",
				includeMethods: true,
			},
			wantType: mustReadTestData("grouped_types_a.go.test"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out, err := CopyType(tt.args.fetcher, tt.args.typeName, tt.args.includeMethods)
			if tt.wantErr != nil {
				assert.ErrorContains(t, err, tt.wantErr.Error())
				return
			}
			assert.NoError(t, err)
			assert.EqualValues(t, string(tt.wantType), strings.TrimSpace(out))
		})
	}
}
