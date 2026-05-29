package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/FunnySneko/ReadPill/server/internal/aggregation"
	"github.com/FunnySneko/ReadPill/server/internal/api"
	"github.com/FunnySneko/ReadPill/server/internal/db"
	"github.com/FunnySneko/ReadPill/server/internal/user_actions"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

func tryErrorOut(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			next.ServeHTTP(w, r)
			return
		}

		switch r.URL.Path {
		case "/login", "/signup":
			next.ServeHTTP(w, r)
			return
		}

		cookie, err := r.Cookie("token")
		if err != nil {
			api.ErrorOut(w, err, "unauthorized", http.StatusUnauthorized)
			return
		}

		claims := jwt.MapClaims{}

		token, err := jwt.ParseWithClaims(
			cookie.Value,
			claims,
			func(t *jwt.Token) (any, error) {
				_, ok := t.Method.(*jwt.SigningMethodHMAC)
				if !ok {
					return nil, fmt.Errorf(
						"unexpected signing method: %v",
						t.Header["alg"],
					)
				}

				return []byte(os.Getenv("JWTSECRET")), nil
			},
		)

		if err != nil || token == nil || !token.Valid {
			api.ErrorOut(w, err, "invalid token", http.StatusUnauthorized)
			return
		}

		userID, ok := claims["user_id"]
		if !ok {
			api.ErrorOut(w, nil, "missing user id", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "userID", userID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func main() {
	err := godotenv.Load()
	tryErrorOut(err)

	db := db.Db{}
	err = db.ConnectToDatabase()
	tryErrorOut(err)
	defer db.CloseConnection()

	agg, err := aggregation.NewAggregator()
	tryErrorOut(err)

	serverHandler := api.NewServerHandler(
		&db,
		agg,
		user_actions.NewActionsHandler(&db),
	)

	mux := http.NewServeMux()

	api.HandlerFromMux(serverHandler, mux)

	mux.Handle(
		"/images/",
		http.StripPrefix("/images/", http.FileServer(http.Dir("./uploads"))),
	)

	protected := AuthMiddleware(mux)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"*"},
		AllowCredentials: true,
		MaxAge:           300,
	})

	handler := c.Handler(protected)

	fmt.Println("Server starting on :8080")
	log.Fatalln(http.ListenAndServe(":8080", handler))
}
