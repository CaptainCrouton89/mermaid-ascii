"""Poetry build hook — compile the Go binary before packaging."""

import platform
import subprocess
from pathlib import Path


def build(setup_kwargs=None):
    """Called by Poetry during build."""
    scripts_dir = Path(__file__).parent
    if platform.system() == "Windows":
        script = scripts_dir / "binary.ps1"
        subprocess.run(
            ["powershell", "-ExecutionPolicy", "Bypass", "-File", str(script)],
            check=True,
        )
    else:
        script = scripts_dir / "binary.sh"
        subprocess.run(["bash", str(script)], check=True)


if __name__ == "__main__":
    build()
