package autentication

import (
    "errors"
    "github.com/astaxie/beego/logs"
    "owlhnode/database"
    "owlhnode/utils"
    "owlhnode/validation"
)

func CreateMasterToken(login map[string]string) (token string, err error) {
    // master uuid
    // user 'admin'
    // pass 'admin'

    //check user/pass
    masterExists := false
    users, err := ndb.GetLoginData()
    if err != nil {
        logs.Error("CreateMasterToken error getting login data: %s", err)
        return "", errors.New("CreateMasterToken error getting login data")
    }
    masters, err := ndb.GetMasters()
    if err != nil {
        logs.Error("CreateMasterToken error getting master data: %s", err)
        return "", errors.New("CreateMasterToken error getting master data")
    }

    for x := range users {
        hashedPassFromMaster, err := validation.CheckPasswordHash(login["pass"], users[x]["pass"])
        if err != nil {
            continue
        }

        if login["user"] == users[x]["user"] && hashedPassFromMaster {
            for masterid := range masters {
                if masters[masterid]["master"] == login["master"] && masters[masterid]["login"] == x {
                    masterExists = true
                    token, err = validation.Encode(masters[masterid]["secret"])
                    if err != nil {
                        logs.Error("CreateMasterToken Error generating existing token for master: %s", err)
                        return "", err
                    }
                    return token, nil
                }
            }
            if !masterExists {
                //add master into db
                uuid := utils.Generate()
                secret := utils.Generate()
                err = ndb.InsertMaster(uuid, "master", login["master"])
                if err != nil {
                    logs.Error("CreateMasterToken Error inserting Master uuid: %s", err)
                    return "", err
                }
                err = ndb.InsertMaster(uuid, "secret", secret)
                if err != nil {
                    logs.Error("CreateMasterToken Error inserting Master secret: %s", err)
                    return "", err
                }
                err = ndb.InsertMaster(uuid, "login", x)
                if err != nil {
                    logs.Error("CreateMasterToken Error inserting Master login credentials: %s", err)
                    return "", err
                }

                token, err = validation.Encode(secret)
                if err != nil {
                    logs.Error("CreateMasterToken Error generating token: %s", err)
                    return "", err
                }
                return token, nil
            }
        }
    }
    return "", errors.New("CreateMasterToken Incorrect Login credentials")
}

func AddUserFromMaster(masterUser map[string]map[string]string) (err error) {
    nodeUsers, err := ndb.GetLoginData()
    if err != nil {
        logs.Error("AddUserFromMaster Error getting Node users: %s", err)
        return err
    }
    //update all masters to "deleted" status
    for w := range nodeUsers {
        if nodeUsers[w]["type"] == "master" {
            err = ndb.UpdateUsers(w, "status", "deleted")
            if err != nil {
                logs.Error("AddUserFromMaster Error updating status before update: %s", err)
                return err
            }
        }
    }
    //update users
    nodeUsers, err = ndb.GetLoginData()
    var existsUser bool
    for y := range masterUser {
        existsUser = false
        for x := range nodeUsers {
            if x == y {
                existsUser = true
            }
        }
        if existsUser {
            err = ndb.UpdateUsers(y, "type", masterUser[y]["type"])
            if err != nil {
                logs.Error("AddUserFromMaster Error updating node user type: %s", err)
                return err
            }
            err = ndb.UpdateUsers(y, "user", masterUser[y]["user"])
            if err != nil {
                logs.Error("AddUserFromMaster Error updating node user name: %s", err)
                return err
            }
            err = ndb.UpdateUsers(y, "masterID", masterUser[y]["masterID"])
            if err != nil {
                logs.Error("AddUserFromMaster Error updating node user masterID: %s", err)
                return err
            }
            err = ndb.UpdateUsers(y, "status", masterUser[y]["status"])
            if err != nil {
                logs.Error("AddUserFromMaster Error updating node user status: %s", err)
                return err
            }
        } else {
            err = ndb.InsertUserData(y, "type", masterUser[y]["type"])
            if err != nil {
                logs.Error("AddUserFromMaster Error inserting type: %s", err)
                return err
            }
            err = ndb.InsertUserData(y, "user", masterUser[y]["user"])
            if err != nil {
                logs.Error("AddUserFromMaster Error inserting user: %s", err)
                return err
            }
            err = ndb.InsertUserData(y, "masterID", masterUser[y]["masterID"])
            if err != nil {
                logs.Error("AddUserFromMaster Error inserting master id: %s", err)
                return err
            }
            err = ndb.InsertUserData(y, "status", masterUser[y]["status"])
            if err != nil {
                logs.Error("AddUserFromMaster Error inserting status id: %s", err)
                return err
            }
        }
    }

    return nil
}

