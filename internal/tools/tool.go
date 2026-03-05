package tools

type Tool interface {
	Name() string
	Description() string
	Run(input string) (string, error)
}
