package AdminRepository

import (
	"database/sql"
	"fmt"
	"log"
	"muxproject1/model/adminModel"
	"os"

	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/lib/pq"
)

type AdminService interface{
    CreateAdmin(adm adminModel.Admin) (error)
    GetAdminByPassword(admin adminModel.Admin) (bool, error)
}
type adminRepoPostgreSql struct {
	DB *sql.DB
}

func NewService() AdminService{
  dataSource := os.Getenv("POSTGRES_URI")
	db, err := sql.Open("postgres", dataSource)
	if err != nil {
	  log.Fatalf("Failed to connect to database: %v", err)
	}	
	driver, err := postgres.WithInstance(db, &postgres.Config{})

	if err != nil {
	  log.Fatalf("Could not start SQL driver: %v", err)
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations", // Correctly specify the file scheme
		"postgres", driver)
	if err != nil {
	  log.Fatalf("Could not start migration: %v", err)
	}
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
	  log.Fatalf("Migration failed: %v", err)
	}
	svc:= &adminRepoPostgreSql{
		DB:db,
	}
	return svc
}
// In AdminRepository package
func (r *adminRepoPostgreSql) CreateAdmin(adm adminModel.Admin)  error {
    _, err := r.DB.Exec("INSERT INTO admin_table (name, password) VALUES ($1, $2, $3, $4)",adm.AdminID, adm.Name, adm.Password,adm.CreatedAt)
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