package utils

import (
	"crypto/md5"
	"fmt"
	"io"
)

func Md5Str(str string) string {
	h := md5.New()
	io.WriteString(h, str)

	psd := fmt.Sprintf("%x", h.Sum(nil))

	return psd
}
