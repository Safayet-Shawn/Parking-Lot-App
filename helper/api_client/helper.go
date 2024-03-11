package api_client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt"
)

func Respond(w http.ResponseWriter, data map[string]interface{}) {
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
func Message(status bool, message string) map[string]interface{} {
	return map[string]interface{}{"status": status, "message": message}
}

type ApiClient struct {
	req        *http.Request
	queryParam map[string]string
}

func NewApiClient(r *http.Request) *ApiClient {
	return &ApiClient{
		req:        r,
		queryParam: ParseQueryParam(r.URL.RawQuery),
	}
}
func ParseQueryParam(str string) map[string]string {
	ret := make(map[string]string)
	if len(str) < 1 {
		return ret
	}
	query := strings.Split(str, "&")
	if len(query) < 1 {
		return ret
	}
	for _, v := range query {
		tmp := strings.Split(v, "=")
		key, value := tmp[0], tmp[1]
		ret[key] = value
	}
	return ret
}

func (api *ApiClient) JsonBind(i interface{}) error {
	b, err := ioutil.ReadAll(api.req.Body)
	if err != nil {
		log.Printf("failed to read response body, Reason: %v", err)
		return err
	}
	err = json.Unmarshal(b, i)
	if err != nil {
		log.Printf("failed to Unmarshal, Reason: %v", err)
		return err
	}
	return nil
}

func GetToken(r *http.Request) (id, utype int) {
	tokenHeader := r.Header.Get("Authorization") // Grab the token from the header
	if tokenHeader == "" {                       // Token is missing, returns with error code 403 Unauthorized
		return
	}

	splitted := strings.Split(tokenHeader, " ") // The token normally comes in format `Bearer {token-body}`, we check if the retrieved token matched this requirement
	if len(splitted) != 2 {
		return
	}

	tokenPart := splitted[1] // Grab the token part, what we are truly interested in
	s, k := GetTokenClaim(tokenPart)
	return s, k
}
func GetTokenClaim(tokenString string) (userID int, userType int) {
	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		fmt.Println("Error parsing JWT token:", err)
		return
	}

	// Access claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		fmt.Println("Error getting claims from token")
		return
	}

	// Convert claims to JSON format
	_, err = json.Marshal(claims)
	if err != nil {
		fmt.Println("Error marshalling claims to JSON:", err)
		return
	}

	// Access individual claim values
	userID = int(claims["UserId"].(float64))
	userType = int(claims["UserType"].(float64))

	return userID, userType
}
