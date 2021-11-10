package repo

import (
	"context"
	"time"

	"github.com/twinj/uuid"
)

type UserInfo struct {
	CreatedAt time.Time `json:"created_at"`
	Name      string    `json:"name"`
	Phone     string    `json:"phone"`
	Email     string    `json:"email"`
	Link      string    `json:"link"`
	Question  string    `json:"question"`
}

func NewUserInfo(name, phone, email, link, q string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), poolAcquireTimeoutMS*time.Millisecond)
	conn, err := connPool.Acquire(ctx)
	if err != nil {
		cancel()
		return "", err
	}
	cancel()
	defer conn.Release()
	ctx, cancel = context.WithTimeout(context.Background(), repoReqTimeoutMS*time.Millisecond)
	defer cancel()
	aUUID := uuid.NewV4().String()
	_, err = conn.Exec(ctx,
		`insert into user_infos (uuid, created_at, updated_at, name, phone, email, link, question)
	values ($1, $2, $3, $4, $5, $6, $7, $8);
	`,
		aUUID,
		time.Now(),
		time.Now(),
		name,
		phone,
		email,
		link,
		q,
	)
	if err != nil {
		return "", err
	}
	return aUUID, nil
}

func GetUserInfoList(l, o int) ([]*UserInfo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), poolAcquireTimeoutMS*time.Millisecond)
	conn, err := connPool.Acquire(ctx)
	if err != nil {
		cancel()
		return nil, err
	}
	cancel()
	defer conn.Release()
	ctx, cancel = context.WithTimeout(context.Background(), repoReqTimeoutMS*time.Millisecond)
	defer cancel()
	rws, err := conn.Query(ctx,
		`select created_at, name, phone, email, link, question  from user_infos 
		order by uuid 
		limit $1 
		offset $2;
		`,
		l,
		o,
	)
	if err != nil {
		return nil, err
	}
	defer rws.Close()

	r := make([]*UserInfo, 0)
	for rws.Next() {
		ui := new(UserInfo)

		err := rws.Scan(
			&ui.CreatedAt,
			&ui.Name,
			&ui.Phone,
			&ui.Email,
			&ui.Link,
			&ui.Question,
		)
		if err != nil {
			return nil, err
		}
		r = append(r, ui)
	}

	return r, rws.Err()
}
