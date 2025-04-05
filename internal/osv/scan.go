package osv

import (
	"context"
	"errors"
	"os"

	"github.com/google/osv-scanner/v2/pkg/osvscanner"
)

// ScanSPDX scans a SBOM of the SPDX format.
func ScanSPDX(ctx context.Context, sbom string) ([]Vulnerability, error) {
	file, err := os.CreateTemp("", "cupdate-scan-sbom-*.spdx.json")
	if err != nil {
		return nil, err
	}

	_, err = file.WriteString(sbom)
	file.Close()
	if err != nil {
		os.Remove(file.Name())
		return nil, err
	}

	actions := osvscanner.ScannerActions{
		SBOMPaths: []string{file.Name()},
	}

	results, err := osvscanner.DoScan(actions)
	os.Remove(file.Name())
	if err != nil && !errors.Is(err, osvscanner.ErrVulnerabilitiesFound) {
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

				vulnerability := Vulnerability{
					ID:               vuln.ID,
					Modified:         vuln.Modified,
					DatabaseSpecific: vuln.DatabaseSpecific,
					Details:          vuln.Details,
					Related:          vuln.Related,
					SchemaVersion:    vuln.SchemaVersion,
					Summary:          vuln.Summary,
				}

				if len(vuln.Affected) > 0 {
					vulnerability.Affected = make([]Affected, 0)
					for _, a := range vuln.Affected {
						res := Affected{
							DatabaseSpecific:  a.DatabaseSpecific,
							EcosystemSpecific: a.EcosystemSpecific,
							Versions:          a.Versions,
						}

						if a.Package.Name != "" {
							res.Package = &AffectedPackage{
								Ecosystem: a.Package.Ecosystem,
								Name:      a.Package.Name,
								Purl:      a.Package.Purl,
							}
						}

						if len(a.Ranges) > 0 {
							res.Ranges = make([]AffectedRange, 0)
							for _, r := range a.Ranges {
								affectedRange := AffectedRange{
									DatabaseSpecific: r.DatabaseSpecific,
									Events:           make([]Event, 0),
									Repo:             r.Repo,
									Type:             string(r.Type),
								}

								for _, event := range r.Events {
									affectedRange.Events = append(affectedRange.Events, Event{
										Introduced:   event.Introduced,
										Fixed:        event.Fixed,
										LastAffected: event.LastAffected,
										Limit:        event.Limit,
									})
								}

								res.Ranges = append(res.Ranges, affectedRange)
							}
						}

						if len(a.Severity) > 0 {
							res.Severities = make([]Severity, 0)
							for _, severity := range a.Severity {
								res.Severities = append(res.Severities, Severity{
									Type:  string(severity.Type),
									Score: severity.Score,
								})
							}
						}
						vulnerability.Affected = append(vulnerability.Affected, res)
					}
				}

				if len(vuln.Credits) > 0 {
					vulnerability.Credits = make([]Credit, 0)
					for _, credit := range vuln.Credits {
						vulnerability.Credits = append(vulnerability.Credits, Credit{
							Contact: credit.Contact,
							Name:    credit.Name,
							Type:    string(credit.Type),
						})
					}
				}

				if !vuln.Published.IsZero() {
					vulnerability.Published = &vuln.Published
				}

				if len(vuln.References) > 0 {
					vulnerability.References = make([]Reference, 0)
					for _, ref := range vulnerability.References {
						vulnerability.References = append(vulnerability.References, Reference{
							Type: ref.Type,
							URL:  ref.URL,
						})
					}
				}

				if len(vuln.Severity) > 0 {
					vulnerability.Severities = make([]Severity, 0)
					for _, severity := range vuln.Severity {
						vulnerability.Severities = append(vulnerability.Severities, Severity{
							Score: severity.Score,
							Type:  string(severity.Type),
						})
					}
				}

				if !vuln.Withdrawn.IsZero() {
					vulnerability.Withdrawn = &vuln.Withdrawn
				}

				vulnerabilities = append(vulnerabilities, vulnerability)
			}
		}
	}

	return vulnerabilities, nil
}
