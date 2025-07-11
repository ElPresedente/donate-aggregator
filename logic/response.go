package logic

type ResponseName string

// вместо enum для LogicResponse.name
const (
	RouletteSpin ResponseName = "roulette-spin-result"
)

type LogicResponse struct {
	name ResponseName
	data any
}
