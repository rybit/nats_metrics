package main

import (
	"log"

	"github.com/nats-io/nats"
	metrics "github.com/rybit/nats_metrics"
)

func main() {
	nc, err := nats.Connect("nats://localhost:4222")
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()
	err = metrics.Init(nc, "metrics")
	if err != nil {
		log.Fatal(err)
	}

	metrics.AddDimension("space", "global")
	metrics.AddDimension("app", "example")
	c := metrics.NewCounter("one-ups", &metrics.DimMap{
		"space": "metric",
		"magic": "value",
	})
	if err != nil {
		log.Fatal(err)
	}
	c.Count(&metrics.DimMap{
		"space":    "instance",
		"instance": "level",
	})
}
