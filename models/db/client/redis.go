// /*
// MIT License
//
// # Copyright (c) 2024 Praromvik
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.
// */

package client

import (
	"bytes"
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"os"
	"os/exec"
)

func ConnectToRedis() (*redis.Client, error) {
	err := runRedisProcess()
	if err != nil {
		return nil, err
	}
	rClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6333",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	return rClient, nil
}

func runRedisProcess() error {
	ps := os.Getenv("REDIS_PROCESS_NAME")
	_ = exec.Command("docker", "stop", ps).Run()
	_ = exec.Command("docker", "rm", ps).Run()

	_ = listDockerProcesses()
	volume := fmt.Sprintf("%s:/data", os.Getenv("REDIS_VOLUME_PATH"))
	port := fmt.Sprintf("%s:6379", os.Getenv("REDIS_PORT"))
	logLevel := fmt.Sprintf("--loglevel %s", os.Getenv("REDIS_LOG_LEVEL"))

	cmd := exec.Command("docker", "run", "--name", ps, "-d", "-p", port, "-v", volume, "redis", "redis-server", logLevel)
	fmt.Printf("Run Redis process = %v \n", cmd.Args)
	return cmd.Run()
}

func listDockerProcesses() error {
	cmd := exec.Command("docker", "ps")
	var b []byte
	buf := bytes.NewBuffer(b)
	cmd.Stdout = buf
	err := cmd.Run()
	if err != nil {
		return err
	}
	// fmt.Printf("out : %s", string(buf.Bytes()))
	return nil
}

func TestRedisConnection(ctx context.Context, rClient *redis.Client) error {
	err := rClient.Set(ctx, "foo", "bar", 0).Err()
	if err != nil {
		return err
	}

	_, err = rClient.Get(ctx, "foo").Result()
	return err
}
