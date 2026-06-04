//
// @project GeniusRabbit adstdendpoints 2018 - 2026
// @author Dmitry Ponomarev <demdxx@gmail.com> 2018 - 2026
//

package dynamic

import (
	"encoding/json"
	"strings"

	"github.com/demdxx/gocast/v2"
	"github.com/valyala/fasthttp"

	"github.com/geniusrabbit/adcorelib/admodels"
	"github.com/geniusrabbit/adcorelib/admodels/types"
	"github.com/geniusrabbit/adcorelib/adquery/bidresponse"
	"github.com/geniusrabbit/adcorelib/adtype"
	"github.com/geniusrabbit/adcorelib/eventtraking/events"
	"github.com/geniusrabbit/adcorelib/httpserver/extensions/endpoint"
)

// Endpoint is a dynamic endpoint
type _endpoint struct {
	urlGen   adtype.URLGenerator
	metaConf MetaConfig
}

// New creates new dynamic endpoint
func New(urlGen adtype.URLGenerator, metaConf MetaConfig) *_endpoint {
	return &_endpoint{urlGen: urlGen, metaConf: metaConf}
}

// Codename of the endpoint
func (e *_endpoint) Codename() string {
	return "dynamic"
}

// Handle request of the dynamic Ad and return response
func (e _endpoint) Handle(source endpoint.Source, request adtype.BidRequester) (response adtype.Response) {
	if request.IsRobot() {
		response = bidresponse.NewEmptyResponse(request, nil, nil)
		_ = e.renderEmpty(request.HTTPRequest(), response)
	} else {
		response = source.Bid(request)
		if err := e.render(request.HTTPRequest(), response); err != nil {
			response = adtype.NewErrorResponse(request, err)
		}
	}
	return response
}

