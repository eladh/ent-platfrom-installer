package server

import (
	"common/kubernetes"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())

	e.Use(middleware.Recover())

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:8080"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	// Route => handle all rest endpoints
	e.GET("/pod", func(c echo.Context) error {
		pods, _ := kubernetes.GetPods()
		return c.JSON( 200 ,pods)
	})

	e.GET("/service", func(c echo.Context) error {
		services, _ := kubernetes.GetServices()
		return c.JSON( 200 ,services)
	})

	e.GET("/deployment", func(c echo.Context) error {
		deployments, _ := kubernetes.GetDeployments()
		return c.JSON( 200 ,deployments)
	})

	e.GET("/ingress", func(c echo.Context) error {
		ingresses, _ := kubernetes.GetIngress()
		return c.JSON( 200 ,ingresses)
	})

	e.GET("/daemonset", func(c echo.Context) error {
		daemonsets, _ := kubernetes.GetDaemonset()
		return c.JSON( 200 ,daemonsets)
	})

	// Start server
	e.Logger.Fatal(e.Start(":1323"))

	// Start monitor K8s cluster
	go kubernetes.Watch()

}