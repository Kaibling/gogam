package gogam

import (
	"io/ioutil"
	//log "github.com/sirupsen/logrus"
	"net/http"
	"net/http/httptest"
	"testing"
	"os"

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


func postRequestGameHandler(command string,url string,testServer *Server,res *httptest.ResponseRecorder) ([]byte, error) {
	a := bytes.NewBuffer([]byte(command))
	req, _ := http.NewRequest("POST", url, a)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if len(res.Result().Cookies()) > 0 {
		req.AddCookie(res.Result().Cookies()[0])
	}
	testServer.gameHandler(res, req)
 	return ioutil.ReadAll(res.Body)
}

func postRequestLoginHandler(command string,url string,testServer *Server,res *httptest.ResponseRecorder) ([]byte, error) {
	a := bytes.NewBuffer([]byte(command))
	req, _ := http.NewRequest("POST", url, a)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	testServer.loginHandler(res, req)
 	return ioutil.ReadAll(res.Body)
}


func TestLoginFailing(t *testing.T) {

	res := httptest.NewRecorder()
	testServer := new(Server)
	testServer.initiateServer()

	content, _ := postRequestLoginHandler("aasd","/login",testServer,res)
	expected := "ynogo"
	if string(content) != expected {
		t.Errorf("Expected %s, got %s.", expected, string(content))
	}
	os.Remove("gogam.db")
}

func TestLoginSuccess(t *testing.T) {

	res := httptest.NewRecorder()
	testServer := new(Server)
	testServer.initiateServer()

	content, _ := postRequestLoginHandler("admin","/login",testServer,res)
	expected := "OK"
	if string(content) != expected {
		t.Errorf("Expected %s, got %s.", expected, string(content))
	}
	os.Remove("gogam.db")
}

func TestNewGameUnautheticated(t *testing.T) {
	//log.SetLevel(log.DebugLevel)
	res := httptest.NewRecorder()
	testServer := new(Server)
	testServer.initiateServer()

	//test game new (unauthenticated)
	content, _ := postRequestGameHandler("game load 1","/game",testServer,res)	
	expected := "ynogo"
	if string(content) != expected {
		t.Errorf("Expected %s, got %s.", expected, string(content))
	}
	os.Remove("gogam.db")
}
func TestNewGameautheticatednoGameID(t *testing.T) {
	//test game new
	res := httptest.NewRecorder()
	testServer := new(Server)
	testServer.initiateServer()

	postRequestLoginHandler("admin","/login",testServer,res)

	content, _ := postRequestGameHandler("game new","/game",testServer,res)
	expected := "not enough arguments"
	if string(content) != expected {
		t.Errorf("Expected %s, got %s.", expected, string(content))
	}
	os.Remove("gogam.db")
}

func TestNewGameautheticatednewGame(t *testing.T) {
	//test game new world1
	res := httptest.NewRecorder()
	testServer := new(Server)
	testServer.initiateServer()

	postRequestLoginHandler("admin","/login",testServer,res)

	content, _ := postRequestGameHandler("game new world1","/game",testServer,res)
	expected := "OK"
	if string(content) != expected {
		t.Errorf("Expected %s, got %s.", expected, string(content))
	}
	os.Remove("gogam.db")
}


func TestGameLoadNonExistingGame(t *testing.T) {
	//test game new world1
	res := httptest.NewRecorder()
	testServer := new(Server)
	testServer.initiateServer()

	//login/get session
	postRequestLoginHandler("admin","/login",testServer,res)
	//load game without create
	//postRequestGameHandler("game load 1","/game",testServer,res)
	content, _ := postRequestGameHandler("game load 1","/game",testServer,res)
	expected := "OK"
	if string(content) != expected {
		t.Errorf("Expected %s, got %s.", expected, string(content))
	}
	os.Remove("gogam.db")
}

func TestGameLoadExistingGame(t *testing.T) {
	//test game new world1
	res := httptest.NewRecorder()
	testServer := new(Server)
	testServer.initiateServer()

	postRequestLoginHandler("admin","/login",testServer,res)

	postRequestGameHandler("game new world1","/game",testServer,res)
	postRequestGameHandler("game load 1","/game",testServer,res)
	content, _ := postRequestGameHandler("game load 0","/game",testServer,res)
	expected := "game already loaded"
	if string(content) != expected {
		t.Errorf("Expected %s, got %s.", expected, string(content))
	}
	os.Remove("gogam.db")
}

func TestGameLoadExistingGameTwice(t *testing.T) {
	
	//test game new world1
	res := httptest.NewRecorder()
	testServer := new(Server)
	testServer.initiateServer()

	postRequestLoginHandler("admin","/login",testServer,res)

	postRequestGameHandler("game new world1","/game",testServer,res)
	postRequestGameHandler("game load 0","/game",testServer,res)
	postRequestGameHandler("game load 0","/game",testServer,res)

	content, _ := postRequestGameHandler("game load 0","/game",testServer,res)
	expected := "game already loaded"
	if string(content) != expected {
		t.Errorf("Expected %s, got %s.", expected, string(content))
	}
	os.Remove("gogam.db")
}
