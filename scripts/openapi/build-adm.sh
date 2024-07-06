
mrcmd_func_openapi_build_adm() {
  local sectionDir="${1:?}" # sample: .../admin-api, .../public-api
  local sharedDir="${2:?}" # sample: .../_shared

  local catalogDir="${sectionDir}/catalog"
  local boxDir="${catalogDir}/box"
  local laminateDir="${catalogDir}/laminate"
  local paperDir="${catalogDir}/paper"

  local controlsDir="${sectionDir}/controls"
  local elementTemplateDir="${controlsDir}/element-template"
  local submitFormDir="${controlsDir}/submit-form"

  local dictionariesDir="${sectionDir}/dictionaries"
  local materialTypeDir="${dictionariesDir}/material-type"
  local paperColorDir="${dictionariesDir}/paper-color"
  local paperFactureDir="${dictionariesDir}/paper-facture"
  local printFormatDir="${dictionariesDir}/print-format"

  local companyPageDir="${sectionDir}/provider-accounts/company-page"

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

    "${boxDir}/tags.yaml"
    "${laminateDir}/tags.yaml"
    "${paperDir}/tags.yaml"

    "${materialTypeDir}/tags.yaml"
    "${paperColorDir}/tags.yaml"
    "${paperFactureDir}/tags.yaml"
    "${printFormatDir}/tags.yaml"

    "${elementTemplateDir}/tags.yaml"
    "${submitFormDir}/tags.yaml"

    "${companyPageDir}/tags.yaml"
  )

  OPENAPI_PATHS=(
    "${sharedDir}/system/paths.yaml"

    "${boxDir}/box_paths.yaml"
    "${laminateDir}/laminate_paths.yaml"
    "${paperDir}/paper_paths.yaml"

    "${materialTypeDir}/material_type_paths.yaml"
    "${paperColorDir}/paper_color_paths.yaml"
    "${paperFactureDir}/paper_facture_paths.yaml"
    "${printFormatDir}/print_format_paths.yaml"

    "${elementTemplateDir}/element_template_paths.yaml"
    "${submitFormDir}/submit_form_paths.yaml"
    "${submitFormDir}/element_paths.yaml"

    "${companyPageDir}/company_page_paths.yaml"
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

    "${sharedDir}/custom/parameters/Custom.Request.Query.Filter.DensityRange.yaml"
    "${sharedDir}/custom/parameters/Custom.Request.Query.Filter.ElementDetailing.yaml"
    "${sharedDir}/custom/parameters/Custom.Request.Query.Filter.HeightRange.yaml"
    "${sharedDir}/custom/parameters/Custom.Request.Query.Filter.LengthRange.yaml"
    "${sharedDir}/custom/parameters/Custom.Request.Query.Filter.MaterialTypeIDs.yaml"
    # "${sharedDir}/custom/parameters/Custom.Request.Query.Filter.PriceRange.yaml"
    "${sharedDir}/custom/parameters/Custom.Request.Query.Filter.WeightRange.yaml"
    "${sharedDir}/custom/parameters/Custom.Request.Query.Filter.WidthRange.yaml"

    "${boxDir}/box_parameters.yaml"
    "${laminateDir}/laminate_parameters.yaml"
    "${paperDir}/paper_parameters.yaml"

    "${materialTypeDir}/material_type_parameters.yaml"
    "${paperColorDir}/paper_color_parameters.yaml"
    "${paperFactureDir}/paper_facture_parameters.yaml"
    "${printFormatDir}/print_format_parameters.yaml"

    "${elementTemplateDir}/element_template_parameters.yaml"
    "${submitFormDir}/submit_form_parameters.yaml"
    "${submitFormDir}/element_parameters.yaml"

    "${companyPageDir}/company_page_parameters.yaml"
  )

  OPENAPI_COMPONENTS_SCHEMAS=(
    # "${sharedDir}/components/enums/App.Enum.Address.HouseType.yaml"
    # "${sharedDir}/components/enums/App.Enum.DeliveryMethod.yaml"
    # "${sharedDir}/components/enums/App.Enum.Gender.yaml"
    "${sharedDir}/components/enums/App.Enum.Status.yaml"

    "${sharedDir}/components/fields/App.Field.Article.yaml"
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
    # "${sharedDir}/components/fields/App.Field.Phone.yaml"
    "${sharedDir}/components/fields/App.Field.RewriteName.yaml"
    "${sharedDir}/components/fields/App.Field.TagVersion.yaml"
    # "${sharedDir}/components/fields/App.Field.Timezone.yaml"
    # "${sharedDir}/components/fields/App.Field.UUID.yaml"
    "${sharedDir}/components/fields/App.Field.VariableCamelCase.yaml"

    # "${sharedDir}/components/fields/measures/App.Field.Measure.Centimeter.yaml"
    "${sharedDir}/components/fields/measures/App.Field.Measure.Gram.yaml"
    "${sharedDir}/components/fields/measures/App.Field.Measure.GramPerMeter2.yaml"
    # "${sharedDir}/components/fields/measures/App.Field.Measure.Kilogram.yaml"
    # "${sharedDir}/components/fields/measures/App.Field.Measure.KilogramPerMeter2.yaml"
    "${sharedDir}/components/fields/measures/App.Field.Measure.Meter.yaml"
    "${sharedDir}/components/fields/measures/App.Field.Measure.Micrometer.yaml"
    # "${sharedDir}/components/fields/measures/App.Field.Measure.Milligram.yaml"
    "${sharedDir}/components/fields/measures/App.Field.Measure.Millimeter.yaml"
    # "${sharedDir}/components/fields/measures/App.Field.Measure.Price.yaml"

    # "${sharedDir}/components/models/App.Request.Model.ChangeFlag.yaml"
    "${sharedDir}/components/models/App.Request.Model.ChangeStatus.yaml"
    "${sharedDir}/components/models/App.Request.Model.MoveItem.yaml"
    # "${sharedDir}/components/models/App.Response.Model.BinaryAnyFile.yaml"
    # "${sharedDir}/components/models/App.Response.Model.BinaryFile.yaml"
    # "${sharedDir}/components/models/App.Response.Model.BinaryImage.yaml"
    # "${sharedDir}/components/models/App.Response.Model.BinaryMedia.yaml"
    "${sharedDir}/components/models/App.Response.Model.Error.yaml"
    # "${sharedDir}/components/models/App.Response.Model.FileInfo.yaml"
    "${sharedDir}/components/models/App.Response.Model.JsonFile.yaml"
    # "${sharedDir}/components/models/App.Response.Model.ImageInfo.yaml"
    # "${sharedDir}/components/models/App.Response.Model.Success.yaml"
    # "${sharedDir}/components/models/App.Response.Model.SuccessCreatedItem.yaml"
    "${sharedDir}/components/models/App.Response.Model.SuccessCreatedItemInt32.yaml"

    # "${sharedDir}/custom/enums/Custom.Enum.CompanyPublicStatus.yaml"
    "${sharedDir}/custom/enums/Custom.Enum.FormElementDetailing.yaml"
    "${sharedDir}/custom/enums/Custom.Enum.FormElementType.yaml"
    "${sharedDir}/custom/enums/Custom.Enum.PaperSides.yaml"
    "${sharedDir}/custom/enums/Custom.Enum.SubmitFormActivityStatus.yaml"

    "${sharedDir}/custom/fields/Custom.Field.Catalog.BoxID.yaml"
    "${sharedDir}/custom/fields/Custom.Field.Catalog.LaminateID.yaml"
    "${sharedDir}/custom/fields/Custom.Field.Catalog.PaperID.yaml"
    "${sharedDir}/custom/fields/Custom.Field.Controls.ElementTemplateID.yaml"
    "${sharedDir}/custom/fields/Custom.Field.Controls.FormElementID.yaml"
    "${sharedDir}/custom/fields/Custom.Field.Controls.SubmitFormID.yaml"
    "${sharedDir}/custom/fields/Custom.Field.Dictionaries.MaterialTypeID.yaml"
    "${sharedDir}/custom/fields/Custom.Field.Dictionaries.PaperColorID.yaml"
    "${sharedDir}/custom/fields/Custom.Field.Dictionaries.PaperFactureID.yaml"
    "${sharedDir}/custom/fields/Custom.Field.Dictionaries.PrintFormatID.yaml"
    # "${sharedDir}/custom/fields/Custom.Field.DoubleSize.yaml"
    # "${sharedDir}/custom/fields/Custom.Field.Fragment.yaml"
    # "${sharedDir}/custom/fields/Custom.Field.RectFormat.yaml"

    "${sharedDir}/system/schemas.yaml"

    "${boxDir}/box_schemas.yaml"
    "${laminateDir}/laminate_schemas.yaml"
    "${paperDir}/paper_schemas.yaml"

    "${materialTypeDir}/material_type_schemas.yaml"
    "${paperColorDir}/paper_color_schemas.yaml"
    "${paperFactureDir}/paper_facture_schemas.yaml"
    "${printFormatDir}/print_format_schemas.yaml"

    "${elementTemplateDir}/element_template_schemas.yaml"
    "${submitFormDir}/submit_form_schemas.yaml"
    "${submitFormDir}/element_schemas.yaml"

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
