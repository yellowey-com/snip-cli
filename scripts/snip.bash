_snip_completions() {
    local cur prev commands
    cur="${COMP_WORDS[COMP_CWORD]}"
    prev="${COMP_WORDS[COMP_CWORD - 1]}"
    commands="add edit find remove run reload"

    case "$prev" in
    run | edit | remove | rm)
        local IFS=$'\n'
        COMPREPLY=($(compgen -W "$(snip _list-descriptions 2>/dev/null)" -- "$cur"))
        return
        ;;
    esac

    COMPREPLY=($(compgen -W "$commands" -- "$cur"))
}

complete -F _snip_completions snip
