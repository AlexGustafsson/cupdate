package oci

import (
	"encoding/json"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSBOMAttestationUnmarshalJSON(t *testing.T) {
	fs, err := os.OpenRoot("./testdata/attestations")
	require.NoError(t, err)

	testCases := []struct {
		// Path is the path within fs
		Path     string
		Expected *SBOMAttestation
	}{
		{
			Path: "sbom.json",
			Expected: &SBOMAttestation{
				Type: SBOMTypeSPDX,
				SBOM: `{
  "spdxVersion": "SPDX-2.3",
  "dataLicense": "CC0-1.0",
  "SPDXID": "SPDXRef-DOCUMENT",
  "name": "sbom",
  "documentNamespace": "https://anchore.com/syft/dir/sbom-38619726-4e9b-4a1e-b50e-c4d649f2368b",
  "creationInfo": {
    "licenseListVersion": "3.25",
    "creators": [
      "Organization: Anchore, Inc",
      "Tool: syft-v1.18.1",
      "Tool: buildkit-v0.20.2"
    ],
    "created": "2025-03-28T09:38:12Z"
  },
  "packages": [
    {
      "name": "github.com/AlexGustafsson/cupdate",
      "SPDXID": "SPDXRef-Package-go-module-github.com-AlexGustafsson-cupdate-381730257ca50558",
      "versionInfo": "v0.19.0-30-g509402d",
      "supplier": "NOASSERTION",
      "downloadLocation": "NOASSERTION",
      "filesAnalyzed": false,
      "sourceInfo": "acquired package info from go module information: /cupdate",
      "licenseConcluded": "NOASSERTION",
      "licenseDeclared": "NOASSERTION",
      "copyrightText": "NOASSERTION",
      "externalRefs": [
        {
          "referenceCategory": "SECURITY",
          "referenceType": "cpe23Type",
          "referenceLocator": "cpe:2.3:a:AlexGustafsson:cupdate:v0.19.0-30-g509402d:*:*:*:*:*:*:*"
        },
        {
          "referenceCategory": "PACKAGE-MANAGER",
          "referenceType": "purl",
          "referenceLocator": "pkg:golang/github.com/AlexGustafsson/cupdate@v0.19.0-30-g509402d"
        }
      ]
    },
    {
      "name": "github.com/caarlos0/env/v11",
      "SPDXID": "SPDXRef-Package-go-module-github.com-caarlos0-env-v11-db608271955ec031",
      "versionInfo": "v11.3.1",
      "supplier": "NOASSERTION",
      "downloadLocation": "NOASSERTION",
      "filesAnalyzed": false,
      "checksums": [
        {
          "algorithm": "SHA256",
          "checksumValue": "700acf582d79856984b7e81693b6018bb9445d35c2be96927428991365f99820"
        }
      ],
      "sourceInfo": "acquired package info from go module information: /cupdate",
      "licenseConcluded": "NOASSERTION",
      "licenseDeclared": "NOASSERTION",
      "copyrightText": "NOASSERTION",
      "externalRefs": [
        {
          "referenceCategory": "SECURITY",
          "referenceType": "cpe23Type",
          "referenceLocator": "cpe:2.3:a:caarlos0:env\\/v11:v11.3.1:*:*:*:*:*:*:*"
        },
        {
          "referenceCategory": "PACKAGE-MANAGER",
          "referenceType": "purl",
          "referenceLocator": "pkg:golang/github.com/caarlos0/env@v11.3.1#v11"
        }
      ]
    }
  ],
  "files": [
    {
      "fileName": "cupdate",
      "SPDXID": "SPDXRef-File-cupdate-ea4a2edd76fc063f",
      "fileTypes": [
        "APPLICATION",
        "BINARY"
      ],
      "checksums": [
        {
          "algorithm": "SHA256",
          "checksumValue": "42241ef392f519bd9869e78ff9b5c65a7420e77a89634d8a33f530e813456fe5"
        }
      ],
      "licenseConcluded": "NOASSERTION",
      "licenseInfoInFiles": [
        "NOASSERTION"
      ],
      "copyrightText": "NOASSERTION",
      "comment": "layerID: sha256:4cfbe772b5f4a27f4ff27e2980fcbcff34fd4e9f0f532012b411f490d7ac1725"
    }
  ],
  "relationships": [
    {
      "spdxElementId": "SPDXRef-Package-go-module-google.golang.org-protobuf-03fb8672af6ad95f",
      "relatedSpdxElement": "SPDXRef-Package-go-module-github.com-AlexGustafsson-cupdate-381730257ca50558",
      "relationshipType": "DEPENDENCY_OF"
    }
  ]
}`,
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Path, func(t *testing.T) {
			file, err := fs.Open(testCase.Path)
			require.NoError(t, err)
			defer file.Close()

			content, err := io.ReadAll(file)
			require.NoError(t, err)

			var attestation SBOMAttestation
			require.NoError(t, json.Unmarshal(content, &attestation))

			assert.Equal(t, testCase.Expected.SBOM, attestation.SBOM)
		})
	}
}
