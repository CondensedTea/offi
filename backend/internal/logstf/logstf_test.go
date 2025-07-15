package logstf

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_filterLogs(t *testing.T) {
	type args struct {
		maps     []string
		logs     []Log
		playedAt time.Time
	}
	tests := []struct {
		name             string
		args             args
		wantMatchLogs    []Log
		wantCombinedLogs []Log
	}{
		{
			name: "all valid",
			args: args{
				maps: []string{"cp_process_final", "cp_gullywash_final1"},
				logs: []Log{
					{ID: 1, Map: "cp_process_final", Date: 1650744000},
					{ID: 2, Map: "cp_gullywash_final1", Date: 1650745800},
				},
				playedAt: time.Date(2022, time.April, 24, 1, 0, 0, 0, time.UTC),
			},
			wantMatchLogs: []Log{
				{ID: 1, Map: "cp_process_final", Date: 1650744000},
				{ID: 2, Map: "cp_gullywash_final1", Date: 1650745800},
			},
		},
		{
			name: "filter combined log by map",
			args: args{
				maps: []string{"cp_process_final", "cp_gullywash_final1"},
				logs: []Log{
					{ID: 1, Map: "cp_process_final", Date: 1650744000},
					{ID: 2, Map: "gully", Date: 1650745800},
				},
				playedAt: time.Date(2022, time.April, 24, 1, 0, 0, 0, time.UTC),
			},
			wantMatchLogs: []Log{
				{ID: 1, Map: "cp_process_final", Date: 1650744000},
			},
			wantCombinedLogs: []Log{
				{ID: 2, Map: "gully", Date: 1650745800},
			},
		},
		{
			name: "drop log as too old",
			args: args{
				maps: []string{"cp_process_final", "cp_gullywash_final1"},
				logs: []Log{
					{ID: 1, Map: "cp_process_final", Date: 1650659400},
					{ID: 2, Map: "cp_gullywash_final1", Date: 1647376200}, //  1650054600
				}, playedAt: time.Unix(1650409200, 0),
			},
			wantMatchLogs: []Log{
				{ID: 1, Map: "cp_process_final", Date: 1650659400},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotMatchLogs, gotCombinedLogs := filterLogs(tt.args.maps, tt.args.logs, tt.args.playedAt)
			assert.Equal(t, tt.wantMatchLogs, gotMatchLogs, "primary logs doesnt match")
			assert.Equal(t, tt.wantCombinedLogs, gotCombinedLogs, "secondary logs doesnt match")
		})
	}
}

func Test_mapIsNotValid(t *testing.T) {
	type args struct {
		maps   []string
		logMap string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "complete map match",
			args: args{
				maps:   []string{"cp_badlands", "cp_process_final1"},
				logMap: "cp_badlands",
			},
			want: false,
		},
		{
			name: "wrong version but still match",
			args: args{
				maps:   []string{"cp_badlands", "cp_process_final2"},
				logMap: "cp_process_rc3",
			},
			want: false,
		},
		{
			name: "no match",
			args: args{
				maps:   []string{"cp_badlands", "cp_process_final2"},
				logMap: "combined log bad + proc",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, mapIsNotValid(tt.args.maps, tt.args.logMap), "mapIsNotValid(%v, %v)", tt.args.maps, tt.args.logMap)
		})
	}
}
