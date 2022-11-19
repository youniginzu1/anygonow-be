package lib

import (
	"testing"
	"time"
)

func TestGetContentType(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "a.jpeg",
			args: args{
				name: "a.jpeg",
			},
			want:    "image/jpeg",
			wantErr: false,
		},
		{
			name: "a.jpeg",
			args: args{
				name: "a.png",
			},
			want:    "image/png",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetContentType(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetContentType() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetContentType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWeekDay(t *testing.T) {
	left, right := GetTimeRange(3)
	t.Log(time.UnixMilli(left))
	t.Log(time.UnixMilli(right))
}
