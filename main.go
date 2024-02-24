package main

import (
	"crypto/tls"
	"fmt"

	"github.com/cbstorm/forward_requests/configs"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"github.com/gofiber/fiber/v2/middleware/proxy"
	"github.com/valyala/fasthttp"
)

func main() {
	app_config := configs.GetConfig()
	app_config.Load()
	app := fiber.New()
	proxy.WithTlsConfig(&tls.Config{
		InsecureSkipVerify: true,
	})
	proxy.WithClient(&fasthttp.Client{
		NoDefaultUserAgentHeader: true,
		DisablePathNormalizing:   true,
	})
	addr := fmt.Sprintf(":%s", app_config.APP_PORT)
	if app_config.ENV == "development" {
		addr = fmt.Sprintf("%s%s", "localhost", addr)
	}
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, X-Token, X-Refresh",
	}))
	app.Use(healthcheck.New())
	app.Use(proxy.Balancer(proxy.Config{
		Servers: app_config.SERVERS,
		ModifyRequest: func(c *fiber.Ctx) error {
			c.Request().Header.Add("X-Forwarded-For", c.IP())
			return nil
		},
		ModifyResponse: func(c *fiber.Ctx) error {
			c.Response().Header.Del(fiber.HeaderServer)
			return nil
		},
	}))
	app.Listen(addr)
}
