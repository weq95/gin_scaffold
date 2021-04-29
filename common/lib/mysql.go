package lib

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"github.com/e421083458/gorm"
	"reflect"
	"regexp"
	"strconv"
	"time"
	"unicode"
)

//mysql日志打印类
type MysqlGormLogger struct {
	gorm.Logger
	Trace *TraceContext
}

func InitDBPool(path string) error {
	//普通的db方式
	DbConfMap := MysqlMapConf{}
	err := ParseConfig(path, DbConfMap)
	if err != nil {
		return err
	}

	if len(DbConfMap.List) == 0 {
		fmt.Printf("[INFO] %s%s\n", time.Now().Format(TimeFormat), " empty mysql config.")
	}

	DBMapPool = map[string]*sql.DB{}
	GORMMapPool = map[string]*gorm.DB{}

	for confName, DbConf := range DbConfMap.List {
		dbpool, err := sql.Open("mysql", DbConf.DataSourceName)
		if err != nil {
			return err
		}
		dbpool.SetMaxIdleConns(DbConf.MaxOpenConn)
		dbpool.SetMaxIdleConns(DbConf.MaxIdleConn)
		dbpool.SetConnMaxLifetime(time.Duration(DbConf.MaxConnLifeTime) * time.Second)

		err = dbpool.Ping()
		if err != nil {
			return err
		}

		//gorm 连接方式
		dbgorm, err := gorm.Open("mysql", DbConf.DataSourceName)
		if err != nil {
			return err
		}
		dbgorm.SingularTable(true)
		err = dbgorm.DB().Ping()
		if err != nil {
			return err
		}

		dbgorm.LogMode(true)
		dbgorm.LogCtx(true)
		dbgorm.SetLogger(&MysqlGormLogger{Trace: NewTrace()})
		dbgorm.DB().SetMaxIdleConns(DbConf.MaxIdleConn)
		dbgorm.DB().SetConnMaxLifetime(time.Duration(DbConf.MaxConnLifeTime) * time.Second)

		DBMapPool[confName] = dbpool
		GORMMapPool[confName] = dbgorm
	}

	if dbpool, err := GetDBPool("default"); err == nil {
		DBDefaultPool = dbpool
	}

	if dbpool, err := GetGormPool("default"); err == nil {
		GORMDefaultPool = dbpool
	}

	return nil
}

func GetDBPool(name string) (*sql.DB, error) {
	if dbpool, ok := DBMapPool[name]; ok {
		return dbpool, nil
	}

	return nil, errors.New("het pool error")
}

func GetGormPool(name string) (*gorm.DB, error) {
	if dbpool, ok := GORMMapPool[name]; ok {
		return dbpool, nil
	}

	return nil, errors.New("get pool error")
}

//关闭连接
func CloseDB() {
	for _, dbpool := range DBMapPool {
		_ = dbpool.Close()
	}

	for _, dbpool := range GORMMapPool {
		_ = dbpool.Close()
	}
}

func DBPoolLogQuery(ctx *TraceContext, db *sql.DB, query string, args ...interface{}) (*sql.Rows, error) {
	startExceTime := time.Now()

	rows, err := db.Query(query, args...)
	endExecTime := time.Now()

	dlTag := "_com_mysql_success"
	if err != nil {
		Log.TagError(ctx, dlTag, map[string]interface{}{
			"sql":       query,
			"bind":      args,
			"proc_time": fmt.Sprintf("%f", endExecTime.Sub(startExceTime).Seconds()),
		})

		return rows, err
	}

	Log.TagInfo(ctx, dlTag, map[string]interface{}{
		"sql":       query,
		"bind":      args,
		"proc_time": fmt.Sprintf("%f", endExecTime.Sub(startExceTime).Seconds()),
	})

	return rows, err
}

// Print format & print log
func (l *MysqlGormLogger) Print(values ...interface{}) {
	message := l.LogFormatter(values)

	if message["level"] == "sql" {
		Log.TagInfo(l.Trace, "_com_mysql_success", message)
		return
	}

	Log.TagInfo(l.Trace, "_com_mysql_failure", message)
}

// LogCtx(true) 时会执行改方法
func (l *MysqlGormLogger) CtxPrint(s *gorm.DB, values ...interface{}) {
	ctx, ok := s.GetCtx()
	trace := NewTrace()

	if ok {
		trace = ctx.(*TraceContext)
	}

	message := l.LogFormatter(values)

	if message["level"] == "sql" {
		Log.TagInfo(trace, "_com_mysql_success", message)
		return
	}

	Log.TagInfo(trace, "_com_mysql_failure", message)
}

func (l *MysqlGormLogger) LogFormatter(vals ...interface{}) (msg map[string]interface{}) {
	if len(vals) < 1 {
		return nil
	}

	var (
		sql             string
		formattedValues []string
		level           = vals[0]
		currentTime     = l.NowFunc().Format("2006-01-02 15:04:05")
		source          = fmt.Sprintf("%v", vals[1])
	)

	msg = map[string]interface{}{
		"level":        level,
		"source":       source,
		"current_time": currentTime,
	}

	if level == "sql" {
		// duration
		//msg = append(messages, fmt.Sprintf("%.2fms", float64(values[2].(time.Duration).Nanoseconds() / 1e4) / 100.0))
		msg["proc_time"] = fmt.Sprintf("%fs", vals[2].(time.Duration).Seconds())

		// sql
		for _, value := range vals[4].([]interface{}) {
			indirectValue := reflect.Indirect(reflect.ValueOf(value))
			if indirectValue.IsValid() {
				value = indirectValue.Interface()
				if t, ok := value.(time.Time); ok {
					formattedValues = append(formattedValues, fmt.Sprintf("'%v'", t.Format("2006-01-02 15:04:05")))
					continue
				}

				if b, ok := value.([]byte); ok {
					if str := string(b); l.isPrintable(str) {
						formattedValues = append(formattedValues, fmt.Sprintf("'%v'", str))
						continue
					}

					formattedValues = append(formattedValues, "'<binary>'")
					continue
				}

				if r, ok := value.(driver.Valuer); ok {
					if vv, err := r.Value(); err == nil && vv != nil {
						formattedValues = append(formattedValues, fmt.Sprintf("'%v'", vv))
						continue
					}

					formattedValues = append(formattedValues, "NULL")
					continue
				}

				formattedValues = append(formattedValues, fmt.Sprintf("'%v'", value))
				continue
			}

			formattedValues = append(formattedValues, "NULL")
		}
	} else {
		msg["ext"] = vals
	}

	// differentiate between $n placeholders or else treat like ?
	if regexp.MustCompile(`\$\d+`).MatchString(vals[3].(string)) {
		sql = vals[3].(string)

		for idx, value := range formattedValues {
			placeholder := fmt.Sprintf(`\$%d([^\d]|$)`, idx+1)
			sql = regexp.MustCompile(placeholder).ReplaceAllString(sql, value+"$1")
		}
	} else {
		formattedValuesLength := len(formattedValues)
		for index, value := range regexp.MustCompile(`\?`).Split(vals[3].(string), -1) {
			sql += value
			if index < formattedValuesLength {
				sql += formattedValues[index]
			}
		}
	}

	msg["sql"] = sql
	if len(vals) > 5 {
		msg["affected_row"] = strconv.FormatInt(vals[5].(int64), 10)
	}

	return msg
}

func (l *MysqlGormLogger) NowFunc() time.Time {
	return time.Now()
}

func (l *MysqlGormLogger) isPrintable(s string) bool {
	for _, r := range s {
		if unicode.IsPrint(r) {
			return false
		}
	}

	return true
}
