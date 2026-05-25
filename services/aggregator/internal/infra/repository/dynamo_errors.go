package repository

import (
	"aggregator/internal/usecase/apperrors"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func mapDynamoWriteError(err error) error {
	if err == nil {
		return nil
	}
	if isConditionalCheckFailure(err) {
		return apperrors.ErrEventAlreadyProcessed
	}
	return err
}

func isConditionalCheckFailure(err error) bool {
	var txErr *types.TransactionCanceledException
	if !errors.As(err, &txErr) {
		return false
	}

	for _, reason := range txErr.CancellationReasons {
		if aws.ToString(reason.Code) == "ConditionalCheckFailed" {
			return true
		}
	}
	return false
}

func mapDynamoConditionalUpdateError(err error) error {
	if err == nil {
		return nil
	}

	var conditionalErr *types.ConditionalCheckFailedException
	if errors.As(err, &conditionalErr) {
		return nil
	}

	return err
}
