package configs

import (
	"fmt"
	"net/url"
)

func (db DBConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		db.Host,
		db.Port,
		db.User,
		db.Password,
		db.Name,
		db.SSLMode,
	)
}

func (db DBConfig) URL() string {
	u := &url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(db.User, db.Password),
		Host:   fmt.Sprintf("%s:%s", db.Host, db.Port),
		Path:   db.Name,
	}

	query := u.Query()
	query.Set("sslmode", db.SSLMode)
	u.RawQuery = query.Encode()

	return u.String()
}
