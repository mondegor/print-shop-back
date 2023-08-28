package mrcontext

import (
    "context"
    "net"
    "net/http"
)

func UserIpFromRequest(r *http.Request) (net.IP, error) {
    ip, _, err := net.SplitHostPort(r.RemoteAddr)

    if err != nil {
        return net.IP{}, ErrHttpRequestUserIP.New(r.RemoteAddr)
    }

    parsedIp := net.ParseIP(ip)

    if parsedIp != nil {
        return parsedIp, nil
    }

    return net.IP{}, ErrHttpRequestParseUserIP.New(ip)
}

func UserIpNewContext(ctx context.Context, userIP net.IP) context.Context {
    return context.WithValue(ctx, ctxUserIPKey, userIP)
}

func GetUserIp(ctx context.Context) net.IP {
    value, ok := ctx.Value(ctxUserIPKey).(net.IP)

    if ok {
        return value
    }

    return net.IP{}
}
