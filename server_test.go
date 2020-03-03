package gogam

import (
	"io/ioutil"
	"strings"
	//log "github.com/sirupsen/logrus"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"
)

/*
	{Text: "game start"},
*/

var testCommandNewGame = "game new"
var testCommandGameLoad = "game load"

func postRequestGameHandler(command string, testServer *Server, res *httptest.ResponseRecorder) ([]byte, error) {
	a := bytes.NewBuffer([]byte(command))
	req, _ := http.NewRequest("POST", "/game", a)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if len(res.Result().Cookies()) > 0 {
		req.AddCookie(res.Result().Cookies()[0])
	}
	testServer.gameHandler(res, req)
	return ioutil.ReadAll(res.Body)
}

func postRequestLoginHandler(command string, testServer *Server, res *httptest.ResponseRecorder) ([]byte, error) {
	a := bytes.NewBuffer([]byte(command))
	req, _ := http.NewRequest("POST", "/login", a)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	testServer.loginHandler(res, req)
	return ioutil.ReadAll(res.Body)
}

func TestLoginFailing(t *testing.T) {

	res := httptest.NewRecorder()
	testServer := new(Server)
	testServer.initiateServer()

	content, _ := postRequestLoginHandler("aasd", testServer, res)
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

	content, _ := postRequestLoginHandler("admin", testServer, res)
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
	content, _ := postRequestGameHandler(testCommandGameLoad+" 1", testServer, res)
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

	postRequestLoginHandler("admin", testServer, res)

	content, _ := postRequestGameHandler(testCommandNewGame, testServer, res)
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
	testServer := new(Server)
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
	testServer := new(Server)
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
	testServer := new(Server)
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

func TestUserNewUnauthenticated(t *testing.T) {
	//test user new user1
	res := httptest.NewRecorder()
	testServer := new(Server)
	testServer.initiateServer()

	content, _ := postRequestGameHandler("user new user1", testServer, res)
	expected := "ynogo"
	if string(content) != expected {
		t.Errorf("Expected %s, got %s.", expected, string(content))
	}
	os.Remove("gogam.db")
}

func TestUserNewauthenticated(t *testing.T) {
	//test user new user1
	res := httptest.NewRecorder()
	testServer := new(Server)
	testServer.initiateServer()

	postRequestLoginHandler("admin", testServer, res)

	content, _ := postRequestGameHandler("user new user1", testServer, res)
	expected := "OK"
	if string(content) != expected {
		t.Errorf("Expected %s, got %s.", expected, string(content))
	}
	os.Remove("gogam.db")
}

func TestGameListNoExistingGames(t *testing.T) {
	//log.SetLevel(log.DebugLevel)
	//test game list
	res := httptest.NewRecorder()
	testServer := new(Server)
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
	testServer := new(Server)
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
	testServer := new(Server)
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
	testServer := new(Server)
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
	testServer := new(Server)
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
	testServer := new(Server)
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
	testServer := new(Server)
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
	testServer := new(Server)
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
	testServer := new(Server)
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
	testServer := new(Server)
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
	testServer := new(Server)
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