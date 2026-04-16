Output docs/demo.gif

Set Shell "zsh"
Set FontSize 22
Set Width 1400
Set Height 800
Set Padding 20
Set BorderRadius 8
Set WindowBar Colorful
Set Theme "Builtin Dark"
Set TypingSpeed 50ms
Set Framerate 30

Hide
Type "cd /Users/sven1103/git/opencode-config-cli/.worktrees/feature/improved-lipgloss-styling"
Enter
Sleep 500ms
Type "clear"
Enter
Show

Type "./occo --version"
Enter
Sleep 1s

Type "./occo --help"
Enter
Sleep 2s

Type "./occo source list"
Enter
Sleep 2s

Type "./occo bundle install test-add --preset default --dry-run"
Enter
Sleep 2s

Type "./occo bundle install --help"
Enter
Sleep 3s

Type "./occo source list --with-presets 2>&1 | head -12"
Enter
Sleep 3s

Type "clear"
Enter
Sleep 500ms

Type "./occo bundle install test-add --preset default --project-root /tmp/demo-test 2>&1"
Enter
Sleep 2s

Type "./occo bundle status --project-root /tmp/demo-test"
Enter
Sleep 3s
