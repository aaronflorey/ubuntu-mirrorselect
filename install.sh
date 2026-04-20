#!/usr/bin/env bash
# Written in [Amber](https://amber-lang.com/)
# version: 0.5.1-alpha
# We cannot import `bash_version` from `env.ab` because it imports `text.ab` making a circular dependency.
# This is a workaround to avoid that issue and the import system should be improved in the future.
bash_version__0_v0() {
    major_11="$(echo "${BASH_VERSINFO[0]}")"
    minor_12="$(echo "${BASH_VERSINFO[1]}")"
    command_2="$(echo "${BASH_VERSINFO[2]}")"
    __status=$?
    patch_13="${command_2}"
    ret_bash_version0_v0=("${major_11}" "${minor_12}" "${patch_13}")
    return 0
}

replace__1_v0() {
    local source=$1
    local search=$2
    local replace=$3
    # Here we use a command to avoid #646
    result_64=""
    bash_version__0_v0 
    left_comp=("${ret_bash_version0_v0[@]}")
    right_comp=(4 3)
    comp="$(
        # Compare if left array >= right array
        len_comp="$( (( "${#left_comp[@]}" < "${#right_comp[@]}" )) && echo "${#left_comp[@]}"|| echo "${#right_comp[@]}")"
        for (( i=0; i<len_comp; i++ )); do
            left="${left_comp[i]:-0}"
            right="${right_comp[i]:-0}"
            if (( "${left}" > "${right}" )); then
                echo 1
                exit
            elif (( "${left}" < "${right}" )); then
                echo 0
                exit
            fi
        done
        (( "${#left_comp[@]}" == "${#right_comp[@]}" || "${#left_comp[@]}" > "${#right_comp[@]}" )) && echo 1 || echo 0
)"
    if [ "${comp}" != 0 ]; then
        result_64="${source//"${search}"/"${replace}"}"
        __status=$?
    else
        result_64="${source//"${search}"/${replace}}"
        __status=$?
    fi
    ret_replace1_v0="${result_64}"
    return 0
}

replace_one__2_v0() {
    local source=$1
    local search=$2
    local replace=$3
    # Here we use a command to avoid #646
    result_10=""
    bash_version__0_v0 
    left_comp=("${ret_bash_version0_v0[@]}")
    right_comp=(4 3)
    comp="$(
        # Compare if left array >= right array
        len_comp="$( (( "${#left_comp[@]}" < "${#right_comp[@]}" )) && echo "${#left_comp[@]}"|| echo "${#right_comp[@]}")"
        for (( i=0; i<len_comp; i++ )); do
            left="${left_comp[i]:-0}"
            right="${right_comp[i]:-0}"
            if (( "${left}" > "${right}" )); then
                echo 1
                exit
            elif (( "${left}" < "${right}" )); then
                echo 0
                exit
            fi
        done
        (( "${#left_comp[@]}" == "${#right_comp[@]}" || "${#left_comp[@]}" > "${#right_comp[@]}" )) && echo 1 || echo 0
)"
    if [ "${comp}" != 0 ]; then
        result_10="${source/"${search}"/"${replace}"}"
        __status=$?
    else
        result_10="${source/"${search}"/${replace}}"
        __status=$?
    fi
    ret_replace_one2_v0="${result_10}"
    return 0
}

__SED_VERSION_UNKNOWN_0=0
__SED_VERSION_GNU_1=1
__SED_VERSION_BUSYBOX_2=2
sed_version__3_v0() {
    # We can't match against a word "GNU" because
    # alpine's busybox sed returns "This is not GNU sed version"
    re='\bCopyright\b.+\bFree Software Foundation\b'; [[ $(sed --version 2>/dev/null) =~ $re ]]
    __status=$?
    if [ "$(( ${__status} == 0 ))" != 0 ]; then
        ret_sed_version3_v0="${__SED_VERSION_GNU_1}"
        return 0
    fi
    # On BSD single `sed` waits for stdin. We must use `sed --help` to avoid this.
    re='\bBusyBox\b'; [[ $(sed --help 2>&1) =~ $re ]]
    __status=$?
    if [ "$(( ${__status} == 0 ))" != 0 ]; then
        ret_sed_version3_v0="${__SED_VERSION_BUSYBOX_2}"
        return 0
    fi
    ret_sed_version3_v0="${__SED_VERSION_UNKNOWN_0}"
    return 0
}

split__5_v0() {
    local text=$1
    local delimiter=$2
    result_15=()
    IFS="${delimiter}" read -rd '' -a result_15 < <(printf %s "$text")
    __status=$?
    ret_split5_v0=("${result_15[@]}")
    return 0
}

trim_left__9_v0() {
    local text=$1
    command_7="$(echo "${text}" | sed -e 's/^[[:space:]]*//')"
    __status=$?
    ret_trim_left9_v0="${command_7}"
    return 0
}

trim_right__10_v0() {
    local text=$1
    command_8="$(echo "${text}" | sed -e 's/[[:space:]]*$//')"
    __status=$?
    ret_trim_right10_v0="${command_8}"
    return 0
}

