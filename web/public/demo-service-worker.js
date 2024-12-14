self.addEventListener('install', () => {
  console.log('Demo service worker installed')
})

self.addEventListener('activate', () => {
  console.log('Demo service worker activated')
})

const imagesResponse = `{
    "images": [
        {
            "reference": "archivebox/archivebox:0.7.2",
            "description": "Official Docker image for the ArchiveBox self-hosted internet archiving tool.",
            "tags": [
                "deployment",
                "replica set"
            ],
            "links": [
                {
                    "type": "oci-registry",
                    "url": "https://docker.io"
                },
                {
                    "type": "docker",
                    "url": "https://hub.docker.com/r/archivebox/archivebox"
                }
            ],
            "vulnerabilities": [],
            "lastModified": "2024-12-14T16:23:40.299978805Z",
            "image": "https://www.gravatar.com/avatar/241d968991f54ef17324ad81e02998a0?s=80&r=g&d=mm"
        },
        {
            "reference": "b4bz/homer:v24.12.1",
            "latestReference": "b4bz/homer:v24.12.1",
            "description": "A dead simple static HOMe for your servER to keep your services on hand from a simple yaml config.",
            "tags": [
                "deployment",
                "replica set",
                "up-to-date"
            ],
            "links": [
                {
                    "type": "oci-registry",
                    "url": "https://docker.io"
                },
                {
                    "type": "docker",
                    "url": "https://hub.docker.com/r/b4bz/homer"
                }
            ],
            "vulnerabilities": [],
            "lastModified": "2024-12-14T17:44:57.641311764Z"
        },
        {
            "reference": "calico/kube-controllers:v3.23.5",
            "tags": [
                "deployment",
                "replica set"
            ],
            "links": [
                {
                    "type": "oci-registry",
                    "url": "https://docker.io"
                },
                {
                    "type": "docker",
                    "url": "https://hub.docker.com/r/calico/kube-controllers"
                }
            ],
            "vulnerabilities": [],
            "lastModified": "2024-12-14T19:45:06.48395342Z"
        },
        {
            "reference": "calico/node:v3.23.5",
            "description": "Calico's per-host DaemonSet container image.  Provides CNI networking and policy for Kubernetes.",
            "tags": [
                "daemon set"
            ],
            "links": [
                {
                    "type": "oci-registry",
                    "url": "https://docker.io"
                },
                {
                    "type": "docker",
                    "url": "https://hub.docker.com/r/calico/node"
                }
            ],
            "vulnerabilities": [],
            "lastModified": "2024-12-14T19:45:06.786202786Z"
        },
        {
            "reference": "coredns/coredns:1.9.3",
            "latestReference": "coredns/coredns:1.12.0",
            "description": "CoreDNS docker repository",
            "tags": [
                "deployment",
                "minor",
                "outdated",
                "replica set"
            ],
            "links": [
                {
                    "type": "oci-registry",
                    "url": "https://docker.io"
                },
                {
                    "type": "docker",
                    "url": "https://hub.docker.com/r/coredns/coredns"
                }
            ],
            "vulnerabilities": [],
            "lastModified": "2024-12-14T16:45:00.467533323Z"
        },
        {
            "reference": "ghcr.io/advplyr/audiobookshelf:2.17.4",
            "latestReference": "ghcr.io/advplyr/audiobookshelf:2.17.5",
            "description": "Self-hosted audiobook and podcast server",
            "tags": [
                "deployment",
                "github",
                "outdated",
                "patch",
                "replica set"
            ],
            "links": [
                {
                    "type": "oci-registry",
                    "url": "https://ghcr.io"
                },
                {
                    "type": "github",
                    "url": "https://github.com/advplyr/audiobookshelf"
                },
                {
                    "type": "ghcr",
                    "url": "https://github.com/users/advplyr/packages/container/package/audiobookshelf"
                }
            ],
            "vulnerabilities": [],
            "lastModified": "2024-12-14T17:45:02.054892928Z"
        },
        {
            "reference": "ghcr.io/alexgustafsson/cupdate",
            "latestReference": "ghcr.io/alexgustafsson/cupdate",
            "description": "A WIP service to keep container images up-to-date in k8s and more",
            "tags": [
                "deployment",
                "github",
                "replica set",
                "up-to-date"
            ],
            "links": [
                {
                    "type": "oci-registry",
                    "url": "https://ghcr.io"
                },
                {
                    "type": "github",
                    "url": "https://github.com/alexgustafsson/cupdate"
                },
                {
                    "type": "ghcr",
                    "url": "https://github.com/users/alexgustafsson/packages/container/package/cupdate"
                }
            ],
            "vulnerabilities": [],
            "lastModified": "2024-12-13T17:00:55.584501239Z"
        },
        {
            "reference": "ghcr.io/alexgustafsson/srdl:0.3.1",
            "latestReference": "ghcr.io/alexgustafsson/srdl:0.3.1",
            "description": "Like ytdl and ytdl-sub but for Sveriges Radio - archive programs from SR",
            "tags": [
                "cron job",
                "github",
                "job",
                "up-to-date"
            ],
            "links": [
                {
                    "type": "oci-registry",
                    "url": "https://ghcr.io"
                },
                {
                    "type": "github",
                    "url": "https://github.com/alexgustafsson/srdl"
                },
                {
                    "type": "ghcr",
                    "url": "https://github.com/users/alexgustafsson/packages/container/package/srdl"
                }
            ],
            "vulnerabilities": [],
            "lastModified": "2024-12-14T20:45:10.240183311Z"
        },
        {
            "reference": "ghcr.io/jmbannon/ytdl-sub:2024.12.08",
            "latestReference": "ghcr.io/jmbannon/ytdl-sub:2024.12.08",
            "description": "Lightweight tool to automate downloading and metadata generation with yt-dlp",
            "tags": [
                "cron job",
                "github",
                "job",
                "up-to-date"
            ],
            "links": [
                {
                    "type": "oci-registry",
                    "url": "https://ghcr.io"
                },
                {
                    "type": "github",
                    "url": "https://github.com/jmbannon/ytdl-sub"
                },
                {
                    "type": "ghcr",
                    "url": "https://github.com/users/jmbannon/packages/container/package/ytdl-sub"
                }
            ],
            "vulnerabilities": [],
            "lastModified": "2024-12-14T18:44:56.678743971Z"
        },
        {
            "reference": "grafana/grafana:11.4.0",
            "latestReference": "grafana/grafana:11.4.0",
            "description": "The official Grafana docker container",
            "tags": [
                "deployment",
                "replica set",
                "up-to-date"
            ],
            "links": [
                {
                    "type": "oci-registry",
                    "url": "https://docker.io"
                },
                {
                    "type": "docker",
                    "url": "https://hub.docker.com/r/grafana/grafana"
                }
            ],
            "vulnerabilities": [],
            "lastModified": "2024-12-14T17:44:55.339094521Z",
            "image": "https://www.gravatar.com/avatar/31cea69afa424609b2d83621b4d47f1d?s=80&r=g&d=mm"
        },
        {
            "reference": "grafana/loki:3.2.2",
            "latestReference": "grafana/loki:3.2.2",
            "description": "Loki - Cloud Native Log Aggregation by Grafana",
            "tags": [
                "deployment",
                "replica set",
                "up-to-date"
            ],
            "links": [
                {
                    "type": "oci-registry",
                    "url": "https://docker.io"
                },
                {
                    "type": "docker",
                    "url": "https://hub.docker.com/r/grafana/loki"
                }
            ],
            "vulnerabilities": [],
            "lastModified": "2024-12-14T17:44:58.347502737Z",
            "image": "https://www.gravatar.com/avatar/31cea69afa424609b2d83621b4d47f1d?s=80&r=g&d=mm"
        },
        {
            "reference": "grafana/promtail:3.2.2",
            "latestReference": "grafana/promtail:3.2.2",
            "tags": [
                "daemon set",
                "up-to-date"
            ],
            "links": [
                {
                    "type": "oci-registry",
                    "url": "https://docker.io"
                },
                {
                    "type": "docker",
                    "url": "https://hub.docker.com/r/grafana/promtail"
                }
            ],
            "vulnerabilities": [],
            "lastModified": "2024-12-14T17:44:53.738665245Z",
            "image": "https://www.gravatar.com/avatar/31cea69afa424609b2d83621b4d47f1d?s=80&r=g&d=mm"
        },
        {
            "reference": "hashicorp/vault:1.18.1",
            "latestReference": "hashicorp/vault:1.18.2",
            "description": "Official vault docker images",
            "tags": [
                "outdated",
                "patch",
                "stateful set"
            ],
            "links": [
                {
                    "type": "oci-registry",
                    "url": "https://docker.io"
                },
                {
                    "type": "docker",
                    "url": "https://hub.docker.com/r/hashicorp/vault"
                }
            ],
            "vulnerabilities": [],
            "lastModified": "2024-12-14T16:45:23.53962115Z",
            "image": "https://www.gravatar.com/avatar/f7304a610df968cef5badc9e9dcd40e0?s=80&r=g&d=mm"
        },
        {
            "reference": "homeassistant/home-assistant:2024.12.2",
            "latestReference": "homeassistant/home-assistant:2024.12.3",
            "description": "Open source home automation that puts local control and privacy first. ",
            "tags": [
                "deployment",
                "outdated",
                "patch",
                "replica set"
            ],
            "links": [
                {
                    "type": "oci-registry",
                    "url": "https://docker.io"
                },
                {
                    "type": "svc",
                    "url": "https://github.com/home-assistant/core"
                },
                {
                    "type": "docker",
                    "url": "https://hub.docker.com/r/homeassistant/home-assistant"
                }
            ],
            "vulnerabilities": [],
            "lastModified": "2024-12-14T18:44:54.471259644Z",
            "image": "https://www.gravatar.com/avatar/461df105cc6cfcf386ebd5af85b925dc?s=80&r=g&d=mm"
        },
        {
            "reference": "intel/intel-gpu-plugin:0.31.1",
            "latestReference": "intel/intel-gpu-plugin:0.31.1",
            "description": "Intel GPU device plugin for Kubernetes",
            "tags": [
                "daemon set",
                "up-to-date"
            ],
            "links": [
                {
                    "type": "oci-registry",
                    "url": "https://docker.io"
                },
                {
                    "type": "docker",
                    "url": "https://hub.docker.com/r/intel/intel-gpu-plugin"
                }
            ],
            "vulnerabilities": [],
            "lastModified": "2024-12-14T18:44:53.094744187Z",
            "image": "https://www.gravatar.com/avatar/9d2eb115ddd65132d6179073c8552e5e?s=80&r=g&d=mm"
        },
        {
            "reference": "jacobalberty/unifi:v8.6.9",
            "latestReference": "jacobalberty/unifi:v8.6.9",
            "description": "Unifi Access Point controller",
            "tags": [
                "deployment",
                "replica set",
                "up-to-date"
            ],
            "links": [
                {
                    "type": "oci-registry",
                    "url": "https://docker.io"
                },
                {
                    "type": "svc",
                    "url": "https://github.com/jacobalberty/unifi-docker"
                },
                {
                    "type": "docker",
                    "url": "https://hub.docker.com/r/jacobalberty/unifi"
                }
            ],
            "vulnerabilities": [],
            "lastModified": "2024-12-14T20:45:10.871097653Z"
        },
        {
            "reference": "lscr.io/linuxserver/jellyfin",
            "latestReference": "lscr.io/linuxserver/jellyfin",
            "tags": [
                "deployment",
                "replica set",
                "up-to-date"
            ],
            "links": [
                {
                    "type": "svc",
                    "url": "https://github.com/linuxserver/docker-jellyfin"
                },
                {
                    "type": "oci-registry",
                    "url": "https://lscr.io"
                }
            ],
            "vulnerabilities": [],
            "lastModified": "2024-12-14T05:00:45.779295332Z"
        },
        {
            "reference": "mongo:6.0.19",
            "description": "MongoDB document databases provide high availability and easy scalability.",
            "tags": [
                "deployment",
                "replica set",
                "vulnerable"
            ],
            "links": [
                {
                    "type": "oci-registry",
                    "url": "https://docker.io"
                },
                {
                    "type": "docker",
                    "url": "https://hub.docker.com/_/mongo"
                }
            ],
            "vulnerabilities": [
                {
                    "id": 1,
                    "severity": "critical",
                    "authority": "Docker Scout",
                    "link": "https://hub.docker.com/layers/library/mongo/6.0.19/images/sha256-ebad181937de72a6226b39a63eb92b26406cf0f3bd44b5d92810264c93b76078"
                },
                {
                    "id": 2,
                    "severity": "critical",
                    "authority": "Docker Scout",
                    "link": "https://hub.docker.com/layers/library/mongo/6.0.19/images/sha256-ebad181937de72a6226b39a63eb92b26406cf0f3bd44b5d92810264c93b76078"
                },
                {
                    "id": 3,
                    "severity": "critical",
                    "authority": "Docker Scout",
                    "link": "https://hub.docker.com/layers/library/mongo/6.0.19/images/sha256-ebad181937de72a6226b39a63eb92b26406cf0f3bd44b5d92810264c93b76078"
                },
                {
                    "id": 4,
                    "severity": "high",
                    "authority": "Docker Scout",
                    "link": "https://hub.docker.com/layers/library/mongo/6.0.19/images/sha256-ebad181937de72a6226b39a63eb92b26406cf0f3bd44b5d92810264c93b76078"
                },
                {
                    "id": 5,
                    "severity": "high",
                    "authority": "Docker Scout",
                    "link": "https://hub.docker.com/layers/library/mongo/6.0.19/images/sha256-ebad181937de72a6226b39a63eb92b26406cf0f3bd44b5d92810264c93b76078"
                },
                {
                    "id": 6,
                    "severity": "high",
                    "authority": "Docker Scout",
                    "link": "https://hub.docker.com/layers/library/mongo/6.0.19/images/sha256-ebad181937de72a6226b39a63eb92b26406cf0f3bd44b5d92810264c93b76078"
                },
                {
                    "id": 7,
                    "severity": "high",
                    "authority": "Docker Scout",
                    "link": "https://hub.docker.com/layers/library/mongo/6.0.19/images/sha256-ebad181937de72a6226b39a63eb92b26406cf0f3bd44b5d92810264c93b76078"
                },
                {
                    "id": 8,
                    "severity": "high",
                    "authority": "Docker Scout",
                    "link": "https://hub.docker.com/layers/library/mongo/6.0.19/images/sha256-ebad181937de72a6226b39a63eb92b26406cf0f3bd44b5d92810264c93b76078"
                },
                {
                    "id": 9,
                    "severity": "medium",
                    "authority": "Docker Scout",
                    "link": "https://hub.docker.com/layers/library/mongo/6.0.19/images/sha256-ebad181937de72a6226b39a63eb92b26406cf0f3bd44b5d92810264c93b76078"
                },
                {
                    "id": 10,
                    "severity": "medium",
                    "authority": "Docker Scout",
                    "link": "https://hub.docker.com/layers/library/mongo/6.0.19/images/sha256-ebad181937de72a6226b39a63eb92b26406cf0f3bd44b5d92810264c93b76078"
                },
                {
                    "id": 11,
                    "severity": "medium",
                    "authority": "Docker Scout",
                    "link": "https://hub.docker.com/layers/library/mongo/6.0.19/images/sha256-ebad181937de72a6226b39a63eb92b26406cf0f3bd44b5d92810264c93b76078"
                },
                {
                    "id": 12,
                    "severity": "medium",
                    "authority": "Docker Scout",
                    "link": "https://hub.docker.com/layers/library/mongo/6.0.19/images/sha256-ebad181937de72a6226b39a63eb92b26406cf0f3bd44b5d92810264c93b76078"
                },
                {
                    "id": 13,
                    "severity": "medium",
                    "authority": "Docker Scout",
                    "link": "https://hub.docker.com/layers/library/mongo/6.0.19/images/sha256-ebad181937de72a6226b39a63eb92b26406cf0f3bd44b5d92810264c93b76078"
                },
                {
                    "id": 14,
                    "severity": "low",
                    "authority": "Docker Scout",
                    "link": "https://hub.docker.com/layers/library/mongo/6.0.19/images/sha256-ebad181937de72a6226b39a63eb92b26406cf0f3bd44b5d92810264c93b76078"
                },
                {
                    "id": 15,
                    "severity": "low",
                    "authority": "Docker Scout",
                    "link": "https://hub.docker.com/layers/library/mongo/6.0.19/images/sha256-ebad181937de72a6226b39a63eb92b26406cf0f3bd44b5d92810264c93b76078"
                },
                {
                    "id": 16,
                    "severity": "low",
                    "authority": "Docker Scout",
                    "link": "https://hub.docker.com/layers/library/mongo/6.0.19/images/sha256-ebad181937de72a6226b39a63eb92b26406cf0f3bd44b5d92810264c93b76078"
                }
            ],
            "lastModified": "2024-12-14T17:44:54.643865327Z",
            "image": "https://www.gravatar.com/avatar/7510e100f7ebeca4a0b8c3c617349295?s=80&r=g&d=mm"
        },
        {
            "reference": "prom/node-exporter:v1.8.2",
            "latestReference": "prom/node-exporter:v1.8.2",
            "description": "prom/node-exporter",
            "tags": [
                "daemon set",
                "up-to-date"
            ],
            "links": [
                {
                    "type": "oci-registry",
                    "url": "https://docker.io"
                },
                {
                    "type": "docker",
                    "url": "https://hub.docker.com/r/prom/node-exporter"
                }
            ],
            "vulnerabilities": [],
            "lastModified": "2024-12-14T17:44:53.036548387Z",
            "image": "https://www.gravatar.com/avatar/63f4dc6944f814f3b2440a3c41dd400b?s=80&r=g&d=mm"
        },
        {
            "reference": "quay.io/jetstack/cert-manager-cainjector:v1.16.2",
            "latestReference": "quay.io/jetstack/cert-manager-cainjector:v1.16.2",
            "tags": [
                "deployment",
                "replica set",
                "up-to-date"
            ],
            "links": [
                {
                    "type": "svc",
                    "url": "https://github.com/cert-manager/cert-manager"
                },
                {
                    "type": "oci-registry",
                    "url": "https://quay.io"
                }
            ],
            "vulnerabilities": [],
            "lastModified": "2024-12-14T19:45:06.81965899Z"
        },
        {
            "reference": "quay.io/jetstack/cert-manager-controller:v1.16.2",
            "latestReference": "quay.io/jetstack/cert-manager-controller:v1.16.2",
            "tags": [
                "deployment",
                "replica set",
                "up-to-date"
            ],
            "links": [
                {
                    "type": "svc",
                    "url": "https://github.com/cert-manager/cert-manager"
                },
                {
                    "type": "oci-registry",
                    "url": "https://quay.io"
                }
            ],
            "vulnerabilities": [],
            "lastModified": "2024-12-14T19:45:06.80501198Z"
        },
        {
            "reference": "quay.io/jetstack/cert-manager-startupapicheck:v1.14.4",
            "latestReference": "quay.io/jetstack/cert-manager-startupapicheck:v1.16.2",
            "tags": [
                "job",
                "minor",
                "outdated"
            ],
            "links": [
                {
                    "type": "svc",
                    "url": "https://github.com/cert-manager/cert-manager"
                },
                {
                    "type": "oci-registry",
                    "url": "https://quay.io"
                }
            ],
            "vulnerabilities": [],
            "lastModified": "2024-12-14T16:45:11.841156462Z"
        },
        {
            "reference": "quay.io/jetstack/cert-manager-webhook:v1.16.2",
            "latestReference": "quay.io/jetstack/cert-manager-webhook:v1.16.2",
            "tags": [
                "deployment",
                "replica set",
                "up-to-date"
            ],
            "links": [
                {
                    "type": "svc",
                    "url": "https://github.com/cert-manager/cert-manager"
                },
                {
                    "type": "oci-registry",
                    "url": "https://quay.io"
                }
            ],
            "vulnerabilities": [],
            "lastModified": "2024-12-14T20:45:08.805075046Z"
        },
        {
            "reference": "registry.gitlab.com/arm-research/smarter/smarter-device-manager:v1.20.11",
            "latestReference": "registry.gitlab.com/arm-research/smarter/smarter-device-manager:v1.20.11",
            "description": "Provides a device manager container that enables access to device drivers on containers for K8s",
            "tags": [
                "daemon set",
                "up-to-date"
            ],
            "links": [
                {
                    "type": "gitlab",
                    "url": "https://gitlab.com/arm-research/smarter/smarter-device-manager"
                },
                {
                    "type": "oci-registry",
                    "url": "https://registry.gitlab.com"
                }
            ],
            "vulnerabilities": [],
            "lastModified": "2024-12-14T19:45:05.811741904Z"
        },
        {
            "reference": "registry.k8s.io/kube-state-metrics/kube-state-metrics:v2.14.0",
            "tags": [
                "deployment",
                "replica set"
            ],
            "links": [
                {
                    "type": "oci-registry",
                    "url": "https://registry.k8s.io"
                }
            ],
            "vulnerabilities": [],
            "lastModified": "2024-12-14T18:44:52.814869095Z"
        },
        {
            "reference": "traefik:v3.2.2",
            "description": "Traefik, The Cloud Native Edge Router",
            "tags": [
                "daemon set",
                "vulnerable"
            ],
            "links": [
                {
                    "type": "oci-registry",
                    "url": "https://docker.io"
                },
                {
                    "type": "svc",
                    "url": "https://github.com/traefik/traefik"
                },
                {
                    "type": "docker",
                    "url": "https://hub.docker.com/_/traefik"
                }
            ],
            "vulnerabilities": [
                {
                    "id": 17,
                    "severity": "critical",
                    "authority": "Docker Scout",
                    "link": "https://hub.docker.com/layers/library/traefik/v3.2.2/images/sha256-f288eb36fde276f07a9827bd95fa55cda0c502112bfe51a09396242c83f3fdb6"
                },
                {
                    "id": 18,
                    "severity": "critical",
                    "authority": "Docker Scout",
                    "link": "https://hub.docker.com/layers/library/traefik/v3.2.2/images/sha256-f288eb36fde276f07a9827bd95fa55cda0c502112bfe51a09396242c83f3fdb6"
                },
                {
                    "id": 19,
                    "severity": "critical",
                    "authority": "Docker Scout",
                    "link": "https://hub.docker.com/layers/library/traefik/v3.2.2/images/sha256-f288eb36fde276f07a9827bd95fa55cda0c502112bfe51a09396242c83f3fdb6"
                },
                {
                    "id": 20,
                    "severity": "high",
                    "authority": "Docker Scout",
                    "link": "https://hub.docker.com/layers/library/traefik/v3.2.2/images/sha256-f288eb36fde276f07a9827bd95fa55cda0c502112bfe51a09396242c83f3fdb6"
                },
                {
                    "id": 21,
                    "severity": "high",
                    "authority": "Docker Scout",
                    "link": "https://hub.docker.com/layers/library/traefik/v3.2.2/images/sha256-f288eb36fde276f07a9827bd95fa55cda0c502112bfe51a09396242c83f3fdb6"
                },
                {
                    "id": 22,
                    "severity": "high",
                    "authority": "Docker Scout",
                    "link": "https://hub.docker.com/layers/library/traefik/v3.2.2/images/sha256-f288eb36fde276f07a9827bd95fa55cda0c502112bfe51a09396242c83f3fdb6"
                },
                {
                    "id": 23,
                    "severity": "high",
                    "authority": "Docker Scout",
                    "link": "https://hub.docker.com/layers/library/traefik/v3.2.2/images/sha256-f288eb36fde276f07a9827bd95fa55cda0c502112bfe51a09396242c83f3fdb6"
                },
                {
                    "id": 24,
                    "severity": "high",
                    "authority": "Docker Scout",
                    "link": "https://hub.docker.com/layers/library/traefik/v3.2.2/images/sha256-f288eb36fde276f07a9827bd95fa55cda0c502112bfe51a09396242c83f3fdb6"
                },
                {
                    "id": 25,
                    "severity": "medium",
                    "authority": "Docker Scout",
                    "link": "https://hub.docker.com/layers/library/traefik/v3.2.2/images/sha256-f288eb36fde276f07a9827bd95fa55cda0c502112bfe51a09396242c83f3fdb6"
                },
                {
                    "id": 26,
                    "severity": "medium",
                    "authority": "Docker Scout",
                    "link": "https://hub.docker.com/layers/library/traefik/v3.2.2/images/sha256-f288eb36fde276f07a9827bd95fa55cda0c502112bfe51a09396242c83f3fdb6"
                },
                {
                    "id": 27,
                    "severity": "medium",
                    "authority": "Docker Scout",
                    "link": "https://hub.docker.com/layers/library/traefik/v3.2.2/images/sha256-f288eb36fde276f07a9827bd95fa55cda0c502112bfe51a09396242c83f3fdb6"
                },
                {
                    "id": 28,
                    "severity": "medium",
                    "authority": "Docker Scout",
                    "link": "https://hub.docker.com/layers/library/traefik/v3.2.2/images/sha256-f288eb36fde276f07a9827bd95fa55cda0c502112bfe51a09396242c83f3fdb6"
                },
                {
                    "id": 29,
                    "severity": "medium",
                    "authority": "Docker Scout",
                    "link": "https://hub.docker.com/layers/library/traefik/v3.2.2/images/sha256-f288eb36fde276f07a9827bd95fa55cda0c502112bfe51a09396242c83f3fdb6"
                },
                {
                    "id": 30,
                    "severity": "low",
                    "authority": "Docker Scout",
                    "link": "https://hub.docker.com/layers/library/traefik/v3.2.2/images/sha256-f288eb36fde276f07a9827bd95fa55cda0c502112bfe51a09396242c83f3fdb6"
                },
                {
                    "id": 31,
                    "severity": "low",
                    "authority": "Docker Scout",
                    "link": "https://hub.docker.com/layers/library/traefik/v3.2.2/images/sha256-f288eb36fde276f07a9827bd95fa55cda0c502112bfe51a09396242c83f3fdb6"
                },
                {
                    "id": 32,
                    "severity": "low",
                    "authority": "Docker Scout",
                    "link": "https://hub.docker.com/layers/library/traefik/v3.2.2/images/sha256-f288eb36fde276f07a9827bd95fa55cda0c502112bfe51a09396242c83f3fdb6"
                }
            ],
            "lastModified": "2024-12-14T18:44:55.219754888Z",
            "image": "https://www.gravatar.com/avatar/7510e100f7ebeca4a0b8c3c617349295?s=80&r=g&d=mm"
        },
        {
            "reference": "victoriametrics/victoria-metrics:v1.108.0",
            "latestReference": "victoriametrics/victoria-metrics:v1.108.0",
            "description": "Single-node version of VictoriaMetrics",
            "tags": [
                "deployment",
                "replica set",
                "up-to-date"
            ],
            "links": [
                {
                    "type": "oci-registry",
                    "url": "https://docker.io"
                },
                {
                    "type": "docker",
                    "url": "https://hub.docker.com/r/victoriametrics/victoria-metrics"
                }
            ],
            "vulnerabilities": [],
            "lastModified": "2024-12-14T17:44:51.651413508Z",
            "image": "https://www.gravatar.com/avatar/a625fcc372d7c43a7943d031fbd88942?s=80&r=g&d=mm"
        },
        {
            "reference": "victoriametrics/vmagent:v1.108.0",
            "latestReference": "victoriametrics/vmagent:v1.108.0",
            "description": "Agent for collecting metrics from various sources, filtering and sending them to VictoriaMetrics",
            "tags": [
                "deployment",
                "replica set",
                "up-to-date"
            ],
            "links": [
                {
                    "type": "oci-registry",
                    "url": "https://docker.io"
                },
                {
                    "type": "docker",
                    "url": "https://hub.docker.com/r/victoriametrics/vmagent"
                }
            ],
            "vulnerabilities": [],
            "lastModified": "2024-12-14T18:44:51.677428614Z",
            "image": "https://www.gravatar.com/avatar/a625fcc372d7c43a7943d031fbd88942?s=80&r=g&d=mm"
        },
        {
            "reference": "yooooomi/your_spotify_client:1.12.0",
            "latestReference": "yooooomi/your_spotify_client:1.12.0",
            "tags": [
                "deployment",
                "replica set",
                "up-to-date"
            ],
            "links": [
                {
                    "type": "oci-registry",
                    "url": "https://docker.io"
                },
                {
                    "type": "docker",
                    "url": "https://hub.docker.com/r/yooooomi/your_spotify_client"
                }
            ],
            "vulnerabilities": [],
            "lastModified": "2024-12-14T17:44:52.344746698Z"
        }
    ],
    "summary": {
        "images": 30,
        "outdated": 6,
        "vulnerable": 2,
        "processing": 0
    },
    "pagination": {
        "total": 30,
        "page": 0,
        "size": 30
    }
}`

const tagsResponse = `["replica set", "outdated", "deployment", "github", "up-to-date", "daemon set", "vulnerable", "minor", "job", "stateful set", "patch", "cron job"]`

/** @type {{matcher: RegExp, response: {body?: BodyInit, options?: ResponseInit}}[]} */
const data = [
  {
    matcher: /^GET \/api\/v1\/images.*$/,
    response: {
      body: imagesResponse,
      options: {
        status: 200,
        statusText: 'OK',
        headers: {},
      },
    },
  },
  {
    matcher: /^GET \/api\/v1\/tags$/,
    response: {
      body: tagsResponse,
      options: {
        status: 200,
        statusText: 'OK',
        headers: {},
      },
    },
  },
]

self.addEventListener('fetch', (event) => {
  const url = new URL(event.request.url)
  if (url.host === self.location.host && url.pathname.startsWith('/api/v1/')) {
    const key = `${event.request.method} ${url.pathname}${url.search}`
    for (const { matcher, response } of data) {
      if (matcher.test(key)) {
        event.respondWith(new Response(response.body, response.options))
        return
      }
    }
  }

  event.respondWith(
    fetch(event.request).catch(function () {
      return caches.match('/offline')
    })
  )
})
