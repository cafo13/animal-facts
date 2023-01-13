package types

type Fact struct {
	Id       string `json:"id"`
	Text     string `json:"text"`
	Category string `json:"category"`
	Source   string `json:"source"`
	Approved bool   `json:"approved"`
}
