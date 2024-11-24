package etf2l

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_parseMatchDate(t *testing.T) {
	timeNow = func() time.Time {
		return time.Date(2022, 4, 2, 12, 0, 0, 0, time.Local)
	}

	type args struct {
		textBlock string
	}
	tests := []struct {
		name    string
		args    args
		want    time.Time
		wantErr bool
	}{
		{
			name: "default",
			args: args{textBlock: "Results submitted: 2 Apr 2022, 12:30"},
			want: time.Date(2022, 4, 2, 12, 30, 0, 0, time.Local),
		},
		{
			name: "today as date",
			args: args{textBlock: "Results submitted: Today, 12:30"},
			want: time.Date(2022, 4, 2, 12, 30, 0, 0, time.Local),
		},
		{
			name: "yesterday as date",
			args: args{textBlock: "Results submitted: Yesterday, 12:30"},
			want: time.Date(2022, 4, 1, 12, 30, 0, 0, time.Local),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseMatchDate(tt.args.textBlock)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseMatchDate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.Equal(t, tt.want, got) {
				t.Errorf("parseMatchDate() got = %v, want %v", got, tt.want)
			}
		})
	}
}
