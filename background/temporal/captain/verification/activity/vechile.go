package activity_verfication

import (
	"context"
	"database/sql"
	"fmt"

	"encore.app/gen/pgdb"
)

func (a *ActivityVerfication) VerifyVehicle(ctx context.Context, captainID int32) (bool, error) {
	verification, err := a.conn.Query.GetCaptainVerificationByCaptainID(ctx, captainID)
	if err != nil {
		fmt.Printf("Failed to get verification record for captain %d: %v\n", captainID, err)
		return false, err
	}

	vehicles, err := a.conn.Query.GetVehiclesByVerificationID(ctx, verification.ID)
	if err != nil {
		fmt.Printf("Failed to get vehicles for verification %d: %v\n", verification.ID, err)
		return false, err
	}

	for _, v := range vehicles {
		_, err := a.conn.Query.UpdateVehicleStatus(ctx, pgdb.UpdateVehicleStatusParams{
			ID: v.ID,
			Status: pgdb.NullDocumentStatusEnum{
				DocumentStatusEnum: pgdb.DocumentStatusEnumAPPROVED,
				Valid:              true,
			},
			AdminComment: sql.NullString{String: "Verified by Temporal Workflow", Valid: true},
			VerifiedBy:   sql.NullInt32{Int32: 1, Valid: true},
		})
		if err != nil {
			fmt.Printf("Failed to update vehicle %d status: %v\n", v.ID, err)
			return false, err
		}
	}

	fmt.Printf("Marked %d Vehicles as verified for captain %d\n", len(vehicles), captainID)
	return true, nil
}
