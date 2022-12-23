package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"tokoku-app-project/barang"
	"tokoku-app-project/config"
	"tokoku-app-project/pegawai"
	"tokoku-app-project/pelanggan"
	"tokoku-app-project/transaksi"
	"tokoku-app-project/transaksibarang"
)

func callClear() { cmd := exec.Command("clear"); cmd.Stdout = os.Stdout; cmd.Run() }

func main() {
	var (
		cfg                     = config.ReadConfig()
		conn                    = config.ConnectSQL(*cfg)
		PegawaiMenu             = pegawai.NewPegawaiMenu(conn)
		BarangMenu              = barang.NewBarangMenu(conn)
		PelangganMenu           = pelanggan.NewPelangganMenu(conn)
		TransaksiMenu           = transaksi.NewTransaksiMenu(conn)
		TransaksiBarangMenu     = transaksibarang.NewTransaksiBarangMenu(conn)
		inputMenu           int = 1
	)
	for inputMenu != 0 {
		fmt.Println("==========================")
		fmt.Println("MENU")
		fmt.Println()
		fmt.Println("1. Login")
		fmt.Println("0. Exit")
		fmt.Println()
		fmt.Print("Pilih menu : ")
		fmt.Scanln(&inputMenu)
		callClear()
		switch inputMenu {
		case 1:
			var inNama, inPassword string
			fmt.Println("==========================")
			fmt.Println("LOGIN")
			fmt.Println()
			fmt.Print("Masukkan username : ")
			fmt.Scanln(&inNama)
			fmt.Print("Masukkan password : ")
			fmt.Scanln(&inPassword)
			callClear()
			fmt.Println("==========================")
			resLogin, err := PegawaiMenu.Login(inNama, inPassword) //mengambil method login dari package pegawai
			if err != nil {
				fmt.Println(err.Error())
			}
			if resLogin.GetID() > 0 {
				fmt.Println("Login sukses, selamat datang", resLogin.GetNama())
				isLogged := true
				var isAdmin bool // false
				if resLogin.GetUsername() == "admin" && inPassword == "admin" {
					isAdmin = !isAdmin
				}
				for isLogged {
					fmt.Println("==========================")
					if isAdmin {
						fmt.Println("1. Tambah Pegawai")
						fmt.Println("2. Hapus Pegawai")
						fmt.Println("3. Hapus Barang")
						fmt.Println("4. Hapus Pelanggan")
						fmt.Println("5. Hapus Transaksi")
					} else {
						fmt.Println("1. Tambah Pelanggan")
						fmt.Println("2. Tambah Barang")
						fmt.Println("3. Edit Barang")
						fmt.Println("4. Transaksi")
						fmt.Println("5. Lihat Transaksi")
						fmt.Println("6. Tambah Stok")
					}
					fmt.Println("9. Log out")
					fmt.Println("0. Exit")
					fmt.Println()
					fmt.Print("Pilih menu : ")
					fmt.Scanln(&inputMenu)
					callClear()
					switch inputMenu {
					case 1:
						///////// menu pegawai
						if isAdmin {
							var newPegawai pegawai.Pegawai
							var tmp string
							reader := bufio.NewReader(os.Stdin)
							fmt.Println("0. Kembali")
							fmt.Println("==========================")
							fmt.Println("TAMBAH PEGAWAI")
							fmt.Println()
							fmt.Print("Masukkan nama : ")
							nama, _ := reader.ReadString('\n')
							nama = nama[:len(nama)-1]
							if nama == "0" {
								callClear()
								continue
							}
							newPegawai.SetNama(nama)
							fmt.Print("Masukkan username : ")
							fmt.Scanln(&tmp)
							if tmp == "0" {
								callClear()
								continue
							}
							newPegawai.SetUsername(tmp)
							fmt.Print("Masukkan password : ")
							fmt.Scanln(&tmp)
							if tmp == "0" {
								callClear()
								continue
							}
							newPegawai.SetPassword(tmp)
							newPegawai.SetIsActive(1)
							callClear()
							fmt.Println("==========================")
							isAdded, err := PegawaiMenu.Register(newPegawai)
							if err != nil {
								fmt.Println(err.Error())
							}
							if isAdded {
								fmt.Println("Sukses menambahkan pegawai")
							} else {
								fmt.Println("Gagal mendaftarn pegawai")
							}
							//////// menu pelanggan
						} else {
							var newPelanggan pelanggan.Pelanggan
							var inHP string
							reader := bufio.NewReader(os.Stdin)
							fmt.Println("0. Kembali")
							fmt.Println("==========================")
							fmt.Println("TAMBAH PELANGGAN")
							fmt.Println()
							fmt.Print("Masukkan nomer hp : ")
							fmt.Scanln(&inHP)
							if inHP == "0" {
								callClear()
								continue
							}
							newPelanggan.SetHP(inHP)
							newPelanggan.SetIDPegawai(resLogin.GetID())
							fmt.Print("Masukkan nama : ")
							nama, _ := reader.ReadString('\n')
							nama = nama[:len(nama)-1]
							if nama == "0" {
								callClear()
								continue
							}
							newPelanggan.SetNamaPelanggan(nama)
							callClear()
							isInserted, err := PelangganMenu.Register(newPelanggan)
							if err != nil {
								fmt.Println(err.Error())
							}
							if isInserted {
								fmt.Println("==========================")
								fmt.Println("Sukses menambahkan pelanggan")
							} else {
								fmt.Println("==========================")
								fmt.Println("Gagal mendaftarn pelanggan")
							}
						}
						/////// menu pegawai
					case 2:
						if isAdmin {
							deleteMode := true
							for deleteMode {
								_, strPegawai, err := PegawaiMenu.Data(resLogin.GetUsername())
								if err != nil {
									fmt.Println(err.Error())
								}
								if len(strPegawai) > 0 {
									fmt.Println("==========================")
									fmt.Println("HAPUS PEGAWAI")
									fmt.Println()
									fmt.Println("ID\t| Nama Pegawai")
									fmt.Println()
									fmt.Print(strPegawai)
									fmt.Println()
									fmt.Print("Masukkan ID pegawai / 0. Kembali: ")
									var inPegawaiID int
									fmt.Scanln(&inPegawaiID)
									callClear()
									fmt.Println("==========================")
									if inPegawaiID == 0 {
										deleteMode = !deleteMode
										callClear()
										continue
									}
									isDeleted, err := PegawaiMenu.Delete(inPegawaiID)
									if err != nil {
										fmt.Println(err.Error())
									}
									if isDeleted {
										fmt.Println("berhasil menghapus pegawai")
									} else {
										fmt.Println("gagal menghapus pegawai")
									}
								} else {
									fmt.Println("==========================")
									fmt.Println("Kak", resLogin.GetNama(), "belum memiliki data pegawai satu pun")
									deleteMode = !deleteMode
								}
							}
							//////menu barang
						} else {
							var newBarang barang.Barang
							var tmp int
							reader := bufio.NewReader(os.Stdin)
							fmt.Println("==========================")
							fmt.Println("TAMBAH BARANG")
							fmt.Println()
							fmt.Print("Masukkan barcode / 0. Kembali : ")
							fmt.Scanln(&tmp)
							if tmp == 0 {
								callClear()
								continue
							}
							newBarang.SetBarcode(tmp)
							newBarang.SetIDPegawai(resLogin.GetID())
							fmt.Print("Masukkan nama barang / 0. Kembali : ")
							nama, _ := reader.ReadString('\n')
							nama = nama[:len(nama)-1]
							if nama == "0" {
								callClear()
								continue
							}
							newBarang.SetNamaBarang(nama)
							fmt.Print("Masukkan stok : ")
							fmt.Scanln(&tmp)
							newBarang.SetStok(tmp)
							fmt.Print("Masukkan harga : ")
							fmt.Scanln(&tmp)
							newBarang.SetHarga(tmp)
							callClear()
							fmt.Println("==========================")
							isAdded, err := BarangMenu.Register(newBarang)
							if err != nil {
								fmt.Println(err.Error())
								callClear()
								continue
							}
							if isAdded {
								fmt.Println("Sukses menambahkan barang")
							} else {
								fmt.Println("Gagal menambahkan barang")
							}
						}
					case 3:
						if isAdmin {
							deleteMode := true
							for deleteMode {
								var inBarcode int
								_, strBarang, err := BarangMenu.Data(inBarcode)
								if err != nil {
									fmt.Println(err.Error())
								}
								if len(strBarang) > 0 {
									fmt.Println("==========================")
									fmt.Println("HAPUS BARANG")
									fmt.Println()
									fmt.Println("Barcode\t| Barang {Stok} [Harga] <Created By>")
									fmt.Println()
									fmt.Print(strBarang)
									fmt.Println()
									fmt.Print("Masukkan barcode / 0. Kembali : ")
									fmt.Scanln(&inBarcode)
									callClear()
									fmt.Println("==========================")
									if inBarcode == 0 {
										deleteMode = !deleteMode
										callClear()
										continue
									}
									isDeleted, err := BarangMenu.Delete(inBarcode)
									if err != nil {
										fmt.Println(err.Error())
									}
									if isDeleted {
										fmt.Println("berhasil menghapus barang")
									} else {
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
								var inBarcode int
								_, strBarang, err := BarangMenu.Data(inBarcode)
								if err != nil {
									fmt.Println(err.Error())
								}
								if len(strBarang) > 0 {
									fmt.Println("==========================")
									fmt.Println("EDIT BARANG")
									fmt.Println()
									fmt.Println("Barcode\t| Barang {Stok} [Harga] <Created By>")
									fmt.Println()
									fmt.Print(strBarang)
									fmt.Println()
									fmt.Print("Masukkan barcode / 0. Kembali : ")
									fmt.Scanln(&inBarcode)
									callClear()
									fmt.Println("==========================")
									if inBarcode == 0 {
										editMode = !editMode
										callClear()
										continue
									}
									arrBarang, strBarang, err := BarangMenu.Data(inBarcode)
									if err != nil {
										fmt.Println(err.Error())
									}
									if len(arrBarang) > 0 {
										idx := strings.Index(strBarang, "<")
										createdBy := strBarang[idx+1 : len(strBarang)-2]
										var inStok, inHarga int
										var upBarang barang.Barang
										reader := bufio.NewReader(os.Stdin)
										fmt.Println("EDIT BARANG")
										fmt.Println()
										fmt.Print("Barcode\t\t:")
										fmt.Println(arrBarang[0].GetBarcode())

										fmt.Print("Nama barang\t:")
										fmt.Println(arrBarang[0].GetNamaBarang())
										fmt.Print("Stok\t\t:")
										fmt.Println(arrBarang[0].GetStok())
										fmt.Print("Harga\t\t:")
										fmt.Println(arrBarang[0].GetHarga())
										fmt.Print("Created by\t:")
										fmt.Println(createdBy)
										fmt.Println()
										fmt.Println("# Kosongkan input jika tidak ingin ada perubahan #")
										fmt.Print("Masukkan nama barang / 0. Kembali : ")
										nama, _ := reader.ReadString('\n')
										nama = nama[:len(nama)-1]
										if nama == "" {
											nama = arrBarang[0].GetNamaBarang()
										} else if nama == "0" {
											callClear()
											continue
										}
										fmt.Print("Masukkan stok barang : ")
										fmt.Scanln(&inStok)
										if inStok == 0 {
											inStok = arrBarang[0].GetStok()
										}
										fmt.Print("Masukkan harga barang : ")
										fmt.Scanln(&inHarga)
										callClear()
										fmt.Println("==========================")
										if inHarga == 0 {
											inHarga = arrBarang[0].GetHarga()
										}
										upBarang.SetBarcode(arrBarang[0].GetBarcode())
										upBarang.SetNamaBarang(nama)
										upBarang.SetStok(inStok)
										upBarang.SetHarga(inHarga)
										isEdited, err := BarangMenu.Update(upBarang)
										if err != nil {
											fmt.Println(err.Error())
										}
										if isEdited {
											fmt.Println("berhasil edit informasi barang")
										} else {
											fmt.Println("gagal edit informasi barang")
										}
									} else {
										fmt.Println("Barang tidak ditemukan")
									}

								}
							}

						}

					case 4:
						////// menu pelanggan
						if isAdmin {
							deleteMode := true
							for deleteMode {
								var nomer_pelanggan string
								_, strPelanggan, err := PelangganMenu.Data(nomer_pelanggan)
								if err != nil {
									fmt.Println(err.Error())
								}
								if len(strPelanggan) > 0 {
									fmt.Println("==========================")
									fmt.Println("HAPUS PELANGGAN")
									fmt.Println()
									fmt.Println("HP\t| Nama Pelanggan <Created By>")
									fmt.Println()
									fmt.Print(strPelanggan)
									fmt.Println()
									fmt.Print("Masukkan nomer hp / 0. Kembali: ")
									var inHP string
									fmt.Scanln(&inHP)
									callClear()
									fmt.Println("==========================")
									if inHP == "0" {
										deleteMode = !deleteMode
										callClear()
										continue
									}
									isDeleted, err := PelangganMenu.Delete(inHP)
									if err != nil {
										fmt.Println(err.Error())
									}
									if isDeleted {
										fmt.Println("==========================")
										fmt.Println("berhasil menghapus pelanggan")
									} else {
										fmt.Println("==========================")
										fmt.Println("gagal menghapus pelanggan")
									}
								} else {
									fmt.Println("==========================")
									fmt.Println("Kak", resLogin.GetNama(), "belum memiliki data pelanggan satu pun")
									deleteMode = !deleteMode
								}
							}
						} else {
							transaksiMode := true
							for transaksiMode {
								var nomer_pelanggan string
								_, strPelanggan, err := PelangganMenu.Data(nomer_pelanggan)
								if err != nil {
									fmt.Println(err.Error())
									fmt.Println("harap masukkan data dengan benar")

								}
								if len(strPelanggan) > 0 {
									var newTransaksi transaksi.Transaksi
									fmt.Println("==========================")
									fmt.Println("TRANSAKSI")
									fmt.Println()
									fmt.Println("Pilih pelanggan")
									fmt.Println()
									fmt.Println("HP\t| Nama Pelanggan <Created By>")
									fmt.Println()
									fmt.Print(strPelanggan)
									fmt.Println()
									fmt.Print("Masukkan no. hp / 0. Kembali : ")
									var inHP string
									fmt.Scanln(&inHP)
									newTransaksi.SetHP(inHP)
									newTransaksi.SetIDPegawai(resLogin.GetID())
									callClear()
									if inHP == "0" {
										callClear()
										transaksiMode = !transaksiMode
										continue
									}
									idInserted, err := TransaksiMenu.Insert(newTransaksi)
									if err != nil {
										fmt.Println(err.Error())
									}
									if idInserted > 0 {
										arrTransaksi, _, err := TransaksiMenu.Data(idInserted)
										if err != nil {
											fmt.Println(err.Error())
										}
										if len(arrTransaksi) > 0 {
											sellMode := true
											for sellMode {
												_, strTransaksiBarang, err := TransaksiBarangMenu.Data(idInserted)
												if err != nil {
													fmt.Println(err.Error())
												}
												var inBarcode, inJumlah int
												fmt.Println("==========================")
												fmt.Println("TRANSAKSI")
												fmt.Println()
												fmt.Print("No. Transaksi\t:")
												fmt.Println(idInserted)
												fmt.Print("Waktu\t\t:")
												fmt.Println(arrTransaksi[0].GetTanggal())
												fmt.Print("Kasir\t\t:")
												fmt.Printf("%d <%s>\n", arrTransaksi[0].GetIDPegawai(), arrTransaksi[0].GetNamaPegawai())
												fmt.Print("Pelanggan\t:")
												fmt.Printf("%s <%s>\n", arrTransaksi[0].GetHP(), arrTransaksi[0].GetNamaPelanggan())
												fmt.Println()
												if len(strTransaksiBarang) > 0 {
													fmt.Print(strTransaksiBarang)
												} else {
													fmt.Println("Belum ada barang yang dipilih")
												}
												fmt.Println()
												amount, err := TransaksiBarangMenu.Amount(arrTransaksi[0].GetID())
												if err != nil {
													fmt.Println(err.Error())
												}
												fmt.Println("Total bayar : ", amount)
												_, strBarang, err := BarangMenu.Data(inBarcode)
												if err != nil {
													fmt.Println(err.Error())
												}
												if len(strBarang) > 0 {
													fmt.Println()
													fmt.Println("==========================")
													fmt.Println("Pilih BARANG")
													fmt.Println()
													fmt.Println("Barcode\t| Barang {Stok} [Harga] <Created By>")
													fmt.Println()
													fmt.Print(strBarang)
													fmt.Println()
													fmt.Print("Masukkan barcode / 0. Kembali : ")
													fmt.Scanln(&inBarcode)
													if inBarcode == 0 {
														callClear()
														sellMode = !sellMode
														transaksiMode = !transaksiMode
														continue
													}
													fmt.Print("Masukkan jumlah beli / 0. Kembali : ")
													fmt.Scanln(&inJumlah)
													callClear()
													if inJumlah == 0 {
														callClear()
														continue
													}
													isSell, err := BarangMenu.Sell(inBarcode, inJumlah)
													if err != nil {
														fmt.Println(err.Error())
													}
													if isSell {
														var newTransaksiBarang transaksibarang.TransaksiBarang
														newTransaksiBarang.SetIDTransaksi(idInserted)
														newTransaksiBarang.SetBarcode(inBarcode)
														newTransaksiBarang.SetJumlah(inJumlah)
														isInserted, err := TransaksiBarangMenu.Insert(newTransaksiBarang)
														if err != nil {
															fmt.Println(err.Error())
														}
														if isInserted {
															fmt.Println("Berhasil memasukkan barang ke transaksi")
														} else {
															fmt.Println("Gagal memasukkan barang ke transaksi")
														}
													} else {
														fmt.Println("Gagal memasukkan barang ke transaksi, barang tidak ada / stok kurang")
													}
												}

											}
										}

									}

								}
							}

						}
					case 5:
						if isAdmin {
							deleteMode := true
							for deleteMode {
								var inID int
								arrTransaksi, _, err := TransaksiMenu.Data(inID)
								if err != nil {
									fmt.Println(err.Error())
								}
								if len(arrTransaksi) > 0 {
									fmt.Println("==========================")
									fmt.Println("HAPUS TRANSAKSI")
									fmt.Println()
									fmt.Println("Pilih BARANG")
									fmt.Println()
									fmt.Println("ID Transaksi\t| Waktu\t\t\t| Kasir\t\t\t| Amount\t|\t Pelanggan")
									fmt.Println()
									for _, v := range arrTransaksi {
										amount, _ := TransaksiBarangMenu.Amount(v.GetID())
										tmp := amount
										count := 0
										for tmp > 0 {
											tmp = tmp / 10
											count++
										}
										count /= 6
										tabAmount := strings.Repeat("\t", 2)
										tabAmount = tabAmount[:len(tabAmount)-count]
										fmt.Printf("%d\t\t| %s\t| %s\t| %d%s|%s\n", v.GetID(), v.GetTanggal(), v.GetNamaPegawai(), amount, tabAmount, v.GetNamaPelanggan())
									}
									fmt.Println()
									fmt.Println("Pilih ID transaksi untuk detail lebih lanjut")
									fmt.Print("Masukkan ID Transaksi / 0. Kembali : ")
									fmt.Scanln(&inID)
									callClear()
									if inID == 0 {
										callClear()
										deleteMode = !deleteMode
										continue
									}
									arrTransaksi, _, err := TransaksiMenu.Data(inID)
									if err != nil {
										fmt.Println(err.Error())
									}
									if len(arrTransaksi) > 0 {
										_, strTransaksiBarang, err := TransaksiBarangMenu.Data(inID)
										if err != nil {
											fmt.Println(err.Error())
										}
										var inHapus string
										inChoice := true
										amount, _ := TransaksiBarangMenu.Amount(inID)
										for inChoice {
											fmt.Println("==========================")
											fmt.Println("HAPUS TRANSAKSI")
											fmt.Println()
											fmt.Print("No. Transaksi\t:")
											fmt.Println(inID)
											fmt.Print("Waktu\t\t:")
											fmt.Println(arrTransaksi[0].GetTanggal())
											fmt.Print("Kasir\t\t:")
											fmt.Printf("%d <%s>\n", arrTransaksi[0].GetIDPegawai(), arrTransaksi[0].GetNamaPegawai())
											fmt.Print("Pelanggan\t:")
											fmt.Printf("%s <%s>\n", arrTransaksi[0].GetHP(), arrTransaksi[0].GetNamaPelanggan())
											fmt.Println()
											if len(strTransaksiBarang) > 0 {
												fmt.Print(strTransaksiBarang)
											} else {
												fmt.Println("Belum ada barang yang dipilih")
											}

											if err != nil {
												fmt.Println(err.Error())
											}
											fmt.Println()
											fmt.Println("Total bayar : ", amount)
											fmt.Println()
											fmt.Print("Yakin ingin menghapus? (ya/tidak) : ")
											fmt.Scanln(&inHapus)
											inHapus = strings.ToLower(inHapus)
											callClear()
											fmt.Println("==========================")
											if inHapus == "ya" {
												isDeleted, err := TransaksiMenu.Delete(inID)
												if err != nil {
													fmt.Println(err.Error())
												}
												if isDeleted {
													fmt.Println("berhasil menghapus transaksi")
													inChoice = !inChoice
												} else {
													fmt.Println("gagal menghapus transaksi")
													inChoice = !inChoice
												}
											} else if inHapus == "tidak" {
												callClear()
												inChoice = !inChoice
												continue
											}
										}
									} else {
										fmt.Println("==========================")
										fmt.Println("Data tidak ditemukan")
									}
								}

							}
						} else {
							viewMode := true
							for viewMode {
								var inID int
								arrTransaksi, _, err := TransaksiMenu.Data(inID)
								if err != nil {
									fmt.Println(err.Error())
								}
								if len(arrTransaksi) > 0 {
									fmt.Println("==========================")
									fmt.Println("LIHAT TRANSAKSI")
									fmt.Println()
									fmt.Println("Pilih BARANG")
									fmt.Println()
									fmt.Println("ID Transaksi\t| Waktu\t\t\t| Kasir\t\t\t| Amount\t|\t Pelanggan")
									fmt.Println()
									for _, v := range arrTransaksi {
										amount, _ := TransaksiBarangMenu.Amount(v.GetID())
										tmp := amount
										count := 0
										for tmp > 0 {
											tmp = tmp / 10
											count++
										}
										count /= 6
										tabAmount := strings.Repeat("\t", 2)
										tabAmount = tabAmount[:len(tabAmount)-count]
										fmt.Printf("%d\t\t| %s\t| %s\t| %d%s|%s\n", v.GetID(), v.GetTanggal(), v.GetNamaPegawai(), amount, tabAmount, v.GetNamaPelanggan())
									}
									fmt.Println()
									fmt.Println("Pilih ID transaksi untuk detail lebih lanjut")
									fmt.Print("Masukkan ID Transaksi / 0. Kembali : ")
									fmt.Scanln(&inID)
									callClear()
									if inID == 0 {
										callClear()
										viewMode = !viewMode
										continue
									}
									arrTransaksi, _, err := TransaksiMenu.Data(inID)
									if err != nil {
										fmt.Println(err.Error())
									}
									if len(arrTransaksi) > 0 {
										_, strTransaksiBarang, err := TransaksiBarangMenu.Data(inID)
										if err != nil {
											fmt.Println(err.Error())
										}
										var tmp int
										detailMode := true
										amount, _ := TransaksiBarangMenu.Amount(inID)
										for detailMode {
											fmt.Println("==========================")
											fmt.Println("HAPUS TRANSAKSI")
											fmt.Println()
											fmt.Print("No. Transaksi\t:")
											fmt.Println(inID)
											fmt.Print("Waktu\t\t:")
											fmt.Println(arrTransaksi[0].GetTanggal())
											fmt.Print("Kasir\t\t:")
											fmt.Printf("%d <%s>\n", arrTransaksi[0].GetIDPegawai(), arrTransaksi[0].GetNamaPegawai())
											fmt.Print("Pelanggan\t:")
											fmt.Printf("%s <%s>\n", arrTransaksi[0].GetHP(), arrTransaksi[0].GetNamaPelanggan())
											fmt.Println()
											if len(strTransaksiBarang) > 0 {
												fmt.Print(strTransaksiBarang)
											} else {
												fmt.Println("Belum ada barang yang dipilih")
											}

											if err != nil {
												fmt.Println(err.Error())
											}
											fmt.Println()
											fmt.Println("Total bayar : ", amount)
											fmt.Println()
											fmt.Println("==========================")
											fmt.Println("MENU")
											fmt.Println()
											fmt.Println("9. Kembali")
											fmt.Println("0. Halaman")
											fmt.Println()
											fmt.Print("Pilih menu : ")
											fmt.Scanln(&tmp)
											callClear()
											if tmp == 9 {
												detailMode = !detailMode
												continue
											}
											if tmp == 0 {
												detailMode = !detailMode
												viewMode = !viewMode
												continue
											}
										}
									} else {
										fmt.Println("==========================")
										fmt.Println("Data tidak ditemukan")
									}
								}
							}
						}
					case 6:
						increaseMode := true
						for increaseMode {
							var inBarcode int
							_, strBarang, err := BarangMenu.Data(inBarcode)
							if err != nil {
								fmt.Println(err.Error())
							}
							if len(strBarang) > 0 {
								fmt.Println("==========================")
								fmt.Println("TAMBAH STOK BARANG")
								fmt.Println()
								fmt.Println("Barcode\t| Barang {Stok} [Harga] <Created By>")
								fmt.Println()
								fmt.Print(strBarang)
								fmt.Println()
								fmt.Print("Masukkan barcode / 0. Kembali : ")
								fmt.Scanln(&inBarcode)
								fmt.Println("==========================")
								if inBarcode == 0 {
									increaseMode = !increaseMode
									callClear()
									continue
								}
								callClear()
								stokMode := true
								for stokMode {
									arrBarang, strBarang, err := BarangMenu.Data(inBarcode)
									if err != nil {
										fmt.Println(err.Error())
									}
									if len(arrBarang) > 0 {
										idx := strings.Index(strBarang, "<")
										createdBy := strBarang[idx+1 : len(strBarang)-2]
										var inStok int
										fmt.Println("==========================")
										fmt.Println("TAMBAH STOK BARANG")
										fmt.Println()
										fmt.Print("Barcode\t\t:")
										fmt.Println(arrBarang[0].GetBarcode())

										fmt.Print("Nama barang\t:")
										fmt.Println(arrBarang[0].GetNamaBarang())
										fmt.Print("Stok\t\t:")
										fmt.Println(arrBarang[0].GetStok())
										fmt.Print("Harga\t\t:")
										fmt.Println(arrBarang[0].GetHarga())
										fmt.Print("Created by\t:")
										fmt.Println(createdBy)
										fmt.Println()
										fmt.Print("Masukkan stok barang / 0. Kembali: ")
										fmt.Scanln(&inStok)
										if inStok == 0 {
											stokMode = !stokMode
											callClear()
											continue
										}
										callClear()
										fmt.Println("==========================")
										isEdited, err := BarangMenu.IncreaseStok(inStok, arrBarang[0].GetBarcode())
										if err != nil {
											fmt.Println(err.Error())
										}
										if isEdited {
											fmt.Println("berhasil tambah stok barang")
										} else {
											fmt.Println("gagal tambah stok barang")
										}
									} else {
										fmt.Println("==========================")
										fmt.Println("Barang tidak ditemukan")
										stokMode = !stokMode
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
