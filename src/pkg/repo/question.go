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
//		select *
//		from questions
//		where questions.difficulty_level = [%d]
//		order by random()
//		limit [%d]
//	) as T1
//	join answers as T2 on T1.uuid = T2.question_uuid
const getRandomQFirstPart = "select T1.uuid as q_uuid, T1.text as q_text, T2.uuid as a_uuid, T2.text as a_text, T2.is_right from ( select * from questions where questions.difficulty_level = "
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

		results := make([]*Result, 0)

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
			var r = new(Result)
			err = rws.Scan(
				&r.QUUID,
				&r.QText,
				&r.AUUID,
				&r.AText,
				&r.IsRight,
			)
			if err != nil {
				return nil, err
			}
			results = append(results, r)
		}

		if err != rws.Err() {
			return nil, err
		}

		m = addResultsToMap(results, m)
	} else {
		// TODO: make it
		return nil, errDistribution
	}

	return deleteFromMapQuestionWithoutAnswer(m), nil
}

func addResultsToMap(results []*Result, m map[string]*RQuestion) map[string]*RQuestion {
	for j := range results {
		r := results[j]
		if _, ok := m[r.QUUID]; ok {
			// double right answer check
			canAppend := (r.IsRight && !m[r.QUUID].HasRightAnswer) || !r.IsRight
			if canAppend {
				m[r.QUUID].Answers = append(m[r.QUUID].Answers,
					RAnswer{
						Text:    r.AText,
						UUID:    r.AUUID,
						IsRight: r.IsRight,
					})
				if r.IsRight {
					m[r.QUUID].HasRightAnswer = true
				}
			}
		} else {
			m[r.QUUID] = &RQuestion{
				HasRightAnswer: r.IsRight,
				Text:           r.QText,
				Answers: []RAnswer{
					{
						Text:    r.AText,
						UUID:    r.AUUID,
						IsRight: r.IsRight,
					},
				},
			}
		}
	}

	return m
}

func deleteFromMapQuestionWithoutAnswer(m map[string]*RQuestion) map[string]*RQuestion {
	for k := range m {
		if !m[k].HasRightAnswer {
			delete(m, k)
		}
	}

	return m
}
