package api

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/grafana/grafana/pkg/api/datasource"
	"github.com/grafana/grafana/pkg/api/pluginproxy"
	"github.com/grafana/grafana/pkg/infra/metrics"
	"github.com/grafana/grafana/pkg/models"
	"github.com/grafana/grafana/pkg/plugins"
)

var reValidPath = regexp.MustCompile("^/api/datasources/proxy/\\d+/api/v1/(query|query_range)")
var reValidQuery = regexp.MustCompile("^([a-zA-Z0-9_]+|{__name__=~'\\w+\\.\\*'})")

// ProxyDataSourceRequest proxies datasource requests
func (hs *HTTPServer) ProxyDataSourceRequest(c *models.ReqContext) {
	c.TimeRequest(metrics.MDataSourceProxyReqTimer)

	dsID := c.ParamsInt64(":id")
	ds, err := hs.DatasourceCache.GetDatasource(dsID, c.SignedInUser, c.SkipCache)
	if err != nil {
		if errors.Is(err, models.ErrDataSourceAccessDenied) {
			c.JsonApiErr(403, "Access denied to datasource", err)
			return
		}
		c.JsonApiErr(500, "Unable to load datasource meta data", err)
		return
	}

	// find plugin
	plugin, ok := plugins.DataSources[ds.Type]
	if !ok {
		c.JsonApiErr(500, "Unable to find datasource plugin", err)
		return
	}

	if !hs.isValidProxyCall(c, err) {
		return
	}

	// macaron does not include trailing slashes when resolving a wildcard path
	proxyPath := ensureProxyPathTrailingSlash(c.Req.URL.Path, c.Params("*"))

	proxy, err := pluginproxy.NewDataSourceProxy(ds, plugin, c, proxyPath, hs.Cfg)
	if err != nil {
		if errors.Is(err, datasource.URLValidationError{}) {
			c.JsonApiErr(400, fmt.Sprintf("Invalid data source URL: %q", ds.Url), err)
		} else {
			c.JsonApiErr(500, "Failed creating data source proxy", err)
		}
		return
	}
	proxy.HandleRequest()
}

func (hs *HTTPServer) isValidProxyCall(c *models.ReqContext, err error) bool {
	if !reValidPath.MatchString(c.Req.URL.Path) {
		c.JsonApiErr(404, "Oooops", err)
		return false
	}
	query := c.Req.URL.Query().Get("query")
	if !reValidQuery.MatchString(query) {
		c.JsonApiErr(404, "Oooops", err)
		return false
	}
	return true
}

// ensureProxyPathTrailingSlash Check for a trailing slash in original path and makes
// sure that a trailing slash is added to proxy path, if not already exists.
func ensureProxyPathTrailingSlash(originalPath, proxyPath string) string {
	if len(proxyPath) > 1 {
		if originalPath[len(originalPath)-1] == '/' && proxyPath[len(proxyPath)-1] != '/' {
			return proxyPath + "/"
		}
	}

	return proxyPath
}
