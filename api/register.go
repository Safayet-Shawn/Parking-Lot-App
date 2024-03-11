package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"

	"github.com/Safayet-Shawn/Parking-Lot/db"
	"github.com/Safayet-Shawn/Parking-Lot/helper/api_client"

	"github.com/nyaruka/phonenumbers"
	"golang.org/x/crypto/bcrypt"
)

func RegistrationUser(w http.ResponseWriter, r *http.Request) {
	// driverID := r.Context().Value("token").(*middlewares.Token).UserId
	// tk := r.Context().Value("token")
	// tok := tk.(*db.Token)
	// log.Printf("User id:%v", tok.UserId)
	// i := tok.UserId
	id, _ := api_client.GetToken(r)
	fmt.Println("=============userid============", id)
	user := &User{}
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		// applog.Errorf(r.Context(), "failed to read response body, Reason: %v", err)
		msg := api_client.Message(false, err.Error())
		api_client.Respond(w, msg)
		return
	}

	err = json.Unmarshal(b, user)
	if err != nil {
		// applog.Errorf(r.Context(), "Failed to unmarshal user[type] info ,Reason:%v", err.Error())
		msg := api_client.Message(false, err.Error())
		api_client.Respond(w, msg)
		return
	}

	num, _ := phonenumbers.Parse(user.Phone, "BANGLADESH")
	regionNumber := phonenumbers.GetRegionCodeForNumber(num)
	countryCode := phonenumbers.GetCountryCodeForRegion(regionNumber)

	if countryCode == 880 {
		// format it using national format
		user.Phone = phonenumbers.Format(num, phonenumbers.NATIONAL)
		user.Phone = strings.Replace(user.Phone, "-", "", -1)
	}

	if user.UserType == 1 || user.UserType == 2 {
		ret := validateUser(*user)
		if ret != "" {
			w.WriteHeader(http.StatusBadRequest)
			msg := api_client.Message(false, err.Error())
			api_client.Respond(w, msg)
			return
		}
		_, err = saveUserToDB(*user)
		if err != nil {
			// applog.Errorf(r.Context(), "failed to save user into database, Reason: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			msg := api_client.Message(false, err.Error())
			api_client.Respond(w, msg)
			return
		}
		resp := api_client.Message(true, "successfully created User")
		api_client.Respond(w, resp)
		return
	} else {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(" user cannot be implemented except user type 1 & 2"))
		return
	}
}

func validateEmail(email string) bool {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`

	re := regexp.MustCompile(pattern)

	// Check if the email matches the pattern
	return re.MatchString(email)
}

func validatePhone(phone string) bool {
	Re := regexp.MustCompile(`^(?:(?:\(?(?:00|\+)([1-4]\d\d|[1-9]\d?)\)?)?[\-\.\ \\\/]?)?((?:\(?\d{1,}\)?[\-\.\ \\\/]?){0,})(?:[\-\.\ \\\/]?(?:#|ext\.?|extension|x)[\-\.\ \\\/]?(\d+))?$`)
	return Re.MatchString(phone) && len(phone) != 12
}

func validateUser(user User) string {
	if user.UserType == 2 {
		if user.Email != "" && !validateEmail(user.Email) {
			return fmt.Sprintf("email empty or invalid")
		}
	} else {
		if user.Address == "" {
			return fmt.Sprintf("need to provide address")
		}
		if user.Name == "" {
			return fmt.Sprintf("need to provide name")
		}
		if user.Email != "" && !validateEmail(user.Email) {
			return fmt.Sprintf("email empty or invalid")
		}

	}

	if user.Phone != "" && !validatePhone(user.Phone) {
		return fmt.Sprintf("phone number empty or invalid")
	}

	if user.Password == "" {
		return fmt.Sprintf("need to provide passoword")
	}
	if user.RepeatPassword == "" {
		return fmt.Sprintf("need to provide repeat_password")
	}

	if user.Password != user.RepeatPassword {
		return fmt.Sprintf("password and repeat password need to match")
	}
	return ""
}

func saveUserToDB(user User) (*db.User, error) {
	d := db.GetDB()
	dUser := &db.User{}

	dUser.Name = user.Name
	dUser.Phone = user.Phone
	dUser.Email = user.Email
	dUser.Address = user.Address
	dUser.UserType = user.UserType
	gPass := hashAndSalt([]byte(user.Password), context.Background())
	dUser.Password = gPass

	res := d.Debug().Create(&dUser)
	if res.Error != nil {
		return nil, res.Error
	}

	err := d.Table("users").Where("user_type= ?", 2).Last(&dUser).Error
	if err != nil {
		return nil, nil
	}
	return dUser, nil

}

func hashAndSalt(pwd []byte, ctx context.Context) string {
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		// applog.Errorf(ctx, "Failed to generate hash password from ,given password ,reason:%v", err)
	}
	// GenerateFromPassword returns a byte slice so we need to
	// convert the bytes to a string and return it
	return string(hash)
}

func comparePasswords(hashedPwd string, plainPwd []byte, r *http.Request) bool {
	// Since we'll be getting the hashed password from the DB it
	// will be a string so we'll need to convert it to a byte slice
	byteHash := []byte(hashedPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, plainPwd)
	if err != nil {
		// applog.Errorf(r.Context(), "hashed and plained password not mached/failed to compare, Reason:%v", err)
		return false
	}

	return true
}
