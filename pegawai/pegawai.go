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
