package activity_verfication

import "encore.app/internal/connection"

type ActivityVerfication struct {
	conn *connection.Connections
}

func Init() (*ActivityVerfication, error) {
	conn, err := connection.InitConnection()
	if err != nil {
		return nil, err
	}
	return &ActivityVerfication{
		conn: conn,
	}, nil
}
