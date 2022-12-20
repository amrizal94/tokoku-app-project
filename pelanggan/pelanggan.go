package pelanggan

type Pelanggan struct {
	hp         string
	id_pegawai int
	nama       string
}

func (p *Pelanggan) SetHP(newHP string) {
	p.hp = newHP
}
func (p *Pelanggan) SetIDPegawai(newIDPegawai int) {
	p.id_pegawai = newIDPegawai
}
func (p *Pelanggan) SetName(newNama string) {
	p.nama = newNama
}

func (p *Pelanggan) GetHP() string {
	return p.hp
}
func (p *Pelanggan) GetIDPegawai() int {
	return p.id_pegawai
}
func (p *Pelanggan) GetName() string {
	return p.nama
}
