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
desc transaksi;
desc transaksi_barang ;

ALTER TABLE tokoku.pegawai ADD isActive BOOL DEFAULT false NOT NULL;
ALTER TABLE tokoku.pegawai ADD nama varchar(255) NOT NULL;
ALTER TABLE tokoku.pegawai MODIFY isActive BOOL DEFAULT false NOT NULL;
ALTER TABLE tokoku.pelanggan ADD isActive BOOL DEFAULT false NOT NULL;
ALTER TABLE tokoku.barang ADD harga int NOT NULL;
ALTER TABLE tokoku.barang ADD isActive BOOL DEFAULT false NOT NULL;
ALTER TABLE tokoku.transaksi_barang ADD jumlah int NOT NULL;


-- INSERT 
INSERT into pegawai (username, password, nama, isActive) values ('admin', 'admin', 'Admin', '1');
INSERT into transaksi_barang (id_transaksi, barcode, jumlah) values (3, 112, 2);
INSERT into transaksi_barang (id_transaksi, barcode, jumlah) values (3, 113, 3);
INSERT into transaksi_barang (id_transaksi, barcode, jumlah) values (3, 777, 4);
INSERT into transaksi_barang (id_transaksi, barcode, jumlah) values (3, 887, 5);
INSERT into transaksi_barang (id_transaksi, barcode, jumlah) values (3, 899, 6);
INSERT into transaksi_barang (id_transaksi, barcode, jumlah) values (1, 899, 6);
INSERT into transaksi_barang (id_transaksi, barcode, jumlah) values (1, 887, 6);
-- UPDATE 
UPDATE pegawai p
SET nama = 'Admin'
WHERE p.id = 1;

UPDATE pegawai p
SET isActive = 1
WHERE p.id = 1;

UPDATE pelanggan  p
SET isActive = 1
WHERE p.hp = 08987;


UPDATE barang
SET stok = stok - 100 
WHERE barcode = 112 and stok > 100;

UPDATE barang  b
SET b.isActive = 1
WHERE b.barcode = 899;

-- DELETE 
DELETE FROM pegawai WHERE id = 2;
DELETE FROM barang
WHERE barcode = 0;


-- SELECT 
SELECT * FROM pegawai;

SELECT * FROM barang;

SELECT b.barcode ,b.id_pegawai ,b.nama ,b.stok ,b.harga ,p.nama 'nama pegawai'
FROM barang b
JOIN pegawai p ON p.id = b.id_pegawai;

SELECT * FROM pelanggan p ;

SELECT p.hp, p.nama, p.id_pegawai, p2.nama  
		FROM pelanggan p  
		JOIN pegawai p2 ON p2.id = p.id_pegawai
		WHERE p.hp = 0811111
		AND p.isActive = 1;

SELECT p.hp, p.nama, p.id_pegawai, p2.nama  
FROM pelanggan p  
JOIN pegawai p2 ON p2.id = p.id_pegawai 
WHERE p.isActive = 1;


SELECT barcode,id_pegawai,nama,stok,harga
FROM barang;

SELECT * FROM transaksi t;

SELECT * FROM transaksi_barang tb ;

SELECT b.nama, tb.jumlah, b.harga, tb.jumlah * b.harga 
FROM barang b 
JOIN transaksi_barang tb ON tb.barcode = b.barcode
WHERE tb.id_transaksi = 3;

SELECT t.id, t.tanggal, t.id_pegawai, p.nama
FROM transaksi t 
JOIN pegawai p ON p.id = t.id_pegawai 

SELECT b.nama, tb.jumlah, b.harga, tb.jumlah * b.harga 
FROM barang b
JOIN transaksi_barang tb ON tb.barcode = b.barcode 
WHERE tb.id_transaksi = 1;

SELECT t.id ,t.tanggal ,t.id_pegawai ,p.nama "Nama Pegawai" , t.hp ,p2.nama as "Nama Pelanggan"
FROM transaksi t 
JOIN pegawai p ON p.id = t.id_pegawai 
JOIN pelanggan p2 ON p2.hp = t.hp;









