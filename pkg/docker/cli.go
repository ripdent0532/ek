package docker

import (
	"encoding/json"
	"fmt"
	"github.com/docker/distribution/context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/ek/pkg/streams"
	"github.com/moby/term"
	"io"
	"strings"
)

type Streams interface {
	In() *streams.In
	Out() *streams.Out
	Err() io.Writer
}

type DockerCli struct {
	cli 			   *client.Client
	in                 *streams.In
	out                *streams.Out
	err                io.Writer
	auth			   map[string]string
}

func (dockerCli *DockerCli) Out() *streams.Out {
	return dockerCli.out
}

func (dockerCli *DockerCli) Err() io.Writer {
	return dockerCli.err
}

func (dockerCli *DockerCli) SetIn(in *streams.In) {
	dockerCli.in = in
}

func (dockerCli *DockerCli) In() *streams.In {
	return dockerCli.in
}

func NewClient(configPath string) *DockerCli {
	// 加载配置
	authMapping := GetAuthMapping(configPath)

	// 初始化客户端
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	_, stdout, _ := term.StdStreams()
	buildBuff := streams.NewOut(stdout)
	return &DockerCli{
		cli: cli,
		out: buildBuff,
		auth: authMapping,
	}
}

func (dockerCli *DockerCli) PullImage(image string) {
	ctx := context.Background()
	cli := dockerCli.cli


	reader, err := cli.ImagePull(ctx, image, types.ImagePullOptions{RegistryAuth: getAuth(image, dockerCli.auth)})
	if err != nil {
		panic(err.Error())
	}

	aux := func(msg jsonmessage.JSONMessage) {
		var result types.BuildResult
		if err := json.Unmarshal(*msg.Aux, &result); err != nil {

		} else {
			fmt.Println("我是ImageID： " + result.ID)
		}
	}

	err = jsonmessage.DisplayJSONMessagesStream(reader, dockerCli.Out(), dockerCli.Out().FD(), dockerCli.Out().IsTerminal(), aux)
	if err != nil {
		if jerr, ok := err.(*jsonmessage.JSONError); ok {
			// If no error code is set, default to 1
			if jerr.Code == 0 {
				jerr.Code = 1
			}
			print(jerr.Message)
		}
	}
	_, _ = fmt.Fprint(dockerCli.Out())
	reader.Close()
}

func (dockerCli *DockerCli) ImageTag(image string, tag string) {

	err := dockerCli.cli.ImageTag(context.Background(), image, tag)
	if err != nil {
		fmt.Errorf(err.Error())
	}
}

func (dockerCli *DockerCli) PushImage(image string) {

	resp, err := dockerCli.cli.ImagePush(context.Background(), image,  types.ImagePushOptions{RegistryAuth: getAuth(image, dockerCli.auth)})
	dockerCli.display(resp, err)

}

func (dockerCli *DockerCli) display(reader io.ReadCloser, err error) {
	if err != nil {
		panic(err.Error())
	}

	aux := func(msg jsonmessage.JSONMessage) {
		var result types.BuildResult
		if err := json.Unmarshal(*msg.Aux, &result); err != nil {

		} else {
			fmt.Println("我是ImageID： " + result.ID)
		}
	}

	err = jsonmessage.DisplayJSONMessagesStream(reader, dockerCli.Out(), dockerCli.Out().FD(), dockerCli.Out().IsTerminal(), aux)
	if err != nil {
		if jerr, ok := err.(*jsonmessage.JSONError); ok {
			// If no error code is set, default to 1
			if jerr.Code == 0 {
				jerr.Code = 1
			}
			print(jerr.Message)
		}
	}
	_, _ = fmt.Fprint(dockerCli.Out())
	reader.Close()
}

func getAuth(image string, authMapping map[string]string) string {
	host := strings.Split(image, "/")[0]
	return authMapping[host]
}