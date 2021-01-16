package api

import (
	"encoding/json"
	"log"
	"fmt"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

type Agent struct {
	ID        bson.ObjectId `bson:"_id,omitempty"`
	Name      string
}

type Command struct {
	ID        bson.ObjectId `bson:"_id,omitempty"`
	Command      string
	Agent      string
	Status      string
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
	err = c.Find(bson.M{"agent": "test"}).All(&results)

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