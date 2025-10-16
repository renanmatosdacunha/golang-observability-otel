package api

type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	CPF      string `json:"cpf"`
	FullName string `json:"full_name" `
}
