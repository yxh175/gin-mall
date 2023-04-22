package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"

	util "mall/pkg/utils"
	"mall/service"
	"mall/types"
)

func ImportSkillProductHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req types.SkillProductImportReq

		if err := ctx.ShouldBind(&req); err == nil {
			// 参数校验
			file, _, _ := ctx.Request.FormFile("file")
			l := service.GetSkillProductSrv()
			resp, err := l.Import(ctx.Request.Context(), file)
			if err != nil {
				util.LogrusObj.Infoln(err)
				ctx.JSON(http.StatusInternalServerError, ErrorResponse(err))
				return
			}
			ctx.JSON(http.StatusOK, resp)
		} else {
			util.LogrusObj.Infoln(err)
			ctx.JSON(http.StatusBadRequest, ErrorResponse(err))
		}
	}
}

func InitSkillProductHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req types.SkillProductImportReq

		if err := ctx.ShouldBind(&req); err == nil {
			// 参数校验
			l := service.GetSkillProductSrv()
			resp, err := l.InitSkillGoods(ctx.Request.Context())
			if err != nil {
				util.LogrusObj.Infoln(err)
				ctx.JSON(http.StatusInternalServerError, ErrorResponse(err))
				return
			}
			ctx.JSON(http.StatusOK, resp)
		} else {
			util.LogrusObj.Infoln(err)
			ctx.JSON(http.StatusBadRequest, ErrorResponse(err))
		}
	}
}

func SkillProductHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req types.SkillProductServiceReq

		if err := ctx.ShouldBind(&req); err == nil {
			// 参数校验
			userId := ctx.Keys["user_id"].(uint)
			l := service.GetSkillProductSrv()
			resp, err := l.SkillProduct(ctx.Request.Context(), userId, &req)
			if err != nil {
				util.LogrusObj.Infoln(err)
				ctx.JSON(http.StatusInternalServerError, ErrorResponse(err))
				return
			}
			ctx.JSON(http.StatusOK, resp)
		} else {
			util.LogrusObj.Infoln(err)
			ctx.JSON(http.StatusBadRequest, ErrorResponse(err))
		}
	}
}
