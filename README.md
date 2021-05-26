更新时间 2021 年 05 月 26 日

版本号: v1.8.0
[网络游戏防沉迷实名认证系统 接口对接技术规范](https://res.wx.qq.com/op_res/DnyOZwFQxaP-SJPdqnAvqrMOI2g4C8Hykha3Br5XOLlt0xc883qI9813oM1aH_4h03B2XT05qRMxIiWSU-ggrw)

[网络游戏防沉迷实名认证系统测试系统说明](https://wlc.nppa.gov.cn/fcm_company/网络游戏防沉迷实名认证系统测试系统说明.pdf)

## 安装

```bash
go get -u github.com/ixugo/gofcm
```

## 通过测试用例

![image-20210526204814544](http://img.golang.space/shot-1622033294800.png)

测试文件中有 9 个测试函数，前 8 个用于通过接口测试。
需要更改 TestMain 函数中的参数 AppId，bizId，secretKey

将官网获取的测试码一次写入 code 数组
执行测试

## 如何使用

```go
import "github.com/ixugo/gofcm"

func main(){
    gofcm.New(appId,bitId,secretKey)
}
```
