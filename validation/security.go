package validation

import (
	"github.com/astaxie/beego/logs"
	// jwt "github.com/dgrijalva/jwt-go"
	// "golang.org/x/crypto/bcrypt"
	// "errors"
	"strings"
    "owlhnode/database"
    // "owlhmaster/utils"
)

// Check user privileges
func UserPrivilegeValidation(uuidUser string, requestType string) (val bool, err error) {
	allRelations, err := ndb.GetUserGroupRoles(); if err != nil {logs.Error("UserPrivilegeValidation error getting permissions: %s",err); return false, err}
	roles, err := ndb.GetUserRole(); if err != nil {logs.Error("UserPrivilegeValidation error getting user roles: %s",err); return false, err}
	for x := range allRelations{
		if allRelations[x]["user"] == uuidUser{
			//Compare with role permissions
			allPermissionsRole := strings.Split(roles[allRelations[x]["role"]]["permissions"], ",")
			for r := range allPermissionsRole{
				if allPermissionsRole[r] == requestType{
					return true, nil
				}
			}	
				
			//Compare with role permissions for groups
			for y := range allRelations{
				if allRelations[x]["group"] == allRelations[y]["group"]{
					allPermissionsRole = strings.Split(roles[allRelations[y]["role"]]["permissions"], ",")
					for r := range allPermissionsRole{
						if allPermissionsRole[r] == requestType{
							return true, nil
						}
					}	
				}
			}
		}
	}

	return false, nil
}