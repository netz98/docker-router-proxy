package proxy

import (
	"../routing"
	"net/http"
	"log"
	"net/http/httputil"
)

type ProxyRouter struct {
	Debug bool
	Cache *routing.Cache
	Domain string
}

func (p *ProxyRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("X-Forwarded-For", r.RemoteAddr)
	w.Header().Set("X-Forwarded-Host", r.Host)
	w.Header().Set("X-Forwarded-Server", "localhost")

	targetUrl := routing.ResolveTargetContainer(r, p.Cache, p.Debug, p.Domain)
	if targetUrl.String() == "" {
		// error
		if p.Debug {
			log.Println("Error resolving container for request:", r.Host)
		}
	} else {
		// proxy
		httputil.NewSingleHostReverseProxy(targetUrl).ServeHTTP(w, r)
	}
}
