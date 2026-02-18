package user_booking

import (
	"sync/atomic"

	"go.temporal.io/sdk/workflow"
)

func UserBooking(ctx workflow.Context, user_id int64) error {
	status := "Finding"
	var rejected atomic.Int32
	var accepted atomic.Bool

	rejectedSignal := workflow.GetSignalChannel(ctx, "rejected")
	acceptedSignal := workflow.GetSignalChannel(ctx, "accepted")
	cancelSignal := workflow.GetSignalChannel(ctx, "cancel")
	selector := workflow.NewSelector(ctx)
	workflow.SetQueryHandler(ctx, "status", func() (string, error) {
		return status, nil
	})
	workflow.SetQueryHandler(ctx, "get_user_id", func() (int64, error) {
		return user_id, nil
	})
	selector.AddReceive(rejectedSignal, func(c workflow.ReceiveChannel, more bool) {
		var rejected_value bool
		c.Receive(ctx, &rejected_value)
		if rejected_value {
			rejected.Add(1)
		}
	})
	selector.AddReceive(acceptedSignal, func(c workflow.ReceiveChannel, more bool) {
		var accepted_value int
		c.Receive(ctx, &accepted_value)
		accepted.Swap(true)
		status = "Accepted"
	})
	selector.AddReceive(cancelSignal, func(c workflow.ReceiveChannel, more bool) {
		var cancel_value bool
		c.Receive(ctx, &cancel_value)
		if cancel_value {
			status = "Cancelled"
		}
	})

	for {
		selector.Select(ctx)
		if accepted.Load() && status != "Cancelled" {
			break
		}
	}
	return nil
}
