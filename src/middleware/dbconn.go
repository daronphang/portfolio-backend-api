package middleware

import (
	"errors"
	"fmt"
	"portfolio_golang/src/config"
	"portfolio_golang/src/dblayer"
	e "portfolio_golang/src/errors"

	"github.com/gin-gonic/gin"
)

func ConnectDB(instance string) gin.HandlerFunc {
	return func(c *gin.Context) {
		db, ok := dblayer.DBInstance(instance)
		if !ok {
			err := e.NewError(
				e.ErrFailedDBConn,
				fmt.Sprintf("%s DB instance does not exist", instance),
				c.HandlerName(),
			)
			c.Error(err).SetMeta(err.Meta())
			c.AbortWithStatusJSON(e.ErrorResponse(err))
		}

		var err error

		if db.DB != nil {
			err = db.DB.Ping()
		}

		if err != nil || db.DB == nil {
			// if cannot ping DB or variable is nil, to reconnect
			switch instance {
			case "MSSQL":
				db, err = dblayer.InitMSSQL(config.AppCfg)
			case "MYSQL":
				db, err = dblayer.InitMYSQL(config.AppCfg)
			default:
				err = errors.New("db instance provided does not exist")
			}

			if err != nil {
				err := e.NewError(
					e.ErrFailedDBConn,
					err.Error(),
					c.HandlerName(),
				)
				c.Error(err).SetMeta(err.Meta())
				c.AbortWithStatusJSON(e.ErrorResponse(err))
			}
		}

		c.Set("db", db)
		c.Next()
	}
}
