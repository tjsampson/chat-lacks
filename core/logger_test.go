package core

import (
	"log"
	"os"
	"reflect"
	"testing"
)

func Test_defaultLogger(t *testing.T) {
	tests := []struct {
		name string
		want *log.Logger
	}{
		{"HappyPath", FakeDefaultLogger()},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := defaultLogger(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("defaultLogger() = %v, want %v", got, tt.want)
			}
		})
	}
}

func FakeDefaultLogger() *log.Logger {
	return log.New(os.Stdout, prefix, loggerFlag)
}

func TestCreateLogger(t *testing.T) {
	type args struct {
		logOutput string
	}
	tests := []struct {
		name  string
		args  args
		want  *log.Logger
		want1 *os.File
	}{
		{"Empty LogOutput", args{logOutput: ""}, FakeDefaultLogger(), nil},
		{"Stdout LogOutput", args{logOutput: "stdout"}, FakeDefaultLogger(), nil},
	}

	t.Log("Given the need to Create a Logger")
	{
		for _, tt := range tests {
			t.Logf("\tTest: %s\tWhen checking if log output = %q then result = %v", tt.name, tt.args.logOutput, tt.want)
			{
				t.Run(tt.name, func(t *testing.T) {
					got, got1 := CreateLogger(tt.args.logOutput)
					if !reflect.DeepEqual(got, tt.want) {
						t.Errorf("Create() got = %v, want %v", got, tt.want)
					}
					if !reflect.DeepEqual(got1, tt.want1) {
						t.Errorf("Create() got1 = %v, want %v", got1, tt.want1)
					}
				})
			}
		}
	}
}
