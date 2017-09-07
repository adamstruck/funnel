package openstack

import (
	"context"
	"github.com/ohsu-comp-bio/funnel/compute/basic"
	"github.com/ohsu-comp-bio/funnel/config"
	"github.com/ohsu-comp-bio/funnel/logger"
	pbs "github.com/ohsu-comp-bio/funnel/proto/scheduler"
	"github.com/ohsu-comp-bio/funnel/proto/tes"
)

// Name of the scheduler backend
const Name = "openstack"

var log = logger.Sub(Name)

// NewBackend returns a new SchedulerBackend instance.
func NewBackend(conf config.Config) (basic.SchedulerBackend, error) {

	// Create a client for talking to the funnel scheduler
	client, err := basic.NewClient(conf.Backends.Basic)
	if err != nil {
		log.Error("Can't connect scheduler client", err)
		return nil, err
	}

	return &Backend{conf, client}, nil
}

// Backend represents the OpenStack backend.
type Backend struct {
	conf   config.Config
	client basic.Client
}

// GetOffer returns an offer based on available OpenStack VM node instances.
func (s *Backend) GetOffer(j *tes.Task) *basic.Offer {
	log.Debug("Running OpenStack backend")

	offers := []*basic.Offer{}
	predicates := append(basic.DefaultPredicates, basic.NodeHasTag("openstack"))

	for _, w := range s.getNodes() {
		// Filter out nodes that don't match the task request.
		// Checks CPU, RAM, disk space, ports, etc.
		if !basic.Match(w, j, predicates) {
			continue
		}

		sc := basic.DefaultScores(w, j)
		/*
				    TODO?
				    if w.State == pbs.NodeState_Alive {
						  sc["startup time"] = 1.0
				    }
			sc = sc.Weighted(s.conf.Backends.OpenStack.Weights)
		*/

		offer := basic.NewOffer(w, j, sc)
		offers = append(offers, offer)
	}

	// No matching nodes were found.
	if len(offers) == 0 {
		return nil
	}

	basic.SortByAverageScore(offers)
	return offers[0]
}

func (s *Backend) getNodes() []*pbs.Node {

	// Get the nodes from the funnel server
	nodes := []*pbs.Node{}
	req := &pbs.ListNodesRequest{}
	resp, err := s.client.ListNodes(context.Background(), req)

	// If there's an error, return an empty list
	if err != nil {
		log.Error("Failed ListNodes request. Recovering.", err)
		return nodes
	}

	nodes = resp.Nodes

	// TODO include unprovisioned node templates from config

	return nodes
}
