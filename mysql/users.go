package mysql

import (
	"database/sql"
	"snippetbox/models"

	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) Insert(name, email, password string) error {
	//creates a new record in users table containing
	//the validated name,email and hashed password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}
	stmt := `insert into users(name,email,hashed_password,created) values(?,?,?,now())`
	_, err = m.DB.Exec(stmt, name, email, string(hashedPassword))
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok {
			if mysqlErr.Number == 1062 {
				return models.ErrDuplicateEmail
			}
		}
	}

	return err
}

func (m *UserModel) Authenticate(email, password string) (int, error) {
	//we retrieve the hashed password associated with the email addr
	//from the users table,if email does not exist we return
	//ErrInvalidCredentials error

	//if email does exist, we want to compare the bcrypt-hashed
	//password to the plain-text password.
	//if they match we return the user's id
	//else ErrInvalidCredentials error
	var id int
	var hashedPassword []byte
	row := m.DB.QueryRow("Select id,hashed_password from users where email=?", email)
	err := row.Scan(&id, &hashedPassword)

	if err == sql.ErrNoRows {
		return 0, models.ErrInvalidCredentials
	} else if err != nil {
		return 0, err
	}
	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return 0, models.ErrInvalidCredentials
	} else if err != nil {
		return 0, err
	}
	return id, nil
}
func (m *UserModel) Get(id int) (*models.User, error) {
	s:=&models.User{}
	stmt:=`select id, name,email,created from users where id=?`
	err:=m.DB.QueryRow(stmt,id).Scan(&s.ID,&s.Name,&s.Email,&s.Created)
	if err == sql.ErrNoRows{
		return nil,models.ErrNoRecord
	}else if err!=nil{
		return nil,err
	}
	return s,nil
}
