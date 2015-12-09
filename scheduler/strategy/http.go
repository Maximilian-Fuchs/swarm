package strategy

import (
	"github.com/docker/swarm/cluster"
	"github.com/docker/swarm/scheduler/node"
	"encoding/json"
	"net/http"
	"bytes"
	"io/ioutil"
	"fmt"
)

// SpreadPlacementStrategy places a container on the node with the fewest running containers.
type HTTPPlacementStrategy struct {
}

// Initialize a SpreadPlacementStrategy.
func (p *HTTPPlacementStrategy) Initialize() error {
	return nil
}

// Name returns the name of the strategy.
func (p *HTTPPlacementStrategy) Name() string {
	return "http"
}

func requestStrategy(config *cluster.ContainerConfig, nodes []*node.Node) ([]* node.Node, error){
	type StrategyRequest struct {
		Config *cluster.ContainerConfig
		Nodes []*node.Node
	}

	// Create JSON struct for request
	reqJson := StrategyRequest{
		Config: config,
		Nodes: nodes,
	}

	// Create JSON String from struct
	reqJsonStr, _ := json.Marshal(reqJson)

	// Build the HTTP Request
	url := "http://localhost:4000/strategy"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqJsonStr))
	req.Header.Set("Content-Type", "application/json")

	// Execute the HTTP Request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read the response
	body, _ := ioutil.ReadAll(resp.Body)

	// Construct the result
	json.Unmarshal(body, &nodes)

	return nodes, nil
}

// RankAndSort sorts nodes based on the spread strategy applied to the container config.
func (p *HTTPPlacementStrategy) RankAndSort(config *cluster.ContainerConfig, nodes []*node.Node) ([]*node.Node, error) {
	ans, err := requestStrategy(config, nodes)

	for idx, node := range(ans){
		fmt.Printf("%d - %s \n", idx, node.Name)
	}

	return ans, err
}
