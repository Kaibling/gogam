package utility

import (
	"strconv"
	"net/http"
	"bytes"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	    "crypto/sha1"
	"encoding/base64"
	"encoding/json"

)
func PrettyJSON(object interface{}) string {
	a, _ := json.MarshalIndent(object, "", " ")
	return string(a)
}

func StringToUint(data string) uint {
		u64, _ := strconv.ParseUint(data, 10,32)
	return uint(u64)
}

func HashPassword(password string) string {
    hasher := sha1.New()
    hasher.Write([]byte(password))
    return base64.URLEncoding.EncodeToString(hasher.Sum(nil))	
}


func PostRequest(url string, bytedata []byte) (string, error) {

	log.Debug("PostRequest: trying post request to: ", url)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bytedata))
	if err != nil {
		log.Debug(err)
	}
	//req.Header.Set("Content-Type", "application/json")
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