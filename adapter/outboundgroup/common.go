package outboundgroup

import (
	"time"

	C "github.com/Dreamacro/clash/constant"
	"github.com/Dreamacro/clash/constant/provider"
)

const (
	defaultGetProxiesDuration = time.Second * 5
)

func getProvidersProxies(providers []provider.ProxyProvider, touch bool, useFilter *Filter) []C.Proxy {
	proxies := []C.Proxy{}
	for _, provider := range providers {
		if touch {
			proxies = append(proxies, provider.ProxiesWithTouch()...)
		} else {
			proxies = append(proxies, provider.Proxies()...)
		}
	}
	if useFilter.IsFilterSet {
		return useFilter.FilterProxiesByName(proxies)
	}
	return proxies
}
