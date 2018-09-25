package APIs

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"golang.org/x/crypto/bcrypt"
	"github.com/csnand/APIProject/APIs/DBConfig"
	"github.com/labstack/echo"
)

type UserName struct {
	UserName string `json:"nickname"`
}

type RegisterParams struct {
	Mobile        string `json:"mobile_num"`
	Code          string `json:"code"`
	Password      string `json:"password"`
	Avatar        string `json:"avatar"`
	Cover_img     string `json:"cover_img"`
	Nickname      string `json:"nickname"`
	Gender        string `json:"gender"`
	Status        string `json:"status"`
	Introduction  string `json:"introduction"`
	Constellation string `json:"constellation"`
	Birthday      string `json:"birthday"`
}

type RegisterInfo struct {
	Mobile    string `json:"mobile_num"`
	Code      string `json:"code"`
	Password  string `json:"password"`
	Avatar    string `json:"avatar"`
	Cover_img string `json:"cover_img"`
	Profile   struct {
		Nickname      string `json:"nickname"`
		Gender        string `json:"gender"`
		Status        string `json:"status"`
		Introduction  string `json:"introduction"`
		Constellation string `json:"constellation"`
		Birthday      string `json:"birthday"`
	} `json:"Profile"`
}

type ProfileMap map[string]interface{}

type RegisterInfoDB struct {
	Mobile       string     `db:"phone"`
	Code         string     `db:"code"`
	PasswordHash string     `db:"password_hash"`
	Avatar       string     `db:"avatar"`
	Cover_img    string     `db:"cover_img"`
	Profile      ProfileMap `db:"profile"`
}

func Register(c echo.Context) (err error) {

	var reginfo RegisterInfo
	var regParams RegisterParams
	if err := c.Bind(&regParams); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if regParams.Mobile == "" || regParams.Password == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "mobile_num and password are required",
		})
	}

	if regParams.Nickname == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "nickname is required",
		})
	}

	regParams.toJSONStruct(&reginfo)

	err = insertUserProfile(reginfo)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "successfully registered",
	})
}

func VerifyUsername(c echo.Context) (err error) {

	var username UserName
	if err := c.Bind(&username); err != nil {
		return c.JSON(http.StatusBadRequest, errors.New("Params Bind Err"))
	}

	if username.UserName == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "nickname is required",
		})
	}

	DBCONN, err := DBConfig.PInitDb()
	if err != nil {
		return c.JSON(http.StatusBadRequest, errors.New("DB Initialise Err"))
	}

	defer DBCONN.Close()

	var nameNum = 999
	stmt, err := DBCONN.Prepare("select COUNT(*) from (SELECT profile->>'nickname' AS nickname FROM users) AS nameTable WHERE nameTable.nickname=$1")
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": err.Error(),
		})
	}

	err = stmt.QueryRow(username.UserName).Scan(&nameNum)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": err.Error(),
		})
	}

	fmt.Println(username.UserName)
	fmt.Println(nameNum)

	if nameNum != 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": errors.New("username exists").Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "username is unique",
	})
}

func insertUserProfile(reginfo RegisterInfo) error {

	//users and auths table need to modified to finish user registration
	//first, user profile need to be stored into user table to get user_id
	//then, using that user_id to store phone password_hash into auths table
	//note: username must be unique

	DBCONN, err := DBConfig.PInitDb()
	if err != nil {
		return err
	}

	defer DBCONN.Close()

	stmt, err := DBCONN.Prepare("INSERT INTO users (profile, cover_img, avatar) VALUES ($1, $2, $3) RETURNING id")
	if err != nil {
		return err
	}

	var userID string
	profileByte, err := json.Marshal(reginfo.Profile)
	if err != nil {
		return errors.New("parsing user profile error")
	}

	err = stmt.QueryRow(profileByte, reginfo.Cover_img, reginfo.Avatar).Scan(&userID)
	if err != nil {
		return err
	}

	stmt, err = DBCONN.Prepare("INSERT INTO auths (user_id, phone, password_hash) VALUES ($1, $2, $3)")
	if err != nil {
		return err
	}

	passHashed, err := bcrypt.GenerateFromPassword([]byte(reginfo.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(userID, reginfo.Mobile, string(passHashed))
	if err != nil {
		return err
	}

	stmt.Close()

	return nil
}

func (r RegisterParams) toJSONStruct(reginfo *RegisterInfo) {

	reginfo.Mobile = r.Mobile
	reginfo.Code = r.Code
	reginfo.Password = r.Password
	reginfo.Avatar = r.Avatar
	reginfo.Cover_img = r.Cover_img
	reginfo.Profile.Nickname = r.Nickname
	reginfo.Profile.Status = r.Status
	reginfo.Profile.Gender = r.Gender
	reginfo.Profile.Introduction = r.Introduction
	reginfo.Profile.Constellation = r.Constellation
	reginfo.Profile.Birthday = r.Birthday
}
