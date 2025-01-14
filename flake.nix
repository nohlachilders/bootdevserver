{
    inputs = {
        nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
        devenv.url = "github:cachix/devenv";
        systems.url = "github:nix-systems/default";
        flake-utils = {
            url = "github:numtide/flake-utils";
            inputs.systems.follows = "systems";
        };
    };

    outputs = { self, nixpkgs, devenv, flake-utils, ... } @ inputs:
        flake-utils.lib.eachDefaultSystem (system: let
            pkgs = nixpkgs.legacyPackages.${system};
        in {
            packages = {
                devenv-up = self.devShells.${system}.default.config.procfileScript;
                devenv-test = self.devShells.${system}.default.config.test;
            };

            devShells.default = devenv.lib.mkShell {
                inherit inputs pkgs;
                modules = [
                    ({pkgs, config, ... }: {
                        # stuff goes here
                        languages.go = {
                            enable = true;
                            enableHardeningWorkaround = true;
                        };
                        services.postgres = {
                            listen_addresses = "127.0.0.1";
                            enable = true;
                            createDatabase = false;
                        };
                        packages = with pkgs; [
                            sqlc
                            gopls
                            delve
                            goose
                            (nixpkgs.legacyPackages.${system}.buildGoModule rec {
                                name = "bootdotdev";
                                src = nixpkgs.legacyPackages.${system}.fetchFromGitHub {
                                    owner = "bootdotdev";
                                    repo = "bootdev";
                                    rev = "b283943";
                                    sha256 = "sha256-ofXMlH1cvhfCFmgjZVMqt/kF8F9ZlD2CPH55d7dkMN8=";
                                };
                                vendorHash = "sha256-jhRoPXgfntDauInD+F7koCaJlX4XDj+jQSe/uEEYIMM=";

                            })
                        ];
                    })
                ];
            };
        }
        );
}
