package models

type CheckingReq struct {
	Text string `json:"text"`
}

type CheckingResp struct {
	OriginalText string   `json:"original_text"`
	Result       string   `json:"result"`
	Sources      []string `json:"sources"`
	Reason       string   `json:"reason"`
}
