
mrcmd_func_openapi_build_adm_controls() {
  local sectionDir="${1:?}" # sample: .../admin-api, .../public-api
  local sharedDir="${2:?}" # sample: .../_shared

  local controlsDir="${sectionDir}/controls"
  local elementTemplateDir="${controlsDir}/element-template"
  local submitFormDir="${controlsDir}/submit-form"

  # OPENAPI_VERSION="3.0.3"

  OPENAPI_HEADERS=(
    "${controlsDir}/header.yaml"
    "${sharedDir}/description-errors.md"
  )

  OPENAPI_SERVERS=(
    "${sectionDir}/servers.yaml"
  )

  OPENAPI_TAGS=(
    "${elementTemplateDir}/tags.yaml"
    "${submitFormDir}/tags.yaml"
  )

  OPENAPI_PATHS=(
    "${elementTemplateDir}/element_template_paths.yaml"
    "${submitFormDir}/submit_form_paths.yaml"
    "${submitFormDir}/element_paths.yaml"
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
    "${sharedDir}/custom/parameters/Custom.Request.Query.Filter.ElementDetailing.yaml"
    # "${sharedDir}/custom/parameters/Custom.Request.Query.Filter.HeightRange.yaml"
    # "${sharedDir}/custom/parameters/Custom.Request.Query.Filter.LengthRange.yaml"
    # "${sharedDir}/custom/parameters/Custom.Request.Query.Filter.MaterialTypeIDs.yaml"
    # "${sharedDir}/custom/parameters/Custom.Request.Query.Filter.PriceRange.yaml"
    # "${sharedDir}/custom/parameters/Custom.Request.Query.Filter.WeightRange.yaml"
    # "${sharedDir}/custom/parameters/Custom.Request.Query.Filter.WidthRange.yaml"

    "${elementTemplateDir}/element_template_parameters.yaml"
    "${submitFormDir}/submit_form_parameters.yaml"
    "${submitFormDir}/element_parameters.yaml"
  )

  OPENAPI_COMPONENTS_SCHEMAS=(
    # "${sharedDir}/components/enums/App.Enum.Address.HouseType.yaml"
    # "${sharedDir}/components/enums/App.Enum.DeliveryMethod.yaml"
    # "${sharedDir}/components/enums/App.Enum.Gender.yaml"
    "${sharedDir}/components/enums/App.Enum.Status.yaml"

    # "${sharedDir}/components/fields/App.Field.Article.yaml"
    # "${sharedDir}/components/fields/App.Field.Boolean.yaml"
    "${sharedDir}/components/fields/App.Field.Caption.yaml"
    "${sharedDir}/components/fields/App.Field.DateTimeCreatedAt.yaml"
    "${sharedDir}/components/fields/App.Field.DateTimeUpdatedAt.yaml"
    # "${sharedDir}/components/fields/App.Field.Date.yaml"
    # "${sharedDir}/components/fields/App.Field.DateTime.yaml"
    # "${sharedDir}/components/fields/App.Field.Email.yaml"
    # "${sharedDir}/components/fields/App.Field.ExternalURL.yaml"
    # "${sharedDir}/components/fields/App.Field.FileURL.yaml"
    # "${sharedDir}/components/fields/App.Field.Float64.yaml"
    # "${sharedDir}/components/fields/App.Field.GEO.yaml"
    # "${sharedDir}/components/fields/App.Field.ImageURL.yaml"
    # "${sharedDir}/components/fields/App.Field.Int16.yaml"
    "${sharedDir}/components/fields/App.Field.Int32.yaml"
    # "${sharedDir}/components/fields/App.Field.Int64.yaml"
    # "${sharedDir}/components/fields/App.Field.JsonData.yaml"
    "${sharedDir}/components/fields/App.Field.ListPager.Total.yaml"
    # "${sharedDir}/components/fields/App.Field.OrderIndex.yaml"
    # "${sharedDir}/components/fields/App.Field.Percent.yaml"
    # "${sharedDir}/components/fields/App.Field.Phone.yaml"
    "${sharedDir}/components/fields/App.Field.RewriteName.yaml"
    # "${sharedDir}/components/fields/App.Field.Size2D.yaml"
    # "${sharedDir}/components/fields/App.Field.Size3D.yaml"
    "${sharedDir}/components/fields/App.Field.TagVersion.yaml"
    # "${sharedDir}/components/fields/App.Field.Timezone.yaml"
    # "${sharedDir}/components/fields/App.Field.Uint.yaml"
    # "${sharedDir}/components/fields/App.Field.UUID.yaml"
    "${sharedDir}/components/fields/App.Field.VariableCamelCase.yaml"

    # "${sharedDir}/components/fields/measures/App.Field.Measure.Centimeter.yaml"
    # "${sharedDir}/components/fields/measures/App.Field.Measure.Gram.yaml"
    # "${sharedDir}/components/fields/measures/App.Field.Measure.GramPerMeter2.yaml"
    # "${sharedDir}/components/fields/measures/App.Field.Measure.Kilogram.yaml"
    # "${sharedDir}/components/fields/measures/App.Field.Measure.KilogramPerMeterS2.yaml"
    # "${sharedDir}/components/fields/measures/App.Field.Measure.Meter.yaml"
    # "${sharedDir}/components/fields/measures/App.Field.Measure.Meter2.yaml"
    # "${sharedDir}/components/fields/measures/App.Field.Measure.Meter3.yaml"
    # "${sharedDir}/components/fields/measures/App.Field.Measure.Micrometer.yaml"
    # "${sharedDir}/components/fields/measures/App.Field.Measure.Milligram.yaml"
    # "${sharedDir}/components/fields/measures/App.Field.Measure.Millimeter.yaml"
    # "${sharedDir}/components/fields/measures/App.Field.Measure.Millimeter2.yaml"
    # "${sharedDir}/components/fields/measures/App.Field.Measure.Millimeter3.yaml"
    # "${sharedDir}/components/fields/measures/App.Field.Measure.Price.yaml"

    # "${sharedDir}/components/models/App.Request.Model.ChangeFlag.yaml"
    "${sharedDir}/components/models/App.Request.Model.ChangeStatus.yaml"
    "${sharedDir}/components/models/App.Request.Model.MoveItem.yaml"
    # "${sharedDir}/components/models/App.Response.Model.BinaryAnyFile.yaml"
    # "${sharedDir}/components/models/App.Response.Model.BinaryFile.yaml"
    "${sharedDir}/components/models/App.Response.Model.BinaryImage.yaml"
    # "${sharedDir}/components/models/App.Response.Model.BinaryMedia.yaml"
    "${sharedDir}/components/models/App.Response.Model.Error.yaml"
    # "${sharedDir}/components/models/App.Response.Model.FileInfo.yaml"
    "${sharedDir}/components/models/App.Response.Model.JsonFile.yaml"
    # "${sharedDir}/components/models/App.Response.Model.ImageInfo.yaml"
    # "${sharedDir}/components/models/App.Response.Model.Success.yaml"
    # "${sharedDir}/components/models/App.Response.Model.SuccessCreatedItem.yaml"
    "${sharedDir}/components/models/App.Response.Model.SuccessCreatedItemUint.yaml"

    # "${sharedDir}/custom/enums/Custom.Enum.CompanyPublicStatus.yaml"
    "${sharedDir}/custom/enums/Custom.Enum.FormElementDetailing.yaml"
    "${sharedDir}/custom/enums/Custom.Enum.FormElementType.yaml"
    # "${sharedDir}/custom/enums/Custom.Enum.FragmentPosition.yaml"
    # "${sharedDir}/custom/enums/Custom.Enum.PaperSides.yaml"
    "${sharedDir}/custom/enums/Custom.Enum.SubmitFormActivityStatus.yaml"

    # "${sharedDir}/custom/fields/Custom.Field.Catalog.BoxID.yaml"
    # "${sharedDir}/custom/fields/Custom.Field.Catalog.LaminateID.yaml"
    # "${sharedDir}/custom/fields/Custom.Field.Catalog.PaperID.yaml"
    "${sharedDir}/custom/fields/Custom.Field.Controls.ElementTemplateID.yaml"
    "${sharedDir}/custom/fields/Custom.Field.Controls.FormElementID.yaml"
    "${sharedDir}/custom/fields/Custom.Field.Controls.SubmitFormID.yaml"
    # "${sharedDir}/custom/fields/Custom.Field.Dictionaries.MaterialTypeID.yaml"
    # "${sharedDir}/custom/fields/Custom.Field.Dictionaries.PaperColorID.yaml"
    # "${sharedDir}/custom/fields/Custom.Field.Dictionaries.PaperFactureID.yaml"
    # "${sharedDir}/custom/fields/Custom.Field.Dictionaries.PrintFormatID.yaml"
    # "${sharedDir}/custom/fields/Custom.Field.Fragment.yaml"
    # "${sharedDir}/custom/fields/Custom.Field.Layout.yaml"
    # "${sharedDir}/custom/fields/Custom.Field.Rect2dFormat.yaml"

    "${elementTemplateDir}/element_template_schemas.yaml"
    "${submitFormDir}/submit_form_schemas.yaml"
    "${submitFormDir}/element_schemas.yaml"
  )

  OPENAPI_COMPONENTS_RESPONSES=(
    "${sharedDir}/components/responses/App.ResponseJson.Errors.yaml"
    "${sharedDir}/components/responses/App.ResponseJson.ErrorsAuth.yaml"
  )

  OPENAPI_SECURITY_SCHEMES=(
    "${sharedDir}/securitySchemes.yaml"
  )
}