func AddRolesFromMaster(masterRole map[string]map[string]string) (err error) {
    nodeRoles, err := ndb.GetUserRole()
    if err != nil {
        logs.Error("AddRolesFromMaster Error getting Node roles: %s", err)
        return err
    }
    //update all masters to "deleted" status
    for w := range nodeRoles {
        if nodeRoles[w]["type"] == "master" {
            err = ndb.UpdateUsers(w, "status", "deleted")
            if err != nil {
                logs.Error("AddRolesFromMaster Error updating status before update: %s", err)
                return err
            }
        }
    }
    //update users
    nodeRoles, err = ndb.GetUserRole()
    var existsRole bool
    for y := range masterRole {
        existsRole = false
        for x := range nodeRoles {
            if x == y {
                existsRole = true
            }
        }
        if existsRole {
            err = ndb.UpdateUserRoles(y, "type", masterRole[y]["type"])
            if err != nil {
                logs.Error("AddRolesFromMaster Error updating node role type: %s", err)
                return err
            }
            err = ndb.UpdateUserRoles(y, "role", masterRole[y]["role"])
            if err != nil {
                logs.Error("AddRolesFromMaster Error updating node role name: %s", err)
                return err
            }
            // err = ndb.UpdateUserRoles(y,"permissions", masterRole[y]["permissions"]); if err != nil {logs.Error("AddRolesFromMaster Error updating node permissions name: %s",err); return err}
            err = ndb.UpdateUserRoles(y, "masterID", masterRole[y]["masterID"])
            if err != nil {
                logs.Error("AddRolesFromMaster Error updating node role masterID: %s", err)
                return err
            }
            err = ndb.UpdateUserRoles(y, "status", masterRole[y]["status"])
            if err != nil {
                logs.Error("AddRolesFromMaster Error updating node role status: %s", err)
                return err
            }
        } else {
            err = ndb.InsertUserRole(y, "type", masterRole[y]["type"])
            if err != nil {
                logs.Error("AddRolesFromMaster Error inserting type: %s", err)
                return err
            }
            err = ndb.InsertUserRole(y, "role", masterRole[y]["role"])
            if err != nil {
                logs.Error("AddRolesFromMaster Error inserting role: %s", err)
                return err
            }
            // err = ndb.InsertUserRole(y, "permissions", masterRole[y]["permissions"]); if err != nil {logs.Error("AddRolesFromMaster Error inserting permissions: %s",err); return err}
            err = ndb.InsertUserRole(y, "masterID", masterRole[y]["masterID"])
            if err != nil {
                logs.Error("AddRolesFromMaster Error inserting master id: %s", err)
                return err
            }
            err = ndb.InsertUserRole(y, "status", masterRole[y]["status"])
            if err != nil {
                logs.Error("AddRolesFromMaster Error inserting status id: %s", err)
                return err
            }
        }
    }

    return nil
}

