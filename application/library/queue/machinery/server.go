package machinery

import (
	"github.com/RichardKnop/machinery/v1"
	"github.com/RichardKnop/machinery/v1/config"
)

func NewServer(tasks map[string]interface{}, configPaths ...interface{}) (*Server, error) {
	s := &Server{}
	err := s.ParseConfig(configPaths...)
	if err != nil {
		return s, err
	}
	err = s.newServer()
	if err != nil {
		return s, err
	}
	if tasks != nil {
		err = s.RegisterTasks(tasks)
	}
	return s, err
}

type Server struct {
	config *config.Config
	*machinery.Server
}

func (s *Server) ParseConfig(configPaths ...interface{}) (err error) {
	var configPath string
	if len(configPaths) > 0 {
		switch c := configPaths[0].(type) {
		case string:
			configPath = c
		case config.Config:
			s.config = &c
			return
		case *config.Config:
			s.config = c
			return
		}
	}
	if len(configPath) > 0 {
		s.config, err = config.NewFromYaml(configPath, true)
		return
	}

	s.config, err = config.NewFromEnvironment()
	return
}

// newServer Create server instance
func (s *Server) newServer() (err error) {
	s.Server, err = machinery.NewServer(s.config)
	return
}
