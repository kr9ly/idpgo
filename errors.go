package idpgo

type BadRequestError string

func (e BadRequestError) Error() string {
	return string(e)
}