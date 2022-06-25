package dnspod

import "fmt"

const (
	methodUserDetail = "User.Detail"
)

type User struct {
	RealName          string `json:"real_name"`
	UserType          string `json:"user_type"`
	Telephone         string `json:"telephone"`
	Im                string `json:"im"`
	Nick              string `json:"nick"`
	Id                string `json:"id"`
	Email             string `json:"email"`
	Status            string `json:"status"`
	EmailVerified     string `json:"email_verified"`
	TelephoneVerified string `json:"telephone_verified"`
	WeixinBinded      string `json:"weixin_binded"`
	AgentPending      bool   `json:"agent_pending"`
	Balance           int    `json:"balance"`
	Smsbalance        int    `json:"smsbalance"`
	UserGrade         string `json:"user_grade"`
}

type Agent struct {
	Discount     string `json:"discount"`
	Points       string `json:"points"`
	BalanceLimit string `json:"balance_limit"`
	Users        string `json:"users"`
}

type UserInfo struct {
	User  User  `json:"user"`
	Agent Agent `json:"agent"`
}

type userWrapper struct {
	Status Status   `json:"status"`
	Info   UserInfo `json:"info"`
}

type UserService struct {
	client *Client
}

func (u *UserService) Profile() (UserInfo, *Response, error) {
	payload := u.client.CommonParams.toPayLoad()

	returnedUserInfo := userWrapper{}

	res, err := u.client.post(methodUserDetail, payload, &returnedUserInfo)
	if err != nil {
		return UserInfo{}, nil, err
	}

	if returnedUserInfo.Status.Code != "1" {
		return UserInfo{}, nil, fmt.Errorf("code: %s, message: %s", returnedUserInfo.Status.Code, returnedUserInfo.Status.Message)
	}

	if u.client.CommonParams.IsInternational && returnedUserInfo.Info.User.UserGrade == "" {
		returnedUserInfo.Info.User.UserGrade = "DP_Free"
	}

	return returnedUserInfo.Info, res, nil
}
