package session

import (
	"testing"
	"time"
)

// TestDelegateRestsDone checks which rests of a delegated child count as the
// task finished. A Claude Code child that ends its turn to wait on its own
// background shell is done on its hooks, but has not finished: it wakes when
// the job exits, so its parent must not be mailed yet.
func TestDelegateRestsDone(t *testing.T) {
	started := int64(1000) // unix ms
	after := (started + 1) * int64(time.Millisecond)
	for _, tc := range []struct {
		name string
		w    WindowState
		want bool
	}{
		{"done after the prompt", WindowState{AgentState: AgentStateDone, AgentStateAt: after}, true},
		{"idle after the prompt", WindowState{AgentState: AgentStateIdle, AgentStateAt: after}, true},
		{"the rest it was in when typed to", WindowState{AgentState: AgentStateDone, AgentStateAt: started * int64(time.Millisecond)}, false},
		{"silent", WindowState{AgentState: AgentStateWorking, AgentReason: AgentReasonSilent, AgentStateAt: after}, false},
		{"done on a background shell", WindowState{AgentState: AgentStateDone, AgentStateAt: after, AgentWaitsOn: "1 background shell"}, false},
		{"idle on a background shell", WindowState{AgentState: AgentStateIdle, AgentStateAt: after, AgentWaitsOn: "1 background shell"}, false},
		{"errored with a shell left", WindowState{AgentState: AgentStateErrored, AgentStateAt: after, AgentWaitsOn: "1 background shell"}, true},
	} {
		if got := delegateRestsDone(tc.w, started); got != tc.want {
			t.Errorf("%s: delegateRestsDone = %v, want %v", tc.name, got, tc.want)
		}
	}
}
