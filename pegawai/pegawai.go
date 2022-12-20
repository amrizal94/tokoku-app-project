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
	nama     string
	isActive int8
}

type AuthPegawai struct {
	db *sql.DB
}

type PegawaiInterface interface {
	Login(username, password string) (Pegawai, error)
	Register(newPegawai Pegawai) (bool, int, error)
	Duplicate(username string) (int, int8)
	Update(newPassword, newName string, isActive, id int) (bool, error)
}

func (p *Pegawai) SetUsername(newUsername string) {
	p.username = newUsername
}
func (p *Pegawai) SetPassword(newPassword string) {
	p.password = newPassword
}
func (p *Pegawai) SetNama(newNama string) {
	p.nama = newNama
}
func (p *Pegawai) SetIsActive(newIsActive int8) {
	p.isActive = newIsActive
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
func (p *Pegawai) GetNama() string {
	return p.nama
}
func (p *Pegawai) GetIsActive() int8 {
	return p.isActive
}

func NewPegawaiMenu(conn *sql.DB) PegawaiInterface {
	return &AuthPegawai{
		db: conn,
	}
}

func (ap *AuthPegawai) Login(username, password string) (Pegawai, error) {
	loginQry, err := ap.db.Prepare(`
	SELECT id, nama
	FROM pegawai
	WHERE username = ? and password = ? and isActive = 1;`)
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
	if err := row.Scan(&res.id, &res.nama); err != nil {
		log.Println("after login query", err.Error())
		return Pegawai{}, errors.New("username atau password salah")
	}
	res.username = username
	return res, nil
}

func (ap *AuthPegawai) Duplicate(username string) (int, int8) {
	res := ap.db.QueryRow("SELECT id, isActive FROM pegawai where username = ?", username)
	var tmp Pegawai
	if err := res.Scan(&tmp.id, &tmp.isActive); err != nil {
		if err.Error() != "sql: no rows in result set" {
			log.Println("Result scan error", err.Error())
			return tmp.id, 1
		}
	}
	if tmp.id > 0 {
		return tmp.id, tmp.isActive
	}
	return tmp.id, tmp.isActive
}

func (ap *AuthPegawai) Register(newPegawai Pegawai) (bool, int, error) {
	registerQry, err := ap.db.Prepare("INSERT INTO pegawai (username, password, nama, isActive) values (?,?,?,1)")
	if err != nil {
		log.Println("prepare insert pegawai registerQry", err.Error())
		return false, 0, errors.New("prepare statement insert pegawai error registerQry")
	}
	isDuplicate, isActive := ap.Duplicate(newPegawai.username)
	if isDuplicate > 0 {
		if isActive > 0 {
			log.Println("duplicated information registerQry")
			return true, 0, errors.New("username sudah digunakan registerQry")
		} else {
			return false, isDuplicate, nil
		}

	}

	// menjalankan query dengan parameter tertentu
	res, err := registerQry.Exec(newPegawai.GetUsername(), newPegawai.GetPassword(), newPegawai.GetNama())
	if err != nil {
		log.Println("insert pegawai registerQry ", err.Error())
		return false, 0, errors.New("insert pegawai error registerQry")
	}
	// Cek berapa baris yang terpengaruh query diatas
	affRows, err := res.RowsAffected()
	if err != nil {
		log.Println("after insert username registerQry ", err.Error())
		return false, 0, errors.New("error setelah insert registerQry")
	}

	if affRows <= 0 {
		log.Println("no record affected registerQry")
		return true, 0, errors.New("no record registerQry")
	}

	return true, 0, nil
}
func (ap *AuthPegawai) Update(newPassword, newName string, isActive, id int) (bool, error) {
	updateQry, err := ap.db.Prepare(`
	UPDATE pegawai
	SET password = ?, nama = ?, isActive = 1
	WHERE id = ?;`)
	if err != nil {
		if isActive == 0 {
			log.Println("prepare insert pegawai updateQry", err.Error())
			return false, errors.New("prepare statement insert pegawai error updateQry")
		} else {
			log.Println("prepare change password pegawai ", err.Error())
			return false, errors.New("prepare statement change password pegawai error updateQry")
		}
	}

	res, err := updateQry.Exec(newPassword, newName, id)
	if err != nil {
		if isActive == 0 {
			log.Println("insert pegawai updateQry", err.Error())
			return false, errors.New("insert pegawai error")
		} else {
			log.Println("update password ", err.Error())
			return false, errors.New("update password error")
		}
	}
	affRow, err := res.RowsAffected()
	if err != nil {
		if isActive == 0 {
			log.Println("after insert pegawai ", err.Error())
			return false, errors.New("after insert pegawai error")
		} else {
			log.Println("after update password ", err.Error())
			return false, errors.New("error setelah update password")
		}
	}
	if affRow <= 0 {
		log.Println("no record affected")
		return false, errors.New("no record")
	}
	return true, nil
}
