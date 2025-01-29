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
                        env.GOOSE_DBSTRING = "postgres://nohlachilders@localhost:5432/chirpy";
                        env.GOOSE_DRIVER = "postgres";

                        languages.go = {
                            enable = true;
                            enableHardeningWorkaround = true;
                        };

                        processes = {
                            go-server.exec = "air";
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
                            air
                            openssl
                            (nixpkgs.legacyPackages.${system}.buildGoModule rec {
                                name = "bootdotdev";
                                src = nixpkgs.legacyPackages.${system}.fetchFromGitHub {
                                    owner = "bootdotdev";
                                    repo = "bootdev";
                                    rev = "950affe77fee1cd8b15bf02dcf9d1bcfc08ebde7";
                                    sha256 = "sha256-e6LUAcG0tCTfRWGkJ85jIfjcr4/1fIP61rPPUTDrkjg=";
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