func AddGroupFromMaster(masterGroup map[string]map[string]string) (err error) {
    nodeGroups, err := ndb.GetUserGroup()
    if err != nil {
        logs.Error("AddGroupFromMaster Error getting Node groups: %s", err)
        return err
    }
    //update all masters to "deleted" status
    for w := range nodeGroups {
        if nodeGroups[w]["type"] == "master" {
            err = ndb.UpdateUserGroup(w, "status", "deleted")
            if err != nil {
                logs.Error("AddGroupFromMaster Error updating status before update: %s", err)
                return err
            }
        }
    }
    //update groups
    nodeGroups, err = ndb.GetUserGroup()
    var existsGroup bool
    for y := range masterGroup {
        existsGroup = false
        for x := range nodeGroups {
            if x == y {
                existsGroup = true
            }
        }
        if existsGroup {
            err = ndb.UpdateUserGroup(y, "type", masterGroup[y]["type"])
            if err != nil {
                logs.Error("AddGroupFromMaster Error updating node group type: %s", err)
                return err
            }
            err = ndb.UpdateUserGroup(y, "group", masterGroup[y]["group"])
            if err != nil {
                logs.Error("AddGroupFromMaster Error updating node group name: %s", err)
                return err
            }
            err = ndb.UpdateUserGroup(y, "masterID", masterGroup[y]["masterID"])
            if err != nil {
                logs.Error("AddGroupFromMaster Error updating node group masterID: %s", err)
                return err
            }
            err = ndb.UpdateUserGroup(y, "status", masterGroup[y]["status"])
            if err != nil {
                logs.Error("AddGroupFromMaster Error updating node group status: %s", err)
                return err
            }
        } else {
            err = ndb.InsertUserGroup(y, "type", masterGroup[y]["type"])
            if err != nil {
                logs.Error("AddGroupFromMaster Error inserting type: %s", err)
                return err
            }
            err = ndb.InsertUserGroup(y, "group", masterGroup[y]["group"])
            if err != nil {
                logs.Error("AddGroupFromMaster Error inserting group: %s", err)
                return err
            }
            err = ndb.InsertUserGroup(y, "masterID", masterGroup[y]["masterID"])
            if err != nil {
                logs.Error("AddGroupFromMaster Error inserting master id: %s", err)
                return err
            }
            err = ndb.InsertUserGroup(y, "status", masterGroup[y]["status"])
            if err != nil {
                logs.Error("AddGroupFromMaster Error inserting status id: %s", err)
                return err
            }
        }
    }

    return nil
}

