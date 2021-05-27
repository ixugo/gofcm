package gofcm

import (
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"
)

const (
	testCheck    = "https://wlc.nppa.gov.cn/test/authentication/check/%s"
	testQuery    = "https://wlc.nppa.gov.cn/test/authentication/query/%s"
	testLoginout = "https://wlc.nppa.gov.cn/test/collection/loginout/%s"
)

var code = []string{
	"MhQ23v", "w6HkUY", "63VT9a", "X6xBeJ",
	"sKFMbR", "N3bAKP", "nz324r", "DpSFTj",
}

var e *Engine

func TestMain(m *testing.M) {
	appId := "6e1645a413f0437a96abab8f46d3aebd"
	bizId := "1101999999"
	secretKey := "0f9193177fb67f9b7fc385a0fa7bc4b1"
	e = New(appId, bizId, secretKey)

	os.Exit(m.Run())
}

// TestCheck1 测试1:实名认证结果返回“认证成功”
func TestCheck1(t *testing.T) {
	uri := fmt.Sprintf(testCheck, code[0])
	data := []CheckReqBody{
		{"100000000000000001", "某一一", "110000190101010001"},
		{"100000000000000002", "某一二", "110000190101020007"},
		{"100000000000000003", "某一三", "110000190101030002"},
		{"100000000000000004", "某一四", "110000190101040008"},
		{"100000000000000005", "某一五", "11000019010101001X"},
		{"100000000000000006", "某一六", "110000190101020015"},
		{"100000000000000007", "某一七", "110000190101030010"},
		{"100000000000000008", "某一八", "110000190101040016"},
	}

	for _, v := range data {
		resp, err := e.check(uri, v)
		if err != nil {
			t.Fatal(err)
		}
		if resp.ErrCode != 0 {
			t.Fatal("状态码错误")
		}

		if resp.Data.Result.Status != 0 {
			t.Fatal("认证结果异常")
		}
	}
}

// TestCheck2 测试2:实名认证结果返回“认证中”
func TestCheck2(t *testing.T) {
	uri := fmt.Sprintf(testCheck, code[1])
	data := []CheckReqBody{
		{"200000000000000001", "某二一", "110000190201010009"},
		{"200000000000000002", "某二二", "110000190201020004"},
		{"200000000000000003", "某二三", "11000019020103000X"},
		{"200000000000000004", "某二四", "110000190201040005"},
		{"200000000000000005", "某二五", "110000190201010017"},
		{"200000000000000006", "某二六", "110000190201020012"},
		{"200000000000000007", "某二七", "110000190201030018"},
		{"200000000000000008", "某二八", "110000190201040013"},
	}

	for _, v := range data {
		resp, err := e.check(uri, v)
		if err != nil {
			t.Fatal(err)
		}
		if resp.ErrCode != 0 {
			t.Fatal("状态码错误")
		}

		if resp.Data.Result.Status != 1 {
			t.Fatal("认证结果异常")
		}
	}
}

// TestCheck3 测试3:实名认证结果返回“认证失败”
func TestCheck3(t *testing.T) {
	uri := fmt.Sprintf(testCheck, code[2])
	data := []CheckReqBody{
		{"300000000000000001", "小S", "110000190201010017"},
		{"300000000000000002", "大S", "110000190201010017"},
	}

	for _, v := range data {
		resp, err := e.check(uri, v)
		if err != nil {
			t.Fatal(err)
		}
		if resp.ErrCode != 0 {
			t.Fatal("状态码错误", resp.ErrCode)
		}

		if resp.Data.Result.Status != 2 {
			t.Fatal("认证结果异常")
		}
	}
}

