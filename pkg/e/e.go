package e

import (
	"errors"
	"fmt"
)

var (
	ErrToolTypeNotFound = fmt.Errorf("tool type not found")
	ErrToolTypeExists   = fmt.Errorf("tool type exists")
	ErrToolTypeIsUsed   = fmt.Errorf("tool type is used")

	ErrToolSetNotFound = fmt.Errorf("tool set not found")
	ErrToolSetExists   = fmt.Errorf("tool set exists")

	ErrTransactionNotFound    = fmt.Errorf("transaction not found")
	ErrTransactionUnfinished  = fmt.Errorf("you have an unfinished issue")
	ErrTransactionAllFinished = fmt.Errorf("you have no open or pending transactions")
	ErrTransactionLimit       = fmt.Errorf("3 unsuccessful scan attempts. Data sent for QA review")
	ErrTransactionCheckQA     = fmt.Errorf("You cannot get new tools while you are being verified QA")

	ErrUserNotFound     = fmt.Errorf("user not found")
	ErrUserExists       = fmt.Errorf("user is exists")
	ErrUserInUse        = fmt.Errorf("the user has outstanding transactions")
	ErrUserRoleNotFound = fmt.Errorf("role not found")

	ErrCvScanNotFound       = fmt.Errorf("cv_scan not found")
	ErrCvScanDetailNotFound = fmt.Errorf("cv_scan detail not found")

	ErrNothingToChange = fmt.Errorf("nothing to change")

	ErrMLServiceNonOK  = errors.New("ML service returned non-OK HTTP status")
	ErrMLServiceDecode = errors.New("failed to decode ML service response")

	ErrInvalidRequestBody = errors.New("invalid request body")
)

func Wrap(msg string, err error) error {
	return fmt.Errorf("%s: %w", msg, err)
}

func WrapWithFunc(funcName, msg string, err error) error {
	return Wrap(fmt.Sprintf("[%s] %s", funcName, msg), err)
}
