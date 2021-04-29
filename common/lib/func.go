package lib

import (
	"bytes"
	"crypto/md5"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	log2 "github.com/gin_scaffiold/common/log"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

var (
	TimeLocaltion *time.Location
	TimeFormat    = "2006-01-02 15:04:05"
	DateFormat    = "2006-01-02"
	LocalIp       = net.ParseIP("127.0.0.1")
)

//函数传入配置文件 InitModule("./conf/dev/")
func InitModule(configPath string) error {
	return initModule(configPath, []string{"base", "mysql", "redis"})
}

func initModule(configPath string, modules []string) error {
	if configPath == "" {
		fmt.Println("input config file like ./conf/dev/")
		os.Exit(1)
	}

	log.Println("------------------------------------------------------------------------")
	log.Printf("[INFO] config=%s\n", configPath)
	log.Printf("[INFO] %s\n", " start loading resources.")

	//设置ip信息,优先设置便于日志打印
	ips := GetLocalIPs()
	if len(ips) > 0 {
		LocalIp = ips[0]
	}

	//解析配置文件目录
	if err := ParseConfPath(configPath); err != nil {
		return err
	}

	//初始化配置文件
	if err := InitViperConf(); err != nil {
		return err
	}

	//加载base配置
	if InArrayString("base", modules) {
		if err := InitBaseConf(GetConfPath("base")); err != nil {
			fmt.Printf("[ERROR] %s%s\n", time.Now().Format(TimeFormat), " InitBaseConf:"+err.Error())
		}
	}

	//加载redis配置
	if InArrayString("redis", modules) {
		if err := InitRedisConf(GetConfPath("redis_map")); err != nil {
			fmt.Printf("[ERROR] %s%s\n", time.Now().Format(TimeFormat), " InitRedisConfig:"+err.Error())
		}
	}

	//加载mysql配置
	if InArrayString("mysql", modules) {
		if err := InitDBPool(GetConfPath("mysql_map")); err != nil {
			fmt.Printf("[ERROR] %s%s\n", time.Now().Format(TimeFormat), " InitDBPool:"+err.Error())
		}
	}

	//设置时区
	location, err := time.LoadLocation(ConfBase.TimeLocation)
	if err != nil {
		return err
	}
	TimeLocaltion = location

	log.Printf("[INFO] %s\n", " success loading resources.")
	log.Println("------------------------------------------------------------------------")
	return nil
}

func GetLocalIPs() (ips []net.IP) {
	interfaceAddr, err := net.InterfaceAddrs()
	if err != nil {
		return nil
	}

	for _, addr := range interfaceAddr {
		ipNet, isValidIpNet := addr.(*net.IPNet)
		if isValidIpNet && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				ips = append(ips, ipNet.IP)
			}
		}
	}

	return ips
}

func InArrayString(s string, arr []string) bool {
	for _, i := range arr {
		if i == s {
			return true
		}
	}

	return false
}

func Substr(str string, start int64, end int64) string {
	lenth := int64(len(str))
	if start < 0 || start > lenth {
		return ""
	}

	if end < 0 {
		return ""
	}

	if end > lenth {
		end = lenth
	}

	return str[start:end]
}

func Destroy() {
	log.Println("------------------------------------------------------------------------")
	log.Printf("[INFO] %s\n", " start destroy resources.")
	CloseDB()
	log2.Close()
	log.Printf("[INFO] %s\n", " success destroy resources.")
}

func HttpGet(trace *TraceContext, url string, params url.Values, msTimeout int, header http.Header) (*http.Response, []byte, error) {
	startTime := time.Now().UnixNano()
	client := http.Client{
		Timeout: time.Duration(msTimeout) * time.Millisecond,
	}

	url = AddGetDataToUrl(url, params)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		go LogTagWarn(trace, url, "GET", "", params, startTime, err, []byte{})

		return nil, nil, err
	}

	if len(header) > 0 {
		req.Header = header
	}

	req = addTrace2Header(req, trace)
	resp, err := client.Do(req)
	if err != nil {
		go LogTagWarn(trace, url, "GET", "", params, startTime, err, []byte{})

		return nil, nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if resp != nil {

		go LogTagWarn(trace, url, "GET", "", params, startTime, err, body)

		return nil, nil, err
	}

	go LogTagInfo(trace, url, "GET", "", params, startTime, err, body)

	return resp, body, nil
}

func HttpPost(trace *TraceContext, url string, params url.Values, msTimeout int, header http.Header, contextType string) (*http.Response, []byte, error) {
	startTime := time.Now().UnixNano()

	client := http.Client{
		Timeout: time.Duration(msTimeout) * time.Millisecond,
	}

	if contextType == "" {
		contextType = "application/x-www-form-urlencoded"
	}

	urlParamsEncode := params.Encode()
	req, err := http.NewRequest("POST", url, strings.NewReader(urlParamsEncode))
	if len(header) > 0 {
		req.Header = header
	}

	req = addTrace2Header(req, trace)
	req.Header.Set("Content-Type", contextType)
	resp, err := client.Do(req)
	if err != nil {

		go LogTagWarn(trace, url, "POST", urlParamsEncode, nil, startTime, err, []byte{})

		return nil, nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		go LogTagWarn(trace, url, "POST", urlParamsEncode, nil, startTime, err, body)

		return nil, nil, err
	}

	go LogTagInfo(trace, url, "POST", urlParamsEncode, nil, startTime, err, body)

	return resp, body, nil
}