trim__11_v0() {
    local text=$1
    trim_right__10_v0 "${text}"
    ret_trim_right10_v0__178_22="${ret_trim_right10_v0}"
    trim_left__9_v0 "${ret_trim_right10_v0__178_22}"
    ret_trim11_v0="${ret_trim_left9_v0}"
    return 0
}

lowercase__12_v0() {
    local text=$1
    command_9="$(echo "${text}" | tr '[:upper:]' '[:lower:]')"
    __status=$?
    ret_lowercase12_v0="${command_9}"
    return 0
}

parse_int__14_v0() {
    local text=$1
    [ -n "${text}" ] && [ "${text}" -eq "${text}" ] 2>/dev/null
    __status=$?
    if [ "${__status}" != 0 ]; then
        ret_parse_int14_v0=''
        return "${__status}"
    fi
    ret_parse_int14_v0="${text}"
    return 0
}

text_contains__17_v0() {
    local source=$1
    local search=$2
    command_10="$(if [[ "${source}" == *"${search}"* ]]; then
    echo 1
  fi)"
    __status=$?
    result_47="${command_10}"
    ret_text_contains17_v0="$([ "_${result_47}" != "_1" ]; echo $?)"
    return 0
}

match_regex__20_v0() {
    local source=$1
    local search=$2
    local extended=$3
    sed_version__3_v0 
    sed_version_63="${ret_sed_version3_v0}"
    replace__1_v0 "${search}" "/" "\\/"
    search="${ret_replace1_v0}"
    output_65=""
    if [ "$(( $(( ${sed_version_63} == ${__SED_VERSION_GNU_1} )) || $(( ${sed_version_63} == ${__SED_VERSION_BUSYBOX_2} )) ))" != 0 ]; then
        # '\b' is supported but not in POSIX standards. Disable it
        replace__1_v0 "${search}" "\\b" "\\\\b"
        search="${ret_replace1_v0}"
    fi
    if [ "${extended}" != 0 ]; then
        # GNU sed versions 4.0 through 4.2 support extended regex syntax,
        # but only via the "-r" option
        if [ "$(( ${sed_version_63} == ${__SED_VERSION_GNU_1} ))" != 0 ]; then
            # '\b' is not in POSIX standards. Disable it
            replace__1_v0 "${search}" "\\b" "\\b"
            search="${ret_replace1_v0}"
            command_11="$(echo "${source}" | sed -r -ne "/${search}/p")"
            __status=$?
            output_65="${command_11}"
        else
            command_12="$(echo "${source}" | sed -E -ne "/${search}/p")"
            __status=$?
            output_65="${command_12}"
        fi
    else
        if [ "$(( $(( ${sed_version_63} == ${__SED_VERSION_GNU_1} )) || $(( ${sed_version_63} == ${__SED_VERSION_BUSYBOX_2} )) ))" != 0 ]; then
            # GNU Sed BRE handle \| as a metacharacter, but it is not POSIX standands. Disable it
            replace__1_v0 "${search}" "\\|" "|"
            search="${ret_replace1_v0}"
        fi
        command_13="$(echo "${source}" | sed -ne "/${search}/p")"
        __status=$?
        output_65="${command_13}"
    fi
    if [ "$([ "_${output_65}" == "_" ]; echo $?)" != 0 ]; then
        ret_match_regex20_v0=1
        return 0
    fi
    ret_match_regex20_v0=0
    return 0
}

starts_with__23_v0() {
    local text=$1
    local prefix=$2
    command_14="$(if [[ "${text}" == "${prefix}"* ]]; then
    echo 1
  fi)"
    __status=$?
    result_9="${command_14}"
    ret_starts_with23_v0="$([ "_${result_9}" != "_1" ]; echo $?)"
    return 0
}

ends_with__24_v0() {
    local text=$1
    local suffix=$2
    command_15="$(if [[ "${text}" == *"${suffix}" ]]; then
    echo 1
  fi)"
    __status=$?
    result_7="${command_15}"
    ret_ends_with24_v0="$([ "_${result_7}" != "_1" ]; echo $?)"
    return 0
}

slice__25_v0() {
    local text=$1
    local index=$2
    local length=$3
    if [ "$(( ${length} == 0 ))" != 0 ]; then
        __length_16="${text}"
        length="$(( ${#__length_16} - ${index} ))"
    fi
    if [ "$(( ${length} <= 0 ))" != 0 ]; then
        ret_slice25_v0=""
        return 0
    fi
    command_17="$(printf "%.${length}s" "${text: ${index}}")"
    __status=$?
    ret_slice25_v0="${command_17}"
    return 0
}

dir_exists__36_v0() {
    local path=$1
    [ -d "${path}" ]
    __status=$?
    ret_dir_exists36_v0="$(( ${__status} == 0 ))"
    return 0
}

file_exists__37_v0() {
    local path=$1
    [ -f "${path}" ]
    __status=$?
    ret_file_exists37_v0="$(( ${__status} == 0 ))"
    return 0
}

dir_create__42_v0() {
    local path=$1
    dir_exists__36_v0 "${path}"
    ret_dir_exists36_v0__87_12="${ret_dir_exists36_v0}"
    if [ "$(( ! ${ret_dir_exists36_v0__87_12} ))" != 0 ]; then
        mkdir -p "${path}"
        __status=$?
        if [ "${__status}" != 0 ]; then
            ret_dir_create42_v0=''
            return "${__status}"
        fi
    fi
}

