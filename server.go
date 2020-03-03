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
	"strconv"
	"strings"
)

var apiResponseForbidden = "ynogo"
var apiResponseOk = "OK"
var apiResponseMissingArguments = "not enough arguments"
var apiResponseGameAlreadyRunning = "game already running"
var apiResponseGameLoaded = "game loaded"
var apiResponseGameIDNotFound = "game not found"
var apiResponseUsernameInUser = "Username already in use"

type user struct {
	gorm.Model
	Nick       string
	Online     bool
	Characters []*character `gorm:"many2many:user_character;association_jointable_foreignkey:character_id"`
}

type Server struct {
	db           *gorm.DB
	BindingPort  string
	Users        []*user
	game         *game
	sessionStore *sessions.CookieStore
}

func (selfServer *Server) initDB() {
	var err error
	selfServer.db, err = gorm.Open("sqlite3", "./gogam.db")
	if err != nil {
		panic("failed to connect database")
	}

	selfServer.db.AutoMigrate(&user{})
	selfServer.db.AutoMigrate(&passive{})
	selfServer.db.AutoMigrate(&skill{})
	selfServer.db.AutoMigrate(&character{})
	selfServer.db.AutoMigrate(&game{})
}

func (selfServer *Server) loginHandler(w http.ResponseWriter, r *http.Request) {

	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	bodyString := string(bodyBytes)
	log.Debug("recieved data :", bodyString)

	username := bodyString
	var loginUser user

	selfServer.db.First(&loginUser, "nick = ?", username)

	if reflect.DeepEqual(user{}, loginUser) {
		log.Debug("no user found in db for :", username)
		log.Debug("sending http.StatusUnauthorized ")
		fmt.Fprintf(w, apiResponseForbidden)
		//http.Error(w, "ynogo" , http.StatusUnauthorized)
		return
	}

	log.Debug("user loaded from db :", username)
	session, err := selfServer.sessionStore.Get(r, "gogam-session")
	session.Values["user"] = loginUser
	log.Debug("Saved user :", username, " in Seesions")
	err = session.Save(r, w)
	if err != nil {
		log.Debug("Erros saving user to session:", err)
	}
	fmt.Fprintf(w, apiResponseOk)

}

func userCommandHandler (command []string,selfServer *Server,w http.ResponseWriter) {
	switch command[1] {
		case "new":
			if len(command) != 3 {
				fmt.Fprintf(w, apiResponseMissingArguments)
				return
			}
			username := command[2]
			var loginUser user
			selfServer.db.First(&loginUser, "nick = ?", username)

			if !reflect.DeepEqual(user{}, loginUser) {
				log.Debug("Username already in use")
				fmt.Fprintf(w, apiResponseUsernameInUser)
				return
			}
			loginUser = user{
				Nick: username,
			}
			selfServer.db.Create(&loginUser)
			log.Debug("Saved user :", username, " in Seesions")
			fmt.Fprintf(w, apiResponseOk)
			return
		}
}

func charCommandHandler (command []string,selfServer *Server,w http.ResponseWriter,sessionUser *user) {
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
			selfServer.db.Find(&loadGame, "id = ?", gameID)
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
			//todo game

			loadGame.addCharacter(newChar)
			selfServer.db.Create(&newChar)
			sessionUser.Characters = append(sessionUser.Characters, newChar)

			//db.Model(&user).Related(&card)
			//log.Info(toJSON(sessionUser))
			selfServer.db.Save(&sessionUser)
			fmt.Fprintf(w, apiResponseOk)
			return

		case "list":

			var testUser user
			selfServer.db.Preload("Characters").First(&testUser, "nick = ?", sessionUser.Nick)

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
			//selfServer.db.Preload("Characters").First(&testUser, "nick = ?", sessionUser.Nick)
			// char id
			// name
			//level
			//health/max health
			//exp
			//game
			//passives
			//skills
			var resultChar character
			selfServer.db.Table("Characters").
			Joins("join user_character on user_character.character_id = characters.id").
			Joins("join users on user_character.user_id = users.id").
			//Joins("join user_character on user_character.character_id = characters.id").
			//Joins("join users on user_character.user_id = users.id").
			Where("users.id = ?", "1").Find(&resultChar)
		
			/*
			record := &struct{ ID uint }{}
			selfServer.db.Debug().Table("users").
				Select("character.id").
				Joins("join user_character on user_character.user_id = user.id").
				Joins("join characters on user_character.character_id = characters.id").
				Joins("join game_character on game_character.character_id = characters.id").
				Joins("join games on game_character.character_id = characters.id").
				Where("game.id = ?", selfServer.game.ID).Scan(record)
			*/
			
			//log.Info(toJSON((resultChar)))

			fmt.Fprintf(w, apiResponseOk)
			return

		}
}

