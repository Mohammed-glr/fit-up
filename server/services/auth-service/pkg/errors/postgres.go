package errors

import "errors"

var (
	ErrFailedToPingPostgress = errors.New("failed to ping postgres database")
	ErrFailedToConnectPostgress = errors.New("failed to connect to postgres database")
	ErrFailedToCreatePostgressPool = errors.New("failed to create postgres connection pool")
	ErrFailedToClosePostgressPool = errors.New("failed to close postgres connection pool")
	ErrFailedToExecutePostgressQuery = errors.New("failed to execute postgres query")
	ErrFailedToPreparePostgressStatement = errors.New("failed to prepare postgres statement")
	ErrFailedToBeginPostgressTransaction = errors.New("failed to begin postgres transaction")
	ErrFailedToCommitPostgressTransaction = errors.New("failed to commit postgres transaction")
	ErrFailedToRollbackPostgressTransaction = errors.New("failed to rollback postgres transaction")
	ErrFailedToScanPostgressRow = errors.New("failed to scan postgres row")
	ErrFailedToGetPostgressRows = errors.New("failed to get postgres rows")
)
