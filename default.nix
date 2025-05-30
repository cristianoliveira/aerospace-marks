{ pkgs ? import <nixpkgs> {} }:

{
  default = pkgs.callPackage ./nix/package-default.nix { inherit pkgs; };
  nightly = pkgs.callPackage ./nix/package-nightly.nix { inherit pkgs; };
  source = pkgs.callPackage ./nix/package-source.nix { inherit pkgs; };
}
