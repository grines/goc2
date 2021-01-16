package web

import (
	"RedMap/internal/app/api"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"text/template"

	"github.com/globalsign/mgo/bson"
	"github.com/julienschmidt/httprouter"
)

type Todo struct {
	Title string
	Done  bool
}

type TodoPageData struct {
	PageTitle string
	ScanID    string
	Todos     []Todo
}

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
	router.GET("/api/reports", apiReports)
	router.GET("/", reports)

	router.GET("/scan", scan)
	router.POST("/api/scan", apiScan)

	//Services Endpoints
	router.GET("/api/services/:name", apiServices)
	router.GET("/services/:name", services)
	router.GET("/services/", redirect)

	//Vulnerabilities Endpoints
	router.GET("/api/vulns/:name", apiVulns)
	router.GET("/vulns/:name", vulns)
	router.GET("/vulns/", redirect)

	//Secrets Endpoints
	router.GET("/api/secrets/:name", apiSecrets)
	router.GET("/secrets/:name", secrets)
	router.GET("/secrets/", redirect)

	//Users Endpoints
	router.GET("/api/users/:name", apiUsers)
	router.GET("/users/:name", users)
	router.GET("/users/", redirect)

	//Sub Endpoints
	router.GET("/api/sub/:name", apiSub)
	router.GET("/sub/:name", sub)
	router.GET("/sub/", redirect)

	fmt.Printf("Starting server at port 8005\n")
	if err := http.ListenAndServe(":8005", router); err != nil {
		log.Fatal(err)
	}
}

func redirect(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	http.Redirect(w, r, "/", 301)
}

func scan(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "text/html")
	tmpl := template.Must(template.ParseFiles("web/template/base.html", "web/template/tmp-scan.html"))
	data := TodoPageData{
		PageTitle: "RedMap - Reports",
	}
	tmpl.ExecuteTemplate(w, "base", data)
}

func apiScan(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}

	domains := r.FormValue("domains")
	git := r.FormValue("git")
	email := r.FormValue("email")
	scanner := r.FormValue("scanner")

	jsond := map[string]interface{}{
		"status": "Scan Started",
	}

	jsondata, err := json.Marshal(jsond)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Fprintf(w, string(jsondata))
	go api.WebScan(domains, git, email, scanner)
}

func apiReports(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	fmt.Fprintf(w, "%s", api.GetReports())
}

func reports(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "text/html")
	tmpl := template.Must(template.ParseFiles("web/template/base.html", "web/template/tmp-reports.html"))
	data := TodoPageData{
		PageTitle: "RedMap - Reports",
		ScanID:    ps.ByName("name"),
	}
	tmpl.ExecuteTemplate(w, "base", data)
}

func apiServices(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	d := api.GetServices(ps.ByName("name"))

	fmt.Fprintf(w, "%s", string(d))
}

func services(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "text/html")
	tmpl := template.Must(template.ParseFiles("web/template/base.html", "web/template/tmp-services.html"))
	data := TodoPageData{
		PageTitle: "RedMap - Services",
		ScanID:    ps.ByName("name"),
		Todos: []Todo{
			{Title: "Task 1", Done: false},
			{Title: "Task 2", Done: true},
			{Title: "Task 3", Done: true},
		},
	}
	tmpl.ExecuteTemplate(w, "base", data)
}

func apiVulns(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	d := api.GetVulns(ps.ByName("name"))

	fmt.Fprintf(w, "%s", string(d))
}

func vulns(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "text/html")
	tmpl := template.Must(template.ParseFiles("web/template/base.html", "web/template/tmp-vulns.html"))
	data := TodoPageData{
		PageTitle: "RedMap - Vulnerabilities",
		ScanID:    ps.ByName("name"),
		Todos: []Todo{
			{Title: "Task 1", Done: false},
			{Title: "Task 2", Done: true},
			{Title: "Task 3", Done: true},
		},
	}
	tmpl.ExecuteTemplate(w, "base", data)
}

func apiSecrets(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	d := api.GetSecrets(ps.ByName("name"))

	fmt.Fprintf(w, "%s", string(d))
}

func secrets(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "text/html")
	tmpl := template.Must(template.ParseFiles("web/template/base.html", "web/template/tmp-secrets.html"))
	data := TodoPageData{
		PageTitle: "RedMap - Secrets",
		ScanID:    ps.ByName("name"),
		Todos: []Todo{
			{Title: "Task 1", Done: false},
			{Title: "Task 2", Done: true},
			{Title: "Task 3", Done: true},
		},
	}
	tmpl.ExecuteTemplate(w, "base", data)
}

func apiUsers(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	d := api.GetEmails(ps.ByName("name"))

	fmt.Fprintf(w, "%s", string(d))
}

func users(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "text/html")
	tmpl := template.Must(template.ParseFiles("web/template/base.html", "web/template/tmp-users.html"))
	data := TodoPageData{
		PageTitle: "RedMap - Users",
		ScanID:    ps.ByName("name"),
		Todos: []Todo{
			{Title: "Task 1", Done: false},
			{Title: "Task 2", Done: true},
			{Title: "Task 3", Done: true},
		},
	}
	tmpl.ExecuteTemplate(w, "base", data)
}

func apiSub(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	d := api.GetSubs(ps.ByName("name"))

	fmt.Fprintf(w, "%s", string(d))
}

func sub(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "text/html")
	tmpl := template.Must(template.ParseFiles("web/template/base.html", "web/template/tmp-sub.html"))
	data := TodoPageData{
		PageTitle: "RedMap - Subdomain Takeovers",
		ScanID:    ps.ByName("name"),
		Todos: []Todo{
			{Title: "Task 1", Done: false},
			{Title: "Task 2", Done: true},
			{Title: "Task 3", Done: true},
		},
	}
	tmpl.ExecuteTemplate(w, "base", data)
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/hello" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	if r.Method != "GET" {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}

	fmt.Fprintf(w, "Hello!")
}

func formHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}
	//fmt.Fprintf(w, "POST request successful")
	//name := r.FormValue("name")
	//address := r.FormValue("address")

	//fmt.Fprintf(w, "Name = %s\n", name)
	//fmt.Fprintf(w, "Address = %s\n", address)
	json := "{'status': 'started'}"
	fmt.Fprintf(w, json)
}
