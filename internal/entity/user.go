package entity

type User struct {
	Name     string `json:"name"`
	Password string `json:"-"`
	Token    string `json:"-"`
}