func gameCommandHandler (command []string,selfServer *Server,w http.ResponseWriter) {
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
			selfServer.db.Create(&newGame)
			log.Debug("game saved: ", newGame.Name)
			fmt.Fprintf(w, apiResponseOk)
			return

		case "load":
			if len(command) != 3 {
				fmt.Fprintf(w, apiResponseMissingArguments)
				return
			}
			//check, if game already loaded
			if selfServer.game != nil {
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
			selfServer.db.Preload("Characters").Find(&loadGame, "id = ?", gameID)
			if reflect.DeepEqual(loadGame, loadGameCheck) {
				log.Debug("game id not found:", gameID)
				fmt.Fprintf(w, apiResponseGameIDNotFound)
				return
			}

			selfServer.game = &loadGame
			err = selfServer.game.loadGameField()
			if err != nil {
				log.Warn(err)
			}
			log.Debug("game loaded:", loadGame.ID)
			fmt.Fprintf(w, apiResponseGameLoaded)
			return

		case "list":
			var gameArray []game
			selfServer.db.Find(&gameArray)
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

			if selfServer.game.InProgress {
				//running game
				//check, if user has character in the game

				//set online

			} else {
				//new game starting
				var loadCharacter character
				//selfServer.db.Debug().Find(&loadCharacter,"game_id = ?", selfServer.game.ID).Related(&sessionUser)
				//record := &struct {ID uint}{}
				record := &struct{ ID uint }{}
				selfServer.db.Debug().Table("users").
					Select("character.id").
					Joins("join user_character on user_character.user_id = user.id").
					Joins("join characters on user_character.character_id = characters.id").
					Joins("join game_character on game_character.character_id = characters.id").
					Joins("join games on game_character.character_id = characters.id").
					Where("game.id = ?", selfServer.game.ID).Scan(record)

				//selfServer.db.Where(&sessionUser).Find(&loadCharacter)
				//db.Model(&user).Related(&card, "CreditCard")
				log.Info("found:", toJSON(record))

				selfServer.game.addCharacter(&loadCharacter)
				selfServer.db.Save(&selfServer.game)
				fmt.Fprintf(w, apiResponseOk)
				return
			}
		}
}
func (selfServer *Server) gameHandler(w http.ResponseWriter, r *http.Request) {
	session, err := selfServer.sessionStore.Get(r, "gogam-session")
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
	case "user":
		userCommandHandler (bodyString,selfServer,w);return
	case "game":
		gameCommandHandler (bodyString,selfServer,w);return
	case "char":
		charCommandHandler (bodyString,selfServer,w,sessionUser);return
		
	}
	fmt.Fprintf(w, apiResponseOk)
}
func (selfServer *Server) initiateServer() {

	if selfServer.BindingPort == "" {
		selfServer.BindingPort = "7070"
	}

	var b [8]byte
	_, err := crypto_rand.Read(b[:])
	if err != nil {
		panic("cannot seed math/rand package with cryptographically secure random number generator")
	}
	math_rand.Seed(int64(binary.LittleEndian.Uint64(b[:])))

	log.SetReportCaller(true)
	log.Debug("Logging started")
	selfServer.initDB()
	gob.Register(&user{})
	loginUser := user{
		Nick: "admin",
	}
	selfServer.db.Create(&loginUser)
	log.Debug("Admin created")
	selfServer.sessionStore = sessions.NewCookieStore([]byte("sessionkey"))
}

func (selfServer *Server) StartServer() {
	selfServer.initiateServer()
	http.HandleFunc("/game", selfServer.gameHandler)
	http.HandleFunc("/login", selfServer.loginHandler)
	log.Debug("Server is listening on Port: ", selfServer.BindingPort)
	http.ListenAndServe(":"+selfServer.BindingPort, nil)
	/*
			selfServer.Users = append(selfServer.Users,user{
				nick:"fritz",
				//CharacterID: 1,
			})
			selfServer.Users[0].Character = &character {
			PlayerID: 1,
			Name: "supermage",
		    Level: 1,
		    Health: 10,
		    MaxHealth: 10,
			Experience: 2,
			}


			game := &Game{
			Name: "hans",
			GameField: BoardParser(),
			}
			game.addCharacter(selfServer.Users[0].Character)
			game.GameField.ShowMap()
			game.characterOverview()


	*/

}
