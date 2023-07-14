package mrhttp

import (
    "fmt"
    "net"
    "path"
    "path/filepath"
)

const (
    ListenTypeSock = "sock"
    ListenTypePort = "port"
)

type (
    ListenOptions struct {
        AppPath string
        Type string
        SockName string
        BindIP string
        Port string
    }
)

func (s *Server) createListener(opt *ListenOptions) (net.Listener, error) {
    var listener net.Listener
    var listenErr error

    if opt.Type == ListenTypeSock {
        s.logger.Info("Detect app real path")
        appDir, err := filepath.Abs(filepath.Dir(opt.AppPath))

        if err != nil {
            s.logger.Fatal(fmt.Errorf("app real path: %w", err))
        }

        socketPath := path.Join(appDir, opt.SockName)
        s.logger.Info("Listen to unix socket: %s", socketPath)

        listener, listenErr = net.Listen("unix", socketPath)
        s.logger.Info("Server is listening unix socket: %s", socketPath)
    } else if opt.Type == ListenTypePort  {
        addr := fmt.Sprintf("%s:%s", opt.BindIP, opt.Port)
        s.logger.Info("Listen to tcp: %s", addr)

        listener, listenErr = net.Listen("tcp", addr)
        s.logger.Info("Server is listening to port %s", addr)
    } else {
        availableValues := fmt.Sprintf("Available values: %s, %s", ListenTypePort, ListenTypeSock)

        if opt.Type == "" {
            listenErr = fmt.Errorf("listen type is required. %s", availableValues)
        } else {
            listenErr = fmt.Errorf("listen type '%s' is unknown. %s", opt.Type, availableValues)
        }
    }

    return listener, listenErr
}
