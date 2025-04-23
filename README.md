<p align="center">
  <img src="https://raw.githubusercontent.com/basedandco/sudoku/main/art/logo.png" width="220"><br>
  <b>SudoKu</b><br>
  <i>root privileges? solve the puzzle first</i>
</p>

[![Build](https://img.shields.io/github/actions/workflow/status/basedandco/sudoku/ci.yml?label=build)](…)
[![Made with Spite](https://img.shields.io/badge/made%20with-spite-red)](#)
[![License: WTFPL](https://img.shields.io/badge/license-WTFPL-blue)](LICENSE)

---

## 🥜 TL;DR

`sudo` → **sudoku** → _root_  
A PAM module written in Go. When you try a privileged command it:

1. Spawns a fresh Sudoku (4×4, 6×6, 9×9, or 🔥 “Diabolical” 16×16).
2. Opens a curses UI (or HTML5 fallback over SSH-kitty).
3. Gives you **90 seconds** to solve it.
4. Logs success or failure to `/var/log/auth.log` with 🌶️ random insults.

No more fat-fingered `rm -rf /`. If you can’t Sudoku, you can’t sudo.

---

## 📸 Demo

```console
$ sudo reboot
┌──────────────────────────────┐
│  Welcome to SudoKu!          │    ⏱ 01:30
│  Fill the 9×9 to earn root.  │
└──────────────────────────────┘
```

---

## 🚀 Installation

> ⚠️ Tested on Arch 🦜 (of course) and Ubuntu. Your distro mileage may vary.

```bash
git clone https://github.com/basedandco/sudoku.git
cd sudoku
make && sudo make install        # drops pam_sudoku.so into /lib/security
sudo make enable                 # appends "auth required pam_sudoku.so" into /etc/pam.d/sudo
```

Disable with:

```bash
sudo make disable
```

---

## 🔧 Configuration

Edit `/etc/sudoku.conf`:

| Key          | Default | Note                           |
| ------------ | ------- | ------------------------------ |
| `GRID_SIZE`  | 9       | 4, 6, 9, 16 (only if you dare) |
| `TIME_LIMIT` | 90      | seconds                        |
| `INSULTS`    | 1       | 0 = polite mode 🤢             |
| `LOCKOUT`    | 3       | fails before _n_-minute ban    |
| `SHOW_TIMER` | 1       | display countdown in UI        |

---

## ❓ FAQ

**Q: I got locked out of prod at 3 AM. Now what?**
A: Maybe finish the puzzle next time, champ.

<!-- **Q: Can I pipe the puzzle into `rofi` or a web UI?**
A: Yep—toggle `UI=rofi` or `UI=http` in the config. -->

**Q: Security risk?**
A: The puzzle seed is cryptographically random; the only risk is to your ego.

---

## 🤝 Contributing

PRs welcome—especially:

- New puzzle generators (Killer, Samurai, Nonograms?)
- Better curses UX (mouse support, Vim bindings)
- Translations of insults (need Klingon)

---

## 🪪 License

WTFPL — \*“Do What The F*\*\* You Want Public License.”* See LICENSE for the exciting two-line legal doc.

---

## ☠️ Disclaimer

SudoKu is comedy software. If you brick your CI/CD pipeline because you couldn’t find a hidden 7, that’s on you. Use in prod only if your sense of humor has root.

---

<p align="center"><sub>© 2025 Based & Co. Ltd. · Go solve some puzzles.</sub></p>
