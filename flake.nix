{
  description = "aerospace-marks: I3wm like marks feature";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs";
    utils.url = "github:numtide/flake-utils";
  };

  outputs = { nixpkgs, utils, ... }: 
    utils.lib.eachDefaultSystem (system:
      let
        pkgs = import nixpkgs { inherit system; };
      in {
        devShells.default = pkgs.mkShell {
          packages = with pkgs; [
            go

            # To create new subcommands, run:
            # cobra-cli add <subcommand-name>
            cobra-cli
          ];
        };

        packages = {
          default = pkgs.callPackage ./nix/package-default.nix {};
          nightly = pkgs.callPackage ./nix/package-nightly.nix {};
          source = pkgs.callPackage ./nix/package-source.nix {};
        };
    });
}
