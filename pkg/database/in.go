package database

import (
	"demo/internal/common"
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/jmoiron/sqlx"
)

func in(query string, args []interface{}) (string, []interface{}, error) {
	r, err := regexp.Compile(`(?i)in\s*\(\s*(\$|\?)`)
	if err != nil {
		return "", nil, common.Wrap(err, "regexp.Compile")
	}

	match := r.FindStringSubmatch(query)

	if len(match) < 2 {
		return query, args, nil
	}

	switch match[1] {
	case "?":
		return questionMark(query, args)
	case "$":
		return sigil(query, args)
	}

	return query, args, nil
}

func questionMark(query string, args []interface{}) (string, []interface{}, error) {
	query, args, err := sqlx.In(query, args...)
	if err != nil {
		return "", nil, common.Wrap(err, "sqlx.In")
	}

	return query, args, nil
}

func sigil(query string, args []interface{}) (string, []interface{}, error) {
	getSigil := func(index int) string {
		return fmt.Sprintf("$%d", index)
	}

	for a, arg := range args {
		vof := reflect.ValueOf(arg)
		if vof.Kind() == reflect.Slice {
			for i := 0; i < vof.Len(); i++ {
				sliceVal := vof.Index(i).Interface()
				if i == 0 {
					args[a] = sliceVal
				} else {
					query = strings.Replace(
						query,
						getSigil(a+1),
						fmt.Sprintf("%s, %s", getSigil(a+1), getSigil(len(args)+1)),
						-1,
					)
					args = append(args, sliceVal)
				}
			}
		}
	}

	return query, args, nil
}
