package tests

import (
	"dds-backend/models"
	"reflect"
	"testing"
)

func TestLoadWeekdays(t *testing.T) {
	type args struct {
		in string
	}
	tests := []struct {
		name    string
		args    args
		want    []models.Weekday
		wantErr bool
	}{
		{"", args{"1,3,5"}, []models.Weekday{models.Monday, models.Wednesday, models.Friday}, false},
		{"", args{""}, []models.Weekday{}, false},
		{"", args{"7,6,5,4,3,2,1"}, []models.Weekday{models.Monday, models.Tuesday, models.Wednesday, models.Thursday, models.Friday, models.Saturday, models.Sunday}, false},
		{"", args{"1,,2"}, nil, true},
		{"", args{","}, nil, true},
		{"", args{"100"}, nil, true},
		{"", args{"3"}, []models.Weekday{models.Wednesday}, false},
		{"", args{"6,6,6"}, []models.Weekday{models.Saturday}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := models.LoadWeekdays(tt.args.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadWeekdays() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LoadWeekdays() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPrettyWeekdays(t *testing.T) {
	type args struct {
		days []models.Weekday
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"", args{[]models.Weekday{models.Monday, models.Wednesday, models.Friday}}, "Monday Wednesday Friday"},
		{"", args{[]models.Weekday{models.Friday}}, "Friday"},
		{"", args{[]models.Weekday{models.Sunday, models.Friday, models.Saturday}}, "Friday Saturday Sunday"},
		{"", args{[]models.Weekday{}}, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := models.PrettyWeekdays(tt.args.days); got != tt.want {
				t.Errorf("PrettyWeekdays() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStoreWeekdays(t *testing.T) {
	type args struct {
		days []models.Weekday
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"", args{[]models.Weekday{models.Monday, models.Wednesday, models.Friday}}, "1,3,5"},
		{"", args{[]models.Weekday{models.Monday, models.Tuesday, models.Wednesday, models.Thursday, models.Friday, models.Saturday, models.Sunday}}, "1,2,3,4,5,6,7"},
		{"", args{[]models.Weekday{models.Monday, models.Friday, models.Wednesday}}, "1,3,5"},
		{"", args{[]models.Weekday{}}, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := models.StoreWeekdays(tt.args.days); got != tt.want {
				t.Errorf("StoreWeekdays() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTimePoint_Before(t *testing.T) {
	type fields struct {
		Hour   int
		Minute int
	}
	type args struct {
		t2 models.TimePoint
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{"tc1", fields{12, 30}, args{models.TimePoint{13, 30}}, true},
		{"tc2", fields{13, 30}, args{models.TimePoint{13, 30}}, false},
		{"tc3", fields{12, 29}, args{models.TimePoint{13, 30}}, true},
		{"tc4", fields{23, 1}, args{models.TimePoint{1, 23}}, false},
		{"tc5", fields{0, 0}, args{models.TimePoint{0, 0}}, false},
		{"tc6", fields{0, 0}, args{models.TimePoint{0, 1}}, true},
		{"tc7", fields{0, 59}, args{models.TimePoint{0, 1}}, false},
		{"tc8", fields{23, 0}, args{models.TimePoint{23, 0}}, false},
		{"tc9", fields{23, 0}, args{models.TimePoint{23, 1}}, true},
		{"1c0", fields{10, 10}, args{models.TimePoint{9, 11}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t1 := &models.TimePoint{
				Hour:   tt.fields.Hour,
				Minute: tt.fields.Minute,
			}
			if got := t1.Before(tt.args.t2); got != tt.want {
				t.Errorf("Before() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTimePoint_IsValid(t1 *testing.T) {
	type fields struct {
		Hour   int
		Minute int
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{"", fields{12, 30}, true},
		{"", fields{23, 59}, true},
		{"", fields{0, 0}, true},
		{"", fields{12, 60}, false},
		{"", fields{-1, 50}, false},
		{"", fields{24, 00}, false},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &models.TimePoint{
				Hour:   tt.fields.Hour,
				Minute: tt.fields.Minute,
			}
			if got := t.IsValid(); got != tt.want {
				t1.Errorf("IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTimePoint_ToStr(t1 *testing.T) {
	type fields struct {
		Hour   int
		Minute int
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{"", fields{12, 30}, "12:30"},
		{"", fields{23, 1}, "23:01"},
		{"", fields{0, 0}, "00:00"},
		{"", fields{13, 31}, "13:31"},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &models.TimePoint{
				Hour:   tt.fields.Hour,
				Minute: tt.fields.Minute,
			}
			if got := t.ToStr(); got != tt.want {
				t1.Errorf("ToStr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWeekday_ToStr(t *testing.T) {
	tests := []struct {
		name string
		w    models.Weekday
		want string
	}{
		{"m", models.Weekday(1), "Monday"},
		{"tu", models.Weekday(2), "Tuesday"},
		{"w", models.Weekday(3), "Wednesday"},
		{"th", models.Weekday(4), "Thursday"},
		{"f", models.Weekday(5), "Friday"},
		{"sa", models.Weekday(6), "Saturday"},
		{"su", models.Weekday(7), "Sunday"},
		{"u1", models.Weekday(0), "ErrorWeekday"},
		{"u2", models.Weekday(8), "ErrorWeekday"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.w.ToStr(); got != tt.want {
				t.Errorf("ToStr() = %v, want %v", got, tt.want)
			}
		})
	}
}
