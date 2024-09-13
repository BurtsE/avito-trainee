package router

import (
	"avito-test/internal/models"
	"avito-test/internal/router/convertion"
	"avito-test/internal/service_errors"
	"encoding/json"
	"errors"
	"strconv"

	"github.com/valyala/fasthttp"
)

type tenderImpl struct {
	r *Router
}

func registerTenderApi(r *Router) {
	tenderImpl := tenderImpl{r}

	r.router.GET("/api/tenders", tenderImpl.tenderList)
	r.router.POST("/api/tenders/new", tenderImpl.createTender)
	r.router.GET("/api/tenders/my", tenderImpl.userTenders)
	r.router.GET("/api/tenders/{tenderId}/status", tenderImpl.getTenderStatus)
	r.router.PUT("/api/tenders/{tenderId}/status", tenderImpl.setTenderStatus)
	r.router.PATCH("/api/tenders/{tenderId}/edit", tenderImpl.editTender)
	r.router.PUT("/api/tenders/{tenderId}/rollback/{version}", tenderImpl.rollback)
}

func (t *tenderImpl) createTender(ctx *fasthttp.RequestCtx) {
	builder := models.TenderBuilder{}
	err := json.Unmarshal(ctx.Request.Body(), &builder)
	if err != nil {
		t.r.logger.Println(err)
		ctx.SetStatusCode(400)
		return
	}
	if builder.Name == "" || builder.Description == "" || builder.ServiceType == nil ||
		builder.OrganizationId == "" || builder.UserName == "" || builder.Status == nil {
		t.r.logger.Println("invalid json")
		ctx.SetStatusCode(400)
		return
	}
	tender, err := t.r.billingService.CreateTender(builder)
	if errors.As(err, &service_errors.AuthError{}) {
		ctx.Response.SetBody([]byte(err.Error()))
		ctx.SetStatusCode(403)
		return
	} else if errors.As(err, &service_errors.UserNotFound{}) {
		ctx.Response.SetBody([]byte(err.Error()))
		ctx.SetStatusCode(401)
		return
	} else if errors.As(err, &service_errors.TenderError{}) {
		ctx.Response.SetBody([]byte(err.Error()))
		ctx.SetStatusCode(404)
		return
	} else if err != nil {
		t.r.logger.Println(err.Error())
		ctx.Response.SetBody([]byte("internal server error"))
		ctx.SetStatusCode(400)
	}

	responce, err := json.Marshal(&tender)
	if err != nil {
		t.r.logger.Println(err)
		ctx.SetStatusCode(400)
		return
	}
	ctx.SetBody(responce)
}

func (t *tenderImpl) editTender(ctx *fasthttp.RequestCtx) {
	tenderId := ctx.UserValue("tenderId").(string)
	username := string(ctx.QueryArgs().Peek("username"))
	builder := models.TenderBuilder{}
	err := json.Unmarshal(ctx.Request.Body(), &builder)
	if err != nil {
		t.r.logger.Println(err)
		ctx.SetStatusCode(400)
		return
	}
	tender, err := t.r.billingService.UpdateTender(tenderId, username, builder)
	if errors.As(err, &service_errors.AuthError{}) {
		ctx.Response.SetBody([]byte(err.Error()))
		ctx.SetStatusCode(403)
		return
	} else if errors.As(err, &service_errors.UserNotFound{}) {
		ctx.Response.SetBody([]byte(err.Error()))
		ctx.SetStatusCode(401)
		return
	} else if errors.As(err, &service_errors.TenderError{}) {
		ctx.Response.SetBody([]byte(err.Error()))
		ctx.SetStatusCode(404)
		return
	} else if err != nil {
		t.r.logger.Println(err.Error())
		ctx.Response.SetBody([]byte("internal server error"))
		ctx.SetStatusCode(400)
	}
	responce, err := json.Marshal(&tender)
	if err != nil {
		t.r.logger.Println(err)
		ctx.SetStatusCode(400)
		return
	}
	ctx.SetBody(responce)
}

func (t *tenderImpl) getTenderStatus(ctx *fasthttp.RequestCtx) {
	tenderId := ctx.UserValue("tenderId").(string)
	username := string(ctx.QueryArgs().Peek("username"))

	tender, err := t.r.billingService.GetTender(tenderId, username)
	if errors.As(err, &service_errors.AuthError{}) {
		ctx.Response.SetBody([]byte(err.Error()))
		ctx.SetStatusCode(403)
		return
	} else if errors.As(err, &service_errors.UserNotFound{}) {
		ctx.Response.SetBody([]byte(err.Error()))
		ctx.SetStatusCode(401)
		return
	} else if errors.As(err, &service_errors.TenderError{}) {
		ctx.Response.SetBody([]byte(err.Error()))
		ctx.SetStatusCode(404)
		return
	} else if err != nil {
		t.r.logger.Println(err.Error())
		ctx.Response.SetBody([]byte("internal server error"))
		ctx.SetStatusCode(400)
	}

	ctx.SetBody([]byte(tender.Status.String()))
}

