{ pkgs, ... }:
  pkgs.buildGoModule rec {
    # name of our derivation
    name = "aerospace-marks";
    # FIXME: once we have the first release, we can use the version
    # version = "v0.0.2";
    version = "ee2d3df";

    # sources that will be used for our derivation.
    src = pkgs.fetchFromGitHub {
      owner = "cristianoliveira";
      repo = "aerospace-marks";
      rev = version;
      sha256 = "sha256-WgAH+ILK7IDrYiBBZDwSEfSTebyjF65WlfCWNsqnV8M=";
    };

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
