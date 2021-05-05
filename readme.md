# Gerp

Grep's dumb but incredibly fast cousin.

# Installation

Gerp is written in go, it compiles with go 1.16 and above so make sure you have a recent go compiler.

You have two options:

### Installation Via Go Tool

Go 1.16 and above:

`go install github.com/insomnimus/gerp@latest # or use a tag like v0.2.2`

Legacy:

`go get -u -v github.com/insomnimus/gerp`

### Installation Via Git Clone

```
git clone https://github.com/insomnimus/gerp
cd gerp
git checkout main # or use any tag/branch you want
go install
```

# FAQ

### Why should I use it?

-	It's fast. Really fast. Try it if you don't believe me.
-	Especially great on SSD's.
-	Full glob support on windows thanks to mattn's amazing [zglob library](https://github.com/mattn/go-zglob).
-	`--glob` flag on UNIX-like systems for even faster globbing.
-	Utilizes all available CPU cores for the maximum efficiency.
-	Friendly to use, very simple set of flags.

### Why should I not use it?

-	Doesn't have an extensive unicode support like ripgrep. But it makes gerp even faster for that reason.
-	Does not have a rich set of command line options, on purpose.
-	Not as universal as grep, since you can find grep on pretty much any UNIX-based os by default.
-	Does not respect your .gitignore files yet, told you it was dumb.

### Does gerp ship with shell autocompletions?

Yes indeed, run `gerp --help-completions` for the instructions.

### I want to uninstall gerp, how do i do it?

No hard feelings, here's how:

bash:
`rm "$(which gerp)"`

powershell:
`rm (where.exe gerp.exe)`

manual:
-	Locate $GOBIN or $GOPATH/bin.
-	find gerp or gerp.exe.
-	delete it.

# How fast are we talking about?

Here's a small test comparing gerp with ripgrep, by no means this is how you should bench and compare stuff though.
Note that, ripgreps globbing library may be the culprit here, overall i still believe ripgrep is faster.

```
/home> $g=measure-command{gerp --hidden --quiet 'package\s[main]{4}' **/*.go}
/home> $g=measure-command{gerp --hidden --quiet 'package\s[main]{4}' **/*.go}
/home> $g=measure-command{gerp --hidden --quiet 'package\s[main]{4}' **/*.go}
/home> $r=measure-command{rg --hidden 'package\s[main]{4}' --glob='**/*.go' -j=4}
/home> $r=measure-command{rg --hidden 'package\s[main]{4}' --glob='**/*.go' -j=4}
/home> $r=measure-command{rg --hidden 'package\s[main]{4}' --glob='**/*.go' -j=4}
/home> $r

Days              : 0
Hours             : 0
Minutes           : 0
Seconds           : 2
Milliseconds      : 583
Ticks             : 25835777
TotalDays         : 2.99025196759259E-05
TotalHours        : 0.000717660472222222
TotalMinutes      : 0.0430596283333333
TotalSeconds      : 2.5835777
TotalMilliseconds : 2583.5777
/home> $g

Days              : 0
Hours             : 0
Minutes           : 0
Seconds           : 1
Milliseconds      : 877
Ticks             : 18770703
TotalDays         : 2.17253506944444E-05
TotalHours        : 0.000521408416666667
TotalMinutes      : 0.031284505
TotalSeconds      : 1.8770703
TotalMilliseconds : 1877.0703
```
