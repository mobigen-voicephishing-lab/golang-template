package main

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	apphttp "github.com/mobigen/golang-web-template/internal/adapter/inbound/http"
	persistence "github.com/mobigen/golang-web-template/internal/adapter/outbound/persistence/gorm"
	"github.com/mobigen/golang-web-template/internal/bootstrap"
	"github.com/mobigen/golang-web-template/internal/infrastructure/config"
	"github.com/mobigen/golang-web-template/internal/infrastructure/db"
	"github.com/mobigen/golang-web-template/internal/infrastructure/logger"

	"github.com/sirupsen/logrus"
)

// Context context of main
type Context struct {
	Env       *config.Environment
	Conf      *config.Configuration
	Log       *logger.Logger
	CM        *config.ConfigManager
	Datastore *db.DataStore
	Router    *apphttp.Router
}

// InitLog Initialize logger
func (c *Context) InitLog() {
	log := logger.Logger{}.GetInstance()
	log.SetLogLevel(logrus.DebugLevel)
	c.Log = log
}

// ReadEnv Read value of the environment
func (c *Context) ReadEnv() error {
	c.Env = new(config.Environment)
	// Get Home
	homePath := os.Getenv("APP_HOME")
	if len(homePath) > 0 {
		c.Log.Errorf("APP_HOME : %s", homePath)
		c.Env.Home = homePath
	} else {
		dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			return err
		}
		c.Env.Home = dir
	}
	// Get Profile
	profile := os.Getenv("PROFILE")
	if len(profile) > 0 {
		c.Log.Errorf("PROFILE : %s", profile)
		c.Env.Profile = profile
	} else {
		c.Env.Profile = "prod"
	}
	// Get Log Level
	logLevel := os.Getenv("LOG_LEVEL")
	if len(logLevel) > 0 {
		c.Log.Errorf("LOG_LEVEL : %s", logLevel)
		_, err := config.CheckLogLevel(logLevel)
		if err != nil {
			return err
		}
		c.Env.LogLevel = logLevel
	} else {
		c.Env.LogLevel = "-"
	}

	c.Log.Errorf("[ Env ] Read ...................................................................... [ OK ]")
	return nil
}

// ReadConfig Read Configuration File By Viper
func (c *Context) ReadConfig() error {
	c.CM = config.ConfigManager{}.New(c.Log.Logger)
	// Write Config File Info
	configPath := c.Env.Home + "/configs"
	configName := c.Env.Profile
	configType := "yaml"
	// Config file struct
	conf := new(config.Configuration)
	// Read
	if err := c.CM.ReadConfig(configPath, configName, configType, conf); err != nil {
		return err
	}
	// Save
	c.Conf = conf

	c.Log.Errorf("[ Configuration ] Read ............................................................ [ OK ]")
	return nil
}

// SetLogger set log level, log output. and etc
func (c *Context) SetLogger() error {
	if c.Env.LogLevel != "-" {
		c.Conf.Log.Level = c.Env.LogLevel
	}
	if err := c.Log.Setting(&c.Conf.Log, c.Env.Home); err != nil {
		return err
	}
	c.Log.Start()
	return nil
}

// InitDatastore Initialize datastore
func (c *Context) InitDatastore() error {
	// Create datastore
	ds, err := db.DataStore{}.New(c.Env.Home, c.Log.Logger)
	if err != nil {
		return err
	}
	// Connect
	if err := ds.Connect(&c.Conf.Datastore); err != nil {
		return err
	}

	// Migrate (모델을 파라미터로 전달)
	if err := ds.Migrate(persistence.SampleModel()); err != nil {
		return err
	}

	c.Datastore = ds
	c.Log.Errorf("[ DataStore ] Initialze ........................................................... [ OK ]")
	return nil
}

// InitRouter Initialize router
func (c *Context) InitRouter() error {
	// init echo framework
	r, err := apphttp.Init(c.Log.Logger, c.Conf.Server.Debug)
	if err != nil {
		return err
	}
	c.Router = r
	c.Log.Errorf("[ Router ] Initialze .............................................................. [ OK ]")
	return nil
}

// Initialize env/config load and sub module init
func Initialize() (*Context, error) {
	c := new(Context)
	c.Conf = new(config.Configuration)

	// 환경 변수, 컨피그를 읽어 들이는 과정에서 로그 출력을 위해
	// 아주 기초적인 부분만을 초기화 한다.
	c.InitLog()

	// Env
	if err := c.ReadEnv(); err != nil {
		return nil, err
	}

	// Read Config
	if err := c.ReadConfig(); err != nil {
		return nil, err
	}

	// Setting Log(from env and conf)
	if err := c.SetLogger(); err != nil {
		return nil, err
	}

	// Datastore
	if err := c.InitDatastore(); err != nil {
		return nil, err
	}

	// Echo Framework Init
	if err := c.InitRouter(); err != nil {
		return nil, err
	}

	// Other Module Init

	c.Log.Errorf("[ ALL ] Initialze ................................................................. [ OK ]")
	return c, nil
}

// InitDepencyInjection sub model Dependency injection and path regi to server
func (c *Context) InitDepencyInjection() error {
	injector := bootstrap.Injector{}.New(c.Router, c.Datastore, c.Log)
	injector.Init()
	return nil
}

// StartSubModules Start SubModule And Waiting Signal / Main Loop
func (c *Context) StartSubModules() {
	// Signal
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGHUP, syscall.SIGTERM)
	c.Log.Errorf("[ Signal ] Listener Start ......................................................... [ OK ]")

	// Echo Framework
	echoServerErr := make(chan error)
	listenAddr := fmt.Sprintf("%s:%d", c.Conf.Server.Host, c.Conf.Server.Port)
	go func() {
		if err := c.Router.Run(listenAddr); err != nil {
			echoServerErr <- err
		}
	}()
	c.Log.Errorf("[ Router ] Listener Start ......................................................... [ OK ]")

	// Start Other Sub Modules

	for {
		select {
		case err := <-echoServerErr:
			c.Log.Errorf("[ SERVER ] ERROR[ %s ]", err.Error())
			c.StopSubModules()
			return
		case sig := <-signalChannel:
			c.Log.Errorf("[ SIGNAL ] Receive [ %s ]", sig.String())
			c.StopSubModules()
			return
		case <-time.After(time.Second * 5):
			c.Log.Errorf("I'm Alive...")
		}
	}
}

// StopSubModules Stop Submodules
func (c *Context) StopSubModules() {
	if err := c.Datastore.Shutdown(); err != nil {
		c.Log.Errorf("[ DataStore ] Shutdown .......................................................... [ Fail ]")
		c.Log.Errorf("[ ERROR ] %s", err.Error())
	} else {
		c.Log.Errorf("[ DataStore ] Shutdown ............................................................ [ OK ]")
	}

	c.Router.Shutdown()
	c.Log.Errorf("[ Router ] Shutdown ............................................................... [ OK ]")

	// 사용하는 서브 모듈(Goroutine)들이 안전하게 종료 될 수 있도록 종료 코드를 추가한다.
}

// @title Golang Web Template API
// @version 1.0.0
// @description This is a golang web template server.

// @contact.name API Support
// @contact.url http://mobigen.com
// @contact.email jblim@mobigen.com

// @host localhost:8080
// @BasePath /
func main() {
	// Initialize Sub module And Read Env, Config
	c, err := Initialize()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// Initialization Interconnection of WebServer Layer
	// controller - application - domain - repository - infrastructures
	c.InitDepencyInjection()

	// Start sub module and main loop
	c.StartSubModules()

	// Bye Bye
	c.Log.Shutdown()
}
