package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/couchbase/indexing/secondary/planner"
)

type RebalanceTestCase struct {
	Comment       string  `json:"comment"`
	Resize        bool    `json:"resize"`
	Plan          string  `json:"plan"`
	AddNode       int     `json:"addNode"`
	DeleteNode    int     `json:"deleteNode"`
	EjectOnly     bool    `json:"ejectOnly"`
	ReplicaRepair bool    `json:"replicaRepair"`
	Threshold     float64 `json:"threshold"`
	Timeout       int     `json:"timeout"`
	Detail        bool    `json:"detail"`
}

func main() {

	log.Printf("-------------------------------------------")
	if len(os.Args) != 2 {
		log.Printf("Usage: ./plannerSim <input_conf_file.json>")
		return
	}
	inpConf := os.Args[1]

	file, err := ioutil.ReadFile(inpConf)
	if err != nil {
		panic(err)
	}

	testcase := RebalanceTestCase{}

	if json.Unmarshal([]byte(file), &testcase) != nil {
		panic(err)
	}
	log.Printf("Test config: %+v", testcase)

	s := planner.NewSimulator()

	plan, err := planner.ReadPlan(testcase.Plan)
	if err != nil {
		panic(err)
	}

	config := planner.DefaultRunConfig()
	config.Resize = testcase.Resize
	config.AddNode = testcase.AddNode
	config.DeleteNode = testcase.DeleteNode
	config.EjectOnly = testcase.EjectOnly
	config.DisableRepair = testcase.ReplicaRepair
	config.Threshold = testcase.Threshold
	config.Timeout = testcase.Timeout
	now := time.Now()
	config.Runtime = &now

	p, _, err := s.RunSingleTest(config, planner.CommandRebalance, nil, plan, nil)
	if err != nil {
		panic(err)
	}

	p.PrintCost()

	if err := planner.ValidateSolution(p.Result); err != nil {
		panic(err)
	}

	p.PrintLayout()
}
