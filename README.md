# Synced PZ Multiplayer Saves 🎮

Program that syncs multiplayer saves between your friends, allowing every one of you to be a host at any time

## Table of Contents

-   [Description](#description) 📄
-   [Dependencies](#dependencies) 📦
-   [Instructions and usage](#instructions-and-usage) 🛠️
-   [Compiling](#compiling) 🖥️
-   [Contributing](#contributing) 🤝
-   [Known Issues](#known-issues) 🐞
-   [Help](#help) ❓
-   [Authors](#authors) 👥
-   [Version History](#version-history) 📜
-   [License](#license) 📄
-   [Acknowledgments](#acknowledgments) 🙏

## Description

This program is a simple way to sync multiplayer saves between your friends. It works by syncing the save files between the players, so that every player can be a host at any time.
This is useful for games like Project Zomboid, where the host has to be online for the other players to join the game.
So, by using this program, you can play with your friends without having to wait for the host to be online. Anyone with a synced save can host the game.

Only works on WINDOWS and has only been tested on servers that use Steam invites (Invite friends via Steam).

WORKS on Cracked versions of Project Zomboid that uses
[Online-fix](https://online-fix.me/)

Warning: This program is still in development and may have bugs. Use at your own risk, always create backups. But it should be safe to use, since the most recent copy of the save is always kept in the cloud (Git repository).

## Dependencies

-   Git
-   Go 1.23.5 (Just if you want to compile the program)

## Instructions and Usage

See [USAGE.md](markdown/USAGE.md)

## Compiling

```bash
go build .
./syncedpz.exe
```

or just

```bash
go run .
```

## Contributing

Always feel free to contribute to this project. Any help is welcome.

1.  Fork it (<https://github.com/JotaEspig/synced-pz-multiplayer-saves/fork>)
2.  Create your feature branch (`git checkout -b feature/fooBar`)
3.  Commit your changes (`git commit -am 'Add some fooBar'`)
4.  Push to the branch (`git push origin feature/fooBar`)
5.  Create a new Pull Request

## Known Issues

The time it takes to commit and push a totally new save is too long. And, for now, there is no something
like a progress bar to show the user that the program is working.

## Help

Contact the authors

## Authors

-   João Vitor Espig ([JotaEspig](https://github.com/JotaEspig))

## Version History

See [CHANGELOG.md](CHANGELOG.md)

## License

This project is licensed under the Apache License 2.0 - see the LICENSE file for details

## Acknowledgments

-   [README-Template.md](https://gist.github.com/DomPizzie/7a5ff55ffa9081f2de27c315f5018afc)
