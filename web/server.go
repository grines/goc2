package web

import (
	"goc2/internal/app/api"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	//"text/template"

	"github.com/globalsign/mgo/bson"
	"github.com/julienschmidt/httprouter"
)

type domainObject struct {
	ID       bson.ObjectId `bson:"_id,omitempty"`
	Type     string
	Hostname string
	Domain   string
	Private  bool
	Ipv4     string
}

//Start the Web Server
func Start() {
	router := httprouter.New()

	//Main Entry
	router.POST("/api/scan", apiScan)

	//Main Entry
	router.GET("/api/test", apiTest)

	//Agents Endpoints
	router.GET("/api/agents", apiAgents)
	router.GET("/agents/", redirect)

	//commands Endpoints
	router.GET("/api/cmds/:name", apiCmds)
	router.GET("/cmds/", redirect)

	fmt.Printf("Starting server at port 8005\n")
	if err := http.ListenAndServe(":8005", router); err != nil {
		log.Fatal(err)
	}
}

func redirect(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	http.Redirect(w, r, "/", 301)
}

func apiScan(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}

	//domains := r.FormValue("domains")
	//git := r.FormValue("git")
	//email := r.FormValue("email")
	//scanner := r.FormValue("scanner")

	jsond := map[string]interface{}{
		"status": "Scan Started",
	}

	jsondata, err := json.Marshal(jsond)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Fprintf(w, string(jsondata))
	//go api.WebScan(domains, git, email, scanner)
}

func apiAgents(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	d := api.GetAgents()

	fmt.Fprintf(w, "%s", string(d))
}


func apiCmds(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	d := api.GetCommands(ps.ByName("name"))

	fmt.Fprintf(w, "%s", string(d))
}

func apiTest(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	//fmt.Fprintf(w, "POST request successful")
	//name := r.FormValue("name")
	//address := r.FormValue("address")

	//fmt.Fprintf(w, "Name = %s\n", name)
	//fmt.Fprintf(w, "Address = %s\n", address)
	json := "{\"status\": \"started\"}"
	fmt.Fprintf(w, json)
}
