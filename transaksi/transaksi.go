package transaksi

import (
	"database/sql"
	"errors"
	"log"
)

type Transaksi struct {
	id         int
	id_pegawai int
	hp         string
	tanggal    string
}

type TransaksiMenu struct {
	db *sql.DB
}

type TransaksiInterface interface {
	Insert(newTransaksi Transaksi) (int, error)
	Select(id int) ([]Transaksi, error)
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

func (tm *TransaksiMenu) Select(id int) ([]Transaksi, error) {
	var (
		selectTransaksiQry *sql.Rows
		err                error
	)
	if id == 0 {
		selectTransaksiQry, err = tm.db.Query(`
		SELECT id_pegawai,hp,tanggal
		FROM transaksi;`)
	} else {
		selectTransaksiQry, err = tm.db.Query(`
		SELECT id_pegawai,hp,tanggal
		FROM transaksi
		WHERE id = ?;`, id)
	}
	if err != nil {
		log.Println("select transaksi", err.Error())
		return nil, errors.New("select transaksi error")
	}

	arrTransaksi := []Transaksi{}
	for selectTransaksiQry.Next() {
		var tmp Transaksi
		err = selectTransaksiQry.Scan(&tmp.id_pegawai, &tmp.hp, &tmp.tanggal)
		if err != nil {
			log.Println("Loop through rows, using Scan to assign column data to struct fields", err.Error())
			return arrTransaksi, err
		}
		arrTransaksi = append(arrTransaksi, tmp)
	}
	return arrTransaksi, nil
}
