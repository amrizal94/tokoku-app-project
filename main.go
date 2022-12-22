package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"tokoku-app-project/barang"
	"tokoku-app-project/config"
	"tokoku-app-project/pegawai"
	"tokoku-app-project/pelanggan"
	"tokoku-app-project/transaksi"
	"tokoku-app-project/transaksibarang"
)

var (
	cfg                 = config.ReadConfig()
	conn                = config.ConnectSQL(*cfg)
	PegawaiMenu         = pegawai.NewPegawaiMenu(conn)
	BarangMenu          = barang.NewBarangMenu(conn)
	PelangganMenu       = pelanggan.NewPelangganMenu(conn)
	TransaksiMenu       = transaksi.NewTransaksiMenu(conn)
	TransaksiBarangMenu = transaksibarang.NewTransaksiBarangMenu(conn)
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

func listTransaksi(id int) ([]transaksi.Transaksi, string, error) {
	arrTransaksi, err := TransaksiMenu.Select(id)

	var strTransaksi string
	if err != nil {
		fmt.Println(err.Error())
		return arrTransaksi, strTransaksi, err
	}
	for _, v := range arrTransaksi {
		arrPegawai, _ := PegawaiMenu.Select(v.GetIDPegawai(), 0)
		strTransaksi += fmt.Sprintf("ID: %d <%s>\n", v.GetID(), arrPegawai[0].GetNama())
	}
	return arrTransaksi, strTransaksi, err
}

func listTransaksiBarang(id int) ([]transaksibarang.TransaksiBarang, string, error) {
	arrTransaksiBarang, err := TransaksiBarangMenu.Select(id)

	var strTransaksiBarang string
	if err != nil {
		fmt.Println(err.Error())
		return arrTransaksiBarang, strTransaksiBarang, err
	}
	for _, v := range arrTransaksiBarang {
		strTransaksiBarang += fmt.Sprintf("%s %d x %d %d\n", v.GetNama(), v.GetJumlah(), v.GetHarga(), v.GetTotal())
	}
	return arrTransaksiBarang, strTransaksiBarang, err
}

func main() {
	var (
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
						fmt.Println("4. Hapus Pelanggan")

					} else {
						fmt.Println("1. Tambah Pelanggan")
						fmt.Println("2. Tambah Barang")
						fmt.Println("3. Edit Barang")
						fmt.Println("4. Transaksi")
					}
					fmt.Println("9. Log out")
					fmt.Println("0. Exit")
					fmt.Scanln(&inputMenu)
					switch inputMenu {
					case 1:
						///////// menu pegawai
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
							newPegawai.SetIsActive(1)
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
							//////// menu pelanggan
						} else {
							var newPelanggan pelanggan.Pelanggan
							var inHP string
							reader := bufio.NewReader(os.Stdin)
							fmt.Println("==========================")
							fmt.Println("TAMBAH PELANGGAN")
							fmt.Print("Masukkan nomer hp : ")
							fmt.Scanln(&inHP)
							newPelanggan.SetHP(inHP)
							newPelanggan.SetIDPegawai(resLogin.GetID())
							fmt.Print("Masukkan nama : ")
							nama, _ := reader.ReadString('\n')
							nama = nama[:len(nama)-1]
							newPelanggan.SetNamaPelanggan(nama)
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
								_, strPegawai, err := listPegawai(0, resLogin.GetID())
								if err != nil {
									fmt.Println(err.Error())
								}
								if len(strPegawai) > 0 {
									fmt.Println("==========================")
									fmt.Println("HAPUS PEGAWAI")
									fmt.Print(strPegawai)
									fmt.Print("Masukkan ID pegawai / 0. Kembali halaman: ")
									var inPegawaiID int
									fmt.Scanln(&inPegawaiID)
									if inPegawaiID == 0 {
										deleteMode = !deleteMode
										continue
									}
									isDeleted, err := PegawaiMenu.Delete(inPegawaiID, 0)
									if err != nil {
										fmt.Println(err.Error())
									}
									if isDeleted {
										fmt.Println("==========================")
										fmt.Println("berhasil menghapus kegiatan")
									} else {
										fmt.Println("==========================")
										fmt.Println("gagal menghapus kegiatan")
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
							fmt.Print("Masukkan barcode : ")
							fmt.Scanln(&tmp)
							newBarang.SetBarcode(tmp)
							newBarang.SetIDPegawai(resLogin.GetID())
							fmt.Print("Masukkan nama barang : ")
							nama, _ := reader.ReadString('\n')
							nama = nama[:len(nama)-1]
							newBarang.SetNamaBarang(nama)
							fmt.Print("Masukkan stok : ")
							fmt.Scanln(&tmp)
							newBarang.SetStok(tmp)
							fmt.Print("Masukkan harga : ")
							fmt.Scanln(&tmp)
							newBarang.SetHarga(tmp)
							fmt.Println("==========================")
							isAdded, err := BarangMenu.Register(newBarang)
							if err != nil {
								fmt.Println(err.Error())
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
									fmt.Print(strBarang)
									fmt.Print("Masukkan barcode / 0. Kembali halaman: ")
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
								var inBarcode int
								_, strBarang, err := BarangMenu.Data(inBarcode)
								if err != nil {
									fmt.Println(err.Error())
								}
								if len(strBarang) > 0 {
									fmt.Println("==========================")
									fmt.Println("EDIT BARANG")
									fmt.Print(strBarang)
									fmt.Print("Masukkan barcode / 0. Kembali halaman: ")
									fmt.Scanln(&inBarcode)
									if inBarcode == 0 {
										editMode = !editMode
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
										reader := bufio.NewReader(os.Stdin)
										fmt.Println("==========================")
										fmt.Println("EDIT BARANG")
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
										fmt.Println("# Kosongkan input jika tidak ingin ada perubahan #")
										fmt.Print("Masukkan nama barang : ")
										nama, _ := reader.ReadString('\n')
										nama = nama[:len(nama)-1]
										if nama == "" {
											nama = arrBarang[0].GetNamaBarang()
										}
										fmt.Print("Masukkan stok barang : ")
										fmt.Scanln(&inStok)
										if inStok == 0 {
											inStok = arrBarang[0].GetStok()
										}
										fmt.Print("Masukkan harga barang : ")
										fmt.Scanln(&inHarga)
										if inHarga == 0 {
											inHarga = arrBarang[0].GetHarga()
										}
										isEdited, err := BarangMenu.Update(arrBarang[0].GetBarcode(), nama, inStok, inHarga)
										if err != nil {
											fmt.Println(err.Error())
										}
										if isEdited {
											fmt.Println("==========================")
											fmt.Println("berhasil edit informasi barang")
										} else {
											fmt.Println("==========================")
											fmt.Println("berhasil edit informasi barang")
										}
									} else {
										fmt.Println("==========================")
										fmt.Println("Kak", resLogin.GetNama(), "belum memiliki data barang satu pun")
										editMode = !editMode
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
									fmt.Print(strPelanggan)
									fmt.Print("Masukkan no. hp / 0. Kembali halaman: ")
									var inHP string
									fmt.Scanln(&inHP)
									if inHP == "0" {
										deleteMode = !deleteMode
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
									fmt.Println("Pilih pelanggan")
									fmt.Print(strPelanggan)
									fmt.Print("Masukkan no. hp / 0. Kembali halaman: ")
									var inHP string
									fmt.Scanln(&inHP)
									newTransaksi.SetHP(inHP)
									newTransaksi.SetIDPegawai(resLogin.GetID())
									if inHP == "0" {
										transaksiMode = !transaksiMode
										continue
									}
									idInserted, err := TransaksiMenu.Insert(newTransaksi)
									if err != nil {
										fmt.Println(err.Error())
									}
									if idInserted > 0 {
										arrTransaksi, strTransaksi, err := listTransaksi(idInserted)
										if err != nil {
											fmt.Println(err.Error())
										}
										if len(arrTransaksi) > 0 {
											idx := strings.Index(strTransaksi, "<")
											kasir := strTransaksi[idx+1 : len(strTransaksi)-2]

											sellMode := true
											for sellMode {
												_, strTransaksiBarang, err := listTransaksiBarang(idInserted)
												if err != nil {
													fmt.Println(err.Error())
												}
												var inBarcode, inJumlah int
												fmt.Println("==========================")
												fmt.Println("TRANSAKSI")
												fmt.Print("No. Transaksi\t:")
												fmt.Println(idInserted)
												fmt.Print("Waktu\t\t:")
												fmt.Println(arrTransaksi[0].GetTanggal())
												fmt.Print("Kasir\t\t:")
												fmt.Printf("%d <%s>\n", arrTransaksi[0].GetIDPegawai(), kasir)
												if len(strTransaksiBarang) > 0 {
													fmt.Print(strTransaksiBarang)
												} else {
													fmt.Println("Belum ada barang yang dipilih")
												}
												_, strBarang, err := BarangMenu.Data(inBarcode)
												if err != nil {
													fmt.Println(err.Error())
												}
												if len(strBarang) > 0 {
													fmt.Println("==========================")
													fmt.Println("Pilih BARANG")
													fmt.Print(strBarang)
													fmt.Println("==========================")
													fmt.Print("Masukkan barcode / 0. Kembali halaman: ")
													fmt.Scanln(&inBarcode)
													if inBarcode == 0 {
														sellMode = !sellMode
														transaksiMode = !transaksiMode
														continue
													}
													fmt.Print("Masukkan jumlah beli / 0. Kembali halaman: ")
													fmt.Scanln(&inJumlah)
													if inJumlah == 0 {
														sellMode = !sellMode
														transaksiMode = !transaksiMode
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
														fmt.Println("Gagal memasukkan barang ke transaksi, stok kurang")
													}
												}

											}
										}

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
