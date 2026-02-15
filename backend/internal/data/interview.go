package data

import (
	"ai-interview/internal/biz"
	"context"
	"database/sql"
	"encoding/json"

	"github.com/go-kratos/kratos/v2/log"
)

type interviewRepo struct {
	data *Data
	log  *log.Helper
}

// NewInterviewRepo 创建面试仓储
func NewInterviewRepo(data *Data, logger log.Logger) biz.InterviewRepo {
	return &interviewRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *interviewRepo) Create(ctx context.Context, interview *biz.Interview) (*biz.Interview, error) {
	result, err := r.data.db.ExecContext(ctx,
		`INSERT INTO interviews (user_id, title, position, status, language,
			llm_provider, llm_model, tts_provider, tts_voice, resume)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		interview.UserID, interview.Title, interview.Position, interview.Status,
		interview.Language, interview.LLMProvider, interview.LLMModel,
		interview.TTSProvider, interview.TTSVoice, interview.Resume,
	)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	interview.ID = id

	return interview, nil
}

func (r *interviewRepo) GetByID(ctx context.Context, id int64) (*biz.Interview, error) {
	interview := &biz.Interview{}
	err := r.data.db.QueryRowContext(ctx,
		`SELECT id, user_id, title, position, status, language,
			llm_provider, llm_model, tts_provider, tts_voice, resume, created_at, updated_at
		FROM interviews WHERE id = ?`, id,
	).Scan(&interview.ID, &interview.UserID, &interview.Title, &interview.Position,
		&interview.Status, &interview.Language, &interview.LLMProvider, &interview.LLMModel,
		&interview.TTSProvider, &interview.TTSVoice, &interview.Resume,
		&interview.CreatedAt, &interview.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return interview, nil
}

func (r *interviewRepo) ListByUserID(ctx context.Context, userID int64, page, pageSize int) ([]*biz.Interview, int, error) {
	var total int
	err := r.data.db.QueryRowContext(ctx,
		"SELECT COUNT(*) FROM interviews WHERE user_id = ?", userID,
	).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	rows, err := r.data.db.QueryContext(ctx,
		`SELECT id, user_id, title, position, status, language, created_at, updated_at
		FROM interviews WHERE user_id = ? ORDER BY created_at DESC LIMIT ? OFFSET ?`,
		userID, pageSize, offset,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var interviews []*biz.Interview
	for rows.Next() {
		i := &biz.Interview{}
		if err := rows.Scan(&i.ID, &i.UserID, &i.Title, &i.Position, &i.Status,
			&i.Language, &i.CreatedAt, &i.UpdatedAt); err != nil {
			return nil, 0, err
		}
		interviews = append(interviews, i)
	}

	return interviews, total, nil
}

func (r *interviewRepo) UpdateStatus(ctx context.Context, id int64, status string) error {
	_, err := r.data.db.ExecContext(ctx,
		"UPDATE interviews SET status = ? WHERE id = ?", status, id,
	)
	return err
}

func (r *interviewRepo) CreateMessage(ctx context.Context, msg *biz.InterviewMessage) (*biz.InterviewMessage, error) {
	result, err := r.data.db.ExecContext(ctx,
		"INSERT INTO interview_messages (interview_id, role, content) VALUES (?, ?, ?)",
		msg.InterviewID, msg.Role, msg.Content,
	)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	msg.ID = id

	return msg, nil
}

func (r *interviewRepo) ListMessages(ctx context.Context, interviewID int64) ([]*biz.InterviewMessage, error) {
	rows, err := r.data.db.QueryContext(ctx,
		"SELECT id, interview_id, role, content, created_at FROM interview_messages WHERE interview_id = ? ORDER BY id ASC",
		interviewID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []*biz.InterviewMessage
	for rows.Next() {
		m := &biz.InterviewMessage{}
		if err := rows.Scan(&m.ID, &m.InterviewID, &m.Role, &m.Content, &m.CreatedAt); err != nil {
			return nil, err
		}
		messages = append(messages, m)
	}

	return messages, nil
}

func (r *interviewRepo) CreateEvaluation(ctx context.Context, eval *biz.Evaluation) (*biz.Evaluation, error) {
	categoriesJSON, err := json.Marshal(eval.Categories)
	if err != nil {
		categoriesJSON = []byte("[]")
	}

	result, err := r.data.db.ExecContext(ctx,
		`INSERT INTO evaluations (interview_id, overall_score, summary, categories, strengths, weaknesses, suggestions)
		VALUES (?, ?, ?, ?, ?, ?, ?)`,
		eval.InterviewID, eval.OverallScore, eval.Summary,
		string(categoriesJSON), eval.Strengths, eval.Weaknesses, eval.Suggestions,
	)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	eval.ID = id

	return eval, nil
}

func (r *interviewRepo) GetEvaluation(ctx context.Context, interviewID int64) (*biz.Evaluation, error) {
	eval := &biz.Evaluation{}
	var categoriesJSON sql.NullString

	err := r.data.db.QueryRowContext(ctx,
		`SELECT id, interview_id, overall_score, summary, categories, strengths, weaknesses, suggestions, created_at
		FROM evaluations WHERE interview_id = ?`, interviewID,
	).Scan(&eval.ID, &eval.InterviewID, &eval.OverallScore, &eval.Summary,
		&categoriesJSON, &eval.Strengths, &eval.Weaknesses, &eval.Suggestions, &eval.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	if categoriesJSON.Valid {
		_ = json.Unmarshal([]byte(categoriesJSON.String), &eval.Categories)
	}

	return eval, nil
}
