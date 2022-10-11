package controller

import (
	"database/sql"
	"net/http"
	"portfolio_golang/src/dblayer"
	e "portfolio_golang/src/errors"
	"portfolio_golang/src/internal"

	"github.com/gin-gonic/gin"
)

type Testing struct {
	ID              sql.NullString `json:"id"`
	NAME            sql.NullString `json:"name"`
	ADDED_TIMESTAMP sql.NullString `json:"added_timestamp"`
}

func TestDB(c *gin.Context) {
	data, ok := internal.StoredContextValues(c, "db")
	if !ok {
		return
	}
	db := data["db"].(*dblayer.DBLayer)

	rows, err := db.Query("testing")
	if err != nil {
		err := e.NewError(e.ErrFailedSQLQuery, err.Error(), c.HandlerName())
		c.Error(err).SetMeta(err.Meta())
		c.AbortWithStatusJSON(e.ErrorResponse(err))
		return
	}

	defer rows.Close()

	var testingList []Testing

	for rows.Next() {
		var testing Testing
		err = rows.Scan(
			&testing.ID,
			&testing.NAME,
			&testing.ADDED_TIMESTAMP,
		)

		if err != nil {
			err := e.NewError(e.ErrRowScan, err.Error(), c.HandlerName())
			c.Error(err).SetMeta(err.Meta())
			c.AbortWithStatusJSON(e.ErrorResponse(err))
			return
		}
		testingList = append(testingList, testing)
	}

	err = rows.Err()
	if err != nil {
		err := e.NewError(e.ErrRowScan, err.Error(), c.HandlerName())
		c.Error(err).SetMeta(err.Meta())
		c.AbortWithStatusJSON(e.ErrorResponse(err))
		return
	}

	res := JSONResponse{
		Response: testingList,
		Message:  "DB query success",
	}
	c.JSON(http.StatusOK, res)
}
