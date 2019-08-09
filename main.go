package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/kataras/go-sessions"
	"golang.org/x/crypto/bcrypt"
)

var db *sql.DB
var err error

func connect_db() {
	db, err = sql.Open("mysql", "root:arrad@tcp(127.0.0.1)/go_db")

	if err != nil {
		log.Fatalln(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalln(err)
	}
}

type user struct {
	ID        int
	Username  string
	FirstName string
	LastName  string
	Password  string
}

// type Artikel struct {
// 	Artikel []Article
// }
type Article struct {
	ID    int
	Judul string
	Isi   string
	Sta   string
}
type Pesan struct {
	ID       int
	Nama     string
	Email    string
	Isipesan string
}
type ResponsePesan struct {
	Pesana []Pesan
	User   string
}
type ResponseArticle struct {
	Articles []Article
	User     string
}

func routes() {
	http.HandleFunc("/", home)
	http.HandleFunc("/register", register)
	http.HandleFunc("/login", login)
	http.HandleFunc("/logout", logout)
	http.HandleFunc("/about", about)
	http.HandleFunc("/articleku", articleku)
	http.HandleFunc("/contact", contact)
	http.HandleFunc("/update", update)
	http.HandleFunc("/updatestatrue", updatestatrue)
	http.HandleFunc("/updatestafalse", updatestafalse)
	http.HandleFunc("/updateaktif", updateaktif)
	http.HandleFunc("/pesan", pesan)
}

func main() {
	connect_db()
	routes()

	defer db.Close()
	http.Handle("/static/",
		http.StripPrefix("/static/",
			http.FileServer(http.Dir("views/assets"))))
	fmt.Println("Server running on port :8000")
	http.ListenAndServe(":8000", nil)
}

func checkErr(w http.ResponseWriter, r *http.Request, err error) bool {
	if err != nil {

		fmt.Println(r.Host + r.URL.Path)

		http.Redirect(w, r, r.Host+r.URL.Path, 301)
		return false
	}
	return true
}

func QueryUser(username string) user {
	var users = user{}
	err = db.QueryRow(`
		SELECT id, 
		username, 
		first_name, 
		last_name, 
		password 
		FROM users WHERE username=?
		`, username).
		Scan(
			&users.ID,
			&users.Username,
			&users.FirstName,
			&users.LastName,
			&users.Password,
		)
	return users
}

func home(w http.ResponseWriter, r *http.Request) {
	session := sessions.Start(w, r)

	var t, err = template.ParseFiles("views/home.html")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	rows, err := db.Query("select id,judul,isi,sta from `article` where sta like 'T' and aktif like 'T' ")
	if err != nil {
		log.Printf("error query:%+v\n", err)
		return
	}
	defer rows.Close()
	var datas []Article
	var datan Article
	for rows.Next() {
		datan = Article{}
		err = rows.Scan(&datan.ID, &datan.Judul, &datan.Isi, &datan.Sta)
		if err != nil {
			log.Printf("scan error: %+v\n", err)
			continue
		}
		datas = append(datas, datan)
	}
	response := ResponseArticle{
		Articles: datas,
		User:     session.GetString("name"),
	}
	err = t.Execute(w, response)
	if err != nil {
		log.Printf("execute template err: %+v\n", err)
	}
	return
}

func pesan(w http.ResponseWriter, r *http.Request) {
	session := sessions.Start(w, r)

	var t, err = template.ParseFiles("views/pesan.html")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	rows, err := db.Query("select id,nama,email,pesan from `pesan` where aktif like 'T' ")
	if err != nil {
		log.Printf("error query:%+v\n", err)
		return
	}
	defer rows.Close()
	var datas []Pesan
	var datan Pesan
	for rows.Next() {
		datan = Pesan{}
		err = rows.Scan(&datan.ID, &datan.Nama, &datan.Email, &datan.Isipesan)
		if err != nil {
			log.Printf("scan error: %+v\n", err)
			continue
		}
		datas = append(datas, datan)
	}
	response := ResponsePesan{
		Pesana: datas,
		User:   session.GetString("name"),
	}
	err = t.Execute(w, response)
	if err != nil {
		log.Printf("execute template err: %+v\n", err)
	}
	return
}

func about(w http.ResponseWriter, r *http.Request) {
	session := sessions.Start(w, r)
	if len(session.GetString("username")) == 0 {
		//http.Redirect(w, r, "/login", 301)
	}

	var data = map[string]string{
		"User": session.GetString("username"),
	}
	var t, err = template.ParseFiles("views/about.html")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	t.Execute(w, data)
	return

}

func update(w http.ResponseWriter, r *http.Request) {
	judulup := r.FormValue("judulup")
	isiup := r.FormValue("isiup")
	idup := r.FormValue("idup")
	if idup != "" {
		log.Printf("upte cuy:%+v\n", idup)
		stmt, err := db.Prepare("update `article` set judul=? , isi=? where id=?")
		if err == nil {
			_, err := stmt.Exec(judulup, isiup, idup)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			// return 'return dihilangkan supaya tidak keluar
		}
	}
	http.Redirect(w, r, "/articleku", 302)
	return
}
func updatestatrue(w http.ResponseWriter, r *http.Request) {
	//update
	idup := r.FormValue("idupsta")
	if idup != "" {
		log.Printf("upte cuy:%+v\n", idup)
		stmt, err := db.Prepare("update `article` set sta=? where id=?")
		if err == nil {
			_, err := stmt.Exec("T", idup)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			// return 'return dihilangkan supaya tidak keluar
		}
	}
	http.Redirect(w, r, "/articleku", 302)
	return
}
func updateaktif(w http.ResponseWriter, r *http.Request) {
	//update
	idup := r.FormValue("idupsta")
	if idup != "" {
		log.Printf("upte cuy:%+v\n", idup)
		stmt, err := db.Prepare("update `article` set aktif=? where id=?")
		if err == nil {
			_, err := stmt.Exec("F", idup)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			// return 'return dihilangkan supaya tidak keluar
		}
	}
	http.Redirect(w, r, "/articleku", 302)
	return
}
func updatestafalse(w http.ResponseWriter, r *http.Request) {
	//update
	idup := r.FormValue("idupsta")
	if idup != "" {
		log.Printf("upte cuy:%+v\n", idup)
		stmt, err := db.Prepare("update `article` set sta=? where id=?")
		if err == nil {
			_, err := stmt.Exec("F", idup)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			// return 'return dihilangkan supaya tidak keluar
		}
	}
	http.Redirect(w, r, "/articleku", 302)
	return
}
func articleku(w http.ResponseWriter, r *http.Request) {
	session := sessions.Start(w, r)
	if len(session.GetString("username")) == 0 {
		http.Redirect(w, r, "/login", 302)
	}
	var t, err = template.ParseFiles("views/article.html")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	// Insert
	// if r.Method != "POST" {
	// 	http.ServeFile(w, r, "views/article.html")
	// 	return
	// }
	judulna := r.FormValue("judul")
	isina := r.FormValue("isi")

	if judulna != "" {

		stmt, err := db.Prepare("insert into `article` ( judul, isi, sta) values ( ?, ?, ?)")
		if err == nil {
			_, err := stmt.Exec(judulna, isina, "F")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			// return 'return dihilangkan supaya tidak keluar
		}
	}

	// Select

	rows, err := db.Query("select id,judul,isi,sta from `article`  where aktif like 'T' ")
	if err != nil {
		log.Printf("error query:%+v\n", err)
		return
	}
	defer rows.Close()
	var datas []Article
	var datan Article
	for rows.Next() {
		datan = Article{}
		err = rows.Scan(&datan.ID, &datan.Judul, &datan.Isi, &datan.Sta)
		if err != nil {
			log.Printf("scan error: %+v\n", err)
			continue
		}
		datas = append(datas, datan)
	}
	response := ResponseArticle{
		Articles: datas,
		User:     session.GetString("username"),
	}
	err = t.Execute(w, response)
	if err != nil {
		log.Printf("execute template err: %+v\n", err)
	}
	return
}

func contact(w http.ResponseWriter, r *http.Request) {
	turingku := r.FormValue("turing")
	turingini := r.FormValue("turingini")
	nama := r.FormValue("nama")
	email := r.FormValue("email")
	pesan := r.FormValue("pesan")
	var Infona string
	Infona = ""
	if nama != "" {
		if turingku == turingini {
			log.Printf("input : %+v\n", nama)
			stmt, err := db.Prepare("insert into `pesan` ( nama, email, pesan) values ( ?, ?, ?)")
			if err == nil {
				_, err := stmt.Exec(nama, email, pesan)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				Infona = "Pesan Anda berhasil terkirim."
				// return 'return dihilangkan supaya tidak keluar
			}
		} else {
			Infona = "captcha yang anda masukan salah"
		}
	}
	session := sessions.Start(w, r)
	if len(session.GetString("username")) == 0 {
		//http.Redirect(w, r, "/login", 301)
	}
	turing := rand.Intn(999)
	turingna := strconv.Itoa(turing)
	var data = map[string]string{
		"User":   session.GetString("username"),
		"Turing": turingna,
		"Info":   Infona,
	}
	var t, err = template.ParseFiles("views/contact.html")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	t.Execute(w, data)
	return

}

func register(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		session := sessions.Start(w, r)
		var data = map[string]string{
			"Info": session.GetString("eror"),
		}
		var t, err = template.ParseFiles("views/register.html")
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		t.Execute(w, data)
		session.Set("eror", "")
		//http.ServeFile(w, r, "views/register.html")
		return
	}

	username := r.FormValue("email")
	first_name := r.FormValue("first_name")
	last_name := r.FormValue("last_name")
	password := r.FormValue("password")

	users := QueryUser(username)

	if (user{}) == users {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

		if len(hashedPassword) != 0 && checkErr(w, r, err) {
			stmt, err := db.Prepare("INSERT INTO users SET username=?, password=?, first_name=?, last_name=?")
			if err == nil {
				_, err := stmt.Exec(&username, &hashedPassword, &first_name, &last_name)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}
		}
	} else {
		session := sessions.Start(w, r)
		session.Set("eror", "Email sudah dipakai")
		http.Redirect(w, r, "/register", 302)
	}
}

func login(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.ServeFile(w, r, "views/login.html")
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")
	if username != "" && password != "" {
		users := QueryUser(username)

		//deskripsi dan compare password
		var password_tes = bcrypt.CompareHashAndPassword([]byte(users.Password), []byte(password))

		if password_tes == nil {
			//login success
			session := sessions.Start(w, r)
			session.Set("username", users.Username)
			session.Set("name", users.FirstName)
			http.Redirect(w, r, "/", 302)
		} else {
			//login failed
			http.Redirect(w, r, "/login", 302)
		}
	}

}
func logout(w http.ResponseWriter, r *http.Request) {
	session := sessions.Start(w, r)
	session.Clear()
	sessions.Destroy(w, r)
	http.Redirect(w, r, "/", 302)
}
