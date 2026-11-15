/**
 * @see https://spdx.github.io/spdx-spec/v2.3.1/
 * @see https://github.com/spdx/spdx-spec/tree/develop/schemas
 */
export interface Document {
  packages: Package[]
  relationships: Relationship[]
}

export interface Package {
  SPDXID: string
  name: string
  description?: string
  versionInfo?: string
  externalRefs?: ExternalRef[]
}

export interface ExternalRef {
  referenceCategory: 'OTHER' | 'PERSISTENT-ID' | 'SECURITY' | 'PACKAGE-MANAGER'
  referenceLocator: string
  referenceType: string
}

export interface Relationship {
  spdxElementId: string
  relatedSpdxElement: string
  relationshipType:
    | 'VARIANT_OF'
    | 'COPY_OF'
    | 'PATCH_FOR'
    | 'TEST_DEPENDENCY_OF'
    | 'CONTAINED_BY'
    | 'DATA_FILE_OF'
    | 'OPTIONAL_COMPONENT_OF'
    | 'ANCESTOR_OF'
    | 'GENERATES'
    | 'CONTAINS'
    | 'OPTIONAL_DEPENDENCY_OF'
    | 'FILE_ADDED'
    | 'REQUIREMENT_DESCRIPTION_FOR'
    | 'DEV_DEPENDENCY_OF'
    | 'DEPENDENCY_OF'
    | 'BUILD_DEPENDENCY_OF'
    | 'DESCRIBES'
    | 'PREREQUISITE_FOR'
    | 'HAS_PREREQUISITE'
    | 'PROVIDED_DEPENDENCY_OF'
    | 'DYNAMIC_LINK'
    | 'DESCRIBED_BY'
    | 'METAFILE_OF'
    | 'DEPENDENCY_MANIFEST_OF'
    | 'PATCH_APPLIED'
    | 'RUNTIME_DEPENDENCY_OF'
    | 'TEST_OF'
    | 'TEST_TOOL_OF'
    | 'DEPENDS_ON'
    | 'SPECIFICATION_FOR'
    | 'FILE_MODIFIED'
    | 'DISTRIBUTION_ARTIFACT'
    | 'AMENDS'
    | 'DOCUMENTATION_OF'
    | 'GENERATED_FROM'
    | 'STATIC_LINK'
    | 'OTHER'
    | 'BUILD_TOOL_OF'
    | 'TEST_CASE_OF'
    | 'PACKAGE_OF'
    | 'DESCENDANT_OF'
    | 'FILE_DELETED'
    | 'EXPANDED_FROM_ARCHIVE'
    | 'DEV_TOOL_OF'
    | 'EXAMPLE_OF'
}
