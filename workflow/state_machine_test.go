package workflow_test

import (
	"github.com/masseelch/wapiti/workflow"
	"github.com/stretchr/testify/require"
	"testing"
)

const (
	start workflow.Place = iota
	one
	two
	three
	four
	five
)

func TestPlaces_Contains(t *testing.T) {
	require.False(t, workflow.Places{}.Contains(one))
	require.False(t, workflow.Places{start}.Contains(one))
	require.True(t, workflow.Places{start}.Contains(start))
	require.True(t, workflow.Places{start, one}.Contains(one))
}

func TestNewStateMachine(t *testing.T) {
	s, err := workflow.NewStateMachine(nil, nil)
	require.Error(t, err)
	require.Nil(t, s)

	s, err = workflow.NewStateMachine(workflow.Places{start}, nil)
	require.NoError(t, err)
	require.NotNil(t, s)

	// One not reachable.
	s, err = workflow.NewStateMachine(workflow.Places{start, one}, nil)
	require.Error(t, err)
	require.Nil(t, s)

	// Transition has no name.
	s, err = workflow.NewStateMachine(workflow.Places{start, one}, workflow.Transitions{{From: start, To: one}})
	require.Error(t, err)
	require.Nil(t, s)

	// Transition invalid from.
	s, err = workflow.NewStateMachine(
		workflow.Places{start, one},
		workflow.Transitions{
			{"Init", start, two, nil},
		},
	)
	require.Error(t, err)
	require.Nil(t, s)

	// Transition invalid to.
	s, err = workflow.NewStateMachine(
		workflow.Places{start, one},
		workflow.Transitions{
			{"Init", two, start, nil},
		},
	)
	require.Error(t, err)
	require.Nil(t, s)

	// All okay.
	s, err = workflow.NewStateMachine(
		workflow.Places{start, one, two},
		workflow.Transitions{
			{"Init", start, one, nil},
			{"One-Two", one, two, nil},
			{"Loop", two, two, nil},
		},
	)
	require.NoError(t, err)
	require.NotNil(t, s)
}

func TestStateMachine_AdjacentPlaces(t *testing.T) {
	s, err := workflow.NewStateMachine(
		workflow.Places{start, one, two, three, four, five},
		workflow.Transitions{
			{"Start-One", start, one, nil},
			{"One-Two", one, two, nil},
			{"Loop", two, two, nil},
			{"Two-Three", two, three, nil},
			{"Two-Four", two, four, nil},
			{"Two-Five", two, five, nil},
			{"Five-One", five, one, nil},
		},
	)
	require.NoError(t, err)
	require.NotNil(t, s)

	for p, e := range map[workflow.Place]workflow.Places{
		start: {one},
		one:   {two},
		two:   {two, three, four, five},
		three: nil,
		four:  nil,
		five:  {one},
	} {
		ps, err := s.AdjacentPlaces(p)
		require.NoError(t, err)
		require.Len(t, ps, len(e))
		require.Subset(t, e, ps)
	}
}

func TestStateMachine_AllowedTransitions(t *testing.T) {
	s, err := workflow.NewStateMachine(
		workflow.Places{start, one, two},
		workflow.Transitions{
			{"Start-One", start, one, nil},
			{"Start-Two", start, two, nil},
			{"One-Two", one, two, nil},
		},
	)
	require.NoError(t, err)
	require.NotNil(t, s)

	ts := s.AllowedTransitions()
	require.Len(t, ts, 2)
	require.Subset(t, workflow.Transitions{
		{"Start-One", start, one, nil},
		{"Start-Two", start, two, nil},
	}, ts)
}
