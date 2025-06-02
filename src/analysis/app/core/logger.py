import logging
import logging.config
import sys
import colorama

from app.core.settings import settings


class ColoredFormatter(logging.Formatter):
    """
    Цветные логи :)

    Не работают в docker logs :(
    """

    COLORS = {
        logging.DEBUG: colorama.Fore.WHITE,
        logging.INFO: colorama.Fore.CYAN,
        logging.WARNING: colorama.Fore.YELLOW,
        logging.ERROR: colorama.Fore.RED,
        logging.CRITICAL: colorama.Fore.MAGENTA,
    }

    def format(self, record: logging.LogRecord) -> str:
        color = self.COLORS.get(record.levelno, colorama.Fore.RESET)
        record.color = color
        message = super().format(record)
        return f"{color}{message}{colorama.Style.RESET_ALL}"


def configure_logger():
    """
    Конфигуратор логгера со всеми его настройками и форматированием.
    """
    colorama.init()

    format_str = (
        "(%(asctime)s [%(levelname)s] %(filename)s:%(lineno)d - "
        "%(message)s)"
    )
    datefmt = "%Y-%m-%dT%H:%M:%S"

    formatters = {
        "console": {
            "()": ColoredFormatter,
            "format": format_str,
            "datefmt": datefmt,
        },
        "file": {"format": format_str, "datefmt": datefmt},
    }

    handlers = {
        "console": {
            "class": "logging.StreamHandler",
            "formatter": "console",
            "stream": sys.stdout,
            "level": settings.log.level,
        }
    }

    if settings.log.save_to_file:
        handlers["file"] = {
            "class": "logging.FileHandler",
            "formatter": "file",
            "filename": "logfile.log",
            "level": settings.log.level,
        }

    logger_config = {
        "version": 1,
        "disable_existing_loggers": False,
        "formatters": formatters,
        "handlers": handlers,
        "loggers": {
            "": {
                "handlers": list(handlers.keys()),
                "level": settings.log.level,
                "propagate": False,
            }
        },
    }

    logging.config.dictConfig(logger_config)


configure_logger()


def get_logger(name: str) -> logging.Logger:
    """
    Возвращает логгер для добавление его в код.
    """
    return logging.getLogger(name)
