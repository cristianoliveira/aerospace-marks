{ pkgs, ... }:
  pkgs.buildGoModule rec {
    # name of our derivation
    name = "aerospace-marks";
    version = "nightly"; # Branch 

    # sources that will be used for our derivation.
    src = pkgs.fetchFromGitHub {
      owner = "cristianoliveira";
      repo = "aerospace-marks";
      rev = version;
      sha256 = "sha256-JCsXHLJhjDkip5+FbjQyoCh8BEq3hA4vd9JUl2DLWxY=";
    };

    vendorHash = "sha256-QG9vWwagqOnB3lgUf9Rzj9BDQefaz+UcV4vH1p1lBZ0=";

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
