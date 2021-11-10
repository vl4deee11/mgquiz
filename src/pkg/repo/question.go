package repo

import (
	"context"
	"errors"
	"magnusquiz/pkg/enum"
	"strconv"
	"strings"
	"time"

	"github.com/twinj/uuid"
)

func NewQuestion(text string, lvl enum.DifficultyLevel) (string, error) {
	if lvl < enum.EASY {
		lvl = enum.EASY
	}
	if lvl > enum.HARD {
		lvl = enum.HARD
	}
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
		`insert into questions (uuid, created_at, updated_at, text, difficulty_level)
	values ($1, $2, $3, $4, $5);
	`,
		aUUID,
		time.Now(),
		time.Now(),
		text,
		lvl,
	)
	if err != nil {
		return "", err
	}
	return aUUID, nil
}

type Result struct {
	QUUID   string
	AUUID   string
	QText   string
	AText   string
	IsRight bool
}

// Query template
//	select
//		T1.uuid as q_uuid,
//		T1.text as q_text,
//		T2.uuid as a_uuid,
//		T2.text as a_text,
//		T2.is_right
//	from (
//		select
//			uuid,
//			text
//		from questions
//		where questions.difficulty_level = [%d]
//		order by random()
//		limit [%d]
//	) as T1
//	join answers as T2 on T1.uuid = T2.question_uuid
const getRandomQFirstPart = "select T1.uuid as q_uuid, T1.text as q_text, T2.uuid as a_uuid, T2.text as a_text, T2.is_right from ( select uuid, text from questions where questions.difficulty_level = "
const getRandomQSecondPart = " order by random() limit "
const getRandomQThirdPart = ") as T1 join answers as T2 on T1.uuid = T2.question_uuid"

var errDistribution = errors.New("it is not possible to generate questions, normally balanced between difficulty levels")

func GetNRandomQuestions(n int) (map[string]*RQuestion, error) {
	m := map[string]*RQuestion{}

	if n%enum.CountDifficultyLevels == 0 {
		batch := n / enum.CountDifficultyLevels
		queries := make([]string, enum.CountDifficultyLevels)
		batchS := strconv.Itoa(batch)
		for i := 1; i <= enum.CountDifficultyLevels; i++ {
			iS := strconv.Itoa(i)
			var sb strings.Builder
			sb.Grow(len(getRandomQFirstPart) + len(getRandomQSecondPart) + len(getRandomQThirdPart) + len(iS) + len(batchS))
			sb.WriteString(getRandomQFirstPart)
			sb.WriteString(iS)
			sb.WriteString(getRandomQSecondPart)
			sb.WriteString(batchS)
			sb.WriteString(getRandomQThirdPart)
			queries[i-1] = sb.String()
		}

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
		rws, err := conn.Query(ctx, strings.Join(queries, " UNION "))
		if err != nil {
			return nil, err
		}
		defer rws.Close()
		for rws.Next() {
			var (
				qUUID   string
				qText   string
				aUUID   string
				aText   string
				isRight bool
			)
			err = rws.Scan(
				&qUUID,
				&qText,
				&aUUID,
				&aText,
				&isRight,
			)
			if err != nil {
				return nil, err
			}

			// add to map of results
			// possible no right answer for question if it's not in db
			if rQuestion, ok := m[qUUID]; ok {
				// double right answer check
				canAppend := (isRight && !rQuestion.HasRightAnswer) || !isRight
				if canAppend {
					rQuestion.Answers = append(rQuestion.Answers,
						RAnswer{
							Text:    aText,
							UUID:    aUUID,
							IsRight: isRight,
						})
					if isRight {
						rQuestion.HasRightAnswer = true
					}
				}
			} else {
				m[qUUID] = &RQuestion{
					HasRightAnswer: isRight,
					Text:           qText,
					Answers: []RAnswer{
						{
							Text:    aText,
							UUID:    aUUID,
							IsRight: isRight,
						},
					},
				}
			}
		}

		if err != rws.Err() {
			return nil, err
		}
	} else {
		// TODO: make it
		return nil, errDistribution
	}

	return m, nil
}
