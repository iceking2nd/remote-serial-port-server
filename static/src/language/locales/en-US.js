export default {
  menu: {
    serial: {
      port: "Port",
      baudRate: "Baud Rate",
      dataBits: "Data Bits",
      stopBits: "Stop Bits",
      parity: "Parity",
      flowControl: "Flow Control",
      placeholders: {
        port: "Please select a port",
        baudRate: "Please enter the baud rate",
        dataBits: "Please select data bits",
        stopBits: "Please select stop bits",
        parity: "Please select parity",
        flowControl: "Please select flow control",
      },
      options: {
        parity: {
          none: "None",
          odd: "Odd",
          even: "Even",
          mark: "Mark",
          space: "Space",
        }
      }
    },
    theme: {
      font: "Font",
      fontSize: "Font Size",
      placeholders: {
        font: "Please select a font",
        fontSize: "Please enter font size",
      },
      buttons: {
        connect: "Connect",
        disconnect: "Disconnect",
      }
    }
  },
  message: {
    no_port: "Please select a port",
    websocket_disconnected: "WebSocket connection is disconnected. Please check the network or reconnect.",
    websocket_error: "WebSocket error occurred. Please check the network or reconnect.",
    websocket_params_changed: "Form parameters have changed, WebSocket connection is disconnected."
  },
  log: {
    api_key_not_available: "API key is not available. Please enter the API key first.",
    websocket_connected: "WebSocket is connected.",
    websocket_disconnected: "WebSocket is disconnected.",
    websocket_error: "WebSocket error occurred:",
    websocket_manually_disconnect: "WebSocket has been manually disconnected.",
  }
}
