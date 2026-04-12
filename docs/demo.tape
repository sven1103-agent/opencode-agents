Output docs/demo.gif

Set Shell "zsh"
Set FontSize 28
Set Width 1100
Set Height 680
Set Padding 18
Set BorderRadius 8
Set WindowBar Colorful
Set Theme "Builtin Dark"
Set TypingSpeed 40ms
Set Framerate 30

Hide
Type "cd /tmp/oc-demo-recording/work"
Enter
Sleep 500ms
Show

Type "mkdir -p demo-project"
Enter
Sleep 700ms

Type "oc --version"
Enter
Sleep 1s

Type "oc source add qbicsoftware/opencode-config-bundle --name qbic"
Enter
Sleep 1.2s

Type "oc source list"
Enter
Sleep 1.2s

Type "oc bundle apply qbic --preset mixed --project-root demo-project"
Enter
Sleep 1.5s

Type "oc bundle status --project-root demo-project"
Enter
Sleep 1.2s

Type "ls demo-project"
Enter
Sleep 2s
