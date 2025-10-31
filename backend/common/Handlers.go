package handlers

import (
	"database/sql"
	"encoding/json"
	"eventplanner/backend/data"
	"net/http"

	"github.com/gorilla/securecookie"
	_ "github.com/lib/pq"
)

var cookieHandler = securecookie.New(
	securecookie.GenerateRandomKey(64),
	securecookie.GenerateRandomKey(32))

var db *sql.DB

func InitDB(dbConn *sql.DB) {
	db = dbConn
}

type User struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

// POST /register
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var u User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if u.UserName == "" || u.Password == "" {
		http.Error(w, "Username and password required", http.StatusBadRequest)
		return
	}

	exists, err := data.UserExists(db, u.UserName)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	if exists {
		http.Error(w, "Username already exists", http.StatusConflict)
		return
	}

	if err := data.CreateUser(db, u.UserName, u.Password); err != nil {
		http.Error(w, "Could not register user", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "User registered successfully"})
}

// POST /login
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var u User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	valid, err := data.ValidateUser(db, u.UserName, u.Password)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	if !valid {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	SetCookie(u.UserName, w)
	json.NewEncoder(w).Encode(map[string]string{"message": "Login successful"})
}

// POST /logout
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	ClearCookie(w)
	json.NewEncoder(w).Encode(map[string]string{"message": "Logged out successfully"})
}

// ---------------- Cookies -----------------

func SetCookie(userName string, w http.ResponseWriter) {
	value := map[string]string{"name": userName}
	if encoded, err := cookieHandler.Encode("session", value); err == nil {
		http.SetCookie(w, &http.Cookie{
			Name:  "session",
			Value: encoded,
			Path:  "/",
		})
	}
}

func ClearCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:   "session",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})
}
