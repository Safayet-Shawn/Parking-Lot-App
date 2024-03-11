package route

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/Safayet-Shawn/Parking-Lot/db"
	"github.com/Safayet-Shawn/Parking-Lot/helper/api_client"
	"github.com/Safayet-Shawn/Parking-Lot/helper/redis"

	"github.com/golang-jwt/jwt"
	// "github.com/adeffi/monolith/admin-dashboard/pkg/db"
	// "github.com/adeffi/monolith/base/go/applog"
	// "github.com/adeffi/monolith/base/go/xcontext"
	// "github.com/adeffi/monolith/helpers/api_client"
	// "github.com/golang-jwt/jwt"
	// "github.com/adeffi/monolith/admin-dashboard/pkg/db"
	// "github.com/adeffi/monolith/base/go/applog"
	// "github.com/adeffi/monolith/base/go/xcontext"
	// "github.com/adeffi/monolith/helpers/api_client"
	// "github.com/golang-jwt/jwt"
)

var (
	hash_token string
)

// driverID := r.Context().Value("token").(*middlewares.Token).UserId

var mp = map[string]string{
	// "/v1/register": "1",
	"/v1/login": "1",
}

func SetHashToken(tk string) {
	hash_token = tk
}
func checkNotauth(r *http.Request) bool {
	if _, ok := mp[r.URL.String()]; ok {
		return true
	}
	if user, pass, ok := r.BasicAuth(); ok && user == "hello" && pass == "hello" {
		return true
	}
	return false
}
func JwtAuthentication(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if checkNotauth(r) {
			next.ServeHTTP(w, r)
			return
		}

		response := make(map[string]interface{})
		tokenHeader := r.Header.Get("Authorization") // Grab the token from the header
		if tokenHeader == "" {                       // Token is missing, returns with error code 403 Unauthorized
			response = api_client.Message(false, "Missing auth token")
			w.WriteHeader(http.StatusUnauthorized)
			w.Header().Add("Content-Type", "application/json")
			api_client.Respond(w, response)
			return
		}

		splitted := strings.Split(tokenHeader, " ") // The token normally comes in format `Bearer {token-body}`, we check if the retrieved token matched this requirement
		if len(splitted) != 2 {
			response = api_client.Message(false, "Invalid/Malformed auth token")
			w.WriteHeader(http.StatusUnauthorized)
			w.Header().Add("Content-Type", "application/json")
			api_client.Respond(w, response)
			return
		}

		tokenPart := splitted[1] // Grab the token part, what we are truly interested in
		log.Printf("jwt token string:%v", tokenPart)

		tk := &db.Token{}
		if hash_token == "" {
			response = api_client.Message(false, "token password is empty")
			w.WriteHeader(http.StatusUnauthorized)
			w.Header().Add("Content-Type", "application/json")
			api_client.Respond(w, response)
			return
		}
		token, err := jwt.ParseWithClaims(tokenPart, tk, func(token *jwt.Token) (interface{}, error) {
			return []byte(hash_token), nil
		})
		if err != nil { // Malformed token, returns with http code 403 as usual
			log.Printf("failed to parse with claims[malformed token] ,reason:%v", err)
			response = api_client.Message(false, "Malformed authentication token")
			w.WriteHeader(http.StatusUnauthorized)
			w.Header().Add("Content-Type", "application/json")
			api_client.Respond(w, response)
			return
		}

		if !token.Valid { // Token is invalid, maybe not signed on this server
			response = api_client.Message(false, "Token is not valid.")
			w.WriteHeader(http.StatusUnauthorized)
			w.Header().Add("Content-Type", "application/json")
			api_client.Respond(w, response)
			return
		}
		ctx := context.Background()

		redToken, err := redis.GetRedis().Get(ctx, "jwt").Result()
		if err != nil {
			log.Printf("failed to value from redis ", err)
		}
		if redToken == "" {
			response = api_client.Message(false, "token is expired")
			w.WriteHeader(http.StatusUnauthorized)
			w.Header().Add("Content-Type", "application/json")
			api_client.Respond(w, response)
			return
		}
		if err != nil {
			response = api_client.Message(false, "failed to get token from redis")
			w.WriteHeader(http.StatusUnauthorized)
			w.Header().Add("Content-Type", "application/json")
			api_client.Respond(w, response)
			return
		}
		next.ServeHTTP(w, r) // proceed in the middleware chain!
	})
}
