package tambah_barang

import (
	"database/sql"
	"errors"
	"log"
)

type Barang struct {
	Barcode    int
	ID_pegawai string
	Nama       string
	Stok       int
}
type BarangMenuN struct {
	DB *sql.DB
}

func (ab *BarangMenuN) Insert(newBarang Barang) (int, error) {
	insertBarang, err := ab.DB.Prepare("INSERT INTO barang (barcode, title, nama, Stok) values (?,?,?,?)")
	if err != nil {
		log.Println("prepare insert barang ", err.Error())
		return 0, errors.New("prepare statement insert user error")
	}

	res, err := insertBarang.Exec(newBarang.Barcode, newBarang.ID_pegawai, newBarang.Nama, newBarang.Stok)

	if err != nil {
		log.Println("insert barang ", err.Error())
		return 0, errors.New("insert barng error")
	}

	affRows, err := res.RowsAffected()

	if err != nil {
		log.Println("after insert barang ", err.Error())
		return 0, errors.New("error setelah insert barang")
	}

	if affRows <= 0 {
		log.Println("no record affected")
		return 0, errors.New("no record")
	}

	id, _ := res.LastInsertId()

	return int(id), nil
}
