import yaml
import os
from pathlib import Path


def get_config_path() -> Path:
    """
    Определяет путь к конфигурационному файлу
    на основе переменной окружения APP_ENV.
    """
    env = os.getenv("APP_ENV", "default")
    filename = "config.test.yml" if env == "test" else "config.yml"
    return Path(__file__).resolve().parent.parent.parent / filename


def load_yaml_config(path: Path = None) -> dict:
    """
    Загружает конфигурацию из YAML файла, если он существует.
    Иначе возвращает пустой словарь.
    """
    path = path or get_config_path()
    if path.exists():
        with path.open("r") as f:
            return yaml.safe_load(f)
    else:
        print(f"[INFO] Config file not found at {path}, using default values.")
        return {}
