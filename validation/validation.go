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
	if err != nil {logs.Error(err); return false, err}
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

func CheckToken(token string, user string, pass string, permission string)(hasPrivileges bool, err error){	
	//check token
	masters,err := ndb.GetMasters()
	for x := range masters{
		tkn, err := Encode(masters[x]["secret"])
		if err != nil {
			logs.Error("Error checking Master token: %s", err); return false, err
		}else{
			if token == tkn {
				status,err := UserPrivilegeValidation(user, permission); if err != nil {logs.Error("Privileges error: %s",err); return false,err}
					if status{
						return true,nil
					}else{
						return false,errors.New("This user has not enough privileges level")
					}
			}else{
				return false,errors.New("The token retrieved is false")
			}
		}		
	}
	return false,errors.New("There are no token.")
}