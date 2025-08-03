package protocol

type ErrInvalidData struct {
	msg string
}

func (e ErrInvalidData) Error() string {
	return e.msg
}

func NewErrInvalidData(msg string) ErrInvalidData {
	return ErrInvalidData{msg}
}

type ErrInvalidCommand struct {
	msg string
}

func (e ErrInvalidCommand) Error() string {
	return e.msg
}

func NewErrInvalidCommand(msg string) ErrInvalidCommand {
	return ErrInvalidCommand{msg}
}
