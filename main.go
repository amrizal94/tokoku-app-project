package main

import (
	"fmt"
	"tokoku-app-project/config"
	"tokoku-app-project/pegawai"
)

func main() {
	var (
		cfg             = config.ReadConfig()
		conn            = config.ConnectSQL(*cfg)
		PegawaiMenu     = pegawai.NewPegawaiMenu(conn)
		inputMenu   int = 1
	)
	for inputMenu != 0 {
		fmt.Println("==========================")
		fmt.Println("1. Register")
		fmt.Println("2. Login")
		fmt.Println("0. Exit")
		fmt.Scanln(&inputMenu)
		switch inputMenu {
		case 1:
		case 2:
			var inNama, inPassword string
			fmt.Println("==========================")
			fmt.Println("LOGIN")
			fmt.Print("Masukkan nama : ")
			fmt.Scanln(&inNama)
			fmt.Print("Masukkan password : ")
			fmt.Scanln(&inPassword)
			resLogin, err := PegawaiMenu.Login(inNama, inPassword)
			if err != nil {
				fmt.Println(err.Error())
			}
			if resLogin.GetID() > 0 {
				fmt.Println("==========================")
				fmt.Println("Login sukses")
				isLogged := true
				var isAdmin bool
				if resLogin.GetUsername() == "admin" && inPassword == "admin" {
					isAdmin = !isAdmin
				}
				for isLogged {
					fmt.Println("==========================")
					if isAdmin {
						fmt.Println("1. Tambah Pegawai")
					} else {
						fmt.Println("1. Tambah Pelanggan")
					}
					fmt.Println("9. Log out")
					fmt.Println("0. Exit")
					fmt.Scanln(&inputMenu)
					switch inputMenu {
					case 1:
						if isAdmin {
							var newPegawai pegawai.Pegawai
							var tmp string
							fmt.Println("==========================")
							fmt.Println("TAMBAH PEGAWAI")
							fmt.Print("Masukkan username : ")
							fmt.Scanln(&tmp)
							newPegawai.SetUsername(tmp)
							fmt.Print("Masukkan password : ")
							fmt.Scanln(&tmp)
							newPegawai.SetPassword(tmp)
							isAdded, err := PegawaiMenu.Register(newPegawai)
							if err != nil {
								fmt.Println(err.Error())
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
							fmt.Println("TAMBAH PEGAWAI")
							fmt.Print("Masukkan username : ")
							fmt.Scanln(&tmp)
							// newPegawai.SetUsername(tmp)
							fmt.Print("Masukkan password : ")
							fmt.Scanln(&tmp)
							// newPegawai.SetPassword(tmp)
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
