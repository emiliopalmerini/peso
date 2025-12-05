{ pkgs }:

pkgs.mkShell {
  buildInputs = with pkgs; [
    go
    sqlite
    pkg-config
  ];
}
