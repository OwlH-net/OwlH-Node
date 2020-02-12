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

func AddUserFromMaster(user map[string]string) (err error) {
	logs.Notice(user)
	// uuid := utils.Generate()
	// err = ndb.InsertUserData(uuid, "type", user["type"]); if err != nil {logs.Error("AddUserFromMaster Error inserting type: %s",err); return err} 
	// err = ndb.InsertUserData(uuid, "user", user["user"]); if err != nil {logs.Error("AddUserFromMaster Error inserting user: %s",err); return err} 
	// err = ndb.InsertUserData(uuid, "master", user["master"]); if err != nil {logs.Error("AddUserFromMaster Error inserting master id: %s",err); return err} 

	return nil
}