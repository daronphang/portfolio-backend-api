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
	UUID              	sql.NullString `json:"uuid"`
	CONTACT_NAME        sql.NullString `json:"contact_name"`
	EMAIL 				sql.NullString `json:"email"`
	SUBJECT 			sql.NullString `json:"subject"`
	MESSAGE 			sql.NullString `json:"message"`
	CREATED_DATETIME 	sql.NullString `json:"created_datetime"`
}

type ContactPayload struct {
	CONTACT_NAME	string	`json:"contact_name" binding:"required,max=255"`
	EMAIL			string	`json:"email" binding:"required,max=255"`
	SUBJECT			string	`json:"subject" binding:"max=255"`
	MESSAGE			string	`json:"message" binding:"max=255"`
}

var sqlReadStr = `
	SELECT
	contact_name,
	email,
	subject,
	message,
	created_datetime
	FROM
	contacts
	%s
	ORDER BY
	created_datetime DESC
`

var sqlCreateStr = `
	INSERT INTO contacts (
		uuid,
		contact_name,
		email,
		subject,
		message
	)
	VALUES (
		?,?,?,?,?
	)
`

var sqlUpdateStr = `
	UPDATE 
	contacts
	SET message = '%s'
	WHERE contact_name = '%s'
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
			&row.CONTACT_NAME,
			&row.EMAIL,
			&row.SUBJECT,
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

	_, err = stmt.Exec(uuid.New(), payload.CONTACT_NAME, payload.EMAIL, payload.SUBJECT, payload.MESSAGE)
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
		fmt.Sprintf("Portfolio Contact: %s", payload.SUBJECT),
		fmt.Sprintf(
			"Name: %s\nEmail: %s\nMessage: %s",
			payload.CONTACT_NAME,
			payload.EMAIL,
			payload.MESSAGE,
		),
	)
}