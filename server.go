package gogam
import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"net/http"
	"io/ioutil"
	//"bytes"
	//"encoding/json"
	"log"
	"fmt"
	"github.com/gorilla/sessions"
	"strings"
	"encoding/gob"
	
)

type user struct {
	gorm.Model
	Nick string
	Character *character
}

type Server struct {
	db *gorm.DB
	BindingPort string
	Users 		[]user
	//playerRepo *playerRepo
	game *game
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
	//defer 
}

func (selfServer *Server) loginHandler(w http.ResponseWriter, r *http.Request) {

	bodyBytes, err := ioutil.ReadAll(r.Body)
    if err != nil {
        log.Fatal(err)
	}

    bodyString := string(bodyBytes)
	log.Println("recieved data :",bodyString)
	

	username := bodyString
	var loginUser user
	selfServer.db.First(&loginUser, "nick = ?", username)
	if loginUser == (user{}) {
		log.Println("no user found in db for :",username)
		log.Println("sending http.StatusUnauthorized ")
		http.Error(w, "ynogo" , http.StatusUnauthorized)
		return
	}

	log.Println("user loaded from db :",username)
	session, err := selfServer.sessionStore.Get(r, "gogam-session")
	session.Values[username] = loginUser
	log.Println("Saved user :",username, " in Seesions")
	err = session.Save(r, w)
	if err != nil {
		log.Println("Erros saving user to session:",err)
	}
	fmt.Fprintf(w, "OK")


}

func (selfServer *Server) gameHandler(w http.ResponseWriter, r *http.Request) {
	session, err := selfServer.sessionStore.Get(r, "gogam-session")
	if err != nil  {
		log.Println("error creating session")
	}

	bodyBytes, err := ioutil.ReadAll(r.Body)
    if err != nil {
        log.Fatal(err)
	}

    bodyString := strings.Split(string(bodyBytes)," ")
	log.Println("recieved data :",bodyString)

	val := session.Values[bodyString[2]]
	var getUser = &user{}
	var ok bool
	authenticated := false
    if getUser, ok = val.(*user); !ok {
		// Handle the case that it's not an expected type
		log.Println("not authenticated")
	} else {
		log.Println("User authenticated:",getUser.Nick)
		authenticated = true
	}
	

	switch bodyString[0] {
		case "user":
			switch bodyString[1] {
				case "new":
					username := bodyString[2]
					var loginUser user
					selfServer.db.First(&loginUser, "nick = ?", username)

					if loginUser != (user{}) {
						log.Println("Username already in use")
						fmt.Fprintf(w, "Username already in use")
						return
					} 
					loginUser = user {
						Nick: username,
					}
					selfServer.db.Create(&loginUser)
					//session.Values[username] = loginUser
					log.Println("Saved user :",username, " in Seesions")
					//err = session.Save(r, w)
					//if err != nil {
					//	log.Println("Erros saving user to session:",err)
					//}
					fmt.Fprintf(w, "OK")
					return
			}
		case "load":
			if authenticated {
				log.Println("sd")
			}
		case "game":
			if bodyString[1] == "new" {
				newGame := &game{
				Name: bodyString[2],
				}
				newGame.GameField = BoardParser()
				selfServer.db.Create(&newGame)

			}else if bodyString[1] == "load" {
				//check, if existing
				//check, if game already loaded
			}
	}
	fmt.Fprintf(w, "OK")
}


func (selfServer *Server) StartServer() {
	if selfServer.BindingPort == "" {
		selfServer.BindingPort = "7070"
	}


    //log.SetPrefix("TRACE: ")
    log.SetFlags(log.Ldate | log.Lmicroseconds | log.Llongfile)
    log.Println("Logging started")

	selfServer.initDB()

	gob.Register(&user{})
	selfServer.sessionStore = sessions.NewCookieStore([]byte("sessionkey"))
	http.HandleFunc("/game", selfServer.gameHandler)
	http.HandleFunc("/login", selfServer.loginHandler)

	log.Println("Server is listening on Port: ",selfServer.BindingPort)
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

