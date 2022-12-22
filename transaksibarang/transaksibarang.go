package transaksibarang

import (
	"database/sql"
	"errors"
	"log"
)

type TransaksiBarang struct {
	id_transaksi int
	nama         string
	barcode      int
	jumlah       int
	harga        int
	total        int
}

type TransaksiBarangMenu struct {
	db *sql.DB
}

type TransaksiBarangInterface interface {
	Insert(newTransaksiBarang TransaksiBarang) (bool, error)
	Select(id_transaksi int) ([]TransaksiBarang, error)
}

func NewTransaksiBarangMenu(conn *sql.DB) TransaksiBarangInterface {
	return &TransaksiBarangMenu{
		db: conn,
	}
}

func (tb *TransaksiBarang) SetIDTransaksi(newIDTransaksi int) {
	tb.id_transaksi = newIDTransaksi
}
func (tb *TransaksiBarang) SetBarcode(newBarcode int) {
	tb.barcode = newBarcode
}
func (tb *TransaksiBarang) SetJumlah(newJumlah int) {
	tb.jumlah = newJumlah
}
func (tb *TransaksiBarang) SetNama(newNama string) {
	tb.nama = newNama
}
func (tb *TransaksiBarang) SetTotal(newTotal int) {
	tb.jumlah = newTotal
}
func (tb *TransaksiBarang) SetHarga(newHarga int) {
	tb.harga = newHarga
}

func (tb *TransaksiBarang) GetIDTransaksi() int {
	return tb.id_transaksi
}
func (tb *TransaksiBarang) GetBarcode() int {
	return tb.barcode
}
func (tb *TransaksiBarang) GetJumlah() int {
	return tb.jumlah
}
func (tb *TransaksiBarang) GetNama() string {
	return tb.nama
}
func (tb *TransaksiBarang) GetTotal() int {
	return tb.total
}
func (tb *TransaksiBarang) GetHarga() int {
	return tb.harga
}

func (tbm *TransaksiBarangMenu) Duplicate(newTB TransaksiBarang) (bool, error) {
	res := tbm.db.QueryRow(`
	SELECT id_transaksi, barcode
	FROM transaksi_barang 
	WHERE id_transaksi = ?
	AND barcode = ?
	`, newTB.id_transaksi, newTB.barcode)
	var id_transaksi, barcode int
	if err := res.Scan(&id_transaksi, &barcode); err != nil {
		if err.Error() != "sql: no rows in result set" {
			log.Println("Result scan error", err.Error())
			return false, err
		}
	}
	if id_transaksi > 0 && barcode > 0 {
		return true, nil
	}
	return false, nil
}

func (tbm *TransaksiBarangMenu) Update(newTB TransaksiBarang) (bool, error) {
	updateQry, err := tbm.db.Prepare(`
	UPDATE transaksi_barang
	SET jumlah = jumlah + ?
	WHERE id_transaksi = ?
	AND barcode = ?;`)
	if err != nil {
		log.Println("prepare update transaksi barang", err.Error())
		return false, errors.New("prepare statement update transaksi barang error")
	}
	res, err := updateQry.Exec(newTB.jumlah, newTB.id_transaksi, newTB.barcode)
	if err != nil {
		log.Println("update transaksi barang", err.Error())
		return false, errors.New("update transaksi barang error")
	}
	affRow, err := res.RowsAffected()
	if err != nil {
		log.Println("after update transaksi barang", err.Error())
		return false, errors.New("error setelah update transaksi barang")
	}
	if affRow <= 0 {
		log.Println("no record affected")
		return false, errors.New("no record")
	}
	return true, nil
}

func (tbm *TransaksiBarangMenu) Insert(newTransaksiBarang TransaksiBarang) (bool, error) {
	isDuplicate, err := tbm.Duplicate(newTransaksiBarang)
	if err != nil {
		log.Println("error duplicate transaksi barang", err.Error())
		return false, errors.New("duplicate transaksi barang error")
	}

	if isDuplicate {
		isUpdated, err := tbm.Update(newTransaksiBarang)
		if err != nil {
			return isUpdated, err
		}

		return isUpdated, err
	}
	insertQry, err := tbm.db.Prepare(`
		INSERT INTO transaksi_barang (id_transaksi, barcode, jumlah) values (?,?,?)`)
	if err != nil {
		log.Println("prepare insert transaksi barang ", err.Error())
		return false, errors.New("prepare statement insert transaksi barang error")
	}
	res, err := insertQry.Exec(newTransaksiBarang.id_transaksi, newTransaksiBarang.barcode, newTransaksiBarang.jumlah)
	if err != nil {
		log.Println("insert transaksi barang ", err.Error())
		return false, errors.New("insert transaksi barang error")
	}
	affRows, err := res.RowsAffected()
	if err != nil {
		log.Println("after insert transaksi barang ", err.Error())
		return false, errors.New("error setelah insert transaksi barang")
	}
	if affRows <= 0 {
		log.Println("no record affected")
		return false, errors.New("no record")
	}

	return true, nil
}

func (tbm *TransaksiBarangMenu) Select(id_transaksi int) ([]TransaksiBarang, error) {
	selectTransaksiBarangQry, err := tbm.db.Query(`
	SELECT b.nama, tb.jumlah, b.harga, tb.jumlah * b.harga 
	FROM barang b 
	JOIN transaksi_barang tb ON tb.barcode = b.barcode
	WHERE tb.id_transaksi = ?;`, id_transaksi)
	if err != nil {
		log.Println("select transaksi barang", err.Error())
		return nil, errors.New("select transaksi barang error")
	}
	arrtransaksibarang := []TransaksiBarang{}
	for selectTransaksiBarangQry.Next() {
		var tmp TransaksiBarang
		err = selectTransaksiBarangQry.Scan(&tmp.nama, &tmp.jumlah, &tmp.harga, &tmp.total)
		if err != nil {
			log.Println("Loop through rows, using Scan to assign column data to struct fields", err.Error())
			return arrtransaksibarang, err
		}
		arrtransaksibarang = append(arrtransaksibarang, tmp)
	}
	return arrtransaksibarang, nil
}
