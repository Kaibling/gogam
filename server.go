package gogam

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"io/ioutil"
	"net/http"
	//"bytes"
	crypto_rand "crypto/rand"
	"encoding/binary"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"github.com/gorilla/sessions"
	log "github.com/sirupsen/logrus"
	math_rand "math/rand"
	"reflect"
	//"strconv"
	//"strings"
)

var apiResponseForbidden = "not allowed"
var apiResponseUnauthenticated = "not authenticated"
var apiResponseOk = "OK"
var apiResponseMissingArguments = "not enough arguments"
var apiResponseGameAlreadyRunning = "game already running"
var apiResponseGameLoaded = "game loaded"
var apiResponseGameIDNotFound = "game not found"
var apiResponseUsernameInUser = "Username already in use"
var apiResponseMalformedJSON = "malformed request"

type returnJSONMessage struct {
	StatusCode int `json:"status_code"`
	Message string `json:"message"`

}

func newReturnMessage(statusCode int,message string) *returnJSONMessage {
	returnMessage := new(returnJSONMessage)
	returnMessage.StatusCode = statusCode
	returnMessage.Message = message
	return returnMessage
}

func (selfReturnJSONMessage *returnJSONMessage) toString() string {
	returnString,err := json.Marshal(&selfReturnJSONMessage)
    if err != nil {
        panic(err)
	}
	return string(returnString)
}



type user struct {
	gorm.Model
	Nick       string
	Online     bool
	Characters []*character `gorm:"many2many:user_character;association_jointable_foreignkey:character_id"`
}

type GameServer struct {
	db           *gorm.DB
	BindingPort  string
	Users        []*user
	game         *game
	sessionStore *sessions.CookieStore
}

func (selfGameServer *GameServer) initDB() {
	var err error
	selfGameServer.db, err = gorm.Open("sqlite3", "./gogam.db")
	if err != nil {
		panic("failed to connect database")
	}

	selfGameServer.db.AutoMigrate(&user{})
	selfGameServer.db.AutoMigrate(&passive{})
	selfGameServer.db.AutoMigrate(&skill{})
	selfGameServer.db.AutoMigrate(&character{})
	selfGameServer.db.AutoMigrate(&game{})
}

func (selfGameServer *GameServer) loginHandler(w http.ResponseWriter, r *http.Request) {

	var data map[string]interface{}
	body, readErr := ioutil.ReadAll(r.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}
	log.Debug("request:",string(body))
    err := json.Unmarshal(body, &data)
    if err != nil {
		log.Warn("Not a valid JSON struct",string(body))
		http.Error(w, newReturnMessage(http.StatusBadRequest,apiResponseMalformedJSON).toString() , http.StatusBadRequest)
		return
	}

	username := data["username"]
	var loginUser user

	selfGameServer.db.First(&loginUser, "nick = ?", username)

	if reflect.DeepEqual(user{}, loginUser) {
		log.Debug("no user found in db for :", username)
		log.Debug("sending http.StatusUnauthorized ")
		http.Error(w, newReturnMessage(http.StatusUnauthorized,apiResponseUnauthenticated).toString() , http.StatusUnauthorized)
		return
	}

	log.Debug("user loaded from db :", username)
	session, err := selfGameServer.sessionStore.Get(r, "gogam-session")
	session.Values["user"] = loginUser
	log.Debug("Saved user :", username, " in Seesions")
	err = session.Save(r, w)
	if err != nil {
		log.Debug("Erros saving user to session:", err)
	}
	http.Error(w, newReturnMessage(http.StatusOK,apiResponseOk).toString() , http.StatusOK)


}

func charCommandHandler (command []string,selfGameServer *GameServer,w http.ResponseWriter,sessionUser *user) {
	switch command[1] {
		case "new":
			if len(command) != 4 {
				fmt.Fprintf(w, apiResponseMissingArguments)
				return
			}

			gameID := command[3]
			characterName := command[2]
			maxHealth := 90
			minHealth := 35
			health := math_rand.Intn(maxHealth - minHealth)
			var loadGame game
			var loadGameCheck game
			loadGameCheck = loadGame
			selfGameServer.db.Find(&loadGame, "id = ?", gameID)
			if reflect.DeepEqual(loadGame, loadGameCheck) {
				log.Debug("game id not found:", gameID)
				fmt.Fprintf(w, apiResponseGameIDNotFound)
				return
			}
			loadGame.loadGameField()
			newChar := &character{
				Name:       characterName,
				Game:       &loadGame,
				Level:      1,
				Health:     health,
				MaxHealth:  health,
				Experience: 0,
			}
			loadGame.addCharacter(newChar)
			selfGameServer.db.Create(&newChar)
			sessionUser.Characters = append(sessionUser.Characters, newChar)

			//db.Model(&user).Related(&card)
			//log.Info(toJSON(sessionUser))
			selfGameServer.db.Save(&sessionUser)
			fmt.Fprintf(w, apiResponseOk)
			return

		case "list":

			var testUser user
			selfGameServer.db.Preload("Characters").First(&testUser, "nick = ?", sessionUser.Nick)

			var returnArray []string
			for _, ga := range testUser.Characters {
				returnArray = append(returnArray, ga.Name)
			}

			byteArray, err := json.Marshal(&returnArray)
			if err != nil {
				log.Debug(err)
			}
			fmt.Fprintf(w, string(byteArray))
			return
		case "stats":

			//var testUser user
			//selfGameServer.db.Preload("Characters").First(&testUser, "nick = ?", sessionUser.Nick)
			// char id
			// name
			//level
			//health/max health
			//exp
			//game
			//passives
			//skills
			var resultChar character
			selfGameServer.db.Table("Characters").
			Joins("join user_character on user_character.character_id = characters.id").
			Joins("join users on user_character.user_id = users.id").
			//Joins("join user_character on user_character.character_id = characters.id").
			//Joins("join users on user_character.user_id = users.id").
			Where("users.id = ?", "1").Find(&resultChar)
		
			/*
			record := &struct{ ID uint }{}
			selfGameServer.db.Debug().Table("users").
				Select("character.id").
				Joins("join user_character on user_character.user_id = user.id").
				Joins("join characters on user_character.character_id = characters.id").
				Joins("join game_character on game_character.character_id = characters.id").
				Joins("join games on game_character.character_id = characters.id").
				Where("game.id = ?", selfGameServer.game.ID).Scan(record)
			*/
			
			//log.Info(toJSON((resultChar)))

			fmt.Fprintf(w, apiResponseOk)
			return

		}
}



