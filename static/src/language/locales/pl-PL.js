export default {
  menu: {
    serial: {
      port: "Port",
      baudRate: "Prędkość (baud)",
      dataBits: "Bity danych",
      stopBits: "Bity stopu",
      parity: "Parzystość",
      flowControl: "Kontrola przepływu",
      placeholders: {
        port: "Wybierz port",
        baudRate: "Wprowadź prędkość",
        dataBits: "Wybierz bity danych",
        stopBits: "Wybierz bity stopu",
        parity: "Wybierz parzystość",
        flowControl: "Wybierz kontrolę przepływu",
      },
      options: {
        parity: {
          none: "Brak",
          odd: "Nieparzysta",
          even: "Parzysta",
          mark: "Znacznik",
          space: "Spacja",
        }
      }
    },
    theme: {
      font: "Czcionka",
      fontSize: "Rozmiar czcionki",
      placeholders: {
        font: "Wybierz czcionkę",
        fontSize: "Wprowadź rozmiar czcionki",
      },
      buttons: {
        connect: "Połącz",
        disconnect: "Rozłącz",
      }
    }
  },
  message: {
    no_port: "Wybierz port",
    websocket_disconnected: "Połączenie WebSocket zostało przerwane. Sprawdź sieć lub połącz ponownie.",
    websocket_error: "Wystąpił błąd WebSocket. Sprawdź sieć lub połącz ponownie.",
    websocket_params_changed: "Parametry formularza uległy zmianie, połączenie WebSocket zostało przerwane."
  },
  log: {
    api_key_not_available: "Klucz API jest niedostępny. Najpierw wprowadź klucz API.",
    websocket_connected: "WebSocket jest połączony.",
    websocket_disconnected: "WebSocket jest rozłączony.",
    websocket_error: "Wystąpił błąd WebSocket:",
    websocket_manually_disconnect: "WebSocket został ręcznie rozłączony.",
  }
}
