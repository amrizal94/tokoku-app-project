package pelanggan

import (
	"database/sql"
	"errors"
	"log"
	"strconv"
)

type Pelanggan struct {
	hp         string
	id_pegawai int
	nama       string
}

type PelangganMenu struct {
	db *sql.DB
}

type PelangganInterface interface {
	Insert(newPelanggan Pelanggan) (bool, error)
	Select(hp string) ([]Pelanggan, error)
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
func (p *Pelanggan) SetIDPegawai(newIDPegawai int) {
	p.id_pegawai = newIDPegawai
}
func (p *Pelanggan) SetNama(newNama string) {
	p.nama = newNama
}

func (p *Pelanggan) GetHP() string {
	return p.hp
}
func (p *Pelanggan) GetIDPegawai() int {
	return p.id_pegawai
}
func (p *Pelanggan) GetNama() string {
	return p.nama
}

func (pm *PelangganMenu) Insert(newPelanggan Pelanggan) (bool, error) {
	insertQry, err := pm.db.Prepare(`
	INSERT INTO pelanggan (hp, id_pegawai, nama) values (?,?,?)`)
	if err != nil {
		log.Println("prepare insert pelanggan ", err.Error())
		return false, errors.New("prepare statement insert pelanggan error")
	}
	res, err := insertQry.Exec(newPelanggan.hp, newPelanggan.id_pegawai, newPelanggan.nama)
	if err != nil {
		log.Println("insert pelanggan ", err.Error())
		return false, errors.New("insert pelanggan error")
	}
	affRows, err := res.RowsAffected()
	if err != nil {
		log.Println("after insert pelanggan ", err.Error())
		return false, errors.New("error setelah insert pelanggan")
	}
	if affRows <= 0 {
		log.Println("no record affected")
		return false, errors.New("no record")
	}
	return true, nil
}

func (pm *PelangganMenu) Select(hp string) ([]Pelanggan, error) {
	var (
		selectPelangganQry *sql.Rows
		err                error
	)
	intHP, err := strconv.Atoi(hp)
	if err != nil {
		log.Println("convert string to integer", err.Error())
		return nil, errors.New("convert string to integer")
	}
	if intHP == 0 {
		selectPelangganQry, err = pm.db.Query(`
		SELECT hp,id_pegawai,nama
		FROM pelanggan;`)
	} else {
		selectPelangganQry, err = pm.db.Query(`
		SELECT hp,id_pegawai,nama
		FROM pelanggan
		WHERE hp = ?;`, intHP)
	}
	if err != nil {
		log.Println("select pelanggan", err.Error())
		return nil, errors.New("select pelanggan error")
	}

	arrPelanggan := []Pelanggan{}
	for selectPelangganQry.Next() {
		var tmp Pelanggan
		err = selectPelangganQry.Scan(&tmp.hp, &tmp.id_pegawai, &tmp.nama)
		if err != nil {
			log.Println("Loop through rows, using Scan to assign column data to struct fields", err.Error())
			return arrPelanggan, err
		}
		arrPelanggan = append(arrPelanggan, tmp)
	}
	return arrPelanggan, nil
}

func (pm *PelangganMenu) Delete(hp string) (bool, error) {
	deletePelangganQry, err := pm.db.Prepare("DELETE FROM pelanggan WHERE hp = ?;")
	if err != nil {
		log.Println("prepare delete pelanggan ", err.Error())
		return false, errors.New("prepare statement delete pelanggan error")
	}
	res, err := deletePelangganQry.Exec(hp)
	if err != nil {
		log.Println("delete pelanggan ", err.Error())
		return false, errors.New("delete pelanggan error")
	}
	affRows, err := res.RowsAffected()
	if err != nil {
		log.Println("after delete pelanggan ", err.Error())
		return false, errors.New("error setelah delete pelanggan")
	}
	if affRows <= 0 {
		log.Println("no record affected")
		return false, errors.New("no record")
	}
	return true, nil
}
