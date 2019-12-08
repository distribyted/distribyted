package command

import "github.com/sirupsen/logrus"

const (
	ServerDescription = "Server Description"
	ServerHelp        = "Help text here"
)

type Server struct {
}

func (c *Server) Execute(args []string) error {
	logrus.Info("starting server command")

	return nil
}
