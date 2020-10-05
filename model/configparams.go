package model

import ()

type ServiceInstance struct {
	Instance_id        string
	Memory_request     string
	Cpu_request        string
	Image_name_and_tag string
	Memory_limit       string
	Cpu_limit          string
	Storage_capacity   string
	Pg_password        string
	Storage_class      string
}
