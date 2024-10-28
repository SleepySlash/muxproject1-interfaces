package AdminRepository

import (
	"database/sql"
	"fmt"
	"log"
	"muxproject1/model/adminModel"
	"os"

	_ "github.com/lib/pq"
)

type AdminService interface{
    CreateAdmin(adm adminModel.Admin) (error)
    GetAdminByPassword(admin adminModel.Admin) (bool, error)
}
type adminRepoPostgreSql struct {
	DB *sql.DB
}

func NewService() (AdminService,error){
  	dataSource := os.Getenv("POSTGRES_URI")
	db, err := sql.Open("postgres", dataSource)
	if err != nil {
	  log.Fatalf("Failed to connect to database: %v", err)
	}	

	if err:= db.Ping(); err!=nil{
		log.Fatal(" Ping failure ",err)
		return nil,err
	}

	svc:= &adminRepoPostgreSql{
		DB:db,
	}
	return svc,nil
}
// In AdminRepository package
func (r *adminRepoPostgreSql) CreateAdmin(adm adminModel.Admin) error {
    _, err := r.DB.Exec(
        "INSERT INTO admin_table (admin_id, name, password, created_at) VALUES ($1, $2, $3, $4)",
        adm.AdminID, adm.Name, adm.Password, adm.CreatedAt, // Ensure adm.CreatedAt is populated
    )
    return err
}



func (r *adminRepoPostgreSql) GetAdminByPassword(admin adminModel.Admin) (bool, error) {
    rows, err := r.DB.Query("SELECT admin_id, name, password, created_at FROM admin_table WHERE password = ($1)",admin.Password)
   if err != nil {
     return false, err
   }
   fmt.Sprintf("%v", rows)
   return true, nil
}