package request

const (
	REQUEST_ID_CONFIG          = 10010
	REQUEST_ID_NODE_REGISTERED = 10020
	REQUEST_ID_AUTH_ERROR      = 10030
	REQUEST_ID_KEEPALIVE       = 10040
)

type MasterResponse struct {
	Id   int    `json:"id"`
	Data string `json:"data"`
}

type NodeRequest struct {
	Id     int    `json:"id"`
	Data   string `json:"data"`
	NodeId string `json:"nodeId"`
	Token  string `json:"token"`
}
