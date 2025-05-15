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
      sha256 = "sha256-TQAcKtIid1+zP6dRixcuH6jZZHko1+ujZH2Kya+Snq8=";
    };

    vendorHash = "sha256-1pCoyPjokM+Qm+mbNQLCOmVVqwHDUxJ10RyvLV/becQ=";

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
