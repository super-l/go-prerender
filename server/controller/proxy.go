package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-prerender/internal"
	"go-prerender/internal/config"
	"go-prerender/internal/storage"
	"go-prerender/utils"
	"strings"
)

func Proxy(c *gin.Context) {
	// 获取当前URI
	uri := c.Request.URL.RequestURI()
	uri = strings.ReplaceAll(uri, "/proxy", "")
	url := fmt.Sprintf("%s/%s", config.GetConfig().System.BaseUrl, uri)

	domain, err := utils.UriLib.GetDomain(config.GetConfig().System.BaseUrl)
	if err != nil {
		internal.SLogger.GetStdoutLogger().Error(err)
		return
	}

	// 根据uri,获取MD5
	nameMd5 := utils.Md5Lib.MD5(uri)

	// 尝试读取缓存，如果缓存并不存在，则实时抓取
	content, err := storage.GetData(domain, nameMd5)
	if err != nil {
		content, err = internal.GetSource(url)
		if err != nil {
			internal.SLogger.GetStdoutLogger().Error(err)
			return
		}

		// 如果开启h5,则添加<!doctype html>
		if config.GetConfig().System.H5 {
			content = "<!doctype html>" + content
		}

		// 设置本地缓存
		_ = storage.SaveData(domain, nameMd5, content)
	}
	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(200, content)
}

func GetQueryParams(c *gin.Context) map[string]any {
	query := c.Request.URL.Query()
	var queryMap = make(map[string]any, len(query))
	for k := range query {
		queryMap[k] = c.Query(k)
	}
	return queryMap
}
