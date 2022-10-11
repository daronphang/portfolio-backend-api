package controller

import (
	"database/sql"
	"fmt"
	"net/http"
	"portfolio_golang/src/dblayer"
	e "portfolio_golang/src/errors"
	"portfolio_golang/src/internal"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Contacts struct {
	UID              	sql.NullString `json:"uid"`
	USERNAME            sql.NullString `json:"username"`
	EMAIL 				sql.NullString `json:"email"`
	COMPANY 			sql.NullString `json:"company"`
	MESSAGE 			sql.NullString `json:"message"`
	CREATED_DATETIME 	sql.NullString `json:"created_datetime"`
}

type ContactPayload struct {
	USERNAME	string	`json:"username" binding:"required,max=250"`
	EMAIL		string	`json:"email" binding:"required,max=250"`
	COMPANY		string	`json:"company" binding:"max=250"`
	MESSAGE		string	`json:"message" binding:"max=250"`
}

var sqlReadStr = `
	SELECT
	*
	FROM
	contact_me
	%s
	ORDER BY
	created_datetime DESC
`

var sqlCreateStr = `
	INSERT INTO contact_me (
		uid,
		username,
		email,
		company,
		message
	)
	VALUES (
		?,?,?,?,?
	)
`

var sqlUpdateStr = `
	UPDATE 
	contact_me
	SET message = '%s'
	WHERE username = '%s'
`

func ReadContacts(c *gin.Context) {
	data, ok := internal.StoredContextValues(c, "db")
	if !ok {
		return
	}
	db := data["db"].(*dblayer.DBLayer)

	// handler is used to get all contacts or specific contact
	// only email is allowed as query clause
	email := c.Param("email")

	if email != "" {
		email = fmt.Sprintf("WHERE email = '%s'", email)
	}

	rows, err := db.Query(fmt.Sprintf(sqlReadStr, email))
	if err != nil {
		err := e.NewError(e.ErrFailedSQLQuery, err.Error(), c.HandlerName())
		c.Error(err).SetMeta(err.Meta())
		c.AbortWithStatusJSON(e.ErrorResponse(err))
		return
	}

	defer rows.Close()

	var contacts []Contacts

	for rows.Next() {
		var row Contacts
		err = rows.Scan(
			&row.UID,
			&row.USERNAME,
			&row.EMAIL,
			&row.COMPANY,
			&row.MESSAGE,
			&row.CREATED_DATETIME,
		)

		if err != nil {
			err := e.NewError(e.ErrRowScan, err.Error(), c.HandlerName())
			c.Error(err).SetMeta(err.Meta())
			c.AbortWithStatusJSON(e.ErrorResponse(err))
			return
		}
		contacts = append(contacts, row)
	}

	err = rows.Err()
	if err != nil {
		err := e.NewError(e.ErrRowScan, err.Error(), c.HandlerName())
		c.Error(err).SetMeta(err.Meta())
		c.AbortWithStatusJSON(e.ErrorResponse(err))
		return
	}

	res := JSONResponse{
		Response: contacts,
		Message:  "query successful",
	}
	c.JSON(http.StatusOK, res)
}

func CreateContact(c *gin.Context) {
	data, ok := internal.StoredContextValues(c, "db", "payload")
	if !ok {
		return
	}
	db := data["db"].(*dblayer.DBLayer)
	payload := data["payload"].(*ContactPayload)

	stmt, err := db.Prepare(sqlCreateStr)
	if err != nil {
		err := e.NewError(e.ErrFailedSQLExec, err.Error(), c.HandlerName())
		c.Error(err).SetMeta(err.Meta())
		c.AbortWithStatusJSON(e.ErrorResponse(err))
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(uuid.New(), payload.USERNAME, payload.EMAIL, payload.COMPANY, payload.MESSAGE)
	if err != nil {
		err := e.NewError(e.ErrFailedSQLExec, err.Error(), c.HandlerName())
		c.Error(err).SetMeta(err.Meta())
		c.AbortWithStatusJSON(e.ErrorResponse(err))
		return
	}

	res := JSONResponse{
		Message:  "create contact successful",
	}
	c.JSON(http.StatusOK, res)

	go SendGmail(
		"daronphang@gmail.com",
		"daronphang@gmail.com",
		"Portfolio Contact", 
		fmt.Sprintf(
			"Name: %s\nEmail: %s\nCompany: %s\nMessage: %s",
			payload.USERNAME,
			payload.EMAIL,
			payload.COMPANY,
			payload.MESSAGE,
		),
	)
}