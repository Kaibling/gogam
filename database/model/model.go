package model

import "github.com/jinzhu/gorm"

type Envelope struct {
	Status string `json:status`
	Error string `json:"error,omitempty"`
	Message interface{} `json:"message,omitempty"`
}

type User struct {
	gorm.Model
	Name      string `gorm:"unique"`
	Password string
	Character []Character //`gorm:"foreignkey:UserID"`//`gorm:"foreignKey:id"`
}

type Character struct {
	gorm.Model
	UserID int
	Name string
	Class string
	Abilities []Ability `gorm:"many2many:character_abilities;"`
}


type Ability struct {
	gorm.Model
	Name string
	Damage int
}