func AddUserGroupRolesFromMaster(masterUgr map[string]map[string]string) (err error) {
    nodeUGR, err := ndb.GetUserGroupRoles()
    if err != nil {
        logs.Error("AddUserGroupRolesFromMaster Error getting Node userGroupRoles: %s", err)
        return err
    }
    //update all masters to "deleted" status
    for w := range nodeUGR {
        if nodeUGR[w]["type"] == "master" {
            err = ndb.UpdateUserGroupRoles(w, "status", "deleted")
            if err != nil {
                logs.Error("AddUserGroupRolesFromMaster Error updating status before update: %s", err)
                return err
            }
        }
    }
    //update groups
    nodeUGR, err = ndb.GetUserGroupRoles()
    var existsUGR bool
    for y := range masterUgr {
        existsUGR = false
        for x := range nodeUGR {
            if x == y {
                existsUGR = true
            }
        }
        if existsUGR {
            if masterUgr[y]["user"] != "" {
                err = ndb.UpdateUserGroupRoles(y, "user", masterUgr[y]["user"])
                if err != nil {
                    logs.Error("AddUserGroupRolesFromMaster Error updating node user name: %s", err)
                    return err
                }
            }
            if masterUgr[y]["role"] != "" {
                err = ndb.UpdateUserGroupRoles(y, "role", masterUgr[y]["role"])
                if err != nil {
                    logs.Error("AddUserGroupRolesFromMaster Error updating node role name: %s", err)
                    return err
                }
            }
            if masterUgr[y]["group"] != "" {
                err = ndb.UpdateUserGroupRoles(y, "group", masterUgr[y]["group"])
                if err != nil {
                    logs.Error("AddUserGroupRolesFromMaster Error updating node group name: %s", err)
                    return err
                }
            }
            err = ndb.UpdateUserGroupRoles(y, "type", masterUgr[y]["type"])
            if err != nil {
                logs.Error("AddUserGroupRolesFromMaster Error updating node group type: %s", err)
                return err
            }
            err = ndb.UpdateUserGroupRoles(y, "masterID", masterUgr[y]["masterID"])
            if err != nil {
                logs.Error("AddUserGroupRolesFromMaster Error updating node group masterID: %s", err)
                return err
            }
            err = ndb.UpdateUserGroupRoles(y, "status", masterUgr[y]["status"])
            if err != nil {
                logs.Error("AddUserGroupRolesFromMaster Error updating node group status: %s", err)
                return err
            }
        } else {
            if masterUgr[y]["user"] != "" {
                err = ndb.InsertUserGroupRole(y, "user", masterUgr[y]["user"])
                if err != nil {
                    logs.Error("AddUserGroupRolesFromMaster Error inserting node user name: %s", err)
                    return err
                }
            }
            if masterUgr[y]["role"] != "" {
                err = ndb.InsertUserGroupRole(y, "role", masterUgr[y]["role"])
                if err != nil {
                    logs.Error("AddUserGroupRolesFromMaster Error inserting node role name: %s", err)
                    return err
                }
            }
            if masterUgr[y]["group"] != "" {
                err = ndb.InsertUserGroupRole(y, "group", masterUgr[y]["group"])
                if err != nil {
                    logs.Error("AddUserGroupRolesFromMaster Error inserting node group name: %s", err)
                    return err
                }
            }
            err = ndb.InsertUserGroupRole(y, "type", masterUgr[y]["type"])
            if err != nil {
                logs.Error("AddUserGroupRolesFromMaster Error inserting type: %s", err)
                return err
            }
            err = ndb.InsertUserGroupRole(y, "masterID", masterUgr[y]["masterID"])
            if err != nil {
                logs.Error("AddUserGroupRolesFromMaster Error inserting master id: %s", err)
                return err
            }
            err = ndb.InsertUserGroupRole(y, "status", masterUgr[y]["status"])
            if err != nil {
                logs.Error("AddUserGroupRolesFromMaster Error inserting status id: %s", err)
                return err
            }
        }
    }

    return nil
}

