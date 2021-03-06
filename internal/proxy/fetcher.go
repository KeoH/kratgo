package proxy

import (
	"sync"

	"github.com/valyala/fasthttp"
)

var setCookieHeader = []byte("Set-Cookie")

var statusRedirect = []int{
	fasthttp.StatusMovedPermanently,
	fasthttp.StatusFound,
	fasthttp.StatusSeeOther,
	fasthttp.StatusTemporaryRedirect,
	fasthttp.StatusPermanentRedirect,
}

// Hop-by-hop headers. These are removed when sent to the backend.
// http://www.w3.org/Protocols/rfc2616/rfc2616-sec13.html
var hopHeaders = []string{
	"Connection",
	"Keep-Alive",
	"Proxy-Authenticate",
	"Proxy-Authorization",
	"Te", // canonicalized version of "TE"
	"Trailers",
	"Transfer-Encoding",
	"Upgrade",
}

var fetcherPool = sync.Pool{
	New: func() interface{} {
		return &fetcher{
			req:  fasthttp.AcquireRequest(),
			resp: fasthttp.AcquireResponse(),
		}
	},
}

func acquireFetcher() *fetcher {
	return fetcherPool.Get().(*fetcher)
}

func releaseFetcher(f *fetcher) {
	f.reset()
	fetcherPool.Put(f)
}

func (f *fetcher) reset() {
	f.executeHeaderRule = false

	f.req.Reset()
	f.resp.Reset()
}

func (f *fetcher) Do(hostClient *fasthttp.HostClient) error {
	return hostClient.Do(f.req, f.resp)
}

func (f *fetcher) processHeaderRules(rules []HeaderRule, params *evalParams) error {
	for _, r := range rules {
		params.reset()

		f.executeHeaderRule = true

		if r.expr != nil {
			for _, p := range r.params {
				params.set(p.name, getEvalValue(f.req, f.resp, p.name, p.subKey))
			}

			result, err := r.expr.Evaluate(params.p)
			if err != nil {
				return err
			}

			f.executeHeaderRule = result.(bool)
		}

		if !f.executeHeaderRule {
			continue
		}

		if r.action == setHeaderAction {
			f.resp.Header.Set(r.name, getEvalValue(f.req, f.resp, r.value.value, r.value.subKey))
		} else {
			f.resp.Header.Del(r.name)
		}
	}

	return nil
}
