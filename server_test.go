package gogam

import (
	"encoding/json"
	"io/ioutil"
	//"strings"
	"bytes"

	log "github.com/sirupsen/logrus"

	//"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"

	//"reflect"
	"testing"
	//"net/http/cookiejar"
)

/*
	{Text: "game start"},
*/

var testCommandNewGame = "game new"
var testCommandGameLoad = "game load"

/*
func postRequestGameHandler(command string, testServer *GameServer, res *httptest.ResponseRecorder) ([]byte, error) {
	a := bytes.NewBuffer([]byte(command))
	req, _ := http.NewRequest("POST", "/game", a)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if len(res.Result().Cookies()) > 0 {
		req.AddCookie(res.Result().Cookies()[0])
	}
	testServer.gameHandler(res, req)
	return ioutil.ReadAll(res.Body)
}
*/

func requestThingy(requestBody []byte, testServer *GameServer, res *httptest.ResponseRecorder,url string, method string) (*httptest.ResponseRecorder,*http.Request)  {
	a := bytes.NewBuffer(requestBody)
	req, _ := http.NewRequest(method, url, a)
	if len(res.Result().Cookies()) > 0 {
		req.AddCookie(res.Result().Cookies()[0])
	}
	return res,req
}

func postRequestLoginHandler(command string, testServer *GameServer, res *httptest.ResponseRecorder) ([]byte, error) {
	a := bytes.NewBuffer([]byte(command))
	req, _ := http.NewRequest("POST", "/login", a)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	testServer.loginHandler(res, req)
	return ioutil.ReadAll(res.Body)
}

func TestUserNewUnauthenticated(t *testing.T) {

	//log.SetLevel(log.DebugLevel)

	type requestJSON struct {
		Username string  `json:"username"`
	}

	res := httptest.NewRecorder()
	testServer := new(GameServer)
	testServer.initiateServer()

	stringJSON := requestJSON{Username: "asd"}
	byteJSON,err := json.Marshal(stringJSON)
		if err != nil {
		t.Errorf(err.Error())
	}

	var req *http.Request
	res,req = requestThingy(byteJSON, testServer, res,"/login", http.MethodPost)
	testServer.loginHandler(res, req)

	content,err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Errorf(err.Error())
	}

	expectedContent := newReturnMessage(401,"not authenticated")
	var recieveMessage returnJSONMessage
	err = json.Unmarshal(content,&recieveMessage)
	if err != nil {
		t.Errorf(err.Error())
	}
	if recieveMessage != *expectedContent {
		t.Errorf("Expected %s, got %s.", expectedContent.toString(), recieveMessage.toString())
	}
	os.Remove("gogam.db")
}

func TestLoginSuccess(t *testing.T) {

	//log.SetLevel(log.DebugLevel)

	type requestJSON struct {
		Username string  `json:"username"`
	}

	res := httptest.NewRecorder()
	testServer := new(GameServer)
	testServer.initiateServer()

	stringJSON := requestJSON{Username: "admin"}
	byteJSON,err := json.Marshal(stringJSON)
		if err != nil {
		t.Errorf(err.Error())
	}

	var req *http.Request
	res,req = requestThingy(byteJSON, testServer, res,"/login", http.MethodPost)
	testServer.loginHandler(res, req)

	content,err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Errorf(err.Error())
	}


	expectedContent := newReturnMessage(200,"OK")
	var recieveMessage returnJSONMessage
	err = json.Unmarshal(content,&recieveMessage)
	if err != nil {
		t.Errorf(err.Error())
	}
	if recieveMessage != *expectedContent {
		t.Errorf("Expected %s, got %s.", expectedContent.toString(), recieveMessage.toString())
	}
	os.Remove("gogam.db")
}

func TestUserNewMalformedJson(t *testing.T) {

	//log.SetLevel(log.DebugLevel)
	//test user new user1
	res := httptest.NewRecorder()
	testServer := new(GameServer)
	testServer.initiateServer()


	var req *http.Request
	res,req = requestThingy([]byte("aaa"), testServer, res,"/user", http.MethodPut)
	testServer.userHandler(res, req)

	content,err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Errorf(err.Error())
	}

	expectedContent := newReturnMessage(400,"malformed request")
	var recieveMessage returnJSONMessage
	err = json.Unmarshal(content,&recieveMessage)
	if err != nil {
		t.Errorf(err.Error())
	}
	if recieveMessage != *expectedContent {
		t.Errorf("Expected %s, got %s.", expectedContent.toString(), recieveMessage.toString())
	}
	os.Remove("gogam.db")
}

