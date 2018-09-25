package APIs

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"time"
	"github.com/csnand/APIProject/APIs/DBConfig"
	"github.com/labstack/echo"
	"github.com/lib/pq"
)

//(SELECT CAST(CASE WHEN ( SELECT COUNT(*) FROM activity_collections WHERE user_id=$1 AND article_id=articles.id ) = 1 THEN true ELSE false END AS boolean) AS isCollected)

type ArticleInfo struct {
	Id          int       `json:"article_id"`
	ActivityId  string    `json:"activity_id"`
	UserId      string    `json:"user_id"`
	Title       string    `json:"title"`
	Content     string    `json:"content"`
	Url         string    `json:"url"`
	CoverImg    string    `json:"cover_img"`
	Desc        string    `json:"desc"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Tags        []string  `json:"tags"`
	Popular     bool      `json:"popular"`
	IsCollected bool      `json:"iscollected"`
}

type ArticleReturn struct {
	Id          int            `json:"article_id"`
	ActivityId  sql.NullString `json:"activity_id"`
	UserId      string         `json:"user_id"`
	Title       string         `json:"title"`
	Content     string         `json:"content"`
	CoverImg    string         `json:"cover_img"`
	Desc        string         `json:"desc"`
	CreatedAt   string         `json:"created_at"`
	UpdatedAt   string         `json:"updated_at"`
	Tags        []string       `json:"tags"`
	Popular     bool           `json:"popular"`
	IsCollected bool           `json:"iscollected"`
}

type ArticleArray struct {
	ArticleList []struct {
		Id          int            `json:"article_id"`
		ActivityId  sql.NullString `json:"activity_id"`
		UserId      string         `json:"user_id"`
		Title       string         `json:"title"`
		Content     string         `json:"content"`
		CoverImg    string         `json:"cover_img"`
		Desc        string         `json:"desc"`
		CreatedAt   string         `json:"created_at"`
		UpdatedAt   string         `json:"updated_at"`
		Tags        []string       `json:"tags"`
		Popular     bool           `json:"popular"`
		IsCollected bool           `json:"iscollected"`
	} `json:"article_list"`
}

type ArticleQuery struct {
	Id         int      `json:"article_id"`
	ActivityId string   `json:"activity_id"`
	UserId     string   `json:"user_id"`
	Title      string   `json:"title"`
	Tags       []string `json:"tags"`
	Popular    bool     `json:"popular"`
	Mode       string   `json:"mode"`
}

type ArticleId struct {
	Id int `json:"article_id"`
}

func UploadArticle(c echo.Context) (err error) {

	var article ArticleInfo

	if err := c.Bind(&article); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if article.UserId == "" || article.Title == "" || article.Content == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "user_id and title are required, article content cannot be null",
		})
	}

	if article.Popular != true && article.Popular != false {
		article.Popular = false
	}

	article.CreatedAt = time.Now()
	article.UpdatedAt = time.Now()

	err = articleModel(article, 1)

	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Article Uploaded Successfully",
	})

}

func UpdateArticle(c echo.Context) (err error) {

	var article ArticleInfo

	if err := c.Bind(&article); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if article.Id == 0 || article.Url == "" || article.Content == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "article_id and url are required, article content cannot be null",
		})
	}

	article.UpdatedAt = time.Now()

	err = articleModel(article, 2)

	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Article Updated Successfully",
	})

}

func DeleteArticle(c echo.Context) (err error) {

	var id ArticleId
	if err := c.Bind(&id); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if id.Id == 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "article_id cannot be zero",
		})
	}

	err = articleDelete(id.Id)

	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Article Deleted",
	})

}

func ReadArticle(c echo.Context) (err error) {

	var aQuery ArticleQuery
	var aList ArticleArray
	if err := c.Bind(&aQuery); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if aQuery.UserId == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "user_id is required",
		})
	}

	err = articleRead(aQuery, &aList)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, aList)
}

func RecommandArticle(c echo.Context) (err error) {

	var aQuery ArticleQuery
	var aList ArticleArray

	if err := c.Bind(&aQuery); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if aQuery.UserId == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "user_id is required",
		})
	}

	err = articleGetRecommand(aQuery, &aList)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, aList)
}

func PopularArticle(c echo.Context) (err error) {

	var aQuery ArticleQuery
	var aList ArticleArray

	if err := c.Bind(&aQuery); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if aQuery.UserId == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "user_id is required",
		})
	}

	err = articleGetPopular(aQuery, &aList)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, aList)
}

func CollectArticle(c echo.Context) (err error) {

	var aQuery ArticleQuery
	if err := c.Bind(&aQuery); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if aQuery.Id == 0 || aQuery.UserId == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "article_id and user_id is required",
		})
	}

	err = articleCollect(aQuery)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "successfully collected article",
	})

}

func GetCollectedArticle(c echo.Context) (err error) {

	var aQuery ArticleQuery
	var aList ArticleArray

	if err := c.Bind(&aQuery); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if aQuery.UserId == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "user_id is required",
		})
	}

	err = articleCollectGet(aQuery, &aList)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, aList)
}

func DeleteArticleCollect(c echo.Context) (err error) {

	var aQuery ArticleQuery
	if err := c.Bind(&aQuery); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": err.Error(),
		})
	}

	if aQuery.Id == 0 || aQuery.UserId == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "article_id and user_id are required",
		})
	}

	err = articleCollectDelete(aQuery)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "collection deleted",
	})
}

func articleCollectDelete(aQuery ArticleQuery) error {

	DBCONN, err := DBConfig.PInitDb()
	if err != nil {
		return err
	}
	defer DBCONN.Close()

	stmt, err := DBCONN.Prepare("DELETE FROM article_collections WHERE article_id=$1 AND user_id=$2")
	if err != nil {
		return err
	}

	res, err := stmt.Exec(aQuery.Id, aQuery.UserId)
	if err != nil {
		return err
	}

	fmt.Println(res)
	stmt.Close()

	return nil
}

func articleCollectGet(aQuery ArticleQuery, aList *ArticleArray) error {

	DBCONN, err := DBConfig.PInitDb()
	if err != nil {
		return err
	}

	defer DBCONN.Close()

	var stmt *sql.Stmt
	var rows *sql.Rows

	stmt, err = DBCONN.Prepare("SELECT id, activity_id, user_id, title, cover_img, created_at, updated_at, content, tags, popular FROM articles WHERE id IN (SELECT article_id FROM article_collections WHERE user_id=$1)")
	if err != nil {
		return nil
	}

	rows, err = stmt.Query(aQuery.UserId)
	if err != nil {
		return nil
	}

	for rows.Next() {
		var ca, ua time.Time
		var articleReturn ArticleReturn
		err = rows.Scan(&articleReturn.Id, &articleReturn.ActivityId,
			&articleReturn.UserId, &articleReturn.Title,
			&articleReturn.CoverImg, &ca, &ua,
			&articleReturn.Content,
			pq.Array(&articleReturn.Tags),
			&articleReturn.Popular)

		if err != nil {
			return err
		}

		articleReturn.CreatedAt = ca.String()
		articleReturn.UpdatedAt = ua.String()

		aList.ArticleList = append(aList.ArticleList, articleReturn)
	}
	return nil
}

func articleCollect(aQuery ArticleQuery) error {

	DBCONN, err := DBConfig.PInitDb()
	if err != nil {
		return err
	}

	defer DBCONN.Close()

	var stmt *sql.Stmt
	stmt, err = DBCONN.Prepare("INSERT INTO article_collections (article_id, user_id) VALUES ($1, $2)")
	if err != nil {
		return err
	}

	res, err := stmt.Exec(aQuery.Id, aQuery.UserId)
	if err != nil {
		return nil
	}

	stmt.Close()
	fmt.Println(res)

	return nil
}

func articleGetRecommand(aQuery ArticleQuery, aList *ArticleArray) error {

	DBCONN, err := DBConfig.PInitDb()
	if err != nil {
		return err
	}

	defer DBCONN.Close()

	var stmt *sql.Stmt
	var rows *sql.Rows

	stmt, err = DBCONN.Prepare("SELECT id, activity_id, user_id, title, cover_img, created_at, updated_at, content, tags, popular, " +
		"(SELECT CAST(CASE WHEN ( SELECT COUNT(*) FROM article_collections WHERE user_id=$1 AND article_id=articles.id ) = 1 THEN true ELSE false END AS boolean) AS isCollected)" +
		" FROM articles WHERE popular=false")
	if err != nil {
		return err
	}

	rows, err = stmt.Query(aQuery.UserId)
	if err != nil {
		return err
	}

	for rows.Next() {

		var ca, ua time.Time
		var articleReturn ArticleReturn
		err = rows.Scan(&articleReturn.Id, &articleReturn.ActivityId,
			&articleReturn.UserId, &articleReturn.Title,
			&articleReturn.CoverImg, &ca, &ua,
			&articleReturn.Content,
			pq.Array(&articleReturn.Tags),
			&articleReturn.Popular,
			&articleReturn.IsCollected)

		if err != nil {
			return err
		}

		articleReturn.CreatedAt = ca.String()
		articleReturn.UpdatedAt = ua.String()

		aList.ArticleList = append(aList.ArticleList, articleReturn)
	}
	rows.Close()
	stmt.Close()

	return nil
}

func articleGetPopular(aQuery ArticleQuery, aList *ArticleArray) error {

	DBCONN, err := DBConfig.PInitDb()
	if err != nil {
		return err
	}

	defer DBCONN.Close()

	var stmt *sql.Stmt
	var rows *sql.Rows

	stmt, err = DBCONN.Prepare("SELECT id, activity_id, user_id, title, cover_img, created_at, updated_at, content, tags, popular, " +
		"(SELECT CAST(CASE WHEN ( SELECT COUNT(*) FROM article_collections WHERE user_id=$1 AND article_id=articles.id ) = 1 THEN true ELSE false END AS boolean) AS isCollected)" +
		" FROM articles WHERE popular=true")
	if err != nil {
		return err
	}

	rows, err = stmt.Query(aQuery.UserId)
	if err != nil {
		return err
	}

	for rows.Next() {

		var ca, ua time.Time
		var articleReturn ArticleReturn
		err = rows.Scan(&articleReturn.Id, &articleReturn.ActivityId,
			&articleReturn.UserId, &articleReturn.Title,
			&articleReturn.CoverImg, &ca, &ua,
			&articleReturn.Content,
			pq.Array(&articleReturn.Tags),
			&articleReturn.Popular,
			&articleReturn.IsCollected)

		if err != nil {
			return err
		}

		articleReturn.CreatedAt = ca.String()
		articleReturn.UpdatedAt = ua.String()

		aList.ArticleList = append(aList.ArticleList, articleReturn)
	}

	rows.Close()
	stmt.Close()

	return nil
}

func articleRead(aQuery ArticleQuery, aList *ArticleArray) error {

	DBCONN, err := DBConfig.PInitDb()
	if err != nil {
		return err
	}

	defer DBCONN.Close()

	var stmt *sql.Stmt
	var rows *sql.Rows

	if aQuery.Id != 0 {
		stmt, err = DBCONN.Prepare("SELECT id, activity_id, user_id, title, cover_img, created_at, updated_at, content, tags, popular, " +
			"(SELECT CAST(CASE WHEN ( SELECT COUNT(*) FROM article_collections WHERE user_id=$1 AND article_id=articles.id ) = 1 THEN true ELSE false END AS boolean) AS isCollected)" +
			" FROM articles WHERE id=$2")
		if err != nil {
			return err
		}

		rows, err = stmt.Query(aQuery.UserId, aQuery.Id)
		if err != nil {
			return err
		}

	} else if aQuery.Mode == "creator" {
		stmt, err = DBCONN.Prepare("SELECT id, activity_id, user_id, title, cover_img, created_at, updated_at, content, tags, popular, " +
			"(SELECT CAST(CASE WHEN ( SELECT COUNT(*) FROM article_collections WHERE user_id=$1 AND article_id=articles.id ) = 1 THEN true ELSE false END AS boolean) AS isCollected)" +
			" FROM articles WHERE user_id=$1")
		if err != nil {
			return err
		}

		rows, err = stmt.Query(aQuery.UserId)
		if err != nil {
			return err
		}

	} else {
		stmt, err = DBCONN.Prepare("SELECT id, activity_id, user_id, title, cover_img, created_at, updated_at, content, tags, popular, " +
			"(SELECT CAST(CASE WHEN ( SELECT COUNT(*) FROM article_collections WHERE user_id=$1 AND article_id=articles.id ) = 1 THEN true ELSE false END AS boolean) AS isCollected)" +
			" FROM articles")
		if err != nil {
			return err
		}

		rows, err = stmt.Query(aQuery.UserId)
		if err != nil {
			return err
		}
	}

	for rows.Next() {

		var ca, ua time.Time
		var articleReturn ArticleReturn
		err = rows.Scan(&articleReturn.Id, &articleReturn.ActivityId,
			&articleReturn.UserId, &articleReturn.Title,
			&articleReturn.CoverImg, &ca, &ua,
			&articleReturn.Content,
			pq.Array(&articleReturn.Tags),
			&articleReturn.Popular,
			&articleReturn.IsCollected)

		if err != nil {
			return err
		}

		articleReturn.CreatedAt = ca.String()
		articleReturn.UpdatedAt = ua.String()

		aList.ArticleList = append(aList.ArticleList, articleReturn)
	}
	rows.Close()
	stmt.Close()

	return nil
}

func articleModel(info ArticleInfo, mode int) error {
	// if mode == 1 -> insert, if mode == 2 -> Update, else error

	var err error
	if mode == 1 {
		err = articleInsert(info)
	} else if mode == 2 {
		err = articleUpdate(info)
	} else {
		err = errors.New("mode must be either 1 or 2")
	}

	if err != nil {
		return err
	}

	return nil
}

func articleUpdate(info ArticleInfo) error {

	DBCONN, err := DBConfig.PInitDb()
	if err != nil {
		return err
	}

	defer DBCONN.Close()

	fmt.Println("# Updating values")

	stmt, err := DBCONN.Prepare("UPDATE articles SET title=$1, content=$2, cover_img=$3, desciption=$4, updated_at=$5 WHERE id=$6")

	if err != nil {
		return err
	}

	res, err := stmt.Exec(info.Title,
		info.Content,
		info.CoverImg,
		info.Desc,
		info.UpdatedAt,
		info.Id)

	fmt.Println(res)

	if err != nil {
		return err
	}

	stmt.Close()
	return nil

}

func articleInsert(info ArticleInfo) error {

	DBCONN, err := DBConfig.PInitDb()
	if err != nil {
		return err
	}

	defer DBCONN.Close()

	fmt.Println("# Inserting values")

	stmt, err := DBCONN.Prepare("INSERT INTO articles (user_id, title, content, cover_img, desciption, created_at, updated_at, tags, popular) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)")

	if err != nil {
		return err
	}

	res, err := stmt.Exec(info.UserId,
		info.Title,
		info.Content,
		info.CoverImg,
		info.Desc,
		info.CreatedAt,
		info.UpdatedAt,
		pq.Array(info.Tags),
		info.Popular)

	fmt.Println(res)

	if err != nil {
		return err
	}

	stmt.Close()

	return nil

}

func articleDelete(id int) error {

	DBCONN, err := DBConfig.PInitDb()
	if err != nil {
		return err
	}

	defer DBCONN.Close()

	fmt.Println("# Deleting values")
	stmt, err := DBCONN.Prepare("DELETE FROM articles WHERE id=$1")

	if err != nil {
		return err
	}

	res, err := stmt.Exec(id)

	fmt.Println(res)

	if err != nil {
		return err
	}

	stmt.Close()
	return nil

}
