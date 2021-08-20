package signal

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/pkg/errors"
)

func WaitForTermination(ctx context.Context, afterTermination ...func() error) (err error) {
	sigChan := make(chan os.Signal, 1)

	signal.Notify(sigChan, _sigINT, _sigTERM)
	defer signal.Stop(sigChan)

	select {
	case sig := <-sigChan:
		err = Error{sig}
	case <-ctx.Done():
		err = ctx.Err()
	}

	if err != nil {
		return errors.Wrap(err, "listening to closing chan")
	}

	for _, f := range afterTermination {
		if err = f(); err != nil {
			return err
		}
	}

	return nil
}

type Error struct{ os.Signal }

func (e Error) Error() string { return fmt.Sprintf("receiving signal %s", e.Signal) }