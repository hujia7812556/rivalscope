package main

import (
	"embed"
	"io/fs"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"rivalscope/internal/dto"
)

// 嵌入前端构建产物 web/dist。
// 说明:dist 内仅提交 .gitkeep 占位保证 go build 可通过;
// 实际部署请先执行 make web 生成真实产物后再构建二进制。
//
//go:embed all:web/dist
var webFS embed.FS

// registerWeb 注册前端静态资源服务与 SPA history 路由回退。
// 未匹配到 API 路由的请求:/api 前缀返回 JSON 404,其余路径尝试静态文件,
// 文件不存在时回退 index.html(交给前端路由处理)。
func registerWeb(r *gin.Engine) {
	dist, err := fs.Sub(webFS, "web/dist")
	if err != nil {
		return
	}
	fileServer := http.FileServer(http.FS(dist))

	r.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path
		if strings.HasPrefix(path, "/api") {
			c.JSON(http.StatusNotFound, dto.Response{
				Code: 404, Message: "接口不存在", Data: gin.H{},
			})
			return
		}
		// 命中真实静态文件(assets/js/css 等)则直接返回
		trimmed := strings.TrimPrefix(path, "/")
		if trimmed != "" {
			if f, err := dist.Open(trimmed); err == nil {
				f.Close()
				fileServer.ServeHTTP(c.Writer, c.Request)
				return
			}
		}
		// 其余路径统一回退到 SPA 入口
		index, err := fs.ReadFile(dist, "index.html")
		if err != nil {
			c.String(http.StatusNotFound, "前端资源未构建,请先执行 make web")
			return
		}
		c.Data(http.StatusOK, "text/html; charset=utf-8", index)
	})
}
