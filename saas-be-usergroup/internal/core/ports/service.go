package ports

import (
	"context"
	"github.com/saas-be-usergroup/internal/core/domain/auth"
)

type (
	AuthService interface {
		RegisterBeforeWithEmail(ctx context.Context, in auth.RegisterBeforeWithEmail) (*auth.RegisterBeforeResponse, error)
		ConfirmationRegister(ctx context.Context, in auth.ConfirmationRegister) (*auth.ConfirmationResponse, error)
	}
)
