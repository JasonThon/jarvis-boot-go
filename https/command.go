package https

type Command interface {
	Check() error
}
