package gh

import (
	"net/url"

	"github.com/webx-top/restyclient"
)

// documention: https://cloud.tencent.com/developer/article/2245869
// https://developers.weixin.qq.com/doc/offiaccount/Account_Management/Generating_a_Parametric_QR_Code.html

const (
	ticketURL = `https://api.weixin.qq.com/cgi-bin/qrcode/create?access_token=`
	qrcodeURL = `https://mp.weixin.qq.com/cgi-bin/showqrcode?ticket=`
)

/*
	{
	    "expire_seconds": 604800,
	    "action_name": "QR_SCENE",
	    "action_info": {
	        "scene": {
	            "scene_id": 123
	        }
	    }
	}
*/
type RequestQRCodeCreate struct {
	ExpireSeconds int        `json:"expire_seconds,omitempty"`
	ActionName    string     `json:"action_name"`
	ActionInfo    ActionInfo `json:"action_info"`
}

type ActionInfo struct {
	Scene Scene `json:"scene"`
}

type Scene struct {
	SceneID  int64  `json:"scene_id,omitempty"`  // action_name=QR_SCENE 场景值ID，临时二维码时为32位非0整型，永久二维码时最大值为100000（目前参数只支持1--100000）
	SceneStr string `json:"scene_str,omitempty"` // action_name=QR_STR_SCENE 场景值ID（字符串形式的ID），字符串类型，长度限制为1到64
}

type ResponseQRCodeCreate struct {
	ExpireSeconds int    `json:"expire_seconds"`
	Ticket        string `json:"ticket"`
	URL           string `json:"url"`
}

func GetQRCodeURL(accessToken string, scene string, expireSeconds ...int) (string, error) {
	req := restyclient.Retryable()
	data := &RequestQRCodeCreate{
		ExpireSeconds: 604800,
		ActionName:    `QR_STR_SCENE`,
		ActionInfo: ActionInfo{
			Scene: Scene{
				SceneStr: scene,
			},
		},
	}
	if len(expireSeconds) > 0 {
		data.ExpireSeconds = expireSeconds[0]
	}
	res := &ResponseQRCodeCreate{}
	_, err := req.SetBody(data).SetResult(res).Post(ticketURL + accessToken)
	if err != nil {
		return ``, err
	}
	return qrcodeURL + url.QueryEscape(res.Ticket), err
}
