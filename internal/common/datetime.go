package common

import (
	"database/sql/driver"
	"regexp"
	"time"
)

type DateTime struct {
	Time  time.Time
	Valid bool
}

func (d DateTime) Value() (driver.Value, error) {
	if !d.Valid {
		return nil, nil
	}
	return d.Time.Format(time.RFC3339), nil
}

func (d *DateTime) Scan(src interface{}) error {
	if t, ok := src.(time.Time); ok {
		d.Time = t
		d.Valid = true
	}
	return nil
}

func (d DateTime) MarshalJSON() ([]byte, error) {
	if d.Valid {
		return []byte(`"` + d.Time.Format(time.RFC3339) + `"`), nil
	} else {
		return []byte(`null`), nil
	}
}

func (d *DateTime) UnmarshalJSON(b []byte) error {
	if b[0] == '"' && b[len(b)-1] == '"' {
		b = b[1 : len(b)-1]
	}

	stringB := string(b)
	if stringB == "null" {
		return nil
	}

	if stringB[0] == '{' {
		var re = regexp.MustCompile(`(\d+[^"]+)`)
		f := re.FindString(stringB)
		if f != "" {
			stringB = f
		}
	}

	date, err := unmarshal(stringB)
	if err == nil {
		d.Time = date
		d.Valid = true
	}

	return err
}

func unmarshal(date string) (time.Time, error) {
	check := make([]string, 0)
	check = append(check, time.DateTime)
	check = append(check, time.DateOnly)
	check = append(check, time.RFC3339)

	var dateTime time.Time
	var err error
	for _, v := range check {
		loc, _ := time.LoadLocation("Europe/Moscow")
		dateTime, err = time.ParseInLocation(v, date, loc)
		if err == nil {
			return dateTime, nil
		}
	}
	if err != nil {
		return dateTime, Wrap(err, "time.ParseInLocation")
	}

	return dateTime, nil
}
