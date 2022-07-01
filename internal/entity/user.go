package entity

type User struct {
	ID       string `json:"ID"`
	Name     string `json:"name"`
	Password []byte `json:"-"`
	Token    string `json:"-"`
}
