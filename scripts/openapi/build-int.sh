
mrcmd_func_openapi_build_int() {
  local sectionDir="${1:?}" # sample: .../admin-api, .../public-api
  local sharedDir="${2:?}" # sample: .../_shared

  local systemDir="${sectionDir}/system"

  # OPENAPI_VERSION="3.0.3"

  OPENAPI_HEADERS=(
    "${sectionDir}/header.yaml"
  )

  OPENAPI_SERVERS=(
    "${sectionDir}/servers.yaml"
  )

  OPENAPI_TAGS=(
    "${systemDir}/tags.yaml"
  )

  OPENAPI_PATHS=(
    "${systemDir}/paths.yaml"
  )

#  OPENAPI_COMPONENTS_HEADERS=(
#    "${sharedDir}/components/headers/"
#  )

  OPENAPI_COMPONENTS_PARAMETERS=(
    "${sharedDir}/components/parameters/App.Request.Header.AcceptLanguage.yaml"
    "${sharedDir}/components/parameters/App.Request.Header.CorrelationID.yaml"
  )

  OPENAPI_COMPONENTS_SCHEMAS=(
    "${sharedDir}/components/models/App.Response.Model.Error.yaml"

    "${systemDir}/schemas.yaml"
  )

  OPENAPI_COMPONENTS_RESPONSES=(
    "${sharedDir}/components/responses/App.ResponseJson.Errors.yaml"
    # "${sharedDir}/components/responses/App.ResponseJson.ErrorsAuth.yaml"
  )

#  OPENAPI_SECURITY_SCHEMES=(
#    "${sharedDir}/securitySchemes.yaml"
#  )
}
