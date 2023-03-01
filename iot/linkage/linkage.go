package linkage

type Action struct {
	Addition []string `json:"addition"`
	Type     string   `json:"type"`
	//一下是作为和类型一起的出现的
	DeviceCommand DeviceCommand `json:"device_command,omitempty"`
	SmnForwarding SmnForwarding `json:"smn_forwarding,omitempty"`
	DeviceAlarm   DeviceAlarm   `json:"device_alarm,omitempty"`
}
type DeviceCommand struct {
	DeviceId string `json:"device_id"`
	Cmd      Cmd    `json:"cmd"`
}
type SmnForwarding struct {
	MessageContent string `json:"message_content"`
	MessageTitle   string `json:"message_title"`
	ProjectId      string `json:"project_id"`
	RegionName     string `json:"region_name"`
	ThemeName      string `json:"theme_name"`
	TopicUrn       string `json:"topic_urn"`
}
type DeviceAlarm struct {
	AlarmStatus string `json:"alarm_status"`
	Description string `json:"description,omitempty"`
	Name        string `json:"name"`
	Severity    string `json:"severity"`
}
type Cmd struct {
	BufferTimeout   int64    `json:"buffer_timeout"`
	CommandBody     []string `json:"command_body"`
	CommandName     string   `json:"command_name"`
	ResponseTimeout int64    `json:"response_timeout"`
	ServiceId       string   `json:"service_id"`
}
type ConditionGroup struct {
	Logic      string      `json:"logic"`
	Conditions []Condition `json:"conditions"`
	TimeRange  TimeRange   `json:"time_range"`
}
type Condition struct {
	Type string `json:"type"`
	//一下作为数组值  匹配一个类型加入到Conditions数组中
	DailyTimerCondition     DailyTimerCondition     `json:"daily_timer_condition,omitempty"`
	SimpleTimerCondition    SimpleTimerCondition    `json:"simple_timer_condition,omitempty"`
	DevicePropertyCondition DevicePropertyCondition `json:"device_property_condition,omitempty"`
}
type TimeRange struct {
	DaysOfWeek string `json:"days_of_week"`
	EndTime    string `json:"end_time"`
	StartTime  string `json:"start_time"`
}
type DailyTimerCondition struct {
	DaysOfWeek string `json:"days_of_week"`
	Time       string `json:"time"`
}
type SimpleTimerCondition struct {
	RepeatCount    int   `json:"repeat_count"`
	RepeatInterval int   `json:"repeat_interval"`
	StartTime      int64 `json:"start_time"`
}
type DevicePropertyCondition struct {
	ProductId string   `json:"product_id"`
	Filters   []Filter `json:"filters"`
}
type Filter struct {
	Operator string   `json:"operator"`
	Path     string   `json:"path"`
	Value    string   `json:"value"`
	Strategy Strategy `json:"strategy"`
}
type Strategy struct {
	EventValidTime int    `json:"event_valid_time"`
	Trigger        string `json:"trigger"`
}
