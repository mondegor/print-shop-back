
mrcmd_func_openapi_build_auth() {
  local sectionDir="${1:?}" # sample: .../admin-api, .../public-api
  local sharedDir="${2:?}" # sample: .../_shared

  local authDir="${sectionDir}/auth"
  local checkDir="${sectionDir}/check"
  local operationDir="${sectionDir}/operation"
  local securityDir="${sectionDir}/security"
  local sessionDir="${sectionDir}/session"

  # OPENAPI_VERSION="3.0.3"

  OPENAPI_HEADERS=(
    "${sectionDir}/header.yaml"
    "${sharedDir}/description-errors.md"
  )

  OPENAPI_SERVERS=(
    "${sectionDir}/servers.yaml"
  )

  OPENAPI_TAGS=(
    "${authDir}/tags.yaml"
    "${checkDir}/tags.yaml"
    "${operationDir}/tags.yaml"
    "${securityDir}/tags.yaml"
    "${sessionDir}/tags.yaml"
  )

  OPENAPI_PATHS=(
    "${authDir}/auth_paths.yaml"
    "${checkDir}/check_paths.yaml"
    "${operationDir}/operation_paths.yaml"
    "${securityDir}/security_paths.yaml"
    "${sessionDir}/session_paths.yaml"
  )

  OPENAPI_COMPONENTS_HEADERS=(
    # "${sharedDir}/components/headers/"

    "${authDir}/auth_headers.yaml"
  )

  OPENAPI_COMPONENTS_PARAMETERS=(
    "${sharedDir}/components/parameters/App.Request.Header.AcceptLanguage.yaml"
    "${sharedDir}/components/parameters/App.Request.Header.CorrelationID.yaml"
    # "${sharedDir}/components/parameters/App.Request.Header.IdempotencyKey.yaml"
    "${sharedDir}/components/parameters/App.Request.Query.Language.yaml"
    # "${sharedDir}/components/parameters/App.Request.Query.Filter.SearchText.yaml"
    # "${sharedDir}/components/parameters/App.Request.Query.Filter.Statuses.yaml"
    # "${sharedDir}/components/parameters/App.Request.Query.ListCursor.yaml"
    # "${sharedDir}/components/parameters/App.Request.Query.ListPager.yaml"
    # "${sharedDir}/components/parameters/App.Request.Query.ListSorter.yaml"

    "${authDir}/auth_parameters.yaml"
  )

  OPENAPI_COMPONENTS_SCHEMAS=(
    # "${sharedDir}/components/enums/App.Enum.Address.HouseType.yaml"
    # "${sharedDir}/components/enums/App.Enum.DeliveryMethod.yaml"
    # "${sharedDir}/components/enums/App.Enum.Gender.yaml"
    # "${sharedDir}/components/enums/App.Enum.Status.yaml"

    # "${sharedDir}/components/fields/App.Field.Article.yaml"
    # "${sharedDir}/components/fields/App.Field.Boolean.yaml"
    # "${sharedDir}/components/fields/App.Field.Caption.yaml"
    "${sharedDir}/components/fields/App.Field.DateTimeCreatedAt.yaml"
    "${sharedDir}/components/fields/App.Field.DateTimeUpdatedAt.yaml"
    # "${sharedDir}/components/fields/App.Field.Date.yaml"
    # "${sharedDir}/components/fields/App.Field.DateTime.yaml"
    # "${sharedDir}/components/fields/App.Field.DoubleSize.yaml"
    # "${sharedDir}/components/fields/App.Field.Email.yaml"
    # "${sharedDir}/components/fields/App.Field.ExternalURL.yaml"
    # "${sharedDir}/components/fields/App.Field.FileURL.yaml"
    # "${sharedDir}/components/fields/App.Field.Float64.yaml"
    # "${sharedDir}/components/fields/App.Field.GEO.yaml"
    # "${sharedDir}/components/fields/App.Field.ImageURL.yaml"
    # "${sharedDir}/components/fields/App.Field.Int16.yaml"
    # "${sharedDir}/components/fields/App.Field.Int32.yaml"
    # "${sharedDir}/components/fields/App.Field.Int64.yaml"
    # "${sharedDir}/components/fields/App.Field.JsonData.yaml"
    # "${sharedDir}/components/fields/App.Field.ListCursor.Cursor.yaml"
    # "${sharedDir}/components/fields/App.Field.ListCursor.HasNext.yaml"
    # "${sharedDir}/components/fields/App.Field.ListPager.Total.yaml"
    # "${sharedDir}/components/fields/App.Field.OrderIndex.yaml"
    # "${sharedDir}/components/fields/App.Field.Percent.yaml"
    # "${sharedDir}/components/fields/App.Field.Phone.yaml"
    # "${sharedDir}/components/fields/App.Field.RewriteName.yaml"
    # "${sharedDir}/components/fields/App.Field.TagVersion.yaml"
    # "${sharedDir}/components/fields/App.Field.Timezone.yaml"
    # "${sharedDir}/components/fields/App.Field.Uint.yaml"
    # "${sharedDir}/components/fields/App.Field.UUID.yaml"
    # "${sharedDir}/components/fields/App.Field.VariableCamelCase.yaml"

    # "${sharedDir}/components/fields/measures/App.Field.Measure.Centimeter.yaml"
    # "${sharedDir}/components/fields/measures/App.Field.Measure.DoubleMillimeter.yaml"
    # "${sharedDir}/components/fields/measures/App.Field.Measure.Gram.yaml"
    # "${sharedDir}/components/fields/measures/App.Field.Measure.GramPerMeter2.yaml"
    # "${sharedDir}/components/fields/measures/App.Field.Measure.Kilogram.yaml"
    # "${sharedDir}/components/fields/measures/App.Field.Measure.KilogramPerMeter2.yaml"
    # "${sharedDir}/components/fields/measures/App.Field.Measure.Meter.yaml"
    # "${sharedDir}/components/fields/measures/App.Field.Measure.Meter2.yaml"
    # "${sharedDir}/components/fields/measures/App.Field.Measure.Meter3.yaml"
    # "${sharedDir}/components/fields/measures/App.Field.Measure.Micrometer.yaml"
    # "${sharedDir}/components/fields/measures/App.Field.Measure.Milligram.yaml"
    # "${sharedDir}/components/fields/measures/App.Field.Measure.Millimeter.yaml"
    # "${sharedDir}/components/fields/measures/App.Field.Measure.Price.yaml"
    # "${sharedDir}/components/fields/measures/App.Field.Measure.TripleMillimeter.yaml"

    # "${sharedDir}/components/models/App.Request.Model.ChangeFlag.yaml"
    # "${sharedDir}/components/models/App.Request.Model.ChangeStatus.yaml"
    # "${sharedDir}/components/models/App.Request.Model.MoveItem.yaml"
    # "${sharedDir}/components/models/App.Response.Model.BinaryAnyFile.yaml"
    # "${sharedDir}/components/models/App.Response.Model.BinaryFile.yaml"
    # "${sharedDir}/components/models/App.Response.Model.BinaryImage.yaml"
    # "${sharedDir}/components/models/App.Response.Model.BinaryMedia.yaml"
    "${sharedDir}/components/models/App.Response.Model.Error.yaml"
    # "${sharedDir}/components/models/App.Response.Model.FileInfo.yaml"
    # "${sharedDir}/components/models/App.Response.Model.JsonFile.yaml"
    # "${sharedDir}/components/models/App.Response.Model.ImageInfo.yaml"
    # "${sharedDir}/components/models/App.Response.Model.Success.yaml"
    # "${sharedDir}/components/models/App.Response.Model.SuccessCreatedItem.yaml"
    # "${sharedDir}/components/models/App.Response.Model.SuccessCreatedItemUint.yaml"
    # "${sharedDir}/components/models/App.Response.Model.SuccessSavedItem.yaml"
    # "${sharedDir}/components/models/App.Response.Model.TextFile.yaml"
    # "${sharedDir}/components/models/App.Response.Model.Volume.yaml"

    "${sharedDir}/custom/enums/Custom.Enum.UserAuth2fa.yaml"
    "${sharedDir}/custom/enums/Custom.Enum.UserStatus.yaml"

    "${sectionDir}/enums/Auth.Enum.ConfirmMethod.yaml"
    "${sectionDir}/enums/Auth.Enum.PasswordStrength.yaml"

    "${authDir}/auth_schemas.yaml"
    "${checkDir}/check_schemas.yaml"
    "${operationDir}/operation_schemas.yaml"
    "${securityDir}/security_schemas.yaml"
    "${sessionDir}/session_schemas.yaml"
  )

  OPENAPI_COMPONENTS_RESPONSES=(
    "${sharedDir}/components/responses/App.ResponseJson.Errors.yaml"
    "${sharedDir}/components/responses/App.ResponseJson.ErrorsAuth.yaml"

    "${authDir}/auth_responses.yaml"
  )

  OPENAPI_SECURITY_SCHEMES=(
    "${sharedDir}/securitySchemes.yaml"
  )
}
