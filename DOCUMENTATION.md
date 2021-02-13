# Gotokey Documentation

## Keylogger

Gotokey uses several keyloggers with varying levels of complexity to provide the user with the best coverage against text loss.

Gotokey provides a protection system to avoid recording password. The principle is to set a sequence of keys as indicating the beginning of a password. The keylogger will then skip the recording keys until either an end-of-password input is registered, or a given number of keyboard inputs have been sent.

### Raw Keylogger

This keylogger is conceptual. It corresponds to the exact list of key inputs received by gotokey. These inputs are not actually logged.

### Base Keylogger

The base keylogger logs every input as an hexadecimal number between 0x00 and 0xFF. It stores all inputs sent during the same minute on the same timestamped line.

It uses the following symbols:

- `[A-Z]` / `[a-z]` a letter key press / release
- `[.':][0-9]` a digit key press / release / immediate
- `[.':][_NTBLRUDCSAW]` a single-letter-named-key press / release / immediate
- `[-^=]([CSAW]_|_[CSAW]|F[1-9A-O])` a two-letter-named-key press / release / immediate
- `(--|^^|==)[A-Z_]+/` any long-named-key press / release / immediate
- `&` repeat the previous event

Immediate means a press followed by a release, without any other event in between. This doesn't necessarily mean that the two happened in a short period of time. Gotokey doesn't use immediates when the press and the release happen in different minutes.

- Single letter names are:
  - `_` - Spacebar
  - `N` - Newline / enter / return
  - `T` - Tabulation
  - `B` - Backspace
  - `L`, `R`, `U`, `D` - Arrows Left, Right, Up, Down
  - (`C`, `S`, `A`, `W` - Control, Shift, Alt, Win)
- Two letter names are:
  - `C_` - Left Control / `_C` - Right Control
  - `S_` - Left Shift / `_S` - Right Shift
  - `A_` - Left Alt / `_A` - Right Alt
  - `W_` - Left Win / `_W` - Right Win (aka. Super)
  - `F1`-`F9`, `FA`, `FB`, `FC` - F1-F9, F10, F11, F12
  - `FD`-`FO` - F13-F24
- Some common long names:
  - `DELETE`
  - `ESCAPE`
  - `PGUP`, `PGDOWN` - page up, page down
  - `BEGIN`, `END` - Goto beginning of line, end of line
  - `F10`-`F24`
  - `BACKTICK`, `MINUS`, `EQUAL`
  - `BRACKETOPEN`, `BRACKETCLOSE`, `BACKSLASH`
  - `SEMICOLON`, `QUOTE`
  - `COMMA`, `PERIOD`, `SLASH`
  - `ANGLEBRACKET`
  - `PRINT` - Print Screen
  - `MENU` - Menu key (sometimes missing, aka. Apps key)
  - `INSERT` - Insert key (sometimes overloaded on Delete)
  - `CAPSLOCK` - Caps lock

If there is any digit between two inputs, this is the rounded elapsed time between the two inputs or since the beginning of the timestamped minute.

Example:

`i love flowers`

```
2021-02-13T08:32 55Ii1:_1LlOoVve1:_FfLl
2021-02-13T08:33 OowEeRrs
```

The Base Keylogger filters out repeated keypresses **only for modifiers** (CSAW).

### Serial Keylogger

This is the simplest key logger. It records all effectful keys with their modifiers. It stores all inputs sent during the same minute on the same timestamped line.

By default, it assumes a QWERTY keyboard. If you use any other keyboard, you'll
have to provide the Serial Keylogger with a description of your layout such as
the following:

```txt
~!@#$%^&*()_+
`1234567890-=
QWERTYUIOP{}|
qwertyuiop[]\
ASDFGHJKL:
asdfghjkl;
ZXCVBNM<>?
zxcvbnm,./
```

Characters from the provided keyboard description will be logged. The high/log version of a row is selected depending on whether the Shift modifier is active.

If the Control, Shift, Alt or Win modifier is active, a named hotkey will be recorded in a verbose form between a dot `.` and a slash `/`.

Key combination examples:

- Ctrl+X, Ctrl+C, Ctrl+V: `.C.X/` `.C.C/` `.C.V/`
- Alt+F: `.A.F/`
- Ctrl+Shift+F: `.CS.F/`
- Enter: `.ENTER/`
- Backspace: `.BACKSPACE/`
- Escape: `.ESCAPE/`
- Ctrl+Win+Shift+Alt+A: `.CWSA.A/`
- Alt+Shift+Win+Ctrl+A: `.CWSA.A/`
- Left, Right, Up, Down: `.LEFT/`, `.RIGHT/`, `.UP/`, `.DOWN/`
- Home, End, PgUp, PgDown: `.HOME/`, `.END/`, `.PGUP/`, `.PGDOWN/`.
- Ctrl+Up, Ctrl+Home, Ctrl+End: `.C.UP/`, `.C.BEGIN/`, `.C.END/`
- Ctrl+PgUp, Ctrl+PgDown `.C.PGUP/`, `.C.PGDOWN/`
- The dot will be inserted as follows: `.DOT/`

The Serial Keylogger filters out:

- empty modifier keypresses (CSA)
- automatic repeated keypresses for modifiers (CWSA).

### Other keylogger...

### Password protection example

Say you use the following four passwords, and want to protect them:

```
uru@#iuiugcglhiaykgn
caokyzstsnyu04092
cathorserainarabia
yjlguduhjydtscqwpsma
```

You'll set the following password protection rules:

- `uru`
- `caok`
- `cathor`
- `yjlg`

The password protection check is done only for keys right after a click, a tabulation press (including shift+tab), or an enter press is registered. The password protection ends after either a click, a tabulation press or an enter press is registered, or after 80 keyboard keys have been received.
