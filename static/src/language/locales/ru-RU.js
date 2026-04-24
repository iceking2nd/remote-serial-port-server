export default {
  menu: {
    serial: {
      port: "Порт",
      baudRate: "Скорость (бод)",
      dataBits: "Биты данных",
      stopBits: "Стоп-биты",
      parity: "Чётность",
      flowControl: "Управление потоком",
      placeholders: {
        port: "Выберите порт",
        baudRate: "Введите скорость",
        dataBits: "Выберите биты данных",
        stopBits: "Выберите стоп-биты",
        parity: "Выберите чётность",
        flowControl: "Выберите управление потоком",
      },
      options: {
        parity: {
          none: "Нет",
          odd: "Нечётная",
          even: "Чётная",
          mark: "Маркер",
          space: "Пробел",
        }
      }
    },
    theme: {
      font: "Шрифт",
      fontSize: "Размер шрифта",
      placeholders: {
        font: "Выберите шрифт",
        fontSize: "Введите размер шрифта",
      },
      buttons: {
        connect: "Подключить",
        disconnect: "Отключить",
      }
    }
  },
  message: {
    no_port: "Выберите порт",
    websocket_disconnected: "Соединение WebSocket разорвано. Проверьте сеть или переподключитесь.",
    websocket_error: "Произошла ошибка WebSocket. Проверьте сеть или переподключитесь.",
    websocket_params_changed: "Параметры формы изменены, соединение WebSocket разорвано."
  },
  log: {
    api_key_not_available: "API-ключ недоступен. Пожалуйста, сначала введите API-ключ.",
    websocket_connected: "WebSocket подключён.",
    websocket_disconnected: "WebSocket отключён.",
    websocket_error: "Ошибка WebSocket:",
    websocket_manually_disconnect: "WebSocket отключён вручную.",
  }
}
