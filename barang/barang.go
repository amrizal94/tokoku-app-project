package barang

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
)

type Barang struct {
	barcode      int
	id_pegawai   int
	nama_barang  string
	nama_pegawai string
	stok         int
	harga        int
	isActive     int8
}
type BarangMenu struct {
	db *sql.DB
}

type BarangInterface interface {
	Register(newBarang Barang) (bool, error)
	Data(barcode int) ([]Barang, string, error)
	Delete(barcode int) (bool, error)
	Update(upBarang Barang) (bool, error)
	Sell(barcode int, jumlah int) (bool, error)
	IncreaseStok(inStok, barcode int) (bool, error)
}

func NewBarangMenu(conn *sql.DB) BarangInterface {
	return &BarangMenu{
		db: conn,
	}
}

func (b *Barang) SetBarcode(newBarcode int) {
	b.barcode = newBarcode
}
func (b *Barang) SetIDPegawai(newID int) {
	b.id_pegawai = newID
}
func (b *Barang) SetNamaBarang(newNama string) {
	b.nama_barang = newNama
}
func (b *Barang) SetNamaPegawai(newNama string) {
	b.nama_pegawai = newNama
}
func (b *Barang) SetStok(newStok int) {
	b.stok = newStok
}
func (b *Barang) SetHarga(newHarga int) {
	b.harga = newHarga
}
func (b *Barang) SetIsActive(newIsActive int8) {
	b.isActive = newIsActive
}

func (b *Barang) GetBarcode() int {
	return b.barcode
}
func (b *Barang) GetIDPegawai() int {
	return b.id_pegawai
}
func (b *Barang) GetNamaBarang() string {
	return b.nama_barang
}
func (b *Barang) GetNamaPegawai() string {
	return b.nama_pegawai
}
func (b *Barang) GetStok() int {
	return b.stok
}
func (b *Barang) GetHarga() int {
	return b.harga
}
func (b *Barang) GetIsActive() int8 {
	return b.isActive
}

func (bm *BarangMenu) Duplicate(newBarang Barang) (bool, error) {
	res := bm.db.QueryRow(`
	SELECT barcode isActive
	FROM barang 
	WHERE barcode = ?
	`, newBarang.barcode)
	var isActive int8
	var barcode int
	if err := res.Scan(&barcode, &isActive); err != nil {
		if err.Error() != "sql: no rows in result set" {
			log.Println("Result scan error", err.Error())
			return false, err
		}
	}
	if isActive == 1 {
		return true, errors.New("barcode sudah digunakan")
	} else if barcode > 0 {
		bm.Update(newBarang)
	}
	return false, nil
}

// ////menu regis
func (bm *BarangMenu) Register(newBarang Barang) (bool, error) {
	regiterQry, err := bm.db.Prepare(`
	INSERT INTO barang
	(barcode, id_pegawai, nama, stok, harga, isActive) values (?,?,?,?,?,1)
	`)
	if err != nil {
		log.Println("prepare register barang ", err.Error())
		return false, errors.New("prepare statement register barang error")
	}
	isDuplicate, err := bm.Duplicate(newBarang)
	if err != nil {
		return false, err
	}
	if isDuplicate {
		log.Println("duplicated information")
		return false, err
	}
	res, err := regiterQry.Exec(newBarang.barcode, newBarang.id_pegawai, newBarang.nama_barang, newBarang.stok, newBarang.harga)

	if err != nil {
		log.Println("register barang ", err.Error())
		return false, errors.New("register barang error")
	}

	affRows, err := res.RowsAffected()

	if err != nil {
		log.Println("after register barang ", err.Error())
		return false, errors.New("error setelah register barang")
	}

	if affRows <= 0 {
		log.Println("no record affected")
		return false, errors.New("no record")
	}
	return true, nil
}

