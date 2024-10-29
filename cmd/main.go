package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"simplewallet/config"
	"simplewallet/router"
	"simplewallet/util/db"
	"syscall"
	"time"

	_ "github.com/lib/pq"
)

func Init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	err := db.InitDb(&config.Config.Db)
	if err != nil {
		panic(err)
	}
	err = db.InitRedis(&config.Config.Redis)
	if err != nil {
		panic(err)
	}
}
func main() {

	Init()

	engine := router.InitRouter()
	addr := config.Config.GinHost
	server := &http.Server{
		Addr:         addr,
		Handler:      engine,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	SignalHandler(server)

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
		panic(err)
	}
}

func SignalHandler(server *http.Server) {
	logID := ""
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGILL, syscall.SIGTRAP, syscall.SIGABRT, syscall.SIGBUS, syscall.SIGFPE, syscall.SIGSEGV, syscall.SIGPIPE, syscall.SIGALRM, syscall.SIGTERM)
	go func() {
		s := <-c
		log.Printf("%s|received signal:%s\n", logID, s.String())
		if s == syscall.SIGTERM || s == syscall.SIGINT {
			maxSecond := time.Duration(30)
			ctx, cancel := context.WithTimeout(context.Background(), maxSecond*time.Second)
			defer cancel()
			if err := server.Shutdown(ctx); err != nil {
				log.Printf("%s|%s timeout,force to shutdown\n", logID, maxSecond.String())
				os.Exit(0)
			}
			log.Printf("%s|service shutdown success\n", logID)
			os.Exit(0)
		}
	}()
}
