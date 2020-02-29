package main

import "fmt"

func main() {

	var a = "hier"
	//var a = "ab"

	switchmal(a)

}


func switchmal(a string) {
	switch a {
	case "hier":
		 machmala(a)
	case "ab":
		fmt.Println("hier is b")
	}
	fmt.Println("ende")
}

func machmala(a string) {
	fmt.Println("hier is a")
	return
}