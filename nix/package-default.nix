{ pkgs, ... }:
  pkgs.buildGoModule rec {
    name = "aerospace-marks";
    version = "v0.2.1";

    go = pkgs.go_1_24;

    # sources that will be used for our derivation.
    src = pkgs.fetchFromGitHub {
      owner = "cristianoliveira";
      repo = "aerospace-marks";
      rev = version;
      sha256 = "sha256-ki5xMo/HfG5WNx+hIgMO5nbgWgmAPKHd7hS82SMUivU=";
    };

    vendorHash = "sha256-gkPFhSZtkgFGpiZBWO8TA+RqejZKfdrjhpxNeOgUuDo=";

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
