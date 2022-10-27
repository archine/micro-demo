package controller

import (
	"context"
	"github.com/gin-gonic/gin"
	"gitlab.avatarworks.com/servers/component/hj-gin/mvc"
	"gitlab.avatarworks.com/servers/component/hj-gin/resp"
	"google.golang.org/grpc"
	"micro-demo/api/order"
	"micro-demo/api/user"
	"micro-demo/user/service"
	"strconv"
)

type UserController struct {
	mvc.Controller
	user.UnimplementedUserServer
	Grpc        *grpc.Server
	UserService *service.UserService
	OrderClient order.OrderClient
}

func init() {
	u := &UserController{}
	u.Prefix("/user").
		GetGroup([]*mvc.ApiInfo{
			{"/:id", u.userInfo, false},
			{"/orders/:id", u.userOrders, false},
		})
	mvc.Register(u)
}

func (u *UserController) PostConstruct() {
	user.RegisterUserServer(u.Grpc, u)
}

// 用户详情
func (u *UserController) userInfo(ctx *gin.Context) {
	userid, err := strconv.Atoi(ctx.Param("id"))
	if resp.ParamInvalid(ctx, err != nil) {
		return
	}
	info, err := u.UserService.FindUserById(userid)
	if resp.BadRequest(ctx, err != nil) {
		return
	}
	resp.Json(ctx, info)
}

// 查询用户的订单
func (u *UserController) userOrders(ctx *gin.Context) {
	userId, err := strconv.Atoi(ctx.Param("id"))
	if resp.ParamInvalid(ctx, err != nil) {
		return
	}
	orderList, err := u.OrderClient.OrderList(context.Background(), &order.OrderListRequest{Userid: int64(userId)})
	if err != nil {
		panic(err)
	}
	resp.Json(ctx, orderList)
}

// UserInfo 用户详情 grpc
func (u *UserController) UserInfo(ctx context.Context, request *user.UserInfoRequest) (*user.UserInfoResponse, error) {
	userid := request.Userid
	userInfo, err := u.UserService.FindUserById(int(userid))
	if err != nil {
		return nil, err
	}
	return &user.UserInfoResponse{
		Username: userInfo.UserName,
		Mobile:   userInfo.Mobile,
	}, nil
}
