package models

type CheckingResult struct {
	OriginalText string   `json:"original_text"`
	Result       string   `json:"result"`
	Sources      []string `json:"sources"`
	Reason       string   `json:"reason"`
}
