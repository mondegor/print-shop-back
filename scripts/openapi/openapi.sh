
function mrcmd_plugins_openapi_method_init() {
  readonly OPENAPI_CAPTION="PrintShop REST API Builder"

  readonly OPENAPI_VARS=(
    "OPENAPI_SRC_DIR"
    "OPENAPI_SRC_SHARED_DIR"
    "OPENAPI_ASSEMBLY_DIR"
    "OPENAPI_FILENAME_PREFIX"
  )

  # default values of array: OPENAPI_VARS
  readonly OPENAPI_VARS_DEFAULT=(
    "${APPX_DIR}/docs/api-src"
    "${APPX_DIR}/docs/api-src/_shared"
    "${APPX_DIR}/docs/api"
    "openapi-"
  )

  mrcore_dotenv_init_var_array OPENAPI_VARS[@] OPENAPI_VARS_DEFAULT[@]
}

function mrcmd_plugins_openapi_method_config() {
  mrcore_dotenv_echo_var_array OPENAPI_VARS[@]
}

function mrcmd_plugins_openapi_method_export_config() {
  mrcore_dotenv_export_var_array OPENAPI_VARS[@]
}

function mrcmd_plugins_openapi_method_exec() {
  local currentCommand="${1:?}" # script name: adm.sh, public.sh, ...
  local sectionName="" # sample: admin-api, public-api, ...
  local fileNamePostfix="" # sample: all, main, ...

  case "${currentCommand}" in

    build-all)
      mrcmd openapi build-adm-all
      mrcmd openapi build-prov-all
      mrcmd openapi build-pub-all
      ;;

    build-adm-all)
      mrcmd openapi build-adm
      mrcmd openapi build-adm-catalog
      mrcmd openapi build-adm-dictionaries
      mrcmd openapi build-adm-controls
      mrcmd openapi build-adm-prv-accounts
      ;;

    build-prov-all)
      mrcmd openapi build-prov
      ;;

    build-pub-all)
      mrcmd openapi build-pub
      mrcmd openapi build-pub-file-station
      ;;

    build-adm)
      sectionName="admin-api"
      ;;

    build-adm-catalog)
      sectionName="admin-api"
      fileNamePostfix="catalog"
      ;;

    build-adm-dictionaries)
      sectionName="admin-api"
      fileNamePostfix="dictionaries"
      ;;

    build-adm-controls)
      sectionName="admin-api"
      fileNamePostfix="controls"
      ;;

    build-adm-prv-accounts)
      sectionName="admin-api"
      fileNamePostfix="provider-accounts"
      ;;

    build-prov)
      sectionName="providers-api"
      ;;

    build-pub)
      sectionName="public-api"
      ;;

    build-pub-file-station)
      sectionName="public-api"
      fileNamePostfix="file-station"
      ;;

    *)
      ${RETURN_FALSE}
      ;;

  esac

  if [ -n "${sectionName}" ]; then
    mrcore_import "${MRCMD_PLUGINS_DIR}/libs/openapi-lib.sh"
    openapi_lib_build_apidoc \
      "openapi/${currentCommand}" \
      "${OPENAPI_SRC_DIR}" \
      "${OPENAPI_SRC_SHARED_DIR}" \
      "${OPENAPI_ASSEMBLY_DIR}" \
      "${sectionName}" \
      "${OPENAPI_FILENAME_PREFIX}" \
      "${fileNamePostfix}"
  fi
}

function mrcmd_plugins_openapi_method_help() {
  #markup:"|-|-|---------|-------|-------|---------------------------------------|"
  echo -e "${CC_YELLOW}Commands:${CC_END}"
  echo -e "    build-all                 Builds all API docs"
  echo -e "    build-adm-all             Builds all admin API docs"
  echo -e "    build-adm                 Builds only full admin API docs"
  echo -e "    build-adm-catalog         Builds admin Catalog API docs"
  echo -e "    build-adm-dictionaries    Builds admin Dictionaries API docs"
  echo -e "    build-adm-prv-accounts    Builds admin Provider accounts API docs"
  echo -e "    build-prov-all            Builds all providers API docs"
  echo -e "    build-pub-all             Builds all public API docs"
}
