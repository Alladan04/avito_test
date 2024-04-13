package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	authDelivery "github.com/Alladan04/avito_test/internal/pkg/auth/delivery/http"
	authRepo "github.com/Alladan04/avito_test/internal/pkg/auth/repo"
	authUsecase "github.com/Alladan04/avito_test/internal/pkg/auth/usecase"
	bannerDelivery "github.com/Alladan04/avito_test/internal/pkg/banner/delivery/http"
	bannerRepo "github.com/Alladan04/avito_test/internal/pkg/banner/repo"
	bannerUsecase "github.com/Alladan04/avito_test/internal/pkg/banner/usecase"
	"github.com/Alladan04/avito_test/internal/pkg/middleware"
	"github.com/redis/go-redis/v9"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println(err)
	}
}
func main() {
	db, err := pgxpool.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Println(err)
		return
	}
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()

	redisOpts, err := redis.ParseURL(os.Getenv("REDIS_URL"))
	if err != nil {
		fmt.Println("redis not connected:", err)
		return
	}
	redisDB := redis.NewClient(redisOpts)
	AuthRepo := authRepo.NewAuthRepo(db)
	AuthUsecase := authUsecase.NewAuthUsecase(AuthRepo)
	AuthDelivery := authDelivery.NewAuthHandler(AuthUsecase)

	CacheRepo := bannerRepo.NewCacheRepo(*redisDB)
	BannerRepo := bannerRepo.NewBannerRepo(db, *conn)
	BannerUsecase := bannerUsecase.NewBannerUsecase(BannerRepo, CacheRepo)
	BannerDelivery := bannerDelivery.NewBannerHandler(BannerUsecase)

	r := mux.NewRouter().PathPrefix("/api").Subrouter()

	r.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		fmt.Println("not found")
	})
	http.Handle("/", r)
	auth := r.PathPrefix("/auth").Subrouter()
	{
		auth.Handle("/signup", http.HandlerFunc(AuthDelivery.SignUp)).Methods(http.MethodPost, http.MethodOptions)
		auth.Handle("/login", http.HandlerFunc(AuthDelivery.SignIn)).Methods(http.MethodPost, http.MethodOptions)
		auth.Handle("/logout", http.HandlerFunc(AuthDelivery.LogOut)).Methods(http.MethodDelete, http.MethodOptions)

	}
	banner := r
	{
		banner.Handle("/banner", middleware.JwtMiddleware(middleware.CheckAdminPermissionMiddleware(http.HandlerFunc(BannerDelivery.AddItem)))).Methods(http.MethodPost, http.MethodOptions)
		banner.Handle("/banner", middleware.JwtMiddleware(middleware.CheckAdminPermissionMiddleware(http.HandlerFunc(BannerDelivery.GetAll)))).Methods(http.MethodGet, http.MethodOptions)
		banner.Handle("/user_banner", middleware.JwtMiddleware(http.HandlerFunc(BannerDelivery.GetOne))).Methods(http.MethodGet, http.MethodOptions)
		banner.Handle("/banner/{id}", middleware.JwtMiddleware(middleware.CheckAdminPermissionMiddleware(http.HandlerFunc(BannerDelivery.UpdateBanner)))).Methods(http.MethodPatch, http.MethodOptions)
		banner.Handle("/banner/{id}", middleware.JwtMiddleware(middleware.CheckAdminPermissionMiddleware(http.HandlerFunc(BannerDelivery.DeleteBanner)))).Methods(http.MethodDelete, http.MethodOptions)

	}

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)

	server := http.Server{
		Handler:           r,
		Addr:              ":8080",
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			fmt.Println("Server Stopped")
		}
	}()
	fmt.Println("Server started ")

	sig := <-signalCh
	fmt.Println("recieved signal: ", sig)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		fmt.Println("Server shutdown failed: " + err.Error())
	}
}
