{ pkgs, ... }:
  pkgs.buildGo124Module rec {
    name = "aerospace-marks";
    version = "v0.1.0";

    # sources that will be used for our derivation.
    src = pkgs.fetchFromGitHub {
      owner = "cristianoliveira";
      repo = "aerospace-marks";
      rev = version;
      sha256 = "sha256-mXYHzCNZOKSb/tAo+6X5Fpq5uf3/R2MxxluRyNagwMw=";
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
