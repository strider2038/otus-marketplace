package messaging

import (
	"context"
	"encoding/json"
	"fmt"

	"billing-service/internal/billing"

	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	"github.com/strider2038/pkg/persistence"
)

type CreateAccrual struct {
	ID          uuid.UUID `json:"id"`
	UserID      uuid.UUID `json:"userId"`
	Amount      float64   `json:"amount"`
	Commission  float64   `json:"commission"`
	Description string    `json:"description"`
}

type AccrualApproved struct {
	ID         uuid.UUID `json:"id"`
	Amount     float64   `json:"amount"`
	Commission float64   `json:"commission"`
}

func (p AccrualApproved) Name() string {
	return "Billing/AccrualApproved"
}

type CreateAccrualProcessor struct {
	accounts           billing.AccountRepository
	operations         billing.OperationRepository
	broker             *billing.BrokerAccount
	transactionManager persistence.TransactionManager
	dispatcher         Dispatcher
}

func NewCreateAccrualProcessor(
	accounts billing.AccountRepository,
	operations billing.OperationRepository,
	broker *billing.BrokerAccount,
	transactionManager persistence.TransactionManager,
	dispatcher Dispatcher,
) *CreateAccrualProcessor {
	return &CreateAccrualProcessor{
		accounts:           accounts,
		operations:         operations,
		broker:             broker,
		transactionManager: transactionManager,
		dispatcher:         dispatcher,
	}
}

func (p *CreateAccrualProcessor) Process(ctx context.Context, message []byte) error {
	var command CreateAccrual
	err := json.Unmarshal(message, &command)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal CreateAccrual command")
	}

	err = p.transactionManager.DoTransactionally(ctx, func(ctx context.Context) error {
		accrual := billing.NewAccrual(command.ID, command.UserID, command.Amount, command.Description)

		return p.addAccrual(ctx, accrual, command.Commission)
	})
	if err != nil {
		return err
	}

	return nil
}

func (p *CreateAccrualProcessor) addAccrual(ctx context.Context, accrual *billing.Operation, commission float64) error {
	account, err := p.accounts.FindByIDForUpdate(ctx, accrual.AccountID)
	if err != nil {
		return errors.WithMessagef(err, `failed to find user account by id "%s"`, accrual.AccountID)
	}

	account.Amount += accrual.Amount - commission

	err = p.operations.Add(ctx, accrual)
	if err != nil {
		return errors.WithMessagef(err, `failed to add accrual "%s" for user "%s"`, accrual.ID, accrual.AccountID)
	}

	err = p.accounts.Save(ctx, account)
	if err != nil {
		return errors.WithMessagef(err, `failed to save user account "%s"`, account.ID)
	}

	err = p.broker.ChargeCommission(
		ctx,
		commission,
		fmt.Sprintf("commission charge from user %s (accrual %s)", account.ID, accrual.ID),
	)
	if err != nil {
		return errors.WithMessagef(err, "failed to charge commission to broker")
	}

	err = p.dispatcher.Dispatch(ctx, AccrualApproved{
		ID:         accrual.ID,
		Amount:     accrual.Amount,
		Commission: commission,
	})
	if err != nil {
		return errors.WithMessage(err, "failed to dispatch AccrualApproved event")
	}

	return nil
}
