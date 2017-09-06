package model

// ResponseInfo 返回值信息
type ResponseInfo struct {
	Status int    `json:"status"`
	Info   string `json:"info"`
}
