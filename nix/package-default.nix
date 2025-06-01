{ pkgs, ... }:
  pkgs.buildGoModule rec {
    name = "aerospace-marks";
    version = "v0.2.0";

    go = pkgs.go_1_24;

    # sources that will be used for our derivation.
    src = pkgs.fetchFromGitHub {
      owner = "cristianoliveira";
      repo = "aerospace-marks";
      rev = version;
      sha256 = "sha256-B3yumy2cJX2CENXNkF1Q3CzP8SD3uzo0eJ6Om4TOdro=";
    };

    vendorHash = "sha256-jBGebNPvSxjoru+CnqpgT3X3hgH8bTa55AhreJ0bqik=";

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
