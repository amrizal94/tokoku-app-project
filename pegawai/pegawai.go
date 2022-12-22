package pegawai

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
)

type Pegawai struct {
	id       int
	username string
	password string
	nama     string
	isActive int8
}

type PegawaiMenu struct {
	db *sql.DB
}

// Buat manggil ke main
type PegawaiInterface interface {
	Login(username, password string) (Pegawai, error)
	Register(newPegawai Pegawai) (bool, error)
	Update(newPegawai Pegawai) (bool, error)
	Data(username string) ([]Pegawai, string, error)
	Delete(id_pegawai int) (bool, error)
}

func NewPegawaiMenu(conn *sql.DB) PegawaiInterface {
	return &PegawaiMenu{
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

// ////menu login
func (pm *PegawaiMenu) Login(username, password string) (Pegawai, error) {
	var (
		row *sql.Row
		res Pegawai
	)

	usernameQry, err := pm.db.Prepare(`
	SELECT id, nama
	FROM pegawai
	WHERE username = ? 
	AND isActive = 1;`)
	if err != nil {
		log.Println("prepare usernameQry pegawai ", err.Error())
		return Pegawai{}, errors.New("prepare statement usernameQry pegawai error")
	}
	row = usernameQry.QueryRow(username)
	if row.Err() != nil {
		log.Println("username query ", row.Err().Error())
		return Pegawai{}, errors.New("select username pegawai error")
	}

	if err := row.Scan(&res.id, &res.nama); err != nil {
		if err.Error() != "sql: no rows in result set" {
			log.Println("after login query")
			return Pegawai{}, err
		}
		return Pegawai{}, errors.New("username belum terdaftar")
	}

	loginQry, err := pm.db.Prepare(`
	SELECT id, nama
	FROM pegawai
	WHERE username = ? 
	AND password = ? 
	AND isActive = 1;`)
	if err != nil {
		log.Println("prepare loginQry pegawai ", err.Error())
		return Pegawai{}, errors.New("prepare statement loginQry pegawai error")
	}
	row = loginQry.QueryRow(username, password)
	if row.Err() != nil {
		log.Println("login query ", row.Err().Error())
		return Pegawai{}, errors.New("select pegawai error")
	}
	if err := row.Scan(&res.id, &res.nama); err != nil {
		if err.Error() != "sql: no rows in result set" {
			log.Println("after login query")
			return Pegawai{}, err
		}
		return Pegawai{}, errors.New("password salah")
	}
	res.username = username
	return res, nil
}

func (pm *PegawaiMenu) Duplicate(username string) (bool, error) {
	res := pm.db.QueryRow(`
	SELECT isActive
	FROM pegawai
	WHERE username = ?
	`, username)
	var isActive int8
	if err := res.Scan(&isActive); err != nil {
		if err.Error() != "sql: no rows in result set" {
			log.Println("Result scan error", err.Error())
			return false, err
		}
	}
	if isActive == 1 {
		return true, nil
	}
	return false, nil
}

// ////menu register
func (pm *PegawaiMenu) Register(newPegawai Pegawai) (bool, error) {
	registerQry, err := pm.db.Prepare("INSERT INTO pegawai (username, password, nama, isActive) values (?,?,?,1)")
	if err != nil {
		log.Println("prepare insert pegawai registerQry", err.Error())
		return false, errors.New("prepare statement insert pegawai error registerQry")
	}
	isDuplicate, err := pm.Duplicate(newPegawai.username)
	if err != nil {
		return false, err
	}
	if isDuplicate {
		log.Println("duplicated information")
		return false, err
	}

	// menjalankan query dengan parameter tertentu
	res, err := registerQry.Exec(newPegawai.username, newPegawai.password, newPegawai.nama)
	if err != nil {
		log.Println("register pegawai ", err.Error())
		return false, errors.New("register pegawai error")
	}
	// Cek berapa baris yang terpengaruh query diatas
	affRows, err := res.RowsAffected()
	if err != nil {
		log.Println("after register pegawai ", err.Error())
		return false, errors.New("error setelah register pegawai")
	}

	if affRows <= 0 {
		log.Println("no record affected")
		return false, errors.New("no record")
	}

	return true, nil
}
func (pm *PegawaiMenu) Update(newPegawai Pegawai) (bool, error) {
	updateQry, err := pm.db.Prepare(`
	UPDATE pegawai
	SET password = ?, nama = ?
	WHERE username = ?;`)
	if err != nil {
		log.Println("prepare updateQry pegawai", err.Error())
		return false, errors.New("prepare statement updateQry pegawai error")
	}

	res, err := updateQry.Exec(newPegawai.password, newPegawai.nama, newPegawai.username)
	if err != nil {
		log.Println("update pegawai ", err.Error())
		return false, errors.New("update pegawai error")
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

func (pm *PegawaiMenu) Data(username string) ([]Pegawai, string, error) {
	var (
		selectPegawaiQry *sql.Rows
		err              error
		cases            int8
		strPegawai       string
	)
	fmt.Println(username)
	if username == "admin" {
		// case 1 untuk memanggil list pegawai yg aktif kecuali admin
		cases = 1
		selectPegawaiQry, err = pm.db.Query(`
		SELECT id, nama
		FROM pegawai
		WHERE username != ?
		AND isActive = 1;`, username)
	} else {
		// case 2 untuk memanggil data pegawai menurut idnya
		cases = 2
		selectPegawaiQry, err = pm.db.Query(`
		SELECT username, password, nama
		FROM pegawai
		WHERE id = ?;`, username)
	}

	if err != nil {
		log.Println("select pegawai", err.Error())
		return nil, strPegawai, errors.New("select pegawai error")
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
			return arrPegawai, strPegawai, err
		}
		strPegawai += fmt.Sprintf("%d\t| %s\n", tmp.id, tmp.nama)
		arrPegawai = append(arrPegawai, tmp)
	}
	return arrPegawai, strPegawai, nil
}

func (pm *PegawaiMenu) Delete(id_pegawai int) (bool, error) {
	isDeleted, err := pm.ChangeIsActive(id_pegawai, 0)
	return isDeleted, err
}

func (pm *PegawaiMenu) ChangeIsActive(id_pegawai int, isActive int8) (bool, error) {
	updateQry, err := pm.db.Prepare(`
	UPDATE pegawai
	SET isActive = ?
	WHERE id = ?;`)
	if err != nil {
		log.Println("prepare change isActive pegawai", err.Error())
		return false, errors.New("prepare statement change isActive error pegawai")
	}
	res, err := updateQry.Exec(isActive, id_pegawai)
	if err != nil {
		log.Println("update isActive pegawai", err.Error())
		return false, errors.New("update isActive pegawai error")
	}
	affRow, err := res.RowsAffected()
	if err != nil {
		log.Println("after update isActive pegawai ", err.Error())
		return false, errors.New("error setelah update isActive pegawai")
	}
	if affRow <= 0 {
		log.Println("no record affected")
		return false, errors.New("no record")
	}
	return true, nil
}
