package utils

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

func ToSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

func ConvertIntSet(set *schema.Set) []int {
	s := make([]int, 0, set.Len())
	for _, v := range set.List() {
		s = append(s, v.(int))
	}
	sort.Ints(s)

	return s
}

func Exists(i int, array []int) bool {

	for _, v := range array {
		if i == v {
			return true
		}
	}
	return false
}

func ConvertArrayInterfaceToArrayString(arrayInt []interface{}) []string {
	arrayStr := make([]string, len(arrayInt))
	for i, v := range arrayInt {
		arrayStr[i] = fmt.Sprint(v)
	}
	return arrayStr
}

func ConvertArrayInterfaceToArrayInt(arrayInterface []interface{}) []int {
	arrayInt := make([]int, len(arrayInterface))
	for i, v := range arrayInterface {
		num, _ := strconv.Atoi(fmt.Sprint(v))
		arrayInt[i] = num
	}
	return arrayInt
}

func ConvertArrayInterfaceToArrayFloat(arrayInterface []interface{}) []float64 {
	arrayInt := make([]float64, len(arrayInterface))
	for i, v := range arrayInterface {
		num, _ := strconv.ParseFloat(fmt.Sprint(v), 64)
		arrayInt[i] = num
	}
	return arrayInt
}

func IsISOTime(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be string", k))
		return warnings, errors
	}

	if _, err := time.Parse("2006-01-02T15:04:05Z0700", v); err != nil {
		errors = append(errors, fmt.Errorf("expected %q to be a valid ISO8601 date, got %q: %+v", k, i, err))
	}

	return warnings, errors
}
