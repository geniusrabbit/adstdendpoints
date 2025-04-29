//
// @project GeniusRabbit adstdendpoints 2018 - 2025
// @author Dmitry Ponomarev <demdxx@gmail.com> 2018 - 2025
//

package dynamic

import (
	"encoding/json"
	"math/rand"

	"github.com/demdxx/gocast/v2"
	"github.com/valyala/fasthttp"

	"github.com/geniusrabbit/adcorelib/admodels"
	"github.com/geniusrabbit/adcorelib/admodels/types"
	"github.com/geniusrabbit/adcorelib/adtype"
	"github.com/geniusrabbit/adcorelib/eventtraking/events"
	"github.com/geniusrabbit/adcorelib/httpserver/extensions/endpoint"
)

// Endpoint is a dynamic endpoint
type _endpoint struct {
	urlGen adtype.URLGenerator
}

// New creates new dynamic endpoint
func New(urlGen adtype.URLGenerator) *_endpoint {
	return &_endpoint{urlGen: urlGen}
}

// Codename of the endpoint
func (e *_endpoint) Codename() string {
	return "dynamic"
}

// Handle request of the dynamic Ad and return response
func (e _endpoint) Handle(source endpoint.Source, request *adtype.BidRequest) (response adtype.Responser) {
	if request.IsRobot() {
		response = adtype.NewEmptyResponse(request, nil, nil)
		_ = e.renderEmpty(request.RequestCtx, response)
	} else {
		response = source.Bid(request)
		if err := e.render(request.RequestCtx, response); err != nil {
			response = adtype.NewErrorResponse(request, err)
		}
	}
	return response
}

func (e _endpoint) render(ctx *fasthttp.RequestCtx, response adtype.Responser) error {
	resp := Response{Version: "1"}

	if response.Request().Debug {
		headers := map[string]string{}
		ctx.Request.Header.VisitAll(func(key, value []byte) {
			headers[string(key)] = string(value)
		})
		resp.Debug = map[string]any{
			"http": map[string]any{
				"uri":     string(ctx.RequestURI()),
				"ip":      string(ctx.RemoteIP()),
				"method":  string(ctx.Method()),
				"query":   ctx.QueryArgs().String(),
				"headers": headers,
			},
		}
	}

	// Process response ad items
	for _, ad := range response.Ads() {
		var (
			assets       []asset
			aditm        = ad.(adtype.ResponserItem)
			url          string
			trackerBlock tracker
		)

		// Generate click URL
		if !aditm.Format().IsProxy() {
			url, _ = e.urlGen.ClickURL(aditm, response)
		}

		trackerBlock = tracker{
			Impressions: []string{
				e.noErrorPixelURL(events.Impression, events.StatusSuccess, aditm.Impression(), aditm, response, false),
			},
			Views: []string{
				e.noErrorPixelURL(events.View, events.StatusSuccess, aditm.Impression(), aditm, response, false),
			},
		}

		// Third-party trackers pixels
		if item, _ := ad.(adtype.ResponserItem); item != nil {
			trackerBlock.Clicks = item.ClickTrackerLinks()
			if links := item.ViewTrackerLinks(); len(links) > 0 {
				trackerBlock.Views = append(trackerBlock.Views, links...)
			}
			if links := item.ImpressionTrackerLinks(); len(links) > 0 {
				trackerBlock.Impressions = append(trackerBlock.Impressions, links...)
			}
		}

		// Process assets if provided
		if baseAssets := aditm.Assets(); len(baseAssets) > 0 {
			assets = make([]asset, 0, len(baseAssets))
			processed := map[string]int{}
			for _, as := range baseAssets {
				if idx, ok := processed[as.Name]; !ok || rand.Float64() > 0.5 {
					nas := asset{
						Name:   as.Name,
						Path:   e.urlGen.CDNURL(as.Path),
						Type:   as.Type.Code(),
						Width:  as.Width,
						Height: as.Height,
						Thumbs: e.thumbsPrepare(as.Thumbs),
					}
					if !ok {
						processed[as.Name] = len(assets)
						assets = append(assets, nas)
					} else {
						assets[idx] = nas
					}
				}
			}
		}

		// Add item to response group by impression ID
		resp.getGroupOrCreate(ad.ImpressionID()).addItem(&item{
			ID:         ad.ID(),
			Type:       ad.PriorityFormatType().Name(),
			URL:        url,
			Content:    aditm.ContentItemString(adtype.ContentItemContent),
			ContentURL: aditm.ContentItemString(adtype.ContentItemIFrameURL),
			Fields:     noEmptyFieldsMap(aditm.ContentFields()),
			Assets:     assets,
			Tracker:    trackerBlock,
			Debug: gocast.IfThenExec(response.Request().Debug,
				func() any { return map[string]any{"adUnit": ad} },
				func() any { return nil }),
		})
	}

	// Add empty group tracking if no items
	req := response.Request()
	for i := range req.Imps {
		imp := &req.Imps[i]
		group := resp.getGroupOrCreate(imp.ID)
		if len(group.Items) == 0 {
			group.CustomTracker = tracker{
				Impressions: []string{
					e.noErrorPixelURL(events.Impression, events.StatusCustom, imp, nil, response, false),
				},
				Views: []string{
					e.noErrorPixelURL(events.View, events.StatusCustom, imp, nil, response, false),
				},
				Clicks: []string{
					e.noErrorPixelURL(events.Click, events.StatusCustom, imp, nil, response, false),
				},
			}
		}
	}

	// Render response to the client as JSONP
	format := string(ctx.QueryArgs().Peek("format"))
	if format == "jsonp" {
		callback := string(ctx.QueryArgs().Peek("callback"))
		if callback == "" {
			callback = "callback"
		}
		ctx.SetStatusCode(fasthttp.StatusOK)
		ctx.SetContentType("application/javascript")
		_, _ = ctx.Write([]byte(callback + "("))
		_ = json.NewEncoder(ctx).Encode(resp)
		_, _ = ctx.Write([]byte(")"))
		return nil
	}

	// Default JSON response
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetContentType("application/json")
	return json.NewEncoder(ctx).Encode(resp)
}

