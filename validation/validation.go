package validation

import (
	"github.com/astaxie/beego/logs"
	jwt "github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"errors"
    "owlhnode/database"
)

func HashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    logs.Info("NEW HASH PASSWD--> "+string(bytes))
    return string(bytes), err
}

func CheckPasswordHash(password string, hash string) (bool, error) {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))	
	if err != nil {logs.Error("validation/CheckPasswordHash "+err.Error()); return false, err}
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
	if err != nil {logs.Error("Encode error: %s", err); return "", err}
	return tokenString, err
}

func VerifyToken(token string, userUuid string)(err error){
	masters,err := ndb.GetMasters(); if err != nil {logs.Error("CheckToken error getting master data: %s", err); return err}
	for x := range masters{
		tkn, err := Encode(masters[x]["secret"])
		if err != nil {
			logs.Error("Error checking Master token: %s", err); return err
		}else{
			if token == tkn {
			// 	if userUuid == "none"{
					return nil
			// 	}else{
			// 		return errors.New("This user has not enough privileges level")
			// 	}
			// }else{
			// 	return errors.New("The token retrieved is false")
			}
		}		
	}
	return errors.New("VerifyToken - Incorrect Token")
}

func VerifyPermissions(uuidUser string, object string, permissions []string)(hasPermissions bool, err error){
	logs.Info("User UUID --> "+uuidUser)
	logs.Info("User Permissions --> %v", permissions)
	for x := range permissions{
		logs.Info("Permission requested --> "+permissions[x])
		status,err := UserPermissionsValidation(uuidUser, permissions[x]); if err != nil {logs.Error("VerifyPermissions error - requestType error: %s",err); return false,err}
		if status{
			return true,nil
		}		
	}
	return false,nil
}