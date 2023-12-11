/* 实现Scan方法，支持在gorm中的值设置为[]string */
package tables

import (
	"database/sql/driver"
	"errors"
	"strings"
)

type StringSlice []string

func (ss *StringSlice) Scan(src interface{}) error {
	asBytes, ok := src.([]byte)
	if !ok {
		return errors.New("scan source was not []bytes")
	}
	str := string(asBytes)
	*ss = strings.Split(str, ",")
	return nil
}

func (ss StringSlice) Value() (driver.Value, error) {
	return strings.Join(ss, ","), nil
}
