package request

type MasterResponse struct {
	Id   string `json:"id"`
	Data string `json:"data"`
}

type NodeRequest struct {
	Id     string `json:"id"`
	Data   string `json:"data"`
	NodeId string `json:"nodeId"`
	Token  string `json:"token"`
}
