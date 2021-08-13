package workflow

import (
	"errors"
	"fmt"
)

type (
	Place      uint
	Places     []Place
	Transition struct {
		Name   string
		From   Place
		To     Place
		Action func() error
	}
	Transitions  []*Transition
	StateMachine struct {
		places      Places
		transitions Transitions
		current     Place
	}
	Equaler interface {
		Equal(Equaler) bool
	}
)

func NewStateMachine(places Places, transitions Transitions) (*StateMachine, error) {
	if len(places) == 0 {
		return nil, errors.New("no places given")
	}
	sm := &StateMachine{
		places:      places,
		transitions: transitions,
		current:     places[0],
	}
	if err := sm.Validate(); err != nil {
		return nil, err
	}
	return sm, nil
}

// Apply tries to find the Transition with the given name and applies it to the machine.
func (s *StateMachine) Apply(n string) error {
	t := s.transitions.Find(n)
	if t == nil {
		return fmt.Errorf(`transition with name "%s" does not exist`, n)
	}
	// Check if the Transition can be applied.
	if !s.Can(t) {
		return fmt.Errorf(`cannot apply transition "%v" on "%v"`, t, s.current)
	}
	if t.Action != nil {
		if err := t.Action(); err != nil {
			return err
		}
	}
	s.current = t.To
	return nil
}

// Can checks if the given Transition can be applied.
func (s StateMachine) Can(t *Transition) bool {
	return t.From == s.current
}

// CurrentPlace returns the place this state machine currently is in.
func (s StateMachine) CurrentPlace() Place {
	return s.current
}

// AdjacentPlaces returns all reachable places for the given place.
func (s StateMachine) AdjacentPlaces(p Place) (Places, error) {
	// If the Place p does not exist return and error.
	if !s.places.Contains(p) {
		return nil, fmt.Errorf(`place does not exist: "%v"`, p)
	}
	m := make(map[Place]struct{})
	for _, t := range s.transitions {
		if t.From == p {
			m[t.To] = struct{}{}
		}
	}
	var ps Places
	for p := range m {
		ps = append(ps, p)
	}
	return ps, nil
}

// AllowedTransitions returns all transitions that can be applied on the current state.
func (s StateMachine) AllowedTransitions() Transitions {
	var ts Transitions
	for _, t := range s.transitions {
		if t.From == s.current {
			ts = append(ts, t)
		}
	}
	return ts
}

// DFS executes a Depth-First-Search on the StateMachine and adds every visited node to the given map.
func (s StateMachine) DFS(v map[Place]struct{}, p Place) error {
	v[p] = struct{}{}
	ps, err := s.AdjacentPlaces(p)
	if err != nil {
		return err
	}
	for _, p := range ps {
		if _, ok := v[p]; !ok {
			if err := s.DFS(v, p); err != nil {
				return err
			}
		}
	}
	return nil
}

// Validate checks the following:
// - places are unique.
// - transitions are unique.
// - start exists.
// - every place is reachable.
func (s StateMachine) Validate() error {
	// Check place uniqueness.
	up := make(map[Place]struct{})
	for _, p := range s.places {
		if _, ok := up[p]; ok {
			return fmt.Errorf(`validate: duplicate place: "%v"`, p)
		}
		up[p] = struct{}{}
	}
	// Check transition validity and uniqueness.
	tp := make(map[string]struct{})
	for _, t := range s.transitions {
		if t.Name == "" {
			return fmt.Errorf(`validate: transition has no name: "%v"`, t)
		}
		if !s.places.Contains(t.From) {
			return fmt.Errorf(`validate: transition from does not exist: "%v"`, t.From)
		}
		if !s.places.Contains(t.To) {
			return fmt.Errorf(`validate: transition to does not exist: "%v"`, t.To)
		}
		if _, ok := tp[t.Name]; ok {
			return fmt.Errorf(`validate: duplicate transition: "%v"`, t)
		}
		tp[t.Name] = struct{}{}
	}
	// Every place is reachable from start. Done with DFS.
	v := make(map[Place]struct{})
	if err := s.DFS(v, s.current); err != nil {
		return fmt.Errorf("validate: %w", err)
	}
	if len(s.places) != len(v) {
		var u Place
		for _, p := range s.places {
			if _, ok := v[p]; !ok {
				u = p
				break
			}
		}
		return fmt.Errorf(`validate: unreachable place detected: "%v"`, u)
	}
	return nil
}

// Contains checks if the given place does exist.
func (ps Places) Contains(p Place) bool {
	for _, e := range ps {
		if e == p {
			return true
		}
	}
	return false
}

func (ts Transitions) Find(n string) *Transition {
	for _, t := range ts {
		if t.Name == n {
			return t
		}
	}
	return nil
}
