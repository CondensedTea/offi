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
					{Id: 1, Map: "cp_process_final", Date: 1650744000},
					{Id: 2, Map: "cp_gullywash_final1", Date: 1650745800},
				},
				playedAt: time.Date(2022, time.April, 24, 1, 0, 0, 0, time.UTC),
			},
			wantMatchLogs: []Log{
				{Id: 1, Map: "cp_process_final", Date: 1650744000},
				{Id: 2, Map: "cp_gullywash_final1", Date: 1650745800},
			},
		},
		{
			name: "filter combined log by map",
			args: args{
				maps: []string{"cp_process_final", "cp_gullywash_final1"},
				logs: []Log{
					{Id: 1, Map: "cp_process_final", Date: 1650744000},
					{Id: 2, Map: "gully", Date: 1650745800},
				},
				playedAt: time.Date(2022, time.April, 24, 1, 0, 0, 0, time.UTC),
			},
			wantMatchLogs: []Log{
				{Id: 1, Map: "cp_process_final", Date: 1650744000},
			},
			wantCombinedLogs: []Log{
				{Id: 2, Map: "gully", Date: 1650745800},
			},
		},
		{
			name: "drop log as too old",
			args: args{
				maps: []string{"cp_process_final", "cp_gullywash_final1"},
				logs: []Log{
					{Id: 1, Map: "cp_process_final", Date: 1650744000},
					{Id: 2, Map: "cp_gullywash_final1", Date: 1650400200},
				},
				playedAt: time.Date(2022, time.April, 24, 1, 0, 0, 0, time.UTC),
			},
			wantMatchLogs: []Log{
				{Id: 1, Map: "cp_process_final", Date: 1650744000},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotMatchLogs, gotCombinedLogs := filterLogs(tt.args.maps, tt.args.logs, tt.args.playedAt)
			assert.Equal(t, gotMatchLogs, tt.wantMatchLogs)
			assert.Equal(t, gotCombinedLogs, tt.wantCombinedLogs)
		})
	}
}
