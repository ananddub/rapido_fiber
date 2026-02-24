package background

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"encore.app/gen/pgdb"
	booking_common "encore.app/internal/common/booking/redis"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/redis/go-redis/v9"
)

type Client struct {
	mqtt         mqtt.Client
	redis        *redis.Client
	pg           *pgdb.Queries
	booking_repo *booking_common.Repo
}

const MQTT_CONNECTED_TOPIC = "$SYS/brokers/+/clients/+/connected"
const MQTT_DISCONNECTED_TOPIC = "$SYS/brokers/+/clients/+/disconnected"

func NewClient(mqtt mqtt.Client, rd *redis.Client, pg *pgdb.Queries) *Client {
	b := booking_common.NewRepo(rd)
	b.Init(context.Background())
	return &Client{
		mqtt:         mqtt,
		redis:        rd,
		pg:           pg,
		booking_repo: booking_common.NewRepo(rd),
	}
}

type Payload struct {
	IPAddress      string                 `json:"ipaddress"`
	ReceiveMaximum int                    `json:"receive_maximum"`
	ConnProps      map[string]interface{} `json:"conn_props"`
	ExpiryInterval int                    `json:"expiry_interval"`
	CleanStart     bool                   `json:"clean_start"`
	SockPort       int                    `json:"sockport"`
	ProtoName      string                 `json:"proto_name"`
	ConnectedAt    int64                  `json:"connected_at"`
	ClientID       string                 `json:"clientid"`
	ClientAttrs    map[string]interface{} `json:"client_attrs"`
	ProtoVer       int                    `json:"proto_ver"`
	Username       string                 `json:"username"`
	Ts             int64                  `json:"ts"`
	Protocol       string                 `json:"protocol"`
	Keepalive      int                    `json:"keepalive"`
}

func (p *Payload) GetKey() (string, error) {
	if strings.HasPrefix(p.ClientID, "user:") {
		return p.ClientID, nil
	}
	if strings.HasPrefix(p.ClientID, "captain:") {
		return p.ClientID, nil
	}
	return "", fmt.Errorf("invalid client id: %s", p.ClientID)
}
func (c *Client) ListnerUser(ctx context.Context) {
	c.mqtt.Subscribe(MQTT_CONNECTED_TOPIC, 1, func(client mqtt.Client, msg mqtt.Message) {
		var p Payload
		err := json.Unmarshal(msg.Payload(), &p)
		if err != nil {
			println("Error unmarshalling message: " + err.Error())
			return
		}
		key, err := p.GetKey()
		if err != nil {
			fmt.Printf("Client connected: Invalid Payload: %s, IP: %s payload: %+v\n", p.ClientID, p.IPAddress, p)
			return
		}
		keyint, err := strconv.Atoi(strings.Split(key, ":")[1])
		if err != nil {
			fmt.Printf("Client connected: Invalid Payload: %s, IP: %s payload: %+v\n", p.ClientID, p.IPAddress, p)
			return
		}
		if strings.HasPrefix(p.ClientID, "captain:") {
			user, err := c.pg.GetCaptainById(ctx, int32(keyint))
			if err != nil {
				fmt.Printf("Client connected: User not found: %s, IP: %s payload: %+v\n", p.ClientID, p.IPAddress, p)
				return
			}
			err = c.booking_repo.CreateCaptain(ctx, booking_common.CaptainData{
				Id:       string(user.ID),
				Name:     user.Name,
				Phone:    user.Phone,
				IsBooked: user.IsBlocked.Bool,
				Status:   string(user.Status.UserStatus),
			})
			if err != nil {
				fmt.Printf("Client connected: Error creating user in redis: %s, IP: %s payload: %+v error: %s\n", p.ClientID, p.IPAddress, p, err.Error())
				return
			}
			fmt.Printf("Client connected: Invalid Payload: %s, IP: %s payload: %+v\n", p.ClientID, p.IPAddress, p)
			return
		}
		if !strings.HasPrefix(p.ClientID, "user:") {
			user, err := c.pg.GetUserById(ctx, int32(keyint))
			if err != nil {
				fmt.Printf("Client connected: User not found: %s, IP: %s payload: %+v\n", p.ClientID, p.IPAddress, p)
				return
			}

			err = c.booking_repo.CreateUser(ctx, booking_common.UserData{
				Id:       string(user.ID),
				Name:     user.Name,
				Phone:    user.Phone,
				IsBooked: user.IsBlocked.Bool,
				Status:   string(user.Status.UserStatus),
			})
			if err != nil {
				fmt.Printf("Client connected: Error creating user in redis: %s, IP: %s payload: %+v error: %s\n", p.ClientID, p.IPAddress, p, err.Error())
				return
			}
			fmt.Printf("Client connected: Invalid Payload: %s, IP: %s payload: %+v\n", p.ClientID, p.IPAddress, p)
			return
		}
		fmt.Printf("Client connected: %s, IP: %s\n", p.ClientID, p.IPAddress)
	})

	c.mqtt.Subscribe(MQTT_DISCONNECTED_TOPIC, 1, func(client mqtt.Client, msg mqtt.Message) {
		var p Payload
		err := json.Unmarshal(msg.Payload(), &p)
		if err != nil {
			println("Error unmarshalling message: " + err.Error())
			return
		}
		key, err := p.GetKey()
		if err != nil {
			fmt.Printf("Client connected: Invalid Payload: %s, IP: %s payload: %+v\n", p.ClientID, p.IPAddress, p)
			return
		}
		keyint, err := strconv.Atoi(strings.Split(key, ":")[1])
		if err != nil {
			fmt.Printf("Client connected: Invalid Payload: %s, IP: %s payload: %+v\n", p.ClientID, p.IPAddress, p)
			return
		}
		if strings.HasPrefix(p.ClientID, "captain:") {
			user, err := c.pg.GetCaptainById(ctx, int32(keyint))
			if err != nil {
				fmt.Printf("Client connected: User not found: %s, IP: %s payload: %+v\n", p.ClientID, p.IPAddress, p)
				return
			}
			err = c.booking_repo.UpdateUserStatus(ctx, string(user.ID), "offline")
			if err != nil {
				fmt.Printf("Client connected: Error creating user in redis: %s, IP: %s payload: %+v error: %s\n", p.ClientID, p.IPAddress, p, err.Error())
				return
			}
			fmt.Printf("Client connected: Invalid Payload: %s, IP: %s payload: %+v\n", p.ClientID, p.IPAddress, p)
			return
		}
		if !strings.HasPrefix(p.ClientID, "user:") {
			user, err := c.pg.GetUserById(ctx, int32(keyint))
			if err != nil {
				fmt.Printf("Client connected: User not found: %s, IP: %s payload: %+v\n", p.ClientID, p.IPAddress, p)
				return
			}
			err = c.booking_repo.UpdateCaptainStatus(ctx, string(user.ID), "offline")
			if err != nil {
				fmt.Printf("Client connected: Error creating user in redis: %s, IP: %s payload: %+v error: %s\n", p.ClientID, p.IPAddress, p, err.Error())
				return
			}
			fmt.Printf("Client connected: Invalid Payload: %s, IP: %s payload: %+v\n", p.ClientID, p.IPAddress, p)
			return
		}
		fmt.Printf("Client disconnected: %s, IP: %s\n", p.ClientID, p.IPAddress)
	})
}
