package APIs

import (
	"database/sql"
	"net/http"
	"github.com/csnand/APIProject/APIs/DBConfig"
	"github.com/labstack/echo"
	_ "github.com/lib/pq"
)

type TagsArray struct {
	Tags []string `json:"tags"`
}

func GetTags(c echo.Context) (err error) {

	var tList TagsArray
	err = tagsGet(&tList)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, tList)

}

func tagsGet(tList *TagsArray) error {

	DBCONN, err := DBConfig.PInitDb()
	if err != nil {
		return err
	}

	defer DBCONN.Close()

	var stmt *sql.Stmt
	var rows *sql.Rows

	stmt, err = DBCONN.Prepare("SELECT name FROM tags")
	if err != nil {
		return err
	}

	rows, err = stmt.Query()
	if err != nil {
		return err
	}

	for rows.Next() {

		var tag string
		err = rows.Scan(&tag)

		if err != nil {
			return err
		}

		tList.Tags = append(tList.Tags, tag)
	}

	rows.Close()

	return nil
}
