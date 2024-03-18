package container

import (
	"strings"

	co "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/go-connections/nat"
	"github.com/go-zoox/cli"
	"github.com/go-zoox/core-utils/fmt"
	"github.com/go-zoox/docker"
	"github.com/go-zoox/docker/container"
)

func Run() *cli.Command {
	return &cli.Command{
		Name:  "run",
		Usage: "Create and run a new container from an image",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "interactive",
				Aliases: []string{"i"},
				Usage:   "Keep STDIN open even if not attached",
			},
			&cli.BoolFlag{
				Name:    "tty",
				Usage:   "Allocate a pseudo-TTY",
				Aliases: []string{"t"},
				Value:   true,
			},
			&cli.StringFlag{
				Name:    "name",
				Aliases: []string{"n"},
				Usage:   "Assign a name to the container",
			},
			&cli.StringSliceFlag{
				Name:    "env",
				Aliases: []string{"e"},
				Usage:   "Set environment variables",
			},
			&cli.StringSliceFlag{
				Name:    "volume",
				Aliases: []string{"v"},
				Usage:   "Bind mount a volume",
			},
			&cli.StringSliceFlag{
				Name:    "publish",
				Aliases: []string{"p"},
				Usage:   "Publish a container's port(s) to the host",
			},
			&cli.StringFlag{
				Name:    "workdir",
				Aliases: []string{"w"},
				Usage:   "Working directory inside the container",
			},
			&cli.StringFlag{
				Name:    "user",
				Aliases: []string{"u"},
				Usage:   "Username or UID (format: <name|uid>[:<group|gid>])",
			},
			&cli.StringFlag{
				Name:  "hostname",
				Usage: "Container host name",
			},
			&cli.StringSliceFlag{
				Name:  "entrypoint",
				Usage: "Overwrite the default ENTRYPOINT of the image",
			},
			&cli.StringFlag{
				Name:  "health-cmd",
				Usage: "Command to run to check health",
			},
			&cli.StringFlag{
				Name:  "health-interval",
				Usage: "Time between running the check (ms|s|m|h) (default 0s)",
			},
			&cli.StringFlag{
				Name:  "health-timeout",
				Usage: "Maximum time to allow one check to run (ms|s|m|h) (default 0s)",
			},
			&cli.IntFlag{
				Name:  "health-retries",
				Usage: "Consecutive failures needed to report unhealthy",
			},
			&cli.StringFlag{
				Name:  "health-start-period",
				Usage: "Start period for the container to initialize before starting health-retries countdown (ms|s|m|h) (default 0s)",
			},
			&cli.StringSliceFlag{
				Name:    "label",
				Aliases: []string{"l"},
				Usage:   "Set meta data on a container",
			},
			&cli.StringFlag{
				Name:  "mac-address",
				Usage: "Container MAC address (e.g. 92:d0:c6:0a:29:33)",
			},
			&cli.StringFlag{
				Name:  "restart",
				Usage: "Restart policy to apply when a container exits (default \"no\")",
			},
			&cli.IntFlag{
				Name:  "restart-max-attempts",
				Usage: "Maximum restart count for the container (default 0)",
			},
			&cli.BoolFlag{
				Name:  "privileged",
				Usage: "Give extended privileges to this container",
			},
			&cli.StringSliceFlag{
				Name:  "link",
				Usage: "Add link to another container",
			},
			&cli.BoolFlag{
				Name:  "rm",
				Usage: "Automatically remove the container when it exits",
			},
			&cli.Int64Flag{
				Name:  "cpu-period",
				Usage: "Limit CPU CFS (Completely Fair Scheduler) period",
			},
			&cli.Int64Flag{
				Name:  "cpu-quota",
				Usage: "Limit CPU CFS (Completely Fair Scheduler) quota",
			},
			&cli.Int64Flag{
				Name:  "cpu-shares",
				Usage: "CPU shares (relative weight)",
			},
			&cli.StringFlag{
				Name:  "cpuset-cpus",
				Usage: "CPUs in which to allow execution (0-3, 0,1)",
			},
			&cli.StringFlag{
				Name:  "cpuset-mems",
				Usage: "MEMs in which to allow execution (0-3, 0,1)",
			},
			&cli.Int64Flag{
				Name:  "cpu-rt-period",
				Usage: "Limit CPU real-time period in microseconds",
			},
			&cli.Int64Flag{
				Name:  "cpu-rt-runtime",
				Usage: "Limit CPU real-time runtime in microseconds",
			},
			&cli.Int64Flag{
				Name:  "cpus",
				Usage: "Number of CPUs",
			},
			&cli.Int64Flag{
				Name:  "cpu-percent",
				Usage: "CPU percent",
			},
			&cli.Int64Flag{
				Name:  "memory",
				Usage: "Memory limit",
			},
			&cli.Int64Flag{
				Name:  "memory-reservation",
				Usage: "Memory soft limit",
			},
			&cli.Int64Flag{
				Name:  "memory-swap",
				Usage: "Swap limit equal to memory plus swap: '-1' to enable unlimited swap",
			},
			&cli.Int64Flag{
				Name:  "memory-swappiness",
				Usage: "Tune container memory swappiness (0 to 100)",
			},
			&cli.StringSliceFlag{
				Name:  "dns",
				Usage: "Set custom DNS servers",
			},
			&cli.StringSliceFlag{
				Name:  "dns-search",
				Usage: "Set custom DNS search domains",
			},
			&cli.StringSliceFlag{
				Name:  "dns-opt",
				Usage: "Set DNS options",
			},
			&cli.StringSliceFlag{
				Name:  "security-opt",
				Usage: "Security Options",
			},
			&cli.StringFlag{
				Name:  "platform",
				Usage: "Set platform if server is multi-platform capable",
			},
			&cli.BoolFlag{
				Name:    "detach",
				Usage:   "Run container in background and print container ID",
				Aliases: []string{"d"},
			},
		},
		Action: func(ctx *cli.Context) error {
			client, err := docker.New()
			if err != nil {
				return err
			}

			image := ctx.Args().First()
			if image == "" {
				return fmt.Errorf("image is required")
			}
			cmd := ctx.Args().Tail()

			return client.Container().Run(ctx.Context, func(opt *container.RunOptions) {
				opt.Name = ctx.String("name")

				//
				opt.Detached = ctx.Bool("detach")

				// container
				opt.Container.Image = image
				opt.Container.Env = ctx.StringSlice("env")
				opt.Container.Cmd = cmd
				opt.Container.Volumes = map[string]struct{}{}
				for _, volumeKV := range ctx.StringSlice("volume") {
					kv := strings.SplitN(volumeKV, ":", 2)
					volumeKey, volumeValue := kv[0], kv[1]
					opt.Host.Mounts = append(opt.Host.Mounts, mount.Mount{
						Type:   mount.TypeBind,
						Source: volumeKey,
						Target: volumeValue,
						// ReadOnly: false,
					})
				}
				opt.Container.ExposedPorts = nat.PortSet{}
				for _, portKV := range ctx.StringSlice("publish") {
					kv := strings.SplitN(portKV, ":", 2)
					portKey, portValue := kv[0], kv[1]
					opt.Container.ExposedPorts[nat.Port(portKey)] = struct{}{}
					opt.Host.PortBindings = nat.PortMap{}
					opt.Host.PortBindings[nat.Port(portKey)] = []nat.PortBinding{
						{
							HostIP:   "0.0.0.0",
							HostPort: portValue,
						},
					}
				}

				opt.Container.WorkingDir = ctx.String("workdir")
				opt.Container.User = ctx.String("user")
				opt.Container.Hostname = ctx.String("hostname")
				opt.Container.Entrypoint = ctx.StringSlice("entrypoint")
				opt.Container.Healthcheck = &co.HealthConfig{
					Test: []string{
						"sh", ctx.String("health-cmd"),
					},
					Interval:    ctx.Duration("health-interval"),
					Timeout:     ctx.Duration("health-timeout"),
					Retries:     ctx.Int("health-retries"),
					StartPeriod: ctx.Duration("health-start-period"),
				}
				opt.Container.Labels = map[string]string{}
				for _, labelKV := range ctx.StringSlice("label") {
					kv := strings.SplitN(labelKV, "=", 2)
					labelKey, labelValue := kv[0], kv[1]
					opt.Container.Labels[labelKey] = labelValue
				}
				opt.Container.MacAddress = ctx.String("mac-address")
				opt.Container.Tty = ctx.Bool("tty")

				// host
				opt.Host.RestartPolicy = co.RestartPolicy{
					Name:              co.RestartPolicyMode(ctx.String("restart")),
					MaximumRetryCount: ctx.Int("restart-max-attempts"),
				}
				opt.Host.Privileged = ctx.Bool("privileged")
				opt.Host.Links = ctx.StringSlice("link")
				opt.Host.AutoRemove = ctx.Bool("rm")
				// Resources
				opt.Host.Resources.CPUPeriod = ctx.Int64("cpu-period")
				opt.Host.Resources.CPUQuota = ctx.Int64("cpu-quota")
				opt.Host.Resources.CPUShares = ctx.Int64("cpu-shares")
				opt.Host.Resources.CpusetCpus = ctx.String("cpuset-cpus")
				opt.Host.Resources.CpusetMems = ctx.String("cpuset-mems")
				opt.Host.Resources.CPURealtimePeriod = ctx.Int64("cpu-rt-period")
				opt.Host.Resources.CPURealtimeRuntime = ctx.Int64("cpu-rt-runtime")
				// Applicable to Windows
				opt.Host.Resources.CPUCount = ctx.Int64("cpus")
				opt.Host.Resources.CPUPercent = ctx.Int64("cpu-percent")
				//
				opt.Host.Resources.Memory = ctx.Int64("memory")
				opt.Host.Resources.MemoryReservation = ctx.Int64("memory-reservation")
				opt.Host.Resources.MemorySwap = ctx.Int64("memory-swap")
				ms := ctx.Int64("memory-swappiness")
				opt.Host.Resources.MemorySwappiness = &ms
				//
				opt.Host.DNS = ctx.StringSlice("dns")
				opt.Host.DNSSearch = ctx.StringSlice("dns-search")
				opt.Host.DNSOptions = ctx.StringSlice("dns-opt")
				//
				opt.Host.SecurityOpt = ctx.StringSlice("security-opt")

				// Network
				opt.Network.EndpointsConfig = map[string]*network.EndpointSettings{}
				opt.Network.EndpointsConfig["bridge"] = &network.EndpointSettings{
					NetworkID: "bridge",
				}

				// platform
				// 	examples:
				//		linux/amd64
				//		linux/arm64
				platform := ctx.String("platform")
				for _, pl := range strings.Split(platform, ",") {
					if pl == "" {
						continue
					}

					kv := strings.Split(pl, "/")
					platformOS, platformArch := kv[0], kv[1]
					var platformVariant string
					if len(kv) >= 3 {
						platformVariant = kv[2]
					}
					opt.Platform.OS = platformOS
					opt.Platform.Architecture = platformArch
					opt.Platform.Variant = platformVariant
				}
			})
		},
	}
}
