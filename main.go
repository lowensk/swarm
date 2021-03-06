package main

import (
	"fmt"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"

	"github.com/docker/swarm/discovery"
	_ "github.com/docker/swarm/discovery/consul"
	_ "github.com/docker/swarm/discovery/etcd"
	_ "github.com/docker/swarm/discovery/file"
	"github.com/docker/swarm/discovery/token"
	_ "github.com/docker/swarm/discovery/zookeeper"
)

func main() {
	app := cli.NewApp()
	app.Name = "swarm"
	app.Usage = "docker clustering"
	app.Version = "0.0.1"
	app.Author = ""
	app.Email = ""

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:   "debug",
			Usage:  "debug mode",
			EnvVar: "DEBUG",
		},
	}

	// logs
	app.Before = func(c *cli.Context) error {
		log.SetOutput(os.Stderr)
		if c.Bool("debug") {
			log.SetLevel(log.DebugLevel)
		}
		return nil
	}

	app.Commands = []cli.Command{
		{
			Name:      "create",
			ShortName: "c",
			Usage:     "create a cluster",
			Action: func(c *cli.Context) {
				discovery := &token.TokenDiscoveryService{}
				discovery.Initialize("", 0)
				token, err := discovery.CreateCluster()
				if err != nil {
					log.Fatal(err)
				}
				fmt.Println(token)
			},
		},
		{
			Name:      "list",
			ShortName: "l",
			Usage:     "list nodes in a cluster",
			Flags:     []cli.Flag{flDiscovery},
			Action: func(c *cli.Context) {
				if c.String("discovery") == "" {
					log.Fatal("--discovery required to list a cluster")
				}

				d, err := discovery.New(c.String("discovery"), 0)
				if err != nil {
					log.Fatal(err)
				}

				nodes, err := d.Fetch()
				if err != nil {
					log.Fatal(err)
				}
				for _, node := range nodes {
					fmt.Println(node)
				}
			},
		},
		{
			Name:      "manage",
			ShortName: "m",
			Usage:     "manage a docker cluster",
			Flags: []cli.Flag{
				flStrategy, flFilter,
				flDiscovery, flHosts, flHeartBeat,
				flTls, flTlsCaCert, flTlsCert, flTlsKey, flTlsVerify,
				flEnableCors},
			Action: manage,
		},
		{
			Name:      "join",
			ShortName: "j",
			Usage:     "join a docker cluster",
			Flags:     []cli.Flag{flDiscovery, flAddr, flHeartBeat},
			Action:    join,
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
