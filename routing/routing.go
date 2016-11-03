package routing

import (
	"strings"
	"os/exec"
	"regexp"
	"net/url"
	"net/http"
	"fmt"
)

type Cache struct {
	registry map[string]*url.URL
}

func (r *Cache) Init() *Cache {
	r.registry = make(map[string]*url.URL)
	return r
}

func (r *Cache) hasCache(hostname string) bool {
	if _, ok := r.registry[hostname]; ok {
		return true
	}
	return false
}

func (r *Cache) setCache(hostname string, targetUrl *url.URL) {
	r.registry[hostname] = targetUrl
}

func (r *Cache) getCache(hostname string) *url.URL {
	return r.registry[hostname]
}

func ResolveTargetContainer(r *http.Request, cache *Cache, debug bool, domain string) *url.URL {

	// get hostname without port
	host := r.Host
	hostinfo := strings.Split(host, ":")
	hostname := hostinfo[0]

	// remove internal domain
	if strings.Contains(hostname, domain) {
		hostname = strings.Replace(hostname, domain, "", 1)
	}

	// define cache key
	cacheKey := hostname

	// try to load results from cache
	targetUrl := &url.URL{}
	if !cache.hasCache(cacheKey) {

		// fallback: exact match
		container := findContainerInProcesslist(hostname)
		if "" == container {
			// fallback 1: search for match with underscrores
			container = findContainerInProcesslist(strings.Replace(hostname, "-", "_", -1))

			if "" == container {
				// fallback 2: searhc for matches with dashes
				container = findContainerInProcesslist(strings.Replace(hostname, "_", "-", -1))
			}
		}

		// determine final URL
		if container != "" {

			// SSL request?
			schema := "http://"
			if r.TLS != nil {
				schema = "https://"
			}

			// get final target URL for this reuqest
			targetUrl, _ = url.Parse(schema + container)
		}

		if debug {
			fmt.Println("Resolving hostname", hostname, "setting cache:", targetUrl)
		}

		// write cache for next request resolving (will be stored until next service restart)
		cache.setCache(cacheKey, targetUrl)
	} else {
		targetUrl = cache.getCache(cacheKey)
		if debug {
			fmt.Println("--> Using cached route", targetUrl, "for host", hostname)
		}
	}

	return targetUrl
}

func findContainerInProcesslist(hostname string) string {
	output, err := exec.Command("docker", "ps", "--filter=\"name=" + hostname + "\"").Output()
	if err != nil {
		return ""
	}
	processlist := strings.Split(string(output), "\n")

	container := ""
	for _, row := range processlist {

		// default container resolving without xdebug preference
		container = getMatch(row, hostname)
		if "" != container {
			break
		}
	}

	return container
}

func getMatch(row, hostname string) string {
	match := ""
	re := regexp.MustCompile(" (?P<host>([0-9.]+)):(?P<port>[0-9]+)->([0-9]+)/tcp(.+)" + hostname)
	matches := re.FindAllStringSubmatch(row, -1)
	if matches != nil {
		match = matches[0][1] + ":" + matches[0][3]
	}

	return match
}

