package app

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
	"github.com/nedokyrill/posts-service/graphql"
	"github.com/nedokyrill/posts-service/internal/resolvers"
	"github.com/nedokyrill/posts-service/internal/service"
	"github.com/nedokyrill/posts-service/internal/storage"
	"github.com/nedokyrill/posts-service/internal/storage/mem"
	"github.com/nedokyrill/posts-service/internal/storage/postgres"
	"github.com/nedokyrill/posts-service/pkg/consts"
	"github.com/nedokyrill/posts-service/pkg/db"
	"github.com/nedokyrill/posts-service/pkg/logger"
	"github.com/nedokyrill/posts-service/pkg/server"
	"github.com/nedokyrill/posts-service/pkg/utils"
)

func Run() {
	// Init LOGGER
	logger.InitLogger()

	// Load ENVIRONMENT VARIABLES
	err := godotenv.Load()
	if err != nil {
		logger.Logger.Fatal("error loading .env file, exiting...")
	}

	// Init REPO layer
	var postStore storage.PostStorage
	var commStore storage.CommentStorage

	if os.Getenv("IN_MEM_STORAGE") == "true" {
		logger.Logger.Info("using memory storage")
		postStore = mem.NewPostStorageMem()
		commStore = mem.NewCommentsStorageMem()
	} else {
		logger.Logger.Info("using postgres storage")
		ctx, cancel := context.WithTimeout(context.Background(), consts.PgxTimeout)
		defer cancel()

		conn, err := db.Connect(ctx)
		if err != nil {
			logger.Logger.Fatal("error connecting to database, exiting...")
		}
		defer conn.Close()

		postStore = postgres.NewPostStorePgx(conn)
		commStore = postgres.NewCommentsStorePgx(conn)
	}

	// Init SERVICE layer
	postServ := service.NewPostService(postStore)
	commServ := service.NewCommentService(commStore, postStore)
	viewerServ := service.NewViewerService()

	// Init ROUTER n start SERVER
	hand := handler.New(graphql.NewExecutableSchema(graphql.Config{Resolvers: &resolvers.Resolver{
		PostService:    postServ,
		CommentService: commServ,
		ViewerService:  viewerServ,
	}}))
	hand.AddTransport(transport.POST{})    // поддержка post
	hand.AddTransport(transport.GET{})     // поддержка get
	hand.AddTransport(transport.Websocket{ // поддержка вебсокетов
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	})

	router := utils.NewGinRouter()

	// Init ENDPOINTS
	router.POST("/query", gin.WrapH(hand))
	router.GET("/query", gin.WrapH(hand))
	router.GET("/", gin.WrapH(playground.Handler("graphQL playground", "/query")))

	srv := server.NewAPIServer(router)

	// START
	go srv.Start()

	// GRACEFUL SHUTDOWN
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Logger.Info("shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err = srv.Shutdown(ctx); err != nil {
		logger.Logger.Fatalw("shutdown error",
			"error", err)
	}
}
