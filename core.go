package gofcm

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"sort"
	"strconv"
	"time"
)

type Engine struct {
	appId     string // 应用标识
	bizId     string // 游戏备案识别码
	secretKey string // 用户秘钥
	client    *http.Client
	keys      []string
	aes       cipher.AEAD
}

func New(appId, bizId, secretKey string) *Engine {
	c := http.Client{
		Timeout: 10 * time.Second,
	}
	b, err := hex.DecodeString(secretKey)
	if err != nil {
		return nil
	}
	block, err := aes.NewCipher(b)
	if err != nil {
		return nil
	}
	aead, err := cipher.NewGCM(block)
	if err != nil {
		return nil
	}

	return &Engine{
		appId:     appId,
		bizId:     bizId,
		secretKey: secretKey,
		client:    &c,
		keys:      []string{"appId", "bizId", "timestamps"},
		aes:       aead,
	}
}

const (
	check    = "https://api.wlc.nppa.gov.cn/idcard/authentication/check"
	query    = "http://api2.wlc.nppa.gov.cn/idcard/authentication/query"
	loginout = "http://api2.wlc.nppa.gov.cn/behavior/collection/loginout"
)

var errParam = errors.New("请求参数有误")

// Check 实名认证接口
func (e *Engine) Check(c CheckReqBody) (*CheckRespBody, error) {
	return e.check(check, c)
}

func (e *Engine) check(uri string, c CheckReqBody) (*CheckRespBody, error) {
	e1 := c.AI == ""
	e2 := c.IdNum == ""
	e3 := c.Name == ""
	if e1 || e2 || e3 {
		return nil, errParam
	}
	h := make(http.Header, 5)
	e.setHeader(&h)
	resp, err := e.post(uri, c, h)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result CheckRespBody
	err = json.Unmarshal(b, &result)
	return &result, err
}

// Query 实名认证结果查询接口
func (e *Engine) Query(ai string) (*CheckRespBody, error) {
	if ai == "" {
		return nil, errParam
	}
	return e.query(query, ai)
}

func (e *Engine) query(uri string, ai string) (*CheckRespBody, error) {
	if ai == "" {
		return nil, errParam
	}
	h := make(http.Header, 6)
	e.setHeader(&h)

	uri = fmt.Sprintf(uri+"?ai=%s", ai)
	query := map[string]string{
		"ai": ai,
	}

	resp, err := e.get(uri, query, h)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result CheckRespBody
	err = json.Unmarshal(b, &result)
	return &result, err
}

// Loginout 游戏用户行为数据上报接口
func (e *Engine) Loginout(l LoginoutReqBody) (*LoginoutRespBody, error) {
	return e.loginout(loginout, l)

}

func (e *Engine) loginout(uri string, l LoginoutReqBody) (*LoginoutRespBody, error) {
	if len(l.Collections) == 0 {
		return nil, errParam
	}
	// 此处仅抽第一个参数判断合法性
	if l.Collections[0].Ct == 0 && l.Collections[0].Pi == "" {
		return nil, errParam
	}
	if l.Collections[0].Ct == 2 && l.Collections[0].Di == "" {
		return nil, errParam
	}

	h := make(http.Header, 5)
	e.setHeader(&h)
	resp, err := e.post(uri, l, h)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result LoginoutRespBody
	err = json.Unmarshal(b, &result)
	return &result, err
}

// getBody 请求报文体数据进行AES-128/GCM + BASE64算法加密
func (e *Engine) getBody(body []byte) (string, error) {
	nonce := make([]byte, e.aes.NonceSize())
	_, err := io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return "", err
	}
	data := append(nonce, e.aes.Seal(nil, nonce, body, nil)...)
	return base64.StdEncoding.EncodeToString(data), nil
}

// getSign 接口签名
func (e *Engine) getSign(h http.Header, b string, query map[string]string) string {
	h.Del("Content-Type")

	keys := e.keys

	for k, v := range query {
		keys = append(keys, k)
		h.Add(k, v)
	}

	sort.Strings(keys)

	result := e.secretKey
	for _, v := range keys {
		// 文档里写 k-v ，实际是 kv
		result += v + h.Get(v)
	}
	result += b

	hash := sha256.New()
	_, _ = hash.Write([]byte(result))
	return hex.EncodeToString(hash.Sum(nil))
}

func (e *Engine) setHeader(h *http.Header) {
	h.Add("Content-Type", "application/json; charset=utf-8")
	h.Add("appId", e.appId)
	h.Add("bizId", e.bizId)
	h.Add("timestamps", strconv.FormatInt(time.Now().Unix()*1000, 10))
}

func (e *Engine) post(uri string, b interface{}, h http.Header) (*http.Response, error) {
	_b, err := json.Marshal(b)
	if err != nil {
		return nil, err
	}
	body, err := e.getBody(_b)
	if err != nil {
		return nil, err
	}
	data := `{"data":"` + body + `"}`
	h.Add("sign", e.getSign(h.Clone(), data, nil))
	r, err := http.NewRequest(http.MethodPost, uri, bytes.NewReader([]byte(data)))
	if err != nil {
		return nil, err
	}
	for k, v := range h {
		r.Header.Set(k, v[0])
	}
	return e.client.Do(r)
}

func (e *Engine) get(uri string, query map[string]string, h http.Header) (*http.Response, error) {
	h.Add("sign", e.getSign(h.Clone(), "", query))
	r, err := http.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		return nil, err
	}
	for k, v := range h {
		r.Header.Set(k, v[0])
	}
	return e.client.Do(r)
}

// GetCNErrMsg 获取中文错误消息，仅包含 check 与 query 业务异常
func (e *Engine) GetCNErrMsg(state int) string {
	switch state {
	case 0:
		return "请求成功"
	case 1001:
		return "系统错误"
	case 1002:
		return "接口请求的资源不存在"
	case 1003:
		return "接口请求方式错误"
	case 1004:
		return "接口请求核心参数缺失"
	case 1005:
		return "接口请求IP地址非法"
	case 1006:
		return "接口请求超出流量限制"
	case 1007:
		return "接口请求过期"
	case 1008:
		return "接口请求方身份非法"
	case 1009:
		return "接口请求方权限未启用"
	case 1010:
		return "接口请求方无该接口权限"
	case 1011:
		return "接口请求方身份核验错误"
	case 1012:
		return "接口请求报文核验失败"
	case 2001:
		return "身份证号格式校验失败"
	case 2002:
		return "实名认证条目已达上限"
	case 2003:
		return "无该编码提交的实名认证记录"
	case 2004:
		return "编码已经被占用"
	}
	return strconv.Itoa(state)
}
