package config

import (
	"github.com/mr-time2028/WebChat/internal/database"
	"github.com/mr-time2028/WebChat/internal/models"
	"log"
	"os"
	"os/signal"
	"syscall"
)

type App struct {
	Domain   string
	Debug    bool
	HTTPPort string
	DB       *database.DB
	Hub      *models.Hub
	Models   *models.ModelManager
	Auth     *models.Auth
}

func (a *App) ListenForShutdown() {
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	close(quit)
	a.Shutdown()
	os.Exit(0)
}

func (a *App) Shutdown() {
	log.Println("cleanup tasks...")
	close(a.Hub.RequestChan)
	close(a.Hub.ResponseChan)
	log.Println("server shutdown gracefully")
}
