package main

import (
    "log"
    "net/http"
    "net/http/httputil"
    "net/url"
    "os"

    "github.com/gin-gonic/gin"
)

type ServiceConfig struct {
    Name string `yaml:"name"`
    URL  string `yaml:"url"`
    Path string `yaml:"path"`
}

type Config struct {
    Services []ServiceConfig `yaml:"services"`
}

func main() {
    router := gin.Default()

    // Auth Service Proxy
    authURL, _ := url.Parse(os.Getenv("AUTH_SERVICE_URL"))
    authProxy := httputil.NewSingleHostReverseProxy(authURL)
    
    // Company Service Proxy
    companyURL, _ := url.Parse(os.Getenv("COMPANY_SERVICE_URL"))
    companyProxy := httputil.NewSingleHostReverseProxy(companyURL)

    // Routes
    router.Any("/api/v1/auth/*path", gin.WrapH(authProxy))
    router.Any("/api/v1/companies/*path", gin.WrapH(companyProxy))

    // Health check
    router.GET("/health", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{"status": "healthy"})
    })

    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    log.Printf("API Gateway starting on port %s", port)
    router.Run(":" + port)
}