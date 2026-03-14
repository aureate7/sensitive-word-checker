package main

import (
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type detectReq struct {
	Text       string         `json:"text"`
	Categories []string       `json:"categories"`
	Options    *DetectOptions `json:"options,omitempty"`
}

func main() {
	r := gin.Default()
	r.Use(cors.Default())

	// 设置模板路径
	r.LoadHTMLGlob(filepath.Join("templates", "*.html"))

	// 初始化敏感词检测器，加载 temp 目录下所有词库
	detector := NewDetector("temp")

	// 主页：渲染 index_new.html
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index_new.html", gin.H{
			"title": "敏感词检测系统 (Go + Gin)",
		})
	})

	// 检测接口
	r.POST("/api/detect", func(c *gin.Context) {
		var req detectReq
		if err := c.ShouldBindJSON(&req); err != nil || req.Text == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "text required"})
			return
		}
		res := detector.DetectWithOptions(req.Text, req.Categories, req.Options)
		c.JSON(http.StatusOK, res)
	})

	// 统计接口
	r.GET("/api/statistics", func(c *gin.Context) {
		c.JSON(http.StatusOK, detector.Statistics())
	})

	// 分类接口
	r.GET("/api/categories", func(c *gin.Context) {
		c.JSON(http.StatusOK, CategoryDisplay)
	})

	// 健康检查
	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	fmt.Println("===================================================")
	fmt.Println("🚀 敏感词检测系统启动成功！")
	fmt.Println("📍 访问地址: http://localhost:8008")
	fmt.Println("📂 词库目录: ./temp")
	fmt.Println("📄 模板目录: ./templates")
	fmt.Println("===================================================")

	_ = r.Run(":8008")
}
