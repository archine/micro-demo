package vo

import "github.com/archine/ioc"

type UserInfo struct {
	UserName string `json:"user_name"`
	Mobile   string `json:"mobile"`
}

type UserInfoMapper struct {}

func (u *UserInfoMapper) CreateBean() ioc.Bean {
	return &UserInfoMapper{}
}
