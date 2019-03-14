package ndb

import (
    "github.com/astaxie/beego/logs"
    "database/sql"
    _ "github.com/mattn/go-sqlite3"
    "errors"
)


type Node struct {
    NId       string
    NName     string
    NIp       string
    NPort     int
    NType     string
    NUUID     string
}

func AddNode(node Node) (err error){
    logs.Info("DB -> Add Node")
    if Db != nil {
        stmt, err := Db.Prepare("INSERT INTO node(node_name, node_ip, node_port, node_type, node_UUID) values(?,?,?,?,?)")
        if err != nil {
            logs.Error("DB NODE -> Add NODE -> Error db.prepare -> maybe db conn lost? ")
            return errors.New("DB NODE -> Add NODE -> Error db.prepare -> maybe db conn lost? ")
        }
        //Validate fields!!!
        _, err = stmt.Exec(node.NName, node.NIp, node.NPort, node.NType, node.NUUID )
        if err != nil {
            logs.Error("DB NODE -> ADD NODE : %s", err.Error())
            return errors.New("DB NODE -> Query is not working: " + err.Error())
        }
        return nil
    } else {
        logs.Error("DB NODE ->Can't acces to database")
        return errors.New("DB NODE -> Can't acces to database")
    }
}

func GetNode(nid string) (n *Node, err error) {
    logs.Info("DB -> Get Node")
    var node Node
    if Db != nil {
//        rows, err := Db.Query("SELECT * FROM master WHERE master_id=1;")
        row := Db.QueryRow("SELECT * FROM node WHERE node_id=%s;",nid)
        logs.Info ("DB -> Row %s", row)
        err = row.Scan(&node.NId, &node.NName, &node.NIp, &node.NPort, &node.NType, &node.NUUID)
        if err == sql.ErrNoRows {
            logs.Warn("DB NODE -> Can't find %s",nid)
            return nil, errors.New("DB NODE -> Can't find %s",nid)
        }
        if err != nil {
            logs.Warn("DB NODE -> Error reading data GetNode")
            return nil, errors.New("DB NODE -> Error reading data GetNode")
        }
        return &node, nil
    } else {
        logs.Info("DB NODE -> Database not exist")
        return nil, errors.New("DB NODE -> Database not exist")
    }
}
