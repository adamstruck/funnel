package gce

// TODO
// - resource tracking via GCP APIs
// - provisioning limits, e.g. don't create more than 100 VMs, or
//   maybe use N VCPUs max, across all VMs
// - act on failed machines?
// - know how to shutdown machines

import (
	"context"
	"fmt"
	"github.com/ohsu-comp-bio/funnel/compute/basic"
	"github.com/ohsu-comp-bio/funnel/config"
	"github.com/ohsu-comp-bio/funnel/logger"
	pbs "github.com/ohsu-comp-bio/funnel/proto/scheduler"
	"github.com/ohsu-comp-bio/funnel/proto/tes"
)

// Name of the scheduler backend
const Name = "gce"

var log = logger.Sub(Name)

// NewBackend returns a new Google Cloud Engine SchedulerBackend instance.
func NewBackend(conf config.Config) (basic.SchedulerBackend, error) {
	// TODO need GCE basic config validation. If zone is missing, nothing works.

	// Create a client for talking to the funnel scheduler
	client, err := basic.NewClient(conf.Backends.Basic)
	if err != nil {
		log.Error("Can't connect scheduler client", err)
		return nil, err
	}

	// Create a client for talking to the GCE API
	gce, gerr := newClientFromConfig(conf)
	if gerr != nil {
		log.Error("Can't connect GCE client", gerr)
		return nil, gerr
	}

	return &Backend{
		conf:   conf,
		client: client,
		gce:    gce,
	}, nil
}

// Backend represents the GCE backend, which provides
// and interface for both scheduling and scaling.
type Backend struct {
	conf   config.Config
	client basic.Client
	gce    Client
}

// GetOffer returns an offer based on available Google Cloud VM node instances.
func (s *Backend) GetOffer(j *tes.Task) *basic.Offer {
	log.Debug("Running GCE backend")

	offers := []*basic.Offer{}
	predicates := append(basic.DefaultPredicates, basic.NodeHasTag("gce"))

	for _, n := range s.getNodes() {
		// Filter out nodes that don't match the task request.
		// Checks CPU, RAM, disk space, ports, etc.
		if !basic.Match(n, j, predicates) {
			continue
		}

		sc := basic.DefaultScores(n, j)
		/*
			    TODO?
			    if w.State == pbs.NodeState_Alive {
					  sc["startup time"] = 1.0
			    }
		*/
		weights := map[string]float32{}
		sc = sc.Weighted(weights)

		offer := basic.NewOffer(n, j, sc)
		offers = append(offers, offer)
	}

	// No matching nodes were found.
	if len(offers) == 0 {
		return nil
	}

	basic.SortByAverageScore(offers)
	return offers[0]
}

// getNodes returns a list of all GCE nodes and appends a set of
// uninitialized nodes, which the scheduler can use to create new node VMs.
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

	// Include unprovisioned (template) nodes.
	// This is how the scheduler can schedule tasks to nodes that
	// haven't been started yet.
	for _, t := range s.gce.Templates() {
		t.Id = basic.GenNodeID("funnel")
		nodes = append(nodes, &t)
	}

	return nodes
}

// ShouldStartNode tells the scaler loop which nodes
// belong to this backend, basically.
func (s *Backend) ShouldStartNode(n *pbs.Node) bool {
	// Only start works that are uninitialized and have a gce template.
	tpl, ok := n.Metadata["gce-template"]
	return ok && tpl != "" && n.State == pbs.NodeState_UNINITIALIZED
}

// StartNode calls out to GCE APIs to start a new node instance.
func (s *Backend) StartNode(n *pbs.Node) error {

	// Get the template ID from the node metadata
	template, ok := n.Metadata["gce-template"]
	if !ok || template == "" {
		return fmt.Errorf("Could not get GCE template ID from metadata")
	}

	return s.gce.StartNode(template, s.conf.Server.RPCAddress(), n.Id)
}