is_mac_os_mktemp__43_v0() {
    # macOS's mktemp does not have --version
    mktemp --version >/dev/null 2>&1
    __status=$?
    if [ "${__status}" != 0 ]; then
        ret_is_mac_os_mktemp43_v0=1
        return 0
    fi
    ret_is_mac_os_mktemp43_v0=0
    return 0
}

temp_dir_create__44_v0() {
    local template=$1
    local auto_delete=$2
    local force_delete=$3
    trim__11_v0 "${template}"
    ret_trim11_v0__113_8="${ret_trim11_v0}"
    if [ "$([ "_${ret_trim11_v0__113_8}" != "_" ]; echo $?)" != 0 ]; then
        echo "The template cannot be an empty string"'!'""
        ret_temp_dir_create44_v0=''
        return 1
    fi
    filename_60=""
    is_mac_os_mktemp__43_v0 
    ret_is_mac_os_mktemp43_v0__119_8="${ret_is_mac_os_mktemp43_v0}"
    if [ "${ret_is_mac_os_mktemp43_v0__119_8}" != 0 ]; then
        # usage: mktemp [-d] [-p tmpdir] [-q] [-t prefix] [-u] template ...
        # mktemp [-d] [-p tmpdir] [-q] [-u] -t prefix
        command_18="$(mktemp -d -p "$TMPDIR" "${template}")"
        __status=$?
        if [ "${__status}" != 0 ]; then
            ret_temp_dir_create44_v0=''
            return "${__status}"
        fi
        filename_60="${command_18}"
    else
        command_19="$(mktemp -d -p "$TMPDIR" -t "${template}")"
        __status=$?
        if [ "${__status}" != 0 ]; then
            ret_temp_dir_create44_v0=''
            return "${__status}"
        fi
        filename_60="${command_19}"
    fi
    if [ "$([ "_${filename_60}" != "_" ]; echo $?)" != 0 ]; then
        echo "Failed to make a temporary directory"
        ret_temp_dir_create44_v0=''
        return 1
    fi
    if [ "${auto_delete}" != 0 ]; then
        if [ "${force_delete}" != 0 ]; then
            trap 'rm -rf '"${filename_60}"'' EXIT
            __status=$?
            if [ "${__status}" != 0 ]; then
                echo "Setting auto deletion fails. You must delete temporary dir ${filename_60}."
            fi
        else
            trap 'rmdir '"${filename_60}"'' EXIT
            __status=$?
            if [ "${__status}" != 0 ]; then
                echo "Setting auto deletion fails. You must delete temporary dir ${filename_60}."
            fi
        fi
    fi
    ret_temp_dir_create44_v0="${filename_60}"
    return 0
}

file_chmod__45_v0() {
    local path=$1
    local mode=$2
    file_exists__37_v0 "${path}"
    ret_file_exists37_v0__153_8="${ret_file_exists37_v0}"
    if [ "${ret_file_exists37_v0__153_8}" != 0 ]; then
        chmod "${mode}" "${path}"
        __status=$?
        if [ "${__status}" != 0 ]; then
            ret_file_chmod45_v0=''
            return "${__status}"
        fi
        ret_file_chmod45_v0=''
        return 0
    fi
    echo "The file ${path} doesn't exist"'!'""
    ret_file_chmod45_v0=''
    return 1
}

