package repo

import (
	"context"
	"time"

	"github.com/twinj/uuid"
)

func NewAnswer(text, qUUID string, isRight bool) (string, error) {
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
		`insert into answers (uuid, created_at, updated_at, is_right, text, question_uuid) 
	values  ($1, $2, $3, $4, $5, $6);
	`,
		aUUID,
		time.Now(),
		time.Now(),
		isRight,
		text,
		qUUID,
	)
	if err != nil {
		return "", err
	}
	return aUUID, nil
}
