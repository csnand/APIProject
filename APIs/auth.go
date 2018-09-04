package APIs

import (
	"fmt"
	"net/http"
	"time"

	"./DBConfig"
	"github.com/labstack/echo"
	"golang.org/x/crypto/bcrypt"
)

type userinfo struct {
	Password string `json:"password"`
	Mobile   string `json:"mobile_num"`
}

type getUser struct {
	isPassword bool
	userID     string
}

func Login(c echo.Context) (err error) {

	var user userinfo

	//Bind with Map works only if Header set into "Content-Type: application/json"
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": err.Error(),
		})
	}

	ok, err := getUserPass(user)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": err.Error(),
		})
	}

	if ok.isPassword {
		// Create token
		token := jwt.New(jwt.SigningMethodHS256)

		// Set claims
		claims := token.Claims.(jwt.MapClaims)

		claims["user_id"] = ok.userID
		// claims["admin"] = false

		//set token expire time
		claims["exp"] = time.Now().Add(time.Hour * 24 * 30).Unix()

		// Generate encoded token and send it as response.
		t, err := token.SignedString([]byte("secrets"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"message": err.Error(),
			})
		}
		return c.JSON(http.StatusOK, map[string]string{
			"token":   t,
			"user_id": ok.userID,
		})
	}

	return echo.ErrUnauthorized
}

func Restricted(c echo.Context) error {
	// user := c.Get("user").(*jwt.Token)
	// claims := user.Claims.(jwt.MapClaims)
	// name := claims["name"].(string)
	return c.String(http.StatusOK, "Welcome!")
}

func ResetPassword(c echo.Context) (err error) {

	var userpass userinfo
	if err := c.Bind(&userpass); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": err.Error(),
		})
	}

	if userpass.Mobile == "" || userpass.Password == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "mobile_num and password are required",
		})
	}

	err = resetUserPassword(userpass)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "password reset successfully",
	})
}

func resetUserPassword(userpass userinfo) error {

	DBCONN, err := DBConfig.PInitDb()
	if err != nil {
		return err
	}
	defer DBCONN.Close()

	passwordHashed, err := bcrypt.GenerateFromPassword([]byte(userpass.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	stmt, err := DBCONN.Prepare("UPDATE auths SET password_hash=$1 WHERE user_id=(SELECT user_id FROM auths WHERE phone=$2)")
	if err != nil {
		return err
	}

	res, err := stmt.Exec(string(passwordHashed), userpass.Mobile)
	if err != nil {
		return err
	}

	fmt.Println(res)

	return nil
}

func getUserPass(user userinfo) (getUser, error) {

	getuser := getUser{false, ""}

	DBCONN, err := DBConfig.PInitDb()
	if err != nil {
		return getuser, nil
	}

	defer DBCONN.Close()

	stmt, err := DBCONN.Prepare("SELECT password_hash, user_id FROM auths WHERE phone=$1")
	if err != nil {
		return getuser, err
	}

	var password_hash string
	err = stmt.QueryRow(user.Mobile).Scan(&password_hash, &getuser.userID)
	if err != nil {
		return getuser, err
	}

	fmt.Println(password_hash)

	err = bcrypt.CompareHashAndPassword([]byte(password_hash), []byte(user.Password))
	if err != nil {
		return getuser, err
	}

	getuser.isPassword = true

	return getuser, nil
}
