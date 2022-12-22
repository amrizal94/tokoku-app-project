package transaksi

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
)

type Transaksi struct {
	id             int
	id_pegawai     int
	hp             string
	nama_pegawai   string
	nama_pelanggan string
	tanggal        string
}

type TransaksiMenu struct {
	db *sql.DB
}

type TransaksiInterface interface {
	Insert(newTransaksi Transaksi) (int, error)
	Data(id int) ([]Transaksi, string, error)
	Delete(id_transaksi int) (bool, error)
}

func NewTransaksiMenu(conn *sql.DB) TransaksiInterface {
	return &TransaksiMenu{
		db: conn,
	}
}

func (t *Transaksi) SetID(newID int) {
	t.id = newID
}
func (t *Transaksi) SetIDPegawai(newIDPegawai int) {
	t.id_pegawai = newIDPegawai
}
func (t *Transaksi) SetHP(newHP string) {
	t.hp = newHP
}
func (t *Transaksi) SetTanggal(newTanggal string) {
	t.tanggal = newTanggal
}
func (t *Transaksi) SetNamaPegawai(newNama string) {
	t.nama_pegawai = newNama
}
func (t *Transaksi) SetNamaPelanggan(newNama string) {
	t.nama_pelanggan = newNama
}
func (t *Transaksi) GetID() int {
	return t.id
}
func (t *Transaksi) GetIDPegawai() int {
	return t.id_pegawai
}
func (t *Transaksi) GetHP() string {
	return t.hp
}
func (t *Transaksi) GetTanggal() string {
	return t.tanggal
}
func (t *Transaksi) GetNamaPegawai() string {
	return t.nama_pegawai
}
func (t *Transaksi) GetNamaPelanggan() string {
	return t.nama_pelanggan
}

func (tm *TransaksiMenu) Insert(newTransaksi Transaksi) (int, error) {
	insertQry, err := tm.db.Prepare(`
	INSERT INTO transaksi (id_pegawai, hp, tanggal) values (?,?,now())`)
	if err != nil {
		log.Println("prepare insert transaksi ", err.Error())
		return 0, errors.New("prepare statement insert transaksi error")
	}
	res, err := insertQry.Exec(newTransaksi.id_pegawai, newTransaksi.hp)
	if err != nil {
		log.Println("insert transaksi ", err.Error())
		return 0, errors.New("insert transaksi error")
	}
	affRows, err := res.RowsAffected()
	if err != nil {
		log.Println("after insert transaksi ", err.Error())
		return 0, errors.New("error setelah insert transaksi")
	}
	if affRows <= 0 {
		log.Println("no record affected")
		return 0, errors.New("no record")
	}
	id, _ := res.LastInsertId()
	return int(id), nil
}

func (tm *TransaksiMenu) Data(id int) ([]Transaksi, string, error) {
	var (
		selectTransaksiQry *sql.Rows
		err                error
		strTransaksi       string
	)
	if id == 0 {
		selectTransaksiQry, err = tm.db.Query(`
		SELECT t.id ,t.tanggal ,t.id_pegawai ,p.nama "Nama Pegawai" , t.hp ,p2.nama as "Nama Pelanggan"
		FROM transaksi t 
		JOIN pegawai p ON p.id = t.id_pegawai 
		JOIN pelanggan p2 ON p2.hp = t.hp;`)
	} else {
		selectTransaksiQry, err = tm.db.Query(`
		SELECT t.id ,t.tanggal ,t.id_pegawai ,p.nama "Nama Pegawai" , t.hp ,p2.nama as "Nama Pelanggan"
		FROM transaksi t 
		JOIN pegawai p ON p.id = t.id_pegawai 
		JOIN pelanggan p2 ON p2.hp = t.hp
		WHERE t.id = ?;`, id)
	}
	if err != nil {
		log.Println("select transaksi", err.Error())
		return nil, strTransaksi, errors.New("select transaksi error")
	}

	arrTransaksi := []Transaksi{}
	for selectTransaksiQry.Next() {
		var tmp Transaksi
		err = selectTransaksiQry.Scan(&tmp.id, &tmp.tanggal, &tmp.id_pegawai, &tmp.nama_pegawai, &tmp.hp, &tmp.nama_pelanggan)
		if err != nil {
			log.Println("Loop through rows, using Scan to assign column data to struct fields selectTransaksiQry", err.Error())
			return arrTransaksi, strTransaksi, err
		}
		strTransaksi += fmt.Sprintf("%d\t\t| %s\t| %s\t| %s\n", tmp.id, tmp.tanggal, tmp.nama_pegawai, tmp.nama_pelanggan)
		arrTransaksi = append(arrTransaksi, tmp)
	}
	return arrTransaksi, strTransaksi, nil
}

func (tm *TransaksiMenu) Delete(id_transaksi int) (bool, error) {
	deleteQry, err := tm.db.Prepare(`
	DELETE FROM transaksi 
	WHERE id = ?;
	`)
	if err != nil {
		log.Println("prepare delete transaksi ", err.Error())
		return false, errors.New("prepare statement delete transaksi error")
	}
	res, err := deleteQry.Exec(id_transaksi)
	if err != nil {
		log.Println("delete transaksi ", err.Error())
		return false, errors.New("delete transaksi error")
	}
	affRows, err := res.RowsAffected()
	if err != nil {
		log.Println("after delete transaksi ", err.Error())
		return false, errors.New("error setelah delete transaksi")
	}

	if affRows <= 0 {
		log.Println("no record affected")
		return false, errors.New("no record")
	}
	return true, nil
}
