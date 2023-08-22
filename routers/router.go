package routers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/nikitamirzani323/wl_agen_backend_api/controllers"
	"github.com/nikitamirzani323/wl_agen_backend_api/middleware"
)

func Init() *fiber.App {
	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error {
		// Set some security headers:
		// c.Set("Content-Security-Policy", "frame-ancestors 'none'")
		// c.Set("X-XSS-Protection", "1; mode=block")
		// c.Set("X-Content-Type-Options", "nosniff")
		// c.Set("X-Download-Options", "noopen")
		// c.Set("Strict-Transport-Security", "max-age=5184000")
		// c.Set("X-Frame-Options", "SAMEORIGIN")
		// c.Set("X-DNS-Prefetch-Control", "off")

		// Go to next middleware:
		return c.Next()
	})
	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(compress.New())
	app.Get("/ipaddress", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":      fiber.StatusOK,
			"message":     "Success",
			"record":      "data",
			"BASEURL":     c.BaseURL(),
			"HOSTNAME":    c.Hostname(),
			"IP":          c.IP(),
			"IPS":         c.IPs(),
			"OriginalURL": c.OriginalURL(),
			"Path":        c.Path(),
			"Protocol":    c.Protocol(),
			"Subdomain":   c.Subdomains(),
		})
	})
	app.Get("/dashboard", monitor.New())

	app.Post("/api/login", controllers.CheckLogin)
	app.Post("/api/valid", middleware.JWTProtected(), controllers.Home)
	app.Post("/api/alladmin", middleware.JWTProtected(), controllers.Adminhome)
	app.Post("/api/saveadmin", middleware.JWTProtected(), controllers.AdminSave)
	app.Post("/api/alladminrule", middleware.JWTProtected(), controllers.Adminrulehome)
	app.Post("/api/saveadminrule", middleware.JWTProtected(), controllers.AdminruleSave)

	app.Post("/api/provider", middleware.JWTProtected(), controllers.Providerhome)
	app.Post("/api/catebank", middleware.JWTProtected(), controllers.CateBankhome)

	app.Post("/api/curr", middleware.JWTProtected(), controllers.Currhome)
	app.Post("/api/bank", middleware.JWTProtected(), controllers.Agenbankhome)
	app.Post("/api/banklist", middleware.JWTProtected(), controllers.Agenbanklist)
	app.Post("/api/banksave", middleware.JWTProtected(), controllers.AgenbankSave)
	app.Post("/api/member", middleware.JWTProtected(), controllers.Memberhome)
	app.Post("/api/membersearch", middleware.JWTProtected(), controllers.Membersearch)
	app.Post("/api/membersave", middleware.JWTProtected(), controllers.MemberSave)
	app.Post("/api/memberbanksave", middleware.JWTProtected(), controllers.MemberBankSave)
	app.Post("/api/memberbankdelete", middleware.JWTProtected(), controllers.MemberBankDelete)

	app.Post("/api/transaksidepowd", middleware.JWTProtected(), controllers.Transdpwdhome)
	app.Post("/api/transaksidepowdsave", middleware.JWTProtected(), controllers.TransdpwdSave)

	return app
}
