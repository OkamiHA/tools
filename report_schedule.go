package schedulerreport

import "time"

type ReportScheduleTimeConfig struct {
	DayOfMonth int `bson:"dayOfMonth,omitempty" json:"dayOfMonth,omitempty" mapstructure:"day_of_month"`
	DayOfWeek  int `bson:"dayOfWeek,omitempty" json:"dayOfWeek,omitempty" mapstructure:"day_of_week"`
	Hour       int `bson:"hour,omitempty" json:"hour,omitempty" mapstructure:"hour"`
	Minute     int `bson:"minute,omitempty" json:"minute,omitempty" mapstructure:"minute"`
	Second     int `bson:"second,omitempty" json:"second,omitempty" mapstructure:"second"`
}

type ScheduleReportConfig struct {
	Name           string                       `bson:"name,omitempty" json:"name,omitempty" mapstructure:"name"`
	Frequency      Frequency                    `bson:"frequency,omitempty" json:"frequency,omitempty" mapstructure:"frequency"`
	CustomerConfig ScheduleReportCustomerConfig `bson:"customerConfig,omitempty" json:"customerConfig,omitempty" mapstructure:"customer_config"`
}

type Frequency string

const (
	FrequencyDaily   Frequency = "daily"
	FrequencyWeekly  Frequency = "weekly"
	FrequencyMonthly Frequency = "monthly"
)

type ScheduleReportCustomerConfig struct {
	StartTimeConfig     ReportScheduleTimeConfig `bson:"startTimeConfig,omitempty" json:"startTimeConfig,omitempty" mapstructure:"start_time_config"`
	EndTimeConfig       ReportScheduleTimeConfig `bson:"endTimeConfig,omitempty" json:"endTimeConfig,omitempty" mapstructure:"end_time_config"`
	IsPreviousFrequency bool                     `bson:"isPreviousFrequency,omitempty" json:"isPreviousFrequency,omitempty" mapstructure:"is_previous_frequency"`
	SendEmailTimeConfig ReportScheduleTimeConfig `bson:"sendEmailTimeConfig,omitempty" json:"sendEmailTimeConfig,omitempty" mapstructure:"send_email_time_config"`
	DeadlineTimeConfig  ReportScheduleTimeConfig `bson:"deadlineTimeConfig,omitempty" json:"deadlineTimeConfig,omitempty" mapstructure:"deadline_time_config"`
	StartTime           int64                    `bson:"startTime,omitempty" json:"startTime,omitempty" mapstructure:"start_time"`
	EndTime             int64                    `bson:"endTime,omitempty" json:"endTime,omitempty" mapstructure:"end_time"`
	ExportReportTime    int64                    `bson:"exportReportTime,omitempty" json:"exportReportTime,omitempty" mapstructure:"export_report_time"`
	SendReportTime      int64                    `bson:"sendReportTime,omitempty" json:"sendReportTime,omitempty" mapstructure:"send_report_time"`
	DeadlineTime        int64                    `bson:"deadlineTime,omitempty" json:"deadlineTime,omitempty" mapstructure:"deadline_time"`
}

func (t ReportScheduleTimeConfig) TotalSeconds() int64 {
	return int64(t.DayOfMonth*24*3600 + t.Hour*3600 + t.Minute*60 + t.Second)
}

func (t ReportScheduleTimeConfig) DecorateTimeConfig() ReportScheduleTimeConfig {
	if t.Second < 0 {
		t.Second += 60
		t.Minute -= 1
	}
	if t.Minute < 0 {
		t.Minute += 60
		t.Hour -= 1
	}
	if t.Hour < 0 {
		t.Hour += 24
		t.DayOfMonth -= 1
		t.DayOfWeek -= 1
	}
	return t
}

func (t ReportScheduleTimeConfig) DecorateDeltaTimeConfigWithFrequency(now time.Time, delta int) ReportScheduleTimeConfig {
	if t.TotalSeconds() >= 0 {
		return t
	}
	newTime := time.Date(now.Year(), now.Month(), 0, 0, 0, 0, 0, now.Location())
	t.DayOfMonth += newTime.Day()
	return t
}

func (t ReportScheduleTimeConfig) Sub(tu ReportScheduleTimeConfig) ReportScheduleTimeConfig {
	return ReportScheduleTimeConfig{
		DayOfMonth: t.DayOfMonth - tu.DayOfMonth,
		DayOfWeek:  t.DayOfWeek - tu.DayOfWeek,
		Hour:       t.Hour - tu.Hour,
		Minute:     t.Minute - tu.Minute,
		Second:     t.Second - tu.Second,
	}
}

func NewReportScheduleTimeConfig(timeValue time.Time) ReportScheduleTimeConfig {
	return ReportScheduleTimeConfig{
		DayOfMonth: timeValue.Day(),
		DayOfWeek:  int(timeValue.Weekday()),
		Hour:       timeValue.Hour(),
		Minute:     timeValue.Minute(),
		Second:     timeValue.Second(),
	}
}

func (t ReportScheduleTimeConfig) DecorateTimeConfigByFrequency(frequency Frequency) ReportScheduleTimeConfig {
	switch frequency {
	case FrequencyDaily:
		t.DayOfMonth = 0
		t.DayOfWeek = 0
	case FrequencyWeekly:
		t.DayOfMonth = 0
	case FrequencyMonthly:
		t.DayOfWeek = 0
	default:
	}
	return t
}

func (t ReportScheduleTimeConfig) GetDeltaSeconds(tu ReportScheduleTimeConfig, now time.Time, isPrevious bool) int64 {
	var deltaTimeConfig ReportScheduleTimeConfig
	if isPrevious {
		deltaTimeConfig = t.Sub(tu).DecorateDeltaTimeConfigWithFrequency(now)
		return deltaTimeConfig.TotalSeconds()
	}
	deltaTimeConfig = tu.Sub(t).DecorateDeltaTimeConfigWithFrequency(now)
	return deltaTimeConfig.TotalSeconds()
}

func (t ReportScheduleTimeConfig) DecorateTimeConfigOfFrequencyMonthlyWithCurrentTime(tu ReportScheduleTimeConfig, now time.Time, isNextTime bool) ReportScheduleTimeConfig {
	if t.DayOfMonth < 29 || now.Day() != 1 {
		return t
	}
	var deltaSeconds int64
	switch isNextTime {
	case true:
		deltaSeconds = t.Sub(tu).TotalSeconds()
	default:
		deltaSeconds = tu.Sub(t).TotalSeconds()
	}
	deltaMonth := 0
	if deltaSeconds < 0 {
		if isNextTime {
			deltaMonth = 1
		} else {
			deltaMonth = -1
		}
	}
	newMonth := int(now.Month()) + deltaMonth
	newTime := time.Date(now.Year(), time.Month(newMonth), t.DayOfMonth, t.Hour, t.Minute, t.Second, 0, time.Local)
	if newMonth != int(newTime.Month()) {
		t.DayOfMonth = time.Date(now.Year(), time.Month(newTime.Month()), 0, 0, 0, 0, 0, time.Local).Day()
		return t
	}
	return t
}
