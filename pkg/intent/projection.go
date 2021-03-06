package intent

import (
	"github.com/bottlerocket-os/bottlerocket-update-operator/pkg/marker"
	"github.com/pkg/errors"
)

var nextLinear map[marker.NodeAction]marker.NodeAction = map[marker.NodeAction]marker.NodeAction{
	// Stabilization from known points.
	"":                         marker.NodeActionStabilize,
	marker.NodeActionStabilize: marker.NodeActionStabilize,
	marker.NodeActionUnknown:   marker.NodeActionStabilize,

	// Linear progression
	marker.NodeActionReset:         marker.NodeActionStabilize,
	marker.NodeActionPrepareUpdate: marker.NodeActionPerformUpdate,
	marker.NodeActionPerformUpdate: marker.NodeActionRebootUpdate,
	// FIN. The actor must know what to do next to bring itself around again if
	// that's what's appropriate.
	marker.NodeActionRebootUpdate: marker.NodeActionRebootUpdate,
}

// TODO: add tests for the expected state machine turns.

func calculateNext(action marker.NodeAction) (marker.NodeAction, error) {
	// TODO: resolve next state if applicable
	next, ok := nextLinear[action]
	if !ok {
		return marker.NodeActionUnknown, errors.Errorf("no next action from %q, resolving as unknown", action)
	}
	return next, nil
}

// FallbackNodeAction is the first recovery action that is reasonble to take
// from an unknown point.
var FallbackNodeAction = func() marker.NodeAction {
	n, _ := calculateNext(marker.NodeActionUnknown)
	return n
}()
