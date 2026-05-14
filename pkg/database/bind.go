package database

import (
	"demo/internal/common"
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

// Формат совместимости в параметрах запроса ? -> $1
func (d *db_) bind(query string) string {
	var result string
	switch d.driver {
	case Pgx:
		r, count := regexp.MustCompile(`\?`), 0
		result = r.ReplaceAllStringFunc(query, func(match string) string {
			count++
			return fmt.Sprintf("$%d", count)
		})
	case Mysql:
		result = query
	}

	return result
}

// Формат совместимости в параметрах запроса namedvars -> bindvars
func bindNamed(query string, arg interface{}) (string, []interface{}, error) {
	var args []interface{}
	var err error

	vof := reflect.ValueOf(arg)
	switch vof.Kind() {
	case reflect.Struct:
		query, args, err = bindStruct(query, arg)
		if err != nil {
			return "", nil, common.Wrap(err, "bindStruct")
		}
	case reflect.Map:
		query, args, err = bindMap(query, arg)
		if err != nil {
			return "", nil, common.Wrap(err, "bindMap")
		}
	default:
		return "", nil, fmt.Errorf("unsupported argument type: %s", vof.Type().String())
	}

	return query, args, nil
}

func bindStruct(query string, arg interface{}) (string, []interface{}, error) {
	tof := reflect.TypeOf(arg)
	vof := reflect.ValueOf(arg)

	mapa := make(map[string]interface{})
	for i := 0; i < tof.NumField(); i++ {
		value := vof.Field(i).Interface()
		tag := tof.Field(i).Tag.Get("db")
		if tag == "" {
			return "", nil, fmt.Errorf("tag 'db' is empty on tag %s", tag)
		}
		if strings.Contains(query, ":"+tag) {
			mapa[tag] = value
		}
	}

	return bindMap(query, mapa)
}

func bindMap(query string, arg interface{}) (string, []interface{}, error) {
	binds, ok := arg.(map[string]interface{})
	if !ok {
		return "", nil, fmt.Errorf("map must be map[string]interface{}")
	}

	i := 1
	args := make([]interface{}, 0, len(binds))
	for name, param := range binds {
		query = regexp.MustCompile(fmt.Sprintf(`(:%s)(,|\s|\)|\:|\]|;|$)`, name)).ReplaceAllString(query, fmt.Sprintf(`$$%d$2`, i))
		args = append(args, param)
		i++
	}

	for i := range args {
		if !strings.Contains(query, fmt.Sprintf("$%d", i+1)) {
			return "", nil, fmt.Errorf(`arg [%d] "%v" is not used`, i, args[i])
		}
	}

	return query, args, nil
}
