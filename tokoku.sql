use tokoku;

DROP table transaksi_barang;
DROP table transaksi;
DROP table barang;
DROP table pelanggan;
DROP table pegawai;





CREATE table pegawai(
	id int auto_increment primary key,
	username varchar(255) NOT NULL,
	password varchar(255) NOT NULL
);

CREATE table pelanggan(
	hp varchar(25) NOT NULL primary key,
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
	hp varchar(25) NOT NULL,
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
desc pelanggan;
desc pegawai;
desc barang;

ALTER TABLE tokoku.pegawai DROP COLUMN isActive;
ALTER TABLE tokoku.pegawai ADD isActive BOOL DEFAULT false NOT NULL;
ALTER TABLE tokoku.pegawai ADD nama varchar(255) NOT NULL;
ALTER TABLE tokoku.pegawai MODIFY isActive BOOL DEFAULT false NOT NULL;
ALTER TABLE tokoku.barang ADD harga int NOT NULL;



-- INSERT 
INSERT into pegawai (username, password, nama, isActive) values ('admin', 'admin', 'Admin', '1');

-- UPDATE 
UPDATE pegawai p
SET nama = 'Admin'
WHERE p.id = 1;

UPDATE pegawai p
SET isActive = 1
WHERE p.id = 1;

-- DELETE 
DELETE FROM pegawai WHERE id = 2;


-- SELECT 
SELECT * FROM pegawai;
SELECT * FROM barang;
SELECT * FROM pelanggan p ;
SELECT barcode,id_pegawai,nama,stok,harga
FROM barang;




