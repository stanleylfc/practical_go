package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strings"
)

func init() {
	go func() {
		sigRecv := make(chan os.Signal, 1)
		//sigs := []os.Signal{syscall.SIGINT, syscall.SIGQUIT}
		signal.Notify(sigRecv)
		for range sigRecv {

		}
	}()
}

func main() {
	for {
		inputReader := bufio.NewReader(os.Stdin)
		fmt.Print("[shell]:")

		input, err := inputReader.ReadString('\n')

		if err != nil {
			log.Fatal(err)
		}

		input = strings.TrimSpace(input)

		if input == "q" || input == "exit" {
			fmt.Println("byebye")
			return
		}

		strs := strings.Split(input, "|")

		cmds := []*exec.Cmd{}

		for _, str := range strs {
			s := strings.Fields(str)
			cmds = append(cmds, exec.Command(s[0], s[1:]...))
		}

		var lenCmd = len(cmds) - 1
		for i := range cmds {

			if i > 0 {
				stdout, err := cmds[i-1].StdoutPipe()
				if err != nil {
					fmt.Fprintln(os.Stderr, err)
					break
				}
				cmds[i].Stdin = stdout
			}

			if i == lenCmd {
				cmds[i].Stdout = os.Stdout
			}
		}

		for _, cmd := range cmds {
			if err := cmd.Start(); err != nil {
				fmt.Fprintln(os.Stderr, err)
				break
			}
		}

		for _, cmd := range cmds {
			if err := cmd.Wait(); err != nil {
				fmt.Fprintln(os.Stderr, err)
				break
			}
		}
	}

}
