package main
import (

	"context"
	"log"
	"net/http"

	"github.com/Youssef-Shehata/talktuah/internal/auth"
)

func (cfg *apiConfig) authMiddleware(next http.HandlerFunc) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		token := auth.GetBearerToken(r.Header)

		userId, err := auth.ValidateJWT(token, cfg.secret)
		if err != nil {
			log.Printf("  ERROR: auth token : %v\n", err.Error())
			http.Error(w, "", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "userid", userId)

		next(w, r.WithContext(ctx))
	})

}

