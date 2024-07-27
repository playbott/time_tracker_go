package middleware

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"reflect"
	"strconv"
	"strings"
	"unicode/utf8"
)

const strMinMaxName = "str_min_max"

func ValidatorMiddleware(v *validator.Validate) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Locals("validator", v)
		err := v.RegisterValidation(strMinMaxName, strMinMax)
		if err != nil {
			return err
		}
		return c.Next()
	}
}

// strMinMax custom validation function
func strMinMax(fl validator.FieldLevel) bool {
	var err error
	var err2 error
	field := fl.Field()
	if field.Kind() == reflect.String {
		if field.String() == "" {
			return true
		}

		params := strings.Fields(fl.Param())
		lengthMin := 0
		lengthMax := 0

		if len(params) == 0 && len(params) > 2 {
			fmt.Println("parameters:", fl.Param())
			panic(strMinMaxName + " must have 2 parameters in field " + fl.FieldName())
		}
		if len(params) > 0 {
			lengthMin, err = strconv.Atoi(params[0])
			lengthMax = lengthMin
		}
		if len(params) == 2 {
			lengthMax, err2 = strconv.Atoi(params[1])
		}
		if err != nil || err2 != nil {
			fmt.Println("parameters:", fl.Param())
			panic(strMinMaxName + " must have 2 integer parameters in field " + fl.FieldName())
		}

		fieldLength := utf8.RuneCountInString(field.String())
		if fieldLength < lengthMin || fieldLength > lengthMax {
			return false
		}
	}
	return true
}
