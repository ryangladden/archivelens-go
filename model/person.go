package model

type Person struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Metadata string `json:"metadata"`
}
