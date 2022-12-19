use tokoku;

DROP table pegawai;
DROP table pelanggan;
DROP table barang;
DROP table transaksi;
DROP table transaksi_barang;

CREATE table pegawai(
	id int auto_increment primary key,
	username varchar(255) NOT NULL,
	password varchar(255) NOT NULL
);

CREATE table pelanggan(
	hp int NOT NULL primary key,
	id_pegawai int NOT NULL,
	nama varchar(255) NOT NULL,
	constraint fk_pegawai_pelanggan foreign key (id_pegawai) references pegawai(id)
);

CREATE table barang(
	barcode int NOT NULL,
	id_pegawai int NOT NULL,
	nama varchar(255) NOT NULL,
	stok int NOT NULL,
	primary key(barcode),
	constraint fk_pegawai_barang foreign key (id_pegawai) references pegawai(id)
);

CREATE table transaksi(
	id int auto_increment,
	id_pegawai int NOT NULL,
	hp int NOT NULL,
	tanggal datetime,
	primary key(id),
	FOREIGN KEY(id_pegawai) REFERENCES pegawai(id),
  	FOREIGN KEY(hp) REFERENCES pelanggan(hp)
);

CREATE table transaksi_barang(
	id_transaksi int NOT NULL,
	barcode int NOT NULL,
	primary key(id_transaksi, barcode),
	FOREIGN KEY(id_transaksi) REFERENCES transaksi(id),
  	FOREIGN KEY(barcode) REFERENCES barang(barcode)
);

-- INSERT 
INSERT into pegawai (username, password) values ('admin', 'admin');



