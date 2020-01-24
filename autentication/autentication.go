package autentication

import (
    "github.com/astaxie/beego/logs"
    "owlhnode/database"
    "owlhnode/validation"
    "owlhnode/utils"
)

func CreateMasterToken() (token string, err error) {
	//use uuid as secret key
	secret := utils.Generate()
	//create token using secret key
	token,err = validation.Encode(secret)
	if err != nil {logs.Error("Error generating token for master: %s",err); return "",err}
	//add token and secret key to DB
	node,err := ndb.GetNodeData()
	for x := range node{
		err = ndb.InsertNodeData(x, "secret", secret); if err != nil {logs.Error("Error inserting master token into DB: %s",err); return "",err}
	}

	//send token to master
	return token, err
}