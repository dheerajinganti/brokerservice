package model

import (
	"encoding/json"
)

//ProvisionDetails: Type used to unmarshal request body for provision operation
type ProvisionDetails struct {
	ServiceID        string          `json:"service_id"`
	PlanID           string          `json:"plan_id"`
	OrganizationGUID string          `json:"organization_guid"`
	SpaceGUID        string          `json:"space_guid"`
	RawContext       json.RawMessage `json:"context,omitempty"`
	RawParameters    json.RawMessage `json:"parameters,omitempty"`
	MaintenanceInfo  MaintenanceInfo `json:"maintenance_info,omitempty"`
}
type MaintenanceInfo struct {
	Public  map[string]string `json:"public,omitempty"`
	Private string            `json:"private,omitempty"`
}
type ProvisionedServiceSpec struct {
	IsAsync       bool
	DashboardURL  string
	OperationData string
}

type GetInstanceDetailsSpec struct {
	ServiceID    string      `json:"service_id"`
	PlanID       string      `json:"plan_id"`
	DashboardURL string      `json:"dashboard_url"`
	Parameters   interface{} `json:"parameters"`
}

type UnbindSpec struct {
	IsAsync       bool
	OperationData string
}

type BindDetails struct {
	AppGUID       string          `json:"app_guid"`
	PlanID        string          `json:"plan_id"`
	ServiceID     string          `json:"service_id"`
	BindResource  *BindResource   `json:"bind_resource,omitempty"`
	RawContext    json.RawMessage `json:"context,omitempty"`
	RawParameters json.RawMessage `json:"parameters,omitempty"`
}

type BindResource struct {
	AppGuid            string `json:"app_guid,omitempty"`
	SpaceGuid          string `json:"space_guid,omitempty"`
	Route              string `json:"route,omitempty"`
	CredentialClientID string `json:"credential_client_id,omitempty"`
}

type UnbindDetails struct {
	PlanID    string `json:"plan_id"`
	ServiceID string `json:"service_id"`
}

type UpdateServiceSpec struct {
	IsAsync       bool
	DashboardURL  string
	OperationData string
}

type DeprovisionServiceSpec struct {
	IsAsync       bool
	OperationData string
}

type DeprovisionDetails struct {
	PlanID    string `json:"plan_id"`
	ServiceID string `json:"service_id"`
}

type UpdateDetails struct {
	ServiceID       string          `json:"service_id"`
	PlanID          string          `json:"plan_id"`
	RawParameters   json.RawMessage `json:"parameters,omitempty"`
	PreviousValues  PreviousValues  `json:"previous_values"`
	RawContext      json.RawMessage `json:"context,omitempty"`
	MaintenanceInfo MaintenanceInfo `json:"maintenance_info,omitempty"`
}

type PreviousValues struct {
	PlanID    string `json:"plan_id"`
	ServiceID string `json:"service_id"`
	OrgID     string `json:"organization_id"`
	SpaceID   string `json:"space_id"`
}

type PollDetails struct {
	ServiceID     string `json:"service_id"`
	PlanID        string `json:"plan_id"`
	OperationData string `json:"operation"`
}

type LastOperation struct {
	State       LastOperationState
	Description string
}

type LastOperationState string

const (
	InProgress LastOperationState = "in progress"
	Succeeded  LastOperationState = "succeeded"
	Failed     LastOperationState = "failed"
)

type Binding struct {
	IsAsync         bool          `json:"is_async"`
	OperationData   string        `json:"operation_data"`
	Credentials     interface{}   `json:"credentials"`
	SyslogDrainURL  string        `json:"syslog_drain_url"`
	RouteServiceURL string        `json:"route_service_url"`
	VolumeMounts    []VolumeMount `json:"volume_mounts"`
}
type VolumeMount struct {
	Driver       string       `json:"driver"`
	ContainerDir string       `json:"container_dir"`
	Mode         string       `json:"mode"`
	DeviceType   string       `json:"device_type"`
	Device       SharedDevice `json:"device"`
}

type SharedDevice struct {
	VolumeId    string                 `json:"volume_id"`
	MountConfig map[string]interface{} `json:"mount_config"`
}
type GetBindingSpec struct {
	Credentials     interface{}
	SyslogDrainURL  string
	RouteServiceURL string
	VolumeMounts    []VolumeMount
	Parameters      interface{}
}
