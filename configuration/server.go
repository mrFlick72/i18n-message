package configuration

import (
	"context"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/logger"
	"github.com/kataras/iris/v12/middleware/recover"
	"github.com/rsocket/rsocket-go"
	"github.com/rsocket/rsocket-go/payload"
	"github.com/rsocket/rsocket-go/rx/mono"
	"github/mrflick72/i18n-message/internal/heath"
	"log"
	"sync"
)

func newWebServer() *iris.Application {
	app := iris.New()
	app.Use(recover.New())
	app.Use(logger.New())
	return app
}

func NewActuatorServer(wg *sync.WaitGroup) {
	app := newWebServer()
	endpoints := heath.HealthEndpoint{}
	endpoints.ResgisterEndpoints(app)
	app.Listen(manager.GetConfigFor("ACTUATOR_PORT"))
	wg.Done()
}

func NewApplicationServer(wg *sync.WaitGroup) {
	app := newWebServer()
	messageRepository := ConfigureMessageRepository()
	ConfigureMessageEndpoints(messageRepository, app)

	// Listen and serve on 0.0.0.0:8080
	app.Listen(manager.GetConfigFor("PRODUCTION_PORT"))
	wg.Done()
}

func NewARSocketServer(wg *sync.WaitGroup) {
	log.Printf("now rsoket listening on %v", manager.GetConfigFor("RSOCKET_PORT"))
	err := rsocket.Receive().
		Acceptor(func(setup payload.SetupPayload, sendingSocket rsocket.CloseableRSocket) (rsocket.RSocket, error) {
			// bind responder
			return rsocket.NewAbstractSocket(
				rsocket.RequestResponse(func(msg payload.Payload) mono.Mono {
					return mono.Just(msg)
				}),
			), nil
		}).
		Transport(rsocket.TCPServer().SetAddr(manager.GetConfigFor("RSOCKET_PORT")).Build()).
		Serve(context.Background())
	log.Fatalln(err)
	wg.Done()
}
