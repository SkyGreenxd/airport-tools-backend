package repository

import "context"

type Station interface {
	Create(ctx context.Context, station *Station) (*Station, error)
	GetById(ctx context.Context, id int64) (*Station, error)
	GetAll(ctx context.Context) ([]*Station, error)
	Delete(ctx context.Context, id int64) error
	Update(ctx context.Context, station *Station) (*Station, error)
}

type Store interface {
	Create(ctx context.Context, station *Store) (*Store, error)
	GetById(ctx context.Context, id int64) (*Store, error)
	GetAll(ctx context.Context) ([]*Store, error)
	Delete(ctx context.Context, id int64) error
	Update(ctx context.Context, station *Store) (*Store, error)
}

type Location interface {
	Create(ctx context.Context, location *Location) (*Location, error)
	GetById(ctx context.Context, id int64) (*Location, error)
	GetAll(ctx context.Context) ([]*Location, error)
	Delete(ctx context.Context, id int64) error
	Update(ctx context.Context, location *Location) (*Location, error)
}

type Tool interface {
	Create(ctx context.Context, tool *Tool) (*Tool, error)
	GetById(ctx context.Context, id int64) (*Tool, error)
	GetAll(ctx context.Context) ([]*Tool, error)
	Delete(ctx context.Context, id int64) error
	Update(ctx context.Context, tool *Tool) (*Tool, error)
}

type ToolType interface {
	Create(ctx context.Context, toolType *ToolType) (*ToolType, error)
	GetById(ctx context.Context, id int64) (*ToolType, error)
	GetAll(ctx context.Context) ([]*ToolType, error)
	Delete(ctx context.Context, id int64) error
	Update(ctx context.Context, toolType *ToolType) (*ToolType, error)
}

type Transaction interface {
	Create(ctx context.Context, station *Transaction) (*Transaction, error)
	GetById(ctx context.Context, id int64) (*Transaction, error)
	GetAll(ctx context.Context) ([]*Transaction, error)
	Delete(ctx context.Context, id int64) error
	Update(ctx context.Context, station *Transaction) (*Transaction, error)
}

type TransactionTool interface {
	Create(ctx context.Context, station *TransactionTool) (*TransactionTool, error)
	GetById(ctx context.Context, id int64) (*TransactionTool, error)
	GetByTransactionId(ctx context.Context, transactionId int64) ([]*TransactionTool, error)
	Update(ctx context.Context, station *TransactionTool) (*TransactionTool, error)
	Delete(ctx context.Context, id int64) error
	GetAll(ctx context.Context) ([]*TransactionTool, error)
	GetUnreturnedByTransactionID(ctx context.Context, transactionId int64) ([]*TransactionTool, error)
}

type User interface {
	Create(ctx context.Context, user *User) (*User, error)
	GetById(ctx context.Context, id int64) (*User, error)
	GetAll(ctx context.Context) ([]*User, error)
	Delete(ctx context.Context, id int64) error
	Update(ctx context.Context, user *User) (*User, error)
}
