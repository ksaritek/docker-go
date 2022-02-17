package main

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/strslice"
	"github.com/docker/docker/client"
	"io/ioutil"
	"log"
	"time"
)

func main() {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}

	im, err := cli.ImagePull(context.Background(), "alpine", types.ImagePullOptions{})
	if err != nil {
		panic(err)
	}

	im.Close()

	//imgSum, err := cli.ImageList(context.Background(), types.ImageListOptions{})
	//if err != nil {
	//	panic(err)
	//}
	//
	//for _, v := range imgSum {
	//	log.Printf("%v\n", v)
	//}

	cli.ContainerRemove(context.Background(), "testing", types.ContainerRemoveOptions{Force: true})

	cont, err := cli.ContainerCreate(context.Background(),
		&container.Config{
			Image: "alpine",
			Env:   []string{"env-1=val-1", "env-2=val-2"},
			Cmd:   strslice.StrSlice(append([]string{"/bin/sh"}, "-c", "env")),
		},
		nil,
		nil,
		nil,
		"testing",
	)
	if err != nil {
		panic(err)
	}

	if err = cli.ContainerStart(context.Background(), cont.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}
	fmt.Printf("Container %s is started\n", cont.ID)

	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}

	if len(containers) > 0 {
		for _, container := range containers {
			fmt.Printf("Container ID: %s\n", container.ID)
		}
	} else {
		fmt.Println("There are no containers running")
	}

	out, err := cli.ContainerLogs(context.Background(), cont.ID, types.ContainerLogsOptions{ShowStdout: true, ShowStderr: true})
	if err != nil {
		panic(err)
	}

	b, err := ioutil.ReadAll(out)
	if err != nil {
		panic(err)
	}
	log.Printf("Rolling log Contener \n%s", string(b))
	time.Sleep(time.Hour)
}
