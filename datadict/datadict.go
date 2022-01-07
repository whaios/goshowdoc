package datadict

type DataDict interface {
	Open() error
	Close() error
	Query() ([]*Table, error)
}