func HttpJSON(trace *TraceContext, url, content string, msTimeout int, header http.Header) (*http.Response, []byte, error) {
	startTime := time.Now().UnixNano()
	client := http.Client{
		Timeout: time.Duration(msTimeout) * time.Millisecond,
	}

	req, err := http.NewRequest("POST", url, strings.NewReader(content))
	if len(header) > 0 {
		req.Header = header
	}

	req = addTrace2Header(req, trace)
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {

		go LogTagWarn(trace, url, "POST", content, nil, startTime, err, []byte{})

		return nil, nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {

		go LogTagWarn(trace, url, "POST", content, nil, startTime, err, []byte{})

		return nil, nil, err
	}
	defer resp.Body.Close()

	go LogTagInfo(trace, url, "POST", content, nil, startTime, err, body)

	return resp, body, nil
}

func LogTagWarn(ctx *TraceContext, url, method, paramStr string, params url.Values, startTime int64, err error, body []byte) {

	Log.TagWarn(ctx, DLTagHTTPFailed,
		map[string]interface{}{
			"url":       url,
			"proc_time": float32(time.Now().UnixNano()-startTime) / 1.0e9,
			"method":    method,
			"get_args":  params,
			"post_args": Substr(paramStr, 0, 1024),
			"result":    Substr(string(body), 0, 1024),
			"err":       err.Error(),
		})
}

func LogTagInfo(ctx *TraceContext, url, method, paramStr string, params url.Values, startTime int64, err error, body []byte) {
	Log.TagWarn(ctx, DLTagHTTPSuccess,
		map[string]interface{}{
			"url":       url,
			"proc_time": float32(time.Now().UnixNano()-startTime) / 1.0e9,
			"method":    method,
			"get_args":  params,
			"post_args": Substr(paramStr, 0, 1024),
			"result":    Substr(string(body), 0, 1024),
			"err":       err.Error(),
		})
}

func AddGetDataToUrl(url string, data url.Values) string {
	url = url + "?"
	if strings.Contains(url, "?") {
		url = url + "&"
	}

	return fmt.Sprintf("%s%s", url, data.Encode())
}

func addTrace2Header(r *http.Request, ctx *TraceContext) *http.Request {
	traceId := ctx.TraceId
	cSpanId := NewSpanId()

	if traceId != "" {
		r.Header.Set("didi-header-rid", traceId)
	}

	if cSpanId != "" {
		r.Header.Set("didi-header-spanid", cSpanId)
	}

	ctx.CSpanId = cSpanId

	return r
}

func NewSpanId() string {
	timestamp := uint32(time.Now().Unix())
	ipToLong := binary.BigEndian.Uint32(LocalIp.To4())
	b := bytes.Buffer{}

	b.WriteString(fmt.Sprintf("%08x", ipToLong^timestamp))
	b.WriteString(fmt.Sprintf("%08x", rand.Int31()))

	return b.String()
}

func GetMd5Hash(text string) string {
	hash := md5.New()
	hash.Write([]byte(text))

	return hex.EncodeToString(hash.Sum(nil))
}

func Encode(data string) (string, error) {
	h := md5.New()
	_, err := h.Write([]byte(data))

	if err != nil {
		return "", err
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}

func ParseServerAddr(serverAddr string) (host, port string) {
	serverInfo := strings.Split(serverAddr, ":")

	if len(serverInfo) == 2 {
		return serverInfo[0], serverInfo[1]
	}

	return serverAddr, ""
}

func NewTrace() *TraceContext {
	trace := &TraceContext{}
	trace.TraceId = GetTraceId()
	trace.SpanId = NewSpanId()

	return trace
}

func GetTraceId() string {
	return calcTraceId(LocalIp.String())
}

func calcTraceId(ip string) string {
	now := time.Now()
	timestamp := uint32(now.Unix())
	timeNano := now.UnixNano()
	pid := os.Getpid()

	b := bytes.Buffer{}
	netIp := net.ParseIP(ip)
	if netIp == nil {
		b.WriteString("00000000")
	} else {
		b.WriteString(hex.EncodeToString(netIp.To4()))
	}

	b.WriteString(fmt.Sprintf("%08x", timestamp&0xffffffff))
	b.WriteString(fmt.Sprintf("%04x", timeNano&0xffff))
	b.WriteString(fmt.Sprintf("%04x", pid&0xffff))
	b.WriteString(fmt.Sprintf("%06x", rand.Int31n(1<<24)))
	b.WriteString("b0") // 末两位标记来源,b0为go

	return b.String()
}
