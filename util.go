package gogam

import (
	"bufio"
	"bytes"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"os"
	//"fmt"
	"encoding/json"
	"net/http/cookiejar"
)

//PostRequest sds
func PostJSONRequest(url string, bytedata []byte) (string, error) {

	log.Debug("PostRequest: trying post request to: ", url)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bytedata))
	if err != nil {
		log.Debug(err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Debug(err)
		return "", err
	}
	//defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Debug(err)
	}
	resp.Body.Close()
	log.Debug("response Body:", string(body))

	return string(body), nil
}
func PostRequest(url string, bytedata []byte, jar *cookiejar.Jar) (string, error) {

	log.Debug("PostRequest: trying post request to: ", url)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bytedata))
	if err != nil {
		log.Debug(err)
	}
	//req.Header.Set("Content-Type", "application/json")
	client := &http.Client{Jar: jar}
	resp, err := client.Do(req)
	if err != nil {
		log.Debug(err)
		return "", err
	}
	//defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Debug(err)
	}
	resp.Body.Close()
	log.Debug("response Body:", string(body))

	return string(body), nil
}

func BoardParser() *gameField {

	file, err := os.Open("map.mp")
	if err != nil {
		log.Debug(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var board [][]tile
	rowCnt := 0
	var startPoints []position

	for scanner.Scan() {

		var tem []tile
		for x, character := range scanner.Text() {

			tempTile := &tile{
				Movable:      false,
				Interactable: false,
				Info:         "Default",
				AsciArt:      character,
			}
			tempTile.initTile()
			tem = append(tem, *tempTile)
			if character == 's' {
				startPoints = append(startPoints, position{X: x, Y: rowCnt})
			}
		}
		board = append(board, tem)
		rowCnt++

		//log.Debug(scanner.Text())
	}
	log.Debug("Rows from File:", rowCnt)
	log.Debug("Startpoints at: ")
	for num, poi := range startPoints {
		log.Debug("StartPoint ", num, ": ", poi.X, ":", poi.Y)
	}
	return &gameField{
		Field:       &board,
		StartPoints: startPoints,
	}

}

func toJSON(thingy interface{}) string {

	byteArray, err := json.Marshal(thingy)
	if err != nil {
		log.Debug(err)
	}
	return string(byteArray)

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