func SyncRolePermissions(masterRolePerm map[string]map[string]string) (err error) {
    nodeRolePerm, err := ndb.GetRolePermissions()
    if err != nil {
        logs.Error("SyncRolePermissions Error getting Node userGroupRoles: %s", err)
        return err
    }
    //update all masters to "deleted" status
    for w := range nodeRolePerm {
        if nodeRolePerm[w]["type"] == "master" {
            err = ndb.UpdateUserGroupRoles(w, "status", "deleted")
            if err != nil {
                logs.Error("SyncRolePermissions Error updating status before update: %s", err)
                return err
            }
        }
    }

    //update or insert all elements
    var existsValue bool
    for y := range masterRolePerm {
        existsValue = false
        for x := range nodeRolePerm {
            if x == y {
                existsValue = true
            }
        }
        if existsValue {
            err = ndb.UpdateRolePermissions(y, "masterID", masterRolePerm[y]["masterID"])
            if err != nil {
                logs.Error("SyncRolePermissions Error updating node group masterID: %s", err)
                return err
            }
            err = ndb.UpdateRolePermissions(y, "type", masterRolePerm[y]["type"])
            if err != nil {
                logs.Error("SyncRolePermissions Error updating node group type: %s", err)
                return err
            }
            err = ndb.UpdateRolePermissions(y, "status", masterRolePerm[y]["status"])
            if err != nil {
                logs.Error("SyncRolePermissions Error updating node group status: %s", err)
                return err
            }
            err = ndb.UpdateRolePermissions(y, "role", masterRolePerm[y]["role"])
            if err != nil {
                logs.Error("SyncRolePermissions Error updating node role: %s", err)
                return err
            }
            err = ndb.UpdateRolePermissions(y, "permissions", masterRolePerm[y]["permissions"])
            if err != nil {
                logs.Error("SyncRolePermissions Error updating node permissions: %s", err)
                return err
            }
            err = ndb.UpdateRolePermissions(y, "object", masterRolePerm[y]["object"])
            if err != nil {
                logs.Error("SyncRolePermissions Error updating node object: %s", err)
                return err
            }
        } else {
            err = ndb.InsertRolePermissions(y, "masterID", masterRolePerm[y]["masterID"])
            if err != nil {
                logs.Error("SyncRolePermissions Error insert node group masterID: %s", err)
                return err
            }
            err = ndb.InsertRolePermissions(y, "type", masterRolePerm[y]["type"])
            if err != nil {
                logs.Error("SyncRolePermissions Error insert node group type: %s", err)
                return err
            }
            err = ndb.InsertRolePermissions(y, "status", masterRolePerm[y]["status"])
            if err != nil {
                logs.Error("SyncRolePermissions Error insert node group status: %s", err)
                return err
            }
            err = ndb.InsertRolePermissions(y, "role", masterRolePerm[y]["role"])
            if err != nil {
                logs.Error("SyncRolePermissions Error insert node role: %s", err)
                return err
            }
            err = ndb.InsertRolePermissions(y, "permissions", masterRolePerm[y]["permissions"])
            if err != nil {
                logs.Error("SyncRolePermissions Error insert node permissions: %s", err)
                return err
            }
            err = ndb.InsertRolePermissions(y, "object", masterRolePerm[y]["object"])
            if err != nil {
                logs.Error("SyncRolePermissions Error insert node object: %s", err)
                return err
            }
        }
    }

    return nil
}

func SyncPermissions(masterPerm map[string]map[string]string) (err error) {
    nodePerm, err := ndb.GetPermissions()
    if err != nil {
        logs.Error("SyncPermissions Error getting Node userGroupRoles: %s", err)
        return err
    }
    //update all masters to "deleted" status
    for w := range nodePerm {
        if nodePerm[w]["type"] == "master" {
            err = ndb.UpdatePermissions(w, "status", "deleted")
            if err != nil {
                logs.Error("SyncPermissions Error updating status before update: %s", err)
                return err
            }
        }
    }

    //update or insert all elements
    var existsValue bool
    for y := range masterPerm {
        existsValue = false
        for x := range nodePerm {
            if x == y {
                existsValue = true
            }
        }
        if existsValue {
            err = ndb.UpdatePermissions(y, "masterID", masterPerm[y]["masterID"])
            if err != nil {
                logs.Error("SyncPermissions Error updating node group masterID: %s", err)
                return err
            }
            err = ndb.UpdatePermissions(y, "type", masterPerm[y]["type"])
            if err != nil {
                logs.Error("SyncPermissions Error updating node group type: %s", err)
                return err
            }
            err = ndb.UpdatePermissions(y, "status", masterPerm[y]["status"])
            if err != nil {
                logs.Error("SyncPermissions Error updating node group status: %s", err)
                return err
            }
            err = ndb.UpdatePermissions(y, "desc", masterPerm[y]["desc"])
            if err != nil {
                logs.Error("SyncPermissions Error updating node desc: %s", err)
                return err
            }
            err = ndb.UpdatePermissions(y, "permisionGroup", masterPerm[y]["permisionGroup"])
            if err != nil {
                logs.Error("SyncPermissions Error updating node permisionGroup: %s", err)
                return err
            }
        } else {
            err = ndb.InsertPermissions(y, "masterID", masterPerm[y]["masterID"])
            if err != nil {
                logs.Error("SyncPermissions Error insert node group masterID: %s", err)
                return err
            }
            err = ndb.InsertPermissions(y, "type", masterPerm[y]["type"])
            if err != nil {
                logs.Error("SyncPermissions Error insert node group type: %s", err)
                return err
            }
            err = ndb.InsertPermissions(y, "status", masterPerm[y]["status"])
            if err != nil {
                logs.Error("SyncPermissions Error insert node group status: %s", err)
                return err
            }
            err = ndb.InsertPermissions(y, "desc", masterPerm[y]["desc"])
            if err != nil {
                logs.Error("SyncPermissions Error insert node desc: %s", err)
                return err
            }
            err = ndb.InsertPermissions(y, "permisionGroup", masterPerm[y]["permisionGroup"])
            if err != nil {
                logs.Error("SyncPermissions Error insert node permisionGroup: %s", err)
                return err
            }
        }
    }

    return nil
}

