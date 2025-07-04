{ pkgs, ... }:
  pkgs.buildGo124Module rec {
    # name of our derivation
    name = "aerospace-marks";
    version = "source";

    go = pkgs.go_1_24;

    # sources that will be used for our derivation.
    src = ../.;

    vendorHash = "sha256-9FNgpx4PtA9VCbeMdo+BQaED0MWGNmX8E0YjFNsSy04=";

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
