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

	ErrTransactionNotFound = fmt.Errorf("transaction not found")

	ErrUserNotFound = fmt.Errorf("user not found")
	ErrUserExists   = fmt.Errorf("user is exists")
	ErrUserInUse    = fmt.Errorf("the user has outstanding transactions")

	ErrCvScanNotFound       = fmt.Errorf("cv_scan not found")
	ErrCvScanDetailNotFound = fmt.Errorf("cv_scan detail not found")

	ErrEmptyFields = errors.New("empty fields")

	ErrNothingToChange = fmt.Errorf("nothing to change")

	ErrMLServiceNonOK  = errors.New("ML service returned non-OK HTTP status")
	ErrMLServiceDecode = errors.New("failed to decode ML service response")
)

func Wrap(msg string, err error) error {
	return fmt.Errorf("%s: %w", msg, err)
}

func WrapWithFunc(funcName, msg string, err error) error {
	return Wrap(fmt.Sprintf("[%s] %s", funcName, msg), err)
}
