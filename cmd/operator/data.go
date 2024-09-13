package main

import "github.com/AlexGustafsson/cupdate/internal/api"

var mockAPI = &api.InMemoryAPI{
	Tags: []api.Tag{
		{
			Name:        "k8s",
			Description: "Kubernetes",
			Color:       "#DBEAFE",
		},
		{
			Name:        "pod",
			Description: "A kubernetes pod",
			Color:       "#FFEDD5",
		},
		{
			Name:        "job",
			Description: "A kubernetes job",
			Color:       "#DBEAFE",
		},
		{
			Name:        "docker",
			Description: "A docker container",
			Color:       "#FEE2E2",
		},
		{
			Name:        "up-to-date",
			Description: "Up-to-date images",
			Color:       "#DCFCE7",
		},
		{
			Name:        "outdated",
			Description: "Outdated images",
			Color:       "#FEE2E2",
		},
	},
	Images: []api.Image{
		{
			Name:           "home-assistant",
			CurrentVersion: "2024.4.4",
			LatestVersion:  "2024.8.3",
			Tags:           []string{"k8s", "pod", "outdated"},
			Links: []api.ImageLink{
				{
					Type: "github",
					URL:  "https://github.com/home-assistant/core",
				},
				{
					Type: "gitlab",
					URL:  "https://gitlab.com/arm-research/smarter/smarter-device-manager",
				},
				{
					Type: "docker",
					URL:  "https://hub.docker.com/r/homeassistant/home-assistant",
				},
				{
					Type: "quay",
					URL:  "https://quay.io/repository/jetstack/cert-manager-webhook?tab=info",
				},
				{
					Type: "git",
					URL:  "https://github.com/home-assistant/core",
				},
			},
			Image: "https://www.gravatar.com/avatar/461df105cc6cfcf386ebd5af85b925dc?s=120&r=g&d=404",
		},
		{
			Name:           "jacobalberty/unifi",
			CurrentVersion: "v7",
			LatestVersion:  "v7",
			Links: []api.ImageLink{
				{
					Type: "github",
					URL:  "https://github.com/home-assistant/core",
				},
				{
					Type: "gitlab",
					URL:  "https://gitlab.com/arm-research/smarter/smarter-device-manager",
				},
				{
					Type: "docker",
					URL:  "https://hub.docker.com/r/homeassistant/home-assistant",
				},
				{
					Type: "quay",
					URL:  "https://quay.io/repository/jetstack/cert-manager-webhook?tab=info",
				},
				{
					Type: "git",
					URL:  "https://github.com/home-assistant/core",
				},
			},
			Tags: []string{"k8s", "pod", "up-to-date"},
		},
		{
			Name:           "yooooomi/your_spotify_server",
			CurrentVersion: "1.11.0",
			LatestVersion:  "1.11.0",
			Links: []api.ImageLink{
				{
					Type: "github",
					URL:  "https://github.com/home-assistant/core",
				},
				{
					Type: "gitlab",
					URL:  "https://gitlab.com/arm-research/smarter/smarter-device-manager",
				},
				{
					Type: "docker",
					URL:  "https://hub.docker.com/r/homeassistant/home-assistant",
				},
				{
					Type: "quay",
					URL:  "https://quay.io/repository/jetstack/cert-manager-webhook?tab=info",
				},
				{
					Type: "git",
					URL:  "https://github.com/home-assistant/core",
				},
			},
			Tags: []string{"k8s", "pod", "up-to-date"},
		},
		{
			Name:           "hashicorp/vault",
			CurrentVersion: "2024.4.4",
			LatestVersion:  "2024.4.8",
			Links: []api.ImageLink{
				{
					Type: "github",
					URL:  "https://github.com/home-assistant/core",
				},
				{
					Type: "gitlab",
					URL:  "https://gitlab.com/arm-research/smarter/smarter-device-manager",
				},
				{
					Type: "docker",
					URL:  "https://hub.docker.com/r/homeassistant/home-assistant",
				},
				{
					Type: "quay",
					URL:  "https://quay.io/repository/jetstack/cert-manager-webhook?tab=info",
				},
				{
					Type: "git",
					URL:  "https://github.com/home-assistant/core",
				},
			},
			Tags: []string{"k8s", "pod", "outdated"},
		},
	},
	Descriptions: map[string]*api.ImageDescription{
		"home-assistant:2024.4.4": {
			HTML: `<h1 id="home-assistant">Home Assistant</h1>
      <p>
        Open source home automation that puts local control and privacy first.
        Powered by a worldwide community of tinkerers and DIY enthusiasts.
        Perfect to run on a Raspberry Pi or a local server.
      </p>
      <p>
        Check out <a href="https://home-assistant.io">home-assistant.io</a> for
        <a href="https://home-assistant.io/demo/">a demo</a>,
        <a href="https://home-assistant.io/getting-started/">
          installation instructions
        </a>
        ,
        <a href="https://home-assistant.io/getting-started/automation-2/">
          tutorials
        </a>
        and <a href="https://home-assistant.io/docs/">documentation</a>.
      </p>
      <p>
        <img
          src="https://raw.github.com/home-assistant/home-assistant/master/docs/screenshots.png"
          alt="screenshot states"
        />
      </p>
      <h2 id="featured-integrations">Featured integrations</h2>
      <p>
        <img
          src="https://raw.github.com/home-assistant/home-assistant/dev/docs/screenshot-components.png"
          alt="screenshot components"
        />
      </p>
      <p>
        If you run into issues while using Home Assistant, check the
        <a href="https://home-assistant.io/help/">
          Home Assistant help section
        </a>
        of our website for further help and information.
      </p>`,
		},
	},
	ReleaseNotes: map[string]*api.ImageReleaseNotes{
		"home-assistant:2024.4.4": {
			Title: "2024.8.8",
			HTML: `<ul>
        <li>
          Fix BTHome validate triggers for device with multiple buttons (
          <a href="https://github.com/thecode">@thecode</a> -
          <a
            href="https://github.com/home-assistant/core/pull/125183"
            data-hovercard-type="pull_request"
            data-hovercard-url="/home-assistant/core/pull/125183/hovercard"
          >
            #125183
          </a>
          ) (
          <a
            href="https://www.home-assistant.io/integrations/bthome/"
            rel="nofollow"
          >
            bthome docs
          </a>
          )
        </li>
        <li>
          Improve play media support in LinkPlay (
          <a href="https://github.com/silamon">@silamon</a> -
          <a href="https://github.com/home-assistant/core/pull/125205">
            #125205
          </a>
          ) (
          <a
            href="https://www.home-assistant.io/integrations/linkplay/"
            rel="nofollow"
          >
            linkplay docs
          </a>
          )
        </li>
        <li>
          Revert
          <a href="https://github.com/home-assistant/core/pull/122676">
            #122676
          </a>
          Yamaha discovery (<a href="https://github.com/joostlek">@joostlek</a>
          -
          <a href="https://github.com/home-assistant/core/pull/125216">
            #125216
          </a>
          ) (
          <a
            href="https://www.home-assistant.io/integrations/yamaha/"
            rel="nofollow"
          >
            yamaha docs
          </a>
          )
        </li>
        <li>
          Fix blocking call in yale_smart_alarm (
          <a href="https://github.com/gjohansson-ST">@gjohansson-ST</a> -
          <a href="https://github.com/home-assistant/core/pull/125255">
            #125255
          </a>
          ) (
          <a
            href="https://www.home-assistant.io/integrations/yale_smart_alarm/"
            rel="nofollow"
          >
            yale_smart_alarm docs
          </a>
          )
        </li>
        <li>
          Don't show input panel if default code provided in envisalink (
          <a href="https://github.com/gjohansson-ST">@gjohansson-ST</a> -
          <a href="https://github.com/home-assistant/core/pull/125256">
            #125256
          </a>
          ) (
          <a
            href="https://www.home-assistant.io/integrations/envisalink/"
            rel="nofollow"
          >
            envisalink docs
          </a>
          )
        </li>
        <li>
          Increase AquaCell timeout and handle timeout exception properly (
          <a href="https://github.com/Jordi1990">@Jordi1990</a> -
          <a href="https://github.com/home-assistant/core/pull/125263">
            #125263
          </a>
          ) (
          <a
            href="https://www.home-assistant.io/integrations/aquacell/"
            rel="nofollow"
          >
            aquacell docs
          </a>
          )
        </li>
        <li>
          Bump aiorussound to 3.0.4 (
          <a href="https://github.com/noahhusby">@noahhusby</a> -
          <a href="https://github.com/home-assistant/core/pull/125285">
            #125285
          </a>
          ) (
          <a
            href="https://www.home-assistant.io/integrations/russound_rio/"
            rel="nofollow"
          >
            russound_rio docs
          </a>
          )
        </li>
        <li>
          Add follower to the PlayingMode enum (
          <a href="https://github.com/silamon">@silamon</a> -
          <a href="https://github.com/home-assistant/core/pull/125294">
            #125294
          </a>
          ) (
          <a
            href="https://www.home-assistant.io/integrations/linkplay/"
            rel="nofollow"
          >
            linkplay docs
          </a>
          )
        </li>
        <li>
          Don't allow templating min, max, step in config entry template number
          (<a href="https://github.com/emontnemery">@emontnemery</a> -
          <a href="https://github.com/home-assistant/core/pull/125342">
            #125342
          </a>
          ) (
          <a
            href="https://www.home-assistant.io/integrations/template/"
            rel="nofollow"
          >
            template docs
          </a>
          )
        </li>
        <li>
          Fix for Hue sending effect None at turn_on command while no effect is
          active (<a href="https://github.com/marcelveldt">@marcelveldt</a> -
          <a href="https://github.com/home-assistant/core/pull/125377">
            #125377
          </a>
          ) (
          <a
            href="https://www.home-assistant.io/integrations/hue/"
            rel="nofollow"
          >
            hue docs
          </a>
          )
        </li>
        <li>
          Lyric: fixed missed snake case conversions (
          <a href="https://github.com/dalinicus">@dalinicus</a> -
          <a href="https://github.com/home-assistant/core/pull/125382">
            #125382
          </a>
          ) (
          <a
            href="https://www.home-assistant.io/integrations/lyric/"
            rel="nofollow"
          >
            lyric docs
          </a>
          )
        </li>
        <li>
          Bump pysmlight to 0.0.14 (
          <a href="https://github.com/tl-sl">@tl-sl</a> -
          <a href="https://github.com/home-assistant/core/pull/125387">
            #125387
          </a>
          ) (
          <a
            href="https://www.home-assistant.io/integrations/smlight/"
            rel="nofollow"
          >
            smlight docs
          </a>
          )
        </li>
        <li>
          Bump pypck to 0.7.22 (
          <a href="https://github.com/alengwenus">@alengwenus</a> -
          <a href="https://github.com/home-assistant/core/pull/125389">
            #125389
          </a>
          ) (
          <a
            href="https://www.home-assistant.io/integrations/lcn/"
            rel="nofollow"
          >
            lcn docs
          </a>
          )
        </li>
        <li>
          Increase coordinator update_interval for fyta (
          <a href="https://github.com/dontinelli">@dontinelli</a> -
          <a
            href="https://github.com/home-assistant/core/pull/125393"
            data-hovercard-type="pull_request"
            data-hovercard-url="/home-assistant/core/pull/125393/hovercard"
          >
            #125393
          </a>
          ) (
          <a
            href="https://www.home-assistant.io/integrations/fyta/"
            rel="nofollow"
          >
            fyta docs
          </a>
          )
        </li>
        <li>
          Fix controlling AC temperature in airtouch5 (
          <a href="https://github.com/danzel">@danzel</a> -
          <a
            href="https://github.com/home-assistant/core/pull/125394"
            data-hovercard-type="pull_request"
            data-hovercard-url="/home-assistant/core/pull/125394/hovercard"
          >
            #125394
          </a>
          ) (
          <a
            href="https://www.home-assistant.io/integrations/airtouch5/"
            rel="nofollow"
          >
            airtouch5 docs
          </a>
          )
        </li>
        <li>
          Bump sfrbox-api to 0.0.10 (
          <a href="https://github.com/AlexT59">@AlexT59</a> -
          <a
            href="https://github.com/home-assistant/core/pull/125405"
            data-hovercard-type="pull_request"
            data-hovercard-url="/home-assistant/core/pull/125405/hovercard"
          >
            #125405
          </a>
          )
        </li>
        <li>
          Improve handling of old firmware versions (
          <a href="https://github.com/tl-sl">@tl-sl</a> -
          <a
            href="https://github.com/home-assistant/core/pull/125406"
            data-hovercard-type="pull_request"
            data-hovercard-url="/home-assistant/core/pull/125406/hovercard"
          >
            #125406
          </a>
          ) (
          <a
            href="https://www.home-assistant.io/integrations/smlight/"
            rel="nofollow"
          >
            smlight docs
          </a>
          )
        </li>
        <li>
          Set min_power similar to max_power to support all inverters from
          apsystems (
          <a href="https://github.com/mawoka-myblock">@mawoka-myblock</a> -
          <a
            href="https://github.com/home-assistant/core/pull/124247"
            data-hovercard-type="pull_request"
            data-hovercard-url="/home-assistant/core/pull/124247/hovercard"
          >
            #124247
          </a>
          ) (
          <a
            href="https://www.home-assistant.io/integrations/apsystems/"
            rel="nofollow"
          >
            apsystems docs
          </a>
          )
        </li>
        <li>
          Update frontend to 20240906.0 (
          <a href="https://github.com/piitaya">@piitaya</a> -
          <a
            href="https://github.com/home-assistant/core/pull/125409"
            data-hovercard-type="pull_request"
            data-hovercard-url="/home-assistant/core/pull/125409/hovercard"
          >
            #125409
          </a>
          ) (
          <a
            href="https://www.home-assistant.io/integrations/frontend/"
            rel="nofollow"
          >
            frontend docs
          </a>
          )
        </li>
        <li>
          Bump pyatv to 0.15.1 (
          <a href="https://github.com/postlund">@postlund</a> -
          <a
            href="https://github.com/home-assistant/core/pull/125412"
            data-hovercard-type="pull_request"
            data-hovercard-url="/home-assistant/core/pull/125412/hovercard"
          >
            #125412
          </a>
          ) (
          <a
            href="https://www.home-assistant.io/integrations/apple_tv/"
            rel="nofollow"
          >
            apple_tv docs
          </a>
          )
        </li>
      </ul>`,
		},
	},
	Graphs: map[string]*api.Graph{
		"home-assistant:2024.4.4": {
			Root: api.GraphNode{
				Domain: "oci",
				Type:   "image",
				Name:   "home-assistant",
				Parents: []api.GraphNode{
					{
						Domain: "kubernetes",
						Type:   "core/v1/container",
						Name:   "home-assistant",
						Parents: []api.GraphNode{
							{
								Domain: "kubernetes",
								Type:   "core/v1/pod",
								Name:   "home-assistant",
								Parents: []api.GraphNode{
									{
										Domain: "kubernetes",
										Type:   "apps/v1/deployment",
										Name:   "home-assistant",
										Parents: []api.GraphNode{
											{
												Domain:  "kubernetes",
												Type:    "core/v1/namespace",
												Name:    "home-assistant",
												Parents: []api.GraphNode{},
											},
										},
									},
								},
							},
						},
					},
					{
						Domain: "kubernetes",
						Type:   "core/v1/container",
						Name:   "home-assistant",
						Parents: []api.GraphNode{
							{
								Domain: "kubernetes",
								Type:   "core/v1/pod",
								Name:   "home-assistant",
								Parents: []api.GraphNode{
									{
										Domain:  "kubernetes",
										Type:    "core/v1/namespace",
										Name:    "home-assistant",
										Parents: []api.GraphNode{},
									},
								},
							},
						},
					},
				},
			},
		},
	},
}
