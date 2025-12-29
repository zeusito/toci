package actions

import "context"

type Service interface {
	SendOTPByEmail(ctx context.Context, code, toEmail string)
}
