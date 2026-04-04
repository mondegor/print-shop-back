
mrcmd_func_openapi_build_pub() {
  local sectionDir="${1:?}" # sample: .../admin-api, .../public-api
  local sharedDir="${2:?}" # sample: .../_shared

  local calculationsDir="${sectionDir}/calculations"
  local calculationsAlgoDir="${calculationsDir}/algo"
  local calculationsQueryHistoryDir="${calculationsDir}/queryhistory"

  local catalogDir="${sectionDir}/catalog"
  local boxDir="${catalogDir}/box"
  local laminateDir="${catalogDir}/laminate"
  local paperDir="${catalogDir}/paper"

  local companyPageDir="${sectionDir}/company-page"

  local controlsDir="${sectionDir}/controls"
  local submitFormDir="${controlsDir}/submit-form"

  local dictionariesDir="${sectionDir}/dictionaries"
  local materialTypeDir="${dictionariesDir}/material-type"
  local paperColorDir="${dictionariesDir}/paper-color"
  local paperFactureDir="${dictionariesDir}/paper-facture"
  local printFormatDir="${dictionariesDir}/print-format"

  local fileStationDir="${sectionDir}/file-station"

  # OPENAPI_VERSION="3.0.3"

  OPENAPI_HEADERS=(
    "${sectionDir}/header.yaml"
    "${sharedDir}/description-errors.md"
  )

  OPENAPI_SERVERS=(
    "${sectionDir}/servers.yaml"
  )

  OPENAPI_TAGS=(
    "${calculationsAlgoDir}/tags.yaml"
    "${calculationsQueryHistoryDir}/tags.yaml"

    "${boxDir}/tags.yaml"
    "${laminateDir}/tags.yaml"
    "${paperDir}/tags.yaml"

    "${companyPageDir}/tags.yaml"

    "${submitFormDir}/tags.yaml"

    "${materialTypeDir}/tags.yaml"
    "${paperColorDir}/tags.yaml"
    "${paperFactureDir}/tags.yaml"
    "${printFormatDir}/tags.yaml"

    "${fileStationDir}/tags.yaml"
  )

  OPENAPI_PATHS=(
    "${calculationsAlgoDir}/box/packinbox/packinbox_paths.yaml"
    "${calculationsAlgoDir}/sheet/cutting/cutting_paths.yaml"
    "${calculationsAlgoDir}/sheet/imposition/imposition_paths.yaml"
    "${calculationsAlgoDir}/sheet/insideoutside/inside_outside_paths.yaml"
    "${calculationsAlgoDir}/sheet/packinstack/packinstack_paths.yaml"
    "${calculationsQueryHistoryDir}/query_history_paths.yaml"

    "${boxDir}/box_paths.yaml"
    "${laminateDir}/laminate_paths.yaml"
    "${paperDir}/paper_paths.yaml"

    "${companyPageDir}/company_page_paths.yaml"

    "${submitFormDir}/submit_form_paths.yaml"

    "${materialTypeDir}/material_type_paths.yaml"
    "${paperColorDir}/paper_color_paths.yaml"
    "${paperFactureDir}/paper_facture_paths.yaml"
    "${printFormatDir}/print_format_paths.yaml"

    "${fileStationDir}/paths.yaml"
  )

#  OPENAPI_COMPONENTS_HEADERS=(
#    "${sharedDir}/components/headers/"
#  )

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

    "${calculationsQueryHistoryDir}/query_history_parameters.yaml"

    # "${boxDir}/box_parameters.yaml"
    # "${laminateDir}/laminate_parameters.yaml"
    # "${paperDir}/paper_parameters.yaml"

    "${companyPageDir}/company_page_parameters.yaml"

    "${submitFormDir}/submit_form_parameters.yaml"

    # "${materialTypeDir}/material_type_parameters.yaml"
    # "${paperColorDir}/paper_color_parameters.yaml"
    # "${paperFactureDir}/paper_facture_parameters.yaml"
    # "${printFormatDir}/print_format_parameters.yaml"

    "${fileStationDir}/parameters.yaml"
  )

  OPENAPI_COMPONENTS_SCHEMAS=(
    # "${sharedDir}/components/enums/App.Enum.Address.HouseType.yaml"
    # "${sharedDir}/components/enums/App.Enum.DeliveryMethod.yaml"
    # "${sharedDir}/components/enums/App.Enum.Gender.yaml"
    # "${sharedDir}/components/enums/App.Enum.Status.yaml"

    "${sharedDir}/components/fields/App.Field.Article.yaml"
    "${sharedDir}/components/fields/App.Field.Boolean.yaml"
    "${sharedDir}/components/fields/App.Field.Caption.yaml"
    "${sharedDir}/components/fields/App.Field.DateTimeCreatedAt.yaml"
    # "${sharedDir}/components/fields/App.Field.DateTimeUpdatedAt.yaml"
    # "${sharedDir}/components/fields/App.Field.Date.yaml"
    # "${sharedDir}/components/fields/App.Field.DateTime.yaml"
    # "${sharedDir}/components/fields/App.Field.Email.yaml"
    # "${sharedDir}/components/fields/App.Field.ExternalURL.yaml"
    # "${sharedDir}/components/fields/App.Field.FileURL.yaml"
    # "${sharedDir}/components/fields/App.Field.Float64.yaml"
    # "${sharedDir}/components/fields/App.Field.GEO.yaml"
    # "${sharedDir}/components/fields/App.Field.ImageURL.yaml"
    # "${sharedDir}/components/fields/App.Field.Int16.yaml"
    # "${sharedDir}/components/fields/App.Field.Int32.yaml"
    # "${sharedDir}/components/fields/App.Field.Int64.yaml"
    "${sharedDir}/components/fields/App.Field.JsonData.yaml"
    # "${sharedDir}/components/fields/App.Field.ListCursor.Cursor.yaml"
    # "${sharedDir}/components/fields/App.Field.ListCursor.HasNext.yaml"
    # "${sharedDir}/components/fields/App.Field.ListPager.Total.yaml"
    # "${sharedDir}/components/fields/App.Field.OrderIndex.yaml"
    "${sharedDir}/components/fields/App.Field.Percent.yaml"
    # "${sharedDir}/components/fields/App.Field.Phone.yaml"
    "${sharedDir}/components/fields/App.Field.RewriteName.yaml"
    "${sharedDir}/components/fields/App.Field.Size2D.yaml"
    # "${sharedDir}/components/fields/App.Field.Size3D.yaml"
    # "${sharedDir}/components/fields/App.Field.TagVersion.yaml"
    # "${sharedDir}/components/fields/App.Field.Timezone.yaml"
    "${sharedDir}/components/fields/App.Field.Uint.yaml"
    # "${sharedDir}/components/fields/App.Field.UUID.yaml"
    # "${sharedDir}/components/fields/App.Field.VariableCamelCase.yaml"

    # "${sharedDir}/components/fields/measures/App.Field.Measure.Centimeter.yaml"
    "${sharedDir}/components/fields/measures/App.Field.Measure.Gram.yaml"
    "${sharedDir}/components/fields/measures/App.Field.Measure.GramPerMeter2.yaml"
    "${sharedDir}/components/fields/measures/App.Field.Measure.Kilogram.yaml"
    # "${sharedDir}/components/fields/measures/App.Field.Measure.KilogramPerMeter2.yaml"
    "${sharedDir}/components/fields/measures/App.Field.Measure.Meter.yaml"
    "${sharedDir}/components/fields/measures/App.Field.Measure.Meter2.yaml"
    "${sharedDir}/components/fields/measures/App.Field.Measure.Meter3.yaml"
    "${sharedDir}/components/fields/measures/App.Field.Measure.Micrometer.yaml"
    # "${sharedDir}/components/fields/measures/App.Field.Measure.Milligram.yaml"
    "${sharedDir}/components/fields/measures/App.Field.Measure.Millimeter.yaml"
    "${sharedDir}/components/fields/measures/App.Field.Measure.Millimeter2.yaml"
    "${sharedDir}/components/fields/measures/App.Field.Measure.Millimeter3.yaml"
    # "${sharedDir}/components/fields/measures/App.Field.Measure.Price.yaml"

    # "${sharedDir}/components/models/App.Request.Model.ChangeFlag.yaml"
    # "${sharedDir}/components/models/App.Request.Model.ChangeStatus.yaml"
    # "${sharedDir}/components/models/App.Request.Model.MoveItem.yaml"
    # "${sharedDir}/components/models/App.Response.Model.BinaryAnyFile.yaml"
    # "${sharedDir}/components/models/App.Response.Model.BinaryFile.yaml"
    "${sharedDir}/components/models/App.Response.Model.BinaryImage.yaml"
    # "${sharedDir}/components/models/App.Response.Model.BinaryMedia.yaml"
    "${sharedDir}/components/models/App.Response.Model.Error.yaml"
    # "${sharedDir}/components/models/App.Response.Model.FileInfo.yaml"
    # "${sharedDir}/components/models/App.Response.Model.JsonFile.yaml"
    # "${sharedDir}/components/models/App.Response.Model.ImageInfo.yaml"
    # "${sharedDir}/components/models/App.Response.Model.Success.yaml"
    "${sharedDir}/components/models/App.Response.Model.SuccessCreatedItem.yaml"
    # "${sharedDir}/components/models/App.Response.Model.SuccessCreatedItemUint.yaml"
    # "${sharedDir}/components/models/App.Response.Model.SuccessSavedItem.yaml"
    # "${sharedDir}/components/models/App.Response.Model.TextFile.yaml"
    # "${sharedDir}/components/models/App.Response.Model.Volume.yaml"

    "${sharedDir}/custom/enums/Custom.Enum.FragmentPosition.yaml"
    "${sharedDir}/custom/enums/Custom.Enum.PaperSides.yaml"

    "${sharedDir}/custom/fields/Custom.Field.Catalog.BoxID.yaml"
    "${sharedDir}/custom/fields/Custom.Field.Catalog.LaminateID.yaml"
    "${sharedDir}/custom/fields/Custom.Field.Catalog.PaperID.yaml"
    "${sharedDir}/custom/fields/Custom.Field.Dictionaries.MaterialTypeID.yaml"
    "${sharedDir}/custom/fields/Custom.Field.Dictionaries.PaperColorID.yaml"
    "${sharedDir}/custom/fields/Custom.Field.Dictionaries.PaperFactureID.yaml"
    "${sharedDir}/custom/fields/Custom.Field.Dictionaries.PrintFormatID.yaml"
    "${sharedDir}/custom/fields/Custom.Field.Fragment.yaml"
    "${sharedDir}/custom/fields/Custom.Field.Layout.yaml"
    "${sharedDir}/custom/fields/Custom.Field.Rect2dFormat.yaml"

    "${calculationsAlgoDir}/box/packinbox/packinbox_schemas.yaml"
    "${calculationsAlgoDir}/sheet/cutting/cutting_schemas.yaml"
    "${calculationsAlgoDir}/sheet/imposition/imposition_schemas.yaml"
    "${calculationsAlgoDir}/sheet/insideoutside/inside_outside_schemas.yaml"
    "${calculationsAlgoDir}/sheet/packinstack/packinstack_schemas.yaml"
    "${calculationsQueryHistoryDir}/query_history_schemas.yaml"

    "${boxDir}/box_schemas.yaml"
    "${laminateDir}/laminate_schemas.yaml"
    "${paperDir}/paper_schemas.yaml"

    "${companyPageDir}/company_page_schemas.yaml"

    "${materialTypeDir}/material_type_schemas.yaml"
    "${paperColorDir}/paper_color_schemas.yaml"
    "${paperFactureDir}/paper_facture_schemas.yaml"
    "${printFormatDir}/print_format_schemas.yaml"

    "${submitFormDir}/submit_form_schemas.yaml"

    "${fileStationDir}/schemas.yaml"
  )

  OPENAPI_COMPONENTS_RESPONSES=(
    "${sharedDir}/components/responses/App.ResponseJson.Errors.yaml"
    "${sharedDir}/components/responses/App.ResponseJson.ErrorsAuth.yaml"
  )

#  OPENAPI_SECURITY_SCHEMES=(
#    "${sharedDir}/securitySchemes.yaml"
#  )
}