file_extract__50_v0() {
    local path=$1
    local target=$2
    file_exists__37_v0 "${path}"
    ret_file_exists37_v0__229_8="${ret_file_exists37_v0}"
    if [ "${ret_file_exists37_v0__229_8}" != 0 ]; then
        match_regex__20_v0 "${path}" "\\.(tar\\.bz2|tbz|tbz2)\$" 1
        ret_match_regex20_v0__231_13="${ret_match_regex20_v0}"
        match_regex__20_v0 "${path}" "\\.(tar\\.gz|tgz)\$" 1
        ret_match_regex20_v0__232_13="${ret_match_regex20_v0}"
        match_regex__20_v0 "${path}" "\\.(tar\\.xz|txz)\$" 1
        ret_match_regex20_v0__233_13="${ret_match_regex20_v0}"
        match_regex__20_v0 "${path}" "\\.bz2\$" 0
        ret_match_regex20_v0__234_13="${ret_match_regex20_v0}"
        match_regex__20_v0 "${path}" "\\.deb\$" 0
        ret_match_regex20_v0__235_13="${ret_match_regex20_v0}"
        match_regex__20_v0 "${path}" "\\.gz\$" 0
        ret_match_regex20_v0__236_13="${ret_match_regex20_v0}"
        match_regex__20_v0 "${path}" "\\.rar\$" 0
        ret_match_regex20_v0__237_13="${ret_match_regex20_v0}"
        match_regex__20_v0 "${path}" "\\.rpm\$" 0
        ret_match_regex20_v0__238_13="${ret_match_regex20_v0}"
        match_regex__20_v0 "${path}" "\\.tar\$" 0
        ret_match_regex20_v0__239_13="${ret_match_regex20_v0}"
        match_regex__20_v0 "${path}" "\\.xz\$" 0
        ret_match_regex20_v0__240_13="${ret_match_regex20_v0}"
        match_regex__20_v0 "${path}" "\\.7z\$" 0
        ret_match_regex20_v0__241_13="${ret_match_regex20_v0}"
        match_regex__20_v0 "${path}" "\\.\\(zip\\|war\\|jar\\)\$" 0
        ret_match_regex20_v0__242_13="${ret_match_regex20_v0}"
        if [ "${ret_match_regex20_v0__231_13}" != 0 ]; then
            tar xvjf "${path}" -C "${target}"
            __status=$?
            if [ "${__status}" != 0 ]; then
                ret_file_extract50_v0=''
                return "${__status}"
            fi
        elif [ "${ret_match_regex20_v0__232_13}" != 0 ]; then
            tar xzf "${path}" -C "${target}"
            __status=$?
            if [ "${__status}" != 0 ]; then
                ret_file_extract50_v0=''
                return "${__status}"
            fi
        elif [ "${ret_match_regex20_v0__233_13}" != 0 ]; then
            tar xJf "${path}" -C "${target}"
            __status=$?
            if [ "${__status}" != 0 ]; then
                ret_file_extract50_v0=''
                return "${__status}"
            fi
        elif [ "${ret_match_regex20_v0__234_13}" != 0 ]; then
            bunzip2 "${path}"
            __status=$?
            if [ "${__status}" != 0 ]; then
                ret_file_extract50_v0=''
                return "${__status}"
            fi
        elif [ "${ret_match_regex20_v0__235_13}" != 0 ]; then
            dpkg-deb -xv "${path}" "${target}"
            __status=$?
            if [ "${__status}" != 0 ]; then
                ret_file_extract50_v0=''
                return "${__status}"
            fi
        elif [ "${ret_match_regex20_v0__236_13}" != 0 ]; then
            gunzip "${path}"
            __status=$?
            if [ "${__status}" != 0 ]; then
                ret_file_extract50_v0=''
                return "${__status}"
            fi
        elif [ "${ret_match_regex20_v0__237_13}" != 0 ]; then
            unrar x "${path}" "${target}"
            __status=$?
            if [ "${__status}" != 0 ]; then
                ret_file_extract50_v0=''
                return "${__status}"
            fi
        elif [ "${ret_match_regex20_v0__238_13}" != 0 ]; then
            rpm2cpio "${path}" | cpio -idm
            __status=$?
            if [ "${__status}" != 0 ]; then
                ret_file_extract50_v0=''
                return "${__status}"
            fi
        elif [ "${ret_match_regex20_v0__239_13}" != 0 ]; then
            tar xf "${path}" -C "${target}"
            __status=$?
            if [ "${__status}" != 0 ]; then
                ret_file_extract50_v0=''
                return "${__status}"
            fi
        elif [ "${ret_match_regex20_v0__240_13}" != 0 ]; then
            xz --decompress "${path}"
            __status=$?
            if [ "${__status}" != 0 ]; then
                ret_file_extract50_v0=''
                return "${__status}"
            fi
        elif [ "${ret_match_regex20_v0__241_13}" != 0 ]; then
            7z -y "${path}" -o "${target}"
            __status=$?
            if [ "${__status}" != 0 ]; then
                ret_file_extract50_v0=''
                return "${__status}"
            fi
        elif [ "${ret_match_regex20_v0__242_13}" != 0 ]; then
            unzip "${path}" -d "${target}"
            __status=$?
            if [ "${__status}" != 0 ]; then
                ret_file_extract50_v0=''
                return "${__status}"
            fi
        else
            echo "Error: Unsupported file type"
            ret_file_extract50_v0=''
            return 3
        fi
    else
        echo "Error: File not found"
        ret_file_extract50_v0=''
        return 2
    fi
}

env_var_get__98_v0() {
    local name=$1
    command_20="$(echo ${!name})"
    __status=$?
    if [ "${__status}" != 0 ]; then
        ret_env_var_get98_v0=''
        return "${__status}"
    fi
    ret_env_var_get98_v0="${command_20}"
    return 0
}

is_command__100_v0() {
    local command=$1
    [ -x "$(command -v "${command}")" ]
    __status=$?
    if [ "${__status}" != 0 ]; then
        ret_is_command100_v0=0
        return 0
    fi
    ret_is_command100_v0=1
    return 0
}

file_download__151_v0() {
    local url=$1
    local path=$2
    is_command__100_v0 "curl"
    ret_is_command100_v0__14_9="${ret_is_command100_v0}"
    is_command__100_v0 "wget"
    ret_is_command100_v0__17_9="${ret_is_command100_v0}"
    is_command__100_v0 "aria2c"
    ret_is_command100_v0__20_9="${ret_is_command100_v0}"
    if [ "${ret_is_command100_v0__14_9}" != 0 ]; then
        curl -L -o "${path}" "${url}" >/dev/null 2>&1
        __status=$?
    elif [ "${ret_is_command100_v0__17_9}" != 0 ]; then
        wget "${url}" -P "${path}" >/dev/null 2>&1
        __status=$?
    elif [ "${ret_is_command100_v0__20_9}" != 0 ]; then
        aria2c "${url}" -d "${path}" >/dev/null 2>&1
        __status=$?
    else
        ret_file_download151_v0=''
        return 1
    fi
}

