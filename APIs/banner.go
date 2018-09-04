package APIs

import (
	"fmt"
	"net/http"

	"./DBConfig"
	"github.com/labstack/echo"
)

type BannerItem struct {
	ActivityId  string `json:"activity_id"`
	CreatorId   string `json:"creator_id"`
	Title       string `json:"title"`
	CoverImg    string `json:"cover_id"`
	CategoryId  string `json:"category_id"`
	IsCollected string `json:"iscollected"`
}

type BannerArray struct {
	BannerList []struct {
		ActivityId  string `json:"activity_id"`
		CreatorId   string `json:"creator_id"`
		Title       string `json:"title"`
		CoverImg    string `json:"cover_id"`
		CategoryId  string `json:"category_id"`
		IsCollected string `json:"iscollected"`
	} `json:"banner_list"`
}

type BannerQuery struct {
	UserID string `json:"user_id"`
}

func Banner(c echo.Context) (err error) {

	var bQuery BannerQuery
	if err := c.Bind(&bQuery); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": err.Error(),
		})
	}

	if bQuery.UserID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "user_id is required",
		})
	}

	DBCONN, err := DBConfig.PInitDb()
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	defer DBCONN.Close()

	stmt, err := DBCONN.Prepare("SELECT id, title, cover_img, category_id, creator_id, " +
		"(SELECT CAST(CASE WHEN ( SELECT COUNT(*) FROM activity_collections WHERE user_id=$1 AND activity_id=activities.id )=1 THEN true ELSE false END AS boolean) AS isCollected) " +
		" FROM activities WHERE weight>=1000000 AND status!='end'")
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	rows, err := stmt.Query(bQuery.UserID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	var banner BannerArray

	for rows.Next() {

		var bItem BannerItem

		err = rows.Scan(&bItem.ActivityId, &bItem.Title,
			&bItem.CoverImg, &bItem.CategoryId,
			&bItem.CreatorId,
			&bItem.IsCollected)

		if err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}

		banner.BannerList = append(banner.BannerList, bItem)
	}

	rows.Close()

	fmt.Println(banner)

	return c.JSON(http.StatusOK, banner)
}
