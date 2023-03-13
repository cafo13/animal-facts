package models

type LoginInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type MessageResponse struct {
	Message string
}

type ErrorResponse struct {
	Error error
}

type Fact struct {
	Text     string `json:"Text,omitempty"`
	Category string `json:"Category,omitempty"`
	Source   string `json:"Source,omitempty"`
}

type FactWithID struct {
	ID uint `json:"ID"`
	Fact
}
