package common

type CheckError struct {
	Code    string
	Message string
}

func (err CheckError) Error() string {
	return err.Message
}
