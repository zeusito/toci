package actions

import "context"

type DefaultActions struct {
}

func NewDefaultActions() Service {
	return &DefaultActions{}
}

func (s *DefaultActions) SendOTPByEmail(ctx context.Context, code, toEmail string) {

}
