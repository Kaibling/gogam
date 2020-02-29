package gogam

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	log "github.com/sirupsen/logrus"
	"math"
)

type tile struct {
	Movable      bool
	Interactable bool
	Info         string
	AsciArt      rune
}

func (selfTile *tile) initTile() {

	switch selfTile.AsciArt {
	case '=':
		//Wall
		selfTile.Info = "There is a Wall"
	case ' ', 'S':
		//Floor
		selfTile.Info = "Only the Floor"
		selfTile.Movable = true

	case 'D':
		//Door
		selfTile.Info = "It's ... a Door"
		selfTile.Movable = true
		selfTile.Interactable = true
	}
}

type gameField struct {
	Name        string
	Field       *[][]tile
	StartPoints []position
}

func (selfGameField *gameField) ShowMap() {
	for _, row := range *selfGameField.Field {
		for _, tilchen := range row {
			fmt.Print(string(tilchen.AsciArt))
		}
		fmt.Println()
	}
}
func (selfGameField *gameField) getStartPosition() (position, error) {
	if len(selfGameField.StartPoints) == 0 {
		return position{}, errors.New("no Starting Points available")

	}
	returnPosition := selfGameField.StartPoints[len(selfGameField.StartPoints)-1]
	selfGameField.StartPoints = selfGameField.StartPoints[:len(selfGameField.StartPoints)-1]
	return returnPosition, nil

}

type position struct {
	X int
	Y int
}

type game struct {
	gorm.Model
	Characters    []*character `gorm:"many2many:game_character;association_jointable_foreignkey:character_id"`
	Name          string
	GameField     *gameField
	GameFiledJSON []byte
	InProgress    bool
}

func (selfGame *game) loadGameField() error {
	return json.Unmarshal(selfGame.GameFiledJSON, &selfGame.GameField)
}
func (selfGame *game) saveGameField() error {
	var err error
	selfGame.GameFiledJSON, err = json.Marshal(&selfGame.GameField)
	return err
}

func (selfGame *game) characterOverview() {
	for cnt, char := range selfGame.Characters {
		JSONChar, _ := json.Marshal(&char)
		fmt.Println(cnt, ": ", string(JSONChar))
	}

}

func (selfGame *game) addCharacter(character *character) {
	var err error
	selfGame.Characters = append(selfGame.Characters, character)
	character.Position, err = selfGame.GameField.getStartPosition()
	if err != nil {
		log.Println(err)
	}
}

type character struct {
	gorm.Model
	Name       string
	Level      int
	Health     int
	MaxHealth  int
	Experience int
	Position   position
	Game       *game      `gorm:"many2many:game_character;association_jointable_foreignkey:game_id"`
	User       *user      `gorm:"many2many:user_character;association_jointable_foreignkey:user_id"`
	Passives   []*passive `gorm:"many2many:character_passive;association_jointable_foreignkey:passive_id"`
	Skills     []*skill   `gorm:"many2many:character_skill;association_jointable_foreignkey:skill_id"`
}

func (selfCharacter *character) getReduction() int {
	totalReduction := 0
	for _, passiveName := range selfCharacter.Passives {
		totalReduction += passiveName.DamageReduction
	}
	return totalReduction
}

type skill struct {
	gorm.Model
	//SkillID int			`gorm:"foreignkey"`
	SkillType  string
	BaseDamage int
}

type passive struct {
	gorm.Model
	//PassiveID int		`gorm:"foreignkey"`
	PassiveType     string
	DamageReduction int
	DamageIncrease  int
}

type classType struct {
	classtype []string
}

func attack(attacker *character, attackID int, defender *character) {
	attackingSkill := attacker.Skills[attackID]
	reduction := defender.getReduction()
	floatDamage := float64(attackingSkill.BaseDamage) - (float64(reduction)/100)*float64(attackingSkill.BaseDamage)
	damage := math.Round(floatDamage)
	log.Println("attackingSkill.baseDamage ", attackingSkill.BaseDamage)
	log.Println(attacker.Name, " -> ", defender.Name)
	log.Println("reduction: ", reduction, "dmg: ", damage)
	log.Println("defender health: ", defender.Health)
	defender.Health -= int(damage)
	log.Println("defender health: ", defender.Health)
}

/*
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

}
*/
