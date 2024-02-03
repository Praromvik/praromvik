/*
MIT License

# Copyright (c) 2024 Praromvik

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package cmd

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/praromvik/praromvik/pkg/client"
	"github.com/praromvik/praromvik/routers"

	"cloud.google.com/go/firestore"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
)

type Server struct {
	router   http.Handler
	fClient  *firestore.Client
	mgClient *mongo.Client
	rdClient *redis.Client
}

func New(ctx context.Context) (*Server, error) {
	fClient, err := client.ConnectToFireStore(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get firestore client: %v", err)
	}

	mgClient, err := client.ConnectToMongoDB(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get MongoDB client: %v", err)
	}

	rdClient, err := client.ConnectToRedis()
	if err != nil {
		return nil, fmt.Errorf("failed to get redis client: %v", err)
	}
	err = client.TestRedisConnection(ctx, rdClient)

	app := &Server{
		router:   routers.LoadRoutes(fClient),
		fClient:  fClient,
		mgClient: mgClient,
		rdClient: rdClient,
	}
	return app, nil
}

func (a *Server) Start(ctx context.Context) error {
	server := &http.Server{
		Addr:    fmt.Sprintf(":%v", port),
		Handler: a.router,
	}

	// Server run context
	serverCtx, serverStopCtx := context.WithCancel(ctx)

	// Listen for syscall signals for process to interrupt/quit
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sig

		// Shutdown signal with grace period of 30 seconds
		shutdownCtx, _ := context.WithTimeout(serverCtx, 10*time.Second)

		go func() {
			<-shutdownCtx.Done()
			if errors.Is(shutdownCtx.Err(), context.DeadlineExceeded) {
				log.Fatal("graceful shutdown timed out.. forcing exit.")
			}
		}()

		// Trigger graceful shutdown
		err := server.Shutdown(shutdownCtx)
		if err != nil {
			log.Fatal(err)
		}
		serverStopCtx()
	}()

	fmt.Println("Listening on port", port)
	err := server.ListenAndServe()
	if err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}
	<-serverCtx.Done()

	return nil
}
