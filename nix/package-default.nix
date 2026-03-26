{ pkgs, ... }:
  pkgs.buildGoModule rec {
    name = "aerospace-marks";
    version = "v1.0.1";

    go = pkgs.go_1_26;

    # sources that will be used for our derivation.
    src = pkgs.fetchFromGitHub {
      owner = "cristianoliveira";
      repo = "aerospace-marks";
      rev = version;
      sha256 = "sha256-Z1Me3kClHs1BExB324Ko3I0qiFCc8cTGPlaDu9huuNE=";
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
