package cli

import (
	"bufio"
	"fmt"
	"os"
	"errors"
	"os/exec"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	for {
		path, err := os.Getwd()
		fmt.Print(path + "-agent0012$ ")
		cmdString, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		err = runCommand(cmdString)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}
}

func runCommand(commandStr string) error {
	commandStr = strings.TrimSuffix(commandStr, "\n")
	arrCommandStr := strings.Fields(commandStr)
	if len(arrCommandStr) < 1 {
		return errors.New("")
	}
	switch arrCommandStr[0] {
	case "cd":
		if len(arrCommandStr) < 1 {
			return errors.New("Required 1 arguments")
		}
		return os.Chdir(arrCommandStr[1])
	case "exit":
		os.Exit(0)
	case "whos":
		cmd := exec.Command("whoami")
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
		//fmt.Fprintln(os.Stdout, output)
		return cmd.Run()
	default:
		cmd := exec.Command(arrCommandStr[0], arrCommandStr[1:]...)
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
		return cmd.Run()
	}
	return nil
}
