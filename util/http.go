package util
import(
	"log"
	"net"
	"net/url"
	"strings"
)

func IsValidSocmedURL(u *url.URL) bool {
	if u.Host == "accounts.google.com" && u.Path == "/o/oauth2/auth" {
		return true
	}
	if u.Host == "open.login.yahooapis.com" && u.Path == "/openid/op/auth" {
		return true
	}

	if u.Host == "www.facebook.com" && u.Path == "/dialog/oauth" {
		return true
	}

	return false
}

func IsTokopediaURL(u string) bool {
	if u == "" {
		return false
	}

	uObject, err := url.Parse(u)
	if err != nil {
		return false
	}

	host, _, err := net.SplitHostPort(uObject.Host)
	if err != nil {
		host = uObject.Host
	}

	// Filter google, yahoo, facebook URL for oauth login.
	if IsValidSocmedURL(uObject) {
		return true
	}

	if !strings.HasSuffix(host, ".tokopedia.com") && !strings.HasSuffix(host, ".tokopedia.net") && !strings.HasSuffix(host, "devel-go.tkpd") && !strings.HasSuffix(host, ".ndvl") && !strings.HasSuffix(host, ".tokocash.com") && !strings.HasSuffix(host, ".tokopedia.id") {
		log.Println("invalid host url[", u, "]")
		return false
	}

	return true
}
