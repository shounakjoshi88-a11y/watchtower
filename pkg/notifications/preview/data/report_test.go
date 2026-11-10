package data

import "testing"

// A preview always delivers its notification, so even when every container it
// shows was updated the session-level suppression must not kick in - otherwise
// `watchtower --preview` could be silenced by the very report it is previewing.
func TestAllUpdatedNoNotifyStaysFalse(t *testing.T) {
	pb := New()
	pb.AddFromState(UpdatedState)
	pb.AddFromState(UpdatedState)

	report := pb.Report()
	if report == nil {
		t.Fatal("expected a report after adding container states")
	}

	if got, want := len(report.Updated()), 2; got != want {
		t.Fatalf("updated containers = %d, want %d", got, want)
	}

	if report.AllUpdatedNoNotify() {
		t.Fatal("AllUpdatedNoNotify() = true, but a preview always delivers its notification")
	}
}
