package pelanggan

import (
	"database/sql"
	"errors"
	"log"
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
