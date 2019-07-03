package main

import (
	"context"
	"fmt"
	"github.com/99designs/gqlgen/handler"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
)


const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	// Echo instance
	e := echo.New()
	graphqlHandler := handler.GraphQL(NewExecutableSchema(Config{Resolvers: &Resolver{}}))
	playgroundHandler := handler.Playground("GraphQL", "/query")

	// Middleware
	e.Use(EchoContextToContextMiddleware)
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Route => handler
	e.GET("/", func(c echo.Context) error {
		cc := c.(*CustomContext)
		return cc.String(http.StatusOK, "Hello, World!\n")
	})
	e.POST("/query", func(c echo.Context) error {
		cc := c.(*CustomContext)
		req := cc.Request()
		res := cc.Response()
		graphqlHandler.ServeHTTP(res, req)
		return nil
	})

	e.GET("/playground", func(c echo.Context) error {
		cc := c.(*CustomContext)
		req := cc.Request()
		res := cc.Response()
		playgroundHandler.ServeHTTP(res, req)
		return nil
	})

	// Start server
	e.Logger.Fatal(e.Start(":"+port))
}

type CustomContext struct {
	echo.Context
	ctx context.Context
}

func (c *CustomContext) Foo() {
	println("foo")
}

func (c *CustomContext) Bar() {
	println("bar")
}

func EchoContextFromContext(ctx context.Context) (*echo.Context, error) {
	echoContext := ctx.Value("EchoContextKey")
	if echoContext == nil {
		err := fmt.Errorf("could not retrieve echo.Context")
		return nil, err
	}

	ec, ok := echoContext.(*echo.Context)
	if !ok {
		err := fmt.Errorf("echo.Context has wrong type")
		return nil, err
	}
	return ec, nil
}

func EchoContextToContextMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := context.WithValue(c.Request().Context(), "EchoContextKey", c)
		c.SetRequest(c.Request().WithContext(ctx))

		cc := &CustomContext{c, ctx}

		return next(cc)
	}
}
