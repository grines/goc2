package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
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
		fmt.Println(readErr)
	}

	var results []Agent
	jsonErr := json.Unmarshal(body, &results)
	if jsonErr != nil {
		fmt.Println(jsonErr)
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
			readline.PcItem("download",
				readline.PcItemDynamic(listFiles(c2, agent)),
			),
			readline.PcItem("upload"),
			readline.PcItem("agent",
				readline.PcItemDynamic(listAgents(c2)),
			),
			readline.PcItem("cd",
				readline.PcItemDynamic(listFiles(c2, agent)),
			),
			readline.PcItem("cat",
				readline.PcItemDynamic(listFiles(c2, agent)),
			),
			readline.PcItem("mv",
				readline.PcItemDynamic(listFiles(c2, agent)),
			),
			readline.PcItem("cp",
				readline.PcItemDynamic(listFiles(c2, agent)),
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
		case strings.HasPrefix(line, "agent "):
			parts := strings.Split(line, " ")
			agent = parts[1]
		case line == "login":
			pswd, err := l.ReadPassword("please enter your password: ")
			if err != nil {
				break
			}
			println("you enter:", strconv.Quote(string(pswd)))
		case line == "history":
			dat, err := ioutil.ReadFile("/tmp/readline.tmp")
			if err != nil {
				break
			}
			fmt.Print(string(dat))
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
				//uuid := shortuuid.New()
				parts := strings.Split(cmdString, " ")
				file := parts[1]
				//copy(file, "/tmp/"+uuid)
				tempfile := uploadFile(file, c2)
				if tempfile == "NotFound" {
					fmt.Println("File not found")
					break
				}
				cmdString = "upload " + tempfile
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
				break
			}

			cmdid := sendCommand(cmdString, agent, c2)
			deadline := time.Now().Add(15 * time.Second)
			for {
				id, output := getOutput(c2+"/api/cmd/output/"+agent+"/"+cmdid, c2, cmdid)
				if strings.Contains(output, "Location:") {
					parts := strings.Split(output, " ")
					path := parts[1]
					file := filepath.Base(path)
					downloadFile("/tmp/"+file, c2+"/files/"+file)
					//fmt.Println("Process Download" + file)
				}
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
			fmt.Println(readErr)
		}

		var results []Cmd
		jsonErr := json.Unmarshal(body, &results)
		if jsonErr != nil {
			fmt.Println(jsonErr)
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
			fmt.Println(readErr)
		}

		var results Agent
		jsonErr := json.Unmarshal(body, &results)
		if jsonErr != nil {
			fmt.Println(jsonErr)
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

func copy(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}

func downloadFile(filepath string, url string) error {

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("err:", err)
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}

func uploadFile(path string, c2 string) string {
	if _, err := os.Stat(path); err == nil {
		extraParams := map[string]string{
			"operator": "none",
		}
		request, err := newfileUploadRequest(c2+"/api/cmd/files", extraParams, "myFile", path)
		if err != nil {
			fmt.Println(err)
		}
		client := &http.Client{}
		resp, err := client.Do(request)
		if err != nil {
			fmt.Println(err)
		} else {
			body := &bytes.Buffer{}
			_, err := body.ReadFrom(resp.Body)
			if err != nil {
				fmt.Println(err)
			}
			resp.Body.Close()
			return body.String()
		}
		return "Found"
	}
	return "NotFound"
}

// Creates a new file upload http request with optional extra params
func newfileUploadRequest(uri string, params map[string]string, paramName, path string) (*http.Request, error) {
	if _, err := os.Stat(path); err == nil {
		file, err := os.Open(path)
		if err != nil {
			return nil, err
		}
		defer file.Close()

		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		part, err := writer.CreateFormFile(paramName, filepath.Base(path))
		if err != nil {
			return nil, err
		}
		_, err = io.Copy(part, file)

		for key, val := range params {
			_ = writer.WriteField(key, val)
		}
		err = writer.Close()
		if err != nil {
			return nil, err
		}

		req, err := http.NewRequest("POST", uri, body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		return req, err
	}
	return nil, nil
}
