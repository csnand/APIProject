package APIs

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"

	"./DBConfig"
	"github.com/labstack/echo"
)

const (
	APIKEY = "your api key"
	SMSURL = "https://sms.yunpian.com/v2/sms/single_send.json"
)

var numsRunes = []rune("0123456789")
var msg string = "your sms template"

type ReceivedStruct struct {
	Code   int     `json:"code"`
	Msg    string  `json:"msg"`
	Count  int     `json:"count"`
	Fee    float64 `json:"fee"`
	Unit   string  `json:"unit"`
	Mobile string  `json:"mobile"`
	Sid    int64   `json:"sid"`
}

type SmsInfo struct {
	Mobile   string    `json:"mobile_num"`
	Code     string    `json:"code"`
	Message  string    `json:"message"`
	ExpireAt time.Time `json:"expire_at"`
}

func SendSMSFunc(c echo.Context) (err error) {

	var smsinfo SmsInfo
	//Bind with Map works only if Header set into "Content-Type: application/json"
	if err := c.Bind(&smsinfo); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if smsinfo.Mobile == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "mobile_num is required",
		})
	}

	smscode := randNumStr()
	msgWithCode := strings.Replace(msg, "code", smscode, 1)

	sendinfo := url.Values{"apikey": {APIKEY}, "mobile": {smsinfo.Mobile}, "text": {msgWithCode}}

	resp, err := http.PostForm(SMSURL, sendinfo)

	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	var receivedinfo ReceivedStruct
	json.Unmarshal(body, &receivedinfo)

	fmt.Println(smscode)

	smsinfo.Code = smscode
	smsinfo.Message = receivedinfo.Msg
	smsinfo.ExpireAt = time.Now().Add(time.Minute * 5)

	err = addSMSRecord(smsinfo)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, smsinfo)
}

func addSMSRecord(smsinfo SmsInfo) error {

	DBCONN, err := DBConfig.PInitDb()
	if err != nil {
		return err
	}

	defer DBCONN.Close()

	fmt.Println("Inserting sms info into db")

	stmt, err := DBCONN.Prepare("INSERT INTO smsverify (mobile_num, code, isVerified, expired_at) VALUES ($1, $2, $3, $4)")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(smsinfo.Mobile, smsinfo.Code, false, smsinfo.ExpireAt)
	if err != nil {
		return err
	}

	return nil
}

func randNumStr() string {
	rand.Seed(time.Now().UnixNano())
	str := make([]rune, 4)
	for i := range str {
		str[i] = numsRunes[rand.Intn(10)]
	}

	return string(str)
}

func VerifySMSFunc(c echo.Context) (err error) {

	var smsinfo SmsInfo
	if err := c.Bind(&smsinfo); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if smsinfo.Mobile == "" || smsinfo.Code == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "mobile_num and code are required",
		})
	}

	err = smsverify(smsinfo.Mobile, smsinfo.Code)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "code was not match or expired",
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "code matched successfully",
	})
}

func smsverify(mobile string, code string) error {

	DBCONN, err := DBConfig.PInitDb()
	if err != nil {
		return err
	}
	defer DBCONN.Close()

	stmt, err := DBCONN.Prepare("DELETE FROM smsverify WHERE expired_at < now()")
	if err != nil {
		return err
	}
	_, err = stmt.Exec()
	if err != nil {
		return err
	}

	stmt, err = DBCONN.Prepare("SELECT code FROM smsverify WHERE mobile_num=$1 ORDER BY expired_at DESC")
	if err != nil {
		return err
	}

	var codeFromDB string
	err = stmt.QueryRow(mobile).Scan(&codeFromDB)
	if err != nil {
		return err
	}

	fmt.Println(code)
	fmt.Println(codeFromDB)

	if code != codeFromDB {
		return errors.New("code was not match or expired")
	}

	stmt.Close()

	return nil
}