func (e _endpoint) render(ctx *fasthttp.RequestCtx, response adtype.Response) error {
	resp := Response{Version: "1"}

	if response.Request().IsDebug() {
		headers := map[string]string{}
		for key, value := range ctx.Request.Header.All() {
			headers[string(key)] = string(value)
		}
		resp.Debug = map[string]any{
			"http": map[string]any{
				"uri":     string(ctx.RequestURI()),
				"ip":      ctx.RemoteIP().String(),
				"method":  string(ctx.Method()),
				"query":   ctx.QueryArgs().String(),
				"headers": headers,
			},
		}
	}

	// Process response ad items
	for _, ad := range response.Ads() {
		var (
			assets       []adAsset
			aditm        = ad.(adtype.ResponseItem)
			cPrep        = adtype.ContentPreparer(response, aditm)
			url          string
			trackerBlock = &tracker{}
		)

		// Generate click URL
		if !aditm.Format().IsProxy() && !aditm.Format().IsDirect() {
			url, _ = e.urlGen.ClickURL(aditm, response)
		} else if aditm.Format().IsDirect() && !aditm.Impression().IsInterstitial() {
			url, _ = e.urlGen.DirectURL(events.Direct, aditm, response)
		}

		// Generate no-error impression and view tracking pixels for interstitial and non-direct formats
		if aditm.Impression().IsInterstitial() || !aditm.Format().IsDirect() {
			if !aditm.Format().IsDirect() {
				// If direct format, then the real impression will be tracked by the direct URL, so no need to add impression pixel here
				trackerBlock.Impressions = append(
					trackerBlock.Impressions,
					e.noErrorPixelURL(events.Impression, events.StatusSuccess,
						aditm.Impression(), aditm, response, false),
				)
			}
			trackerBlock.Views = append(trackerBlock.Views,
				e.noErrorPixelURL(events.View, events.StatusSuccess,
					aditm.Impression(), aditm, response, false))
		}

		// Third-party trackers pixels
		trackerBlock.Clicks = listContentPrepare(aditm.ClickTrackerLinks(), cPrep)
		if links := aditm.ViewTrackerLinks(); len(links) > 0 {
			trackerBlock.Views = append(trackerBlock.Views, listContentPrepare(links, cPrep)...)
		}
		if links := aditm.ImpressionTrackerLinks(); len(links) > 0 {
			trackerBlock.Impressions = append(trackerBlock.Impressions, listContentPrepare(links, cPrep)...)
		}

		// Process assets if provided in the ad item. This includes generating CDN URLs for asset paths and preparing thumbnails. The assets are collected into a slice of adAsset structs, which will be included in the response item definition.
		if baseAssets := aditm.Assets(); len(baseAssets) > 0 {
			assets = make([]adAsset, 0, len(baseAssets))
			processed := map[string]int{}
			for _, as := range baseAssets {
				if as.URL == "" {
					continue
				}
				if idx, ok := processed[as.Name]; !ok {
					nas := adAsset{
						Name:     as.Name,
						Path:     cPrep.Replace(e.urlGen.CDNURL(as.URL)),
						Type:     as.Type.Code(),
						Width:    as.Width,
						Height:   as.Height,
						Duration: as.Duration,
						Thumbs:   e.thumbsPrepare(as.Thumbs, cPrep),
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

		// Determine ad type for response item
		adType := ad.PriorityFormatType()

		// For non-direct formats, check for IFrame URL or HTML content in the ad item content fields and add as asset if available. For direct interstitial formats, generate direct URL and add as iframe_url asset.
		if adFormat := aditm.Format(); !adFormat.IsDirect() {
			if contentURL := aditm.ContentItemString(adtype.ContentItemIFrameURL); contentURL != "" {
				assets = append(assets, adAsset{
					Name:   "main",
					Type:   "iframe_url",
					Path:   cPrep.Replace(contentURL),
					Width:  adFormat.Width,
					Height: adFormat.Height,
				})
			} else if content := aditm.ContentItemString(adtype.ContentItemContent); content != "" {
				assets = append(assets, adAsset{
					Name:   "main",
					Type:   "html",
					Path:   cPrep.Replace(content),
					Width:  adFormat.Width,
					Height: adFormat.Height,
				})
			}
		} else if aditm.Impression().IsInterstitial() {
			adType = types.FormatProxyType
			directURL, _ := e.urlGen.DirectURL(events.Direct, aditm, response)
			assets = append(assets, adAsset{
				Name: "main",
				Type: "iframe_url",
				Path: directURL,
			})
		}

		// Ad item definition for response. It includes the ad item ID, type, click URL, assets, and tracking pixels. The ad type is determined based on the priority format type of the ad item, with special handling for interstitial direct formats to set the type as "proxy" and include the direct URL as an iframe asset.
		adItemObj := &item{
			ID:      ad.ID(),
			Type:    adType.Name(),
			URL:     url,
			Fields:  noEmptyFieldsMap(aditm.ContentFields(), cPrep),
			Assets:  assets,
			Tracker: trackerBlock,
			AdInfo:  e.prepareItemAdInfo(aditm, response, cPrep),
			Debug: gocast.IfThenExec(response.Request().IsDebug(),
				func() any { return map[string]any{"adUnit": ad} },
				func() any { return nil }),
		}

		// Add item to response group by impression ID
		resp.getGroupOrCreate(ad.TargetCodename()).addItem(adItemObj)

		// Add source info
		e.addSourceInfo(&resp, aditm.Source())
	}

	// Add empty group tracking if no items
	req := response.Request()
	for _, imp := range req.Impressions() {
		group := resp.getGroupOrCreate(imp.TargetCodename())
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

func (e _endpoint) addSourceInfo(resp *Response, source adtype.Source) {
	if source == nil {
		return
	}
	sourceID := gocast.Str(source.ID())
	if resp.hasSource(sourceID) {
		return
	}
	info := source.Info()
	if info == nil {
		info = &adtype.SourceInfo{ID: sourceID, Protocol: source.Protocol()}
	}
	resp.AdSources = append(resp.AdSources, &adSourceInfo{
		ID:          sourceID,
		Name:        info.Name,
		Description: info.Description,
		Domain:      info.Domain,
		IconURL:     info.IconURL,
		LogoURL:     info.LogoURL,
		URL:         info.URL,
		Metadata:    info.Metadata,
	})
}

func (e _endpoint) prepareItemAdInfo(item adtype.ResponseItem, _ adtype.Response, cPrep *strings.Replacer) *itemMetaInfo {
	sourceID := ""
	if source := item.Source(); source != nil {
		sourceID = u64ID2Str(source.ID())
	}
	adInfo := &itemMetaInfo{
		Advertiser: &itemMetaAdvertiserInfo{
			ID: u64ID2Str(item.AccountID()),
		},
		Ad: &itemMetaAdInfo{
			ID:         item.AdID(),
			CampaignID: u64ID2Str(item.CampaignID()),
			AdSourceID: sourceID,
		},
	}
	if e.metaConf.ComplaintAdURL != "" || e.metaConf.AboutAdURL != "" {
		if e.metaConf.AboutAdURL != "" {
			adInfo.Ad.AboutURL = cPrep.Replace(e.metaConf.AboutAdURL)
		}
		if e.metaConf.ComplaintAdURL != "" {
			adInfo.addAction("complaint", "Report this Ad", "", cPrep.Replace(e.metaConf.ComplaintAdURL))
		}
	}
	return adInfo
}

func (e _endpoint) renderEmpty(ctx *fasthttp.RequestCtx, response adtype.Response) error {
	resp := Response{Version: "1"}

	// Add empty group tracking
	req := response.Request()
	for _, imp := range req.Impressions() {
		group := resp.getGroupOrCreate(imp.TargetCodename())
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

func (e _endpoint) thumbsPrepare(thumbs []admodels.AdFileAssetThumb, cPrep *strings.Replacer) []adAssetThumb {
	nthumbs := make([]adAssetThumb, 0, len(thumbs))
	for _, th := range thumbs {
		nthumbs = append(nthumbs, adAssetThumb{
			Path:   cPrep.Replace(e.urlGen.CDNURL(th.URL)),
			Type:   th.Type.Code(),
			Width:  th.Width,
			Height: th.Height,
		})
	}
	return nthumbs
}

func (e _endpoint) noErrorPixelURL(event events.Type, status uint8, imp *adtype.Impression, item adtype.ResponseItem, response adtype.Response, js bool) string {
	if item == nil {
		if imp == nil {
			imp = &adtype.Impression{Target: &adtype.TargetEmpty{}}
		}
		formats := response.Request().Formats().List()
		item = &bidresponse.ResponseItemBlank{
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

//go:inline
func listContentPrepare(arr []string, prep *strings.Replacer) []string {
	if prep == nil {
		return arr
	}
	prepared := make([]string, 0, len(arr))
	for _, s := range arr {
		prepared = append(prepared, prep.Replace(s))
	}
	return prepared
}

func noEmptyFieldsMap(m map[string]any, prep *strings.Replacer) map[string]any {
	if len(m) == 0 {
		return nil
	}
	for k, v := range m {
		switch val := v.(type) {
		case string:
			if val == "" {
				delete(m, k)
			} else if prep != nil {
				m[k] = prep.Replace(val)
			}
		case []string:
			if len(val) == 0 {
				delete(m, k)
			} else if prep != nil {
				m[k] = listContentPrepare(val, prep)
			}
		case nil:
			delete(m, k)
		}
	}
	return m
}

//go:inline
func u64ID2Str(id uint64) string {
	if id == 0 {
		return ""
	}
	return gocast.Str(id)
}
