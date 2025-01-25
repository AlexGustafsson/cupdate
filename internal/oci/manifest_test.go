package oci

import (
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestManifestFromBlob(t *testing.T) {
	testCases := []struct {
		Name        string
		ContentType string
		JSON        string
		Expected    any
		ExpectError bool
	}{
		{
			Name:        "Image Manifest Version 2, Schema 1",
			ContentType: "application/vnd.docker.distribution.manifest.v1+json",
			JSON: `{
   "name": "hello-world",
   "tag": "latest",
   "architecture": "amd64",
   "fsLayers": [
      {
         "blobSum": "sha256:5f70bf18a086007016e948b04aed3b82103a36bea41755b6cddfaf10ace3c6ef"
      },
      {
         "blobSum": "sha256:5f70bf18a086007016e948b04aed3b82103a36bea41755b6cddfaf10ace3c6ef"
      },
      {
         "blobSum": "sha256:cc8567d70002e957612902a8e985ea129d831ebe04057d88fb644857caa45d11"
      },
      {
         "blobSum": "sha256:5f70bf18a086007016e948b04aed3b82103a36bea41755b6cddfaf10ace3c6ef"
      }
   ],
   "history": [
      {
         "v1Compatibility": "{\"id\":\"e45a5af57b00862e5ef5782a9925979a02ba2b12dff832fd0991335f4a11e5c5\",\"parent\":\"31cbccb51277105ba3ae35ce33c22b69c9e3f1002e76e4c736a2e8ebff9d7b5d\",\"created\":\"2014-12-31T22:57:59.178729048Z\",\"container\":\"27b45f8fb11795b52e9605b686159729b0d9ca92f76d40fb4f05a62e19c46b4f\",\"container_config\":{\"Hostname\":\"8ce6509d66e2\",\"Domainname\":\"\",\"User\":\"\",\"Memory\":0,\"MemorySwap\":0,\"CpuShares\":0,\"Cpuset\":\"\",\"AttachStdin\":false,\"AttachStdout\":false,\"AttachStderr\":false,\"PortSpecs\":null,\"ExposedPorts\":null,\"Tty\":false,\"OpenStdin\":false,\"StdinOnce\":false,\"Env\":[\"PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin\"],\"Cmd\":[\"/bin/sh\",\"-c\",\"#(nop) CMD [/hello]\"],\"Image\":\"31cbccb51277105ba3ae35ce33c22b69c9e3f1002e76e4c736a2e8ebff9d7b5d\",\"Volumes\":null,\"WorkingDir\":\"\",\"Entrypoint\":null,\"NetworkDisabled\":false,\"MacAddress\":\"\",\"OnBuild\":[],\"SecurityOpt\":null,\"Labels\":null},\"docker_version\":\"1.4.1\",\"config\":{\"Hostname\":\"8ce6509d66e2\",\"Domainname\":\"\",\"User\":\"\",\"Memory\":0,\"MemorySwap\":0,\"CpuShares\":0,\"Cpuset\":\"\",\"AttachStdin\":false,\"AttachStdout\":false,\"AttachStderr\":false,\"PortSpecs\":null,\"ExposedPorts\":null,\"Tty\":false,\"OpenStdin\":false,\"StdinOnce\":false,\"Env\":[\"PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin\"],\"Cmd\":[\"/hello\"],\"Image\":\"31cbccb51277105ba3ae35ce33c22b69c9e3f1002e76e4c736a2e8ebff9d7b5d\",\"Volumes\":null,\"WorkingDir\":\"\",\"Entrypoint\":null,\"NetworkDisabled\":false,\"MacAddress\":\"\",\"OnBuild\":[],\"SecurityOpt\":null,\"Labels\":null},\"architecture\":\"amd64\",\"os\":\"linux\",\"Size\":0}\n"
      },
      {
         "v1Compatibility": "{\"id\":\"e45a5af57b00862e5ef5782a9925979a02ba2b12dff832fd0991335f4a11e5c5\",\"parent\":\"31cbccb51277105ba3ae35ce33c22b69c9e3f1002e76e4c736a2e8ebff9d7b5d\",\"created\":\"2014-12-31T22:57:59.178729048Z\",\"container\":\"27b45f8fb11795b52e9605b686159729b0d9ca92f76d40fb4f05a62e19c46b4f\",\"container_config\":{\"Hostname\":\"8ce6509d66e2\",\"Domainname\":\"\",\"User\":\"\",\"Memory\":0,\"MemorySwap\":0,\"CpuShares\":0,\"Cpuset\":\"\",\"AttachStdin\":false,\"AttachStdout\":false,\"AttachStderr\":false,\"PortSpecs\":null,\"ExposedPorts\":null,\"Tty\":false,\"OpenStdin\":false,\"StdinOnce\":false,\"Env\":[\"PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin\"],\"Cmd\":[\"/bin/sh\",\"-c\",\"#(nop) CMD [/hello]\"],\"Image\":\"31cbccb51277105ba3ae35ce33c22b69c9e3f1002e76e4c736a2e8ebff9d7b5d\",\"Volumes\":null,\"WorkingDir\":\"\",\"Entrypoint\":null,\"NetworkDisabled\":false,\"MacAddress\":\"\",\"OnBuild\":[],\"SecurityOpt\":null,\"Labels\":null},\"docker_version\":\"1.4.1\",\"config\":{\"Hostname\":\"8ce6509d66e2\",\"Domainname\":\"\",\"User\":\"\",\"Memory\":0,\"MemorySwap\":0,\"CpuShares\":0,\"Cpuset\":\"\",\"AttachStdin\":false,\"AttachStdout\":false,\"AttachStderr\":false,\"PortSpecs\":null,\"ExposedPorts\":null,\"Tty\":false,\"OpenStdin\":false,\"StdinOnce\":false,\"Env\":[\"PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin\"],\"Cmd\":[\"/hello\"],\"Image\":\"31cbccb51277105ba3ae35ce33c22b69c9e3f1002e76e4c736a2e8ebff9d7b5d\",\"Volumes\":null,\"WorkingDir\":\"\",\"Entrypoint\":null,\"NetworkDisabled\":false,\"MacAddress\":\"\",\"OnBuild\":[],\"SecurityOpt\":null,\"Labels\":null},\"architecture\":\"amd64\",\"os\":\"linux\",\"Size\":0}\n"
      }
   ],
   "schemaVersion": 1,
   "signatures": [
      {
         "header": {
            "jwk": {
               "crv": "P-256",
               "kid": "OD6I:6DRK:JXEJ:KBM4:255X:NSAA:MUSF:E4VM:ZI6W:CUN2:L4Z6:LSF4",
               "kty": "EC",
               "x": "3gAwX48IQ5oaYQAYSxor6rYYc_6yjuLCjtQ9LUakg4A",
               "y": "t72ge6kIA1XOjqjVoEOiPPAURltJFBMGDSQvEGVB010"
            },
            "alg": "ES256"
         },
         "signature": "XREm0L8WNn27Ga_iE_vRnTxVMhhYY0Zst_FfkKopg6gWSoTOZTuW4rK0fg_IqnKkEKlbD83tD46LKEGi5aIVFg",
         "protected": "eyJmb3JtYXRMZW5ndGgiOjY2MjgsImZvcm1hdFRhaWwiOiJDbjAiLCJ0aW1lIjoiMjAxNS0wNC0wOFQxODo1Mjo1OVoifQ"
      }
   ]
}`,
			Expected: &ImageManifest{
				ContentType:   "application/vnd.docker.distribution.manifest.v1+json",
				SchemaVersion: 1,
				Digest:        "sha256:44d1afbd01feff689f7f24808298e2fba2647292c17cb6ef9f2041cafb8fe496",
				Platform: &Platform{
					Architecture: "amd64",
				},
				Annotations: Annotations{},
			},
		},
		{
			Name:        "Image Manifest Version 2, Schema 1, Signed",
			ContentType: "application/vnd.docker.distribution.manifest.v1+prettyjws",
			JSON: `{
   "name": "hello-world",
   "tag": "latest",
   "architecture": "amd64",
   "fsLayers": [
      {
         "blobSum": "sha256:5f70bf18a086007016e948b04aed3b82103a36bea41755b6cddfaf10ace3c6ef"
      },
      {
         "blobSum": "sha256:5f70bf18a086007016e948b04aed3b82103a36bea41755b6cddfaf10ace3c6ef"
      },
      {
         "blobSum": "sha256:cc8567d70002e957612902a8e985ea129d831ebe04057d88fb644857caa45d11"
      },
      {
         "blobSum": "sha256:5f70bf18a086007016e948b04aed3b82103a36bea41755b6cddfaf10ace3c6ef"
      }
   ],
   "history": [
      {
         "v1Compatibility": "{\"id\":\"e45a5af57b00862e5ef5782a9925979a02ba2b12dff832fd0991335f4a11e5c5\",\"parent\":\"31cbccb51277105ba3ae35ce33c22b69c9e3f1002e76e4c736a2e8ebff9d7b5d\",\"created\":\"2014-12-31T22:57:59.178729048Z\",\"container\":\"27b45f8fb11795b52e9605b686159729b0d9ca92f76d40fb4f05a62e19c46b4f\",\"container_config\":{\"Hostname\":\"8ce6509d66e2\",\"Domainname\":\"\",\"User\":\"\",\"Memory\":0,\"MemorySwap\":0,\"CpuShares\":0,\"Cpuset\":\"\",\"AttachStdin\":false,\"AttachStdout\":false,\"AttachStderr\":false,\"PortSpecs\":null,\"ExposedPorts\":null,\"Tty\":false,\"OpenStdin\":false,\"StdinOnce\":false,\"Env\":[\"PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin\"],\"Cmd\":[\"/bin/sh\",\"-c\",\"#(nop) CMD [/hello]\"],\"Image\":\"31cbccb51277105ba3ae35ce33c22b69c9e3f1002e76e4c736a2e8ebff9d7b5d\",\"Volumes\":null,\"WorkingDir\":\"\",\"Entrypoint\":null,\"NetworkDisabled\":false,\"MacAddress\":\"\",\"OnBuild\":[],\"SecurityOpt\":null,\"Labels\":null},\"docker_version\":\"1.4.1\",\"config\":{\"Hostname\":\"8ce6509d66e2\",\"Domainname\":\"\",\"User\":\"\",\"Memory\":0,\"MemorySwap\":0,\"CpuShares\":0,\"Cpuset\":\"\",\"AttachStdin\":false,\"AttachStdout\":false,\"AttachStderr\":false,\"PortSpecs\":null,\"ExposedPorts\":null,\"Tty\":false,\"OpenStdin\":false,\"StdinOnce\":false,\"Env\":[\"PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin\"],\"Cmd\":[\"/hello\"],\"Image\":\"31cbccb51277105ba3ae35ce33c22b69c9e3f1002e76e4c736a2e8ebff9d7b5d\",\"Volumes\":null,\"WorkingDir\":\"\",\"Entrypoint\":null,\"NetworkDisabled\":false,\"MacAddress\":\"\",\"OnBuild\":[],\"SecurityOpt\":null,\"Labels\":null},\"architecture\":\"amd64\",\"os\":\"linux\",\"Size\":0}\n"
      },
      {
         "v1Compatibility": "{\"id\":\"e45a5af57b00862e5ef5782a9925979a02ba2b12dff832fd0991335f4a11e5c5\",\"parent\":\"31cbccb51277105ba3ae35ce33c22b69c9e3f1002e76e4c736a2e8ebff9d7b5d\",\"created\":\"2014-12-31T22:57:59.178729048Z\",\"container\":\"27b45f8fb11795b52e9605b686159729b0d9ca92f76d40fb4f05a62e19c46b4f\",\"container_config\":{\"Hostname\":\"8ce6509d66e2\",\"Domainname\":\"\",\"User\":\"\",\"Memory\":0,\"MemorySwap\":0,\"CpuShares\":0,\"Cpuset\":\"\",\"AttachStdin\":false,\"AttachStdout\":false,\"AttachStderr\":false,\"PortSpecs\":null,\"ExposedPorts\":null,\"Tty\":false,\"OpenStdin\":false,\"StdinOnce\":false,\"Env\":[\"PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin\"],\"Cmd\":[\"/bin/sh\",\"-c\",\"#(nop) CMD [/hello]\"],\"Image\":\"31cbccb51277105ba3ae35ce33c22b69c9e3f1002e76e4c736a2e8ebff9d7b5d\",\"Volumes\":null,\"WorkingDir\":\"\",\"Entrypoint\":null,\"NetworkDisabled\":false,\"MacAddress\":\"\",\"OnBuild\":[],\"SecurityOpt\":null,\"Labels\":null},\"docker_version\":\"1.4.1\",\"config\":{\"Hostname\":\"8ce6509d66e2\",\"Domainname\":\"\",\"User\":\"\",\"Memory\":0,\"MemorySwap\":0,\"CpuShares\":0,\"Cpuset\":\"\",\"AttachStdin\":false,\"AttachStdout\":false,\"AttachStderr\":false,\"PortSpecs\":null,\"ExposedPorts\":null,\"Tty\":false,\"OpenStdin\":false,\"StdinOnce\":false,\"Env\":[\"PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin\"],\"Cmd\":[\"/hello\"],\"Image\":\"31cbccb51277105ba3ae35ce33c22b69c9e3f1002e76e4c736a2e8ebff9d7b5d\",\"Volumes\":null,\"WorkingDir\":\"\",\"Entrypoint\":null,\"NetworkDisabled\":false,\"MacAddress\":\"\",\"OnBuild\":[],\"SecurityOpt\":null,\"Labels\":null},\"architecture\":\"amd64\",\"os\":\"linux\",\"Size\":0}\n"
      }
   ],
   "schemaVersion": 1,
   "signatures": [
      {
         "header": {
            "jwk": {
               "crv": "P-256",
               "kid": "OD6I:6DRK:JXEJ:KBM4:255X:NSAA:MUSF:E4VM:ZI6W:CUN2:L4Z6:LSF4",
               "kty": "EC",
               "x": "3gAwX48IQ5oaYQAYSxor6rYYc_6yjuLCjtQ9LUakg4A",
               "y": "t72ge6kIA1XOjqjVoEOiPPAURltJFBMGDSQvEGVB010"
            },
            "alg": "ES256"
         },
         "signature": "XREm0L8WNn27Ga_iE_vRnTxVMhhYY0Zst_FfkKopg6gWSoTOZTuW4rK0fg_IqnKkEKlbD83tD46LKEGi5aIVFg",
         "protected": "eyJmb3JtYXRMZW5ndGgiOjY2MjgsImZvcm1hdFRhaWwiOiJDbjAiLCJ0aW1lIjoiMjAxNS0wNC0wOFQxODo1Mjo1OVoifQ"
      }
   ]
}`,
			Expected: &ImageManifest{
				ContentType:   "application/vnd.docker.distribution.manifest.v1+prettyjws",
				SchemaVersion: 1,
				Digest:        "sha256:44d1afbd01feff689f7f24808298e2fba2647292c17cb6ef9f2041cafb8fe496",
				Platform: &Platform{
					Architecture: "amd64",
				},
				Annotations: Annotations{},
			},
		},
		{
			Name:        "Docker Image Manifest Version 2, Schema 2",
			ContentType: "application/vnd.docker.distribution.manifest.v2+json",
			JSON: `{
    "schemaVersion": 2,
    "mediaType": "application/vnd.docker.distribution.manifest.v2+json",
    "config": {
        "mediaType": "application/vnd.docker.container.image.v1+json",
        "digest": "sha256:b5b2b2c507a0944348e0303114d8d93aaaa081732b86451d9bce1f432a537bc7",
        "size": 7023
    },
    "layers": [
        {
            "mediaType": "application/vnd.docker.image.rootfs.diff.tar.gzip",
            "digest": "sha256:e692418e4cbaf90ca69d05a66403747baa33ee08806650b51fab815ad7fc331f",
            "size": 32654
        },
        {
            "mediaType": "application/vnd.docker.image.rootfs.diff.tar.gzip",
            "digest": "sha256:3c3a4604a545cdc127456d94e421cd355bca5b528f4a9c1905b15da2eb4a4c6b",
            "size": 16724
        },
        {
            "mediaType": "application/vnd.docker.image.rootfs.diff.tar.gzip",
            "digest": "sha256:ec4b8955958665577945c89419d1af06b5f7636b4ac3da7f12184802ad867736",
            "size": 73109
        }
    ]
}`,
			Expected: &ImageManifest{
				ContentType:   "application/vnd.docker.distribution.manifest.v2+json",
				SchemaVersion: 2,
				MediaType:     "application/vnd.docker.distribution.manifest.v2+json",
				Digest:        "sha256:db333c0ff517ff987cf018987d61f86b96e9d4144293d5b2dbcff24f2d409915",
				Annotations:   Annotations{},
			},
		},
		{
			Name:        "Docker Image Manifest List Version 2, Schema 2",
			ContentType: "application/vnd.docker.distribution.manifest.list.v2+json",
			JSON: `{
  "schemaVersion": 2,
  "mediaType": "application/vnd.docker.distribution.manifest.list.v2+json",
  "manifests": [
    {
      "mediaType": "application/vnd.docker.distribution.manifest.v2+json",
      "digest": "sha256:e692418e4cbaf90ca69d05a66403747baa33ee08806650b51fab815ad7fc331f",
      "size": 7143,
      "platform": {
        "architecture": "ppc64le",
        "os": "linux"
      }
    },
    {
      "mediaType": "application/vnd.docker.distribution.manifest.v2+json",
      "digest": "sha256:5b0bcabd1ed22e9fb1310cf6c2dec7cdef19f0ad69efa1f392e94a4333501270",
      "size": 7682,
      "platform": {
        "architecture": "amd64",
        "os": "linux",
        "features": [
          "sse4"
        ]
      }
    }
  ]
}`,
			Expected: &ImageIndex{
				ContentType:   "application/vnd.docker.distribution.manifest.list.v2+json",
				SchemaVersion: 2,
				MediaType:     "application/vnd.docker.distribution.manifest.list.v2+json",
				Digest:        "sha256:471d59459bb32757fc0dc4b1db84186a105f000def20f59e2f41599f8099e827",
				Annotations:   Annotations{},
				Manifests: []ImageManifest{
					{
						MediaType: "application/vnd.docker.distribution.manifest.v2+json",
						Digest:    "sha256:e692418e4cbaf90ca69d05a66403747baa33ee08806650b51fab815ad7fc331f",
						Platform: &Platform{
							Architecture: "ppc64le",
							OS:           "linux",
						},
						Annotations: Annotations{},
					},
					{
						MediaType: "application/vnd.docker.distribution.manifest.v2+json",
						Digest:    "sha256:5b0bcabd1ed22e9fb1310cf6c2dec7cdef19f0ad69efa1f392e94a4333501270",
						Platform: &Platform{
							Architecture: "amd64",
							OS:           "linux",
						},
						Annotations: Annotations{},
					},
				},
			},
		},
		{
			Name:        "OCI Image Index Example",
			ContentType: "application/vnd.oci.image.index.v1+json",
			JSON: `{
  "schemaVersion": 2,
  "mediaType": "application/vnd.oci.image.index.v1+json",
  "manifests": [
    {
      "mediaType": "application/vnd.oci.image.manifest.v1+json",
      "size": 7143,
      "digest": "sha256:e692418e4cbaf90ca69d05a66403747baa33ee08806650b51fab815ad7fc331f",
      "platform": {
        "architecture": "ppc64le",
        "os": "linux"
      }
    },
    {
      "mediaType": "application/vnd.oci.image.manifest.v1+json",
      "size": 7682,
      "digest": "sha256:5b0bcabd1ed22e9fb1310cf6c2dec7cdef19f0ad69efa1f392e94a4333501270",
      "platform": {
        "architecture": "amd64",
        "os": "linux"
      }
    }
  ],
  "annotations": {
    "com.example.key1": "value1",
    "com.example.key2": "value2"
  }
}`,
			Expected: &ImageIndex{
				ContentType:   "application/vnd.oci.image.index.v1+json",
				SchemaVersion: 2,
				MediaType:     "application/vnd.oci.image.index.v1+json",
				Digest:        "sha256:05bbc206ed9c460127a40c1fcdbccf21bf01fe892a80dfff2014f9cd26aebf4a",
				Manifests: []ImageManifest{
					{
						MediaType: "application/vnd.oci.image.manifest.v1+json",
						Digest:    "sha256:e692418e4cbaf90ca69d05a66403747baa33ee08806650b51fab815ad7fc331f",
						Platform: &Platform{
							Architecture: "ppc64le",
							OS:           "linux",
						},
						Annotations: Annotations{},
					},
					{
						MediaType: "application/vnd.oci.image.manifest.v1+json",
						Digest:    "sha256:5b0bcabd1ed22e9fb1310cf6c2dec7cdef19f0ad69efa1f392e94a4333501270",
						Platform: &Platform{
							Architecture: "amd64",
							OS:           "linux",
						},
						Annotations: Annotations{},
					},
				},
				Annotations: Annotations{
					"com.example.key1": "value1",
					"com.example.key2": "value2",
				},
			},
		},
		{
			ContentType: "application/vnd.oci.image.manifest.v1+json",
			Name:        "OCI Image Manifest Example",
			JSON: `{
  "schemaVersion": 2,
  "mediaType": "application/vnd.oci.image.manifest.v1+json",
  "config": {
    "mediaType": "application/vnd.oci.image.config.v1+json",
    "digest": "sha256:b5b2b2c507a0944348e0303114d8d93aaaa081732b86451d9bce1f432a537bc7",
    "size": 7023
  },
  "layers": [
    {
      "mediaType": "application/vnd.oci.image.layer.v1.tar+gzip",
      "digest": "sha256:9834876dcfb05cb167a5c24953eba58c4ac89b1adf57f28f2f9d09af107ee8f0",
      "size": 32654
    },
    {
      "mediaType": "application/vnd.oci.image.layer.v1.tar+gzip",
      "digest": "sha256:3c3a4604a545cdc127456d94e421cd355bca5b528f4a9c1905b15da2eb4a4c6b",
      "size": 16724
    },
    {
      "mediaType": "application/vnd.oci.image.layer.v1.tar+gzip",
      "digest": "sha256:ec4b8955958665577945c89419d1af06b5f7636b4ac3da7f12184802ad867736",
      "size": 73109
    }
  ],
  "subject": {
    "mediaType": "application/vnd.oci.image.manifest.v1+json",
    "digest": "sha256:5b0bcabd1ed22e9fb1310cf6c2dec7cdef19f0ad69efa1f392e94a4333501270",
    "size": 7682
  },
  "annotations": {
    "com.example.key1": "value1",
    "com.example.key2": "value2"
  }
}`,
			Expected: &ImageManifest{
				ContentType:   "application/vnd.oci.image.manifest.v1+json",
				SchemaVersion: 2,
				MediaType:     "application/vnd.oci.image.manifest.v1+json",
				Digest:        "sha256:d701f04d63badebf6532354e76d8d36fdba8fed2c5f1ab538cb1847baab3de22",
				Annotations: Annotations{
					"com.example.key1": "value1",
					"com.example.key2": "value2",
				},
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			blob := newBlobResponse(
				io.NopCloser(strings.NewReader(testCase.JSON)),
				BlobInfo{
					ContentType: testCase.ContentType,
				},
			)

			actual, err := manifestFromBlob(blob)
			assert.Equal(t, testCase.Expected, actual)
			if testCase.ExpectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestManifestsMaybeEqual(t *testing.T) {
	testCases := []struct {
		Name       string
		ManifestA  any
		ManifestB  any
		Platform   *Platform
		MaybeEqual bool
	}{
		{
			Name: "Identical image manifests",
			ManifestA: &ImageManifest{
				Digest: "01ba4719c80b6fe911b091a7c05124b64eeece964e09c058ef8f9805daca546b",
			},
			ManifestB: &ImageManifest{
				Digest: "01ba4719c80b6fe911b091a7c05124b64eeece964e09c058ef8f9805daca546b",
			},
			Platform:   nil,
			MaybeEqual: true,
		},
		{
			Name: "Different image manifests",
			ManifestA: &ImageManifest{
				Digest: "01ba4719c80b6fe911b091a7c05124b64eeece964e09c058ef8f9805daca546b",
			},
			ManifestB: &ImageManifest{
				Digest: "4ed993c7c28d18538eb6e42e6727f637e0021108f8901e616511a87671400468",
			},
			Platform:   nil,
			MaybeEqual: false,
		},
		{
			Name: "Image manifest found in index",
			ManifestA: &ImageManifest{
				Digest: "01ba4719c80b6fe911b091a7c05124b64eeece964e09c058ef8f9805daca546b",
			},
			ManifestB: &ImageIndex{
				Digest: "5564ee5fbe988d6ff476b7b037343c46a114e2638959e49654cd876c60ac6661",
				Manifests: []ImageManifest{
					{
						Digest: "01ba4719c80b6fe911b091a7c05124b64eeece964e09c058ef8f9805daca546b",
					},
				},
			},
			Platform:   nil,
			MaybeEqual: true,
		},
		{
			Name: "Image index contains in manifest",
			ManifestA: &ImageIndex{
				Digest: "5564ee5fbe988d6ff476b7b037343c46a114e2638959e49654cd876c60ac6661",
				Manifests: []ImageManifest{
					{
						Digest: "01ba4719c80b6fe911b091a7c05124b64eeece964e09c058ef8f9805daca546b",
					},
				},
			},
			ManifestB: &ImageManifest{
				Digest: "01ba4719c80b6fe911b091a7c05124b64eeece964e09c058ef8f9805daca546b",
			},
			Platform:   nil,
			MaybeEqual: true,
		},
		{
			Name: "Image indexes equal",
			ManifestA: &ImageIndex{
				Digest: "5564ee5fbe988d6ff476b7b037343c46a114e2638959e49654cd876c60ac6661",
			},
			ManifestB: &ImageIndex{
				Digest: "5564ee5fbe988d6ff476b7b037343c46a114e2638959e49654cd876c60ac6661",
			},
			Platform:   nil,
			MaybeEqual: true,
		},
		{
			Name: "Different image indexes contain identical images",
			ManifestA: &ImageIndex{
				Digest: "5564ee5fbe988d6ff476b7b037343c46a114e2638959e49654cd876c60ac6661",
				Manifests: []ImageManifest{
					{
						Digest: "2cbffe11c9853bdb44577d480183e8f1ede76cc9bfb01168b9838a626cbb6026",
					},
					{
						Digest: "68734cfc63ae7b386b47571ef3f4f5ef1395a73823db56c54919e0adb8dd2c2a",
					},
				},
			},
			ManifestB: &ImageIndex{
				Digest: "bd179a032d9b85aa3577e155a536ab388904bacbb905cf7b1a946360e8ce565c",
				Manifests: []ImageManifest{
					{
						Digest: "2cbffe11c9853bdb44577d480183e8f1ede76cc9bfb01168b9838a626cbb6026",
					},
					{
						Digest: "68734cfc63ae7b386b47571ef3f4f5ef1395a73823db56c54919e0adb8dd2c2a",
					},
				},
			},
			Platform:   nil,
			MaybeEqual: true,
		},
		{
			Name: "Different image indexes contain different images",
			ManifestA: &ImageIndex{
				Digest: "5564ee5fbe988d6ff476b7b037343c46a114e2638959e49654cd876c60ac6661",
				Manifests: []ImageManifest{
					{
						Digest: "2cbffe11c9853bdb44577d480183e8f1ede76cc9bfb01168b9838a626cbb6026",
					},
					{
						Digest: "68734cfc63ae7b386b47571ef3f4f5ef1395a73823db56c54919e0adb8dd2c2a",
					},
				},
			},
			ManifestB: &ImageIndex{
				Digest: "bd179a032d9b85aa3577e155a536ab388904bacbb905cf7b1a946360e8ce565c",
				Manifests: []ImageManifest{
					{
						Digest: "4ee5e80bb16a134ef95c50314b22c295068507ee1dd8bf5bb4e1dc8c5448cc0c",
					},
					{
						Digest: "d628021d8d8c29f934568f844fc02a7602e5fbd690e45030be1552b096583000",
					},
				},
			},
			Platform:   nil,
			MaybeEqual: false,
		},
		{
			Name: "Different image indexes contain identical images for platform",
			ManifestA: &ImageIndex{
				Digest: "5564ee5fbe988d6ff476b7b037343c46a114e2638959e49654cd876c60ac6661",
				Manifests: []ImageManifest{
					{
						Digest: "2cbffe11c9853bdb44577d480183e8f1ede76cc9bfb01168b9838a626cbb6026",
						Platform: &Platform{
							Architecture: "arm64",
							OS:           "linux",
							Variant:      "v8",
						},
					},
					{
						// Only different for the amd64 platform, but checking for arm
						Digest: "68734cfc63ae7b386b47571ef3f4f5ef1395a73823db56c54919e0adb8dd2c2a",
						Platform: &Platform{
							Architecture: "amd64",
							OS:           "linux",
						},
					},
				},
			},
			ManifestB: &ImageIndex{
				Digest: "bd179a032d9b85aa3577e155a536ab388904bacbb905cf7b1a946360e8ce565c",
				Manifests: []ImageManifest{
					{
						Digest: "2cbffe11c9853bdb44577d480183e8f1ede76cc9bfb01168b9838a626cbb6026",
						Platform: &Platform{
							Architecture: "arm64",
							OS:           "linux",
							Variant:      "v8",
						},
					},
					{
						// Only different for the amd64 platform, but checking for arm
						Digest: "d628021d8d8c29f934568f844fc02a7602e5fbd690e45030be1552b096583000",
						Platform: &Platform{
							Architecture: "amd64",
							OS:           "linux",
						},
					},
				},
			},
			Platform: &Platform{
				Architecture: "arm64",
				OS:           "linux",
				Variant:      "v8",
			},
			MaybeEqual: true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			actual := ManifestsMaybeEqual(testCase.ManifestA, testCase.ManifestB, testCase.Platform)
			assert.Equal(t, testCase.MaybeEqual, actual)
		})
	}
}

func TestFilterManifestsByPlatform(t *testing.T) {
	manifests := []ImageManifest{
		{
			Digest:   "1",
			Platform: nil,
		},
		{
			Digest:   "2",
			Platform: &Platform{},
		},
		{
			Digest: "3",
			Platform: &Platform{
				Architecture: "amd64",
			},
		},
		{
			Digest: "4",
			Platform: &Platform{
				Architecture: "amd64",
				OS:           "linux",
			},
		},
		{
			Digest: "5",
			Platform: &Platform{
				Architecture: "unknown",
				OS:           "unknown",
			},
		},
		{
			Digest: "6",
			Platform: &Platform{
				Architecture: "arm64",
			},
		},
		{
			Digest: "7",
			Platform: &Platform{
				Architecture: "arm64",
				OS:           "linux",
			},
		},
		{
			Digest: "8",
			Platform: &Platform{
				Architecture: "arm64",
				OS:           "linux",
				Variant:      "v8",
			},
		},
	}

	testCases := []struct {
		Name           string
		Platform       *Platform
		MatchedDigests []string
	}{
		{
			Name:     "No platform specified",
			Platform: nil,
			MatchedDigests: []string{
				"1", "2", "3", "4", "5", "6", "7", "8",
			},
		},
		{
			Name:     "Empty values specified",
			Platform: &Platform{},
			MatchedDigests: []string{
				"2", "3", "4", "5", "6", "7", "8",
			},
		},
		{
			Name: "Architecture specified (arm64)",
			Platform: &Platform{
				Architecture: "arm64",
			},
			MatchedDigests: []string{
				"6", "7", "8",
			},
		},
		{
			Name: "Architecture specified (amd64)",
			Platform: &Platform{
				Architecture: "amd64",
			},
			MatchedDigests: []string{
				"3", "4",
			},
		},
		{
			Name: "OS specified",
			Platform: &Platform{
				OS: "linux",
			},
			MatchedDigests: []string{
				"4", "7", "8",
			},
		},
		{
			Name: "Everything specified",
			Platform: &Platform{
				Architecture: "arm64",
				OS:           "linux",
				Variant:      "v8",
			},
			MatchedDigests: []string{
				"8",
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			matchedDigests := make([]string, 0)
			for _, matched := range filterManifestsByPlatform(manifests, testCase.Platform) {
				matchedDigests = append(matchedDigests, matched.Digest)
			}

			assert.Equal(t, testCase.MatchedDigests, matchedDigests)
		})
	}
}