func SyncRoleGroups(masterPerm map[string]map[string]string) (err error) {
    nodePerm, err := ndb.GetRoleGroups()
    if err != nil {
        logs.Error("SyncRoleGroups Error getting Node userGroupRoles: %s", err)
        return err
    }
    //update all masters to "deleted" status
    for w := range nodePerm {
        if nodePerm[w]["type"] == "master" {
            err = ndb.UpdateRoleGroups(w, "status", "deleted")
            if err != nil {
                logs.Error("SyncRoleGroups Error updating status before update: %s", err)
                return err
            }
        }
    }

    //update or insert all elements
    var existsValue bool
    for y := range masterPerm {
        existsValue = false
        for x := range nodePerm {
            if x == y {
                existsValue = true
            }
        }
        if existsValue {
            err = ndb.UpdateRoleGroups(y, "masterID", masterPerm[y]["masterID"])
            if err != nil {
                logs.Error("SyncRoleGroups Error updating node group masterID: %s", err)
                return err
            }
            err = ndb.UpdateRoleGroups(y, "type", masterPerm[y]["type"])
            if err != nil {
                logs.Error("SyncRoleGroups Error updating node group type: %s", err)
                return err
            }
            err = ndb.UpdateRoleGroups(y, "status", masterPerm[y]["status"])
            if err != nil {
                logs.Error("SyncRoleGroups Error updating node group status: %s", err)
                return err
            }
            err = ndb.UpdateRoleGroups(y, "group", masterPerm[y]["group"])
            if err != nil {
                logs.Error("SyncRoleGroups Error updating node group: %s", err)
                return err
            }
            err = ndb.UpdateRoleGroups(y, "desc", masterPerm[y]["desc"])
            if err != nil {
                logs.Error("SyncRoleGroups Error updating node desc: %s", err)
                return err
            }
        } else {
            err = ndb.InsertRoleGroups(y, "masterID", masterPerm[y]["masterID"])
            if err != nil {
                logs.Error("SyncRoleGroups Error insert node group masterID: %s", err)
                return err
            }
            err = ndb.InsertRoleGroups(y, "type", masterPerm[y]["type"])
            if err != nil {
                logs.Error("SyncRoleGroups Error insert node group type: %s", err)
                return err
            }
            err = ndb.InsertRoleGroups(y, "status", masterPerm[y]["status"])
            if err != nil {
                logs.Error("SyncRoleGroups Error insert node group status: %s", err)
                return err
            }
            err = ndb.InsertRoleGroups(y, "desc", masterPerm[y]["desc"])
            if err != nil {
                logs.Error("SyncRoleGroups Error insert node desc: %s", err)
                return err
            }
            err = ndb.InsertRoleGroups(y, "group", masterPerm[y]["group"])
            if err != nil {
                logs.Error("SyncRoleGroups Error insert node group: %s", err)
                return err
            }
        }
    }

    return nil
}
