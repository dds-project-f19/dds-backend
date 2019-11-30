package tests

import (
	"dds-backend/models"
	"dds-backend/services"
	"testing"
)

func TestGetCronString(t *testing.T) {
	type args struct {
		time models.TimePoint
		days []models.Weekday
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"Basic", args{models.TimePoint{10, 20}, []models.Weekday{models.Monday, models.Tuesday}}, "20 10 * * 1,2", false},
		{"None", args{models.TimePoint{10, 20}, []models.Weekday{}}, "", true},
		{"Single", args{models.TimePoint{10, 20}, []models.Weekday{models.Wednesday}}, "20 10 * * 3", false},
		{"Order", args{models.TimePoint{10, 20}, []models.Weekday{models.Sunday, models.Thursday, models.Tuesday}}, "20 10 * * 2,4,7", false},
		{"Duplicate", args{models.TimePoint{10, 20}, []models.Weekday{models.Friday, models.Friday}}, "20 10 * * 5", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := services.GetCronString(tt.args.time, tt.args.days)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetCronString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetCronString() got = %v, want %v", got, tt.want)
			}
		})
	}
}
