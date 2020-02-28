package gogam
import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"net/http"
	"io/ioutil"
	//"bytes"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"fmt"
	"github.com/gorilla/sessions"
	"strings"
	"encoding/gob"
	crypto_rand "crypto/rand"
    "encoding/binary"
	math_rand "math/rand"
	"reflect"
	
)

type user struct {
	gorm.Model
	Nick 		string
	Online		bool
	Characters 	[]*character `gorm:"many2many:user_character;association_jointable_foreignkey:character_id"`
}

type Server struct {
	db 				*gorm.DB
	BindingPort 	string
	Users 			[]*user
	game 			*game
	sessionStore 	*sessions.CookieStore
	
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


	//selfServerdb.Model(&character).Related(&skill)

}

func (selfServer *Server) loginHandler(w http.ResponseWriter, r *http.Request) {

	bodyBytes, err := ioutil.ReadAll(r.Body)
    if err != nil {
        log.Fatal(err)
	}

    bodyString := string(bodyBytes)
	log.Debug("recieved data :",bodyString)
	

	username := bodyString
	var loginUser user

	selfServer.db.First(&loginUser, "nick = ?", username)
	
	if reflect.DeepEqual(user{}, loginUser) {
		log.Debug("no user found in db for :",username)
		log.Debug("sending http.StatusUnauthorized ")
		fmt.Fprintf(w, "ynogo")
		//http.Error(w, "ynogo" , http.StatusUnauthorized)
		return
	}

	log.Debug("user loaded from db :",username)
	session, err := selfServer.sessionStore.Get(r, "gogam-session")
	session.Values["user"] = loginUser
	log.Debug("Saved user :",username, " in Seesions")
	err = session.Save(r, w)
	if err != nil {
		log.Debug("Erros saving user to session:",err)
	}
	fmt.Fprintf(w, "OK")


}

func (selfServer *Server) gameHandler(w http.ResponseWriter, r *http.Request) {
	session, err := selfServer.sessionStore.Get(r, "gogam-session")
	if err != nil  {
		log.Debug("error creating session")
	}

	bodyBytes, err := ioutil.ReadAll(r.Body)
    if err != nil {
        log.Fatal(err)
	}

    bodyString := strings.Split(string(bodyBytes)," ")
	log.Debug("recieved data :",bodyString)
	
	val := session.Values["user"]
	var sessionUser = &user{}
	var ok bool
	authenticated := false
    if sessionUser, ok = val.(*user); !ok {
		// Handle the case that it's not an expected type
		log.Debug("not authenticated:",session.Values["user"])
		fmt.Fprintf(w, "ynogo")
		//http.Error(w, "ynogo" , http.StatusUnauthorized)
		return
	} 
	log.Debug("User authenticated:",sessionUser.Nick)
	authenticated = true
	
	switch bodyString[0] {
		case "user":
			switch bodyString[1] {
				case "new":
					if len(bodyString) != 3 {
						fmt.Fprintf(w, "not enough arguments")
						return
					}
					username := bodyString[2]
					var loginUser user
					selfServer.db.First(&loginUser, "nick = ?", username)
						
					if !reflect.DeepEqual(user{}, loginUser) {
						log.Debug("Username already in use")
						fmt.Fprintf(w, "Username already in use")
						return
					} 
					loginUser = user {
						Nick: username,
					}
					selfServer.db.Create(&loginUser)
					log.Debug("Saved user :",username, " in Seesions")
					fmt.Fprintf(w, "OK")
					return
			}
		case "load":
			if authenticated {
				log.Debug("sd")
			}
		case "game":
			if bodyString[1] == "new" {
				if len(bodyString) != 3 {
					fmt.Fprintf(w, "not enough arguments")
					return
					}
				newGame := &game{
				Name: bodyString[2],
				InProgress: false,
				}
				newGame.GameField = BoardParser()
				selfServer.db.Create(&newGame)

			}else if bodyString[1] == "load" {
				if len(bodyString) != 3 {
					fmt.Fprintf(w, "not enough arguments")
					return
				}
				//check, if game already loaded
				if reflect.DeepEqual(game{}, selfServer.game) {
					log.Debug("game is already running")
					fmt.Fprintf(w, "game is already running")
					return
				}
				var loadGame game
				selfServer.db.Preload("Characters").Preload("GameField").First(&loadGame, "id = ?", bodyString[2])
				selfServer.game = &loadGame
				log.Debug("game loaded:",loadGame.ID)
				fmt.Fprintf(w, "game loaded")
				return

			}else if bodyString[1] == "list" {
				var gameArray []game
				selfServer.db.Find(&gameArray)
				var returnArray [] string
				for _,ga := range gameArray{
					returnArray = append(returnArray,ga.Name)
				}
				byteArray,err := json.Marshal(&returnArray)
				if err != nil {
					log.Debug(err)
				}
				fmt.Fprintf(w, string(byteArray))
				return
			}else if bodyString[1] == "join" {
				if len(bodyString) != 3 {
					fmt.Fprintf(w, "not enough arguments")
					return
				}

				if selfServer.game.InProgress{
					//running game
					//check, if user has character in the game
			
					//set online

				} else {
					//new game starting
					var loadCharacter character
					//selfServer.db.First(&loadCharacter,"name=?",characterName)
					selfServer.game.addCharacter(&loadCharacter)
					selfServer.db.Save(&selfServer.game)
					fmt.Fprintf(w, "OK")
					return
				}

			}else {
				log.Debug("Unknown Command :",bodyString)
			}
		case "char":
			if bodyString[1] == "new" {
				if len(bodyString) != 4 {
					fmt.Fprintf(w, "not enough arguments")
					return
				}

				gameID := bodyString[3]
				characterName := bodyString[2]
				maxHealth := 90
				minHealth := 35
				health:= math_rand.Intn(maxHealth-minHealth)
				var loadGame game
				selfServer.db.First(&loadGame, "id = ?", gameID)


				newChar := &character{
					Name: characterName,
					Game: &loadGame,
					Level: 1,
					Health: health,
					MaxHealth: health,
					Experience: 0,
				}
				//log.Debug(toJSON(newChar))

				sessionUser.Characters = append(sessionUser.Characters,newChar)
				selfServer.db.Save(&sessionUser)
				
				//selfServer.db.Update(&sessionUser)

			}else if bodyString[1] == "list"{

				var testUser user
				selfServer.db.Preload("Characters").First(&testUser, "nick = ?", sessionUser.Nick)

				var returnArray [] string
				for _,ga := range testUser.Characters{
					returnArray = append(returnArray,ga.Name)
				}

				byteArray,err := json.Marshal(&returnArray)
				if err != nil {
					log.Debug(err)
				}
				fmt.Fprintf(w, string(byteArray))
				return

				/*
				for _,b := range testUser.Characters {
					var testCharacter character
					selfServer.db.Preload("Passives").Preload("Skills").First(&testCharacter, "name = ?", b.Name)
					byteArray,err := json.Marshal(&testCharacter)
					if err != nil {
						log.Debug(err)
					}
					log.Debug(string(byteArray))
				*/

				//fmt.Fprintf(w, string(byteArray))

				//selfServer.db.Find(&userChar, user{nick: 20})
				//selfServer.db.Preload("Passives").Preload("Skills").First(&userChar, "name = ?", "player2")
				//selfServer.db.Find(&userChar)
				

			}else {
				log.Debug("da is mal nix")

			}

	}
	fmt.Fprintf(w, "OK")
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

	
	loginUser := user {
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
	log.Debug("Server is listening on Port: ",selfServer.BindingPort)
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

