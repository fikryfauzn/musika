package models

import (
	"database/sql/driver"
	"errors"
	"time"
)

// DateOnly is a custom type for handling dates in "DD-MM-YYYY" format
type DateOnly struct {
	time.Time
}

const indonesianDateFormat = "02-01-2006"

// MarshalJSON formats the date as "DD-MM-YYYY" for API responses
func (d DateOnly) MarshalJSON() ([]byte, error) {
	return []byte(`"` + d.Time.Format(indonesianDateFormat) + `"`), nil
}

// UnmarshalJSON parses the date from "DD-MM-YYYY" in API requests
func (d *DateOnly) UnmarshalJSON(data []byte) error {
	parsedTime, err := time.Parse(`"`+indonesianDateFormat+`"`, string(data))
	if err != nil {
		return errors.New("invalid date format, use DD-MM-YYYY")
	}
	d.Time = parsedTime
	return nil
}

// Value implements the driver.Valuer interface for database storage
func (d DateOnly) Value() (driver.Value, error) {
	if d.IsZero() {
		return nil, nil
	}
	return d.Time.Format("2006-01-02"), nil // Store in DB as "YYYY-MM-DD"
}

// Scan implements the sql.Scanner interface for reading from the database
func (d *DateOnly) Scan(value interface{}) error {
	if value == nil {
		*d = DateOnly{}
		return nil
	}
	switch v := value.(type) {
	case time.Time:
		d.Time = v
		return nil
	case []byte:
		parsedTime, err := time.Parse("2006-01-02", string(v))
		if err != nil {
			return err
		}
		d.Time = parsedTime
		return nil
	case string:
		parsedTime, err := time.Parse("2006-01-02", v)
		if err != nil {
			return err
		}
		d.Time = parsedTime
		return nil
	default:
		return errors.New("unsupported type for DateOnly")
	}
}

// SwaggerType tells Swagger to treat DateOnly as a string
func (DateOnly) SwaggerType() string {
	return "string"
}

// SwaggerFormat tells Swagger the format of the string
func (DateOnly) SwaggerFormat() string {
	return "date"
}
