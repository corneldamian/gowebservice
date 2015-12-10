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
	routeFactory map[string]func(*httpway.Router),
	bootstrap func(server *httpway.Server, logger *golog.Logger, router *httpway.Router) error) (service.Service, error) {

	p := &program{
		middlewareFactory: middlewareFactory,
		routeFactory:      routeFactory,
		bootstrap:         bootstrap,
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

	bootstrap func(server *httpway.Server, logger *golog.Logger, router *httpway.Router) error

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

	p.server = p.initWebServer()

	p.router.Logger = golog.GetLogger("general")
	p.router.SessionManager = httpwaymid.NewSessionManager(c.Session_Timeout, c.Session_Expiration, golog.GetLogger("general"))

	if err := p.bootstrap(p.server, golog.GetLogger("general"), p.router); err != nil {
		return err
	}

	p.router = p.router.Middleware(httpwaymid.JSONRenderer("jsonData", "statusCode"))

	if c.Template_Dir != "" {
		p.router = p.router.Middleware(httpwaymid.TemplateRenderer(c.Template_Dir, "templateName", "templateData", "statusCode"))
	}

	p.router.NotFound = httpwaymid.NotFound(p.router)
	p.router.MethodNotAllowed = httpwaymid.MethodNotAllowed(p.router)

	middlewares := p.middlewareFactory(p.router)

	for middlewareName, routeFactory := range p.routeFactory {
		middleware, ok := middlewares[middlewareName]
		if !ok {
			panic("No middleware chain registered for " + middlewareName)
		}

		routeFactory(middleware)
	}

	if err := p.server.Start(); err != nil {
		return err
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

func (p *program) initWebServer() *httpway.Server {

	server := httpway.NewServer(nil)
	server.Addr = GetConfig().WebServerConfig().Address
	server.Handler = p.router

	return server
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
		verbosity = verbosity | golog.LFile | golog.LMicroseconds
	}
	golog.NewLogger("general", systemLogFile, &golog.LoggerConfig{
		Level:          golog.ToLogLevel(c.Log_Level),
		Verbosity:      verbosity,
		FileRotateSize: (2 << 23), /*16MB*/
		FileDepth:      5,
	})

	if c.Enable_Access_Log {
		accessLogFile := filepath.Join(c.Log_Dir, c.Access_Log)
		golog.NewLogger("access", accessLogFile, &golog.LoggerConfig{
			MessageQueueSize: 100000,
			Level:            golog.INFO,
			Verbosity:        golog.LHeaderFooter,
			HeaderWriter:     httpwaymid.AccessLogHeaderWriter,
			FileRotateSize:   (2 << 23), /*16MB*/
		})
	}

	return nil
}
