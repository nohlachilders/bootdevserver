{
  description = "A basic flake with a shell";
  inputs.nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
  inputs.systems.url = "github:nix-systems/default";
  inputs.flake-utils = {
    url = "github:numtide/flake-utils";
    inputs.systems.follows = "systems";
  };

  outputs =
    { nixpkgs, flake-utils, ... }:
    flake-utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
        bootdotdev = nixpkgs.legacyPackages.${system}.buildGoModule rec {
            name = "bootdotdev";
            src = nixpkgs.legacyPackages.${system}.fetchFromGitHub {
                owner = "bootdotdev";
                repo = "bootdev";
                rev = "20d9060";
                sha256 = "sha256-XSdTkr4UwZuC35pG3Rgc2o4IEpoVfY93tfHE36H8GDc=";
            };
            vendorHash = "sha256-jhRoPXgfntDauInD+F7koCaJlX4XDj+jQSe/uEEYIMM=";

        };


      in
      {
        devShells.default = pkgs.mkShell { packages = [ bootdotdev pkgs.go pkgs.jq ]; };
      }
    );
}
