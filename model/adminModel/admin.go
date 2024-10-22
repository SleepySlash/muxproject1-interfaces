package adminModel

type Admin struct {
	AdminID  string `json:"admin_id,omitempty"`
	Name     string `json:"name,omitempty"`
	Password string `json:"password,omitempty"`
}
