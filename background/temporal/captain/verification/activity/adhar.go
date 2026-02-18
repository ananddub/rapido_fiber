package activity_verfication

import (
	"context"
	"database/sql"
	"fmt"

	"encore.app/gen/pgdb"
)

func (a *ActivityVerfication) VerifyAdhar(ctx context.Context, captainID int32) (bool, error) {
	verification, err := a.conn.Query.GetCaptainVerificationByCaptainID(ctx, captainID)
	if err != nil {
		fmt.Printf("Failed to get verification record for captain %d: %v\n", captainID, err)
		return false, err
	}

	_, err = a.conn.Query.UpdateAadharStatus(ctx, pgdb.UpdateAadharStatusParams{
		VerificationID: verification.ID,
		Status: pgdb.NullDocumentStatusEnum{
			DocumentStatusEnum: pgdb.DocumentStatusEnumAPPROVED,
			Valid:              true,
		},
		AdminComment: sql.NullString{String: "Verified by Temporal Workflow", Valid: true},
		VerifiedBy:   sql.NullInt32{Int32: 1, Valid: true}, // Mock Admin ID
	})
	if err != nil {
		fmt.Printf("Failed to update Aadhaar status: %v\n", err)
		return false, err
	}
	fmt.Printf("Marked Aadhaar as verified for captain %d (VerificationID: %d)\n", captainID, verification.ID)
	return true, nil
}
