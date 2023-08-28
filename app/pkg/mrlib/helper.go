package mrlib

import (
    "context"
    "fmt"
    "io"
    "os"
    "os/signal"
    "print-shop-back/pkg/mrapp"
    "print-shop-back/pkg/mrerr"
    "syscall"
)

type Helper struct {
    logger mrapp.Logger
}

func NewHelper(logger mrapp.Logger) *Helper {
    return &Helper{
        logger: logger,
    }
}

func (h *Helper) ExitOnError(err error) {
    if err != nil {
        h.logger.Fatal(err)
    }
}

func (h *Helper) Close(c io.Closer) {
    err := c.Close()

    if err != nil {
        h.logger.Error(mrerr.ErrInternalFailedToClose.Caller(1).Wrap(err, fmt.Sprintf("%v", c)))
    }
}

func (h *Helper) GracefulShutdown(cancel context.CancelFunc) {
    signalAppChan := make(chan os.Signal, 1)

    signal.Notify(
        signalAppChan,
        syscall.SIGABRT,
        syscall.SIGQUIT,
        syscall.SIGHUP,
        os.Interrupt,
        syscall.SIGTERM,
    )

    signalApp := <-signalAppChan
    h.logger.Info("Application shutdown, signal: " + signalApp.String())
    cancel()
}
