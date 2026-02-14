package osv

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"os"
	"os/exec"
)

var ErrScannerNotFound = errors.New("osv: osv-scanner executable was not found")

// ScanSPDX scans a SBOM of the SPDX format.
func ScanSPDX(ctx context.Context, sbom string) ([]Vulnerability, error) {
	file, err := os.CreateTemp("", "cupdate-scan-sbom-*.spdx.json")
	if err != nil {
		return nil, err
	}
	defer os.Remove(file.Name())

	_, err = file.WriteString(sbom)
	if err != nil {
		file.Close()
		return nil, err
	}

	if err := file.Close(); err != nil {
		return nil, err
	}

	cmd := exec.CommandContext(ctx, "osv-scanner", "scan", "--verbosity", "error", "--format", "json", "--lockfile", file.Name())

	var buffer bytes.Buffer
	cmd.Stderr = &buffer
	cmd.Stdout = &buffer

	err = cmd.Run()
	if _, ok := errors.AsType[*exec.ExitError](err); ok {
		// Ignore - if vulnerabilities are found the exit code is 1 but there is
		// still output
	} else if _, ok := errors.AsType[*exec.Error](err); ok {
		return nil, ErrScannerNotFound
	} else if err != nil {
		return nil, err
	}

	var results struct {
		Results []struct {
			Packages []struct {
				Vulnerabilities []Vulnerability `json:"vulnerabilities"`
			} `json:"packages"`
		} `json:"results"`
	}
	if err := json.NewDecoder(&buffer).Decode(&results); err != nil {
		// TODO: Likely an execution error, not a decode error
		return nil, err
	}

	vulnerabilities := make([]Vulnerability, 0)
	seen := make(map[string]struct{})
	for _, result := range results.Results {
		for _, pkg := range result.Packages {
			for _, vuln := range pkg.Vulnerabilities {
				// Assume same information available in each instance
				if _, ok := seen[vuln.ID]; ok {
					continue
				}

				seen[vuln.ID] = struct{}{}

				vulnerabilities = append(vulnerabilities, vuln)
			}
		}
	}

	return vulnerabilities, nil
}
