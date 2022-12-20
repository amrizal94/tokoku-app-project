package main

import (
	"bufio"
	"fmt"
	"os"
	"tokoku-app-project/barang"
	"tokoku-app-project/config"
	"tokoku-app-project/pegawai"
	// "tokoku-app-project/pelanggan"
)

var (
	cfg         = config.ReadConfig()
	conn        = config.ConnectSQL(*cfg)
	PegawaiMenu = pegawai.NewPegawaiMenu(conn)
	BarangMenu  = barang.NewBrangMenu(conn)
)

func listPegawai(id, id_logged int) ([]pegawai.Pegawai, string, error) {
	arrPegawai, err := PegawaiMenu.Select(id, id_logged)
	var strPegawai string
	if err != nil {
		fmt.Println(err.Error())
		return arrPegawai, strPegawai, err
	}
	for _, v := range arrPegawai {
		strPegawai += fmt.Sprintln("ID :", v.GetID(), v.GetNama())
	}

	return arrPegawai, strPegawai, err
}

func listBarang() ([]barang.Barang, string, error) {
	arrBarang, err := BarangMenu.Select()
	var strBarang string
	if err != nil {
		fmt.Println(err.Error())
		return arrBarang, strBarang, err
	}
	for _, v := range arrBarang {
		strBarang += fmt.Sprintln("Barcode :", v.GetBarcode(), v.GetNama())
	}

	return arrBarang, strBarang, err
}

func main() {
	var (
		cfg         = config.ReadConfig()
		conn        = config.ConnectSQL(*cfg)
		PegawaiMenu = pegawai.NewPegawaiMenu(conn)
		// Pelanggan       =pelanggan.NewPelangganMenu
		inputMenu int = 1
	)
	for inputMenu != 0 {
		fmt.Println("==========================")
		fmt.Println("1. Login")
		fmt.Println("0. Exit")
		fmt.Scanln(&inputMenu)
		switch inputMenu {
		case 1:
			var inNama, inPassword string
			fmt.Println("==========================")
			fmt.Println("LOGIN")
			fmt.Print("Masukkan username : ")
			fmt.Scanln(&inNama)
			fmt.Print("Masukkan password : ")
			fmt.Scanln(&inPassword)
			resLogin, err := PegawaiMenu.Login(inNama, inPassword)
			if err != nil {
				fmt.Println(err.Error())
			}
			if resLogin.GetID() > 0 {
				fmt.Println("==========================")
				fmt.Println("Login sukses, selamat datang", resLogin.GetNama())
				isLogged := true
				var isAdmin bool
				if resLogin.GetUsername() == "admin" && inPassword == "admin" {
					isAdmin = !isAdmin
				}
				for isLogged {
					fmt.Println("==========================")
					if isAdmin {
						fmt.Println("1. Tambah Pegawai")
						fmt.Println("2. Hapus Pegawai")
						fmt.Println("3. Hapus Barang")
					} else {
						fmt.Println("1. Tambah Pelanggan")
						fmt.Println("2. Tambah Barang")
						fmt.Println("3. Edit Barang")
					}
					fmt.Println("9. Log out")
					fmt.Println("0. Exit")
					fmt.Scanln(&inputMenu)
					switch inputMenu {
					case 1:
						if isAdmin {
							var newPegawai pegawai.Pegawai
							var tmp string
							reader := bufio.NewReader(os.Stdin)
							fmt.Println("==========================")
							fmt.Println("TAMBAH PEGAWAI")
							fmt.Print("Masukkan nama : ")
							nama, _ := reader.ReadString('\n')
							nama = nama[:len(nama)-1]
							newPegawai.SetNama(nama)
							fmt.Print("Masukkan username : ")
							fmt.Scanln(&tmp)
							newPegawai.SetUsername(tmp)
							fmt.Print("Masukkan password : ")
							fmt.Scanln(&tmp)
							newPegawai.SetPassword(tmp)
							isAdded, isActive, err := PegawaiMenu.Register(newPegawai)
							if err != nil {
								fmt.Println(err.Error())
							} else {
								if isActive > 0 {
									isAdded, err = PegawaiMenu.Update(newPegawai.GetPassword(), newPegawai.GetNama(), int(newPegawai.GetIsActive()), isActive)
									if err != nil {
										fmt.Println(err.Error())
									}
								}

							}
							if isAdded {
								fmt.Println("==========================")
								fmt.Println("Sukses menambahkan pegawai")
							} else {
								fmt.Println("==========================")
								fmt.Println("Gagal mendaftarn pegawai")
							}
						} else {
							var tmp string
							fmt.Println("==========================")
							fmt.Println("TAMBAH PELANGGAN")
							fmt.Print("Masukkan nomer hp : ")
							fmt.Scanln(&tmp)
							// newPegawai.SetUsername(tmp)
							fmt.Print("Masukkan password : ")
							fmt.Scanln(&tmp)
							// newPegawai.SetPassword(tmp)

						}
					case 3:
						if isAdmin {
							deleteMode := true
							for deleteMode {
								listBarang()
								_, strBarang, err := listBarang()
								if err != nil {
									fmt.Println(err.Error())
								}
								if len(strBarang) > 0 {
									fmt.Println("==========================")
									fmt.Println("HAPUS BARANG")
									fmt.Print(strBarang)
									fmt.Print("Masukkan barcode barang / 0. Kembali halaman: ")
									var inBarcode int
									fmt.Scanln(&inBarcode)
									if inBarcode == 0 {
										deleteMode = !deleteMode
										continue
									}
									isDeleted, err := BarangMenu.Delete(inBarcode)
									if err != nil {
										fmt.Println(err.Error())
									}
									if isDeleted {
										fmt.Println("==========================")
										fmt.Println("berhasil menghapus barang")
									} else {
										fmt.Println("==========================")
										fmt.Println("gagal menghapus barang")
									}
								} else {
									fmt.Println("==========================")
									fmt.Println("Kak", resLogin.GetNama(), "belum memiliki data barang satu pun")
									deleteMode = !deleteMode
								}
							}
						} else {
							editMode := true
							for editMode {
								_, strBarang, err := listBarang()
								if err != nil {
									fmt.Println(err.Error())
								}
								if len(strBarang) > 0 {
									fmt.Println("==========================")
									fmt.Println("EDIT BARANG")
									fmt.Print(strBarang)
									fmt.Print("Masukkan barcode barang / 0. Kembali halaman: ")
									var inBarcode int
									fmt.Scanln(&inBarcode)
									if inBarcode == 0 {
										editMode = !editMode
										continue
									}
								}
							}
						}
					case 9:
						isLogged = !isLogged
					case 0:
						isLogged = !isLogged
					}
				}
			}
		}
	}

}
