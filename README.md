Tethealla Server Checker
========================
Connects to Teth servers and notifies you via Telegram bot when one or more have gone offline.

Roadmap
-------
1. ~~Patch server~~ ✔️
2. ~~Login server~~ ✔️
3. **Ship server**

Requirements
------------
  - Working Tethealla-based installation of PSOBB servers
  - golang (if building from source or installing on Linux)

Installation
------------
All installations require making a "server-checker.conf" file.  A sample one is provided below.
Lines that begin with a hash are not read by the program.
This file must be in the same directory as the binary.

```
# Sample configuration file. Substitute your own values for each line.
# Bot key
bot key
# Message ID
message id
# Patch server (whatever its bound to; this may not be localhost!)
127.0.0.1:11000
```

### Windows
Download the executable.

### Linux (or building from source)
Run `go build` in the top-level directory, and then optionally `go install`.