func TestUserNewauthenticated(t *testing.T) {
	//log.SetLevel(log.DebugLevel)
	//test user new user1
	res := httptest.NewRecorder()
	testServer := new(GameServer)
	testServer.initiateServer()

	type requestJSON struct {
		Username string  `json:"username"`
	}

	stringJSON := requestJSON{Username: "admin"}
	byteJSON,err := json.Marshal(stringJSON)
		if err != nil {
		t.Errorf(err.Error())
	}

	var req *http.Request
	res,req = requestThingy(byteJSON, testServer, res,"/login", http.MethodPost)
	testServer.loginHandler(res, req)
	ioutil.ReadAll(res.Body)

	stringJSON = requestJSON{Username: "hans"}
	byteJSON,err = json.Marshal(stringJSON)
		if err != nil {
		t.Errorf(err.Error())
	}

	res,req = requestThingy(byteJSON, testServer, res,"/user", http.MethodPut)
	testServer.userHandler(res, req)

	content,err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Errorf(err.Error())
	}

	expectedContent := newReturnMessage(201,"")
	var recieveMessage returnJSONMessage
	err = json.Unmarshal(content,&recieveMessage)
	if err != nil {
		t.Errorf(err.Error())
	}
	if recieveMessage != *expectedContent {
		t.Errorf("Expected %s, got %s.", expectedContent.toString(), recieveMessage.toString())
	}
	os.Remove("gogam.db")
}



func TestNewGameUnautheticated(t *testing.T) {
	//log.SetLevel(log.DebugLevel)
	res := httptest.NewRecorder()
	testServer := new(GameServer)
	testServer.initiateServer()

	type requestJSON struct {
		Username string  `json:"username"`
	}
	type requestGameJSON struct {
		GameName	string  `json:"gameName"`
		Command 	string `json:"command"`
		GameID 		string `json:"gameID"`
	}

	stringJSON := requestJSON{Username: "admin"}
	byteJSON,err := json.Marshal(stringJSON)
		if err != nil {
		t.Errorf(err.Error())
	}
	var req *http.Request
	res,req = requestThingy(byteJSON, testServer, res,"/login", http.MethodPost)
	testServer.loginHandler(res, req)
	ioutil.ReadAll(res.Body)


	gameStringJSON := requestGameJSON{GameName: "GameNAmemitid1"}
	byteJSON,err = json.Marshal(gameStringJSON)
		if err != nil {
		t.Errorf(err.Error())
	}

	res,req = requestThingy(byteJSON, testServer, res,"/game", http.MethodPut)
	testServer.gameHandler(res, req)
	content,err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Errorf(err.Error())
	}

	expectedContent := newReturnMessage(200,"OK")
	var recieveMessage returnJSONMessage
	err = json.Unmarshal(content,&recieveMessage)
	if err != nil {
		t.Errorf(err.Error())
	}
	if recieveMessage != *expectedContent {
		t.Errorf("Expected %s, got %s.", expectedContent.toString(), recieveMessage.toString())
	}
	os.Remove("gogam.db")
}


