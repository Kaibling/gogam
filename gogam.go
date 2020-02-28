package gogam

import (
    "fmt"
    log "github.com/sirupsen/logrus"
    "math"
  	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"errors"
	"encoding/json"
)


type tile struct {
	Movable bool
	Interactable bool
	Info	string
	AsciArt rune
}

func (selfTile *tile) initTile() {

	switch selfTile.AsciArt {
		case '=':
			//Wall
			selfTile.Info ="There is a Wall"
		case 'S':
			//Startpoint
			selfTile.Info ="Only the Floor"
			selfTile.Movable = true
		case ' ':
			//Floor
			selfTile.Info ="Only the Floor"
			selfTile.Movable = true
		
		case 'D':
			//Door
			selfTile.Info ="It's ... a Door"
			selfTile.Movable = true
			selfTile.Interactable = true
	}
}

type gameField struct {
	Field *[][]tile
	startPoints []position
}

func (selfGameField *gameField) ShowMap() {
	for _,row := range *selfGameField.Field {
		for _,tilchen := range row {
			fmt.Print(string(tilchen.AsciArt))
		}
		fmt.Println()
	}
}
func (selfGameField *gameField) getStartPosition() (position,error) {
	if len(selfGameField.startPoints) == 0 {
		return position{},errors.New("no Starting Points available")

	}
	returnPosition := selfGameField.startPoints[len(selfGameField.startPoints)-1]
	selfGameField.startPoints = selfGameField.startPoints[:len(selfGameField.startPoints)-1]
	return returnPosition, nil

}


type position struct {
	X int
	Y int
}

type game struct {
	gorm.Model			
	Characters 	[]*character	`gorm:"many2many:game_character;association_jointable_foreignkey:character_id"`
	Name		string
	GameField 	*gameField
	InProgress	bool
}

func (selfGame *game) characterOverview() {
	for cnt, char := range selfGame.Characters {
		JSONChar,_ := json.Marshal(&char)
		fmt.Println(cnt, ": ",string(JSONChar))
	}

}

func (selfGame *game) addCharacter(character *character) {
	var err error
	selfGame.Characters = append(selfGame.Characters,character)
	character.Position,err = selfGame.GameField.getStartPosition()
	if err != nil {
		log.Println(err)
	}
}

type character struct {
	gorm.Model
	Name 		string
    Level 		int
    Health 		int
    MaxHealth 	int
	Experience 	int
	Position 	position
	Game		*game			`gorm:"many2many:game_character;association_jointable_foreignkey:game_id"`
	User 		*user	 		`gorm:"many2many:user_character;association_jointable_foreignkey:user_id"`
	Passives 	[]*passive 		`gorm:"many2many:character_passive;association_jointable_foreignkey:passive_id"`
	Skills 		[]*skill		`gorm:"many2many:character_skill;association_jointable_foreignkey:skill_id"`
}

func (selfCharacter *character) getReduction() int {
    totalReduction := 0
    for _,passiveName := range selfCharacter.Passives {
        totalReduction += passiveName.DamageReduction
    }
    return totalReduction
}


type skill struct {
	gorm.Model
    //SkillID int			`gorm:"foreignkey"`
    SkillType string
	BaseDamage int
}

type passive struct {
	gorm.Model
    //PassiveID int		`gorm:"foreignkey"`
	PassiveType string
	DamageReduction int
	DamageIncrease int
}

type classType struct {
	classtype []string
}

func attack(attacker *character,attackID int,defender *character) {
    attackingSkill := attacker.Skills[attackID]
    reduction := defender.getReduction()
    //todo: dafaque
    floatDamage := float64(attackingSkill.BaseDamage) - (float64(reduction) / 100) * float64(attackingSkill.BaseDamage)
    damage := math.Round(floatDamage)
    log.Println("attackingSkill.baseDamage ",attackingSkill.BaseDamage )
    log.Println(attacker.Name," -> ",defender.Name)
    log.Println("reduction: ",reduction, "dmg: ",damage)
    log.Println("defender health: ",defender.Health )
    defender.Health -= int(damage)
    log.Println("defender health: ",defender.Health )
}


//classTypes := []string{"dark","light"}
/*
func main() {
	server := new(server)
	server.StartServer()
	//server.initDB()
	//db := server.db
	//db.Close()



/*
	ab := true

	if ab == true {



	db.Create(&Player {
		Name: "player2",
        Level: 2,
        Health: 50,
        MaxHealth: 50,
		Experience: 2,
	})

	db.Create(&Passive {
		PassiveID: 3,
		PassiveType: "dark",
        DamageReduction: 13,
    })

	db.Create(&Skill {
    SkillID: 1,
    SkillType: "dark",
    BaseDamage: 4,
    })
	}

	 db.Create(&Game {
	Name: "jaja",
	})

	var spiel Game
	db.First(&spiel, "name = ?", "jaja")
	// add things
	
	// update that shit
	//db.Save(&player1)

	
	var player1 Player
	//db.Preload("Passives").Preload("Skills").First(&player1, "name = ?", "player2")
	db.First(&player1, "name = ?", "player2")
	spiel.Player = append(spiel.Player,&player1)
	db.Save(&spiel)
	a,_ := json.Marshal(&spiel)
	fmt.Println(string(a))




/*

	var da Skill
	db.First(&da, "skill_id = ?", 1)
	db.First(&da, 1)
	fmt.Println("sds")

	a, err := json.Marshal(&da)
	fmt.Println(string(a))
	//fmt.Println(da.SkillType)

	
	var pass1 Passive
	db.First(&pass1, "passive_id = ?", 3)
		a, err = json.Marshal(&pass1)
	fmt.Println("passive -->  ",string(a))


	var player1 Player


	//get user
	//db.First(&player1, "name = ?", "player2")
	// add things
	//player1.Skills = append(player1.Skills,&da)
	// update that shit
	//db.Save(&player1)

	
	
	db.Preload("Passives").Preload("Skills").First(&player1, "name = ?", "player2")
	db.First(&player1, "name = ?", "player2")

	a, err = json.Marshal(&player1)
	fmt.Println(string(a))




	
	playerRepo := new(playerRepo)

	playerRepo.playerArray = append(playerRepo.playerArray,initPlayer1(1))
	playerRepo.playerArray = append(playerRepo.playerArray,initPlayer(1))

    fmt.Println(playerRepo.playerArray[0].health)
    attack(playerRepo.playerArray[1],1,playerRepo.playerArray[0])
    fmt.Println(playerRepo.playerArray[0].health)
    attack(playerRepo.playerArray[1],1,playerRepo.playerArray[0])
	fmt.Println(playerRepo.playerArray[0].health)

	gm
	login
	game status
	game list
	game load <id>


	
	login -> shows gameid
	stats -> shows own stats



					pass1 :=&passive {
					//PassiveID: 3,
					PassiveType: "dark",
					DamageReduction: 13,
				}
				//selfServer.db.Create(&pass1)

				skill1 :=&skill {
				//SkillID: 1,
				SkillType: "dark",
				BaseDamage: 4,
				}
				//selfServer.db.Create(skill1)

				newChar.Passives = append(newChar.Passives,pass1)
				newChar.Skills = append(newChar.Skills,skill1)






}
*/