# #!/usr/bin/env amber
__DEFAULT_HOST_3="github.com"
__DEFAULT_SLUG_4="aaronflorey/ubuntu-mirrorselect"
__BINARY_NAME_5="mirrorselect"
default_repo__162_v0() {
    ret_default_repo162_v0=("${__DEFAULT_HOST_3}" "${__DEFAULT_SLUG_4}")
    return 0
}

strip_git_suffix__163_v0() {
    local value=$1
    ends_with__24_v0 "${value}" ".git"
    ret_ends_with24_v0__17_8="${ret_ends_with24_v0}"
    if [ "${ret_ends_with24_v0__17_8}" != 0 ]; then
        __length_22="${value}"
        slice__25_v0 "${value}" 0 "$(( ${#__length_22} - 4 ))"
        ret_strip_git_suffix163_v0="${ret_slice25_v0}"
        return 0
    fi
    ret_strip_git_suffix163_v0="${value}"
    return 0
}

strip_port__164_v0() {
    local authority=$1
    starts_with__23_v0 "${authority}" "["
    ret_starts_with23_v0__25_8="${ret_starts_with23_v0}"
    if [ "${ret_starts_with23_v0__25_8}" != 0 ]; then
        ret_strip_port164_v0="${authority}"
        return 0
    fi
    split__5_v0 "${authority}" ":"
    parts_20=("${ret_split5_v0[@]}")
    __length_23=("${parts_20[@]}")
    if [ "$(( ${#__length_23[@]} != 2 ))" != 0 ]; then
        ret_strip_port164_v0="${authority}"
        return 0
    fi
    parse_int__14_v0 "${parts_20[1]}"
    __status=$?
    if [ "${__status}" != 0 ]; then
        ret_strip_port164_v0="${authority}"
        return 0
    fi
    ret_strip_port164_v0="${parts_20[0]}"
    return 0
}

join_from_index__165_v0() {
    local parts=("${!1}")
    local start=$2
    result_22=""
    index_24=0;
    for part_23 in "${parts[@]}"; do
        if [ "$(( ${index_24} < ${start} ))" != 0 ]; then
            continue
        fi
        if [ "$([ "_${result_22}" != "_" ]; echo $?)" != 0 ]; then
            result_22="${part_23}"
            continue
        fi
        result_22+="/${part_23}"
        (( index_24++ )) || true
    done
    ret_join_from_index165_v0="${result_22}"
    return 0
}

parse_remote__166_v0() {
    local remote=$1
    trim__11_v0 "${remote}"
    ret_trim11_v0__61_36="${ret_trim11_v0}"
    strip_git_suffix__163_v0 "${ret_trim11_v0__61_36}"
    cleaned_8="${ret_strip_git_suffix163_v0}"
    starts_with__23_v0 "${cleaned_8}" "git@"
    ret_starts_with23_v0__63_8="${ret_starts_with23_v0}"
    if [ "${ret_starts_with23_v0__63_8}" != 0 ]; then
        replace_one__2_v0 "${cleaned_8}" "git@" ""
        without_user_14="${ret_replace_one2_v0}"
        split__5_v0 "${without_user_14}" ":"
        parts_16=("${ret_split5_v0[@]}")
        __length_24=("${parts_16[@]}")
        if [ "$(( ${#__length_24[@]} != 2 ))" != 0 ]; then
            ret_parse_remote166_v0=()
            return 1
        fi
        ret_parse_remote166_v0=("${parts_16[0]}" "${parts_16[1]}")
        return 0
    fi
    starts_with__23_v0 "${cleaned_8}" "ssh://"
    ret_starts_with23_v0__73_8="${ret_starts_with23_v0}"
    if [ "${ret_starts_with23_v0__73_8}" != 0 ]; then
        replace_one__2_v0 "${cleaned_8}" "ssh://" ""
        without_scheme_17="${ret_replace_one2_v0}"
        split__5_v0 "${without_scheme_17}" "/"
        path_parts_18=("${ret_split5_v0[@]}")
        __length_26=("${path_parts_18[@]}")
        if [ "$(( ${#__length_26[@]} < 3 ))" != 0 ]; then
            ret_parse_remote166_v0=()
            return 1
        fi
        split__5_v0 "${path_parts_18[0]}" "@"
        authority_parts_19=("${ret_split5_v0[@]}")
        __length_27=("${authority_parts_19[@]}")
        strip_port__164_v0 "${authority_parts_19[$(( ${#__length_27[@]} - 1 ))]}"
        host_21="${ret_strip_port164_v0}"
        join_from_index__165_v0 path_parts_18[@] 1
        slug_25="${ret_join_from_index165_v0}"
        ret_parse_remote166_v0=("${host_21}" "${slug_25}")
        return 0
    fi
    starts_with__23_v0 "${cleaned_8}" "https://"
    ret_starts_with23_v0__86_8="${ret_starts_with23_v0}"
    starts_with__23_v0 "${cleaned_8}" "http://"
    ret_starts_with23_v0__86_44="${ret_starts_with23_v0}"
    if [ "$(( ${ret_starts_with23_v0__86_8} || ${ret_starts_with23_v0__86_44} ))" != 0 ]; then
        starts_with__23_v0 "${cleaned_8}" "https://"
        ret_starts_with23_v0__87_30="${ret_starts_with23_v0}"
        replace_one__2_v0 "${cleaned_8}" "https://" ""
        ret_replace_one2_v0__88_18="${ret_replace_one2_v0}"
        replace_one__2_v0 "${cleaned_8}" "http://" ""
        ret_replace_one2_v0__89_18="${ret_replace_one2_v0}"
        without_scheme_26="$(if [ "${ret_starts_with23_v0__87_30}" != 0 ]; then echo "${ret_replace_one2_v0__88_18}"; else echo "${ret_replace_one2_v0__89_18}"; fi)"
        split__5_v0 "${without_scheme_26}" "/"
        path_parts_27=("${ret_split5_v0[@]}")
        __length_29=("${path_parts_27[@]}")
        if [ "$(( ${#__length_29[@]} < 3 ))" != 0 ]; then
            ret_parse_remote166_v0=()
            return 1
        fi
        join_from_index__165_v0 path_parts_27[@] 1
        slug_28="${ret_join_from_index165_v0}"
        ret_parse_remote166_v0=("${path_parts_27[0]}" "${slug_28}")
        return 0
    fi
    ret_parse_remote166_v0=()
    return 1
}

