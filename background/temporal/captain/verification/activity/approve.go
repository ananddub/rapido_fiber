package activity_verfication

import (
	"context"
	"database/sql"
	"fmt"

	"encore.app/gen/pgdb"
)

func (a *ActivityVerfication) ApproveCaptain(ctx context.Context, captainID int32) (bool, error) {
	verification, err := a.conn.Query.GetCaptainVerificationByCaptainID(ctx, captainID)
	if err != nil {
		fmt.Printf("Failed to get verification record for captain %d: %v\n", captainID, err)
		return false, err
	}

	_, err = a.conn.Query.UpdateCaptainVerificationStatus(ctx, pgdb.UpdateCaptainVerificationStatusParams{
		ID:            verification.ID,
		OverallStatus: pgdb.NullVerificationStatusEnum{VerificationStatusEnum: pgdb.VerificationStatusEnumAPPROVED, Valid: true},
		CurrentStage:  pgdb.NullVerificationStageEnum{VerificationStageEnum: pgdb.VerificationStageEnumFINAL, Valid: true},
	})
	if err != nil {
		fmt.Printf("Failed to update verification status: %v\n", err)
		return false, err
	}

	err = a.conn.Query.UpdateCaptainStatus(ctx, pgdb.UpdateCaptainStatusParams{
		ID:         captainID,
		IsVerified: sql.NullBool{Bool: true, Valid: true},
		IsActive:   sql.NullBool{Bool: true, Valid: true},
	})
	if err != nil {
		fmt.Printf("Failed to update Captain table status: %v\n", err)
		return false, err
	}

	fmt.Printf("Successfully Approved Captain %d (VerificationID: %d)\n", captainID, verification.ID)
	return true, nil
}