func TestNewGameautheticatednoGameID(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	res := httptest.NewRecorder()
	testServer := new(GameServer)
	testServer.initiateServer()

	type requestJSON struct {
		Username string  `json:"username"`
	}
	type requestGameJSON struct {
		GameName	string  `json:"gameName"`
		Command 	string `json:"command"`
		GameID 		string `json:"gameID"`
	}

	stringJSON := requestJSON{Username: "admin"}
	byteJSON,err := json.Marshal(stringJSON)
		if err != nil {
		t.Errorf(err.Error())
	}
	var req *http.Request
	res,req = requestThingy(byteJSON, testServer, res,"/login", http.MethodPost)
	testServer.loginHandler(res, req)
	ioutil.ReadAll(res.Body)


	gameStringJSON := requestGameJSON{GameName: "GameNAmemitid1"}
	byteJSON,err = json.Marshal(gameStringJSON)
		if err != nil {
		t.Errorf(err.Error())
	}

	res,req = requestThingy(byteJSON, testServer, res,"/game", http.MethodPut)
	testServer.gameHandler(res, req)
	content,err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Errorf(err.Error())
	}

	expectedContent := newReturnMessage(200,"OK")
	var recieveMessage returnJSONMessage
	err = json.Unmarshal(content,&recieveMessage)
	if err != nil {
		t.Errorf(err.Error())
	}
	if recieveMessage != *expectedContent {
		t.Errorf("Expected %s, got %s.", expectedContent.toString(), recieveMessage.toString())
	}
	os.Remove("gogam.db")
}
/*
func TestNewGameautheticatednewGame(t *testing.T) {
	//test game new world1
	res := httptest.NewRecorder()
	testServer := new(GameServer)
	testServer.initiateServer()

	postRequestLoginHandler("admin", testServer, res)

	content, _ := postRequestGameHandler(testCommandNewGame+" world1", testServer, res)
	expected := "OK"
	if string(content) != expected {
		t.Errorf("Expected %s, got %s.", expected, string(content))
	}
	os.Remove("gogam.db")
}

func TestGameLoadNonExistingGame(t *testing.T) {

	//test game new world1
	res := httptest.NewRecorder()
	testServer := new(GameServer)
	testServer.initiateServer()

	//login/get session
	postRequestLoginHandler("admin", testServer, res)
	//load game without create
	content, _ := postRequestGameHandler(testCommandGameLoad+" 2", testServer, res)
	expected := "game not found"
	if string(content) != expected {
		t.Errorf("Expected %s, got %s.", expected, string(content))
	}
	os.Remove("gogam.db")

}

func TestGameLoadExistingGame(t *testing.T) {
	//test game new world1
	res := httptest.NewRecorder()
	testServer := new(GameServer)
	testServer.initiateServer()

	postRequestLoginHandler("admin", testServer, res)

	postRequestGameHandler(testCommandNewGame+" world1", testServer, res)
	content, _ := postRequestGameHandler(testCommandGameLoad+" 1", testServer, res)
	expected := "game loaded"
	if string(content) != expected {
		t.Errorf("Expected %s, got %s.", expected, string(content))
	}
	os.Remove("gogam.db")
}

func TestGameLoadExistingGameTwice(t *testing.T) {
	//log.SetLevel(log.DebugLevel)
	//test game new world1
	res := httptest.NewRecorder()
	testServer := new(GameServer)
	testServer.initiateServer()

	postRequestLoginHandler("admin", testServer, res)

	postRequestGameHandler(testCommandNewGame+" world1", testServer, res)
	postRequestGameHandler(testCommandGameLoad+" 1", testServer, res)
	content, _ := postRequestGameHandler(testCommandGameLoad+" 1", testServer, res)
	expected := "game loaded"
	if string(content) != expected {
		t.Errorf("Expected %s, got %s.", expected, string(content))
	}
	os.Remove("gogam.db")
	//log.SetLevel(log.InfoLevel)
}



func TestGameListNoExistingGames(t *testing.T) {
	//log.SetLevel(log.DebugLevel)
	//test game list
	res := httptest.NewRecorder()
	testServer := new(GameServer)
	testServer.initiateServer()

	postRequestLoginHandler("admin", testServer, res)

	content, _ := postRequestGameHandler("game list", testServer, res)
	var gameNameArray []string
	json.Unmarshal([]byte(content), &gameNameArray)
	t.Log(gameNameArray)
	expected := []string{}

	if reflect.DeepEqual(gameNameArray, expected) {
		t.Errorf("Expected %s, got %s.", expected, string(content))
	}
	os.Remove("gogam.db")
	//log.SetLevel(log.InfoLevel)
}

func TestGameListTwoGames(t *testing.T) {

	//test game list
	res := httptest.NewRecorder()
	testServer := new(GameServer)
	testServer.initiateServer()

	postRequestLoginHandler("admin", testServer, res)

	postRequestGameHandler(testCommandNewGame+" world1", testServer, res)
	postRequestGameHandler(testCommandNewGame+" world2", testServer, res)

	content, _ := postRequestGameHandler("game list", testServer, res)
	expected := []string{"world1", "world2"}
	//expected := ["world1","world2"]
	contentStrngArray := strings.Split(string(content), " ")

	if reflect.DeepEqual(contentStrngArray, expected) {
		t.Errorf("Expected %s, got %s.", expected, string(content))
	}
	os.Remove("gogam.db")

}

func TestCharNewWithoutGame(t *testing.T) {

	//char new char1 0
	res := httptest.NewRecorder()
	testServer := new(GameServer)
	testServer.initiateServer()

	postRequestLoginHandler("admin", testServer, res)

	content, _ := postRequestGameHandler("char new char1 1", testServer, res)
	expected := "game not found"

	if string(content) != expected {
		t.Errorf("Expected %s, got %s.", expected, string(content))
	}
	os.Remove("gogam.db")

}

func TestCharNewwithExistingGame(t *testing.T) {
	//char new char1 0
	res := httptest.NewRecorder()
	testServer := new(GameServer)
	testServer.initiateServer()

	postRequestLoginHandler("admin", testServer, res)
	postRequestGameHandler(testCommandNewGame+" world1", testServer, res)

	content, _ := postRequestGameHandler("char new char1 1", testServer, res)
	expected := "OK"

	if string(content) != expected {
		t.Errorf("Expected %s, got %s.", expected, string(content))
	}
	os.Remove("gogam.db")

}

func TestCharListWithoutExistingChar(t *testing.T) {
	//char new char1 0
	res := httptest.NewRecorder()
	testServer := new(GameServer)
	testServer.initiateServer()

	postRequestLoginHandler("admin", testServer, res)
	content, _ := postRequestGameHandler("char list", testServer, res)
	expected := []string{}
	contentStrngArray := strings.Split(string(content), " ")

	if reflect.DeepEqual(contentStrngArray, expected) {
		t.Errorf("Expected %s, got %s.", expected, string(content))
	}
	os.Remove("gogam.db")

}

func TestCharListWithExistingChar(t *testing.T) {
	//char new char1 0
	res := httptest.NewRecorder()
	testServer := new(GameServer)
	testServer.initiateServer()

	postRequestLoginHandler("admin", testServer, res)
	postRequestGameHandler(testCommandNewGame+" world1", testServer, res)
	postRequestGameHandler("char new char1 1", testServer, res)
	content, _ := postRequestGameHandler("char list", testServer, res)
	expected := []string{"char1"}
	contentStrngArray := strings.Split(string(content), " ")

	if reflect.DeepEqual(contentStrngArray, expected) {
		t.Errorf("Expected %s, got %s.", expected, string(content))
	}
	os.Remove("gogam.db")

}

func TestCharStats(t *testing.T) {
	//char new char1 0
	res := httptest.NewRecorder()
	testServer := new(GameServer)
	testServer.initiateServer()

	postRequestLoginHandler("admin", testServer, res)
	postRequestGameHandler(testCommandNewGame+" world1", testServer, res)
	postRequestGameHandler("char new char1 1", testServer, res)
	content, _ := postRequestGameHandler("char stats 1", testServer, res)
	expected := []string{"char1"}
	contentStrngArray := strings.Split(string(content), " ")

	if reflect.DeepEqual(contentStrngArray, expected) {
		t.Errorf("Expected %s, got %s.", expected, string(content))
	}
	os.Remove("gogam.db")

}

/*
func TestGameJoinNonloadedGame(t *testing.T) {
	//char new char1 0
	res := httptest.NewRecorder()
	testServer := new(GameServer)
	testServer.initiateServer()

	postRequestLoginHandler("admin", testServer, res)
	postRequestGameHandler(testCommandNewGame+" world1", testServer, res)
	postRequestGameHandler("char new char1 1", testServer, res)
	content, _ := postRequestGameHandler("game join", testServer, res)
	expected := "no game loaded"

	if string(content) != expected {
		t.Errorf("Expected %s, got %s.", expected, string(content))
	}
	os.Remove("gogam.db")

}

func TestGameJoinLoadedGame(t *testing.T) {
	//char new char1 0
	res := httptest.NewRecorder()
	testServer := new(GameServer)
	testServer.initiateServer()

	postRequestLoginHandler("admin", testServer, res)
	postRequestGameHandler(testCommandNewGame+" world1", testServer, res)
	postRequestGameHandler("char new char1 1", testServer, res)
	postRequestGameHandler(testCommandGameLoad+" 1", testServer, res)
	content, _ := postRequestGameHandler("game join", testServer, res)
	expected := "OK"

	if string(content) != expected {
		t.Errorf("Expected %s, got %s.", expected, string(content))
	}
	os.Remove("gogam.db")

}

func TestGameJoinLoadedGameWithoutChar(t *testing.T) {
	//char new char1 0
	res := httptest.NewRecorder()
	testServer := new(GameServer)
	testServer.initiateServer()

	postRequestLoginHandler("admin", testServer, res)
	postRequestGameHandler(testCommandNewGame+" world1", testServer, res)
	postRequestGameHandler(testCommandGameLoad+" 1", testServer, res)
	content, _ := postRequestGameHandler("game join", testServer, res)
	expected := "no character for this game"

	if string(content) != expected {
		t.Errorf("Expected %s, got %s.", expected, string(content))
	}
	os.Remove("gogam.db")

}

func TestGameJoinLoadedGameWithWrongChar(t *testing.T) {
	//char new char1 0
	res := httptest.NewRecorder()
	testServer := new(GameServer)
	testServer.initiateServer()

	postRequestLoginHandler("admin", testServer, res)
	postRequestGameHandler(testCommandNewGame+" world1", testServer, res)
	postRequestGameHandler(testCommandNewGame+" world2", testServer, res)
	postRequestGameHandler("char new char1 1", testServer, res)
	postRequestGameHandler(testCommandGameLoad+" 2", testServer, res)
	content, _ := postRequestGameHandler("game join", testServer, res)
	expected := "no character for this world"

	if string(content) != expected {
		t.Errorf("Expected %s, got %s.", expected, string(content))
	}
	os.Remove("gogam.db")

}

*/