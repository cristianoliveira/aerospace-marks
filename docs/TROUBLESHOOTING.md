# Trouble Shooting

- PATH may not contain aerospace and aerospace-marks
  - Solution: aad `exec.env-vars.PATH = "\$HOME/golang/bin:/run/current-system/sw/bin:\${PATH}";`
