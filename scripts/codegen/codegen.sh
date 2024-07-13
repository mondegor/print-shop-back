
function mrcmd_plugins_codegen_method_init() {
  readonly CODEGEN_CAPTION="PrintShop code Generator"

  readonly CODEGEN_VARS=(
    "CODEGEN_SRC_DIR"
    "CODEGEN_SRC_SHARED_DIR"
    "CODEGEN_ASSEMBLY_DIR"
  )

  # default values of array: CODEGEN_VARS
  readonly CODEGEN_VARS_DEFAULT=(
    "${APPX_DIR}/codesrc/"
    "${APPX_DIR}/codesrc/_shared"
    "${APPX_DIR}/codesrc/_generated"
  )

  mrcore_dotenv_init_var_array CODEGEN_VARS[@] CODEGEN_VARS_DEFAULT[@]
}

function mrcmd_plugins_codegen_method_config() {
  mrcore_dotenv_echo_var_array CODEGEN_VARS[@]
}

function mrcmd_plugins_codegen_method_export_config() {
  mrcore_dotenv_export_var_array CODEGEN_VARS[@]
}

function mrcmd_plugins_codegen_method_exec() {
  local currentCommand="${1:?}" # script name: adm.sh, public.sh, ...
  local sectionName="" # sample: admin-api, public-api, ...

  case "${currentCommand}" in

    build-all)
      mrcmd codegen build-adm
      ;;

    build-adm)
      sectionName="admin-api"
      ;;

    *)
      ${RETURN_FALSE}
      ;;

  esac

  if [ -n "${sectionName}" ]; then
    mrcore_import "${MRCMD_PLUGINS_DIR}/libs/codegen-lib.sh"
    codegen_lib_build_apidoc \
      "codegen/${currentCommand}" \
      "${CODEGEN_SRC_DIR}" \
      "${CODEGEN_SRC_SHARED_DIR}" \
      "${CODEGEN_ASSEMBLY_DIR}" \
      "${sectionName}" \
      "${CODEGEN_FILENAME_PREFIX}" \
      "${fileNamePostfix}"
  fi
}

function mrcmd_plugins_codegen_method_help() {
  #markup:"|-|-|---------|-------|-------|---------------------------------------|"
  echo -e "${CC_YELLOW}Commands:${CC_END}"
  echo -e "    build-all         Generate all code"
  echo -e "    build-adm-all     Builds all admin code"
}
