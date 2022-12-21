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
	harga      int
}
type BarangMenu struct {
	db *sql.DB
}

type BarangInterface interface {
	Insert(newBarang Barang) (bool, error)
	Select(barcode int) ([]Barang, error)
	Delete(barcode int) (bool, error)
	Update(barcode int, nama string, stok int, harga int) (bool, error)
	Sell(barcode int, jumlah int) (bool, error)
}

func NewBarangMenu(conn *sql.DB) BarangInterface {
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
func (b *Barang) SetHarga(newHarga int) {
	b.harga = newHarga
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
func (b *Barang) GetHarga() int {
	return b.harga
}

func (bm *BarangMenu) Insert(newBarang Barang) (bool, error) {
	insertBarang, err := bm.db.Prepare("INSERT INTO barang (barcode, id_pegawai, nama, stok, harga) values (?,?,?,?,?)")
	if err != nil {
		log.Println("prepare insert barang ", err.Error())
		return false, errors.New("prepare statement insert barang error")
	}

	res, err := insertBarang.Exec(newBarang.barcode, newBarang.id_pegawai, newBarang.nama, newBarang.stok, newBarang.harga)

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

func (bm *BarangMenu) Select(barcode int) ([]Barang, error) {
	var (
		selectBarangQry *sql.Rows
		err             error
	)
	if barcode == 0 {
		selectBarangQry, err = bm.db.Query(`
		SELECT barcode,id_pegawai,nama,stok,harga
		FROM barang;`)
	} else {
		selectBarangQry, err = bm.db.Query(`
		SELECT barcode,id_pegawai,nama,stok,harga
		FROM barang
		WHERE barcode = ?;`, barcode)
	}

	if err != nil {
		log.Println("select barang", err.Error())
		return nil, errors.New("select barang error")
	}
	arrBarang := []Barang{}
	for selectBarangQry.Next() {
		var tmp Barang
		err = selectBarangQry.Scan(&tmp.barcode, &tmp.id_pegawai, &tmp.nama, &tmp.stok, &tmp.harga)
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

func (bm *BarangMenu) Update(barcode int, nama string, stok int, harga int) (bool, error) {
	updateQry, err := bm.db.Prepare(`
	UPDATE barang
	SET nama = ?, stok = ?, harga = ?
	WHERE barcode = ?;`)
	if err != nil {
		log.Println("prepare update barang", err.Error())
		return false, errors.New("prepare statement update barang error")
	}
	res, err := updateQry.Exec(nama, stok, harga, barcode)
	if err != nil {
		log.Println("update barang ", err.Error())
		return false, errors.New("update barang error")
	}
	affRow, err := res.RowsAffected()
	if err != nil {
		log.Println("after update barang ", err.Error())
		return false, errors.New("error setelah update barang")
	}
	if affRow <= 0 {
		log.Println("no record affected")
		return false, errors.New("no record")
	}
	return true, nil
}

func (bm *BarangMenu) Sell(barcode int, jumlah int) (bool, error) {
	updateQry, err := bm.db.Prepare(`
	UPDATE barang
	SET stok = stok - ? 
	WHERE barcode = ? and stok > ?;`)
	if err != nil {
		log.Println("prepare update barang", err.Error())
		return false, errors.New("prepare statement update barang error")
	}
	res, err := updateQry.Exec(jumlah, barcode, jumlah)
	if err != nil {
		log.Println("update barang ", err.Error())
		return false, errors.New("update barang error")
	}
	affRow, err := res.RowsAffected()
	if err != nil {
		log.Println("after update barang ", err.Error())
		return false, errors.New("error setelah update barang")
	}
	if affRow <= 0 {
		log.Println("no record affected stok kurang")
		return false, errors.New("no record stok kurang")
	}
	return true, nil
}
