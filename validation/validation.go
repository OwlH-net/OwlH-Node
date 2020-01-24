package validation

import (
	"github.com/astaxie/beego/logs"
	jwt "github.com/dgrijalva/jwt-go"
	// "golang.org/x/crypto/bcrypt"
	"errors"
    "owlhnode/database"
)

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

func CheckToken(token string)(err error){
	node,err := ndb.GetNodeData()
	for x := range node{
		tkn, err := Encode(node[x]["secret"])
		if err != nil {
			logs.Error("Error checking Master token: %s", err); return err
		}else{
			if token == tkn {
				return nil
			}else{
				return errors.New("The token retrieved is false")
			}
		}		
	}
	return errors.New("There are not token. Error checking token Token")
}