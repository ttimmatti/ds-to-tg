package tg_msgs

import (
	"errors"
	"fmt"
	"log"

	errror "github.com/ttimmatti/discord-tg_parser/errors"
)

func HandleErrors(errs []error) error {
	for _, err := range errs {

		var ierr *errror.Error
		if errors.As(err, &ierr) {

			if ierr.Code() == errror.ErrorCodeFailure {

				// notify
				errMsg := DefaultReply(
					ADMIN_ID,
					fmt.Sprintf("I'm down! Fatal Error: %s", err),
					"",
				)
				if err1 := sendMsg(*errMsg); err1 != nil {
					log.Printf("[VERY FATAL!!!] Couldn't send errMsg to admin :: %s", err)
				} else {
					log.Printf("[INFO] Notified admin!")
				}

				return fmt.Errorf("FATAL [ERROR] :: %w", err)
			}
		}

		log.Printf("[WARN] Err: %s; Code: %d", err, ierr.Code())

	}

	return nil
}
