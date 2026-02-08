package common

import (
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
	"github.com/amanhigh/go-fun/common/util"
	"github.com/amanhigh/go-fun/components/fun-app/handlers"
	"github.com/amanhigh/go-fun/components/fun-app/publisher"
	"github.com/golobby/container/v3"
)

func (fi *FunAppInjector) registerMessagingCore() {
	container.MustSingleton(fi.di, util.NewStdWatermillLogger)
	container.MustSingleton(fi.di, util.NewGoChannel)
	container.MustSingleton(fi.di, providePublisher)
	container.MustSingleton(fi.di, provideSubscriber)
	container.MustSingleton(fi.di, provideRouter)
	container.MustSingleton(fi.di, util.NewWatermillController)
}

func (fi *FunAppInjector) registerCommandHandlers() {
	// DI registrations via providers that return interface types (no Fill)
	container.MustSingleton(fi.di, handlers.NewEnrollmentCommandHandler)
	container.MustSingleton(fi.di, handlers.NewSeatCommandHandler)
	// Provide MessagingServer (depends on handlers) and expose its router for controller
	container.MustSingleton(fi.di, handlers.NewMessagingServer)
}

func (fi *FunAppInjector) registerPublishers() {
	container.MustSingleton(fi.di, publisher.NewBasePublisher)
	container.MustSingleton(fi.di, publisher.NewSeatAllocationPublisher)
	container.MustSingleton(fi.di, publisher.NewEnrollmentPublisher)
}

func providePublisher(channel *gochannel.GoChannel) message.Publisher {
	return channel
}

func provideSubscriber(channel *gochannel.GoChannel) message.Subscriber {
	return channel
}

// provideRouter adapts MessagingServer into a *message.Router for WatermillController.
func provideRouter(srv *handlers.MessagingServer) *message.Router {
	return srv.Router()
}
