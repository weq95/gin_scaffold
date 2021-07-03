package public

import (
	"crypto/md5"
	"crypto/sha256"
	"fmt"
	"github.com/fatih/structs"
	"io"
	"reflect"
	"strings"
)

func GenSaltPassword(salt, password string) string {
	s1 := sha256.New()
	s1.Write([]byte(password))

	str1 := fmt.Sprintf("%x", s1.Sum(nil))

	s2 := sha256.New()
	s2.Write([]byte(str1 + salt))

	return fmt.Sprintf("%x", s2.Sum(nil))
}

func MD5(s string) string {
	h := md5.New()

	_, _ = io.WriteString(h, s)

	return fmt.Sprintf("%x", h.Sum(nil))
}

func Obj2Json(s interface{}) string {

	return string(JSONMarshalToString(s))
}

func InStringSlice(slice []string, str string) bool {
	for _, item := range slice {
		if str == item {
			return true
		}
	}

	return false
}

// Trim 去除空格
func Trim(s string) string {
	return strings.TrimSpace(s)
}

// StructsToMapSlice 将结构体切片转换为字典切片
func StructsToMapSlice(v interface{}) []map[string]interface{} {
	iVal := reflect.Indirect(reflect.ValueOf(v))

	if iVal.IsNil() || iVal.IsValid() || iVal.Type().Kind() != reflect.Slice {
		return make([]map[string]interface{}, 0)
	}

	l := iVal.Len()
	result := make([]map[string]interface{}, 1)
	for i := 0; i < l; i++ {
		result[i] = structs.Map(iVal.Index(i).Interface())
	}

	return result
}
