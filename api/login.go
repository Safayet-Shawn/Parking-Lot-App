package api

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/Safayet-Shawn/Parking-Lot/config"
	"github.com/Safayet-Shawn/Parking-Lot/db"
	"github.com/Safayet-Shawn/Parking-Lot/helper/api_client"
	"github.com/Safayet-Shawn/Parking-Lot/helper/redis"

	"github.com/golang-jwt/jwt"
	"github.com/nyaruka/phonenumbers"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var TokenPassword = "987654"

func UserLogin(w http.ResponseWriter, r *http.Request) {
	login := Login{}
	account := &db.User{}
	err := api_client.NewApiClient(r).JsonBind(&login)
	if err != nil {
		log.Printf("failed to bind json, Reason: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		api_client.Respond(w, api_client.Message(false, err.Error()))
		return
	}

	if len(login.Phone) > 0 {
		login.Identity = login.Phone
	}
	phone, email, err := checkEmailOrPhone(login.Identity)
	if err != nil {
		log.Printf("failed to set email or phone number, reason: %s", err)
		api_client.Respond(w, api_client.Message(false, err.Error()))
		return
	}

	query := db.GetDB().Debug()
	if len(phone) > 0 {
		query = query.Where("phone = ?", phone)
	}
	if len(email) > 0 {
		query = query.Where("email = ?", email)
	}
	err = query.Order("id DESC").First(account).Error
	if err != nil && err == gorm.ErrRecordNotFound {
		log.Printf("failed to find user %s in db, reason: %s", login.Phone, err)
		msg := api_client.Message(false, "phone or email not found")
		api_client.Respond(w, msg)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(login.Password))
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword { // Password does not match!
		log.Printf("password doesn't match for user %s, reason: %s", login.Phone, err)
		msg := api_client.Message(false, "Invalid login credentials. Please try again")
		api_client.Respond(w, msg)
		return
	}

	// Password matched, Need to check if the user type is expected
	v := validateLoginData(login, account)
	if v != nil {
		log.Printf("user type not match for user %s , reason: %s", login.Phone, v)
		writeBodyBadRequest(w, v.Error())
		return
	}

	resp := LoginProcess(account)
	api_client.Respond(w, resp)
}

func validateLoginData(user Login, dbUser *db.User) error {
	if user.UserType != dbUser.UserType {
		return errors.New("UserType didn't matched")
	}
	return nil
}

func writeBodyBadRequest(w http.ResponseWriter, msg string) {
	w.WriteHeader(http.StatusBadRequest)
	m := api_client.Message(false, msg)
	api_client.Respond(w, m)
}

func checkEmailOrPhone(identity string) (string, string, error) {
	if validatePhone(identity) {
		num, _ := phonenumbers.Parse(identity, "BANGLADESH")
		regionNumber := phonenumbers.GetRegionCodeForNumber(num)
		countryCode := phonenumbers.GetCountryCodeForRegion(regionNumber)

		if countryCode == 880 {
			// format it using national format
			identity = phonenumbers.Format(num, phonenumbers.NATIONAL)
			identity = strings.Replace(identity, "-", "", -1)
		}
		return identity, "", nil
	}

	if validateEmail(identity) {
		return "", identity, nil
	}

	return "", "", errors.New("invalid email or phone number provided")
}
func LoginProcess(account *db.User) map[string]interface{} {
	// We are already logged in
	account.Password = ""
	// Create JWT token
	tk := &db.Token{
		UserId:   account.Id,
		UserType: account.UserType,
		Name:     account.Name,
		Address:  account.Address,
	}

	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(config.JwtSecret()))
	ctx := context.Background()
	err := redis.GetRedis().Set(ctx, "jwt", tokenString, time.Hour*24)
	if err != nil {
		log.Printf(" set  JWT token while login into redis where  err:%v", err)
	}
	resp := api_client.Message(true, "Logged In")
	resp["account"] = *account
	resp["token"] = tokenString
	return resp
}