func (selfGameServer *GameServer) userHandler (w http.ResponseWriter, r *http.Request) {

	var data map[string]interface{}
	
	body, readErr := ioutil.ReadAll(r.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}
	log.Debug("request:",string(body))
    err := json.Unmarshal(body, &data)
    if err != nil {
		log.Warn("Not a valid JSON struct",string(body))
		http.Error(w, newReturnMessage(http.StatusBadRequest,apiResponseMalformedJSON).toString() , http.StatusBadRequest)
		return
	}
	var message string
	returnStatusCode := http.StatusInternalServerError

		switch r.Method {
		case http.MethodPut:
			if data["username"] == nil {
				message = "username is missing"
				returnStatusCode = http.StatusConflict
				break
			}
			username := data["username"].(string)
			var loginUser user
			selfGameServer.db.First(&loginUser, "nick = ?", username)

			if !reflect.DeepEqual(user{}, loginUser) {
				message = "Username already in use"
				log.Debug(message)
				returnStatusCode = http.StatusConflict
				break
			}
			loginUser = user{
				Nick: username,
			}
			selfGameServer.db.Create(&loginUser)
			message = ""
			log.Debug("User ",username," Sucessful created")
			returnStatusCode = http.StatusCreated
			break
		default:
			log.Debug(r.Method," is not supported ")
		}
		w.WriteHeader(returnStatusCode)
		returnMessage := &returnJSONMessage{ Message: message,StatusCode:returnStatusCode}
		fmt.Fprintf(w,returnMessage.toString())
		
}

func (selfGameServer *GameServer) charHandler (w http.ResponseWriter, r *http.Request) {

	var data map[string]interface{}
	
	body, readErr := ioutil.ReadAll(r.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}
	log.Debug("request:",string(body))
    err := json.Unmarshal(body, &data)
    if err != nil {
        panic(err)
	}
	var message string
	returnStatusCode := http.StatusInternalServerError

		switch r.Method {
		case http.MethodPut:
			if data["username"] == nil {
				message = "username is missing"
				returnStatusCode = http.StatusConflict
				break
			}
			if data["game"] == nil {
				message = "game is missing"
				returnStatusCode = http.StatusConflict
				break
			}
			username := data["username"].(string)
			//gameID := data["game"].(string)
			var loginUser user
			selfGameServer.db.First(&loginUser, "nick = ?", username)

			if !reflect.DeepEqual(user{}, loginUser) {
				message = "Username already in use"
				log.Debug(message)
				returnStatusCode = http.StatusConflict
				break
			}
			loginUser = user{
				Nick: username,
			}
			selfGameServer.db.Create(&loginUser)
			message = ""
			log.Debug("User ",username," Sucessful created")
			returnStatusCode = http.StatusCreated
			break
		default:
			log.Debug(r.Method," is not supported ")
		}
		w.WriteHeader(returnStatusCode)
		returnMessage := &returnJSONMessage{ Message: message,StatusCode:returnStatusCode}
		fmt.Fprintf(w,returnMessage.toString())
		
}

