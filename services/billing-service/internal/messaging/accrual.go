package messaging

import (
	"context"
	"encoding/json"

	"billing-service/internal/billing"

	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	"github.com/strider2038/pkg/persistence"
)

type CreateAccrual struct {
	ID          uuid.UUID `json:"id"`
	UserID      uuid.UUID `json:"userId"`
	Amount      float64   `json:"amount"`
	Description string    `json:"description"`
}

type AccrualApproved struct {
	ID uuid.UUID `json:"id"`
}

func (p AccrualApproved) Name() string {
	return "Billing/AccrualApproved"
}

type CreateAccrualProcessor struct {
	accounts           billing.AccountRepository
	operations         billing.OperationRepository
	transactionManager persistence.TransactionManager
	dispatcher         Dispatcher
}

func NewCreateAccrualProcessor(
	accounts billing.AccountRepository,
	operations billing.OperationRepository,
	transactionManager persistence.TransactionManager,
	dispatcher Dispatcher,
) *CreateAccrualProcessor {
	return &CreateAccrualProcessor{
		accounts:           accounts,
		operations:         operations,
		transactionManager: transactionManager,
		dispatcher:         dispatcher,
	}
}

func (processor *CreateAccrualProcessor) Process(ctx context.Context, message []byte) error {
	var command CreateAccrual
	err := json.Unmarshal(message, &command)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal CreateAccrual command")
	}

	accrual := billing.NewAccrual(command.ID, command.UserID, command.Amount, command.Description)
	err = processor.transactionManager.DoTransactionally(ctx, processor.addAccrual(accrual))
	if err != nil {
		return err
	}

	return nil
}

func (processor *CreateAccrualProcessor) addAccrual(accrual *billing.Operation) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		account, err := processor.accounts.FindByIDForUpdate(ctx, accrual.AccountID)
		if err != nil {
			return errors.WithMessagef(err, `failed to find user account by id "%s"`, accrual.AccountID)
		}

		account.Amount += accrual.Amount

		err = processor.operations.Add(ctx, accrual)
		if err != nil {
			return errors.WithMessagef(err, `failed to add accrual "%s" for user "%s"`, accrual.ID, accrual.AccountID)
		}

		err = processor.accounts.Save(ctx, account)
		if err != nil {
			return errors.WithMessagef(err, `failed to save user account "%s"`, account.ID)
		}

		err = processor.dispatcher.Dispatch(ctx, AccrualApproved{ID: accrual.ID})
		if err != nil {
			return errors.WithMessage(err, "failed to dispatch AccrualApproved event")
		}

		return nil
	}
}
