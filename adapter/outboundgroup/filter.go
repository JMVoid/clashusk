package outboundgroup

import (
	C "github.com/Dreamacro/clash/constant"
	"regexp"
)

type Filter struct {
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
	reg, err := regexp.Compile(expr)
	if err != nil {
		f.IsFilterSet = false
		return err
	}
	f.IsFilterSet = true
	f.regExp = reg
	return nil
}
