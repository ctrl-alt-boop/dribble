package request

import (
	"fmt"

	"github.com/ctrl-alt-boop/dribble/database"
)

var _ database.Response = (*Response)(nil)
var _ database.Response = (*ChainResponse)(nil)

type (
	Response struct {
		Status        Status
		RequestID     int64
		RequestTarget string
		Body          any
		Error         error
	}
	BatchResponse []*Response
	ChainResponse []*Response
)

func (r Response) Code() int {
	return int(r.Status)
}

func (r Response) Message() string {
	return fmt.Sprintf("%d: %s", r.Code(), r.Status.String())
}

func (c *ChainResponse) Code() int {
	for _, r := range *c {
		if r.Error != nil {
			return int(ErrorChainExecute)
		}
	}
	return int(SuccessChainExecute)
}

func (c *ChainResponse) Message() string {
	return fmt.Sprintf("%d: %s", c.Code(), Status(c.Code()))
}

func (c *BatchResponse) Code() int {
	for _, r := range *c {
		if r.Error != nil {
			return int(ErrorChainExecute)
		}
	}
	return int(SuccessChainExecute)
}

func (c *BatchResponse) Message() string {
	return fmt.Sprintf("%d: %s", c.Code(), Status(c.Code()))
}
