package outboundgroup

import (
	"time"

	C "github.com/Dreamacro/clash/constant"
	types "github.com/Dreamacro/clash/constant/provider"
)

const (
	defaultGetProxiesDuration = time.Second * 5
)

func getProvidersProxies(providers []types.ProxyProvider, touch bool, useFilter *Filter) []C.Proxy {
	var comProxies []C.Proxy
	var setProxies []C.Proxy
	for _, provider := range providers {
		if touch {
			if provider.VehicleType() == types.Compatible {
				comProxies = append(comProxies, provider.ProxiesWithTouch()...)
			} else {
				setProxies = append(setProxies, provider.ProxiesWithTouch()...)
			}
		} else {
			if provider.VehicleType() == types.Compatible {
				comProxies = append(comProxies, provider.Proxies()...)
			} else {
				setProxies = append(setProxies, provider.Proxies()...)
			}
		}
	}
	if useFilter.IsFilterSet {
		return append(comProxies, useFilter.FilterProxiesByName(setProxies)...)
	}
	return append(comProxies, setProxies...)
}
