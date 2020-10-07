package model

//SvcCatalog ...
type CatalogResponse struct {
	Services []Service `json:"services"`
}

//LastOperationStatus type for responding to LastOperation operation
type LastOperationStatus struct {
	State                    string `json:"state"`
	Description              string `json:"description"`
	AsyncPollIntervalSeconds int    `json:"async_poll_interval_seconds, omitempty"`
}

// type CreateServiceInstanceResponse struct {
// 	DashboardUrl  string               `json:"dashboard_url"`
// 	LastOperation *LastOperationStatus `json:"last_operation, omitempty"`
// }

type ProvisioningResponse struct {
	DashboardURL  string `json:"dashboard_url,omitempty"`
	OperationData string `json:"operation,omitempty"`
}
