package controller

import (
	"context"
	"github.com/archine/gin-plus/v2/mvc"
	"github.com/archine/gin-plus/v2/resp"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"micro-demo/api/order"
	"micro-demo/api/user"
	"micro-demo/user/service"
	"strconv"
)

// UserController
// @BasePath("/user")
type UserController struct {
	mvc.Controller
	user.UnimplementedUserServer
	Grpc        *grpc.Server
	UserService *service.UserService
	OrderClient order.OrderClient
}

func init() {
	mvc.Register(&UserController{})
}

func (u *UserController) PostConstruct() {
	user.RegisterUserServer(u.Grpc, u)
}

// FindUserById
// @GET(path="/:id", globalFunc=false) 用户详情
func (u *UserController) FindUserById(ctx *gin.Context) {
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

// FindOrdersById
// @GET(path="/orders/:id", globalFunc=false) 查询用户的订单
func (u *UserController) FindOrdersById(ctx *gin.Context) {
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
