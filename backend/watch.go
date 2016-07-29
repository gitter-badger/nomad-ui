package main

import (
	"log"
	"time"

	"github.com/gorilla/websocket"
	"github.com/hashicorp/nomad/api"
)

const (
	allocsEvent = "fetched_allocs"
	allocEvent  = "fetched_alloc"

	evalsEvent = "fetched_evals"
	evalEvent  = "fetched_eval"

	jobsEvent = "fetched_jobs"
	jobEvent  = "fetched_job"

	nodesEvent = "fetched_nodes"
	nodeEvent  = "fetched_node"
)

type event struct {
	Kind    string
	Payload interface{}
}

type watch struct {
	client    *api.Client
	websocket *websocket.Conn

	commandCh  chan struct{}
	shutdownCh chan struct{}
}

func NewWatch(url string, ws *websocket.Conn, timeout time.Duration) *watch {
	config := api.DefaultConfig()
	config.Address = url
	config.WaitTime = timeout

	client, err := api.NewClient(config)
	if err != nil {
		log.Fatalf("Could not create client: %s", err)
	}

	return &watch{
		client:     client,
		websocket:  ws,
		commandCh:  make(chan struct{}),
		shutdownCh: make(chan struct{}),
	}
}

func (w *watch) Start() {

	// Start the default long-running watchs
	go w.Allocs()
	go w.Evals()
	go w.Nodes()
	go w.Jobs()

	for {
		select {
		case command := <-w.commandCh:
			log.Printf("Received command!", command)
		case <-w.shutdownCh:
			log.Printf("Shutting down!")
			return
		}
	}
}

func (w *watch) Stop() {
	close(w.shutdownCh)
}

func (w *watch) Allocs() {
	q := &api.QueryOptions{WaitIndex: 0}

	for {
		select {
		case <-w.shutdownCh:
			return
		default:
			allocs, meta, err := w.client.Allocations().List(q)
			if err != nil {
				log.Fatalf("watch: unable to fetch allocations: %s", err)
				break
			}
			q = &api.QueryOptions{WaitIndex: meta.LastIndex}
			w.websocket.WriteJSON(event{
				Kind:    allocsEvent,
				Payload: allocs,
			})
		}
	}
}

func (w *watch) Evals() {
	q := &api.QueryOptions{WaitIndex: 0}

	for {
		select {
		case <-w.shutdownCh:
			return
		default:
			evals, meta, err := w.client.Evaluations().List(q)
			if err != nil {
				log.Fatalf("watch: unable to fetch evaluations: %s", err)
				break
			}
			q = &api.QueryOptions{WaitIndex: meta.LastIndex}
			w.websocket.WriteJSON(event{
				Kind:    evalsEvent,
				Payload: evals,
			})
		}
	}
}

func (w *watch) Jobs() {
	q := &api.QueryOptions{WaitIndex: 0}

	for {
		select {
		case <-w.shutdownCh:
			return
		default:
			jobs, meta, err := w.client.Jobs().List(q)
			if err != nil {
				log.Fatalf("watch: unable to fetch jobs: %s", err)
				break
			}
			q = &api.QueryOptions{WaitIndex: meta.LastIndex}
			w.websocket.WriteJSON(event{
				Kind:    jobsEvent,
				Payload: jobs,
			})
		}
	}
}

func (w *watch) Job(jobID string) {
	q := &api.QueryOptions{WaitIndex: 0}

	for {
		select {
		case <-w.shutdownCh:
			return
		default:
			job, meta, err := w.client.Jobs().Info(jobID, q)
			if err != nil {
				log.Fatalf("watch: unable to watch job %q: %s", jobID, err)
				break
			}
			q = &api.QueryOptions{WaitIndex: meta.LastIndex}
			w.websocket.WriteJSON(event{
				Kind:    jobEvent,
				Payload: job,
			})
		}
	}
}

func (w *watch) Nodes() {
	q := &api.QueryOptions{WaitIndex: 0}

	for {
		select {
		case <-w.shutdownCh:
			return
		default:
			nodes, meta, err := w.client.Nodes().List(q)
			if err != nil {
				log.Fatalf("watch: unable to fetch nodes: %s", err)
				break
			}
			q = &api.QueryOptions{WaitIndex: meta.LastIndex}
			w.websocket.WriteJSON(event{
				Kind:    nodesEvent,
				Payload: nodes,
			})
		}
	}
}
