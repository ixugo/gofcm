package gofcm

// CheckReqBody 请求体
type CheckReqBody struct {
	AI    string `json:"ai,omitempty"`    // 游戏内部成员标识
	Name  string `json:"name,omitempty"`  // 姓名
	IdNum string `json:"idNum,omitempty"` // 身份证
}

// CheckRespBody check 和 query 的响应内容
type CheckRespBody struct {
	ErrCode int           `json:"errcode,omitempty"` // 状态码 0:成功; x>1000:失败;
	ErrMsg  string        `json:"errmsg,omitempty"`  // 状态描述
	Data    checkRespBody `json:"data,omitempty"`
}

type checkRespBody struct {
	Result struct {
		Status int    `json:"status,omitempty"` // 0:认证成功; 1: 认证中; 2:认证失败;
		Pi     string `json:"pi,omitempty"`     // 用户唯一标识
	} `json:"result,omitempty"`
}

type LoginoutReqBody struct {
	Collections []Collection `json:"collections,omitempty"`
}

type Collection struct {
	No int    `json:"no"` // 条码编码
	Si string `json:"si"` // 游戏内部会话标识
	Bt int    `json:"bt"` // 用户行为类型  0:下线; 1:上线;
	Ot int64  `json:"ot"` // 行为发生时间，单位秒
	Ct int    `json:"ct"` // 上报类型 0:已认证用户; 2: 游客用户
	Di string `json:"di"` // 设备标识，游客用户必填
	Pi string `json:"pi"` // 用户唯一标识 已认证用户必填
}

type LoginoutRespBody struct {
	ErrCode int               `json:"errcode,omitempty"` // 状态码 0:成功; x>1000:失败;
	ErrMsg  string            `json:"errmsg,omitempty"`  // 状态描述
	Data    *loginoutRespBody `json:"data,omitempty"`
}

type loginoutRespBody struct {
	Result []struct {
		No     int    `json:"no,omitempty"`     // 条目编码
		Status int    `json:"status,omitempty"` // 0:认证成功; 1: 认证中; 2:认证失败;
		Pi     string `json:"pi,omitempty"`     // 用户唯一标识
	} `json:"result,omitempty"`
}
