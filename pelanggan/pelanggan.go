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

// Register implements PelangganInterface
type PelangganInterface interface {
	Register(newPelanggan Pelanggan) (bool, int, error)
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
func (p *Pelanggan) SetName(newNama string) {
	p.nama = newNama
}

func (p *Pelanggan) GetHP() string {
	return p.hp
}
func (p *Pelanggan) GetIDPegawai() int {
	return p.id_pegawai
}
func (p *Pelanggan) GetName() string {
	return p.nama
}

func (ap *PelangganMenu) Register(newPelanggan Pelanggan) (bool, int, error) {
	registerQry, err := ap.db.Prepare("INSERT INTO pelanggan (hp, id_pegawai, nama) values (?,?,?)")
	if err != nil {
		log.Println("prepare insert pegawai registerQry", err.Error())
		return false, 0, errors.New("prepare statement insert pelanggan error registerQry")
	}

	// menjalankan query dengan parameter tertentu
	res, err := registerQry.Exec(newPelanggan.GetHP(), newPelanggan.GetIDPegawai(), newPelanggan.GetName())
	if err != nil {
		log.Println("insert pelanggan registerQry ", err.Error())
		return false, 0, errors.New("insert pelanggan error registerQry")
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
