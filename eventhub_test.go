package broadcast

import (
	"reflect"
	"testing"
)

func TestNewBroadcaster(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"check new broadcaster",
			args{name: "broadcaster"},
			false,
		},
		{
			"check existed broadcaster",
			args{name: "broadcaster"},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := NewBroadcaster(tt.args.name); (err != nil) != tt.wantErr {
				t.Errorf("NewBroadcaster() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetBroadcaster(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name string
		args args
		want *eventHub
	}{
		{
			"get broadcaster",
			args{name: "broadcaster"},
			false,
		},
		{
			"check existed broadcaster",
			args{name: "broadcaster"},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetBroadcaster(tt.args.name); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetBroadcaster() = %v, want %v", got, tt.want)
			}
		})
	}
}
