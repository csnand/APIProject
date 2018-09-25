package APIs

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"github.com/csnand/APIProject/APIs/DBConfig"
	"github.com/labstack/echo"
)

type UserProfieParams struct {
	ID            string `json:"user_id"`
	Avatar        string `json:"avatar"`
	Cover_img     string `json:"cover_img"`
	Premium       bool   `json:"premium"`
	Nickname      string `json:"nickname"`
	Gender        string `json:"gender"`
	Status        string `json:"status"`
	Introduction  string `json:"introduction"`
	Constellation string `json:"constellation"`
	Birthday      string `json:"birthday"`
}

type UserProfie struct {
	ID        string `json:"user_id"`
	Avatar    string `json:"avatar"`
	Cover_img string `json:"cover_img"`
	Premium   bool   `json:"premium"`
	Profile   struct {
		Nickname      string `json:"nickname"`
		Gender        string `json:"gender"`
		Status        string `json:"status"`
		Introduction  string `json:"introduction"`
		Constellation string `json:"constellation"`
		Birthday      string `json:"birthday"`
	} `json:"Profile"`
	Created int `json:"created"`
	Joined  int `json:"joined"`
}

type UserAvatar struct {
	ID     string `json:"user_id"`
	Avatar string `json:"avatar"`
}

type UserCoverIMG struct {
	ID        string `json:"user_id"`
	Cover_img string `json:"cover_img"`
}

type UserPremium struct {
	ID      string `json:"user_id"`
	Premium bool   `json:"premium"`
}

type UserProfileDB struct {
	ID        string     `db:"id"`
	Premium   bool       `db:"premium"`
	Avatar    string     `db:"avatar"`
	Cover_img string     `db:"cover_img"`
	Profile   ProfileMap `db:"profile"`
}

type UserID struct {
	UserID string `json:"user_id"`
}

type KickUserInfo struct {
	ID           string `json:"user_id"`
	KickedUserID string `json:"kicked_user_id"`
	ActivityID   string `json:"activity_id"`
}

func GetProfile(c echo.Context) (err error) {

	var userID UserID
	if err := c.Bind(&userID); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "user_id is required",
		})
	}

	var userProfile UserProfie
	err = getUserProfile(userID.UserID, &userProfile)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": err.Error(),
		})
	}

	err = getUserStats(userID.UserID, &userProfile)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, userProfile)
}

func UpdateProfile(c echo.Context) (err error) {

	var userProfileParams UserProfieParams
	var userProfile UserProfie
	if err := c.Bind(&userProfileParams); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if userProfileParams.ID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "user_id is required",
		})
	}

	userProfileParams.toJSONStruct(&userProfile)

	err = updateUserProfile(userProfile)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "profile updated successfully",
	})
}

func UpdateCover(c echo.Context) (err error) {

	var userCover UserCoverIMG
	if err := c.Bind(&userCover); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if userCover.ID == "" || userCover.Cover_img == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "user_id and cover_img are required",
		})
	}

	err = updateUserCover(userCover)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "cover_img updated successfully",
	})
}

func UpdatePremium(c echo.Context) (err error) {

	var userPremium UserPremium
	if err := c.Bind(&userPremium); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if userPremium.ID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "user_id and premium are required",
		})
	}

	err = updateUserPremium(userPremium)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "premium updated successfully",
	})
}

func KickOutUser(c echo.Context) (err error) {

	var kickout KickUserInfo
	if err := c.Bind(&kickout); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": err.Error(),
		})
	}

	if kickout.ID == "" || kickout.ActivityID == "" || kickout.KickedUserID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "user_id, activity_id and kicked_user_id are required",
		})
	}

	err = kickUserOut(kickout)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "kick out successfully",
	})
}

func kickUserOut(kickout KickUserInfo) error {

	DBCONN, err := DBConfig.PInitDb()
	if err != nil {
		return err
	}
	defer DBCONN.Close()

	stmt, err := DBCONN.Prepare("SELECT creator_id FROM activities WHERE id=$1")
	if err != nil {
		return err
	}

	var creator string
	err = stmt.QueryRow(kickout.ActivityID).Scan(&creator)
	if err != nil {
		return err
	}

	if kickout.ID != creator {
		return errors.New("only creator can kick user out")
	}

	stmt, err = DBCONN.Prepare("INSERT INTO activity_kick_histories (user_id, activity_id, kick_user_id) VALUES ($1, $2, $3)")
	if err != nil {
		return err
	}

	res, err := stmt.Exec(kickout.ID, kickout.ActivityID, kickout.KickedUserID)
	if err != nil {
		return err
	}
	fmt.Println(res)

	stmt, err = DBCONN.Prepare("DELETE FROM activity_members WHERE user_id=$1 AND activity_id=$2")
	if err != nil {
		return err
	}

	res, err = stmt.Exec(kickout.KickedUserID, kickout.ActivityID)
	if err != nil {
		return err
	}
	fmt.Println(res)

	stmt.Close()

	return nil
}

func updateUserPremium(userPremium UserPremium) error {

	DBCONN, err := DBConfig.PInitDb()
	if err != nil {
		return err
	}
	defer DBCONN.Close()

	stmt, err := DBCONN.Prepare("UPDATE users SET premium=$1 WHERE id=$2")
	if err != nil {
		return err
	}

	res, err := stmt.Exec(userPremium.Premium, userPremium.ID)
	if err != nil {
		return err
	}

	fmt.Println(res)

	stmt.Close()

	return nil
}

