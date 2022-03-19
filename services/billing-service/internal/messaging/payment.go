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

type CreatePayment struct {
	ID          uuid.UUID `json:"id"`
	UserID      uuid.UUID `json:"userId"`
	Amount      float64   `json:"amount"`
	Commission  float64   `json:"commission"`
	Description string    `json:"description"`
}

type PaymentSucceeded struct {
	ID         uuid.UUID `json:"id"`
	Amount     float64   `json:"amount"`
	Commission float64   `json:"commission"`
}

func (p PaymentSucceeded) Name() string {
	return "Billing/PaymentSucceeded"
}

type PaymentDeclined struct {
	ID         uuid.UUID `json:"id"`
	Amount     float64   `json:"amount"`
	Commission float64   `json:"commission"`
	Reason     string    `json:"reason"`
}

func (p PaymentDeclined) Name() string {
	return "Billing/PaymentDeclined"
}

type CreatePaymentProcessor struct {
	accounts           billing.AccountRepository
	operations         billing.OperationRepository
	broker             *billing.BrokerAccount
	transactionManager persistence.TransactionManager
	dispatcher         Dispatcher
}

func NewCreatePaymentProcessor(
	accounts billing.AccountRepository,
	operations billing.OperationRepository,
	broker *billing.BrokerAccount,
	transactionManager persistence.TransactionManager,
	dispatcher Dispatcher,
) *CreatePaymentProcessor {
	return &CreatePaymentProcessor{
		accounts:           accounts,
		operations:         operations,
		broker:             broker,
		transactionManager: transactionManager,
		dispatcher:         dispatcher,
	}
}

func (p *CreatePaymentProcessor) Process(ctx context.Context, message []byte) error {
	var command CreatePayment
	err := json.Unmarshal(message, &command)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal CreatePayment command")
	}

	err = p.transactionManager.DoTransactionally(ctx, func(ctx context.Context) error {
		payment := billing.NewPayment(command.ID, command.UserID, command.Amount, command.Description)

		return p.addPayment(ctx, payment, command.Commission)
	})
	if err != nil {
		return err
	}

	return nil
}

func (p *CreatePaymentProcessor) addPayment(ctx context.Context, payment *billing.Operation, commission float64) error {
	account, err := p.accounts.FindByIDForUpdate(ctx, payment.AccountID)
	if err != nil {
		return errors.WithMessagef(err, `failed to find user account by id "%s"`, payment.AccountID)
	}

	account.Amount -= payment.Amount + commission
	if account.Amount < 0 {
		err = p.dispatcher.Dispatch(ctx, PaymentDeclined{
			ID:         payment.ID,
			Amount:     payment.Amount,
			Commission: commission,
			Reason:     "not enough money",
		})
		if err != nil {
			return errors.WithMessage(err, "failed to dispatch PaymentDeclined event")
		}

		return nil
	}

	err = p.operations.Add(ctx, payment.WithCommission(commission))
	if err != nil {
		return errors.WithMessagef(err, `failed to add payment "%s" for user "%s"`, payment.ID, payment.AccountID)
	}

	err = p.accounts.Save(ctx, account)
	if err != nil {
		return errors.WithMessagef(err, `failed to save user account "%s"`, account.ID)
	}

	err = p.broker.ChargeCommission(
		ctx,
		commission,
		fmt.Sprintf("commission charge from user %s (payment %s)", account.ID, payment.ID),
	)
	if err != nil {
		return errors.WithMessagef(err, "failed to charge commission to broker")
	}

	err = p.dispatcher.Dispatch(ctx, PaymentSucceeded{
		ID:         payment.ID,
		Amount:     payment.Amount,
		Commission: commission,
	})
	if err != nil {
		return errors.WithMessage(err, "failed to dispatch PaymentSucceeded event")
	}

	return nil
}
