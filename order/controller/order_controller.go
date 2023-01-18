package controller

import (
	"context"
	"github.com/archine/gin-plus/v2/mvc"
	"github.com/archine/gin-plus/v2/resp"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"micro-demo/api/order"
	"micro-demo/api/user"
	"micro-demo/order/controller/vo"
	"micro-demo/order/service"
	"strconv"
)

// OrderController
// @BasePath("/order")
type OrderController struct {
	mvc.Controller
	order.UnimplementedOrderServer
	Grpc         *grpc.Server
	OrderService *service.OrderService
	UserClient   user.UserClient
}

func init() {
	mvc.Register(&OrderController{})
}

func (o *OrderController) PostConstruct() {
	// 当前controller加入到grpc服务器中
	order.RegisterOrderServer(o.Grpc, o)
}

// OrderList 用户订单列表 grpc调用
func (o *OrderController) OrderList(ctx context.Context, request *order.OrderListRequest) (*order.OrderListResponse, error) {
	return &order.OrderListResponse{
		Data: o.OrderService.FindOrderListGrpc(request.Userid),
	}, nil
}

// OrderInfo
// @GET(path="/:id", globalFunc=false) 订单详情
func (o *OrderController) OrderInfo(ctx *gin.Context) {
	orderId, err := strconv.Atoi(ctx.Param("id"))
	if resp.ParamInvalid(ctx, err != nil) {
		return
	}
	goods := o.OrderService.FindOrderInfo(int64(orderId))
	data := vo.OrderDetail{
		Goods: goods,
	}
	info, err := o.UserClient.UserInfo(context.Background(), &user.UserInfoRequest{Userid: 2})
	if err != nil {
		panic(err)
	}
	data.Mobile = info.GetMobile()
	data.Username = info.GetUsername()
	resp.Json(ctx, data)
}