detect_repo__167_v0() {
    is_command__100_v0 "git"
    ret_is_command100_v0__103_12="${ret_is_command100_v0}"
    if [ "$(( ! ${ret_is_command100_v0__103_12} ))" != 0 ]; then
        default_repo__162_v0 
        ret_detect_repo167_v0=("${ret_default_repo162_v0[@]}")
        return 0
    fi
    command_32="$(git remote get-url origin)"
    __status=$?
    if [ "${__status}" != 0 ]; then
        command_31="$(git remote get-url upstream)"
        __status=$?
        if [ "${__status}" != 0 ]; then
            default_repo__162_v0 
            ret_detect_repo167_v0=("${ret_default_repo162_v0[@]}")
            return 0
        fi
        upstream_6="${command_31}"
        parse_remote__166_v0 "${upstream_6}"
        __status=$?
        if [ "${__status}" != 0 ]; then
            default_repo__162_v0 
            ret_detect_repo167_v0=("${ret_default_repo162_v0[@]}")
            return 0
        fi
        ret_detect_repo167_v0=("${ret_parse_remote166_v0[@]}")
        return 0
    fi
    remote_29="${command_32}"
    parse_remote__166_v0 "${remote_29}"
    __status=$?
    if [ "${__status}" != 0 ]; then
        default_repo__162_v0 
        ret_detect_repo167_v0=("${ret_default_repo162_v0[@]}")
        return 0
    fi
    ret_detect_repo167_v0=("${ret_parse_remote166_v0[@]}")
    return 0
}

normalize_os__168_v0() {
    local raw=$1
    trim__11_v0 "${raw}"
    ret_trim11_v0__123_27="${ret_trim11_v0}"
    lowercase__12_v0 "${ret_trim11_v0__123_27}"
    value_35="${ret_lowercase12_v0}"
    starts_with__23_v0 "${value_35}" "msys"
    ret_starts_with23_v0__128_9="${ret_starts_with23_v0}"
    starts_with__23_v0 "${value_35}" "mingw"
    ret_starts_with23_v0__128_39="${ret_starts_with23_v0}"
    starts_with__23_v0 "${value_35}" "cygwin"
    ret_starts_with23_v0__128_70="${ret_starts_with23_v0}"
    if [ "$([ "_${value_35}" != "_linux" ]; echo $?)" != 0 ]; then
        ret_normalize_os168_v0="linux"
        return 0
    elif [ "$([ "_${value_35}" != "_darwin" ]; echo $?)" != 0 ]; then
        ret_normalize_os168_v0="darwin"
        return 0
    elif [ "$(( $(( ${ret_starts_with23_v0__128_9} || ${ret_starts_with23_v0__128_39} )) || ${ret_starts_with23_v0__128_70} ))" != 0 ]; then
        ret_normalize_os168_v0="windows"
        return 0
    else
        ret_normalize_os168_v0=''
        return 1
    fi
}

normalize_arch__169_v0() {
    local raw=$1
    trim__11_v0 "${raw}"
    ret_trim11_v0__134_27="${ret_trim11_v0}"
    lowercase__12_v0 "${ret_trim11_v0__134_27}"
    value_37="${ret_lowercase12_v0}"
    starts_with__23_v0 "${value_37}" "arm64"
    ret_starts_with23_v0__141_9="${ret_starts_with23_v0}"
    if [ "$([ "_${value_37}" != "_x86_64" ]; echo $?)" != 0 ]; then
        ret_normalize_arch169_v0="amd64"
        return 0
    elif [ "$([ "_${value_37}" != "_amd64" ]; echo $?)" != 0 ]; then
        ret_normalize_arch169_v0="amd64"
        return 0
    elif [ "$([ "_${value_37}" != "_aarch64" ]; echo $?)" != 0 ]; then
        ret_normalize_arch169_v0="arm64"
        return 0
    elif [ "$([ "_${value_37}" != "_arm64" ]; echo $?)" != 0 ]; then
        ret_normalize_arch169_v0="arm64"
        return 0
    elif [ "${ret_starts_with23_v0__141_9}" != 0 ]; then
        ret_normalize_arch169_v0="arm64"
        return 0
    else
        ret_normalize_arch169_v0=''
        return 1
    fi
}

