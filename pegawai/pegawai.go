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
	Register(newPegawai Pegawai) (bool, error)
	Duplicate(username string) (int, int8)
	Update(newPegawai Pegawai) (bool, error)
	Select(id, id_logged int) ([]Pegawai, error)
	Delete(id, isActive int) (bool, error)
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

func (ap *AuthPegawai) Register(newPegawai Pegawai) (bool, error) {
	registerQry, err := ap.db.Prepare("INSERT INTO pegawai (username, password, nama, isActive) values (?,?,?,1)")
	if err != nil {
		log.Println("prepare insert pegawai registerQry", err.Error())
		return false, errors.New("prepare statement insert pegawai error registerQry")
	}
	isDuplicate, isActive := ap.Duplicate(newPegawai.username)
	if isDuplicate > 0 {
		if isActive > 0 {
			log.Println("duplicated information registerQry")
			return false, errors.New("username sudah digunakan registerQry")
		} else {
			newPegawai.id = isDuplicate
			res, err := ap.Update(newPegawai)

			return res, err
		}

	}

	// menjalankan query dengan parameter tertentu
	res, err := registerQry.Exec(newPegawai.username, newPegawai.password, newPegawai.nama)
	if err != nil {
		log.Println("insert pegawai registerQry ", err.Error())
		return false, errors.New("insert pegawai error registerQry")
	}
	// Cek berapa baris yang terpengaruh query diatas
	affRows, err := res.RowsAffected()
	if err != nil {
		log.Println("after insert username registerQry ", err.Error())
		return false, errors.New("error setelah insert registerQry")
	}

	if affRows <= 0 {
		log.Println("no record affected registerQry")
		return true, errors.New("no record registerQry")
	}

	return true, nil
}
func (ap *AuthPegawai) Update(newPegawai Pegawai) (bool, error) {

	resSelect, err := ap.Select(newPegawai.id, 0)
	if err != nil {
		log.Println("res Select")
		return false, errors.New("data pegawai tidak ada")
	}
	if newPegawai.password == "" {
		newPegawai.password = resSelect[0].password
	}
	if newPegawai.nama == "" {
		newPegawai.nama = resSelect[0].nama
	}

	updateQry, err := ap.db.Prepare(`
	UPDATE pegawai
	SET password = ?, nama = ?, isActive = ?
	WHERE id = ?;`)
	if err != nil {
		if newPegawai.isActive == 0 {
			log.Println("prepare insert pegawai updateQry", err.Error())
			return false, errors.New("prepare statement insert pegawai error updateQry")
		} else {
			log.Println("prepare change password pegawai ", err.Error())
			return false, errors.New("prepare statement change password pegawai error updateQry")
		}
	}

	res, err := updateQry.Exec(newPegawai.password, newPegawai.nama, newPegawai.isActive, newPegawai.id)
	if err != nil {
		if newPegawai.isActive == 0 {
			log.Println("insert pegawai updateQry", err.Error())
			return false, errors.New("insert pegawai error")
		} else {
			log.Println("update password ", err.Error())
			return false, errors.New("update password error")
		}
	}
	affRow, err := res.RowsAffected()
	if err != nil {
		if newPegawai.isActive == 0 {
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

func (ap *AuthPegawai) Select(id, id_logged int) ([]Pegawai, error) {
	var (
		selectPegawaiQry *sql.Rows
		err              error
		cases            int8
	)
	if id == 0 && id_logged != 0 {
		cases = 1
		selectPegawaiQry, err = ap.db.Query(`
		SELECT id, nama
		FROM pegawai
		WHERE id != ?
		AND isActive = 1;`, id_logged)
	} else if id != 0 && id_logged == 0 {
		cases = 2
		selectPegawaiQry, err = ap.db.Query(`
		SELECT username, password, nama
		FROM pegawai
		WHERE id = ?;`, id)
	}
	if err != nil {
		log.Println("select pegawai", err.Error())
		return nil, errors.New("select pegawai error")
	}
	arrPegawai := []Pegawai{}
	for selectPegawaiQry.Next() {
		var tmp Pegawai
		switch cases {
		case 1:
			err = selectPegawaiQry.Scan(&tmp.id, &tmp.nama)
		case 2:
			err = selectPegawaiQry.Scan(&tmp.username, &tmp.password, &tmp.nama)
		}

		if err != nil {
			log.Println("Loop through rows, using Scan to assign column data to struct fields", err.Error())
			return arrPegawai, err
		}
		arrPegawai = append(arrPegawai, tmp)
	}
	return arrPegawai, nil
}

func (ap *AuthPegawai) Delete(id, isActive int) (bool, error) {
	newPegawai := Pegawai{}
	newPegawai.id = id
	newPegawai.isActive = int8(isActive)
	resupdate, err := ap.Update(newPegawai)
	return resupdate, err
}
