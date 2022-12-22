package pelanggan

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
)

type Pelanggan struct {
	hp             string
	id_pegawai     int
	nama_pelanggan string
	nama_pegawai   string
	isActive       int8
}

type PelangganMenu struct {
	db *sql.DB
}

type PelangganInterface interface {
	Register(newPelanggan Pelanggan) (bool, error)
	Data(hp string) ([]Pelanggan, string, error)
	Delete(hp string) (bool, error)
}

func NewPelangganMenu(conn *sql.DB) PelangganInterface {
	return &PelangganMenu{
		db: conn,
	}
}

func (p *Pelanggan) SetHP(newHP string) {
	p.hp = newHP
}
func (p *Pelanggan) SetIDPegawai(newID int) {
	p.id_pegawai = newID
}
func (p *Pelanggan) SetNamaPelanggan(newNama string) {
	p.nama_pelanggan = newNama
}
func (p *Pelanggan) SetNamaPegawi(newNama string) {
	p.nama_pegawai = newNama
}
func (p *Pelanggan) SetIsActive(newIsActive int8) {
	p.isActive = newIsActive
}

func (p *Pelanggan) GetHP() string {
	return p.hp
}
func (p *Pelanggan) GetIDPegawai() int {
	return p.id_pegawai
}
func (p *Pelanggan) GetNamaPelanggan() string {
	return p.nama_pelanggan
}
func (p *Pelanggan) GetNamaPegawai() string {
	return p.nama_pegawai
}
func (p *Pelanggan) GetIsActive() int8 {
	return p.isActive
}

func (pm *PelangganMenu) Duplicate(nomer_hp string) (bool, error) {
	res := pm.db.QueryRow(`
	SELECT hp, isActive
	FROM pelanggan 
	WHERE hp = ?
	`, nomer_hp)
	var tmp Pelanggan
	if err := res.Scan(&tmp.hp, tmp.isActive); err != nil {
		if err.Error() != "sql: no rows in result set" {
			log.Println("Result scan error", err.Error())
			return false, err
		}
	}
	if len(tmp.hp) > 0 && tmp.isActive == 1 {
		return true, errors.New("nomer hp sudah digunakan")
	}
	return false, nil
}

func (pm *PelangganMenu) Register(newPelanggan Pelanggan) (bool, error) {
	registerQry, err := pm.db.Prepare(`
	INSERT INTO pelanggan
	(hp, id_pegawai, nama, isActive) values (?,?,?,1);
	`)
	if err != nil {
		log.Println("prepare register pelanggan ", err.Error())
		return false, errors.New("prepare statement register pelanggan error")
	}
	isDuplicate, err := pm.Duplicate(newPelanggan.hp)
	if err != nil {
		log.Println("error duplicated information")
		return false, err
	}
	if isDuplicate {
		log.Println("duplicated information")
		return false, err
	}
	res, err := registerQry.Exec(newPelanggan.hp, newPelanggan.id_pegawai, newPelanggan.nama_pelanggan)
	if err != nil {
		log.Println("register pelanggan", err.Error())
		return false, errors.New("register pelanggan error")
	}
	affRows, err := res.RowsAffected()
	if err != nil {
		log.Println("after register pelanggan", err.Error())
		return false, errors.New("error setelah register pelanggan")
	}
	if affRows <= 0 {
		log.Println("no record affected")
		return true, errors.New("no record")
	}
	return true, nil
}

func (pm *PelangganMenu) Data(nomer_hp string) ([]Pelanggan, string, error) {
	var (
		selectPelangganQry *sql.Rows
		err                error
		strPelanggan       string
	)
	if len(nomer_hp) > 0 {
		selectPelangganQry, err = pm.db.Query(`
		SELECT p.hp "Nomer HP", p.nama "Nama Pelanggan", p.id_pegawai "ID Pegawai", p2.nama "Nama Pegawai"  
		FROM pelanggan p  
		JOIN pegawai p2 ON p2.id = p.id_pegawai
		WHERE p.hp = ?
		AND p.isActive = 1;`, nomer_hp)
	} else {
		selectPelangganQry, err = pm.db.Query(`
		SELECT p.hp "Nomer HP", p.nama "Nama Pelanggan", p.id_pegawai "ID Pegawai", p2.nama "Nama Pegawai"  
		FROM pelanggan p  
		JOIN pegawai p2 ON p2.id = p.id_pegawai 
		WHERE p.isActive = 1;`)
	}
	if err != nil {
		log.Println("select query data pelanggan", err.Error())
		return nil, strPelanggan, errors.New("select query data pelanggan error")
	}

	arrPelanggan := []Pelanggan{}
	for selectPelangganQry.Next() {
		var tmp Pelanggan
		err = selectPelangganQry.Scan(&tmp.hp, &tmp.nama_pelanggan, &tmp.id_pegawai, &tmp.nama_pegawai)
		if err != nil {
			log.Println("Loop through rows, using Scan to assign column data to struct fields", err.Error())
			return nil, strPelanggan, err
		}
		strPelanggan += fmt.Sprintf("%s\t| %s <%s>\n", tmp.hp, tmp.nama_pelanggan, tmp.nama_pegawai)
		arrPelanggan = append(arrPelanggan, tmp)
	}

	return arrPelanggan, strPelanggan, nil
}

func (pm *PelangganMenu) ChangeIsActive(nomer_hp string, isActive int8) (bool, error) {
	updateQry, err := pm.db.Prepare(`
	UPDATE pelanggan
	SET isActive = ?
	WHERE hp = ?;`)
	if err != nil {
		log.Println("prepare change isActive pelanggan", err.Error())
		return false, errors.New("prepare statement change isActive error pelanggan")

	}
	res, err := updateQry.Exec(isActive, nomer_hp)
	if err != nil {
		log.Println("update isActive pelanggan", err.Error())
		return false, errors.New("update isActive pelanggan error")
	}
	affRow, err := res.RowsAffected()
	if err != nil {
		log.Println("after update isActive pelanggan ", err.Error())
		return false, errors.New("error setelah update isActive pelanggan")
	}
	if affRow <= 0 {
		log.Println("no record affected")
		return false, errors.New("no record")
	}
	return true, nil
}

func (pm *PelangganMenu) Delete(nomer_hp string) (bool, error) {
	isChanged, err := pm.ChangeIsActive(nomer_hp, 0)
	if err != nil {
		return isChanged, err
	}

	return isChanged, nil
}
