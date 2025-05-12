import logging
import logging.config
import sys
import colorama  # type: ignore

from configs.config import LOG_LEVEL, SAVE_LOG_TO_FILE


class ColoredFormatter(logging.Formatter):
    """Цветные логи :)"""
    COLORS = {
        logging.DEBUG: colorama.Fore.WHITE,
        logging.INFO: colorama.Fore.CYAN,
        logging.WARNING: colorama.Fore.YELLOW,
        logging.ERROR: colorama.Fore.RED,
        logging.CRITICAL: colorama.Fore.MAGENTA,
    }

    def format(self, record):
        color = self.COLORS.get(record.levelno, colorama.Fore.RESET)
        record.color = color
        message = super().format(record)
        return f"{color}{message}{colorama.Style.RESET_ALL}"


def configure():
    """
    Конфигуратор логирования.
    Включает цветные логи в консоль и возможность сохранения в файл.
    """
    colorama.init()

    formatters = {
        "fileFormatter": {
            "format": "(%(asctime)s [%(levelname)s] %(filename)s:%(lineno)d - "
                      "%(message)s)",
            "datefmt": "%Y-%m-%dT%H:%M:%S%Z"
        },
        "consoleFormatter": {
            "class": "utils.logger_config"
                     ".ColoredFormatter",
            "format": "(%(asctime)s [%(levelname)s] %(filename)s:%(lineno)d - "
                      "%(message)s)",
            "datefmt": "%Y-%m-%dT%H:%M:%S%Z"
        }
    }

    console_handler = {
        "class": "logging.StreamHandler",
        "level": "DEBUG",
        "formatter": "consoleFormatter",
        "stream": sys.stdout
    }

    file_handler = {
        "class": "logging.FileHandler",
        "level": "DEBUG",
        "formatter": "fileFormatter",
        "filename": "logfile.log"
    }

    handlers = {"consoleHandler": console_handler}
    if SAVE_LOG_TO_FILE:
        handlers["fileHandler"] = file_handler

    loggers = {
        "root": {
            "level": LOG_LEVEL,
            "handlers": ["consoleHandler"]
        },
        "analysisLogger": {
            "level": LOG_LEVEL,
            "handlers": ["consoleHandler"],
            "propagate": False
        }
    }

    if SAVE_LOG_TO_FILE:
        loggers["analysisLogger"]["handlers"].append("fileHandler")

    dict_config = {
        "version": 1,
        "disable_existing_loggers": False,
        "formatters": formatters,
        "handlers": handlers,
        "loggers": loggers
    }

    logging.config.dictConfig(dict_config)
