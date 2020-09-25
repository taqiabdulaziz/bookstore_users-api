package mysql_utils

import (
	"github.com/go-sql-driver/mysql"
	"github.com/taqiabdulaziz/bookstore_users-api/utils"
	"strings"
)

const (
	errorNoRows = "no rows in result set"
)

func ParseError(err error) *utils.RestErr {
	sqlErr, ok := err.(*mysql.MySQLError)
	if !ok {
		if strings.Contains(err.Error(), errorNoRows) {
			return utils.NewBadRequestError("no record matching given id")
		}
		return utils.NewInternalServerError("error parsing database response")
	}

	switch sqlErr.Number {
	case 1062:
		return utils.NewBadRequestError("invalid data")
	}

	return utils.NewInternalServerError("error processing request")
}
