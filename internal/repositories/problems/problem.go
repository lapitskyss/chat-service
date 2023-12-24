package problemsrepo

import (
	"github.com/lapitskyss/chat-service/internal/store"
	"github.com/lapitskyss/chat-service/internal/types"
)

type Problem struct {
	ID     types.ProblemID
	ChatID types.ChatID
}

func adaptProblem(p *store.Problem) Problem {
	return Problem{
		p.ID,
		p.ChatID,
	}
}

func adaptProblems(pp []*store.Problem) []Problem {
	result := make([]Problem, len(pp))
	for i, p := range pp {
		result[i] = adaptProblem(p)
	}
	return result
}
