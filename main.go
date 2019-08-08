package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"math/rand" // "os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/kataras/go-sessions"
	"golang.org/x/crypto/bcrypt"
)

var db *sql.DB
var err error

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

type ResponseArticle struct {
	Articles []Article
	User     string
}

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

func routes() {
	http.HandleFunc("/", home)
	http.HandleFunc("/register", register)
	http.HandleFunc("/login", login)
	http.HandleFunc("/logout", logout)
	http.HandleFunc("/about", about)
	http.HandleFunc("/article", article)
	http.HandleFunc("/contact", contact)
	http.HandleFunc("/update", update)
	http.HandleFunc("/updatestatrue", updatestatrue)
	http.HandleFunc("/updatestafalse", updatestafalse)
	http.HandleFunc("/updateaktif", updateaktif)
}

func main() {
	connect_db()
	routes()

	defer db.Close()

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
	if len(session.GetString("username")) == 0 {
		//http.Redirect(w, r, "/login", 301)
	}

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
		User:     session.GetString("username"),
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
		"username": session.GetString("username"),
		"message":  "Welcome to the Go !",
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
	//update
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
	http.Redirect(w, r, "/article", 302)
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
	http.Redirect(w, r, "/article", 302)
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
	http.Redirect(w, r, "/article", 302)
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
	http.Redirect(w, r, "/article", 302)
	return
}
func article(w http.ResponseWriter, r *http.Request) {
	session := sessions.Start(w, r)
	if len(session.GetString("username")) == 0 {
		//http.Redirect(w, r, "/login", 301)
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
	//datas := Queryarticle()
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
	log.Printf("siapkan : %+v\n", nama)
	if turingku == turingini && nama != "" {
		log.Printf("input : %+v\n", nama)
		stmt, err := db.Prepare("insert into `pesan` ( nama, email, pesan) values ( ?, ?, ?)")
		if err == nil {
			_, err := stmt.Exec(nama, email, pesan)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			// return 'return dihilangkan supaya tidak keluar
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
		http.ServeFile(w, r, "views/register.html")
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
		http.Redirect(w, r, "/register", 302)
	}
}

func login(w http.ResponseWriter, r *http.Request) {
	session := sessions.Start(w, r)
	if len(session.GetString("username")) != 0 && checkErr(w, r, err) {
		http.Redirect(w, r, "/", 302)
	}
	if r.Method != "POST" {
		http.ServeFile(w, r, "views/login.html")
		return
	}
	username := r.FormValue("username")
	password := r.FormValue("password")

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
func logout(w http.ResponseWriter, r *http.Request) {
	session := sessions.Start(w, r)
	session.Clear()
	sessions.Destroy(w, r)
	http.Redirect(w, r, "/", 302)
}
