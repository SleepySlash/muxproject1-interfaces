package adminModel

import (
	"time"

	"github.com/google/uuid"
)

type Admin struct {
	AdminID   string `json:"admin_id,omitempty"`
	Name      string `json:"name,omitempty"`
	Password  string `json:"password,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}

func NewAdmin(name, password string) *Admin {
	return &Admin{
		AdminID: uuid.NewString(),
		Name:name,
		Password: password,
		CreatedAt: time.Now().UTC(),
	}
}