package middleware

import (
	"net"
	"strings"

	"github.com/gin-gonic/gin"
)

var (
	// googleLoadBalancerCIDRs covers the public address ranges used by Google Front Ends.
	googleLoadBalancerCIDRs = []*net.IPNet{
		mustParseCIDR("35.191.0.0/16"),
		mustParseCIDR("130.211.0.0/22"),
	}
)

func mustParseCIDR(cidr string) *net.IPNet {
	_, network, err := net.ParseCIDR(cidr)
	if err != nil {
		panic(err)
	}
	return network
}

func parseCIDRList(entries []string) []*net.IPNet {
	if len(entries) == 0 {
		return nil
	}
	cidrs := make([]*net.IPNet, 0, len(entries))
	for _, entry := range entries {
		value := strings.TrimSpace(entry)
		if value == "" {
			continue
		}
		if !strings.Contains(value, "/") {
			ip := net.ParseIP(value)
			if ip == nil {
				continue
			}
			if ip.To4() != nil {
				value += "/32"
			} else {
				value += "/128"
			}
		}
		if _, network, err := net.ParseCIDR(value); err == nil {
			cidrs = append(cidrs, network)
		}
	}
	return cidrs
}

func isTrustedForwarder(ip net.IP, additional []*net.IPNet) bool {
	if ip == nil {
		return false
	}
	if ip.IsLoopback() || ip.IsPrivate() || ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() {
		return true
	}
	for _, cidr := range googleLoadBalancerCIDRs {
		if cidr.Contains(ip) {
			return true
		}
	}
	for _, cidr := range additional {
		if cidr.Contains(ip) {
			return true
		}
	}
	return false
}

func firstForwardedIP(header string) string {
	if header == "" {
		return ""
	}
	for _, part := range strings.Split(header, ",") {
		candidate := strings.TrimSpace(part)
		if candidate == "" {
			continue
		}
		ip := net.ParseIP(candidate)
		if ip != nil {
			return ip.String()
		}
	}
	return ""
}

func resolveClientIP(c *gin.Context, trusted []*net.IPNet) string {
	if c == nil || c.Request == nil {
		return ""
	}

	remoteIP := net.ParseIP(c.RemoteIP())
	if remoteIP == nil {
		return ""
	}

	if isTrustedForwarder(remoteIP, trusted) {
		if ip := firstForwardedIP(c.Request.Header.Get("X-Forwarded-For")); ip != "" {
			return ip
		}
		if ip := strings.TrimSpace(c.Request.Header.Get("X-Real-IP")); ip != "" {
			if parsed := net.ParseIP(ip); parsed != nil {
				return parsed.String()
			}
		}
	}

	return remoteIP.String()
}

func newTrustedProxyList(additional []string) []*net.IPNet {
	list := make([]*net.IPNet, len(googleLoadBalancerCIDRs))
	copy(list, googleLoadBalancerCIDRs)
	if extra := parseCIDRList(additional); len(extra) > 0 {
		list = append(list, extra...)
	}
	return list
}
