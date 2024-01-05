package main

import (
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"timetracker/benchmark"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

func prometheusHandler() gin.HandlerFunc {
	h := promhttp.Handler()

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func main() {
	gin.SetMode(gin.DebugMode)
	router := gin.Default()
	router.UseRawPath = true
	router.UnescapePathValues = false

	router.Use(gin.RecoveryWithWriter(os.Stdout))
	router.Use(gin.LoggerWithWriter(os.Stdout))

	metrics := benchmark.NewMetrics()

	router.GET("/metrics", prometheusHandler())
	encodeBenchmarker := benchmark.NewEncodeBenchmarker(metrics)
	decodeBenchmarker := benchmark.NewDecodeBenchmarker(metrics)
	router.POST("/bench/encode", func(ctx *gin.Context) {
		nRaw := ctx.Query("n")
		n, err := strconv.ParseInt(nRaw, 10, 64)
		if err != nil {
			_ = ctx.AbortWithError(http.StatusBadRequest, err)
			return
		}

		log.Println("Start benchmarking encode")
		res, err := encodeBenchmarker.EncodeBenchmark(int(n))
		if err != nil {
			_ = ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		ctx.JSON(http.StatusOK, res)
	})
	router.POST("/bench/decode", func(ctx *gin.Context) {
		nRaw := ctx.Query("n")
		n, err := strconv.ParseInt(nRaw, 10, 64)
		if err != nil {
			_ = ctx.AbortWithError(http.StatusBadRequest, err)
			return
		}

		log.Println("Start benchmarking decode")
		res, err := decodeBenchmarker.DecodeBenchmark(int(n))
		if err != nil {
			_ = ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		ctx.JSON(http.StatusOK, res)
	})

	s := http.Server{
		Addr:    ":8081",
		Handler: router,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		err := s.ListenAndServe()
		if err != nil && !errors.Is(http.ErrServerClosed, err) {
			panic(err)
		}
	}()

	<-quit
	println("Shutdown Server ...")

	err := s.Close()
	if err != nil {
		panic(err)
	}
	println("Server exited")
}
