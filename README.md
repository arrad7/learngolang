# ProjectGOlang
Aplikasi untuk memenuhi Final Project dari GOlang Course

# How to Run ?
#1.Create The Database with name go_db on mysql server
Buat Database dengan nama go_db atau bisa menyesuaikan dengan kriteria anda
dengan mengubah line ini:

Main.go line 21 'db, err = sql.Open("mysql", "root:arrad@tcp(127.0.0.1)/go_db")'

root = username mysql server

arrad = password mysql server

go_db = database yang digunakan

#2. Create Table 

CREATE TABLE  `go_db`.`users` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `username` varchar(50) DEFAULT NULL,
  `first_name` varchar(200) NOT NULL,
  `last_name` varchar(200) NOT NULL,
  `password` varchar(120) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=latin1;

CREATE TABLE  `go_db`.`article` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `judul` varchar(45) NOT NULL,
  `isi` text NOT NULL,
  `sta` varchar(5) NOT NULL,
  `aktif` varchar(1) NOT NULL DEFAULT 'T',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=latin1;

CREATE TABLE  `go_db`.`pesan` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `nama` varchar(45) NOT NULL,
  `email` varchar(45) NOT NULL,
  `pesan` text NOT NULL,
  `aktif` varchar(1) NOT NULL DEFAULT 'T',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=latin1;

#run
you should import the library first

	go get database/sql

	go get golang.org/x/crypto/bcrypt

	go get github.com/go-sql-driver/mysql
	go get github.com/kataras/go-sessions

### And here we go 
	go run main.go
' Baca keterangan Aplikasi di halaman about
