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
      sha256 = "sha256-5wUy6/G9EaGZ1EFLuJGy8v2AfU3PP3ynaxFV3VJh5XU=";
    };

    vendorHash = "sha256-KJdVwjOSDHRJMI7D4jay3MHX/tajSIpVf/o2Zpt4bB0=";

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
