package entity

type User struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Password []byte `json:"-"`
	Token    string `json:"-"`
}
