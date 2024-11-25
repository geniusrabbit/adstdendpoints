//go:build htmltemplates
// +build htmltemplates

//
// @project GeniusRabbit adstdendpoints 2022
// @author Dmitry Ponomarev <demdxx@gmail.com> 2022
//

package proxy

import (
	"github.com/geniusrabbit/adcorelib/adtype"

	"github.com/geniusrabbit/adstdendpoints"
	"github.com/geniusrabbit/adstdendpoints/templates"
)

type _endpoint struct{}

func New() *_endpoint { return &_endpoint{} }

func (e *_endpoint) Codename() string {
	return "proxy"
}

func (e *_endpoint) Handle(source adstdendpoints.Source, request *adtype.BidRequest) adtype.Responser {
	request.RequestCtx.SetContentType("text/html; charset=UTF-8")
	templates.WriteAdRenderDinamicProxyBanner(request.RequestCtx, request)
	return nil
}
