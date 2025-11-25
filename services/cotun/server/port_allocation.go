package chserver

import "time"

type PortStatus string

const (
	Freed     PortStatus = "freed"
	Allocated PortStatus = "allocated"
	Connected PortStatus = "connected"
)

// PortAllocation 统一的端口分配数据结构
type PortAllocation struct {
	ClientId      string     `json:"clientId"`            //客户端机器ID
	UserId        string     `json:"userId"`              //用户ID
	AppName       string     `json:"appName"`             //应用名称
	ClientVersion string     `json:"clientVersion"`       //客户端cotun版本
	ClientPort    int        `json:"clientPort"`          //应用端口
	MappingPort   int        `json:"mappingPort"`         //映射端口
	Status        PortStatus `json:"status"`              //状态
	AllocTime     *time.Time `json:"allocTime,omitempty"` //端口分配时间
	StartTime     *time.Time `json:"startTime,omitempty"` //隧道建立时间
}

// PortAllocationRequest 端口分配请求
type PortAllocationRequest struct {
	ClientId   string `json:"clientId"`
	UserId     string `json:"userId,omitempty"`
	AppName    string `json:"appName"`
	ClientPort int    `json:"clientPort,omitempty"`
}

// PortAllocationResponse 端口分配响应
type PortAllocationResponse struct {
	ClientId    string `json:"clientId"`
	UserId      string `json:"userId"`
	AppName     string `json:"appName"`
	ClientPort  int    `json:"clientPort"`
	MappingPort int    `json:"mappingPort"`
}
