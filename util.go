package gogam

import (
	"log"
	"bytes"
	 "io/ioutil"
	 "net/http"
	 "os"
	 "bufio"
	 //"fmt"
 	"net/http/cookiejar"
)

//PostRequest sds
func PostJSONRequest(url string, bytedata []byte) (string,error) {

	log.Println("PostRequest: trying post request to: ", url)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bytedata))
	if err != nil {
		log.Println(err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return "",err
	}
	//defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}
	resp.Body.Close()
	log.Println("response Body:", string(body))

	return string(body),nil
}
func PostRequest(url string, bytedata []byte, jar *cookiejar.Jar) (string,error) {

	log.Println("PostRequest: trying post request to: ", url)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bytedata))
	if err != nil {
		log.Println(err)
	}
	//req.Header.Set("Content-Type", "application/json")
	client := &http.Client{Jar: jar}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return "",err
	}
	//defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}
	resp.Body.Close()
	log.Println("response Body:", string(body))

	return string(body),nil
}

func BoardParser() *gameField{ 

	file , err := os.Open("map.mp")
	if err != nil {
		log.Println(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var board [][]tile
	rowCnt := 0
	var startPoints []position


	for scanner.Scan() {

		var tem []tile
		for x,character := range scanner.Text() {

			tempTile := &tile{
				movable: false,
				interactable: false,
				info: "Default",
				asciArt: character,

			}
			tempTile.initTile()
			tem = append(tem,*tempTile)
			if character == 's' {
				startPoints = append(startPoints,position {X: x,Y:rowCnt}) 
			}
		}
		board = append(board,tem)
		rowCnt++

		//log.Println(scanner.Text())
	}
	log.Println("Rows from File:",rowCnt)
	log.Println("Startpoints at: ")
	for num,poi := range startPoints {
		log.Println("StartPoint ",num,": ",poi.X,":",poi.Y)
	}
	return &gameField {
		Field:&board,
		startPoints: startPoints,
	}

}

/*
func interpretFileMap(char rune) *Tile {
	returnTile := &Tile{
		movable: false,
		interactable: false,
		info: "Default",
	}
	switch char {
	case '=':
		//Wall
		returnTile.info ="There is a Wall"
	case ' ':
		//Floor
		returnTile.info ="Only the Floor"
		returnTile.movable = true
	
	case 'D':
		//Door
		returnTile.info ="It's ... a Door"
		returnTile.movable = true
		returnTile.interactable = true
	}
	return returnTile
}
*/