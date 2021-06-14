package engine

import (
	"httparty/collections"
	"httparty/errorrs"
	"httparty/httputils"
)

func MakeStepMap(steps []*Step) map[string]*Step {
	stepMap := map[string]*Step{}

	for _, step := range steps {
		stepMap[step.Name] = step
	}

	return stepMap
}

func MakeGraph(steps []*Step) map[string][]string {
	graph := map[string][]string{}

	for _, step := range steps {
		if len(step.Requires) == 0 {
			graph[step.Name] = []string{}
		} else {
			graph[step.Name] = append(graph[step.Name], step.Requires...)
		}
	}

	return graph
}

func Dfs(graph map[string][]string, start string) []string {
	stack := collections.NewStringStack(start)
	path := collections.NewStringSet()

	for !stack.Empty() {
		vertex := stack.Pop()
		if path.Contains(vertex) {
			continue
		}

		path.Add(vertex)

		for _, rule := range graph[vertex] {
			stack.Push(rule)
		}
	}

	return path.List()
}

func GetStepsInOrder(stepper map[string]*Step, stepNames []string) []*Step {
	steps := make([]*Step, len(stepNames))

	for i, name := range stepNames {
		steps[i] = stepper[name]
	}

	return steps
}

func GetBodyString(request Request) (string, error) {
	var bodyStr string

	switch request.Type {
	case JSON:
		return httputils.MakeJSONBody(request.Body)
	case FormEncoded:
		return httputils.MakeFormEncodedBody(request.Body), nil
	default:
		return bodyStr, errorrs.UnknownContentType
	}
}
