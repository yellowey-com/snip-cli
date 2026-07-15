set -l commands add edit find remove run reload

complete -c snip -f
complete -c snip -n "not __fish_seen_subcommand_from $commands" -a add -d "Add new snippet"
complete -c snip -n "not __fish_seen_subcommand_from $commands" -a edit -d "Edit snippet command"
complete -c snip -n "not __fish_seen_subcommand_from $commands" -a find -d "Search snippets by query"
complete -c snip -n "not __fish_seen_subcommand_from $commands" -a remove -d "Remove snippet"
complete -c snip -n "not __fish_seen_subcommand_from $commands" -a run -d "Execute snippet by description"
complete -c snip -n "not __fish_seen_subcommand_from $commands" -a reload -d "Rebuild and restart"

complete -c snip -n "__fish_seen_subcommand_from run edit remove rm" -f -a "(snip _list-descriptions)"
