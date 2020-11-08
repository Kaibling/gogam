package model

import "github.com/jinzhu/gorm"

type User struct {
	gorm.Model
	Name      string `gorm:"unique"`
	Password string
	Character []Character //`gorm:"foreignkey:UserID"`//`gorm:"foreignKey:id"`
}

type Character struct {
	gorm.Model
	UserID int
	CharacterName string
	Class string
	Abilities []Ability `gorm:"many2many:character_abilities;"`
}


type Ability struct {
	gorm.Model
	Name string
	Damage int
}