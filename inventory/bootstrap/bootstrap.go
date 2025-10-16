package bootstrap

import "go.uber.org/dig"

func NewBootstrap() *dig.Container {
	c := dig.New()
	c = NewApplication(c)
	c = NewRepository(c)
	c = NewHandler(c)
	c = NewService(c)

	return c
}
