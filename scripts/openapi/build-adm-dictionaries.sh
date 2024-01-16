
mrcmd_func_openapi_build_adm_dictionaries() {
  local sectionDir="${1:?}" # sample: .../admin-api, .../public-api
  local sharedDir="${2:?}" # sample: .../_shared
  local moduleDir="${sectionDir}/dictionaries"

  local laminateTypesDir="${moduleDir}/laminate-types"
  local paperColorsDir="${moduleDir}/paper-colors"
  local paperFacturesDir="${moduleDir}/paper-factures"
  local printFormatsDir="${moduleDir}/print-formats"

  # OPENAPI_VERSION="3.0.3"

  OPENAPI_HEADERS=(
    "${moduleDir}/header.yaml"
    "${sharedDir}/description-errors.md"
  )

  OPENAPI_SERVERS=(
    "${sectionDir}/servers.yaml"
  )

  OPENAPI_TAGS=(
    "${laminateTypesDir}/tags.yaml"
    "${paperColorsDir}/tags.yaml"
    "${paperFacturesDir}/tags.yaml"
    "${printFormatsDir}/tags.yaml"
  )

  OPENAPI_PATHS=(
    "${laminateTypesDir}/paths.yaml"
    "${paperColorsDir}/paths.yaml"
    "${paperFacturesDir}/paths.yaml"
    "${printFormatsDir}/paths.yaml"
  )

#  OPENAPI_COMPONENTS_HEADERS=(
#    "${sharedDir}/components/headers/"
#  )

  OPENAPI_COMPONENTS_PARAMETERS=(
    "${sharedDir}/components/parameters/App.Request.Header.AcceptLanguage.yaml"
    "${sharedDir}/components/parameters/App.Request.Header.CorrelationID.yaml"
    # "${sharedDir}/components/parameters/App.Request.Header.CurrentPage.yaml"
    "${sharedDir}/components/parameters/App.Request.Query.Filter.SearchText.yaml"
    "${sharedDir}/components/parameters/App.Request.Query.Filter.Statuses.yaml"
    "${sharedDir}/components/parameters/App.Request.Query.ListPager.yaml"
    "${sharedDir}/components/parameters/App.Request.Query.ListSorter.yaml"

    # "${sharedDir}/custom/parameters/Custom.Request.Query.Filter.DensityRange.yaml"
    # "${sharedDir}/custom/parameters/Custom.Request.Query.Filter.DepthRange.yaml"
    # "${sharedDir}/custom/parameters/Custom.Request.Query.Filter.ElementDetailing.yaml"
    "${sharedDir}/custom/parameters/Custom.Request.Query.Filter.LengthRange.yaml"
    # "${sharedDir}/custom/parameters/Custom.Request.Query.Filter.WeightRange.yaml"
    "${sharedDir}/custom/parameters/Custom.Request.Query.Filter.WidthRange.yaml"

    "${laminateTypesDir}/components-parameters.yaml"
    "${paperColorsDir}/components-parameters.yaml"
    "${paperFacturesDir}/components-parameters.yaml"
    "${printFormatsDir}/components-parameters.yaml"
  )

  OPENAPI_COMPONENTS_SCHEMAS=(
    # "${sharedDir}/components/schemas/enums/App.Enum.Address.HouseType.yaml"
    # "${sharedDir}/components/schemas/enums/App.Enum.DeliveryMethod.yaml"
    # "${sharedDir}/components/schemas/enums/App.Enum.Gender.yaml"
    "${sharedDir}/components/schemas/enums/App.Enum.Status.yaml"

    # "${sharedDir}/components/schemas/fields/App.Field.Article.yaml"
    # "${sharedDir}/components/schemas/fields/App.Field.Boolean.yaml"
    "${sharedDir}/components/schemas/fields/App.Field.Caption.yaml"
    "${sharedDir}/components/schemas/fields/App.Field.Date.CreatedAt.yaml"
    "${sharedDir}/components/schemas/fields/App.Field.Date.UpdatedAt.yaml"
    # "${sharedDir}/components/schemas/fields/App.Field.Date.yaml"
    # "${sharedDir}/components/schemas/fields/App.Field.Datetime.yaml"
    # "${sharedDir}/components/schemas/fields/App.Field.Email.yaml"
    # "${sharedDir}/components/schemas/fields/App.Field.GEO.yaml"
    # "${sharedDir}/components/schemas/fields/App.Field.ImageURL.yaml"
    "${sharedDir}/components/schemas/fields/App.Field.IntegerID.yaml"
    "${sharedDir}/components/schemas/fields/App.Field.ListPager.Total.yaml"
    # "${sharedDir}/components/schemas/fields/App.Field.Phone.yaml"
    # "${sharedDir}/components/schemas/fields/App.Field.StringID.yaml"
    "${sharedDir}/components/schemas/fields/App.Field.TagVersion.yaml"
    # "${sharedDir}/components/schemas/fields/App.Field.Timezone.yaml"
    # "${sharedDir}/components/schemas/fields/App.Field.UUID.yaml"
    # "${sharedDir}/components/schemas/fields/App.Field.VariableCamelCase.yaml"

    # "${sharedDir}/components/schemas/measures/App.Measure.Gram.yaml"
    # "${sharedDir}/components/schemas/measures/App.Measure.GramPerMeter2.yaml"
    "${sharedDir}/components/schemas/measures/App.Measure.Micrometer.yaml"
    # "${sharedDir}/components/schemas/measures/App.Measure.Price.yaml"

    # "${sharedDir}/components/schemas/App.Request.Model.ChangeFlag.yaml"
    "${sharedDir}/components/schemas/App.Request.Model.ChangeStatus.yaml"
    # "${sharedDir}/components/schemas/App.Request.Model.MoveItem.yaml"
    # "${sharedDir}/components/schemas/App.Response.Model.BinaryFile.yaml"
    "${sharedDir}/components/schemas/App.Response.Model.Error.yaml"
    # "${sharedDir}/components/schemas/App.Response.Model.FileInfo.yaml"
    # "${sharedDir}/components/schemas/App.Response.Model.ImageInfo.yaml"
    # "${sharedDir}/components/schemas/App.Response.Model.Success.yaml"
    "${sharedDir}/components/schemas/App.Response.Model.SuccessCreatedItem.yaml"

    # "${sharedDir}/custom/enums/Custom.Enum.CompanyPublicStatus.yaml"
    # "${sharedDir}/custom/enums/Custom.Enum.FormElementDetailing.yaml"
    # "${sharedDir}/custom/enums/Custom.Enum.FormElementType.yaml"
    # "${sharedDir}/custom/enums/Custom.Enum.PaperSides.yaml"

    # "${sharedDir}/custom/fields/Custom.Request.Query.Field.ParamName.yaml"

    "${laminateTypesDir}/components-schemas.yaml"
    "${paperColorsDir}/components-schemas.yaml"
    "${paperFacturesDir}/components-schemas.yaml"
    "${printFormatsDir}/components-schemas.yaml"
  )

  OPENAPI_COMPONENTS_RESPONSES=(
    "${sharedDir}/components/responses/App.ResponseJson.Errors.yaml"
    "${sharedDir}/components/responses/App.ResponseJson.ErrorsAuth.yaml"

    "${laminateTypesDir}/components-responses.yaml"
    "${paperColorsDir}/components-responses.yaml"
    "${paperFacturesDir}/components-responses.yaml"
    "${printFormatsDir}/components-responses.yaml"
  )

  OPENAPI_SECURITY_SCHEMES=(
    "${sharedDir}/securitySchemes.yaml"
  )
}
