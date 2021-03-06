/*
 * Billing service
 *
 * No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)
 *
 * API version: 1.0.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package api

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"strconv"
	"time"

	"billing-service/internal/billing"

	"github.com/bsm/redislock"
	"github.com/gofrs/uuid"
	"github.com/muonsoft/validation"
	"github.com/pkg/errors"
	"github.com/strider2038/pkg/persistence"
)

// BillingApiService is a service that implents the logic for the BillingApiServicer
// This service should implement the business logic for every endpoint for the BillingApi API.
// Include any external packages or services that will be required by this service.
type BillingApiService struct {
	accounts           billing.AccountRepository
	operations         billing.OperationRepository
	transactionManager persistence.TransactionManager
	validator          *validation.Validator
	locker             *redislock.Client
	lockTimeout        time.Duration
}

// NewBillingApiService creates a default api service
func NewBillingApiService(
	accounts billing.AccountRepository,
	operations billing.OperationRepository,
	transactionManager persistence.TransactionManager,
	validator *validation.Validator,
	locker *redislock.Client,
	lockTimeout time.Duration,
) BillingApiServicer {
	return &BillingApiService{
		accounts:           accounts,
		operations:         operations,
		transactionManager: transactionManager,
		validator:          validator,
		locker:             locker,
		lockTimeout:        lockTimeout,
	}
}

// GetBillingAccount -
func (s *BillingApiService) GetBillingAccount(ctx context.Context, id uuid.UUID) (ImplResponse, string, error) {
	account, err := s.accounts.FindByID(ctx, id)
	if err != nil {
		return Response(http.StatusInternalServerError, nil), "", err
	}

	key := makeAccountIdempotenceKey(account)

	return Response(http.StatusOK, account), key, nil
}

// DepositMoney -
func (s *BillingApiService) DepositMoney(ctx context.Context, operation BillingOperation) (ImplResponse, error) {
	if operation.AccountID.IsNil() {
		return newUnauthorizedResponse(), nil
	}
	if operation.IdempotenceKey == "" {
		return newPreconditionRequiredResponse(), nil
	}

	err := s.validator.ValidateValidatable(ctx, operation)
	if err != nil {
		return newUnprocessableEntityResponse(err.Error()), nil
	}

	lock, err := s.locker.Obtain(ctx, operation.AccountID.String()+":deposit", s.lockTimeout, nil)
	if errors.Is(err, redislock.ErrNotObtained) {
		return newConflictResponse(), nil
	}
	if err != nil {
		return newErrorResponsef(err, "failed to obtain lock")
	}
	defer lock.Release(ctx)

	err = s.transactionManager.DoTransactionally(ctx, func(ctx context.Context) error {
		account, err := s.accounts.FindByIDForUpdate(ctx, operation.AccountID)
		if err != nil {
			return errors.WithMessagef(err, "failed to find account %s", operation.AccountID)
		}

		if operation.IdempotenceKey != makeAccountIdempotenceKey(account) {
			return errOutdated
		}

		account.Amount += operation.Amount

		err = s.operations.Add(ctx, billing.NewDeposit(
			operation.AccountID,
			operation.Amount,
			"money deposit by user",
		))
		if err != nil {
			return errors.WithMessagef(err, "failed to add account %s operation", operation.AccountID)
		}

		err = s.accounts.Save(ctx, account)
		if err != nil {
			return errors.WithMessagef(err, "failed to save account %s", operation.AccountID)
		}

		return nil
	})
	if errors.Is(err, errOutdated) {
		return newPreconditionFailedResponse(), nil
	}
	if err != nil {
		return newErrorResponsef(err, "failed to deposit")
	}

	return Response(http.StatusNoContent, nil), nil
}

// WithdrawMoney -
func (s *BillingApiService) WithdrawMoney(ctx context.Context, operation BillingOperation) (ImplResponse, error) {
	if operation.AccountID.IsNil() {
		return newUnauthorizedResponse(), nil
	}
	if operation.IdempotenceKey == "" {
		return newPreconditionRequiredResponse(), nil
	}

	err := s.validator.ValidateValidatable(ctx, operation)
	if err != nil {
		return newUnprocessableEntityResponse(err.Error()), nil
	}

	lock, err := s.locker.Obtain(ctx, operation.AccountID.String()+":withdraw", s.lockTimeout, nil)
	if errors.Is(err, redislock.ErrNotObtained) {
		return newConflictResponse(), nil
	}
	if err != nil {
		return newErrorResponsef(err, "failed to obtain lock")
	}
	defer lock.Release(ctx)

	err = s.transactionManager.DoTransactionally(ctx, func(ctx context.Context) error {
		account, err := s.accounts.FindByIDForUpdate(ctx, operation.AccountID)
		if err != nil {
			return errors.WithMessagef(err, "failed to find account %s", operation.AccountID)
		}

		if operation.IdempotenceKey != makeAccountIdempotenceKey(account) {
			return errOutdated
		}

		account.Amount -= operation.Amount
		if account.Amount < 0 {
			return billing.ErrNotEnoughMoney
		}

		err = s.operations.Add(ctx, billing.NewWithdrawal(
			operation.AccountID,
			operation.Amount,
			"money withdrawal by user",
		))
		if err != nil {
			return errors.WithMessagef(err, "failed to add account %s operation", operation.AccountID)
		}

		err = s.accounts.Save(ctx, account)
		if err != nil {
			return errors.WithMessagef(err, "failed to save account %s", operation.AccountID)
		}

		return nil
	})
	if errors.Is(err, errOutdated) {
		return newPreconditionFailedResponse(), nil
	}
	if errors.Is(err, billing.ErrNotEnoughMoney) {
		return newUnprocessableEntityResponse("Not enough money on the account"), nil
	}
	if err != nil {
		return newErrorResponsef(err, "failed to withdraw")
	}

	return Response(http.StatusNoContent, nil), nil
}

// GetBillingOperations -
func (s *BillingApiService) GetBillingOperations(ctx context.Context, accountID uuid.UUID) (ImplResponse, error) {
	if accountID.IsNil() {
		return newUnauthorizedResponse(), nil
	}

	operations, err := s.operations.FindByAccount(ctx, accountID)
	if err != nil {
		return newErrorResponsef(err, "failed to find account")
	}

	return Response(http.StatusOK, operations), nil
}

func makeAccountIdempotenceKey(account *billing.Account) string {
	hash := sha256.Sum256([]byte(account.ID.String() + ":" + strconv.FormatInt(account.UpdatedAt.UnixNano(), 10)))

	return hex.EncodeToString(hash[:])
}
