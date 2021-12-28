package outboundgroup

import (
	C "github.com/Dreamacro/clash/constant"
	"regexp"
	"sync"
)

type Filter struct {
	mu          sync.Mutex
	IsFilterSet bool
	regExp      *regexp.Regexp
}

func (f *Filter) FilterProxiesByName(proxies []C.Proxy) []C.Proxy {
	filteredProxies := []C.Proxy{}
	for _, proxy := range proxies {
		if f.regExp.MatchString(proxy.Name()) {
			filteredProxies = append(filteredProxies, proxy)
		}
	}
	return filteredProxies
}

func (f *Filter) SetExpr(expr string) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	reg, err := regexp.Compile(expr)
	if err != nil {
		f.IsFilterSet = false
		return err
	}
	f.IsFilterSet = true
	f.regExp = reg
	return nil
}
