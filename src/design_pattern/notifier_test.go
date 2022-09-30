package design_pattern

import "testing"

func TestNotifier(t *testing.T) {
	t.Run("testNotifier", func(t *testing.T) {
		tom := &WorkerTom{}
		jerry := &WorkerJerry{}
		p := &Packer{}
		sp := Speaker{}
		sp.Attach(tom)
		sp.Attach(jerry)
		sp.Attach(p)
		sp.Notify(&Event{
			Notifier: &sp,
			Msg:      "pipeline1",
		})
	})
}
