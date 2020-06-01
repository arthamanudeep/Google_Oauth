package controllers

import (
	"github.com/globalsign/mgo/bson"
	"github.com/thedevsaddam/govalidator"
	"net/http"
	"task/helpers"
	"task/models"
	"time"
)

var Collections = helpers.MongoConnection()

func GetUserList(res http.ResponseWriter, req *http.Request) {
	var userData []interface{}
	err := Collections["user"].Find(bson.M{}).All(&userData)
	if err != nil {
		helpers.ReplyBack(res, 400, "error while getting data")
		return
	}
	helpers.ReplyBack(res, 200, userData)
	return
}

func AddUser(res http.ResponseWriter, req *http.Request) {
	rules := govalidator.MapData{
		"email":    []string{"required","email"},
	}
	requestBody,_ := helpers.GetBody(res,req)
	errs := helpers.Validate(rules, requestBody, nil)
	if errs != nil {
		helpers.ReplyBack(res, 400, errs)
		return
	}
	email := requestBody["email"].(string)
	count, _ := Collections["user"].Find(bson.M{"email": email}).Count()
	if count != 0 {
		helpers.ReplyBack(res, 400, "email already white Listed")
		return
	}
	userData := models.UserModel{
		ID:        bson.NewObjectId(),
		Email:     email,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	insertionErr := Collections["user"].Insert(userData)
	if insertionErr != nil {
		helpers.ReplyBack(res, 400, "error while adding data")
		return
	}
	helpers.ReplyBack(res, 200, "user added successfully")
	return
}

func DeleteUser(res http.ResponseWriter, req *http.Request) {
	rules := govalidator.MapData{
		"user_id":    []string{"required","objectID"},
	}
	requestBody,_ := helpers.GetBody(res,req)
	errs := helpers.Validate(rules, requestBody, nil)
	if errs != nil {
		helpers.ReplyBack(res, 400, errs)
		return
	}
	userId := requestBody["user_id"].(string)
	loggedUserId := req.Context().Value("user_id").(string)
	if userId == loggedUserId{
		helpers.ReplyBack(res,400,"you can not delete you account")
		return
	}
	_, err6 := Collections["token"].RemoveAll(bson.M{"user_id": userId})
	if err6 != nil {
		helpers.ReplyBack(res,400,"error while logout")
		return
	}
	err := Collections["user"].Remove(bson.M{"_id":bson.ObjectIdHex(userId)})
	if err != nil{
		helpers.ReplyBack(res,400,"can not delete user")
	}
	helpers.ReplyBack(res, 200, "user deleted successfully")
	return
}

func Logout(res http.ResponseWriter, req *http.Request) {
	userId := bson.ObjectIdHex(req.Context().Value("user_id").(string))
	_, err6 := Collections["token"].RemoveAll(bson.M{"user_id": userId})
	if err6 != nil {
		helpers.ReplyBack(res,400,"error while logout")
		return
	}
	helpers.ReplyBack(res,201,"logout")
	return
}