func (t *tenderImpl) setTenderStatus(ctx *fasthttp.RequestCtx) {
	tenderId := ctx.UserValue("tenderId").(string)
	username := string(ctx.QueryArgs().Peek("username"))
	status := string(ctx.QueryArgs().Peek("status"))
	statusModel, err := convertion.StatusFromString(status)
	if err != nil {
		t.r.logger.Println(err)
		ctx.SetStatusCode(400)
		return
	}
	tender, err := t.r.billingService.SetTenderStatus(tenderId, username, statusModel)
	if errors.As(err, &service_errors.AuthError{}) {
		ctx.Response.SetBody([]byte(err.Error()))
		ctx.SetStatusCode(403)
		return
	} else if errors.As(err, &service_errors.UserNotFound{}) {
		ctx.Response.SetBody([]byte(err.Error()))
		ctx.SetStatusCode(401)
		return
	} else if errors.As(err, &service_errors.TenderError{}) {
		ctx.Response.SetBody([]byte(err.Error()))
		ctx.SetStatusCode(404)
		return
	} else if err != nil {
		t.r.logger.Println(err.Error())
		ctx.Response.SetBody([]byte("internal server error"))
		ctx.SetStatusCode(400)
	}
	responce, err := json.Marshal(&tender)
	if err != nil {
		t.r.logger.Println(err)
		ctx.SetStatusCode(400)
		return
	}
	ctx.SetBody(responce)
}

func (t *tenderImpl) rollback(ctx *fasthttp.RequestCtx) {
	tenderId := ctx.UserValue("tenderId").(string)
	versionString := ctx.UserValue("version").(string)
	version, err := strconv.Atoi(versionString)
	if err != nil {
		t.r.logger.Println("invalid version value")
		ctx.SetStatusCode(400)
		return
	}
	username := string(ctx.QueryArgs().Peek("username"))

	tender, err := t.r.billingService.RollbackTender(tenderId, username, version)
	if errors.As(err, &service_errors.AuthError{}) {
		ctx.Response.SetBody([]byte(err.Error()))
		ctx.SetStatusCode(403)
		return
	} else if errors.As(err, &service_errors.UserNotFound{}) {
		ctx.Response.SetBody([]byte(err.Error()))
		ctx.SetStatusCode(401)
		return
	} else if errors.As(err, &service_errors.TenderError{}) {
		ctx.Response.SetBody([]byte(err.Error()))
		ctx.SetStatusCode(404)
		return
	} else if err != nil {
		t.r.logger.Println(err.Error())
		ctx.Response.SetBody([]byte("internal server error"))
		ctx.SetStatusCode(400)
		return
	}
	responce, err := json.Marshal(&tender)
	if err != nil {
		t.r.logger.Println(err)
		ctx.SetStatusCode(400)
		return
	}
	ctx.SetBody(responce)
}

func (t *tenderImpl) tenderList(ctx *fasthttp.RequestCtx) {
	var (
		serviceTypes []models.EnumService
	)
	limit, err := ctx.QueryArgs().GetUint("limit")
	if err != nil {
		limit = 5
	}
	offset, err := ctx.QueryArgs().GetUint("offset")
	if err != nil {
		offset = 0
	}

	serviceBytes := ctx.QueryArgs().PeekMultiBytes([]byte("service_type"))
	for _, value := range serviceBytes {
		service, err := convertion.ServiceFromString(string(value))
		if err != nil {
			t.r.logger.Println(err)
			ctx.SetStatusCode(400)
			ctx.Response.SetBody([]byte("Неверныe параметры запроса"))
			return
		}
		serviceTypes = append(serviceTypes, service)
	}
	if len(serviceTypes) == 0 {
		serviceTypes = []models.EnumService{models.Construction, models.Delivery, models.Manufacture}
	}
	tenders, err := t.r.billingService.GetTenders(limit, offset, serviceTypes)
	if errors.As(err, &service_errors.AuthError{}) {
		ctx.Response.SetBody([]byte(err.Error()))
		ctx.SetStatusCode(403)
		return
	} else if errors.As(err, &service_errors.UserNotFound{}) {
		ctx.Response.SetBody([]byte(err.Error()))
		ctx.SetStatusCode(401)
		return
	} else if errors.As(err, &service_errors.TenderError{}) {
		ctx.Response.SetBody([]byte(err.Error()))
		ctx.SetStatusCode(404)
		return
	} else if err != nil {
		t.r.logger.Println(err.Error())
		ctx.Response.SetBody([]byte("internal server error"))
		ctx.SetStatusCode(400)
	}
	responce, err := json.Marshal(&tenders)
	if err != nil {
		t.r.logger.Println(err)
		ctx.SetStatusCode(400)
		return
	}
	ctx.SetBody(responce)
}

func (t *tenderImpl) userTenders(ctx *fasthttp.RequestCtx) {
	username := string(ctx.QueryArgs().Peek("username"))
	limit, err := ctx.QueryArgs().GetUint("limit")
	if err != nil {
		limit = 5
	}
	offset, err := ctx.QueryArgs().GetUint("offset")
	if err != nil {
		offset = 0
	}

	tenders, err := t.r.billingService.GetUserTenders(limit, offset, username)
	if err != nil {
		t.r.logger.Println(err)
		ctx.SetStatusCode(400)
		return
	}
	responce, err := json.Marshal(&tenders)
	if err != nil {
		t.r.logger.Println(err)
		ctx.SetStatusCode(400)
		return
	}
	ctx.SetBody(responce)
}
