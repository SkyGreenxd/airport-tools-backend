package e

import "fmt"

var (
	ErrStationNotFound = fmt.Errorf("station not found")
	ErrStationExists   = fmt.Errorf("station is exists")

	ErrStoreNotFound = fmt.Errorf("store not found")
	ErrStoreExists   = fmt.Errorf("store exists")

	ErrLocationNotFound = fmt.Errorf("location not found")
	ErrLocationExists   = fmt.Errorf("location is exists")

	ErrToolNotFound = fmt.Errorf("tool not found")
	ErrToolExists   = fmt.Errorf("tool is exists")

	ErrToolTypeNotFound = fmt.Errorf("tool type not found")
	ErrToolTypeExists   = fmt.Errorf("tool type exists")
	ErrToolTypeIsUsed   = fmt.Errorf("tool type is used")

	ErrTransactionNotFound = fmt.Errorf("transaction not found")

	ErrTransactionToolNotFound = fmt.Errorf("there is no such entry in the table")
	ErrTransactionToolExists   = fmt.Errorf("such an instrument already exists in the transaction")

	ErrUserNotFound = fmt.Errorf("user not found")
	ErrUserExists   = fmt.Errorf("user is exists")
	ErrUserInUse    = fmt.Errorf("the user has outstanding transactions")

	ErrNothingToChange = fmt.Errorf("nothing to change")
	ErrIncorrectDate   = fmt.Errorf("incorrect date")
)

func Wrap(msg string, err error) error {
	return fmt.Errorf("%s: %w", msg, err)
}

func WrapWithFunc(funcName, msg string, err error) error {
	return Wrap(fmt.Sprintf("[%s] %s", funcName, msg), err)
}
