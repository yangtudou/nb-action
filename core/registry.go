package core

type Registry struct {
	actions map[string]Action
}

func NewRegistry() *Registry {

	return &Registry{
		actions: make(map[string]Action),
	}
}

func (r *Registry) Register(items ...Action) {

	for _, action := range items {
		r.actions[action.Name()] = action
	}
}

func (r *Registry) GetAction(name string) (Action, bool) {

	action, ok := r.actions[name]

	return action, ok
}

func (r *Registry) List() []Action {
	list := make([]Action, 0, len(r.actions))

	for _, action := range r.actions {
		list = append(list, action)
	}

	return list
}
