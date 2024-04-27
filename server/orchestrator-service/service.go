package orchestratorService

type CommandRequest struct {
	Org    string `json:"org"`
	Repo   string `json:"repo"`
	Env    string `json:"env"`
	B64Str string `json:"b64String"`
}
