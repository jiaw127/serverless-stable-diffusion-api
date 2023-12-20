package main

import (
	"flag"
	"github.com/devsapp/serverless-stable-diffusion-api/pkg/config"
	"github.com/devsapp/serverless-stable-diffusion-api/pkg/datastore"
	"github.com/devsapp/serverless-stable-diffusion-api/pkg/server"
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	defaultPort       = "8010"
	defaultDBType     = datastore.SQLite
	shutdownTimeout   = 5 * time.Second // 5s
	defaultConfigPath = "config.json"
)

func handleSignal() {
	// Wait for interrupt signal to gracefully shutdown the server with
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logrus.Info("Shutting down server...")
}

func logInit(logLevel, logFile string) {
	switch logLevel {
	case "debug":
		logrus.SetLevel(logrus.DebugLevel)
		// include function and file
		logrus.SetReportCaller(true)
	case "dev":
		logrus.SetLevel(logrus.InfoLevel)
	default:
		logrus.SetLevel(logrus.WarnLevel)
	}
	logrus.SetOutput(&lumberjack.Logger{
		Filename:   logFile, // log output dir
		MaxSize:    50,      // MB, max logfile size
		MaxBackups: 3,       // max num logfile
		MaxAge:     3,       // max day logfile
		Compress:   false,   // is compress or not
	})
}

func main() {
	port := flag.String("port", defaultPort, "server listen port, default 8010")
	dbType := flag.String("dbType", string(defaultDBType), "db type default sqlite")
	configFile := flag.String("config", defaultConfigPath, "default config path")
	mode := flag.String("mode", "dev", "service work mode debug|dev|product")
	sdShell := flag.String("sd", "", "sd start shell")
	logFile := flag.String("log", "", "log output dir")
	flag.Parse()
	// init log
	logInit(*mode, *logFile)
	logrus.Info("agent start")

	// init config
	if err := config.InitConfig(*configFile); err != nil {
		logrus.Fatal(err.Error())
	}
	logrus.Info(*sdShell)
	config.ConfigGlobal.SdShell = *sdShell

	// init server and start
	agent, err := server.NewAgentServer(*port, datastore.DatastoreType(*dbType), *mode)
	if err != nil {
		logrus.Fatal("agent server init fail")
	}
	go agent.Start()

	// wait shutdown signal
	handleSignal()

	if err := agent.Close(shutdownTimeout); err != nil {
		logrus.Fatal("Shutdown server fail")
	}

	logrus.Info("Server exited")
}
