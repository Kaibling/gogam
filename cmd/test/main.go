package main

import "fmt"

func main() {
	a := "h"
	b := 2
	

	fmt.Println(getIt(a,b))
}


func getIt(a string, b int) string {
	var res string
		switch a {
	case a:
		if b == 2 {
			fmt.Println("da is")
			res = "b is 2"
			return res
		}
		res = "a is"
	case "b":
		res  = "b is"
	}
	return res
}