package controllers

import (
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"github.com/globalsign/mgo/bson"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"io/ioutil"
	"net/http"
	"os"
	"task/helpers"
	"task/models"
	"time"
)

//Loading .env file
var (
	googleOauthConfig = &oauth2.Config{
		RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}
	randomState = "random"
)

func HomePage(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/page/login.html", http.StatusTemporaryRedirect)
}

func GoogleLogin(w http.ResponseWriter, r *http.Request) {
	url := googleOauthConfig.AuthCodeURL(randomState)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func HandleCallBack(w http.ResponseWriter, r *http.Request) {
	Collections := helpers.MongoConnection()
	if r.FormValue("state") != randomState {
		return
	}
	token, err := googleOauthConfig.Exchange(oauth2.NoContext, r.FormValue("code"))
	if err != nil {
		http.Redirect(w, r, "/login?error=no_data_from_google", http.StatusTemporaryRedirect)
		return
	}
	resp, err2 := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err2 != nil {
		http.Redirect(w, r, "/login?error=no_data_from_google", http.StatusTemporaryRedirect)
		return
	}
	defer resp.Body.Close()
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Redirect(w, r, "/login?error=no_data_from_google", http.StatusTemporaryRedirect)
		return
	}

	var userGoogleData map[string]interface{}
	err = json.Unmarshal([]byte(content), &userGoogleData)
	if err != nil {
		panic(err)
	}
	var uData map[string]interface{}
	_ = Collections["user"].Find(bson.M{"email": userGoogleData["email"]}).One(&uData)
	if uData == nil{
		http.Redirect(w, r, "/page/login.html?error=Account%20is%20not%20activated", http.StatusTemporaryRedirect)
		return
	}
	jwtSecret := os.Getenv("JWT_SECRET")
	claims := make(jwt.MapClaims)
	claims["user"] = uData
	claims["exp"] = time.Now().Add(time.Hour * 100072).Unix()
	loginToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err3 := loginToken.SignedString([]byte(jwtSecret))
	outputPayload := make(map[string]interface{})
	outputPayload["token"] = tokenStr
	outputPayload["user"] = uData
	if err3 != nil {
		http.Redirect(w, r, "/page/login.html?error="+err3.Error(), http.StatusTemporaryRedirect)
		return
	}
	tokenModel := models.TokenModel{
		ID:        bson.NewObjectId(),
		UserId:    uData["_id"].(bson.ObjectId),
		Token:     tokenStr,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err4 := Collections["token"].Insert(tokenModel)
	if err4 != nil {
		http.Redirect(w, r, "/page/login.html?error="+err4.Error(), http.StatusTemporaryRedirect)
		return
	}
	http.Redirect(w, r, "/page/userList.html?token="+tokenStr, http.StatusTemporaryRedirect)
	return
}
