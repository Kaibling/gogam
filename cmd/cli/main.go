package main

import (
	"encoding/json"
	"fmt"
	"github.com/Kaibling/gogam"
	"github.com/c-bata/go-prompt"
	"net/http/cookiejar"
	//"net/url"
	//"os"
	"net/http"
	"bytes"
	//"strings"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
)

type clientCli struct {
	jar           *cookiejar.Jar
	url           string
	username      string
	charactername string
	gamename      string
}

func completer(d prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{
		{Text: "login", Description: "login <username>: for authentication"},                        //working
		{Text: "user new", Description: "user new <username>: creates new user for the server"},     //working
		{Text: "char new", Description: "char new <CharacterName> <gameID>: creates new Character"}, //working
		{Text: "char list"},
		{Text: "char stats", Description: "char stats <CharacterID>: shows stats and info of char"},
		{Text: "game load", Description: "game load <gameID>: loads game into server"}, //working
		{Text: "game start"},
		{Text: "game join", Description: "game join: joins game"},
		{Text: "game new", Description: "game new <gameName>: creates new game"}, //working
		{Text: "game list", Description: "game list: shows all games on server"}, //working
		{Text: "quit"}, //working
	}
	return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
}

func (selfclientCli *clientCli) login(username string) (error) {

	data := []byte(username)
	fmt.Println(selfclientCli.url + "login")
	response, err := gogam.PostRequest(selfclientCli.url+"login", data, selfclientCli.jar)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response)
	return err
}


func (selfclientCli *clientCli) createUser(userName string) error {
	type requestJSON struct {
		Username string  `json:"username"`
	}
	urlUser := clientObject.url + "user"
	log.Debug("createUser: PostRequest: trying post request to: ", urlUser)


	a := requestJSON {Username:userName}
	
	bytedata,_ :=json.Marshal(&a)
	log.Debug("JSON: ",string(bytedata))
	req, err := http.NewRequest(http.MethodPut, urlUser, bytes.NewBuffer(bytedata))
	if err != nil {
		log.Debug(err)
	}
	//req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Jar: selfclientCli.jar}
	resp, err := client.Do(req)
	if err != nil {
		log.Debug(err)
		return err
	}
	if resp.StatusCode == 201 {
		return nil
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Debug(err)
	}
	resp.Body.Close()
	log.Debug("response Body:", string(body))

	return err
}
/*
func executor(in string) {

	in = strings.TrimSpace(in)
	command := strings.Split(in, " ")
	switch command[0] {
	case "q", "quit":
		os.Exit(0)
	case "login":
		if len(command) >= 2 {
			clientObject.login(command[1])
		} else {
			fmt.Println("missing parameter: login <username>")
		}

	case "user":
		if command[1] == "new" {
			clientObject.createUser(command[2])
		}

	case "game":
		if command[1] == "join" || command[1] == "load" || command[1] == "new" {
			clientObject.sendCommand(in)
		}
		if command[1] == "list" {
			jsonString := clientObject.sendCommand(in)
			var gameNameArray []string
			json.Unmarshal([]byte(jsonString), &gameNameArray)
			for cnt, gameName := range gameNameArray {
				fmt.Println(cnt, ": ", gameName)
			}
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
		
	case "char":
		if command[1] == "new" || command[1] == "stats" {
			clientObject.sendCommand(in)
		}

	case "":
	default:
		fmt.Println("unknown command")
	}

}
*/

var promtPrefix string
var gameID int
var characterName string
var clientObject *clientCli

func livePrefix() (string, bool) {
	return promtPrefix + "> ", true
}

func main() {
	log.SetLevel(log.DebugLevel)
	log.SetReportCaller(true)
	log.Debug("Logging started")

	//fmt.Println("preload")
	promtPrefix = ""
	clientObject = new(clientCli)

	jar, err := cookiejar.New(&cookiejar.Options{})
	if err != nil {
		fmt.Println(err)
	}
	clientObject.jar = jar
	clientObject.url = "http://localhost:7070/"
	//fmt.Println("preload finished")

	err = clientObject.createUser("hans")
	if err != nil {
		fmt.Println(err)
	}
	err = clientObject.login("han")
	if err != nil {
		fmt.Println(err)
	}


	//PostRequest("http://localhost:7070/user",clientObject.jar)
	//type requestJSON struct {
	//	Username string  `json:"username"`
	//}
	//a := requestJSON {Username:"userName"}
	
	//bytedata,_ :=json.Marshal(&a)
	//log.Info("JSON: ",string(bytedata))

	//gogam.PostJSONRequest("http://localhost:7070/user",bytedata)

	//login
	//clientObject.login("admin")
	//create new game
	//clientObject.sendCommand("game new welt1")
	//clientObject.sendCommand("user welt1")
	//clientObject.sendCommand("game new welt2")
	//load new game
	//clientObject.sendCommand("game load 1")
	//new character
	//clientObject.sendCommand("char new char1 1")
	//clientObject.sendCommand("char new char2 1")
	//clientObject.sendCommand("char stats 1")

	//clientObject.sendCommand("char new char2 2")
	//join game
	//clientObject.sendCommand("game join")

	/*

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



		login
		game -> new
		gmae -> load
		player -> stats
		action -> view env
		action -> interact
		action -> move ->

		game status
		- character info
		- map

		login -> game new -> game load -> game join


	*/
}
