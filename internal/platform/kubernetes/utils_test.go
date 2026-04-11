package kubernetes

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetImageReference(t *testing.T) {
	testCases := []struct {
		SpecImage     string
		StatusImage   string
		StatusImageID string
		Expected      string
		ExpectErr     bool
	}{
		// Happy paths, edge cases
		{
			SpecImage: "mongo:4",
			// No status present - use image from spec
			Expected: "mongo:4",
		},
		{
			SpecImage:   "mongo@sha256:9c20e607b82fc66a0b81a45c04d6ccd8fd056add3a3adacc0bb7a6b99460fdb0",
			StatusImage: "mongo",
			// Status image present, but not as detailed - use image from spec
			Expected: "mongo@sha256:9c20e607b82fc66a0b81a45c04d6ccd8fd056add3a3adacc0bb7a6b99460fdb0",
		},
		{
			SpecImage:     "mongo:4",
			StatusImageID: "mongo",
			// Status image id present, but not as detailed - use image from spec
			Expected: "mongo:4",
		},
		{
			SpecImage:     "mongo:4",
			StatusImageID: "mongo@sha256:9c20e607b82fc66a0b81a45c04d6ccd8fd056add3a3adacc0bb7a6b99460fdb0",
			// Status image id present and valid, but without tag - use tag from spec
			// and sha from status
			Expected: "mongo:4@sha256:9c20e607b82fc66a0b81a45c04d6ccd8fd056add3a3adacc0bb7a6b99460fdb0",
		},
		{
			SpecImage:     "mongo:4",
			StatusImageID: "mongo:4@sha256:9c20e607b82fc66a0b81a45c04d6ccd8fd056add3a3adacc0bb7a6b99460fdb0",
			// Status image id present and valid, with tag - use image from status
			Expected: "mongo:4@sha256:9c20e607b82fc66a0b81a45c04d6ccd8fd056add3a3adacc0bb7a6b99460fdb0",
		},
		{
			SpecImage:     "mongo",
			StatusImage:   "mongo:latest",
			StatusImageID: "mongo@sha256:9c20e607b82fc66a0b81a45c04d6ccd8fd056add3a3adacc0bb7a6b99460fdb0",
			// Status image id present and valid, with tag - use image from status
			Expected: "mongo:latest@sha256:9c20e607b82fc66a0b81a45c04d6ccd8fd056add3a3adacc0bb7a6b99460fdb0",
		},
		{
			SpecImage:     "traefik:3.6.7",
			StatusImage:   "docker.io/library/traefik:3.6.7",
			StatusImageID: "sha256:91528df1690f7da08360dcbbcb92b3ea483eeceb9f136d309f17506a5bd3f58d",
			// Status image id only contains the digest, but along with either of the
			// spec or status image, a full reference can be created
			Expected: "traefik:3.6.7@sha256:91528df1690f7da08360dcbbcb92b3ea483eeceb9f136d309f17506a5bd3f58d",
		},
		// Failure cases
		{
			SpecImage:     "mongo:4",
			StatusImageID: "aW52YWxpZA==",
			// Status image id present, but not a valid reference - use image from spec
			Expected: "mongo:4",
		},
		{
			SpecImage:     "mongo:4",
			StatusImageID: "aW52YWxpZA==",
			StatusImage:   "aW52YWxpZA==",
			// Status image and image id present, but not a valid reference - use image from spec
			Expected: "mongo:4",
		},
		{
			SpecImage: "",
			// Image not present in spec
			ExpectErr: true,
		},
		{
			SpecImage: "aW52YWxpZA==",
			// Image present in spec, but invalid
			ExpectErr: true,
		},
		// Real-world cases
		{
			SpecImage:     "intel/intel-gpu-plugin:0.35.0",
			StatusImage:   "docker.io/intel/intel-gpu-plugin:0.35.0",
			StatusImageID: "docker.io/intel/intel-gpu-plugin@sha256:34697f9c286857da986381595ac2a693524a83c831955247dae47dfda4d2f526",
			Expected:      "intel/intel-gpu-plugin:0.35.0@sha256:34697f9c286857da986381595ac2a693524a83c831955247dae47dfda4d2f526",
		},
		{
			SpecImage:     "prom/node-exporter:v1.11.0-distroless",
			StatusImage:   "docker.io/prom/node-exporter:v1.11.0-distroless",
			StatusImageID: "docker.io/prom/node-exporter@sha256:9ae87be1b066b133100b88d65cac5f83b9ed07bce56c2fc7ea58462f313c6bb1",
			Expected:      "prom/node-exporter:v1.11.0-distroless@sha256:9ae87be1b066b133100b88d65cac5f83b9ed07bce56c2fc7ea58462f313c6bb1",
		},
		{
			SpecImage:     "ghcr.io/alexgustafsson/srdl:0.4.4",
			StatusImage:   "ghcr.io/alexgustafsson/srdl:0.4.4",
			StatusImageID: "ghcr.io/alexgustafsson/srdl@sha256:919fea57db383a9b161302c08cd5e60c9e3609a8e097b3b3a156f22ae4d9c680",
			Expected:      "ghcr.io/alexgustafsson/srdl:0.4.4@sha256:919fea57db383a9b161302c08cd5e60c9e3609a8e097b3b3a156f22ae4d9c680",
		},
		{
			SpecImage:     "jacobalberty/unifi:v10.0.162",
			StatusImage:   "docker.io/jacobalberty/unifi:v10.0.162",
			StatusImageID: "docker.io/jacobalberty/unifi@sha256:896c0ab82d33300694dae82982fd7094497afcbea0be92cadc1e94bfead731d3",
			Expected:      "jacobalberty/unifi:v10.0.162@sha256:896c0ab82d33300694dae82982fd7094497afcbea0be92cadc1e94bfead731d3",
		},
		{
			SpecImage:     "dockurr/samba:4.23.5",
			StatusImage:   "docker.io/dockurr/samba:4.23.5",
			StatusImageID: "docker.io/dockurr/samba@sha256:1f0de2ded42bded18d8c5f7b2e9f40521b5d7ddc7cf903ea54a0239cc4984a4e",
			Expected:      "dockurr/samba:4.23.5@sha256:1f0de2ded42bded18d8c5f7b2e9f40521b5d7ddc7cf903ea54a0239cc4984a4e",
		},
		{
			SpecImage:     "yooooomi/your_spotify_client:1.19.0",
			StatusImage:   "docker.io/yooooomi/your_spotify_client:1.19.0",
			StatusImageID: "docker.io/yooooomi/your_spotify_client@sha256:935717b748f56536bd1f4e3bd2c83b71efbcdfb419fbda0345e7a17099a8d30e",
			Expected:      "yooooomi/your_spotify_client:1.19.0@sha256:935717b748f56536bd1f4e3bd2c83b71efbcdfb419fbda0345e7a17099a8d30e",
		},
		{
			SpecImage:     "victoriametrics/victoria-metrics:v1.139.0",
			StatusImage:   "docker.io/victoriametrics/victoria-metrics:v1.139.0",
			StatusImageID: "docker.io/victoriametrics/victoria-metrics@sha256:67c689e152183138ff68eb237ca53d236e92bbef314545733dcb40324829b7c4",
			Expected:      "victoriametrics/victoria-metrics:v1.139.0@sha256:67c689e152183138ff68eb237ca53d236e92bbef314545733dcb40324829b7c4",
		},
		{
			SpecImage:     "archivebox/archivebox:0.7.3",
			StatusImage:   "docker.io/archivebox/archivebox:0.7.3",
			StatusImageID: "docker.io/archivebox/archivebox@sha256:fdf2936192aa1e909b0c3f286f60174efa24078555be4b6b90a07f2cef1d4909",
			Expected:      "archivebox/archivebox:0.7.3@sha256:fdf2936192aa1e909b0c3f286f60174efa24078555be4b6b90a07f2cef1d4909",
		},
		{
			SpecImage:     "registry.home.internal/fixture-calendar/fetcher:latest@sha256:b66ed4ef5979436fc45503a1ae80bd1c493d19e82886227a512b3f5417b4b65f",
			StatusImage:   "sha256:25f89538ba1dcb294f32ce1898dbc881d92a74fc72f863eea31719187fa521a3",
			StatusImageID: "registry.home.internal/fixture-calendar/fetcher@sha256:b66ed4ef5979436fc45503a1ae80bd1c493d19e82886227a512b3f5417b4b65f",
			Expected:      "registry.home.internal/fixture-calendar/fetcher:latest@sha256:b66ed4ef5979436fc45503a1ae80bd1c493d19e82886227a512b3f5417b4b65f",
		},
		{
			SpecImage:     "docker.io/calico/node:v3.28.1",
			StatusImage:   "docker.io/calico/node:v3.28.1",
			StatusImageID: "docker.io/calico/node@sha256:d8c644a8a3eee06d88825b9a9fec6e7cd3b7c276d7f90afa8685a79fb300e7e3",
			Expected:      "calico/node:v3.28.1@sha256:d8c644a8a3eee06d88825b9a9fec6e7cd3b7c276d7f90afa8685a79fb300e7e3",
		},
		{
			SpecImage:     "grafana/grafana:12.4.2",
			StatusImage:   "docker.io/grafana/grafana:12.4.2",
			StatusImageID: "docker.io/grafana/grafana@sha256:83749231c3835e390a3144e5e940203e42b9589761f20ef3169c716e734ad505",
			Expected:      "grafana/grafana:12.4.2@sha256:83749231c3835e390a3144e5e940203e42b9589761f20ef3169c716e734ad505",
		},
		{
			SpecImage:     "ghcr.io/jmbannon/ytdl-sub:2026.03.19",
			StatusImage:   "ghcr.io/jmbannon/ytdl-sub:2026.03.19",
			StatusImageID: "ghcr.io/jmbannon/ytdl-sub@sha256:cd548a095f0e0e72c6c94a9498e41a5d9f122ae8a43d122d26d6463a6a25cd19",
			Expected:      "ghcr.io/jmbannon/ytdl-sub:2026.03.19@sha256:cd548a095f0e0e72c6c94a9498e41a5d9f122ae8a43d122d26d6463a6a25cd19",
		},
		{
			SpecImage:     "yooooomi/your_spotify_server:1.19.0",
			StatusImage:   "docker.io/yooooomi/your_spotify_server:1.19.0",
			StatusImageID: "docker.io/yooooomi/your_spotify_server@sha256:a45776f2c1c24ebcd957f18de4432263907d6f0031c9b25fab3e95f25d15da0d",
			Expected:      "yooooomi/your_spotify_server:1.19.0@sha256:a45776f2c1c24ebcd957f18de4432263907d6f0031c9b25fab3e95f25d15da0d",
		},
		{
			SpecImage:     "hashicorp/vault:1.21.4",
			StatusImage:   "docker.io/hashicorp/vault:1.21.4",
			StatusImageID: "docker.io/hashicorp/vault@sha256:4e33b126a59c0c333b76fb4e894722462659a6bec7c48c9ee8cea56fccfd2569",
			Expected:      "hashicorp/vault:1.21.4@sha256:4e33b126a59c0c333b76fb4e894722462659a6bec7c48c9ee8cea56fccfd2569",
		},
		{
			SpecImage:     "grafana/promtail:3.5.8",
			StatusImage:   "docker.io/grafana/promtail:3.5.8",
			StatusImageID: "docker.io/grafana/promtail@sha256:2a7c5469d687377de5cb7c8356cf96090c0069814d90ef35f9874db445999609",
			Expected:      "grafana/promtail:3.5.8@sha256:2a7c5469d687377de5cb7c8356cf96090c0069814d90ef35f9874db445999609",
		},
		{
			SpecImage:     "mongo:6.0.26",
			StatusImage:   "docker.io/library/mongo:6.0.26",
			StatusImageID: "sha256:bdc4e039b30b99ae50c1068903fab53aeb30859fecc8624099cceecb2d840190",
			Expected:      "mongo:6.0.26@sha256:bdc4e039b30b99ae50c1068903fab53aeb30859fecc8624099cceecb2d840190",
		},
		{
			SpecImage:     "deluan/navidrome:0.61.1",
			StatusImage:   "docker.io/deluan/navidrome:0.61.1",
			StatusImageID: "docker.io/deluan/navidrome@sha256:1e1660054a856cc09f227d6929252e45a519fdb16004b464dd637f7294ca3ec1",
			Expected:      "deluan/navidrome:0.61.1@sha256:1e1660054a856cc09f227d6929252e45a519fdb16004b464dd637f7294ca3ec1",
		},
		{
			SpecImage:     "traefik:3.6.7",
			StatusImage:   "docker.io/library/traefik:3.6.7",
			StatusImageID: "sha256:91528df1690f7da08360dcbbcb92b3ea483eeceb9f136d309f17506a5bd3f58d",
			Expected:      "traefik:3.6.7@sha256:91528df1690f7da08360dcbbcb92b3ea483eeceb9f136d309f17506a5bd3f58d",
		},
		{
			SpecImage:     "quay.io/jetstack/cert-manager-controller:v1.20.1",
			StatusImage:   "quay.io/jetstack/cert-manager-controller:v1.20.1",
			StatusImageID: "quay.io/jetstack/cert-manager-controller@sha256:9f9556b4b131554694c67c8229d231b1f7d69b882b5f061a56bafa465f3b22fc",
			Expected:      "quay.io/jetstack/cert-manager-controller:v1.20.1@sha256:9f9556b4b131554694c67c8229d231b1f7d69b882b5f061a56bafa465f3b22fc",
		},
		{
			SpecImage:     "ghcr.io/alexgustafsson/cupdate:0.24.3",
			StatusImage:   "ghcr.io/alexgustafsson/cupdate:0.24.3",
			StatusImageID: "ghcr.io/alexgustafsson/cupdate@sha256:1da1b8a1976acafa082b76ee87856b579e8f12175e2a13a4eff367a8b317e00d",
			Expected:      "ghcr.io/alexgustafsson/cupdate:0.24.3@sha256:1da1b8a1976acafa082b76ee87856b579e8f12175e2a13a4eff367a8b317e00d",
		},
		{
			SpecImage:     "coredns/coredns:1.10.1",
			StatusImage:   "docker.io/coredns/coredns:1.10.1",
			StatusImageID: "docker.io/coredns/coredns@sha256:a0ead06651cf580044aeb0a0feba63591858fb2e43ade8c9dea45a6a89ae7e5e",
			Expected:      "coredns/coredns:1.10.1@sha256:a0ead06651cf580044aeb0a0feba63591858fb2e43ade8c9dea45a6a89ae7e5e",
		},
		{
			SpecImage:     "grafana/tempo:2.10.3",
			StatusImage:   "docker.io/grafana/tempo:2.10.3",
			StatusImageID: "docker.io/grafana/tempo@sha256:cac9de2ac9f6da8efca5b64b690a7cb8c786a0c49cac7b4517dd1b0089a6c703",
			Expected:      "grafana/tempo:2.10.3@sha256:cac9de2ac9f6da8efca5b64b690a7cb8c786a0c49cac7b4517dd1b0089a6c703",
		},
		{
			SpecImage:     "restic/restic:0.18.1",
			StatusImage:   "docker.io/restic/restic:0.18.1",
			StatusImageID: "docker.io/restic/restic@sha256:39d9072fb5651c80d75c7a811612eb60b4c06b32ffe87c2e9f3c7222e1797e76",
			Expected:      "restic/restic:0.18.1@sha256:39d9072fb5651c80d75c7a811612eb60b4c06b32ffe87c2e9f3c7222e1797e76",
		},
		{
			SpecImage:     "registry.k8s.io/kube-state-metrics/kube-state-metrics:v2.18.0",
			StatusImage:   "registry.k8s.io/kube-state-metrics/kube-state-metrics:v2.18.0",
			StatusImageID: "registry.k8s.io/kube-state-metrics/kube-state-metrics@sha256:1545919b72e3ae035454fc054131e8d0f14b42ef6fc5b2ad5c751cafa6b2130e",
			Expected:      "registry.k8s.io/kube-state-metrics/kube-state-metrics:v2.18.0@sha256:1545919b72e3ae035454fc054131e8d0f14b42ef6fc5b2ad5c751cafa6b2130e",
		},
		{
			SpecImage:     "ghcr.io/alexta69/metube:2026.04.05",
			StatusImage:   "ghcr.io/alexta69/metube:2026.04.05",
			StatusImageID: "ghcr.io/alexta69/metube@sha256:55d4c89fec18669ff56486fe3569eda0d4534319f6bfc46d8c506224e8e735b3",
			Expected:      "ghcr.io/alexta69/metube:2026.04.05@sha256:55d4c89fec18669ff56486fe3569eda0d4534319f6bfc46d8c506224e8e735b3",
		},
		{
			SpecImage:     "curlimages/curl:8.18.0",
			StatusImage:   "docker.io/curlimages/curl:8.18.0",
			StatusImageID: "docker.io/curlimages/curl@sha256:d94d07ba9e7d6de898b6d96c1a072f6f8266c687af78a74f380087a0addf5d17",
			Expected:      "curlimages/curl:8.18.0@sha256:d94d07ba9e7d6de898b6d96c1a072f6f8266c687af78a74f380087a0addf5d17",
		},
		{
			SpecImage:     "b4bz/homer:v26.4.1",
			StatusImage:   "docker.io/b4bz/homer:v26.4.1",
			StatusImageID: "docker.io/b4bz/homer@sha256:659b488ebc52be44ca050b9a990e44c152d99eaa6af0225809168a57f09b67a8",
			Expected:      "b4bz/homer:v26.4.1@sha256:659b488ebc52be44ca050b9a990e44c152d99eaa6af0225809168a57f09b67a8",
		},
		{
			SpecImage:     "grafana/loki:3.7.1",
			StatusImage:   "docker.io/grafana/loki:3.7.1",
			StatusImageID: "docker.io/grafana/loki@sha256:73e905b51a7f917f7a1075e4be68759df30226e03dcb3cd2213b989cc0dc8eb4",
			Expected:      "grafana/loki:3.7.1@sha256:73e905b51a7f917f7a1075e4be68759df30226e03dcb3cd2213b989cc0dc8eb4",
		},
		{
			SpecImage:     "ghcr.io/alexgustafsson/wg-tunnel:latest",
			StatusImage:   "ghcr.io/alexgustafsson/wg-tunnel:latest",
			StatusImageID: "ghcr.io/alexgustafsson/wg-tunnel@sha256:68bce5af155b891ca381bdbee4f3876ec3daac999d615b5ede7b826f0d1ed8e9",
			Expected:      "ghcr.io/alexgustafsson/wg-tunnel:latest@sha256:68bce5af155b891ca381bdbee4f3876ec3daac999d615b5ede7b826f0d1ed8e9",
		},
		{
			SpecImage:     "lipanski/docker-static-website:latest",
			StatusImage:   "docker.io/lipanski/docker-static-website:latest",
			StatusImageID: "docker.io/lipanski/docker-static-website@sha256:66a530684a934a9b94f65a90f286cba291a7daf4dd7d55dcc17f217915056cd5",
			Expected:      "lipanski/docker-static-website:latest@sha256:66a530684a934a9b94f65a90f286cba291a7daf4dd7d55dcc17f217915056cd5",
		},
		{
			SpecImage:     "homeassistant/home-assistant:2026.4.1",
			StatusImage:   "docker.io/homeassistant/home-assistant:2026.4.1",
			StatusImageID: "docker.io/homeassistant/home-assistant@sha256:8848691147f01a6eee7753de2ade21b04d6168fcd2e2a7089f6f84e3b7b86960",
			Expected:      "homeassistant/home-assistant:2026.4.1@sha256:8848691147f01a6eee7753de2ade21b04d6168fcd2e2a7089f6f84e3b7b86960",
		},
		{
			SpecImage:     "alpine:3.23.3",
			StatusImage:   "docker.io/library/alpine:3.23.3",
			StatusImageID: "alpine@sha256:25109184c71bdad752c8312a8623239686a9a2071e8825f20acb8f2198c3f659",
			Expected:      "alpine:3.23.3@sha256:25109184c71bdad752c8312a8623239686a9a2071e8825f20acb8f2198c3f659",
		},
		{
			SpecImage:     "ghcr.io/alexgustafsson/abcde-ui:latest",
			StatusImage:   "ghcr.io/alexgustafsson/abcde-ui:latest",
			StatusImageID: "ghcr.io/alexgustafsson/abcde-ui@sha256:3340d29ed7a0ea599fd71d989e6a10b73cf68768d4dee538201c66e9149fd567",
			Expected:      "ghcr.io/alexgustafsson/abcde-ui:latest@sha256:3340d29ed7a0ea599fd71d989e6a10b73cf68768d4dee538201c66e9149fd567",
		},
		{
			SpecImage:     "ghcr.io/advplyr/audiobookshelf:2.33.1",
			StatusImage:   "ghcr.io/advplyr/audiobookshelf:2.33.1",
			StatusImageID: "ghcr.io/advplyr/audiobookshelf@sha256:a4a5841bba093d81e5f4ad1eaedb4da3fda6dbb2528c552349da50ad1f7ae708",
			Expected:      "ghcr.io/advplyr/audiobookshelf:2.33.1@sha256:a4a5841bba093d81e5f4ad1eaedb4da3fda6dbb2528c552349da50ad1f7ae708",
		},
		{
			SpecImage:     "ghcr.io/alexgustafsson/grapevine:latest",
			StatusImage:   "ghcr.io/alexgustafsson/grapevine:latest",
			StatusImageID: "ghcr.io/alexgustafsson/grapevine@sha256:dcc75b1b13a8d1e4dfa0a6c0516c5c02f4cb4a20fc2f151e9d962bbb78e3a8af",
			Expected:      "ghcr.io/alexgustafsson/grapevine:latest@sha256:dcc75b1b13a8d1e4dfa0a6c0516c5c02f4cb4a20fc2f151e9d962bbb78e3a8af",
		},
		{
			SpecImage:     "rhasspy/wyoming-whisper:3.1.0",
			StatusImage:   "docker.io/rhasspy/wyoming-whisper:3.1.0",
			StatusImageID: "docker.io/rhasspy/wyoming-whisper@sha256:9501d2659eee83b6eead98d53842193e5fed011eda6c5b1c3ad36f3146b28fed",
			Expected:      "rhasspy/wyoming-whisper:3.1.0@sha256:9501d2659eee83b6eead98d53842193e5fed011eda6c5b1c3ad36f3146b28fed",
		},
	}

	for i, testCase := range testCases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			actual, err := getImageReference(testCase.SpecImage, testCase.StatusImage, testCase.StatusImageID)
			assert.Equal(t, testCase.Expected, actual.String())
			if testCase.ExpectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
