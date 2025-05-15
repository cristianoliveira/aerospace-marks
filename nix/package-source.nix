{ pkgs, ... }:
  pkgs.buildGoModule rec {
    # name of our derivation
    name = "aerospace-marks";
    version = "source";

    # sources that will be used for our derivation.
    src = ../.;

    vendorHash = "sha256-0s4aCxaWRolYdLadouszxnrlT+9x+PpJiOaQ4pUPhAM=";

    ldflags = [
      "-s" "-w"
      "-X main.VERSION=${version}"
    ];

    meta = with pkgs.lib; {
      description = "aerospace-marks: I3wm like marks feature";
      homepage = "https://github.com/cristianoliveira/aerospace-marks";
      license = licenses.mit;
      maintainers = with maintainers; [ cristianoliveira ];
    };
  }
