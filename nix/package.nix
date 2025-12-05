{ pkgs }:

pkgs.buildGoModule {
  pname = "peso";
  version = "0.1.0";
  src = pkgs.lib.cleanSource ../.;

  vendorHash = "sha256-aLRlPX8k5E4GhwB6Jx2t9K3mN7vL6Q0R4S1T5U8V9W0=";

  # Enable CGO for SQLite support
  env.CGO_ENABLED = "1";

  # Build dependencies for SQLite
  buildInputs = with pkgs; [
    sqlite
    pkg-config
  ];

  # Copy migrations and static assets to output
  postInstall = ''
    mkdir -p $out/share/peso
    cp -r ${../.}/migrations $out/share/peso/
    cp -r ${../.}/templates $out/share/peso/ || true
    cp -r ${../.}/static $out/share/peso/ || true
  '';

  meta = with pkgs.lib; {
    description = "A weight tracking application built with Go";
    homepage = "https://github.com/emiliopalmerini/peso";
    license = licenses.mit;
    maintainers = [ ];
  };
}
