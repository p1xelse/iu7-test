package time_tracker

import (
	"fmt"
	"timetracker/cmd/time_tracker/flags"
	_authDelivery "timetracker/internal/Auth/delivery"
	authRepPostgres "timetracker/internal/Auth/repository/postgres"
	authRep "timetracker/internal/Auth/repository/redis"
	authUsecase "timetracker/internal/Auth/usecase"
	_entryDelivery "timetracker/internal/Entry/delivery"
	entryRep "timetracker/internal/Entry/repository/postgres"
	entryUsecase "timetracker/internal/Entry/usecase"
	_friendDelivery "timetracker/internal/Friends/delivery"
	friendRep "timetracker/internal/Friends/repository/postgres"
	friendUsecase "timetracker/internal/Friends/usecase"
	_goalDelivery "timetracker/internal/Goal/delivery"
	goalRep "timetracker/internal/Goal/repository/postgres"
	goalUsecase "timetracker/internal/Goal/usecase"
	_projectDelivery "timetracker/internal/Project/delivery"
	projectRep "timetracker/internal/Project/repository/postgres"
	projectUsecase "timetracker/internal/Project/usecase"
	_tagDelivery "timetracker/internal/Tag/delivery"
	tagRep "timetracker/internal/Tag/repository/postgres"
	tagUsecase "timetracker/internal/Tag/usecase"
	_userDelivery "timetracker/internal/User/delivery"
	userRep "timetracker/internal/User/repository/postgres"
	userUsecase "timetracker/internal/User/usecase"
	"timetracker/internal/cache"
	"timetracker/internal/middleware"

	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
)

type TimeTracker struct {
	base
	PostgresClient            flags.PostgresFlags `toml:"postgres-client"`
	RedisSessionClient        flags.RedisFlags    `toml:"redis-client"`
	RedisProjectStorageClient flags.RedisFlags    `toml:"redis-project-storage-client"`
	Server                    flags.ServerFlags   `toml:"server"`
}

func (tt TimeTracker) Run(sessionDB string) error {
	e := echo.New()
	services, err := tt.Init(e)

	logger := services.Logger

	if err != nil {
		return fmt.Errorf("can not init services: %w", err)
	}

	postgresClient, err := tt.PostgresClient.Init()

	if err != nil {
		logger.Error("can not connect to Postgres client: %w", err)
		return err
	} else {
		logger.Info("Success conect to postgres")
	}

	redisSessionClient, err := tt.RedisSessionClient.Init()

	if err != nil {
		logger.Error("can not connect to Redis session client: %w", err)
		return err
	} else {
		logger.Info("Success conect to redis")
	}

	redisCacheClient, err := tt.RedisProjectStorageClient.Init()

	if err != nil {
		logger.Error("can not connect to Redis cache client: %w", err)
		return err
	} else {
		logger.Info("Success conect to redis")
	}

	entryRepo := entryRep.NewEntryRepository(postgresClient)
	userRepo := userRep.NewUserRepository(postgresClient)
	tagRepo := tagRep.NewTagRepository(postgresClient)
	goalRepo := goalRep.NewGoalRepository(postgresClient)
	projectRepo := projectRep.NewProjectRepository(postgresClient)
	authRepo := authRep.NewAuthRepository(redisSessionClient)
	authPostgresRepo := authRepPostgres.NewAuthRepositoryPostgres(postgresClient)
	friendRepo := friendRep.NewFriendRepository(postgresClient)
	cacheStorage := cache.NewStorageRedis(redisCacheClient)

	entryUC := entryUsecase.New(entryRepo, tagRepo, userRepo)
	goalUC := goalUsecase.New(goalRepo)
	projectUC := projectUsecase.New(projectRepo, cacheStorage)
	tagUC := tagUsecase.New(tagRepo)

	authUC := authUsecase.New(userRepo, authRepo)
	if sessionDB == "postgres" {
		authUC = authUsecase.New(userRepo, authPostgresRepo)
	}

	userUC := userUsecase.New(userRepo)
	friendUC := friendUsecase.New(friendRepo, userRepo)

	aclMiddleware := middleware.NewAclMiddleware(friendUC)

	_entryDelivery.NewDelivery(e, entryUC, aclMiddleware)
	_goalDelivery.NewDelivery(e, goalUC, aclMiddleware)
	_projectDelivery.NewDelivery(e, projectUC, aclMiddleware)
	_tagDelivery.NewDelivery(e, tagUC, aclMiddleware)
	_authDelivery.NewDelivery(e, authUC)
	_userDelivery.NewDelivery(e, userUC, aclMiddleware)
	_friendDelivery.NewDelivery(e, friendUC, aclMiddleware)

	e.Use(echoMiddleware.LoggerWithConfig(echoMiddleware.LoggerConfig{
		Format: tt.Logger.LogHttpFormat,
		Output: logger.Output(),
	}))

	e.Use(echoMiddleware.Recover())
	authMiddleware := middleware.NewMiddleware(authUC)
	e.Use(authMiddleware.Auth)

	httpServer := tt.Server.Init(e)
	server := Server{*httpServer}
	if err := server.Start(); err != nil {
		logger.Fatal(err)
	}
	return nil
}
