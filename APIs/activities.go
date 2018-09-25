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

type ActivityInfo struct {
	ID           string    `json:"activity_id"`
	UserID       string    `json:"user_id"`
	Title        string    `json:"title"`
	Content      string    `json:"content"`
	CoverImg     string    `json:"cover_img"`
	LocationName string    `json:"location_name"`
	Status       string    `json:"status"`
	Latitude     float64   `json:"latitude"`
	Longitude    float64   `json:"longitude"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Tags         []string  `json:"tags"`
	City         string    `json:"city"`
	Popular      bool      `json:"popular"`
}

type ActivityReturn struct {
	ID string `json:"activity_id"`
	// UserID string `json:"user_id"`

	CreatorInfo struct {
		UserID   string `json:"user_id"`
		Avatar   string `json:"avatar"`
		Nickname string `json:"nickname"`
	} `json:"creator_info"`

	Title        string   `json:"title"`
	Content      string   `json:"content"`
	CoverImg     string   `json:"cover_img"`
	LocationName string   `json:"location_name"`
	Status       string   `json:"status"`
	Latitude     float64  `json:"latitude"`
	Longitude    float64  `json:"longitude"`
	CreatedAt    string   `json:"created_at"`
	UpdatedAt    string   `json:"updated_at"`
	DistanceKM   float64  `json:"distance_km"`
	Tags         []string `json:"tags"`
	City         string   `json:"city"`
	Popular      bool     `json:"popular"`
	IsCollected  bool     `json:"iscollected"`

	ParticipantList []struct {
		UserID   string `json:"user_id"`
		Avatar   string `json:"avatar"`
		Nickname string `json:"nickname"`
	} `json:"participants"`
}

type ActivityArray struct {
	ActivityList []struct {
		ID string `json:"activity_id"`
		// UserID string `json:"user_id"`

		CreatorInfo struct {
			UserID   string `json:"user_id"`
			Avatar   string `json:"avatar"`
			Nickname string `json:"nickname"`
		} `json:"creator_info"`

		Title        string   `json:"title"`
		Content      string   `json:"content"`
		CoverImg     string   `json:"cover_img"`
		LocationName string   `json:"location_name"`
		Status       string   `json:"status"`
		Latitude     float64  `json:"latitude"`
		Longitude    float64  `json:"longitude"`
		CreatedAt    string   `json:"created_at"`
		UpdatedAt    string   `json:"updated_at"`
		DistanceKM   float64  `json:"distance_km"`
		Tags         []string `json:"tags"`
		City         string   `json:"city"`
		Popular      bool     `json:"popular"`
		IsCollected  bool     `json:"iscollected"`

		ParticipantList []struct {
			UserID   string `json:"user_id"`
			Avatar   string `json:"avatar"`
			Nickname string `json:"nickname"`
		} `json:"participants"`
	} `json:"activity_list"`
}

type ActivityQuery struct {
	ActivityID string   `json:"activity_id"`
	UserID     string   `json:"user_id"`
	Title      string   `json:"title"`
	Latitude   float64  `json:"latitude"`
	Longitude  float64  `json:"longitude"`
	Tags       []string `json:"tags"`
	City       string   `json:"city"`
	Popular    bool     `json:"popular"`
	Mode       string   `json:"mode"`
}

type ActivityId struct {
	ID string `json:"activity_id"`
}

type pList struct {
	UserID   string `json:"user_id"`
	Avatar   string `json:"avatar"`
	Nickname string `json:"nickname"`
}

//this query is needed to check if it's collected
//select id, title, cover_img, (SELECT CAST(CASE WHEN ( SELECT COUNT(*) FROM article_collections WHERE user_id='dd247d2a-c607-4ac1-9528-0cc51d15843e' AND article_id=articles.id ) = 1 THEN true ELSE false END AS boolean) AS isCollected) from articles;
// (SELECT CAST(CASE WHEN ( SELECT COUNT(*) FROM activity_collections WHERE user_id=$1 AND article_id=articles.id ) = 1 THEN true ELSE false END AS boolean) AS isCollected)

func PostActivity(c echo.Context) (err error) {

	var activity ActivityInfo

	if err := c.Bind(&activity); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if activity.Longitude == 0 || activity.Latitude == 0 || activity.City == "" || activity.LocationName == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "latitude, longitude, city and location_name are required",
		})
	}

	if activity.UserID == "" || activity.Title == "" || activity.Content == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "user_id and title are required, content cannot be null",
		})
	}

	if activity.Popular != true && activity.Popular != false {
		activity.Popular = false
	}

	activity.CreatedAt = time.Now()
	activity.UpdatedAt = time.Now()

	err = activityModel(activity, 1)

	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Activity Post Successfully",
	})

}

func UpdateActivity(c echo.Context) (err error) {

	var activity ActivityInfo

	if err := c.Bind(&activity); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if activity.ID == "" || activity.Content == "" || activity.Title == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "activity_id is required, title and content cannot be null",
		})
	}

	activity.UpdatedAt = time.Now()
	err = activityModel(activity, 2)

	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Activity Updated Successfully",
	})

}

func DeleteActivity(c echo.Context) (err error) {

	var id ActivityId
	if err := c.Bind(&id); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if id.ID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "activity_id cannot be null",
		})
	}

	err = activityDelete(id.ID)

	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Activity Deleted",
	})

}

func GetActivity(c echo.Context) (err error) {

	var aQuery ActivityQuery
	var aList ActivityArray
	if err := c.Bind(&aQuery); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if aQuery.UserID == "" && aQuery.ActivityID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "user_id and/or activity_id is required",
		})
	}

	err = activityGet(aQuery, &aList)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, aList)
}

func GetActivityByDistance(c echo.Context) (err error) {

	var aQuery ActivityQuery
	var aList ActivityArray
	if err := c.Bind(&aQuery); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if aQuery.Latitude == 0 || aQuery.Longitude == 0 || aQuery.City == "" || aQuery.UserID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "latitude, longitude, city and user_id are required",
		})
	}

	err = activityReadDistance(aQuery, &aList)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, aList)
}

func JoinActivity(c echo.Context) (err error) {

	var aQuery ActivityQuery

	if err := c.Bind(&aQuery); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if aQuery.UserID == "" || aQuery.ActivityID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "activity_id or user_id is required",
		})
	}

	err = activityJoin(aQuery)

	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "successfully joined activity",
	})

}

func GetJoinedActivity(c echo.Context) (err error) {

	var aQuery ActivityQuery
	var aList ActivityArray
	if err := c.Bind(&aQuery); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if aQuery.UserID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "user_id is required",
		})
	}

	err = activityJoinGet(aQuery, &aList)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, aList)

}

func CollectActivity(c echo.Context) (err error) {

	var aQuery ActivityQuery

	if err := c.Bind(&aQuery); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if aQuery.UserID == "" || aQuery.ActivityID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "activity_id and user_id is required",
		})
	}

	err = activityCollect(aQuery)

	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "successfully collected activity",
	})

}

func GetCollectedActivity(c echo.Context) (err error) {

	var aQuery ActivityQuery
	var aList ActivityArray
	if err := c.Bind(&aQuery); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if aQuery.UserID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "user_id is required",
		})
	}

	err = activityCollectGet(aQuery, &aList)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, aList)

}

func RecommandActivity(c echo.Context) (err error) {

	var aList ActivityArray
	var aQuery ActivityQuery

	if err := c.Bind(&aQuery); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if aQuery.City == "" || aQuery.UserID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "city and user_id are required",
		})
	}

	err = activityGetRecommand(aQuery, &aList)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, aList)
}

func PopularActivity(c echo.Context) (err error) {

	var aList ActivityArray
	var aQuery ActivityQuery

	if err := c.Bind(&aQuery); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if aQuery.City == "" || aQuery.UserID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "city and user_id are required",
		})
	}

	err = activityGetPopular(aQuery, &aList)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, aList)
}

func DeleteActivityCollect(c echo.Context) (err error) {

	var aQuery ActivityQuery
	if err := c.Bind(&aQuery); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": err.Error(),
		})
	}

	if aQuery.ActivityID == "" || aQuery.UserID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "activity_id and user_id are required",
		})
	}

	err = activityCollectDelete(aQuery)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "collection deleted",
	})
}

func activityCollectDelete(aQuery ActivityQuery) error {

	DBCONN, err := DBConfig.PInitDb()
	if err != nil {
		return err
	}
	defer DBCONN.Close()

	stmt, err := DBCONN.Prepare("DELETE FROM activity_collections WHERE activity_id=$1 AND user_id=$2")
	if err != nil {
		return err
	}

	res, err := stmt.Exec(aQuery.ActivityID, aQuery.UserID)
	if err != nil {
		return err
	}

	fmt.Println(res)
	stmt.Close()

	return nil
}

func activityCollectGet(aQuery ActivityQuery, aList *ActivityArray) error {

	DBCONN, err := DBConfig.PInitDb()
	if err != nil {
		return err
	}

	defer DBCONN.Close()

	var stmt *sql.Stmt
	var rows *sql.Rows

	stmt, err = DBCONN.Prepare("SELECT id, creator_id, title, content, cover_img, location_name, status, latitude, longitude, created_at, updated_at, tags, city, popular FROM activities WHERE id IN (SELECT activity_id FROM activity_collections WHERE user_id=$1)")
	if err != nil {
		return err
	}

	rows, err = stmt.Query(aQuery.UserID)
	if err != nil {
		return err
	}

	for rows.Next() {

		var ca, ua time.Time
		var activityReturn ActivityReturn
		err = rows.Scan(&activityReturn.ID,
			&activityReturn.CreatorInfo.UserID,
			&activityReturn.Title,
			&activityReturn.Content,
			&activityReturn.CoverImg,
			&activityReturn.LocationName,
			&activityReturn.Status,
			&activityReturn.Latitude,
			&activityReturn.Longitude,
			&ca, &ua,
			pq.Array(&activityReturn.Tags),
			&activityReturn.City,
			&activityReturn.Popular)

		if err != nil {
			return err
		}

		activityReturn.CreatedAt = ca.String()
		activityReturn.UpdatedAt = ua.String()

		err = getCreatorInfo(&activityReturn)
		if err != nil {
			return err
		}

		err = getParticipantsList(&activityReturn)
		if err != nil {
			return err
		}

		aList.ActivityList = append(aList.ActivityList, activityReturn)
	}

	rows.Close()
	stmt.Close()

	return nil
}

func activityCollect(aQuery ActivityQuery) error {

	DBCONN, err := DBConfig.PInitDb()
	if err != nil {
		return err
	}

	defer DBCONN.Close()

	stmt, err := DBCONN.Prepare("INSERT INTO activity_collections (activity_id, user_id) VALUES ($1, $2)")
	if err != nil {
		return err
	}

	res, err := stmt.Exec(aQuery.ActivityID, aQuery.UserID)
	if err != nil {
		return err
	}

	stmt.Close()
	fmt.Println(res)

	return nil
}

func activityJoinGet(aQuery ActivityQuery, aList *ActivityArray) error {

	DBCONN, err := DBConfig.PInitDb()
	if err != nil {
		return err
	}

	defer DBCONN.Close()

	var stmt *sql.Stmt
	var rows *sql.Rows

	stmt, err = DBCONN.Prepare("SELECT id, creator_id, title, content, cover_img, location_name, status, latitude, " +
		"longitude, created_at, updated_at, tags, city, popular, " +
		"(SELECT CAST(CASE WHEN ( SELECT COUNT(*) FROM activity_collections WHERE user_id=$1 AND activity_id=activities.id )=1 THEN true ELSE false END AS boolean) AS isCollected) " +
		"FROM activities WHERE id IN (SELECT activity_id FROM activity_members WHERE user_id=$1)")

	if err != nil {
		return err
	}

	rows, err = stmt.Query(aQuery.UserID)
	if err != nil {
		return err
	}

	for rows.Next() {

		var ca, ua time.Time
		var activityReturn ActivityReturn
		err = rows.Scan(&activityReturn.ID,
			&activityReturn.CreatorInfo.UserID,
			&activityReturn.Title,
			&activityReturn.Content,
			&activityReturn.CoverImg,
			&activityReturn.LocationName,
			&activityReturn.Status,
			&activityReturn.Latitude,
			&activityReturn.Longitude,
			&ca, &ua,
			pq.Array(&activityReturn.Tags),
			&activityReturn.City,
			&activityReturn.Popular,
			&activityReturn.IsCollected)

		if err != nil {
			return err
		}

		activityReturn.CreatedAt = ca.String()
		activityReturn.UpdatedAt = ua.String()

		err = getCreatorInfo(&activityReturn)
		if err != nil {
			return err
		}

		err = getParticipantsList(&activityReturn)
		if err != nil {
			return err
		}

		aList.ActivityList = append(aList.ActivityList, activityReturn)
	}

	rows.Close()
	stmt.Close()

	return nil
}

func activityJoin(aQuery ActivityQuery) error {

	DBCONN, err := DBConfig.PInitDb()
	if err != nil {
		return err
	}

	defer DBCONN.Close()

	var isKicked int
	stmt, err := DBCONN.Prepare("SELECT COUNT(*) FROM activity_kick_histories WHERE kick_user_id=$1")
	if err != nil {
		return err
	}

	err = stmt.QueryRow(aQuery.UserID).Scan(&isKicked)
	if err != nil {
		return err
	}

	if isKicked != 0 {
		return errors.New("cannot join (kicked out)")
	}

	stmt, err = DBCONN.Prepare("INSERT INTO activity_members (activity_id, user_id) VALUES ($1, $2)")
	if err != nil {
		return err
	}

	res, err := stmt.Exec(aQuery.ActivityID, aQuery.UserID)
	if err != nil {
		return err
	}

	stmt.Close()
	fmt.Println(res)

	return nil
}

func activityReadDistance(aQuery ActivityQuery, aList *ActivityArray) error {

	DBCONN, err := DBConfig.PInitDb()
	if err != nil {
		return err
	}

	defer DBCONN.Close()

	var stmt *sql.Stmt
	var rows *sql.Rows

	stmt, err = DBCONN.Prepare("SELECT id, creator_id, title, cover_img, content, location_name, status, created_at, updated_at, tags, city, popular, " +
		"ST_Distance(ST_MakePoint($1, $2), location, false)/1000 as distance_km, " +
		"(SELECT CAST(CASE WHEN ( SELECT COUNT(*) FROM activity_collections WHERE user_id=$3 AND activity_id=activities.id )=1 THEN true ELSE false END AS boolean) AS isCollected) " +
		"FROM activities WHERE city=$4 ORDER BY distance_km ASC")
	if err != nil {
		return err
	}

	rows, err = stmt.Query(aQuery.Latitude, aQuery.Longitude, aQuery.UserID, aQuery.City)
	if err != nil {
		return err
	}

	for rows.Next() {

		var ca, ua time.Time
		var activityReturn ActivityReturn
		err = rows.Scan(&activityReturn.ID,
			&activityReturn.CreatorInfo.UserID,
			&activityReturn.Title,
			&activityReturn.CoverImg,
			&activityReturn.Content,
			&activityReturn.LocationName,
			&activityReturn.Status,
			&ca, &ua,
			pq.Array(&activityReturn.Tags),
			&activityReturn.City,
			&activityReturn.Popular,
			&activityReturn.DistanceKM,
			&activityReturn.IsCollected)

		if err != nil {
			return err
		}

		activityReturn.CreatedAt = ca.String()
		activityReturn.UpdatedAt = ua.String()

		err = getCreatorInfo(&activityReturn)
		if err != nil {
			return err
		}

		err = getParticipantsList(&activityReturn)
		if err != nil {
			return err
		}

		aList.ActivityList = append(aList.ActivityList, activityReturn)
	}

	stmt.Close()
	rows.Close()

	return nil
}

func activityGet(aQuery ActivityQuery, aList *ActivityArray) error {

	DBCONN, err := DBConfig.PInitDb()
	if err != nil {
		return err
	}

	defer DBCONN.Close()

	var stmt *sql.Stmt
	var rows *sql.Rows

	if aQuery.ActivityID != "" {
		stmt, err = DBCONN.Prepare("SELECT id, creator_id, title, content, cover_img, location_name, status, " +
			"latitude, longitude, created_at, updated_at, tags, city, popular, " +
			"(SELECT CAST(CASE WHEN ( SELECT COUNT(*) FROM activity_collections WHERE user_id=$1 AND activity_id=activities.id )=1 THEN true ELSE false END AS boolean) AS isCollected) " +
			" FROM activities WHERE id=$2")
		if err != nil {
			return err
		}

		rows, err = stmt.Query(aQuery.UserID, aQuery.ActivityID)
		if err != nil {
			return err
		}

	} else if aQuery.Mode == "creator" {
		stmt, err = DBCONN.Prepare("SELECT id, creator_id, title, content, cover_img, location_name, status, " +
			"latitude, longitude, created_at, updated_at, tags, city, popular, " +
			"(SELECT CAST(CASE WHEN ( SELECT COUNT(*) FROM activity_collections WHERE user_id=$1 AND activity_id=activities.id )=1 THEN true ELSE false END AS boolean) AS isCollected) " +
			" FROM activities WHERE creator_id=$1")
		if err != nil {
			return err
		}

		rows, err = stmt.Query(aQuery.UserID)
		if err != nil {
			return err
		}

	} else {
		stmt, err = DBCONN.Prepare("SELECT id, creator_id, title, content, cover_img, location_name, status, " +
			"latitude, longitude, created_at, updated_at, tags, city, popular, " +
			"(SELECT CAST(CASE WHEN ( SELECT COUNT(*) FROM activity_collections WHERE user_id=$1 AND activity_id=activities.id )=1 THEN true ELSE false END AS boolean) AS isCollected) " +
			" FROM activities")
		if err != nil {
			return err
		}

		rows, err = stmt.Query(aQuery.UserID)
		if err != nil {
			return err
		}
	}

	for rows.Next() {

		var ca, ua time.Time
		var activityReturn ActivityReturn
		err = rows.Scan(&activityReturn.ID,
			&activityReturn.CreatorInfo.UserID,
			&activityReturn.Title,
			&activityReturn.Content,
			&activityReturn.CoverImg,
			&activityReturn.LocationName,
			&activityReturn.Status,
			&activityReturn.Latitude,
			&activityReturn.Longitude,
			&ca, &ua,
			pq.Array(&activityReturn.Tags),
			&activityReturn.City,
			&activityReturn.Popular,
			&activityReturn.IsCollected)

		if err != nil {
			return err
		}

		activityReturn.CreatedAt = ca.String()
		activityReturn.UpdatedAt = ua.String()

		err = getCreatorInfo(&activityReturn)
		if err != nil {
			return err
		}

		err = getParticipantsList(&activityReturn)
		if err != nil {
			return err
		}

		aList.ActivityList = append(aList.ActivityList, activityReturn)
	}

	rows.Close()
	stmt.Close()

	return nil
}

func activityGetPopular(aQuery ActivityQuery, aList *ActivityArray) error {

	DBCONN, err := DBConfig.PInitDb()
	if err != nil {
		return err
	}

	defer DBCONN.Close()

	var stmt *sql.Stmt
	var rows *sql.Rows

	stmt, err = DBCONN.Prepare("SELECT id, creator_id, title, content, cover_img, location_name, status, " +
		"latitude, longitude, created_at, updated_at, tags, city, popular, " +
		"(SELECT CAST(CASE WHEN ( SELECT COUNT(*) FROM activity_collections WHERE user_id=$1 AND activity_id=activities.id )=1 THEN true ELSE false END AS boolean) AS isCollected) " +
		" FROM activities WHERE city=$2 AND popular=true")
	if err != nil {
		return err
	}

	rows, err = stmt.Query(aQuery.UserID, aQuery.City)
	if err != nil {
		return err
	}

	for rows.Next() {

		var ca, ua time.Time
		var activityReturn ActivityReturn
		err = rows.Scan(&activityReturn.ID,
			&activityReturn.CreatorInfo.UserID,
			&activityReturn.Title,
			&activityReturn.Content,
			&activityReturn.CoverImg,
			&activityReturn.LocationName,
			&activityReturn.Status,
			&activityReturn.Latitude,
			&activityReturn.Longitude,
			&ca, &ua,
			pq.Array(&activityReturn.Tags),
			&activityReturn.City,
			&activityReturn.Popular,
			&activityReturn.IsCollected)

		if err != nil {
			return err
		}

		activityReturn.CreatedAt = ca.String()
		activityReturn.UpdatedAt = ua.String()

		err = getCreatorInfo(&activityReturn)
		if err != nil {
			return err
		}

		err = getParticipantsList(&activityReturn)
		if err != nil {
			return err
		}

		aList.ActivityList = append(aList.ActivityList, activityReturn)
	}

	rows.Close()
	stmt.Close()

	return nil
}

func activityGetRecommand(aQuery ActivityQuery, aList *ActivityArray) error {

	DBCONN, err := DBConfig.PInitDb()
	if err != nil {
		return err
	}

	defer DBCONN.Close()

	var stmt *sql.Stmt
	var rows *sql.Rows

	stmt, err = DBCONN.Prepare("SELECT id, creator_id, title, content, cover_img, location_name, status, " +
		"latitude, longitude, created_at, updated_at, tags, city, popular, " +
		"(SELECT CAST(CASE WHEN ( SELECT COUNT(*) FROM activity_collections WHERE user_id=$1 AND activity_id=activities.id )=1 THEN true ELSE false END AS boolean) AS isCollected) " +
		" FROM activities WHERE city=$2 AND popular=false")
	if err != nil {
		return err
	}

	rows, err = stmt.Query(aQuery.UserID, aQuery.City)
	if err != nil {
		return err
	}

	for rows.Next() {

		var ca, ua time.Time
		var activityReturn ActivityReturn
		err = rows.Scan(&activityReturn.ID,
			&activityReturn.CreatorInfo.UserID,
			&activityReturn.Title,
			&activityReturn.Content,
			&activityReturn.CoverImg,
			&activityReturn.LocationName,
			&activityReturn.Status,
			&activityReturn.Latitude,
			&activityReturn.Longitude,
			&ca, &ua,
			pq.Array(&activityReturn.Tags),
			&activityReturn.City,
			&activityReturn.Popular,
			&activityReturn.IsCollected)

		if err != nil {
			return err
		}

		activityReturn.CreatedAt = ca.String()
		activityReturn.UpdatedAt = ua.String()

		err = getCreatorInfo(&activityReturn)
		if err != nil {
			return err
		}

		err = getParticipantsList(&activityReturn)
		if err != nil {
			return err
		}

		aList.ActivityList = append(aList.ActivityList, activityReturn)
	}

	rows.Close()
	stmt.Close()

	return nil
}

func activityModel(info ActivityInfo, mode int) error {
	// if mode == 1 -> insert, if mode == 2 -> Update, else error

	var err error
	if mode == 1 {
		err = activityInsert(info)
	} else if mode == 2 {
		err = activityUpdate(info)
	} else {
		err = errors.New("mode must be either 1 or 2")
	}

	if err != nil {
		return err
	}

	return nil
}

func activityUpdate(info ActivityInfo) error {

	DBCONN, err := DBConfig.PInitDb()
	if err != nil {
		return err
	}

	defer DBCONN.Close()

	fmt.Println("# Updating values")

	stmt, err := DBCONN.Prepare("UPDATE activities SET title=$1, content=$2, cover_img=$3, updated_at=$4 WHERE id=$5")

	if err != nil {
		return err
	}

	res, err := stmt.Exec(info.Title,
		info.Content,
		info.CoverImg,
		info.UpdatedAt,
		info.ID)

	fmt.Println(res)

	if err != nil {
		return err
	}

	stmt.Close()
	return nil

}

func activityInsert(info ActivityInfo) error {

	DBCONN, err := DBConfig.PInitDb()
	if err != nil {
		return err
	}

	defer DBCONN.Close()
	fmt.Println("# Calculating geography value")

	stmt, err := DBCONN.Prepare("SELECT ST_MakePoint($1, $2) as geo")
	if err != nil {
		return err
	}

	var geoinfo string

	err = stmt.QueryRow(info.Latitude, info.Longitude).Scan(&geoinfo)
	if err != nil {
		return err
	}

	fmt.Println("# Inserting values")

	stmt, err = DBCONN.Prepare("INSERT INTO activities (creator_id, title, content, cover_img, status, location_name, latitude, longitude, created_at, updated_at, location, tags, city, popular) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14) RETURNING id")
	if err != nil {
		return err
	}
	var actID string

	err = stmt.QueryRow(info.UserID,
		info.Title,
		info.Content,
		info.CoverImg,
		info.Status,
		info.LocationName,
		info.Latitude,
		info.Longitude,
		info.CreatedAt,
		info.UpdatedAt,
		geoinfo,
		pq.Array(info.Tags),
		info.City,
		info.Popular).Scan(&actID)

	if err != nil {
		return err
	}

	stmt, err = DBCONN.Prepare("INSERT INTO activity_members (activity_id, user_id) VALUES ($1, $2)")
	if err != nil {
		return err
	}

	res, err := stmt.Exec(actID, info.UserID)
	if err != nil {
		return err
	}
	fmt.Println(res)

	stmt.Close()

	return nil

}

func activityDelete(id string) error {

	DBCONN, err := DBConfig.PInitDb()
	if err != nil {
		return err
	}

	defer DBCONN.Close()

	fmt.Println("# Deleting values")
	stmt, err := DBCONN.Prepare("DELETE FROM activities WHERE id=$1")

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

func getCreatorInfo(a *ActivityReturn) error {

	DBCONN, err := DBConfig.PInitDb()
	if err != nil {
		return err
	}
	defer DBCONN.Close()

	stmt, err := DBCONN.Prepare("SELECT avatar, profile->>'nickname' AS nickname FROM users WHERE id=$1")
	if err != nil {
		return nil
	}

	err = stmt.QueryRow(a.CreatorInfo.UserID).Scan(&a.CreatorInfo.Avatar, &a.CreatorInfo.Nickname)
	if err != nil {
		return err
	}

	return nil
}

func getParticipantsList(a *ActivityReturn) error {

	DBCONN, err := DBConfig.PInitDb()
	if err != nil {
		return err
	}
	defer DBCONN.Close()

	stmt, err := DBCONN.Prepare("SELECT id, avatar, profile->>'nickname' AS nickname FROM users WHERE id IN (SELECT user_id FROM activity_members WHERE activity_id=$1) ")
	if err != nil {
		return err
	}

	rows, err := stmt.Query(a.ID)
	if err != nil {
		return err
	}

	for rows.Next() {

		var p pList
		err = rows.Scan(&p.UserID, &p.Avatar, &p.Nickname)
		if err != nil {
			return err
		}

		a.ParticipantList = append(a.ParticipantList, p)

	}

	return nil
}
