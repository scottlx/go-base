package design_pattern

import "fmt"

type Event struct {
	Notifier Notifier
	Msg      string
}

type Observer interface {
	UpdateOnEvent(event *Event)
}

type Notifier interface {
	Attach(observer Observer)
	Detach(observer Observer)
	Notify(event *Event)
}

type Speaker struct {
	workers []Observer
}

func (s *Speaker) Attach(observer Observer) {
	s.workers = append(s.workers, observer)
}

func (s *Speaker) Detach(observer Observer) {
	for i, worker := range s.workers {
		if worker == observer {
			s.workers = append(s.workers[:i], s.workers[i+1:]...)
			break
		}
	}
}

func (s *Speaker) Notify(event *Event) {
	for _, worker := range s.workers {
		worker.UpdateOnEvent(event)
	}
}

type WorkerTom struct {
}

func (w *WorkerTom) UpdateOnEvent(event *Event) {
	if event.Msg == "pipeline1" {
		fmt.Println(" Tom is processing product in pipeline1")
		fmt.Println("......done")
		fmt.Println("send product to pipeline2")
		event.Notifier.Notify(&Event{
			Notifier: event.Notifier,
			Msg:      "pipeline2",
		})
	}
}

type WorkerJerry struct {
}

func (w *WorkerJerry) UpdateOnEvent(event *Event) {
	if event.Msg == "pipeline2" {
		fmt.Println(" Jerry is processing product in pipeline2")
		fmt.Println("......done")
		fmt.Println("send product to pipeline3")
		event.Notifier.Notify(&Event{
			Notifier: event.Notifier,
			Msg:      "pipeline3",
		})
	}
}

type Packer struct {
}

func (w *Packer) UpdateOnEvent(event *Event) {
	if event.Msg == "pipeline3" {
		fmt.Println(" Packer is packing up product in pipeline3")
		fmt.Println("......done")
	}
}
