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

type bidsImpl struct {
	r *Router
}

func registerBidsApi(r *Router) {
	bidsImpl := bidsImpl{r}

	r.router.GET("/api/bids/{tenderId}/list", bidsImpl.bidList)
	r.router.POST("/api/bids/new", bidsImpl.createBid)
	r.router.GET("/api/bids/my", bidsImpl.userBids)
	r.router.GET("/api/bids/{bidId}/status", bidsImpl.getBidStatus)
	r.router.PUT("/api/bids/{bidId}/status", bidsImpl.setBidStatus)
	r.router.PATCH("/api/bids/{bidId}/edit", bidsImpl.editBid)
	r.router.PUT("/api/bids/{bidId}/rollback/{version}", bidsImpl.rollback)
	r.router.PUT("/api/bids/{bidId}/submit_decision", bidsImpl.submitDecision)
	r.router.PUT("/api/bids/{bidId}/feedback", bidsImpl.addFeedback)
}

func (t *bidsImpl) createBid(ctx *fasthttp.RequestCtx) {
	builder := models.BidsBuilder{}
	err := json.Unmarshal(ctx.Request.Body(), &builder)
	if err != nil {
		t.r.logger.Println(err)
		ctx.SetStatusCode(400)
		return
	}
	if builder.Name == "" || builder.Description == "" || builder.TenderId == "" ||
		builder.AuthorType == nil || builder.AuthorId == "" {
		t.r.logger.Println("invalid json")
		ctx.SetStatusCode(400)
		return
	}
	bid, err := t.r.billingService.CreateBid(builder)
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

	responce, err := json.Marshal(&bid)
	if err != nil {
		t.r.logger.Println(err)
		ctx.SetStatusCode(400)
		return
	}
	ctx.SetBody(responce)
}

func (t *bidsImpl) editBid(ctx *fasthttp.RequestCtx) {
	bidId := ctx.UserValue("bidId").(string)
	username := string(ctx.QueryArgs().Peek("username"))
	builder := models.BidsBuilder{}
	err := json.Unmarshal(ctx.Request.Body(), &builder)
	if err != nil {
		t.r.logger.Println(err)
		ctx.SetStatusCode(400)
		return
	}
	bid, err := t.r.billingService.UpdateBid(bidId, username, builder)
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
	responce, err := json.Marshal(&bid)
	if err != nil {
		t.r.logger.Println(err)
		ctx.SetStatusCode(400)
		return
	}
	ctx.SetBody(responce)
}

func (t *bidsImpl) getBidStatus(ctx *fasthttp.RequestCtx) {
	bidId := ctx.UserValue("bidId").(string)
	username := string(ctx.QueryArgs().Peek("username"))

	bid, err := t.r.billingService.GetBid(bidId, username)
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

	ctx.SetBody([]byte(bid.Status.String()))
}

func (t *bidsImpl) setBidStatus(ctx *fasthttp.RequestCtx) {
	bidId := ctx.UserValue("bidId").(string)
	username := string(ctx.QueryArgs().Peek("username"))
	status := string(ctx.QueryArgs().Peek("status"))
	statusModel, err := convertion.BidStatusFromString(status)
	if err != nil {
		t.r.logger.Println(err)
		ctx.SetStatusCode(400)
		return
	}
	bid, err := t.r.billingService.SetBidStatus(bidId, username, statusModel)
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
	responce, err := json.Marshal(&bid)
	if err != nil {
		t.r.logger.Println(err)
		ctx.SetStatusCode(400)
		return
	}
	ctx.SetBody(responce)
}

func (t *bidsImpl) rollback(ctx *fasthttp.RequestCtx) {
	bidId := ctx.UserValue("bidId").(string)
	versionString := ctx.UserValue("version").(string)
	version, err := strconv.Atoi(versionString)
	if err != nil {
		t.r.logger.Println("invalid version value")
		ctx.SetStatusCode(400)
		return
	}
	username := string(ctx.QueryArgs().Peek("username"))

	bid, err := t.r.billingService.RollbackBid(bidId, username, version)
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
	responce, err := json.Marshal(&bid)
	if err != nil {
		t.r.logger.Println(err)
		ctx.SetStatusCode(400)
		return
	}
	ctx.SetBody(responce)
}

func (t *bidsImpl) bidList(ctx *fasthttp.RequestCtx) {
	tenderId := ctx.UserValue("tenderId").(string)
	username := string(ctx.QueryArgs().Peek("username"))
	limit, err := ctx.QueryArgs().GetUint("limit")
	if err != nil {
		limit = 5
	}
	offset, err := ctx.QueryArgs().GetUint("offset")
	if err != nil {
		offset = 0
	}

	bids, err := t.r.billingService.GetBidsForTender(tenderId, username, limit, offset)
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
	responce, err := json.Marshal(&bids)
	if err != nil {
		t.r.logger.Println(err)
		ctx.SetStatusCode(400)
		return
	}
	ctx.SetBody(responce)
}

func (t *bidsImpl) userBids(ctx *fasthttp.RequestCtx) {
	username := string(ctx.QueryArgs().Peek("username"))
	limit, err := ctx.QueryArgs().GetUint("limit")
	if err != nil {
		limit = 5
	}
	offset, err := ctx.QueryArgs().GetUint("offset")
	if err != nil {
		offset = 0
	}

	bids, err := t.r.billingService.GetUserBids(limit, offset, username)
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
	responce, err := json.Marshal(&bids)
	if err != nil {
		t.r.logger.Println(err)
		ctx.SetStatusCode(400)
		return
	}
	ctx.SetBody(responce)
}

func (t *bidsImpl) submitDecision(ctx *fasthttp.RequestCtx) {
	bidId := ctx.UserValue("bidId").(string)
	username := string(ctx.QueryArgs().Peek("username"))
	decision := string(ctx.QueryArgs().Peek("decision"))
	decisionModel, err := convertion.BidDecisionFromString(decision)
	if err != nil {
		t.r.logger.Println(err)
		ctx.SetStatusCode(400)
		return
	}
	bid, err := t.r.billingService.SubmitDecision(bidId, username, decisionModel)
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
	responce, err := json.Marshal(&bid)
	if err != nil {
		t.r.logger.Println(err)
		ctx.SetStatusCode(400)
		return
	}
	ctx.SetBody(responce)
}

func (t *bidsImpl) addFeedback(ctx *fasthttp.RequestCtx) {
	bidId := ctx.UserValue("bidId").(string)
	username := string(ctx.QueryArgs().Peek("username"))
	bidFeedback := string(ctx.QueryArgs().Peek("bidFeedback"))
	if bidFeedback == "" {
		t.r.logger.Println("no feedback provided")
		ctx.Response.SetBody([]byte("no feedback provided"))
		ctx.SetStatusCode(400)
		return
	}
	bid, err := t.r.billingService.AddFeedback(bidId, username, bidFeedback)
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
	responce, err := json.Marshal(&bid)
	if err != nil {
		t.r.logger.Println(err)
		ctx.SetStatusCode(400)
		return
	}
	ctx.SetBody(responce)
}
