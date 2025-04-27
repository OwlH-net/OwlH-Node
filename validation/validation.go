package validation

import (
	"errors"

	ndb "github.com/OwlH-net/OwlH-Node/database"
	"github.com/astaxie/beego/logs"
	jwt "github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPasswordHash(password string, hash string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		logs.Error("validation/CheckPasswordHash " + err.Error())
		return false, err
	}
	return true, nil
}

// Generates a jwt token.
func Encode(secret string) (val string, err error) {
	claims := jwt.StandardClaims{
		ExpiresAt: 15000,
		Issuer:    "OwlH",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		logs.Error("Encode error: %s", err)
		return "", err
	}
	return tokenString, err
}

func VerifyToken(token string, user string) (err error) {
	masters, err := ndb.GetMasters()
	if err != nil {
		logs.Error("CheckToken error getting master data: %s", err)
		return err
	}
	for x := range masters {
		tkn, err := Encode(masters[x]["secret"])
		if err != nil {
			logs.Error("Error checking Master token: %s", err)
			return err
		} else {
			if token == tkn {
				return nil
			}
		}
	}
	return errors.New("VerifyToken - Incorrect Token")
}

func VerifyPermissions(user string, object string, permissions []string) (hasPermissions bool, err error) {
	for x := range permissions {
		status, err := UserPermissionsValidation(user, permissions[x])
		if err != nil {
			logs.Error("VerifyPermissions error - requestType error: %s", err)
			return false, err
		}
		if status {
			return true, nil
		}
	}
	return false, nil
}

func GetUserParamValue(uuid, param string) (value string, err error) {
	value, err = ndb.GetUserParamValue(uuid, param)
	if err != nil {
		return "", err
	}
	return value, nil
}

func GetLocalUserPassword(userName string) (encPassword string, err error) {
	uuid, err := ndb.GetLocalUserID(userName)
	if err != nil {
		return "", err
	}
	uPass, err := GetUserParamValue(uuid, "pass")
	if err != nil {
		return "", err
	}
	return uPass, nil
}

func VerifyLocalUser(user map[string]string) bool {
	uPassword, err := GetLocalUserPassword(user["name"])
	if err != nil {
		return false
	}
	valid, err := CheckPasswordHash(user["password"], uPassword)
	if err != nil {
		return false
	}
	return valid
}

func ChangeLocalUserPassword(user map[string]string) (err error) {
	uuid := ""
	uuid, err = ndb.GetLocalUserID(user["name"])
	if err != nil {
		return err
	}
	err = ndb.UpdateUsers(uuid, "pass", user["newpassword"])
	if err != nil {
		return err
	}
	return nil
}
