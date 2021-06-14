package httparty

import (
	"fmt"
	"httparty/collections"
	"httparty/engine"
	"httparty/validator"
)

type Httparty struct {
	executor Executor
	store    *collections.OrderedMap
	validator validator.Validator
}

type HttpartyOption func(*Httparty)

func NewHttparty(opts ...HttpartyOption) *Httparty {
	executor := DefaultHttpExecutor()
	store := collections.NewOrderedMap()

	hp := &Httparty{
		executor: executor,
		store:    store,
		validator: validator.NewHttpValidator(),
	}

	for _, optfn := range opts {
		optfn(hp)
	}

	return hp
}

func WithExecutor(executor Executor) HttpartyOption {
	return func(hp *Httparty) {
		hp.executor = executor
	}
}

func WithStore(store *collections.OrderedMap) HttpartyOption {
	return func(hp *Httparty) {
		hp.store = store
	}
}

func (hp *Httparty) RunAll(fleets []engine.Fleet) ([]*Httparty, error) {
	response := make(chan *Httparty, len(fleets))
	errc := make(chan error)

	results := []*Httparty{}

	defer close(response)
	defer close(errc)

	for _, fleet := range fleets {
		go func(fleet engine.Fleet) {
			ehp, err := hp.RunIt(fleet)
			if err != nil {
				errc <- err
				return
			}

			response <- ehp
		}(fleet)
	}

	for {
		select {
		case err := <-errc:
			return nil, err
		case result := <-response:
			results = append(results, result)
			if len(results) == len(fleets) {
				return results, nil
			}
		}
	}
}

func (hp *Httparty) RunIt(fleet engine.Fleet) (*Httparty, error) {
	if err := hp.validator.Validate(fleet); err != nil {
		return nil, err
	}

	graph := engine.MakeGraph(fleet.Steps)
	stepper := engine.MakeStepMap(fleet.Steps)

	final := fleet.Steps[fleet.ExitPoint]

	steps := engine.GetStepsInOrder(stepper, engine.Dfs(graph, final.Name))

	for _, step := range steps {
		rstep := step.RenderHeadersAndBody(hp.store)
		resp, err := hp.executor.RunIt(rstep)
		if err != nil {
			return nil, err
		}

		hp.store.Add(fmt.Sprintf("%s::%s", engine.RespPrefix, step.Name), resp)
	}

	return hp, nil
}

func (hp *Httparty) GetRespValue(key string) interface{} {
	val, _ := hp.store.Get(fmt.Sprintf("%s::%s", engine.RespPrefix, key))

	return val
}
