package autentication

import (
    "github.com/astaxie/beego/logs"
    "owlhnode/database"
    "owlhnode/validation"
	"owlhnode/utils"
	"errors"
)

func CreateMasterToken(login map[string]string) (token string, err error) {
	//check master uuid with login uuid
		//EXIST --> create token and send back
		//NO --> insert master uuid
			   //create and insert secret key
			   //insert login uuid
			   //create token and send back

	//check user/pass
	masterExists := false
	users,err := ndb.GetLoginData()
	if err != nil {logs.Error("CreateMasterToken error getting login data: %s", err); return "",errors.New("CreateMasterToken error getting login data")}

	for x := range users{
		hashedPassFromMaster, err := validation.CheckPasswordHash(login["pass"], users[x]["pass"])
		if err != nil {logs.Error("CreateMasterToken Node pass encode error: %s", err); return "",err}

		if login["user"] == users[x]["user"] && hashedPassFromMaster{			
			masters,err := ndb.GetMasters()
			for masterid := range masters {
				if masters[masterid]["master"] == login["master"] && masters[masterid]["login"] == x{
					masterExists = true
					token,err = validation.Encode(masters[masterid]["secret"])
					if err != nil {logs.Error("CreateMasterToken Error generating existing token for master: %s",err); return "",err}
					return token, nil
				}
			}
			if !masterExists{
				//add master into db
				uuid := utils.Generate()
				secret := utils.Generate()
				err = ndb.InsertMaster(uuid, "master", login["master"]); if err != nil {logs.Error("CreateMasterToken Error inserting Master uuid: %s",err); return "",err} 
				err = ndb.InsertMaster(uuid, "secret", secret); if err != nil {logs.Error("CreateMasterToken Error inserting Master secret: %s",err); return "",err} 
				err = ndb.InsertMaster(uuid, "login", x); if err != nil {logs.Error("CreateMasterToken Error inserting Master login credentials: %s",err); return "",err} 

				token,err = validation.Encode(secret)
				if err != nil {logs.Error("CreateMasterToken Error generating token: %s",err); return "",err} 
				return token, nil
			}
		}
	}
	return "", errors.New("CreateMasterToken Incorrect Login credentials")
}

func AddUserFromMaster(masterUser map[string]map[string]string) (err error) {
	nodeUsers, err := ndb.GetLoginData(); if err != nil {logs.Error("AddUserFromMaster Error getting Node users: %s",err); return err} 
	//update all masters to "deleted" status
	for w := range nodeUsers{
		if nodeUsers[w]["type"] == "master"{
			err = ndb.UpdateUsers(w,"status", "deleted")
			if err != nil {logs.Error("AddUserFromMaster Error updating status before update: %s",err); return err} 
		}
	}
	//update users
	nodeUsers, err = ndb.GetLoginData()
	var existsUser bool
	for y := range masterUser{	
		existsUser = false
		for x := range nodeUsers{
			if x == y{existsUser = true}
		}
		if existsUser {
			err = ndb.UpdateUsers(y,"type", masterUser[y]["type"]); if err != nil {logs.Error("AddUserFromMaster Error updating node user type: %s",err); return err} 
			err = ndb.UpdateUsers(y,"user", masterUser[y]["user"]); if err != nil {logs.Error("AddUserFromMaster Error updating node user name: %s",err); return err} 
			err = ndb.UpdateUsers(y,"masterID", masterUser[y]["masterID"]); if err != nil {logs.Error("AddUserFromMaster Error updating node user masterID: %s",err); return err} 
			err = ndb.UpdateUsers(y,"status", masterUser[y]["status"]); if err != nil {logs.Error("AddUserFromMaster Error updating node user status: %s",err); return err} 
		}else{
			err = ndb.InsertUserData(y, "type", masterUser[y]["type"]); if err != nil {logs.Error("AddUserFromMaster Error inserting type: %s",err); return err} 
			err = ndb.InsertUserData(y, "user", masterUser[y]["user"]); if err != nil {logs.Error("AddUserFromMaster Error inserting user: %s",err); return err} 
			err = ndb.InsertUserData(y, "masterID", masterUser[y]["masterID"]); if err != nil {logs.Error("AddUserFromMaster Error inserting master id: %s",err); return err} 
			err = ndb.InsertUserData(y, "status", masterUser[y]["status"]); if err != nil {logs.Error("AddUserFromMaster Error inserting status id: %s",err); return err} 
		}
	}


	return nil
}