detect_target__170_v0() {
    command_33="$(uname -s)"
    __status=$?
    if [ "${__status}" != 0 ]; then
        ret_detect_target170_v0=()
        return "${__status}"
    fi
    raw_os_33="${command_33}"
    command_34="$(uname -m)"
    __status=$?
    if [ "${__status}" != 0 ]; then
        ret_detect_target170_v0=()
        return "${__status}"
    fi
    raw_arch_34="${command_34}"
    normalize_os__168_v0 "${raw_os_33}"
    __status=$?
    if [ "${__status}" != 0 ]; then
        trim__11_v0 "${raw_os_33}"
        ret_trim11_v0__151_46="${ret_trim11_v0}"
        echo "Unsupported operating system: ${ret_trim11_v0__151_46}"
        ret_detect_target170_v0=()
        return 1
    fi
    os_36="${ret_normalize_os168_v0}"
    if [ "$([ "_${os_36}" != "_windows" ]; echo $?)" != 0 ]; then
        echo "Windows hosts are not supported by this installer because it installs into ~/.local/bin"
        ret_detect_target170_v0=()
        return 1
    fi
    normalize_arch__169_v0 "${raw_arch_34}"
    __status=$?
    if [ "${__status}" != 0 ]; then
        trim__11_v0 "${raw_arch_34}"
        ret_trim11_v0__161_42="${ret_trim11_v0}"
        echo "Unsupported architecture: ${ret_trim11_v0__161_42}"
        ret_detect_target170_v0=()
        return 1
    fi
    arch_38="${ret_normalize_arch169_v0}"
    ret_detect_target170_v0=("${os_36}" "${arch_38}")
    return 0
}

release_base_url__171_v0() {
    local host=$1
    local slug=$2
    env_var_get__98_v0 "MIRRORSELECT_INSTALLER_BASE_URL"
    __status=$?
    if [ "${__status}" != 0 ]; then
        ret_release_base_url171_v0="https://${host}/${slug}"
        return 0
    fi
    ret_release_base_url171_v0="${ret_env_var_get98_v0}"
    return 0
}

resolve_latest_release_url__172_v0() {
    local host=$1
    local slug=$2
    is_command__100_v0 "python3"
    ret_is_command100_v0__175_12="${ret_is_command100_v0}"
    if [ "$(( ! ${ret_is_command100_v0__175_12} ))" != 0 ]; then
        echo "Need python3 to resolve the latest release URL"
        ret_resolve_latest_release_url172_v0=''
        return 1
    fi
    release_base_url__171_v0 "${host}" "${slug}"
    base_url_42="${ret_release_base_url171_v0}"
    release_url_43="${base_url_42}/releases/latest"
    command_36="$(python3 -c "import sys, urllib.request as u; req = u.Request(sys.argv[1], headers=dict([('User-Agent', 'mirrorselect-installer')])); response = u.urlopen(req); print(response.geturl()); response.close()" "${release_url_43}")"
    __status=$?
    if [ "${__status}" != 0 ]; then
        ret_resolve_latest_release_url172_v0=''
        return "${__status}"
    fi
    resolved_44="${command_36}"
    ret_resolve_latest_release_url172_v0="${resolved_44}"
    return 0
}

extract_tag__173_v0() {
    local latest_url=$1
    trim__11_v0 "${latest_url}"
    cleaned_46="${ret_trim11_v0}"
    ends_with__24_v0 "${cleaned_46}" "/releases"
    ret_ends_with24_v0__188_8="${ret_ends_with24_v0}"
    ends_with__24_v0 "${cleaned_46}" "/releases/latest"
    ret_ends_with24_v0__188_43="${ret_ends_with24_v0}"
    if [ "$(( ${ret_ends_with24_v0__188_8} || ${ret_ends_with24_v0__188_43} ))" != 0 ]; then
        ret_extract_tag173_v0=''
        return 2
    fi
    text_contains__17_v0 "${cleaned_46}" "/releases/tag/"
    ret_text_contains17_v0__192_12="${ret_text_contains17_v0}"
    if [ "$(( ! ${ret_text_contains17_v0__192_12} ))" != 0 ]; then
        ret_extract_tag173_v0=''
        return 1
    fi
    split__5_v0 "${cleaned_46}" "/"
    parts_48=("${ret_split5_v0[@]}")
    __length_37=("${parts_48[@]}")
    tag_49="${parts_48[$(( ${#__length_37[@]} - 1 ))]}"
    __length_38=("${parts_48[@]}")
    if [ "$(( $([ "_${tag_49}" != "_" ]; echo $?) && $(( ${#__length_38[@]} > 1 )) ))" != 0 ]; then
        __length_39=("${parts_48[@]}")
        tag_49="${parts_48[$(( ${#__length_39[@]} - 2 ))]}"
    fi
    if [ "$([ "_${tag_49}" != "_" ]; echo $?)" != 0 ]; then
        ret_extract_tag173_v0=''
        return 1
    fi
    ret_extract_tag173_v0="${tag_49}"
    return 0
}

