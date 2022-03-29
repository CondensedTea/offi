package logstf

// import (
// 	"reflect"
// 	"testing"
// )
//
// func Test_filterLogs(t *testing.T) {
// 	type args struct {
// 		maps []string
// 		logs []Log
// 	}
// 	tests := []struct {
// 		name string
// 		args args
// 		want []Log
// 	}{
// 		{
// 			name: "remove last elem",
// 			args: args{
// 				maps: []string{"cp_gullywash_final1", "cp_granary_pro_rc8"},
// 				logs: []Log{{Map: "cp_gullywash_final1"}, {Map: "cp_granary_pro_rc8"}, {Map: "cp_process_final"}},
// 			},
// 			want: []Log{{Map: "cp_gullywash_final1"}, {Map: "cp_granary_pro_rc8"}},
// 		},
// 		{
// 			name: "remove middle elem",
// 			args: args{
// 				maps: []string{"cp_gullywash_final1", "cp_process_final"},
// 				logs: []Log{{Map: "cp_gullywash_final1"}, {Map: "cp_granary_pro_rc8"}, {Map: "cp_process_final"}},
// 			},
// 			want: []Log{{Map: "cp_gullywash_final1"}, {Map: "cp_process_final"}},
// 		},
// 		{
// 			name: "remove first elem",
// 			args: args{
// 				maps: []string{"cp_granary_pro_rc8", "cp_process_final"},
// 				logs: []Log{{Map: "cp_gullywash_final1"}, {Map: "cp_granary_pro_rc8"}, {Map: "cp_process_final"}},
// 			},
// 			want: []Log{{Map: "cp_granary_pro_rc8"}, {Map: "cp_process_final"}},
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if got := filterLogs(tt.args.maps, tt.args.logs); !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("filterLogs() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }
