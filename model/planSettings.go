package model

//PlanSetting ...
type PlanSetting struct {
	CPULimit      string `json:"cpu_limit"`
	CPURequest    string `json:"cpu_request"`
	ID            string `json:"id"`
	ImageName     string `json:"image_name"`
	MemoryLimit   string `json:"memory_limit"`
	MemoryRequest string `json:"memory_request"`
	Storage       string `json:"storage"`
}
//Plans ...
type Plans struct {
	PlanSettings []PlanSetting `json:"plan_settings"`
}
