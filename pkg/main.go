package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Youssef-Shehata/talktuah/internal/auth"
	"github.com/Youssef-Shehata/talktuah/internal/database"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"

	"encoding/json"
	"time"
)

type User struct {
	Id        int64 `json:"id"`
	Username  string    `json:"username"`
	Token     string    `json:"token"`
	CreatedAt time.Time `json:"created_at"`
}

type Message struct {
	Id       int64 `json:"id"`
	ChatId   int64 `json:"chat_id"`
	SenderId int64 `json:"sender_id"`
	SentAt   time.Time `json:"sent_at"`
	Content  string    `json:"content"`
}

type Chat struct {
	Id           int64 `json:"id"`
	CreationDate time.Time `json:"creation_date"`
}

type ChatMembers struct {
	ChatId   int64 `json:"chat_id"`
	UserId   int64 `json:"user_id"`
	JoinDate time.Time `json:"join_date"`
}

type apiConfig struct {
	platform string
	rdb      *redis.Client
	query    *database.Queries
	ctx      context.Context
	secret   string
	logger   *logger
}

func Init() (*apiConfig, *mux.Router) {

	godotenv.Load()
	// Initialize the logger
	logFile := "server.log"
	logger, err := NewLogger(logFile)
	if err != nil {
		log.Printf("Could not initialize logger: %v\n", err)
	}

	ctx := context.Background()
	router := mux.NewRouter()
	db, err := sql.Open("sqlite3", "../sql/schema/chat.db")
	if err != nil {
		log.Fatal("ERROR: connecting to db")
	}

	query := database.New(db)
	rdb := newRedisClient()
	cfg := &apiConfig{ctx: ctx, platform: os.Getenv("PLATFORM"), secret: os.Getenv("SECRET"), query: query, rdb: rdb, logger: logger}
	return cfg, router
}

func (cfg *apiConfig) signUp(w http.ResponseWriter, r *http.Request) {

	cfg.logger.Log(INFO, fmt.Errorf("SIGNING YOU UPPP "))
	type params struct {
		Password         string `json:"password"`
		Username         string `json:"username"`
		ExpiresInSeconds int    `json:"expires_in"`
	}
	var p params
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		cfg.logger.Log(ERROR, errors.Wrap(err, "Decoding Request Json"))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if p.Username == "" {
		cfg.logger.Log(ERROR, fmt.Errorf("Empty Username"))
		http.Error(w, "Empty Username", http.StatusBadRequest)
		return

	}

	hashedPass, err := auth.HashPassword(p.Password)
	if err != nil {
		cfg.logger.Log(ERROR, errors.Wrap(err, "Hashing Password"))
		http.Error(w, "Hashing Password", http.StatusInternalServerError)
		return
	}

	// HOW TO KNOW USERNAME IS TAKEN WITH THIS VAGUE ASS ERROR SHIT
	user, err := cfg.query.CreateUser(cfg.ctx, database.CreateUserParams{Username: p.Username, Password: hashedPass})
	if err != nil {
		//for timing attacks
		time.Sleep(time.Second)
		cfg.logger.Log(ERROR, errors.Wrap(err, "Creating User"))
		http.Error(w, "Creating User", http.StatusInternalServerError)
		return
	}

	token, err := auth.MakeJWT(user.ID, cfg.secret, p.ExpiresInSeconds)
	if err != nil {
		cfg.logger.Log(ERROR, errors.Wrap(err, "Creating User Token"))
		http.Error(w, "Creating User Token", http.StatusInternalServerError)
		return
	}
	respondWithJSON(w, http.StatusOK, User{
		Id:        user.ID,
		CreatedAt: user.CreatedAt,
		Token:     token,
		Username:  user.Username,
	})
}
func (cfg *apiConfig) logIn(w http.ResponseWriter, r *http.Request) {

	type params struct {
		Username         string `json:"username"`
		Password         string `json:"password"`
		ExpiresInSeconds int    `json:"expires_in_seconds"`
	}

	var p params
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		cfg.logger.Log(ERROR, errors.Wrap(err, "Decoding Request Json"))
		http.Error(w, "Decoding Request Json", http.StatusBadRequest)
		return
	}
	if p.Username == "" {
		cfg.logger.Log(ERROR, fmt.Errorf("Empty Username"))
		http.Error(w, "Empty Username", http.StatusBadRequest)
	}
	user, err := cfg.query.GetUserByUsername(cfg.ctx, p.Username)
	if err != nil {
		cfg.logger.Log(ERROR, errors.Wrap(err, "Wrong Username or Password"))
		http.Error(w, "Wrong Username or Password", http.StatusUnauthorized)
		return
	}

	if error := auth.CheckHashedPassword(p.Password, user.Password); error != nil {
		cfg.logger.Log(ERROR, errors.Wrap(error, "Wrong Username or Password"))
		http.Error(w, "Wrong Username or Password", http.StatusUnauthorized)
		return
	}

	token, err := auth.MakeJWT(user.ID, cfg.secret, p.ExpiresInSeconds)
	if err != nil {
		cfg.logger.Log(ERROR, errors.Wrap(err, "Creating User Token"))
		http.Error(w, "Creating User Token", http.StatusInternalServerError)
		return

	}
	respondWithJSON(w, http.StatusOK, User{
		Id:        user.ID,
		CreatedAt: user.CreatedAt,
		Token:     token,
		Username:  user.Username,
	})
}

func (cfg *apiConfig) newChat(w http.ResponseWriter, r *http.Request) {
}

func (cfg *apiConfig) sendMessage(w http.ResponseWriter, r *http.Request) {
	//vars := mux.Vars(r)
	//chat_id := vars["id"]
}

func (cfg *apiConfig) chatMessages(w http.ResponseWriter, r *http.Request) {
	//vars := mux.Vars(r)
	//chat_id := vars["id"]
}
func main() {
	godotenv.Load()
	cfg, router := Init()
	defer cfg.logger.Close()

	// all of these end points are authed and they get the user id from the auth middleware
	//create a new chat room , AUTHENTICATED users only , return chat room id
	router.HandleFunc("/api/v1/chat", cfg.newChat).Methods("POST")

	// send messages to chat room, Authed, in chat room members
	router.HandleFunc("/api/v1/chat/{id}", cfg.sendMessage).Methods("POST")

	// get messages of the chat
	router.HandleFunc("/api/v1/chat/{id}", cfg.chatMessages).Methods("GET")

	router.HandleFunc("/api/v1/signup", cfg.signUp).Methods("POST")
	router.HandleFunc("/api/v1/login", cfg.logIn).Methods("POST")

	server := http.Server{
        Handler: router,
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
		Addr:         "localhost:8080",
    }

	err := server.ListenAndServe()
	if err != nil {
		cfg.logger.Log(ERROR, errors.Wrap(err, "Strating Server"))
		return
	}

}

func respondWithJSON(w http.ResponseWriter, status int, payload interface{}) {
	res, err := json.Marshal(payload)
	if err != nil {
		log.Printf("  ERROR : couldn't parse json : %v\n", err)
		http.Error(w, "Failed To Marshal Json", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(res)
}
