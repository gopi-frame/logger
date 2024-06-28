package contract

type Driver interface {
	Open(map[string]any) (Logger, error)
}