build_asset_name__174_v0() {
    local tag=$1
    local os=$2
    local arch=$3
    starts_with__23_v0 "${tag}" "v"
    ret_starts_with23_v0__210_19="${ret_starts_with23_v0}"
    slice__25_v0 "${tag}" 1 0
    ret_slice25_v0__210_46="${ret_slice25_v0}"
    version_52="$(if [ "${ret_starts_with23_v0__210_19}" != 0 ]; then echo "${ret_slice25_v0__210_46}"; else echo "${tag}"; fi)"
    ret_build_asset_name174_v0="${__BINARY_NAME_5}_${version_52}_${os}_${arch}.tar.gz"
    return 0
}

ensure_bin_dir__175_v0() {
    env_var_get__98_v0 "HOME"
    __status=$?
    if [ "${__status}" != 0 ]; then
        echo "HOME is not set"
        ret_ensure_bin_dir175_v0=''
        return 1
    fi
    home_56="${ret_env_var_get98_v0}"
    bin_dir_57="${home_56}/.local/bin"
    dir_exists__36_v0 "${bin_dir_57}"
    ret_dir_exists36_v0__221_12="${ret_dir_exists36_v0}"
    if [ "$(( ! ${ret_dir_exists36_v0__221_12} ))" != 0 ]; then
        dir_create__42_v0 "${bin_dir_57}"
        __status=$?
        if [ "${__status}" != 0 ]; then
            ret_ensure_bin_dir175_v0=''
            return "${__status}"
        fi
    fi
    ret_ensure_bin_dir175_v0="${bin_dir_57}"
    return 0
}

detect_repo__167_v0 
repo_30=("${ret_detect_repo167_v0[@]}")
host_31="${repo_30[0]}"
slug_32="${repo_30[1]}"
detect_target__170_v0 
__status=$?
if [ "${__status}" != 0 ]; then
    exit "${__status}"
fi
target_39=("${ret_detect_target170_v0[@]}")
os_40="${target_39[0]}"
arch_41="${target_39[1]}"
resolve_latest_release_url__172_v0 "${host_31}" "${slug_32}"
__status=$?
if [ "${__status}" != 0 ]; then
    echo "Failed to resolve the latest release for ${host_31}/${slug_32}"
    exit 1
fi
latest_url_45="${ret_resolve_latest_release_url172_v0}"
extract_tag__173_v0 "${latest_url_45}"
__status=$?
if [ "${__status}" != 0 ]; then
code_50="${__status}"
    if [ "$(( ${code_50} == 2 ))" != 0 ]; then
        echo "No published release was found for ${host_31}/${slug_32}"
        exit 1
    fi
    echo "Resolved latest release URL has an unsupported layout: ${latest_url_45}"
    exit 1
fi
tag_51="${ret_extract_tag173_v0}"
build_asset_name__174_v0 "${tag_51}" "${os_40}" "${arch_41}"
archive_name_53="${ret_build_asset_name174_v0}"
release_base_url__171_v0 "${host_31}" "${slug_32}"
base_url_54="${ret_release_base_url171_v0}"
download_url_55="${base_url_54}/releases/download/${tag_51}/${archive_name_53}"
ensure_bin_dir__175_v0 
__status=$?
if [ "${__status}" != 0 ]; then
    exit "${__status}"
fi
bin_dir_58="${ret_ensure_bin_dir175_v0}"
install_path_59="${bin_dir_58}/${__BINARY_NAME_5}"
temp_dir_create__44_v0 "mirrorselect-release.XXXXXXXXXX" 1 1
__status=$?
if [ "${__status}" != 0 ]; then
    exit "${__status}"
fi
temp_dir_61="${ret_temp_dir_create44_v0}"
archive_path_62="${temp_dir_61}/${archive_name_53}"
echo "Downloading ${download_url_55}"
file_download__151_v0 "${download_url_55}" "${archive_path_62}"
__status=$?
if [ "${__status}" != 0 ]; then
    echo "Failed to download archive: ${archive_name_53}"
    exit 1
fi
file_extract__50_v0 "${archive_path_62}" "${temp_dir_61}"
__status=$?
if [ "${__status}" != 0 ]; then
    exit "${__status}"
fi
extracted_binary_66="${temp_dir_61}/${__BINARY_NAME_5}"
file_exists__37_v0 "${extracted_binary_66}"
ret_file_exists37_v0__269_12="${ret_file_exists37_v0}"
if [ "$(( ! ${ret_file_exists37_v0__269_12} ))" != 0 ]; then
    echo "Archive did not contain ${__BINARY_NAME_5} at the expected path"
    exit 1
fi
cp "${extracted_binary_66}" "${install_path_59}"
__status=$?
if [ "${__status}" != 0 ]; then
    exit "${__status}"
fi
file_chmod__45_v0 "${install_path_59}" "755"
__status=$?
if [ "${__status}" != 0 ]; then
    exit "${__status}"
fi
echo "Installed ${__BINARY_NAME_5} ${tag_51} to ${install_path_59}"
