package main

import (
	"net/http"

	"./APIs"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/ping", Ping)
	e.POST("/login", APIs.Login)
	e.POST("/register", APIs.Register)
	e.POST("/smscode", APIs.SendSMSFunc)
	e.POST("/smsverify", APIs.VerifySMSFunc)
	e.POST("/usernameverify", APIs.VerifyUsername)
	e.POST("/resetpassword", APIs.ResetPassword)

	// user group
	user := e.Group("/user")
	user.Use(middleware.JWT([]byte("secrets")))
	user.POST("/getprofile", APIs.GetProfile)
	user.POST("/updateprofile", APIs.UpdateProfile)
	user.POST("/updatepremium", APIs.UpdatePremium)
	user.POST("/updatecover", APIs.UpdateCover)
	user.POST("/kickoutuser", APIs.KickOutUser)

	//article group

	article := e.Group("/article")
	article.Use(middleware.JWT([]byte("secrets")))
	article.POST("/post", APIs.UploadArticle)
	article.POST("/read", APIs.ReadArticle)
	article.DELETE("/delete", APIs.DeleteArticle)
	article.POST("/update", APIs.UpdateArticle)
	article.POST("/recommand", APIs.RecommandArticle)
	article.POST("/popular", APIs.PopularArticle)
	article.POST("/collect", APIs.CollectArticle)
	article.POST("/getcollected", APIs.GetCollectedArticle)
	article.POST("/deletecollected", APIs.DeleteArticleCollect)

	//activity group
	activity := e.Group("/activity")
	activity.Use(middleware.JWT([]byte("secrets")))
	activity.POST("/post", APIs.PostActivity)
	activity.POST("/get", APIs.GetActivity)
	activity.DELETE("/delete", APIs.DeleteActivity)
	activity.POST("/update", APIs.UpdateActivity)

	activity.POST("/getbydistance", APIs.GetActivityByDistance)

	activity.POST("/join", APIs.JoinActivity)
	activity.POST("/getjoined", APIs.GetJoinedActivity)

	activity.POST("/collect", APIs.CollectActivity)
	activity.POST("/getcollected", APIs.GetCollectedActivity)
	activity.POST("/deletecollected", APIs.DeleteActivityCollect)
	activity.POST("/recommand", APIs.RecommandActivity)
	activity.POST("/popular", APIs.PopularActivity)

	tags := e.Group("/tags")
	tags.Use(middleware.JWT([]byte("secrets")))
	tags.POST("/get", APIs.GetTags)

	e.POST("/banner", APIs.Banner)

	// Restricted group - auth required to access
	r := e.Group("/auth")
	r.Use(middleware.JWT([]byte("secrets")))
	r.POST("/test", APIs.Restricted)

	e.Logger.Fatal(e.Start(":1323"))

}

func Ping(c echo.Context) error {
	return c.String(http.StatusOK, "Pong!")
}
