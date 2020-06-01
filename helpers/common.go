package helpers

import (
	"encoding/json"
	"fmt"
	"github.com/globalsign/mgo"
	"io/ioutil"
	"net/http"
	"os"
	"task/packages/govalidator"
	"time"
)

func MongoConnection() map[string]*mgo.Collection{
	var collections = make(map[string]*mgo.Collection)
	userSession, err := mgo.Dial("mongodb://"+os.Getenv("USER_DB_HOST"))
	if userSession != nil {
		collections["user"] = userSession.DB(os.Getenv("USER_DB_NAME")).C("users")
		collections["token"] = userSession.DB(os.Getenv("USER_DB_NAME")).C("token")
		userSession.SetSocketTimeout(1 * time.Second)
	} else {
		fmt.Println("cannot connect to DB exiting ...", err)
		os.Exit(3)
	}
	return collections
}

func GetBody(res http.ResponseWriter,req *http.Request)(map[string]interface{},error){
	var requestBody map[string]interface{}
	bodyBytes, _ := ioutil.ReadAll(req.Body)
	err := json.Unmarshal(bodyBytes, &requestBody)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return nil,err
	}
	return requestBody,nil
}

func Validate(rules map[string][]string, body map[string]interface{}, messages map[string][]string) map[string][]string {
	opts := govalidator.Options{
		Body:            body,
		Rules:           rules,
		Messages:        messages,
		RequiredDefault: false,
	}
	goVal := govalidator.New(opts)
	errs := goVal.Validate()

	if errs != nil && len(errs) > 0 {
		return errs
	} else {
		return nil
	}
}

func ReplyBack(request http.ResponseWriter, responseCode int, responseBody interface{}) bool {
	bytesResponse, _ := json.Marshal(responseBody)
	request.WriteHeader(responseCode)
	_, _ = request.Write(bytesResponse)
	return true
}
