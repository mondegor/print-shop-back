
mrcmd_func_openapi_build_prov() {
  local sectionDir="${1:?}" # sample: .../admin-api, .../public-api
  local sharedDir="${2:?}" # sample: .../_shared

  local companyPageDir="${sectionDir}/company-page"

  # OPENAPI_VERSION="3.0.3"

  OPENAPI_HEADERS=(
    "${sectionDir}/header.yaml"
    "${sharedDir}/description-errors.md"
  )

  OPENAPI_SERVERS=(
    "${sectionDir}/servers.yaml"
  )

  OPENAPI_TAGS=(
    "${sharedDir}/system/tags.yaml"

    "${companyPageDir}/tags.yaml"
  )

  OPENAPI_PATHS=(
    "${sharedDir}/system/paths.yaml"

    "${companyPageDir}/company_page_paths.yaml"
  )

#  OPENAPI_COMPONENTS_HEADERS=(
#    "${sharedDir}/components/headers/"
#  )

  OPENAPI_COMPONENTS_PARAMETERS=(
    "${sharedDir}/components/parameters/App.Request.Header.AcceptLanguage.yaml"
    "${sharedDir}/components/parameters/App.Request.Header.CorrelationID.yaml"
    # "${sharedDir}/components/parameters/App.Request.Header.CurrentPage.yaml"
    # "${sharedDir}/components/parameters/App.Request.Query.Filter.SearchText.yaml"
    # "${sharedDir}/components/parameters/App.Request.Query.Filter.Statuses.yaml"
    # "${sharedDir}/components/parameters/App.Request.Query.ListPager.yaml"
    # "${sharedDir}/components/parameters/App.Request.Query.ListSorter.yaml"

    # "${sharedDir}/custom/parameters/Custom.Request.Query.Filter.DensityRange.yaml"
    # "${sharedDir}/custom/parameters/Custom.Request.Query.Filter.ElementDetailing.yaml"
    # "${sharedDir}/custom/parameters/Custom.Request.Query.Filter.HeightRange.yaml"
    # "${sharedDir}/custom/parameters/Custom.Request.Query.Filter.LengthRange.yaml"
    # "${sharedDir}/custom/parameters/Custom.Request.Query.Filter.MaterialTypeIDs.yaml"
    # "${sharedDir}/custom/parameters/Custom.Request.Query.Filter.PriceRange.yaml"
    # "${sharedDir}/custom/parameters/Custom.Request.Query.Filter.WeightRange.yaml"
    # "${sharedDir}/custom/parameters/Custom.Request.Query.Filter.WidthRange.yaml"

    "${companyPageDir}/company_page_parameters.yaml"
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
    # "${sharedDir}/components/fields/App.Field.Email.yaml"
    # "${sharedDir}/components/fields/App.Field.GEO.yaml"
    # "${sharedDir}/components/fields/App.Field.ImageURL.yaml"
    # "${sharedDir}/components/fields/App.Field.Int16.yaml"
    # "${sharedDir}/components/fields/App.Field.Int32.yaml"
    # "${sharedDir}/components/fields/App.Field.ListPager.Total.yaml"
    # "${sharedDir}/components/fields/App.Field.Phone.yaml"
    # "${sharedDir}/components/fields/App.Field.RewriteName.yaml"
    # "${sharedDir}/components/fields/App.Field.TagVersion.yaml"
    # "${sharedDir}/components/fields/App.Field.Timezone.yaml"
    # "${sharedDir}/components/fields/App.Field.UUID.yaml"
    # "${sharedDir}/components/fields/App.Field.VariableCamelCase.yaml"

    # "${sharedDir}/components/fields/measures/App.Field.Measure.Gram.yaml"
    # "${sharedDir}/components/fields/measures/App.Field.Measure.GramPerMeter2.yaml"
    # "${sharedDir}/components/fields/measures/App.Field.Measure.Micrometer.yaml"
    # "${sharedDir}/components/fields/measures/App.Field.Measure.Milligram.yaml"
    # "${sharedDir}/components/fields/measures/App.Field.Measure.Price.yaml"

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

    "${sharedDir}/custom/enums/Custom.Enum.CompanyPublicStatus.yaml"
    # "${sharedDir}/custom/enums/Custom.Enum.FormElementDetailing.yaml"
    # "${sharedDir}/custom/enums/Custom.Enum.FormElementType.yaml"
    # "${sharedDir}/custom/enums/Custom.Enum.PaperSides.yaml"
    # "${sharedDir}/custom/enums/Custom.Enum.SubmitFormActivityStatus.yaml"

    # "${sharedDir}/custom/fields/Custom.Field.Catalog.BoxID.yaml"
    # "${sharedDir}/custom/fields/Custom.Field.Catalog.LaminateID.yaml"
    # "${sharedDir}/custom/fields/Custom.Field.Catalog.PaperID.yaml"
    # "${sharedDir}/custom/fields/Custom.Field.Controls.ElementTemplateID.yaml"
    # "${sharedDir}/custom/fields/Custom.Field.Controls.FormElementID.yaml"
    # "${sharedDir}/custom/fields/Custom.Field.Controls.SubmitFormID.yaml"
    # "${sharedDir}/custom/fields/Custom.Field.Dictionaries.MaterialTypeID.yaml"
    # "${sharedDir}/custom/fields/Custom.Field.Dictionaries.PaperColorID.yaml"
    # "${sharedDir}/custom/fields/Custom.Field.Dictionaries.PaperFactureID.yaml"
    # "${sharedDir}/custom/fields/Custom.Field.Dictionaries.PrintFormatID.yaml"

    "${sharedDir}/system/schemas.yaml"

    "${companyPageDir}/company_page_schemas.yaml"
  )

  OPENAPI_COMPONENTS_RESPONSES=(
    "${sharedDir}/components/responses/App.ResponseJson.Errors.yaml"
    "${sharedDir}/components/responses/App.ResponseJson.ErrorsAuth.yaml"
  )

  OPENAPI_SECURITY_SCHEMES=(
    "${sharedDir}/securitySchemes.yaml"
  )
}
