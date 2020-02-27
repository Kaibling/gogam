package main

import (
	"fmt"
	"github.com/c-bata/go-prompt"
  "strings"
	"github.com/Kaibling/gogam"
		 "net/http/cookiejar"
		"os"
		"net/url"
)

type clientCli struct {
	jar *cookiejar.Jar
	url string
	username string
	charactername string
	gamename string

}

func completer(d prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{
		{Text: "login"},
		{Text: "user new"},
		{Text: "char new"},
		{Text: "char stats"},
		{Text: "load game"},
		{Text: "game start"},
		{Text: "game new"},
		{Text: "quit"},
	}
	return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
}

func (selfclientCli *clientCli)login(username string) {
	
	data := []byte(username)
	fmt.Println(selfclientCli.url+"login")
	response, err := gogam.PostRequest(selfclientCli.url+"login",data,selfclientCli.jar)
	if err != nil {
		fmt.Println(err)
	}
	if response == "OK" {
		selfclientCli.username = username
	}
	
}

func (selfclientCli *clientCli)sendCommand(command string ) {
	
	data := []byte(command)
	response, err := gogam.PostRequest(selfclientCli.url+"game",data,selfclientCli.jar)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response)
}


func executor(in string) {

	in = strings.TrimSpace(in)

	command := strings.Split(in, " ")
	switch command[0] {
		case "q","quit":
				os.Exit(0)
		case "login":
			if len(command) >= 2 {
				clientObject.login(command[1])
			} else {
				fmt.Println("missing parameter: login <username>")
			}
			
		case "user":
			if command[1] == "new" {
				clientObject.sendCommand(in)
			}
		case "game":
			if command[1] == "new" {
				clientObject.sendCommand(in)
			}
		case "ls":
			fmt.Println(clientObject.url)
				baseURL, err := url.Parse(clientObject.url)
				if err != nil {
					fmt.Println("Malformed URL: ", err.Error())
					return
				}
			fmt.Println(clientObject.jar.Cookies(baseURL))
			fmt.Println(clientObject.username)
		case "":
		default:
			fmt.Println("unknown command")
	}

}

var promtPrefix string
var clientObject *clientCli
func livePrefix() (string, bool) {
	return promtPrefix + "> ", true
}


func main() {

		fmt.Println("preload")
	promtPrefix = ""
	clientObject = new(clientCli)
	
	jar, err := cookiejar.New(&cookiejar.Options{}) 
    if err != nil {
        fmt.Println(err)
	}
	clientObject.jar = jar
	clientObject.url = "http://localhost:7070/"
	fmt.Println("preload finished")
	
	p := prompt.New(
		executor,
		completer,
		prompt.OptionPrefix("sds"+"> "),
		prompt.OptionLivePrefix(livePrefix),
		prompt.OptionPrefixTextColor(prompt.DarkGreen),
		prompt.OptionPreviewSuggestionTextColor(prompt.Blue),
		prompt.OptionSelectedSuggestionBGColor(prompt.LightGray),
		prompt.OptionSuggestionBGColor(prompt.DarkGray),
		prompt.OptionTitle("gogam client"),
	)
	p.Run()


/*
	login
	game -> new
	gmae -> load
	player -> stats
	action -> view env
	action -> interact
	action -> move -> 


	*/
}