// /// menu data
func (bm *BarangMenu) Data(barcode int) ([]Barang, string, error) {
	var (
		selectBarangQry *sql.Rows
		err             error
		strBarang       string
	)
	if barcode == 0 {
		selectBarangQry, err = bm.db.Query(`
		SELECT b.barcode ,b.id_pegawai ,b.nama "Nama Barang" ,b.stok ,b.harga ,p.nama 'Nama Pegawai'
		FROM barang b
		JOIN pegawai p ON p.id = b.id_pegawai
		WHERE b.isActive = 1;`)
	} else {
		selectBarangQry, err = bm.db.Query(`
		SELECT b.barcode ,b.id_pegawai ,b.nama "Nama Barang" ,b.stok ,b.harga ,p.nama 'Nama Pegawai'
		FROM barang b
		JOIN pegawai p ON p.id = b.id_pegawai
		WHERE b.barcode = ?
		AND b.isActive = 1;`, barcode)
	}

	if err != nil {
		log.Println("select barang", err.Error())
		return nil, strBarang, errors.New("select barang error")
	}
	arrBarang := []Barang{}
	for selectBarangQry.Next() {
		var tmp Barang
		err = selectBarangQry.Scan(&tmp.barcode, &tmp.id_pegawai, &tmp.nama_barang, &tmp.stok, &tmp.harga, &tmp.nama_pegawai)
		if err != nil {
			log.Println("Loop through rows, using Scan to assign column data to struct fields", err.Error())
			return arrBarang, strBarang, err
		}
		strBarang += fmt.Sprintf("%d\t| %s {%d} [%d] <%s>\n", tmp.barcode, tmp.nama_barang, tmp.stok, tmp.harga, tmp.nama_pegawai)
		arrBarang = append(arrBarang, tmp)
	}
	return arrBarang, strBarang, nil
}

func (bm *BarangMenu) Update(upBarang Barang) (bool, error) {
	updateQry, err := bm.db.Prepare(`
	UPDATE barang
	SET nama = ?, stok = ?, harga = ?
	WHERE barcode = ?;`)
	if err != nil {
		log.Println("prepare update barang", err.Error())
		return false, errors.New("prepare statement update barang error")
	}
	res, err := updateQry.Exec(upBarang.nama_barang, upBarang.stok, upBarang.harga, upBarang.barcode)
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
	WHERE barcode = ?
	AND stok > ?;`)
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
		log.Println("no record affected barang")
		return false, errors.New("no record barang")
	}
	return true, nil
}

func (bm *BarangMenu) ChangeIsActive(barcode int, isActive int8) (bool, error) {
	updateQry, err := bm.db.Prepare(`
	UPDATE barang
	SET isActive = ?
	WHERE barcode = ?;`)
	if err != nil {
		log.Println("prepare change isActive barang", err.Error())
		return false, errors.New("prepare statement change isActive error barang")
	}
	res, err := updateQry.Exec(isActive, barcode)
	if err != nil {
		log.Println("update isActive barang", err.Error())
		return false, errors.New("update isActive barang error")
	}
	affRow, err := res.RowsAffected()
	if err != nil {
		log.Println("after update isActive barang ", err.Error())
		return false, errors.New("error setelah update isActive barang")
	}
	if affRow <= 0 {
		log.Println("no record affected")
		return false, errors.New("no record")
	}
	return true, nil
}

func (bm *BarangMenu) Delete(barcode int) (bool, error) {
	isChanged, err := bm.ChangeIsActive(barcode, 0)
	if err != nil {
		return isChanged, err
	}

	return isChanged, err
}

func (bm *BarangMenu) IncreaseStok(inStok, barcode int) (bool, error) {
	increaseQry, err := bm.db.Prepare(`
	UPDATE barang
	SET stok = stok + ?
	WHERE barcode = ?;`)
	if err != nil {
		log.Println("prepare change tambah stok barang", err.Error())
		return false, errors.New("prepare statement change tambah stok error barang")
	}
	res, err := increaseQry.Exec(inStok, barcode)
	if err != nil {
		log.Println("update tambah stok barang", err.Error())
		return false, errors.New("update tambah stok barang error")
	}
	affRow, err := res.RowsAffected()
	if err != nil {
		log.Println("after update tambah stok barang ", err.Error())
		return false, errors.New("error setelah update tambah stok barang")
	}
	if affRow <= 0 {
		log.Println("no record affected")
		return false, errors.New("no record")
	}
	return true, nil
}
