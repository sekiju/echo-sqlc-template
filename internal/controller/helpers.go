package controller

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"net/http"
	"reflect"
	"slices"
	"strconv"
)

var ctx = context.Background()
var validate = validator.New(validator.WithRequiredStructEnabled())

func Body[T interface{}](c echo.Context) (T, error) {
	var body T
	if err := c.Bind(&body); err != nil {
		return body, echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	err := validate.Struct(body)
	if err != nil {
		var validationErrors validator.ValidationErrors
		errors.As(err, &validationErrors)

		return body, echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Validation failed: %s", validationErrors.Error()))
	}

	return body, nil
}

func HideKeys(i interface{}, keys ...string) interface{} {
	result := make(map[string]interface{})
	val := reflect.ValueOf(i)

	for i := 0; i < val.NumField(); i++ {
		field := val.Type().Field(i)

		exists := slices.Contains(keys, field.Name)
		if exists || field.PkgPath != "" {
			continue
		}

		key := field.Tag.Get("json")
		if key == "" {
			key = field.Name
		}

		result[key] = val.Field(i).Interface()
	}

	return result
}

func RouteID(c echo.Context) (int32, error) {
	param := c.Param("id")
	id, err := strconv.Atoi(param)
	if err != nil {
		return 0, echo.NewHTTPError(http.StatusNotFound, "Invalid route :id parameter")
	}

	return int32(id), nil
}