func updateUserCover(userCover UserCoverIMG) error {

	DBCONN, err := DBConfig.PInitDb()
	if err != nil {
		return err
	}
	defer DBCONN.Close()

	stmt, err := DBCONN.Prepare("UPDATE users SET cover_img=$1 WHERE id=$2")
	if err != nil {
		return err
	}

	res, err := stmt.Exec(userCover.Cover_img, userCover.ID)
	if err != nil {
		return err
	}

	fmt.Println(res)

	stmt.Close()

	return nil
}

func updateUserProfile(userProfile UserProfie) error {

	DBCONN, err := DBConfig.PInitDb()
	if err != nil {
		return err
	}
	defer DBCONN.Close()

	stmt, err := DBCONN.Prepare("UPDATE users SET profile=$1, avatar=$2 WHERE id=$3")
	if err != nil {
		return err
	}

	profileByte, err := json.Marshal(userProfile.Profile)
	if err != nil {
		return errors.New("parsing user profile error")
	}

	res, err := stmt.Exec(profileByte, userProfile.Avatar, userProfile.ID)
	if err != nil {
		return err
	}

	fmt.Println(res)
	stmt.Close()

	return nil
}

func getUserProfile(userID string, userProfile *UserProfie) error {

	DBCONN, err := DBConfig.PInitDb()
	if err != nil {
		return err
	}

	defer DBCONN.Close()

	stmt, err := DBCONN.Prepare("SELECT id, premium, avatar, cover_img, profile FROM users WHERE id=$1")
	if err != nil {
		return err
	}

	var userProfileDB UserProfileDB
	err = stmt.QueryRow(userID).Scan(&userProfileDB.ID,
		&userProfileDB.Premium,
		&userProfileDB.Avatar,
		&userProfileDB.Cover_img,
		&userProfileDB.Profile)

	userProfileDB.toJSONStruct(userProfile)

	stmt.Close()

	return nil
}

func getUserStats(userID string, userProfile *UserProfie) error {

	DBCONN, err := DBConfig.PInitDb()
	if err != nil {
		return err
	}
	defer DBCONN.Close()

	var created, joined int
	stmt, err := DBCONN.Prepare("SELECT COUNT(*) FROM activities WHERE creator_id=$1")
	if err != nil {
		return err
	}

	err = stmt.QueryRow(userID).Scan(&created)
	if err != nil {
		return err
	}

	stmt, err = DBCONN.Prepare("SELECT COUNT(*) FROM activity_members WHERE user_id=$1")
	if err != nil {
		return err
	}

	err = stmt.QueryRow(userID).Scan(&joined)
	if err != nil {
		return err
	}

	userProfile.Created = created
	userProfile.Joined = joined

	return nil
}

func (p ProfileMap) Value() (driver.Value, error) {
	j, err := json.Marshal(p)
	return j, err
}

func (p *ProfileMap) Scan(src interface{}) error {
	source, ok := src.([]byte)
	if !ok {
		return errors.New("Type assertion .([]byte) failed.")
	}

	var i interface{}
	err := json.Unmarshal(source, &i)
	if err != nil {
		return err
	}

	*p, ok = i.(map[string]interface{})
	if !ok {
		return errors.New("Type assertion .(map[string]interface{}) failed.")
	}

	return nil
}

func (userProfileDB UserProfileDB) toJSONStruct(userProfile *UserProfie) {

	userProfile.ID = userProfileDB.ID
	userProfile.Premium = userProfileDB.Premium
	userProfile.Avatar = userProfileDB.Avatar
	userProfile.Cover_img = userProfileDB.Cover_img

	if str, ok := userProfileDB.Profile["nickname"].(string); ok {
		userProfile.Profile.Nickname = str
	} else {
		userProfile.Profile.Nickname = ""
	}

	if str, ok := userProfileDB.Profile["gender"].(string); ok {
		userProfile.Profile.Gender = str
	} else {
		userProfile.Profile.Gender = ""
	}

	if str, ok := userProfileDB.Profile["status"].(string); ok {
		userProfile.Profile.Status = str
	} else {
		userProfile.Profile.Status = ""
	}

	if str, ok := userProfileDB.Profile["introduction"].(string); ok {
		userProfile.Profile.Introduction = str
	} else {
		userProfile.Profile.Introduction = ""
	}

	if str, ok := userProfileDB.Profile["constellation"].(string); ok {
		userProfile.Profile.Constellation = str
	} else {
		userProfile.Profile.Constellation = ""
	}

	if str, ok := userProfileDB.Profile["birthday"].(string); ok {
		userProfile.Profile.Birthday = str
	} else {
		userProfile.Profile.Birthday = ""
	}
}

func (u UserProfieParams) toJSONStruct(userProfile *UserProfie) {

	userProfile.ID = u.ID
	userProfile.Avatar = u.Avatar
	userProfile.Profile.Nickname = u.Nickname
	userProfile.Profile.Status = u.Status
	userProfile.Profile.Gender = u.Gender
	userProfile.Profile.Introduction = u.Introduction
	userProfile.Profile.Constellation = u.Constellation
	userProfile.Profile.Birthday = u.Birthday
}
