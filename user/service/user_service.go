package service

import (
	"errors"
	ioc "gitlab.avatarworks.com/servers/component/hj-ioc"
	"micro-demo/user/controller/vo"
)

type UserService struct{}

func (u *UserService) CreateBean() ioc.Bean {
	return &UserService{}
}

// FindUserById 查询指定用户
func (u *UserService) FindUserById(userid int) (*vo.UserInfo, error) {
	if userid == 1 {
		return &vo.UserInfo{
			UserName: "张三",
			Mobile:   "15424556662",
		}, nil
	}
	if userid == 2 {
		return &vo.UserInfo{
			UserName: "里斯",
			Mobile:   "13300001111",
		}, nil
	}
	return nil, errors.New("用户不存在")
}
