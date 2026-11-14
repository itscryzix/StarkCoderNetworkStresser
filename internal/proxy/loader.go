package proxy

import (
	"bufio"
	"net"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/sammwyy/mikumikubeam/internal/engine"
)

var reAuth = regexp.MustCompile(`^([^:]+):([^@]+)@(.+)$`)

// splitHostPort parses "host", "host:port", or "[ipv6]:port" forms.
// It never fails: on parse error it returns the input as host with port 0.
func splitHostPort(rest string) (string, int, error) {
	if h, p, err := net.SplitHostPort(rest); err == nil {
		pi, _ := strconv.Atoi(p)
		return h, pi, nil
	}
	// Bare [ipv6] without port.
	if strings.HasPrefix(rest, "[") && strings.HasSuffix(rest, "]") {
		return strings.Trim(rest, "[]"), 0, nil
	}
	// host:port with a single colon (IPv4/hostname). Multiple colons
	// without brackets is a bare IPv6 address, not host:port.
	if strings.Count(rest, ":") == 1 {
		hp := strings.SplitN(rest, ":", 2)
		if p, err := strconv.Atoi(hp[1]); err == nil {
			return hp[0], p, nil
		}
		return hp[0], 0, nil
	}
	return rest, 0, nil
}

func LoadProxies(path string) ([]engine.Proxy, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var out []engine.Proxy
	s := bufio.NewScanner(f)
	for s.Scan() {
		line := strings.TrimSpace(s.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		proto := "http"
		rest := line
		if strings.Contains(line, "://") {
			parts := strings.SplitN(line, "://", 2)
			proto = parts[0]
			rest = parts[1]
		}
		user := ""
		pass := ""
		if m := reAuth.FindStringSubmatch(rest); len(m) == 4 {
			user, pass, rest = m[1], m[2], m[3]
		}
		host := rest
		port := 8080
		if h, p, err := splitHostPort(rest); err == nil {
			host = h
			if p != 0 {
				port = p
			}
		}
		out = append(out, engine.Proxy{Username: user, Password: pass, Protocol: proto, Host: host, Port: port})
	}
	return out, s.Err()
}

func LoadUserAgents(path string) ([]string, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(b), "\n")
	uas := make([]string, 0, len(lines))
	for _, ln := range lines {
		ln = strings.TrimSpace(ln)
		if ln == "" || strings.HasPrefix(ln, "#") {
			continue
		}
		uas = append(uas, ln)
	}
	return uas, nil
}

// Filter proxies by method. For now, mimic original behavior.
func FilterByMethod(proxies []engine.Proxy, method engine.AttackKind) []engine.Proxy {
	allowed := map[engine.AttackKind]map[string]bool{
		engine.AttackHTTPFlood:     {"http": true, "https": true, "socks4": true, "socks5": true},
		engine.AttackHTTPBurst:     {"http": true, "https": true, "socks4": true, "socks5": true},
		engine.AttackHTTPBypass:    {"http": true, "https": true, "socks4": true, "socks5": true},
		engine.AttackHTTPSlowloris: {"http": true, "https": true, "socks4": true, "socks5": true},
		engine.AttackTCPFlood:      {"http": true, "https": true, "socks4": true, "socks5": true},
		engine.AttackTCPBurst:      {"http": true, "https": true, "socks4": true, "socks5": true},
		engine.AttackTCPSlowloris:  {"http": true, "https": true, "socks4": true, "socks5": true},
		engine.AttackMinecraftPing: {"http": true, "https": true, "socks4": true, "socks5": true},
	}
	allowedSet, ok := allowed[method]
	if !ok {
		// Unknown method: don't silently drop all proxies.
		out := make([]engine.Proxy, len(proxies))
		copy(out, proxies)
		return out
	}
	out := make([]engine.Proxy, 0, len(proxies))
	for _, p := range proxies {
		if allowedSet[p.Protocol] {
			out = append(out, p)
		}
	}
	return out
}
