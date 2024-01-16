
mrcmd_func_openapi_build_system() {
  local sectionDir="${1:?}" # sample: .../admin-api, .../public-api
  local sharedDir="${2:?}" # sample: .../_shared

  # OPENAPI_VERSION="3.0.3"

  OPENAPI_HEADERS=(
    "${sectionDir}/header.yaml"
  )

  OPENAPI_SERVERS=(
    "${sectionDir}/servers.yaml"
  )

  OPENAPI_TAGS=(
    "${sectionDir}/tags.yaml"
  )

  OPENAPI_PATHS=(
    "${sectionDir}/paths.yaml"
  )

#  OPENAPI_COMPONENTS_HEADERS=(
#    "${sharedDir}/components/headers/"
#  )

#  OPENAPI_COMPONENTS_PARAMETERS=(
#    "${sectionDir}/components-parameters.yaml"
#  )

  OPENAPI_COMPONENTS_SCHEMAS=(
    "${sharedDir}/components/schemas/App.Response.Model.Error.yaml"

    "${sectionDir}/components-schemas.yaml"
  )

  OPENAPI_COMPONENTS_RESPONSES=(
    "${sharedDir}/components/responses/App.ResponseJson.Errors.yaml"

    # "${sectionDir}/components-responses.yaml"
  )

#  OPENAPI_SECURITY_SCHEMES=(
#    "${sectionDir}/securitySchemes.yaml"
#  )
}