func (e _endpoint) renderEmpty(ctx *fasthttp.RequestCtx, response adtype.Responser) error {
	resp := Response{Version: "1"}

	// Add empty group tracking
	req := response.Request()
	for i := range req.Imps {
		imp := &req.Imps[i]
		group := resp.getGroupOrCreate(imp.ID)
		if len(group.Items) == 0 {
			group.CustomTracker = tracker{
				Impressions: []string{
					e.noErrorPixelURL(events.Impression, events.StatusCustom, imp, nil, response, false),
				},
				Views: []string{
					e.noErrorPixelURL(events.View, events.StatusCustom, imp, nil, response, false),
				},
				Clicks: []string{
					e.noErrorPixelURL(events.Click, events.StatusCustom, imp, nil, response, false),
				},
			}
		}
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetContentType("application/json")
	return json.NewEncoder(ctx).Encode(resp)
}

func (e _endpoint) thumbsPrepare(thumbs []admodels.AdAssetThumb) []assetThumb {
	nthumbs := make([]assetThumb, 0, len(thumbs))
	for _, th := range thumbs {
		nthumbs = append(nthumbs, assetThumb{
			Path:   e.urlGen.CDNURL(th.Path),
			Type:   th.Type.Code(),
			Width:  th.Width,
			Height: th.Height,
		})
	}
	return nthumbs
}

func (e _endpoint) noErrorPixelURL(event events.Type, status uint8, imp *adtype.Impression, item adtype.ResponserItem, response adtype.Responser, js bool) string {
	if item == nil {
		if imp == nil {
			imp = &adtype.Impression{Target: &admodels.Zone{}}
		}
		formats := response.Request().Formats()
		item = &adtype.ResponseItemBlank{
			Imp: imp,
			Src: &adtype.SourceEmpty{},
			FormatVal: gocast.IfThenExec(len(formats) > 0,
				func() *types.Format { return formats[0] },
				func() *types.Format { return &types.Format{} }),
		}
	}
	url, _ := e.urlGen.PixelURL(event, status, item, response, js)
	return url
}

func noEmptyFieldsMap(m map[string]any) map[string]any {
	if len(m) == 0 {
		return nil
	}
	for k, v := range m {
		switch val := v.(type) {
		case string:
			if val == "" {
				delete(m, k)
			}
		case []string:
			if len(val) == 0 {
				delete(m, k)
			}
		case nil:
			delete(m, k)
		}
	}
	return m
}