// TestQuery4 测试4:实名认证结果返回“认证成功”
func TestQuery4(t *testing.T) {
	uri := fmt.Sprintf(testQuery, code[3])
	// uri = "https://wlc.nppa.gov.cn/test/authentication/query"
	data := []CheckReqBody{
		{AI: "100000000000000001"},
		{AI: "100000000000000002"},
		{AI: "100000000000000003"},
		{AI: "100000000000000004"},
		{AI: "100000000000000005"},
		{AI: "100000000000000006"},
		{AI: "100000000000000007"},
		{AI: "100000000000000008"},
	}

	for _, v := range data {
		resp, err := e.query(uri, v.AI)
		if err != nil {
			t.Fatal(err)
		}
		if resp.ErrCode != 0 {
			t.Fatal("状态码错误 code:", resp.ErrCode)
		}

		if resp.Data.Result.Status != 0 {
			t.Fatal("认证结果异常")
		}
	}
}

// TestQuery5 测试5:实名认证结果返回“认证中”
func TestQuery5(t *testing.T) {
	uri := fmt.Sprintf(testQuery, code[4])
	data := []CheckReqBody{
		{AI: "200000000000000001"},
		{AI: "200000000000000002"},
		{AI: "200000000000000003"},
		{AI: "200000000000000004"},
		{AI: "200000000000000005"},
		{AI: "200000000000000006"},
		{AI: "200000000000000007"},
		{AI: "200000000000000008"},
	}

	for _, v := range data {
		resp, err := e.query(uri, v.AI)
		if err != nil {
			t.Fatal(err)
		}
		if resp.ErrCode != 0 {
			t.Fatal("状态码错误 code:", resp.ErrCode)
		}

		if resp.Data.Result.Status != 1 {
			t.Fatal("认证结果异常")
		}
	}
}

// TestQuery6 测试6:实名认证结果返回“认证失败”
func TestQuery6(t *testing.T) {
	uri := fmt.Sprintf(testQuery, code[5])
	data := []CheckReqBody{
		{AI: "300000000000000001"},
		{AI: "300000000000000002"},
	}

	for _, v := range data {
		resp, err := e.query(uri, v.AI)
		if err != nil {
			t.Fatal(err)
		}
		if resp.ErrCode != 0 {
			t.Fatal("状态码错误 code:", resp.ErrCode)
		}

		if resp.Data.Result.Status != 2 {
			t.Fatal("认证结果异常")
		}
	}
}

// TestLoginout7 测试7:模拟“游客模式”下游戏用户行为 数据上报场景
func TestLoginout7(t *testing.T) {
	uri := fmt.Sprintf(testLoginout, code[6])
	data := LoginoutReqBody{
		Collections: []Collection{
			{
				No: 2,
				Si: "100000",
				Bt: 1,
				Ot: time.Now().Add(-5 * time.Second).Unix(),
				Ct: 2,
				Di: "12121212121212121212121212121212",
			},
		},
	}

	resp, err := e.loginout(uri, data)
	if err != nil {
		t.Fatal(err)
	}
	if resp.ErrCode != 0 {
		t.Fatal("状态码错误 code:", resp.ErrCode, " ", resp.Data)
	}

}

// TestLoginout8 测试8:模拟“已认证”游戏用户的行为数 据上报场景
func TestLoginout8(t *testing.T) {
	uri := fmt.Sprintf(testLoginout, code[7])
	data := LoginoutReqBody{
		Collections: []Collection{
			{
				No: 1,
				Si: "100000000000000008",
				Bt: 1,
				Ot: time.Now().Add(-5 * time.Second).Unix(),
				Ct: 0,
				Pi: "1fffbjzos82bs9cnyj1dna7d6d29zg4esnh99u",
			},
		},
	}

	resp, err := e.loginout(uri, data)
	if err != nil {
		t.Fatal(err)
	}
	if resp.ErrCode != 0 {
		t.Fatal("状态码错误 code:", resp.ErrCode)
	}

}

func TestGetSign(t *testing.T) {
	h := make(http.Header, 6)
	e.setHeader(&h)
	h.Set("timestamps", "1622028535743")

	b := map[string]string{
		"ai": "100000000000000001",
	}
	fmt.Println(e.getSign(h, "", b))
}
