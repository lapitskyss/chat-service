package problemsrepo

import (
	"github.com/lapitskyss/chat-service/internal/store"
	"github.com/lapitskyss/chat-service/internal/types"
)

type Problem struct {
	ID        types.ProblemID
	ChatID    types.ChatID
	ManagerID types.UserID
}

func adaptProblem(p *store.Problem) Problem {
	return Problem{
		ID:        p.ID,
		ChatID:    p.ChatID,
		ManagerID: p.ManagerID,
	}
}

func adaptProblems(pp []*store.Problem) []Problem {
	result := make([]Problem, len(pp))
	for i, p := range pp {
		result[i] = adaptProblem(p)
	}
	return result
}
