package authmiddleware

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Auth interface {
	AuthorizeRequest(authList []string, fn func(ctx *gin.Context)) func(ctx *gin.Context)
	GetRequestId() string
	GetToken() string
	GetPlatform() string
}

type auth struct {
	authList []string
	c        *gin.Context
}

func Init() Auth {
	return &auth{}
}

func (a *auth) AuthorizeRequest(authList []string, fn func(ctx *gin.Context)) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		a.c = ctx
		if err := a.validate(ctx, authList...); err != nil {
			ctx.JSON(http.StatusOK, err.Error())
			return
		}

		if err := a.check(ctx); err != nil {
			ctx.JSON(http.StatusUnauthorized, err.Error())
			return
		}

		fn(ctx)
	}
}

func (a *auth) validate(ctx *gin.Context, conf ...string) error {
	m := map[string]string{
		RequestId: RequestId,
		Token:     Token,
		Platform:  Platform,
	}
	for _, c := range conf {
		if _, exists := m[c]; !exists {
			return errors.New("auth not found")
		}
	}
	a.authList = conf
	return nil
}

func (a *auth) check(ctx *gin.Context) error {
	for _, c := range a.authList {
		switch c {
		case RequestId:
			if err := a.requestId(ctx); err != nil {
				return err
			}
		case Token:
			if err := a.token(ctx); err != nil {
				return err
			}
		case Platform:
			if err := a.platform(ctx); err != nil {
				return err
			}
		}
	}
	return nil
}

func (a *auth) requestId(ctx *gin.Context) error {
	h := ctx.GetHeader(RequestId)
	if h == emptyValue {
		return errors.New(ErrorNotSendRequestId)
	}

	a.c.Header(RequestId, h)
	return nil
}

func (a *auth) token(ctx *gin.Context) error {
	h := ctx.GetHeader(Token)
	if h == emptyValue {
		return errors.New(ErrorNotSendUserId)
	}

	a.c.Header(Token, h)
	return nil
}

func (a *auth) platform(ctx *gin.Context) error {
	h := ctx.GetHeader(Platform)
	if h == emptyValue {
		return errors.New(ErrorNotSendPlatform)
	}

	a.c.Header(Platform, h)
	return nil
}

func (a *auth) GetRequestId() string {
	return a.c.GetHeader(RequestId)
}

func (a *auth) GetToken() string {
	return a.c.GetHeader(Token)
}

func (a *auth) GetPlatform() string {
	return a.c.GetHeader(Platform)
}
