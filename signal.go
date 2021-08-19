package signal

import (
	"os"
	"os/signal"
)

func WaitForTermination(afterTermination ...func() error) error {
	intChan := make(chan os.Signal, 1)
	signal.Notify(intChan, _sigINT, _sigTERM)
	<-intChan

	for _, f := range afterTermination {
		if err := f(); err != nil {
			return err
		}
	}

	return nil
}
