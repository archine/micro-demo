package service

import (
	"micro-demo/api/order"
)

type OrderService struct {
}

func (o *OrderService) FindOrderListGrpc(userid int64) []*order.OrderListResponse_OrderInfo {
	if userid == 1 {
		return []*order.OrderListResponse_OrderInfo{
			{No: "10001", Price: 22.5, Status: 1},
			{No: "10002", Price: 37, Status: 3},
		}
	}
	return []*order.OrderListResponse_OrderInfo{
		{No: "10003", Price: 34.5, Status: 2},
		{No: "10004", Price: 56, Status: 3},
	}
}

func (o *OrderService) FindOrderInfo(orderId int64) []string {
	return []string{"烧的坏水壶", "吹热风的风扇"}
}
