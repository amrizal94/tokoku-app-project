package transaksibarang

import (
	"database/sql"
	"errors"
	"log"
)

type Transaksibarang struct {
	transaksi_id int
	barcode_id   int
}

type TransaksibarangMenu struct {
	db *sql.DB
}

type TransaksibarangInterface interface {
	Invite(barcode, id_transaksi int) (bool, error)
}

func NewTransaksibarangMenu(conn *sql.DB) TransaksibarangInterface {
	return &TransaksibarangMenu{
		db: conn,
	}
}

func (ua *Transaksibarang) TransaksiID(newUserID int) {
	ua.transaksi_id = newUserID
}
func (ua *Transaksibarang) BarcodeID(newBarcodeIDID int) {
	ua.barcode_id = newBarcodeIDID
}

func (ua *Transaksibarang) GetTransaksiID() int {
	return ua.transaksi_id
}
func (ua *Transaksibarang) GetBarcodeID() int {
	return ua.transaksi_id
}
func (ua *TransaksibarangMenu) Invite(transaksi_id, barcode_id int) (bool, error) {
	inviteQry, err := ua.db.Prepare("INSERT INTO user_activities (transaksi_id, barcode_id, due_date) values (?, ?, now())")
	if err != nil {
		log.Println("prepare invite user transaksi ", err.Error())
		return false, errors.New("prepare statement invite user transaksi error")
	}
	res, err := inviteQry.Exec(transaksi_id, barcode_id)
	if err != nil {
		log.Println("invite user transaksi ", err.Error())
		return false, errors.New("invite user transaksi error")
	}
	affRows, err := res.RowsAffected()
	if err != nil {
		log.Println("after insert user transaksi ", err.Error())
		return false, errors.New("error setelah insert usertransaksi")
	}
	if affRows <= 0 {
		log.Println("no record affected")
		return false, errors.New("no record")
	}

	return true, nil

}
