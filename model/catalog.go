package model

//Catalog ...
type Catalog struct {
	BasicAuthPass string `json:"basic_auth_pass"`
	BasicAuthUser string `json:"basic_auth_user"`
	PlanSettings  []struct {
		CPULimit      string `json:"cpu_limit"`
		CPURequest    string `json:"cpu_request"`
		ID            string `json:"id"`
		ImageName     string `json:"image_name"`
		MemoryLimit   string `json:"memory_limit"`
		MemoryRequest string `json:"memory_request"`
		Storage       string `json:"storage"`
	} `json:"plan_settings"`
	Services []struct {
		Bindable    bool   `json:"bindable"`
		Description string `json:"description"`
		ID          string `json:"id"`
		Metadata    struct {
			Listing struct {
				Blurb    string      `json:"blurb"`
				ImageURL interface{} `json:"imageUrl"`
			} `json:"listing"`
			Provider struct {
				Name interface{} `json:"name"`
			} `json:"provider"`
		} `json:"metadata"`
		Name  string `json:"name"`
		Plans []struct {
			Description string `json:"description"`
			ID          string `json:"id"`
			Metadata    struct {
				Bullets     []string      `json:"bullets"`
				Costs       []interface{} `json:"costs"`
				DisplayName string        `json:"displayName"`
			} `json:"metadata"`
			Name string `json:"name"`
		} `json:"plans"`
		Tags     []string `json:"tags"`
		Requires []string `json:"requires"`
	} `json:"services"`
}