/*
func gameCommandHandler (command []string,selfGameServer *GameServer,w http.ResponseWriter) {
	switch command[1] {
		case "new":
			if len(command) != 3 {
				fmt.Fprintf(w, apiResponseMissingArguments)
				return
			}
			newGame := &game{
				Name:       command[2],
				InProgress: false,
				GameField:  BoardParser(),
			}
			newGame.saveGameField()
			selfGameServer.db.Create(&newGame)
			log.Debug("game saved: ", newGame.Name)
			fmt.Fprintf(w, apiResponseOk)
			return

		case "load":
			if len(command) != 3 {
				fmt.Fprintf(w, apiResponseMissingArguments)
				return
			}
			//check, if game already loaded
			if selfGameServer.game != nil {
				//game loaded
				log.Debug("game already loaded")
				fmt.Fprintf(w, apiResponseGameLoaded)
				return
			}
			var err error
			gameID, err := strconv.ParseUint(command[2], 10, 32)
			if err != nil {
				log.Debug("error thingy: ", err)
			}
			log.Debug("loading game :", gameID)
			var loadGame game
			var loadGameCheck game
			loadGameCheck = loadGame
			selfGameServer.db.Preload("Characters").Find(&loadGame, "id = ?", gameID)
			if reflect.DeepEqual(loadGame, loadGameCheck) {
				log.Debug("game id not found:", gameID)
				fmt.Fprintf(w, apiResponseGameIDNotFound)
				return
			}

			selfGameServer.game = &loadGame
			err = selfGameServer.game.loadGameField()
			if err != nil {
				log.Warn(err)
			}
			log.Debug("game loaded:", loadGame.ID)
			fmt.Fprintf(w, apiResponseGameLoaded)
			return

		case "list":
			var gameArray []game
			selfGameServer.db.Find(&gameArray)
			var returnArray []string
			for _, ga := range gameArray {
				returnArray = append(returnArray, ga.Name)
			}
			byteArray, err := json.Marshal(&returnArray)
			if err != nil {
				log.Debug(err)
			}
			fmt.Fprintf(w, string(byteArray))
			return
		case "join":
			if len(command) != 2 {
				fmt.Fprintf(w, apiResponseMissingArguments)
				return
			}

			if selfGameServer.game.InProgress {
				//running game
				//check, if user has character in the game

				//set online

			} else {
				//new game starting
				var loadCharacter character
				//selfGameServer.db.Debug().Find(&loadCharacter,"game_id = ?", selfGameServer.game.ID).Related(&sessionUser)
				//record := &struct {ID uint}{}
				record := &struct{ ID uint }{}
				selfGameServer.db.Debug().Table("users").
					Select("character.id").
					Joins("join user_character on user_character.user_id = user.id").
					Joins("join characters on user_character.character_id = characters.id").
					Joins("join game_character on game_character.character_id = characters.id").
					Joins("join games on game_character.character_id = characters.id").
					Where("game.id = ?", selfGameServer.game.ID).Scan(record)

				//selfGameServer.db.Where(&sessionUser).Find(&loadCharacter)
				//db.Model(&user).Related(&card, "CreditCard")
				log.Info("found:", toJSON(record))

				selfGameServer.game.addCharacter(&loadCharacter)
				selfGameServer.db.Save(&selfGameServer.game)
				fmt.Fprintf(w, apiResponseOk)
				return
			}
		}
}

func (selfGameServer *GameServer) gameHandler(w http.ResponseWriter, r *http.Request) {
	session, err := selfGameServer.sessionStore.Get(r, "gogam-session")
	if err != nil {
		log.Debug("error creating session")
	}

	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	bodyString := strings.Split(string(bodyBytes), " ")
	log.Debug("recieved data :", bodyString)

	val := session.Values["user"]
	var sessionUser = &user{}
	var ok bool
	//authenticated := false
	if sessionUser, ok = val.(*user); !ok {
		log.Debug("not authenticated:", session.Values["user"])
		fmt.Fprintf(w, apiResponseForbidden)
		return
	}
	log.Debug("User authenticated:", sessionUser.Nick)
	//authenticated = true

	switch bodyString[0] {
	case "game":
		gameCommandHandler (bodyString,selfGameServer,w);return
	case "char":
		charCommandHandler (bodyString,selfGameServer,w,sessionUser);return
		
	}
	fmt.Fprintf(w, apiResponseOk)
}
*/
func (selfGameServer *GameServer) initiateServer() {

	if selfGameServer.BindingPort == "" {
		selfGameServer.BindingPort = "7070"
	}

	var b [8]byte
	_, err := crypto_rand.Read(b[:])
	if err != nil {
		panic("cannot seed math/rand package with cryptographically secure random number generator")
	}
	math_rand.Seed(int64(binary.LittleEndian.Uint64(b[:])))

	log.SetReportCaller(true)
	log.Debug("Logging started")
	selfGameServer.initDB()
	gob.Register(&user{})
	loginUser := user{
		Nick: "admin",
	}
	selfGameServer.db.Create(&loginUser)
	log.Debug("Admin created")
	selfGameServer.sessionStore = sessions.NewCookieStore([]byte("sessionkey"))
}

func (selfGameServer *GameServer) StartServer() {
	selfGameServer.initiateServer()
	//http.HandleFunc("/game", selfGameServer.gameHandler)
	http.HandleFunc("/login", selfGameServer.loginHandler)
	http.HandleFunc("/user", selfGameServer.userHandler)
	http.HandleFunc("/char", selfGameServer.charHandler)
	log.Debug("GameServer is listening on Port: ", selfGameServer.BindingPort)
	http.ListenAndServe(":"+selfGameServer.BindingPort, nil)

}
