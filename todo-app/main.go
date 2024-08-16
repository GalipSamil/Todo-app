package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

// Todo yapısı bir görevi tanımlar
type Todo struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

var db *sql.DB

// Veritabanını başlat ve tabloları oluştur
func initDB() {
	var err error
	db, err = sql.Open("sqlite3", "./todos.db")
	if err != nil {
		log.Fatal(err)
	}

	statement, err := db.Prepare(`CREATE TABLE IF NOT EXISTS todos (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		completed BOOLEAN NOT NULL CHECK (completed IN (0, 1))
	)`)
	if err != nil {
		log.Fatal(err)
	}
	statement.Exec()
}

// Görev ekleme fonksiyonu (Veritabanına kaydetme)
func addTodoToDB(title string) error {
	statement, err := db.Prepare("INSERT INTO todos (title, completed) VALUES (?, ?)")
	if err != nil {
		return err
	}
	_, err = statement.Exec(title, false)
	return err
}

// Görevleri veritabanından listeleme fonksiyonu
func listTodosFromDB() ([]Todo, error) {
	rows, err := db.Query("SELECT id, title, completed FROM todos")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var todos []Todo
	for rows.Next() {
		var todo Todo
		err := rows.Scan(&todo.ID, &todo.Title, &todo.Completed)
		if err != nil {
			return nil, err
		}
		todos = append(todos, todo)
	}
	return todos, nil
}

// Görevi tamamlanmış olarak işaretleme fonksiyonu (Veritabanında güncelleme)
func completeTodoInDB(id int) error {
	statement, err := db.Prepare("UPDATE todos SET completed = ? WHERE id = ?")
	if err != nil {
		return err
	}
	_, err = statement.Exec(true, id)
	return err
}

// Görev silme fonksiyonu (Veritabanından silme)
func deleteTodoFromDB(id int) error {
	statement, err := db.Prepare("DELETE FROM todos WHERE id = ?")
	if err != nil {
		return err
	}
	_, err = statement.Exec(id)
	return err
}

// Görev ekleme endpointi
func addTodoHandler(w http.ResponseWriter, r *http.Request) {
	var todo Todo
	if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err := addTodoToDB(todo.Title)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

// Görevleri listeleme endpointi
func listTodosHandler(w http.ResponseWriter, _ *http.Request) {
	todos, err := listTodosFromDB()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todos)
}

// Görevi tamamlanmış olarak işaretleme endpointi
func completedTodoHandler(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(strings.TrimSuffix(strings.TrimPrefix(r.URL.Path, "/todos/"), "/complete"))
	err := completeTodoInDB(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// Görev silme endpointi
func deleteTodoHandler(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/todos/"))
	err := deleteTodoFromDB(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// Ana fonksiyon - HTTP sunucusunu başlatır
func main() {
	initDB()
	defer db.Close()

	http.HandleFunc("/todos", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			addTodoHandler(w, r)
		} else if r.Method == http.MethodGet {
			listTodosHandler(w, r)
		}
	})

	http.HandleFunc("/todos/", func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/complete") {
			completedTodoHandler(w, r)
		} else {
			deleteTodoHandler(w, r)
		}
	})

	// Statik dosyalara erişim için (HTML, CSS, JS)
	http.Handle("/", http.FileServer(http.Dir("./static")))

	// Sunucuyu başlat
	fmt.Println("Server is running at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
