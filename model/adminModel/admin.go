package adminModel

type Admin struct {
	AdminID  string `json:"id"`
	Name     string `json:"name"`
	Password string `json:"password"`
}
