package main

import (
	"flag"
	"fmt"
	"github.com/devsapp/serverless-stable-diffusion-api/pkg/config"
	"github.com/devsapp/serverless-stable-diffusion-api/pkg/datastore"
	"github.com/devsapp/serverless-stable-diffusion-api/pkg/server"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	defaultPort       = "9000"
	defaultDBType     = datastore.TableStore
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
	//logrus.SetOutput(&lumberjack.Logger{
	//	Filename:   logFile, // log output dir
	//	MaxSize:    50,      // MB, max logfile size
	//	MaxBackups: 3,       // max num logfile
	//	MaxAge:     3,       // max day logfile
	//	Compress:   false,   // is compress or not
	//})
}

func main() {
	port := flag.String("port", defaultPort, "server listen port, default 8080")
	dbType := flag.String("dbType", string(defaultDBType), "db type default sqlite")
	configFile := flag.String("config", defaultConfigPath, "default config path")
	mode := flag.String("mode", "dev", "service work mode debug|dev|product")
	logFile := flag.String("log", "", "log output dir")
	flag.Parse()

	logInit(*mode, *logFile)
	logrus.Info(fmt.Sprintf("%s start", os.Getenv(config.SERVER_NAME)))

	// init config
	if err := config.InitConfig(*configFile); err != nil {
		logrus.Fatal(err.Error())
	}

	// init server and start
	proxy, err := server.NewProxyServer(*port, datastore.DatastoreType(*dbType), *mode)
	if err != nil {
		logrus.Fatal("proxy server init fail")
	}
	go proxy.Start()

	// wait shutdown signal
	handleSignal()

	if err := proxy.Close(shutdownTimeout); err != nil {
		logrus.Fatal("Shutdown server fail")
	}

	logrus.Info("Server exited")
}
