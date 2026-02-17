package main

import (
	"context"
	"log"
	"miniflux-mcp/pkg/mcp"
	"net"
	"net/http"
	"os"
	"time"

	mcpSdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

func main() {
	minifluxUrl, exists := os.LookupEnv("MINIFLUX_URL")
	if !exists {
		log.Fatal("Require MINIFLUX_URL environment variable set")
	}

	host, exists := os.LookupEnv("HOST")
	if !exists {
		host = "127.0.0.1"
	}

	port, exists := os.LookupEnv("PORT")
	if !exists {
		port = "8080"
	}

	server := mcpSdk.NewServer(&mcpSdk.Implementation{Name: "miniflux"}, nil)

	server.AddReceivingMiddleware(createLoggingMiddleware())

	mcp.RegisterTools(server, minifluxUrl)

	handler := mcpSdk.NewStreamableHTTPHandler(func(r *http.Request) *mcpSdk.Server {
		return server
	}, nil)

	listenAddr := net.JoinHostPort(host, port)

	log.Printf("listening on %s", listenAddr)

	if err := http.ListenAndServe(listenAddr, handler); err != nil {
		log.Fatal(err)
	}
}

func createLoggingMiddleware() mcpSdk.Middleware {
	return func(next mcpSdk.MethodHandler) mcpSdk.MethodHandler {
		return func(
			ctx context.Context,
			method string,
			req mcpSdk.Request,
		) (mcpSdk.Result, error) {
			start := time.Now()
			sessionID := req.GetSession().ID()

			// Log request details.
			log.Printf("[REQUEST] Session: %s | Method: %s",
				sessionID,
				method)

			// Call the actual handler.
			result, err := next(ctx, method, req)

			// Log response details.
			duration := time.Since(start)

			if err != nil {
				log.Printf("[RESPONSE] Session: %s | Method: %s | Status: ERROR | Duration: %v | Error: %v",
					sessionID,
					method,
					duration,
					err)
			} else {
				log.Printf("[RESPONSE] Session: %s | Method: %s | Status: OK | Duration: %v",
					sessionID,
					method,
					duration)
			}

			return result, err
		}
	}
}
