package service

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"sync"
	"time"

	"github.com/go-redis/redis"
	logging "github.com/sirupsen/logrus"

	"mall/pkg/e"
	"mall/repository/cache"
	"mall/repository/db/dao"
	"mall/repository/db/model"
	"mall/serializer"
	"mall/types"
)

const OrderTimeKey = "OrderTime"

var OrderSrvIns *OrderSrv
var OrderSrvOnce sync.Once

type OrderSrv struct {
}

func GetOrderSrv() *OrderSrv {
	OrderSrvOnce.Do(func() {
		OrderSrvIns = &OrderSrv{}
	})
	return OrderSrvIns
}

func (s *OrderSrv) OrderCreate(ctx context.Context, id uint, req *types.OrderServiceReq) (serializer.Response, error) {
	code := e.SUCCESS

	order := &model.Order{
		UserID:    id,
		ProductID: req.ProductID,
		BossID:    req.BossID,
		Num:       int(req.Num),
		Money:     float64(req.Money),
		Type:      1,
	}
	addressDao := dao.NewAddressDao(ctx)
	address, err := addressDao.GetAddressByAid(req.AddressID)
	if err != nil {
		logging.Info(err)
		code = e.ErrorDatabase
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}, err
	}

	order.AddressID = address.ID
	number := fmt.Sprintf("%09v", rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(1000000000))
	productNum := strconv.Itoa(int(req.ProductID))
	userNum := strconv.Itoa(int(id))
	number = number + productNum + userNum
	orderNum, _ := strconv.ParseUint(number, 10, 64)
	order.OrderNum = orderNum

	orderDao := dao.NewOrderDao(ctx)
	err = orderDao.CreateOrder(order)
	if err != nil {
		logging.Info(err)
		code = e.ErrorDatabase
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}, err
	}

	// 订单号存入Redis中，设置过期时间
	data := redis.Z{
		Score:  float64(time.Now().Unix()) + 15*time.Minute.Seconds(),
		Member: orderNum,
	}
	cache.RedisClient.ZAdd(OrderTimeKey, data)
	return serializer.Response{
		Status: code,
		Msg:    e.GetMsg(code),
	}, nil
}

func (s *OrderSrv) OrderList(ctx context.Context, uId uint, req *types.OrderServiceReq) (serializer.Response, error) {
	var orders []*model.Order
	var total int64
	code := e.SUCCESS
	if req.PageSize == 0 {
		req.PageSize = 5
	}

	orderDao := dao.NewOrderDao(ctx)
	condition := make(map[string]interface{})
	condition["user_id"] = uId

	if req.Type == 0 {
		condition["type"] = 0
	} else {
		condition["type"] = req.Type
	}
	orders, total, err := orderDao.ListOrderByCondition(condition, req.BasePage)
	if err != nil {
		code = e.ErrorDatabase
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}, err
	}

	return serializer.BuildListResponse(serializer.BuildOrders(ctx, orders), uint(total)), nil
}

func (s *OrderSrv) OrderShow(ctx context.Context, uId uint, req *types.OrderServiceReq) (serializer.Response, error) {
	code := e.SUCCESS

	orderDao := dao.NewOrderDao(ctx)
	order, _ := orderDao.GetOrderById(uId, req.OrderId)

	addressDao := dao.NewAddressDao(ctx)
	address, err := addressDao.GetAddressByAid(order.AddressID)
	if err != nil {
		logging.Info(err)
		code = e.ErrorDatabase
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}, err
	}

	productDao := dao.NewProductDao(ctx)
	product, err := productDao.GetProductById(order.ProductID)
	if err != nil {
		logging.Info(err)
		code = e.ErrorDatabase
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}, err
	}

	return serializer.Response{
		Status: code,
		Msg:    e.GetMsg(code),
		Data:   serializer.BuildOrder(order, product, address),
	}, nil
}

func (s *OrderSrv) OrderDelete(ctx context.Context, uId uint, req *types.OrderServiceReq) (serializer.Response, error) {
	code := e.SUCCESS

	orderDao := dao.NewOrderDao(ctx)
	err := orderDao.DeleteOrderById(req.OrderId, uId)
	if err != nil {
		logging.Info(err)
		code = e.ErrorDatabase
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}, err
	}

	return serializer.Response{
		Status: code,
		Msg:    e.GetMsg(code),
	}, nil
}
