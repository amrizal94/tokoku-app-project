package barang

import (
	"database/sql"
	"errors"
	"log"
)

type Barang struct {
	barcode    int
	id_pegawai int
	nama       string
	stok       int
}
type BarangMenu struct {
	db *sql.DB
}

type BarangInterface interface {
	Insert(newBarang Barang) (bool, error)
	Select() ([]Barang, error)
	Delete(barcode int) (bool, error)
}

func NewBrangMenu(conn *sql.DB) BarangInterface {
	return &BarangMenu{
		db: conn,
	}
}

func (b *Barang) SetBarcode(newBarcode int) {
	b.barcode = newBarcode
}
func (b *Barang) SetIDPegawai(newIDPegawai int) {
	b.id_pegawai = newIDPegawai
}
func (b *Barang) SetNama(newNama string) {
	b.nama = newNama
}
func (b *Barang) SetStok(newStok int) {
	b.stok = newStok
}

func (b *Barang) GetBarcode() int {
	return b.barcode
}
func (b *Barang) GetIDPegawai() int {
	return b.id_pegawai
}
func (b *Barang) GetNama() string {
	return b.nama
}
func (b *Barang) GetStok() int {
	return b.stok
}

func (bm *BarangMenu) Insert(newBarang Barang) (bool, error) {
	insertBarang, err := bm.db.Prepare("INSERT INTO barang (barcode, id_pegawai, nama, stok) values (?,?,?,?)")
	if err != nil {
		log.Println("prepare insert barang ", err.Error())
		return false, errors.New("prepare statement insert barang error")
	}

	res, err := insertBarang.Exec(newBarang.barcode, newBarang.id_pegawai, newBarang.nama, newBarang.stok)

	if err != nil {
		log.Println("insert barang ", err.Error())
		return false, errors.New("insert barng error")
	}

	affRows, err := res.RowsAffected()

	if err != nil {
		log.Println("after insert barang ", err.Error())
		return false, errors.New("error setelah insert barang")
	}

	if affRows <= 0 {
		log.Println("no record affected")
		return false, errors.New("no record")
	}
	return true, nil
}

func (bm *BarangMenu) Select() ([]Barang, error) {
	selectBarangQry, err := bm.db.Query(`
	SELECT barcode, nama, stok
	FROM barang;`)
	if err != nil {
		log.Println("select barang", err.Error())
		return nil, errors.New("select barang error")
	}
	arrBarang := []Barang{}
	for selectBarangQry.Next() {
		var tmp Barang
		err = selectBarangQry.Scan(&tmp.barcode, &tmp.nama, &tmp.stok)
		if err != nil {
			log.Println("Loop through rows, using Scan to assign column data to struct fields", err.Error())
			return arrBarang, err
		}
		arrBarang = append(arrBarang, tmp)
	}
	return arrBarang, nil
}

func (bm *BarangMenu) Delete(barcode int) (bool, error) {
	deleteBarangQry, err := bm.db.Prepare("DELETE FROM barang WHERE barcode = ?;")
	if err != nil {
		log.Println("prepare delete barang ", err.Error())
		return false, errors.New("prepare statement delete barang error")
	}
	res, err := deleteBarangQry.Exec(barcode)
	if err != nil {
		log.Println("delete barang ", err.Error())
		return false, errors.New("delete barang error")
	}
	affRows, err := res.RowsAffected()
	if err != nil {
		log.Println("after delete barang ", err.Error())
		return false, errors.New("error setelah delete barang")
	}
	if affRows <= 0 {
		log.Println("no record affected")
		return false, errors.New("no record")
	}
	return true, nil
}
