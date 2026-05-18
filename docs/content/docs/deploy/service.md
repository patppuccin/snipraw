---
title: System Service
---

# Running snipraw as a system service

Running snipraw as a system service ensures it starts on boot and restarts on failure. There is no built-in service management, so use your platform's native tooling.

## Linux (systemd)

Create `/etc/systemd/system/snipraw.service`:

```ini
[Unit]
Description=Snipraw snippet server
After=network.target

[Service]
Type=simple
User=your-user
ExecStart=/usr/local/bin/snipraw --dir /home/your-user/snippets
Restart=on-failure
RestartSec=5

[Install]
WantedBy=multi-user.target
```

Enable and start:

```bash
sudo systemctl daemon-reload
sudo systemctl enable --now snipraw
```

Check logs:

```bash
journalctl -u snipraw -f
```

## macOS (launchd)

Create `~/Library/LaunchAgents/com.snipraw.plist`:

```xml
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
  <key>Label</key>
  <string>com.snipraw</string>
  <key>ProgramArguments</key>
  <array>
    <string>/usr/local/bin/snipraw</string>
    <string>--dir</string>
    <string>/Users/you/snippets</string>
  </array>
  <key>RunAtLoad</key>
  <true/>
  <key>KeepAlive</key>
  <true/>
</dict>
</plist>
```

Load it:

```bash
launchctl load ~/Library/LaunchAgents/com.snipraw.plist
```

## Windows (NSSM)

[NSSM](https://github.com/dkxce/NSSM/releases/latest) (Non-Sucking Service Manager) wraps any executable as a Windows service. Install it via [Scoop](https://scoop.sh/#/apps?q=nssm) or [Chocolatey](https://community.chocolatey.org/packages?q=nssm) or [WinGet](https://winget.run/pkg/NSSM/NSSM) or download the zip from the GitHub releases page.

Register snipraw as a service (requires admin shell):

```powershell
nssm install snipraw "C:\path\to\snipraw.exe"
nssm set snipraw AppParameters "--dir C:\Users\you\snippets"
nssm start snipraw
```

To remove:

```powershell
nssm stop snipraw
nssm remove snipraw confirm
```

::: tip
NSSM can write stdout and stderr to a log file. Use `nssm set snipraw AppStdout` and `AppStderr` to configure log paths and keep tabs on what snipraw is doing.
:::
