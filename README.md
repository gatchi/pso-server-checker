Tethealla Server Checker
========================
Connects to local PSOBB servers and notifies you via a Telegram bot when one or more have gone offline.

Roadmap
-------
1. ~~Patch server~~
2. ~~Login server~~
3. ~~Ship server~~
4. Auto-reconnect

Requirements
------------
  - Working Tethealla-based installation of PSOBB servers
  - golang (if building from source or installing on Linux)

Installation
------------
Pull the repo using `git pull` or `go get github.com/gatchi/server-checker`.
Then run `go build` in the top-level directory, and then optionally `go install`.
Make sure to fill out the config file template and either keep it where it is or
(preferably) move it to your local config folder (/usr/local/etc/psobb-server-checker/server-checker.conf).

If you use Windows (64-bit only) you can also optionally download the pre-compiled binary included in this repo
instead of pulling the whole repo and installing golang.  I will try to keep it as up-to-date as the rest of
the code.
