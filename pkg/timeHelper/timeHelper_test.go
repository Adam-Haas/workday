package timeHelper

import (
	"testing"
	"time"
)

func TestIsDayBeforeToday(t *testing.T) {
	tests := []struct {
		name string
		time1 time.Time
		time2 time.Time
		want bool
	}{
		{
			name: "time is day before",
			time1: time.Date(2021, 9, 11, 9, 10, 0, 0, time.Local),
			time2: time.Date(2021, 9, 12, 9, 15, 0, 0, time.Local),
			want: true,
		},
		{
			name: "time is 3 days before",
			time1: time.Date(2021, 9, 9, 9, 10, 0, 0, time.Local),
			time2: time.Date(2021, 9, 12, 9, 15, 0, 0, time.Local),
			want: true,
		},
		{
			name: "time is same day but time before",
			time1: time.Date(2021, 9, 12, 9, 10, 0, 0, time.Local),
			time2: time.Date(2021, 9, 12, 9, 15, 0, 0, time.Local),
			want: false,
		},
		{
			name: "time is same day but time after",
			time1: time.Date(2021, 9, 12, 9, 30, 0, 0, time.Local),
			time2: time.Date(2021, 9, 12, 9, 15, 0, 0, time.Local),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := TimeIsDayBeforeTime(tt.time1, tt.time2); got != tt.want {
				t.Errorf("IsDayBeforeToday() = %v, want %v", got, tt.want)
			}
		})
	}
}
