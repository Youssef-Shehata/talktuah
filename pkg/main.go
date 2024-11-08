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
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"

	"encoding/json"
	"time"
)

type User struct {
    Id uuid.UUID `json:"id"`;
    Username string `json:"username"`;
    Token string `json:"token"`;
    CreatedAt time.Time`json:"created_at"`;
}

type Message struct {
    Id uuid.UUID `json:"id"`;
    ChatId uuid.UUID `json:"chat_id"`;
    SenderId uuid.UUID `json:"sender_id"`;
    SentAt time.Time`json:"sent_at"`;
    Content string `json:"content"`;
}

type Chat struct {
    Id uuid.UUID `json:"id"`;
    CreationDate time.Time`json:"creation_date"`;
}

type ChatMembers struct {
    ChatId uuid.UUID `json:"chat_id"`;
    UserId uuid.UUID `json:"user_id"`;
    JoinDate time.Time`json:"join_date"`;
}

type apiConfig struct {
	platform string
    rdb      *redis.Client
	query    *database.Queries
	ctx      context.Context
	secret   string
}

func Init() (*apiConfig, *http.ServeMux) {

	godotenv.Load()
	ctx := context.Background()
	mux := http.NewServeMux()
	db, err := sql.Open("postgres", os.Getenv("DB_URL"))
    log.Printf("url : %v" , os.Getenv("DB_URL"))
	query := database.New(db)
	if err != nil {
		log.Fatal("ERROR: connecting to db")
	}

    rdb := newRedisClient();
	cfg := &apiConfig{ctx: ctx, platform: os.Getenv("PLATFORM"),   secret: os.Getenv("SECRET"), query: query,  rdb: rdb}
	return cfg, mux
}




func (cfg *apiConfig) signUp(w http.ResponseWriter, r *http.Request) {

	type params struct {
		Email            string `json:"email"`
		Password         string `json:"password"`
		Username         string `json:"username"`
		ExpiresInSeconds int    `json:"expires_in"`
	}
	var p params
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		log.Printf("  ERROR bad request api/createUser: %v \n", err)
		http.Error(w, fmt.Sprint("bad request : ", err.Error()), http.StatusBadRequest)
		return
	}
	if p.Email == "" {
		log.Printf("  ERROR bad request to api/createUser: empty email field\n")
		http.Error(w, "email can't be empty", http.StatusBadRequest)
		return

	}
	if p.Username== "" {
		log.Printf("  ERROR bad request to api/createUser: empty username field\n")
		http.Error(w, "username can't be empty", http.StatusBadRequest)
		return

	}


	hashedPass, err := auth.HashPassword(p.Password)
	if err != nil {
		log.Printf("  ERROR failed to create new user : %v \n", err.Error())
		http.Error(w, "failure hashing password",http.StatusInternalServerError)
		return
	}

	user, err := cfg.query.CreateUser(cfg.ctx, database.CreateUserParams{Email: p.Email, Password: hashedPass})


    // HOW TO KNOW USERNAME IS TAKEN WITH THIS VAGE ERROR SHIT 
	if err != nil {
		time.Sleep(time.Second)
		log.Printf("  ERROR failed to create new user : %v \n", err)
		http.Error(w, "couldnt create new user", http.StatusInternalServerError)
		return
	}

	token, err := auth.MakeJWT(user.ID, cfg.secret, p.ExpiresInSeconds)
	if err != nil {
		log.Printf("  ERROR making token : %v", err.Error())
		http.Error(w, "", http.StatusInternalServerError)
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
		Email            string `json:"email"`
		Password         string `json:"password"`
		ExpiresInSeconds int    `json:"expires_in_seconds"`
	}

	var p params
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		log.Printf("  ERROR bad request api/createUser: %v \n", err)
		http.Error(w, fmt.Sprint("bad request : ", err.Error()), http.StatusBadRequest)
		return
	}

	if p.Email == "" {
		log.Printf("  ERROR bad request to api/createUser: empty email field\n")
		http.Error(w, "email can't be empty", http.StatusBadRequest)
		return
	}
	user, err := cfg.query.GetUserByUsername(cfg.ctx, p.Email)
	if err != nil {
		log.Printf("  ERROR Wrong email or password \n %v", err)
		http.Error(w, "Wrong Email or Password", http.StatusUnauthorized)
		return
	}

	if error := auth.CheckHashedPassword(p.Password, user.Password); err != nil {
		log.Printf("  ERROR Wrong email or password \n %v", error)
		http.Error(w, "Wrong Email or Password", http.StatusUnauthorized)
		return
	}

	token, err := auth.MakeJWT(user.ID, cfg.secret, p.ExpiresInSeconds)
	if err != nil {
		log.Printf("  ERROR making token : %v", err.Error())
		http.Error(w, "", http.StatusInternalServerError)
		return

	}
	respondWithJSON(w, http.StatusOK, User{
		Id:        user.ID,
		CreatedAt: user.CreatedAt,
		Token:     token,
		Username:  user.Username,
	})
}





func main() {
	cfg, mux := Init()


	mux.HandleFunc("POST /api/v1/signup", cfg.signUp)
	mux.HandleFunc("POST /api/v1/login", cfg.logIn)


	server := http.Server{Handler: mux, WriteTimeout: 10 * time.Second, ReadTimeout: 10 * time.Second, Addr: "localhost:"+ os.Getenv("PORT")}
	err := server.ListenAndServe()

	if err != nil {
		log.Print("  ERROR: starting server:", err)
	}

}

func respondWithJSON(w http.ResponseWriter, status int, payload interface{}) {
	res, err := json.Marshal(payload)
	if err != nil {
		log.Printf("  ERROR: couldn't parse json : %v\n", err)
		http.Error(w, "failed to marshal json", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(res)
}
