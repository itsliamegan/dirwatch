package watch

import (
	"time"

	"github.com/itsliamegan/watch/changeset"
	"github.com/itsliamegan/watch/snapshot"
)

const pollInterval = 100 * time.Millisecond

func Start(dir string) (chan changeset.Change, chan error) {
	changes := make(chan changeset.Change)
	errs := make(chan error)

	go func() {
		current, err := snapshot.Take(dir)
		if err != nil {
			errs <- err
		}

		for {
			new, err := snapshot.Take(dir)
			if err != nil {
				errs <- err
			}

			diff := snapshot.Compare(current, new)
			for _, change := range diff.All() {
				changes <- change
			}

			current = new
			time.Sleep(pollInterval)
		}
	}()

	return changes, errs
}
