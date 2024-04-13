package tests

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/Alladan04/avito_test/internal/models"
	bannerDelivery "github.com/Alladan04/avito_test/internal/pkg/banner/delivery/http"
	bannerRepo "github.com/Alladan04/avito_test/internal/pkg/banner/repo"
	bannerUsecase "github.com/Alladan04/avito_test/internal/pkg/banner/usecase"
	"github.com/Alladan04/avito_test/internal/pkg/middleware"
	"github.com/Alladan04/avito_test/internal/pkg/utils"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/suite"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println(err)
	}
}

type APITestSuite struct {
	suite.Suite

	db      *pgxpool.Pool
	conn    *pgx.Conn
	redisdb *redis.Client
	handler *bannerDelivery.BannerHandler
	uc      *bannerUsecase.BannerUsecase
	repo    *bannerRepo.BannerRepo
	//cache   *bannerRepo.CacheRepo
}

func TestAPISuite(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	suite.Run(t, new(APITestSuite))
}

func (s *APITestSuite) SetupSuite() {
	db, err := pgxpool.Connect(context.Background(), os.Getenv("TEST_DB"))
	if err != nil {
		s.FailNow("Failed to connect to pgxpool", err)
	} else {
		s.db = db
	}
	conn, err := pgx.Connect(context.Background(), os.Getenv("TEST_DB"))
	if err != nil {
		s.FailNow("Failed to connect to pgx", err)

	} else {
		s.conn = conn
	}
	redisOpts, err := redis.ParseURL(os.Getenv("TEST_REDIS"))
	if err != nil {
		s.FailNow("Failed to connect to redis", err)
	} else {
		redisDB := redis.NewClient(redisOpts)
		s.redisdb = redisDB
	}

	s.initDeps()

	if err := s.populateDB(); err != nil {
		s.FailNow("Failed to populate DB", err)
	}
}

func (s *APITestSuite) TearDownSuite() {
	s.db.Close() //nolint:errcheck

}

func (s *APITestSuite) initDeps() {
	// Init domain deps
	repo := bannerRepo.NewBannerRepo(s.db, *s.conn)
	cacherepo := bannerRepo.NewCacheRepo(*s.redisdb)
	uc := bannerUsecase.NewBannerUsecase(repo, cacherepo)
	h := bannerDelivery.NewBannerHandler(uc)
	s.repo = repo
	s.uc = uc
	s.handler = h

}

func (s *APITestSuite) TestUserBanner() {
	router := mux.NewRouter().PathPrefix("/api").Subrouter()
	router.Handle("/user_banner", middleware.JwtMiddleware(http.HandlerFunc(s.handler.GetOne))).Methods(http.MethodGet, http.MethodOptions)

	r := s.Require()

	jwt, _ := utils.GenToken(models.User{
		Id:         1,
		Username:   "testuser",
		Password:   "1234",
		CreateTime: time.Now().UTC(),
		IsAdmin:    false,
	}, time.Hour*24)
	req, _ := http.NewRequestWithContext(context.Background(), "GET", "/api/user_banner?feature_id=1&tag_id=1", nil)
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("Token", "Bearer "+jwt)

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	r.Equal(http.StatusOK, resp.Code)

}

func (s *APITestSuite) populateDB() error {
	const (
		insertBanner  = "INSERT INTO banner (title, banner_data, feature_id, url, create_time, update_time, is_active) VALUES ('some title', 'some data', 1, 'some url', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, true);"
		insertFeature = "INSERT INTO feature (id) VALUES (DEFAULT);"
		insertTag     = "INSERT INTO tag(id) VALUES (DEFAULT);"
		insertBT      = "INSERT INTO banner_tag (banner_id, tag_id,feature_id) VALUES (1,1,1);"
	)
	_, err := s.db.Exec(context.Background(), insertFeature)
	if err != nil {
		return err
	}
	_, err = s.db.Exec(context.Background(), insertTag)
	if err != nil {
		return err
	}
	_, err = s.db.Exec(context.Background(), insertBanner)
	if err != nil {
		return err
	}
	_, err = s.db.Exec(context.Background(), insertBT)
	if err != nil {
		return err
	}

	return nil
}
