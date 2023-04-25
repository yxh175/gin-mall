package service

import (
	"context"
	"sync"

	util "mall/pkg/utils"
	"mall/pkg/utils/ctl"
	"mall/repository/db/dao"
	"mall/types"
)

var CarouselSrvIns *CarouselSrv
var CarouselSrvOnce sync.Once

type CarouselSrv struct {
}

func GetCarouselSrv() *CarouselSrv {
	CarouselSrvOnce.Do(func() {
		CarouselSrvIns = &CarouselSrv{}
	})
	return CarouselSrvIns
}

// ListCarousel 列表
func (s *CarouselSrv) ListCarousel(ctx context.Context, req *types.ListCarouselsServiceReq) (resp interface{}, err error) {
	carousels, err := dao.NewCarouselDao(ctx).ListCarousel()
	if err != nil {
		util.LogrusObj.Error(err)
		return
	}
	return ctl.RespSuccessWithData(carousels), nil
}
