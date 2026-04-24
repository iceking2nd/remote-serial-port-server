export default {
  menu: {
    serial: {
      port: "Anschluss",
      baudRate: "Baudrate",
      dataBits: "Datenbits",
      stopBits: "Stoppbits",
      parity: "Parität",
      flowControl: "Flusssteuerung",
      placeholders: {
        port: "Bitte Anschluss auswählen",
        baudRate: "Bitte Baudrate eingeben",
        dataBits: "Bitte Datenbits auswählen",
        stopBits: "Bitte Stoppbits auswählen",
        parity: "Bitte Parität auswählen",
        flowControl: "Bitte Flusssteuerung auswählen",
      },
      options: {
        parity: {
          none: "Keine",
          odd: "Ungerade",
          even: "Gerade",
          mark: "Markierung",
          space: "Leerzeichen",
        }
      }
    },
    theme: {
      font: "Schriftart",
      fontSize: "Schriftgröße",
      placeholders: {
        font: "Bitte Schriftart auswählen",
        fontSize: "Bitte Schriftgröße eingeben",
      },
      buttons: {
        connect: "Verbinden",
        disconnect: "Trennen",
      }
    }
  },
  message: {
    no_port: "Bitte Anschluss auswählen",
    websocket_disconnected: "WebSocket-Verbindung wurde getrennt. Bitte Netzwerk prüfen oder neu verbinden.",
    websocket_error: "WebSocket-Fehler aufgetreten. Bitte Netzwerk prüfen oder neu verbinden.",
    websocket_params_changed: "Formularparameter wurden geändert, WebSocket-Verbindung wurde getrennt."
  },
  log: {
    api_key_not_available: "API-Schlüssel nicht verfügbar. Bitte zuerst den API-Schlüssel eingeben.",
    websocket_connected: "WebSocket ist verbunden.",
    websocket_disconnected: "WebSocket ist getrennt.",
    websocket_error: "WebSocket-Fehler aufgetreten:",
    websocket_manually_disconnect: "WebSocket wurde manuell getrennt.",
  }
}
