//
// @project GeniusRabbit sspserver 2018 - 2019, 2025
// @author Dmitry Ponomarev <demdxx@gmail.com> 2018 - 2019, 2025
//

package direct

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/opentracing/opentracing-go/ext"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"

	"github.com/geniusrabbit/adcorelib/admodels/types"
	"github.com/geniusrabbit/adcorelib/adtype"
	"github.com/geniusrabbit/adcorelib/context/ctxlogger"
	"github.com/geniusrabbit/adcorelib/eventtraking/events"
	"github.com/geniusrabbit/adcorelib/eventtraking/eventstream"
	"github.com/geniusrabbit/adcorelib/gtracing"
	"github.com/geniusrabbit/adstdendpoints"
)

// Error list...
var (
	ErrMultipleDirectNotSupported = errors.New("direct: multiple direct responses not supported")
	ErrInvalidResponseType        = errors.New("direct: invalid response type")
)

type _endpoint struct {
	formats          types.FormatsAccessor
	superFailoverURL string
}

func New(formats types.FormatsAccessor, superFailoverURL string) *_endpoint {
	return &_endpoint{
		formats:          formats,
		superFailoverURL: superFailoverURL,
	}
}

func (e *_endpoint) Codename() string {
	return "direct"
}

// Handle processes the bid request and sends the appropriate direct response.
func (e *_endpoint) Handle(source adstdendpoints.Source, request adtype.BidRequester) adtype.Response {
	request.ImpressionUpdate(func(imp *adtype.Impression) bool {
		imp.Width, imp.Height = -1, -1
		imp.FormatTypes.Reset().Set(types.FormatDirectType)
		return true
	})
	newRequest := request.WithFormats(e.formats)
	response := source.Bid(newRequest)
	if err := e.execDirect(newRequest.HTTPRequest(), response); err != nil {
		ctxlogger.Get(newRequest.Context()).Error("exec direct", zap.Error(err))
	} else {
		e.sendViewEvent(response)
	}
	return response
}

func (e *_endpoint) execDirect(req *fasthttp.RequestCtx, response adtype.Response) (err error) {
	var (
		id              string
		zoneID          uint64
		impID           string
		link            string
		alternativeLink = false
	)

	if span, _ := gtracing.StartSpanFromFastContext(req, "render"); span != nil {
		ext.Component.Set(span, "endpoint.direct")
		defer span.Finish()
	}

	if response == nil || response.Count() == 0 {
		if response != nil {
			if imps := response.Request().Impressions(); len(imps) > 0 {
				impID = imps[0].ID
				if imps[0].Target != nil {
					link = imps[0].Target.AlternativeAdCode("direct")
					zoneID = uint64(imps[0].TargetID())
					alternativeLink = link != ""
				}
			}
		}
	} else if err = response.Validate(); err == nil {
		if response.Count() > 1 {
			err = ErrMultipleDirectNotSupported
		} else {
			adv := response.Ads()[0]
			impID = adv.ImpressionID()
			if adv.Impression() != nil {
				zoneID = uint64(adv.Impression().TargetID())
			}

			switch ad := adv.(type) {
			case adtype.ResponseItem:
				id = ad.AdID()
				if !ad.IsDirect() {
					err = ErrInvalidResponseType
				} else {
					link = adtype.PrepareURL(ad.ActionURL(), response, ad)
				}
			case adtype.ResponseMultipleItem:
				err = ErrMultipleDirectNotSupported
			default:
				// ...
			}
		}
	}

	switch {
	case response != nil && response.Request().IsDebug() && req.QueryArgs().Has("noredirect"):
		req.SetStatusCode(http.StatusOK)
		req.SetContentType("application/json")
		_ = json.NewEncoder(req).Encode(debugResponse{
			ID:                id,
			ZoneID:            zoneID,
			ImpressionID:      impID,
			AuctionID:         response.Request().AuctionID(),
			IsAlternativeLink: alternativeLink,
			Link:              link,
			Superfailover:     e.superFailoverURL,
			Error:             err,
			IsEmpty:           response.Count() < 1,
		})
	case link != "":
		req.Response.Header.Set("X-Status-Alternative", "1")
		req.Redirect(link, http.StatusFound)
	case e.superFailoverURL == "":
		req.Success("text/plain", []byte("Please add superfailover link"))
	default:
		req.Response.Header.Set("X-Status-Failover", "1")
		req.Redirect(e.superFailoverURL, http.StatusFound)
	}
	return err
}

func (e *_endpoint) sendViewEvent(response adtype.Response) {
	if response == nil || response.Error() != nil || len(response.Ads()) == 0 {
		return
	}
	if response.Request().IsDebug() && response.Request().HTTPRequest().QueryArgs().Has("noredirect") {
		ctxlogger.Get(response.Context()).Info("skip event log",
			zap.String("request_id", response.Request().ID()))
		return
	}
	var (
		err    error
		stream = eventstream.StreamFromContext(response.Context())
	)

	switch ad := response.Ads()[0].(type) {
	case adtype.ResponseItem:
		err = stream.Send(events.Direct, events.StatusSuccess, response, ad)
	case adtype.ResponseMultipleItem:
		if len(ad.Ads()) > 0 {
			err = stream.Send(events.Direct, events.StatusSuccess, response, ad.Ads()[0])
		}
	default:
		// Invalid ad type
	}
	if err != nil {
		ctxlogger.Get(response.Context()).Error("send direct event", zap.Error(err))
	}
}
