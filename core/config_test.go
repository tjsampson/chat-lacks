package core

import (
	"reflect"
	"testing"
)

func Test_defaultConfig(t *testing.T) {
	tests := []struct {
		name string
		want *AppConfig
	}{
		{"HappyPath", defaultConfig()},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := defaultConfig(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("defaultConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getConf(t *testing.T) {
	tests := []struct {
		name string
		want *AppConfig
	}{
		{"HappyPath", defaultConfig()},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getConf(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getConf() = %v, want %v", got, tt.want)
			}
		})
	}
}
