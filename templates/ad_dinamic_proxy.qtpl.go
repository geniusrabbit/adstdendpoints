// Code generated by qtc from "ad_dinamic_proxy.qtpl". DO NOT EDIT.
// See https://github.com/valyala/quicktemplate for details.

//line private/templates/ad_dinamic_proxy.qtpl:2
package templates

//line private/templates/ad_dinamic_proxy.qtpl:2
import (
	"github.com/geniusrabbit/adcorelib/adtype"
)

//line private/templates/ad_dinamic_proxy.qtpl:7
import (
	qtio422016 "io"

	qt422016 "github.com/valyala/quicktemplate"
)

//line private/templates/ad_dinamic_proxy.qtpl:7
var (
	_ = qtio422016.Copy
	_ = qt422016.AcquireByteBuffer
)

//line private/templates/ad_dinamic_proxy.qtpl:7
func StreamAdRenderDinamicProxyBanner(qw422016 *qt422016.Writer, request *adtype.BidRequest) {
//line private/templates/ad_dinamic_proxy.qtpl:8
	streamadHeader(qw422016)
//line private/templates/ad_dinamic_proxy.qtpl:9
	streampreloader(qw422016)
//line private/templates/ad_dinamic_proxy.qtpl:10
	streamadRenderNativeCSS(qw422016)
//line private/templates/ad_dinamic_proxy.qtpl:11
	var script = URLGen.LibURL("/embedded.js")

//line private/templates/ad_dinamic_proxy.qtpl:11
	qw422016.N().S(`<ins id="element_`)
//line private/templates/ad_dinamic_proxy.qtpl:12
	qw422016.N().D(int(request.TargetID()))
//line private/templates/ad_dinamic_proxy.qtpl:12
	qw422016.N().S(`"></ins><script type="text/javascript" src="`)
//line private/templates/ad_dinamic_proxy.qtpl:13
	qw422016.N().S(script)
//line private/templates/ad_dinamic_proxy.qtpl:13
	qw422016.N().S(`"></script><script type="text/javascript">!(function(){(new EmbeddedAd({`)
//line private/templates/ad_dinamic_proxy.qtpl:16
	if Debug {
//line private/templates/ad_dinamic_proxy.qtpl:16
		qw422016.N().S(`JSONPLink: '//`)
//line private/templates/ad_dinamic_proxy.qtpl:17
		qw422016.N().S(request.ServiceDomain())
//line private/templates/ad_dinamic_proxy.qtpl:17
		qw422016.N().S(`/b/dynamic/{<id>}?format=jsonp&',`)
//line private/templates/ad_dinamic_proxy.qtpl:17
	}
//line private/templates/ad_dinamic_proxy.qtpl:17
	qw422016.N().S(`element: "element_`)
//line private/templates/ad_dinamic_proxy.qtpl:18
	qw422016.N().D(int(request.TargetID()))
//line private/templates/ad_dinamic_proxy.qtpl:18
	qw422016.N().S(`",zone_id:`)
//line private/templates/ad_dinamic_proxy.qtpl:19
	qw422016.N().D(int(request.TargetID()))
//line private/templates/ad_dinamic_proxy.qtpl:19
	qw422016.N().S(`})).on('render', function() {var loader = window.document.getElementById('loadingBlock');if (loader) {loader.parentElement.removeChild(loader);}}).on('error', function(err) {console.log(err);}).render();})();</script>`)
//line private/templates/ad_dinamic_proxy.qtpl:30
	streamadFooter(qw422016)
//line private/templates/ad_dinamic_proxy.qtpl:31
}

//line private/templates/ad_dinamic_proxy.qtpl:31
func WriteAdRenderDinamicProxyBanner(qq422016 qtio422016.Writer, request *adtype.BidRequest) {
//line private/templates/ad_dinamic_proxy.qtpl:31
	qw422016 := qt422016.AcquireWriter(qq422016)
//line private/templates/ad_dinamic_proxy.qtpl:31
	StreamAdRenderDinamicProxyBanner(qw422016, request)
//line private/templates/ad_dinamic_proxy.qtpl:31
	qt422016.ReleaseWriter(qw422016)
//line private/templates/ad_dinamic_proxy.qtpl:31
}

//line private/templates/ad_dinamic_proxy.qtpl:31
func AdRenderDinamicProxyBanner(request *adtype.BidRequest) string {
//line private/templates/ad_dinamic_proxy.qtpl:31
	qb422016 := qt422016.AcquireByteBuffer()
//line private/templates/ad_dinamic_proxy.qtpl:31
	WriteAdRenderDinamicProxyBanner(qb422016, request)
//line private/templates/ad_dinamic_proxy.qtpl:31
	qs422016 := string(qb422016.B)
//line private/templates/ad_dinamic_proxy.qtpl:31
	qt422016.ReleaseByteBuffer(qb422016)
//line private/templates/ad_dinamic_proxy.qtpl:31
	return qs422016
//line private/templates/ad_dinamic_proxy.qtpl:31
}
