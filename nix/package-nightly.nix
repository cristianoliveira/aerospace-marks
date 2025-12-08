{ pkgs, ... }:
  pkgs.buildGoModule rec {
    # name of our derivation
    name = "aerospace-marks";
    version = "nightly"; # Branch 

    go = pkgs.go_1_24;

    # sources that will be used for our derivation.
    src = pkgs.fetchFromGitHub {
      owner = "cristianoliveira";
      repo = "aerospace-marks";
      rev = version;
      sha256 = "sha256-XlV8d37n8rdm+u2pV99HY0/wJ0GpYzETK1syBSMA4yY=";
    };

    vendorHash = "sha256-MKIcIfx0ScHZiW6qtoNy5l0kAFqWK9SXkVVEUL3P2tg=";

    ldflags = [
      "-s" "-w"
      "-X github.com/cristianoliveira/aerospace-marks/cmd.VERSION=${version}"
    ];

    meta = with pkgs.lib; {
      description = "aerospace-marks: I3wm like marks feature";
      homepage = "https://github.com/cristianoliveira/aerospace-marks";
      license = licenses.mit;
      maintainers = with maintainers; [ cristianoliveira ];
    };
  }
