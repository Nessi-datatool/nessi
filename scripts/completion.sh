#!/bin/bash
# Bash completion script for Nessi

_nessi_completion() {
    local cur prev opts
    COMPREPLY=()
    cur="${COMP_WORDS[COMP_CWORD]}"
    prev="${COMP_WORDS[COMP_CWORD-1]}"

    # Main commands
    commands="scan report validate schema health optimize transactions view license help version"

    # Flags
    global_flags="--help --version --config --verbose --interactive --dry-run"

    # Handle command-specific options
    case "${prev}" in
        nessi)
            COMPREPLY=( $(compgen -W "${commands}" -- ${cur}) )
            return 0
            ;;
        scan)
            local scan_opts="--format --output --verbose"
            COMPREPLY=( $(compgen -W "${scan_opts}" -- ${cur}) )
            return 0
            ;;
        report)
            local report_opts="--format --output --rules --template"
            COMPREPLY=( $(compgen -W "${report_opts}" -- ${cur}) )
            return 0
            ;;
        validate)
            local validate_opts="--rules --threshold --output"
            COMPREPLY=( $(compgen -W "${validate_opts}" -- ${cur}) )
            return 0
            ;;
        schema)
            local schema_opts="history create compare"
            COMPREPLY=( $(compgen -W "${schema_opts}" -- ${cur}) )
            return 0
            ;;
        view)
            local view_opts="--version --timestamp --format"
            COMPREPLY=( $(compgen -W "${view_opts}" -- ${cur}) )
            return 0
            ;;
        license)
            local license_opts="start-trial status"
            COMPREPLY=( $(compgen -W "${license_opts}" -- ${cur}) )
            return 0
            ;;
        --format)
            local format_opts="html pdf json csv text"
            COMPREPLY=( $(compgen -W "${format_opts}" -- ${cur}) )
            return 0
            ;;
        --output)
            # Suggest directories
            COMPREPLY=( $(compgen -d -- ${cur}) )
            return 0
            ;;
        --rules)
            # Suggest yaml files
            COMPREPLY=( $(compgen -f -X '!*.y?(a)ml' -- ${cur}) )
            return 0
            ;;
    esac

    # Handle partial command completion
    if [[ ${cur} == -* ]]; then
        COMPREPLY=( $(compgen -W "${global_flags}" -- ${cur}) )
    else
        COMPREPLY=( $(compgen -W "${commands}" -- ${cur}) )
    fi

    return 0
}

complete -F _nessi_completion nessi

# Installation instructions:
# 1. Save this file to a location on your system (e.g., ~/.nessi/completion.sh)
# 2. Add the following line to your ~/.bashrc or ~/.bash_profile:
#    source ~/.nessi/completion.sh
#
# For ZSH users:
# 1. Make sure you have bash completion compatibility enabled:
#    autoload -U +X compinit && compinit
#    autoload -U +X bashcompinit && bashcompinit
# 2. Then source this file:
#    source ~/.nessi/completion.sh
