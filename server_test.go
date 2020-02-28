package gogam

import (
	"io/ioutil"
	 log "github.com/sirupsen/logrus"
	"net/http"
	"net/http/httptest"
	"testing"
	

	//"strings"
	"bytes"
)

/*
		{Text: "user new",Description: "user new <username>: creates new user for the server"},		//working
		{Text: "char new",Description: "char new <CharacterName> <gameID>: creates new Character"},	//working
		{Text: "char list"},
		{Text: "game load",Description: "game load <gameID>: loads game into server"},				//working
		{Text: "game start"},
		{Text: "game join",Description: "game join: joins game"},
		{Text: "game list",Description: "game list: shows all games on server"},					//working
*/


func postRequestGameHandler(command string,url string,testserver *Server,res *httptest.ResponseRecorder) ([]byte, error) {
	a := bytes.NewBuffer([]byte(command))
	req, _ := http.NewRequest("POST", url, a)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(res.Result().Cookies()[0])
	testServer.gameHandler(res, req)
 	return ioutil.ReadAll(res.Body)
}



func TestLoginFailing(t *testing.T) {
	res := httptest.NewRecorder()
	a := bytes.NewBuffer([]byte("login"))
	req, _ := http.NewRequest("POST", "/login", a)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	testServer := new(Server)
	testServer.initiateServer()
	testServer.loginHandler(res, req)

	content, _ := ioutil.ReadAll(res.Body)
	expected := "ynogo"
	if string(content) != expected {
		t.Errorf("Expected %s, got %s.", expected, string(content))
	}
}

func TestLoginSuccess(t *testing.T) {

	res := httptest.NewRecorder()
	a := bytes.NewBuffer([]byte("admin"))
	req, _ := http.NewRequest("POST", "/login", a)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	testServer := new(Server)
	testServer.initiateServer()
	testServer.loginHandler(res, req)

	content, _ := ioutil.ReadAll(res.Body)
	expected := "OK"
	if string(content) != expected {
		t.Errorf("Expected %s, got %s.", expected, string(content))
	}
}

func TestNewGameUnautheticated(t *testing.T) {

	res := httptest.NewRecorder()
	testServer := new(Server)
	testServer.initiateServer()

	//test game new (unauthenticated)
	a := bytes.NewBuffer([]byte("game new"))
	req, _ := http.NewRequest("POST", "/game", a)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	testServer.gameHandler(res, req)
	content, _ := ioutil.ReadAll(res.Body)
	expected := "ynogo"
	if string(content) != expected {
		t.Errorf("Expected %s, got %s.", expected, string(content))
	}
}
func TestNewGameautheticatednoGameID(t *testing.T) {
	//test game new
	res := httptest.NewRecorder()
	testServer := new(Server)
	testServer.initiateServer()

	a := bytes.NewBuffer([]byte("admin"))
	req, _ := http.NewRequest("POST", "/login", a)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	testServer.loginHandler(res, req)
 	ioutil.ReadAll(res.Body)


	a = bytes.NewBuffer([]byte("game new"))
	req, _ = http.NewRequest("POST", "/game", a)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(res.Result().Cookies()[0])

	testServer.gameHandler(res, req)
	content, _ := ioutil.ReadAll(res.Body)
	expected := "not enough arguments"
	if string(content) != expected {
		t.Errorf("Expected %s, got %s.", expected, string(content))
	}
}

func TestNewGameautheticatednewGame(t *testing.T) {
	//test game new world1
	res := httptest.NewRecorder()
	testServer := new(Server)
	testServer.initiateServer()

	a := bytes.NewBuffer([]byte("admin"))
	req, _ := http.NewRequest("POST", "/login", a)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	testServer.loginHandler(res, req)
 	ioutil.ReadAll(res.Body)


	a = bytes.NewBuffer([]byte("game new world1"))
	req, _ = http.NewRequest("POST", "/game", a)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(res.Result().Cookies()[0])

	testServer.gameHandler(res, req)
	content, _ := ioutil.ReadAll(res.Body)
	expected := "OK"
	if string(content) != expected {
		t.Errorf("Expected %s, got %s.", expected, string(content))
	}
}


