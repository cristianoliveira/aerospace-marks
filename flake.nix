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
            go_1_24

            golangci-lint

            # To create new subcommands, run:
            # cobra-cli add <subcommand-name>
            cobra-cli

            # To generate the mock for the interfaces, run:
            # mockgen -source=./pkg/cli/cli.go -destination=./pkg/cli/mock/mock_cli.go -package=mock
            mockgen
          ];
        };

        packages = import ./default.nix {
          inherit pkgs;
        };
    });
}
