package gowebservice

import (
	"github.com/corneldamian/golog"
	"github.com/corneldamian/httpway"
	"github.com/corneldamian/httpwaymid"
	"github.com/kardianos/service"

	"errors"
	"log"
	"os"
	"path/filepath"
	"time"
)

func DoInit(cfg WebServerConfiger,
	serviceConfig *service.Config,
	middlewareFactory func(*httpway.Router) map[string]*httpway.Router,
	routeFactory map[string]func(*httpway.Router)) (service.Service, error) {

	p := &program{
		middlewareFactory: middlewareFactory,
		routeFactory:      routeFactory,
	}

	s, err := service.New(p, serviceConfig)
	if err != nil {
		log.Fatal(err)
	}

	logger, err := s.Logger(nil)
	if err != nil {
		log.Fatal(err)
	}
	p.systemLog = logger

	if err := initConfig(cfg); err != nil {
		p.systemLog.Error(err)
		return nil, err
	}

	if err := p.initLogger(); err != nil {
		p.systemLog.Error(err)
		return nil, err
	}

	return s, nil
}

type program struct {
	server    *httpway.Server
	systemLog service.Logger
	router    *httpway.Router

	middlewareFactory func(*httpway.Router) map[string]*httpway.Router
	routeFactory      map[string]func(*httpway.Router)
}

func (p *program) Start(srv service.Service) error {
	p.router = httpway.New()

	c := GetConfig().WebServerConfig()
	if c.Enable_Access_Log {
		accessLogger := golog.GetLogger("access").Info
		p.router = p.router.Middleware(httpwaymid.AccessLog(accessLogger))
	}

	if golog.ToLogLevel(c.Log_Level) != golog.DEBUG {
		p.router = p.router.Middleware(httpwaymid.PanicCatcher)
	}

	p.router.NotFound = httpwaymid.NotFound(p.router)
	p.router.MethodNotAllowed = httpwaymid.MethodNotAllowed(p.router)

	p.router.SessionManager = httpwaymid.NewSessionManager()
	p.router.Logger = golog.GetLogger("general")

	middlewares := p.middlewareFactory(p.router)

	for middlewareName, routeFactory := range p.routeFactory {
		middleware, ok := middlewares[middlewareName]
		if !ok {
			panic("No middleware chain registered for " + middlewareName)
		}

		routeFactory(middleware)
	}

	if s, err := p.initWebServer(); err != nil {
		return err
	} else {
		p.server = s
	}

	return nil
}

func (p *program) Stop(s service.Service) error {
	if err := p.server.Stop(); err != nil {
		return err
	}

	if err := p.server.WaitStop(5 * time.Second); err != nil {
		return err
	}
	if err := golog.Stop(2 * time.Second); err != nil {
		return err
	}

	return nil
}

func (p *program) initWebServer() (*httpway.Server, error) {

	server := httpway.NewServer(nil)
	server.Addr = ":8080"
	server.Handler = p.router

	if err := server.Start(); err != nil {
		return nil, err
	}

	return server, nil
}

func (p *program) initLogger() error {
	c := GetConfig().WebServerConfig()

	f, err := os.Stat(c.Log_Dir)
	if err == nil {
		if !f.IsDir() {
			return errors.New("Log path exists but it's file not dir")
		}
	} else {
		if os.IsNotExist(err) {
			if err := os.MkdirAll(c.Log_Dir, 0755); err != nil {
				return err
			}
		} else {
			return err
		}
	}



	systemLogFile := filepath.Join(c.Log_Dir, c.System_Log)
	verbosity := golog.LDefault | golog.LHeaderFooter
	if golog.ToLogLevel(c.Log_Level) == golog.DEBUG {
		verbosity = verbosity | golog.LFile
	}
	golog.NewLogger("general", systemLogFile, &golog.LoggerConfig{
		Level:     golog.ToLogLevel(c.Log_Level),
		Verbosity: verbosity,
	})

	if c.Enable_Access_Log {
		accessLogFile := filepath.Join(c.Log_Dir, c.Access_Log)
		golog.NewLogger("access", accessLogFile, &golog.LoggerConfig{
			MessageQueueSize: 100000,
			Level:            golog.INFO,
			Verbosity:        golog.LHeaderFooter,
			HeaderWriter:     httpwaymid.AccessLogHeaderWriter,
		})
	}

	return nil
}
