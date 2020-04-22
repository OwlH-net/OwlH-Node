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

func UserPermissionsValidation(uuidUser string, permissionRequest string) (val bool, err error) {
	allRelations, err := ndb.GetUserGroupRoles(); if err != nil {logs.Error("UserPermissionsValidation error getting permissions: %s",err); return false, err}
	rolePerm, err := ndb.GetRolePermissions(); if err != nil {logs.Error("UserPermissionsValidation error getting user rolePermissions: %s",err); return false, err}
	allPerm, err := ndb.GetPermissions(); if err != nil {logs.Error("UserPermissionsValidation error getting user GetPermissions: %s",err); return false, err}
	
	logs.Notice(uuidUser+"   --   "+permissionRequest)
	
	//check if permission exists
	permExists := false
	for x := range allPerm{
		if x == permissionRequest {permExists = true}
	}
	if !permExists {logs.Error("Permissions validation error - This permission don't exists"); return false,err}

	logs.Notice("Check all relations")
	for x := range allRelations{
		logs.Info("Searching user: "+allRelations[x]["user"]+"  // Role: "+allRelations[x]["role"]+"  //  Group: "+allRelations[x]["group"])
		if allRelations[x]["user"] == uuidUser{
			logs.Notice(allRelations[x]["user"]+" equals "+uuidUser)
			//Check if user role has admin permissions
			for w := range rolePerm{
				logs.Info("Searching role: "+rolePerm[w]["role"]+"  // Permissions: "+rolePerm[w]["permissions"])
				logs.Notice("Role check: "+rolePerm[w]["role"])
				if allRelations[x]["role"] == rolePerm[w]["role"] && allRelations[x]["role"] != ""{
					//split permissions
					permissionsList := strings.Split(rolePerm[w]["permissions"], ",")
					//check if user is admin
					for p := range permissionsList{
						if permissionsList[p] == "admin" {
							return true, nil
						}					
					}					
					//check for permissions requested
					for s := range permissionsList{
						if permissionsList[s] == permissionRequest {
							return true, nil
						}					
					}					
				}
			}

			//check group roles for user groups
			for s := range allRelations{
				logs.Notice("Group check: "+allRelations[s]["role"])
				if allRelations[s]["group"] == allRelations[x]["group"] && allRelations[x]["group"] != ""{
					for roleID := range rolePerm{
						logs.Notice("Group role check: "+rolePerm[roleID]["role"])					
						if rolePerm[roleID]["role"] == allRelations[s]["role"]{
							//split permissions
							permissionsList := strings.Split(rolePerm[roleID]["permissions"], ",")
							//check if user is admin
							for p := range permissionsList{
								if permissionsList[p] == "admin" {
									return true, nil
								}					
							}					
							//check for permissions requested
							for s := range permissionsList{
								if permissionsList[s] == permissionRequest {
									return true, nil
								}					
							}	
						}
					}
				}
			}
		}
	}

	logs.Error("Permission DENIED!!!!!!")

	return false, nil
}