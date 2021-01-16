package cli

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

var timeoutSetting = 1
var c2 = "http://localhost:8005"
var agent = "test"

//ok
type Cmd struct {
	ID      string
	Command string
	Agent   string
	Status  string
	Cmdid   string
	Output  string
}

func Start() {
	reader := bufio.NewReader(os.Stdin)
	timeout := time.Duration(timeoutSetting) * time.Second
	ticker := time.NewTicker(timeout)
	quit := make(chan struct{})
	for {
		fmt.Print("agent0012$ ")
		select {
		case <-ticker.C:
			cmdString, err := reader.ReadString('\n')
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
			err = sendCommand(cmdString)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
			deadline := time.Now().Add(10 * time.Second)
			for {
				data := getJSON(c2 + "/api/cmd/output/" + agent)
				if data == "True" || cmdString == "\n" {
					break
				}
				if time.Now().After(deadline) {
					fmt.Fprintln(os.Stderr, "Command Timed Out!")
					break
				}
			}
		case <-quit:
			return
		}
	}
}

func sendCommand(cmd string) error {
	resp, err := http.PostForm(c2+"/api/cmd/new",
		url.Values{"cmd": {cmd}})

	if err != nil {
		panic(err)
	}
	if resp.Body != nil {
		defer resp.Body.Close()
	}
	return nil
}

func getJSON(url string) string {

	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	if resp.Body != nil {
		defer resp.Body.Close()
	}

	body, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	var results []Cmd
	jsonErr := json.Unmarshal(body, &results)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	for _, d := range results {
		if len(d.Output) > 0 {
			//fmt.Println(d.Output)
			fmt.Fprintln(os.Stderr, d.Output)
			updateCmdStatus(d.Cmdid)
			return "True"
		}
		//fmt.Println(d.Cmdid + ": " + d.Output)
		//updateCmdStatus(d.Cmdid)
	}

	// Print the HTTP response status.
	//fmt.Println("Response status:", resp.Status)
	return "False"
}

func updateCmdStatus(cmdid string) {
	resp, err := http.PostForm(c2+"/api/cmd/update/output",
		url.Values{"id": {cmdid}, "client_status": {"1"}})

	if err != nil {
		panic(err)
	}
	if resp.Body != nil {
		defer resp.Body.Close()
	}
}
