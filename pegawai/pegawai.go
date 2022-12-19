package pegawai

import (
	"database/sql"
	"errors"
	"log"
)

type Pegawai struct {
	id       int
	username string
	password string
}

type AuthPegawai struct {
	db *sql.DB
}

type PegawaiInterface interface {
	Login(username, password string) (Pegawai, error)
	Register(newPegawai Pegawai) (bool, error)
	Duplicate(username string) bool
}

func NewPegawaiMenu(conn *sql.DB) PegawaiInterface {
	return &AuthPegawai{
		db: conn,
	}
}

func (p *Pegawai) SetUsername(newUsername string) {
	p.username = newUsername
}
func (p *Pegawai) SetPassword(newPassword string) {
	p.password = newPassword
}

func (p *Pegawai) GetID() int {
	return p.id
}
func (p *Pegawai) GetUsername() string {
	return p.username
}
func (p *Pegawai) GetPassword() string {
	return p.password
}

func (ap *AuthPegawai) Login(username, password string) (Pegawai, error) {
	loginQry, err := ap.db.Prepare(`
	SELECT id
	FROM pegawai
	WHERE username = ? and password = ?;`)
	if err != nil {
		log.Println("prepare login pegawai ", err.Error())
		return Pegawai{}, errors.New("prepare statement login pegawai error")
	}
	row := loginQry.QueryRow(username, password)
	if row.Err() != nil {
		log.Println("login query ", row.Err().Error())
		return Pegawai{}, errors.New("select pegawai error")
	}
	res := Pegawai{}
	err = row.Scan(&res.id)
	if err != nil {
		log.Println("after login query", err.Error())
		return Pegawai{}, errors.New("username atau password salah")
	}
	res.username = username
	return res, nil
}

func (ap *AuthPegawai) Duplicate(username string) bool {
	res := ap.db.QueryRow("SELECT id FROM pegawai where username = ?", username)
	var idExist int
	err := res.Scan(&idExist)
	if err != nil {
		if err.Error() != "sql: no rows in result set" {
			log.Println("Result scan error", err.Error())
			return true
		}
	}
	if idExist > 0 {
		return true
	}
	return false
}

func (ap *AuthPegawai) Register(newPegawai Pegawai) (bool, error) {
	registerQry, err := ap.db.Prepare("INSERT INTO pegawai (username, password) values (?,?)")
	if err != nil {
		log.Println("prepare insert pegawai ", err.Error())
		return false, errors.New("prepare statement insert pegawai error")
	}

	if ap.Duplicate(newPegawai.GetUsername()) {
		log.Println("duplicated information")
		return false, errors.New("username sudah digunakan")
	}

	// menjalankan query dengan parameter tertentu
	res, err := registerQry.Exec(newPegawai.GetUsername(), newPegawai.GetPassword())
	if err != nil {
		log.Println("insert user ", err.Error())
		return false, errors.New("insert username error")
	}
	// Cek berapa baris yang terpengaruh query diatas
	affRows, err := res.RowsAffected()
	if err != nil {
		log.Println("after insert username ", err.Error())
		return false, errors.New("error setelah insert")
	}

	if affRows <= 0 {
		log.Println("no record affected")
		return true, errors.New("no record")
	}

	return true, nil
}