func TestGameLoadNonExistingGame(t *testing.T) {
	//test game new world1
	res := httptest.NewRecorder()
	testServer := new(Server)
	testServer.initiateServer()

	a := bytes.NewBuffer([]byte("admin"))
	req, _ := http.NewRequest("POST", "/login", a)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	testServer.loginHandler(res, req)
 	ioutil.ReadAll(res.Body)


	a = bytes.NewBuffer([]byte("game new world1"))
	req, _ = http.NewRequest("POST", "/game", a)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(res.Result().Cookies()[0])
	ioutil.ReadAll(res.Body)

	a = bytes.NewBuffer([]byte("game load 1"))
	req, _ = http.NewRequest("POST", "/game", a)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(res.Result().Cookies()[0])
	testServer.gameHandler(res, req)

	content, _ := postRequestGameHandler("game load 1","/game",testserver,res)

	//content, _ := ioutil.ReadAll(res.Body)
	expected := "OK"
	if string(content) != expected {
		t.Errorf("Expected %s, got %s.", expected, string(content))
	}
}

func TestGameLoadExistingGame(t *testing.T) {
	//test game new world1
	res := httptest.NewRecorder()
	testServer := new(Server)
	testServer.initiateServer()

	a := bytes.NewBuffer([]byte("admin"))
	req, _ := http.NewRequest("POST", "/login", a)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	testServer.loginHandler(res, req)
 	ioutil.ReadAll(res.Body)


	a = bytes.NewBuffer([]byte("game new world1"))
	req, _ = http.NewRequest("POST", "/game", a)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(res.Result().Cookies()[0])
	ioutil.ReadAll(res.Body)

	a = bytes.NewBuffer([]byte("game load 0"))
	req, _ = http.NewRequest("POST", "/game", a)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(res.Result().Cookies()[0])

	testServer.gameHandler(res, req)
	content, _ := ioutil.ReadAll(res.Body)
	expected := "game loaded"
	if string(content) != expected {
		t.Errorf("Expected %s, got %s.", expected, string(content))
	}
}

func TestGameLoadExistingGameTwice(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	//test game new world1
	res := httptest.NewRecorder()
	testServer := new(Server)
	testServer.initiateServer()

	a := bytes.NewBuffer([]byte("admin"))
	req, _ := http.NewRequest("POST", "/login", a)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	testServer.loginHandler(res, req)
 	ioutil.ReadAll(res.Body)


	a = bytes.NewBuffer([]byte("game new world1"))
	req, _ = http.NewRequest("POST", "/game", a)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(res.Result().Cookies()[0])
	ioutil.ReadAll(res.Body)

	a = bytes.NewBuffer([]byte("game load 0"))
	req, _ = http.NewRequest("POST", "/game", a)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(res.Result().Cookies()[0])
	ioutil.ReadAll(res.Body)

	a = bytes.NewBuffer([]byte("game load 0"))
	req, _ = http.NewRequest("POST", "/game", a)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(res.Result().Cookies()[0])
	ioutil.ReadAll(res.Body)

		a = bytes.NewBuffer([]byte("game load 0"))
	req, _ = http.NewRequest("POST", "/game", a)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(res.Result().Cookies()[0])
	ioutil.ReadAll(res.Body)

	testServer.gameHandler(res, req)
	content, _ := ioutil.ReadAll(res.Body)
	expected := "game loaded"
	if string(content) != expected {
		t.Errorf("Expected %s, got %s.", expected, string(content))
	}
}

/*
		{Text: "login",Description: "login <username>: for authentication"}, 						//working
		{Text: "user new",Description: "user new <username>: creates new user for the server"},		//working
		{Text: "char new",Description: "char new <CharacterName> <gameID>: creates new Character"},	//working
		{Text: "char list"},
		{Text: "game load",Description: "game load <gameID>: loads game into server"},				//working
		{Text: "game start"},
		{Text: "game join",Description: "game join: joins game"},
		{Text: "game new",Description: "game new <gameName>: creates new game"},					//working
		{Text: "game list",Description: "game list: shows all games on server"},					//working
		{Text: "quit"},		
*/