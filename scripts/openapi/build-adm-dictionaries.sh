
mrcmd_func_openapi_build_adm_dictionaries() {
  local sectionDir="${1:?}" # sample: .../admin-api, .../public-api
  local sharedDir="${2:?}" # sample: .../_shared

  local dictionariesDir="${sectionDir}/dictionaries"
  local laminateTypeDir="${dictionariesDir}/laminate-type"
  local paperColorDir="${dictionariesDir}/paper-color"
  local paperFactureDir="${dictionariesDir}/paper-facture"
  local printFormatDir="${dictionariesDir}/print-format"

  # OPENAPI_VERSION="3.0.3"

  OPENAPI_HEADERS=(
    "${dictionariesDir}/header.yaml"
    "${sharedDir}/description-errors.md"
  )

  OPENAPI_SERVERS=(
    "${sectionDir}/servers.yaml"
  )

  OPENAPI_TAGS=(
    "${laminateTypeDir}/tags.yaml"
    "${paperColorDir}/tags.yaml"
    "${paperFactureDir}/tags.yaml"
    "${printFormatDir}/tags.yaml"
  )

  OPENAPI_PATHS=(
    "${laminateTypeDir}/laminate_type_paths.yaml"
    "${paperColorDir}/paper_color_paths.yaml"
    "${paperFactureDir}/paper_facture_paths.yaml"
    "${printFormatDir}/print_format_paths.yaml"
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
    # "${sharedDir}/custom/parameters/Custom.Request.Query.Filter.PriceRange.yaml"
    # "${sharedDir}/custom/parameters/Custom.Request.Query.Filter.WeightRange.yaml"
    "${sharedDir}/custom/parameters/Custom.Request.Query.Filter.WidthRange.yaml"

    "${laminateTypeDir}/laminate_type_parameters.yaml"
    "${paperColorDir}/paper_color_parameters.yaml"
    "${paperFactureDir}/paper_facture_parameters.yaml"
    "${printFormatDir}/print_format_parameters.yaml"
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
    # "${sharedDir}/components/fields/App.Field.GEO.yaml"
    # "${sharedDir}/components/fields/App.Field.ImageURL.yaml"
    # "${sharedDir}/components/fields/App.Field.Int16.yaml"
    # "${sharedDir}/components/fields/App.Field.Int32.yaml"
    "${sharedDir}/components/fields/App.Field.ListPager.Total.yaml"
    # "${sharedDir}/components/fields/App.Field.Phone.yaml"
    # "${sharedDir}/components/fields/App.Field.StringID.yaml"
    "${sharedDir}/components/fields/App.Field.TagVersion.yaml"
    # "${sharedDir}/components/fields/App.Field.Timezone.yaml"
    # "${sharedDir}/components/fields/App.Field.UUID.yaml"
    # "${sharedDir}/components/fields/App.Field.VariableCamelCase.yaml"

    # "${sharedDir}/components/fields/measures/App.Field.Measure.Gram.yaml"
    # "${sharedDir}/components/fields/measures/App.Field.Measure.GramPerMeter2.yaml"
    "${sharedDir}/components/fields/measures/App.Field.Measure.Micrometer.yaml"
    # "${sharedDir}/components/fields/measures/App.Field.Measure.Price.yaml"

    # "${sharedDir}/components/models/App.Request.Model.ChangeFlag.yaml"
    "${sharedDir}/components/models/App.Request.Model.ChangeStatus.yaml"
    # "${sharedDir}/components/models/App.Request.Model.MoveItem.yaml"
    # "${sharedDir}/components/models/App.Response.Model.BinaryFile.yaml"
    "${sharedDir}/components/models/App.Response.Model.Error.yaml"
    # "${sharedDir}/components/models/App.Response.Model.FileInfo.yaml"
    # "${sharedDir}/components/models/App.Response.Model.ImageInfo.yaml"
    # "${sharedDir}/components/models/App.Response.Model.Success.yaml"
    "${sharedDir}/components/models/App.Response.Model.SuccessCreatedItem.yaml"

    # "${sharedDir}/custom/enums/Custom.Enum.CompanyPublicStatus.yaml"
    # "${sharedDir}/custom/enums/Custom.Enum.FormElementDetailing.yaml"
    # "${sharedDir}/custom/enums/Custom.Enum.FormElementType.yaml"
    # "${sharedDir}/custom/enums/Custom.Enum.PaperSides.yaml"

    # "${sharedDir}/custom/fields/Custom.Field.Catalog.BoxID.yaml"
    # "${sharedDir}/custom/fields/Custom.Field.Catalog.LaminateID.yaml"
    # "${sharedDir}/custom/fields/Custom.Field.Catalog.PaperID.yaml"
    # "${sharedDir}/custom/fields/Custom.Field.Controls.ElementID.yaml"
    # "${sharedDir}/custom/fields/Custom.Field.Controls.FormID.yaml"
    # "${sharedDir}/custom/fields/Custom.Field.Controls.TemplateID.yaml"
    "${sharedDir}/custom/fields/Custom.Field.Dictionaries.ColorID.yaml"
    "${sharedDir}/custom/fields/Custom.Field.Dictionaries.FactureID.yaml"
    "${sharedDir}/custom/fields/Custom.Field.Dictionaries.LaminateTypeID.yaml"
    "${sharedDir}/custom/fields/Custom.Field.Dictionaries.PrintFormatID.yaml"
    # "${sharedDir}/custom/fields/Custom.Field.ParamName.yaml"

    "${laminateTypeDir}/laminate_type_schemas.yaml"
    "${paperColorDir}/paper_color_schemas.yaml"
    "${paperFactureDir}/paper_facture_schemas.yaml"
    "${printFormatDir}/print_format_schemas.yaml"
  )

  OPENAPI_COMPONENTS_RESPONSES=(
    "${sharedDir}/components/responses/App.ResponseJson.Errors.yaml"
    "${sharedDir}/components/responses/App.ResponseJson.ErrorsAuth.yaml"
  )

  OPENAPI_SECURITY_SCHEMES=(
    "${sharedDir}/securitySchemes.yaml"
  )
}
