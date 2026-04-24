export default {
  menu: {
    serial: {
      port: "連接埠",
      baudRate: "鮑率",
      dataBits: "資料位元",
      stopBits: "停止位元",
      parity: "同位檢查",
      flowControl: "流量控制",
      placeholders: {
        port: "請選擇連接埠",
        baudRate: "請輸入鮑率",
        dataBits: "請選擇資料位元",
        stopBits: "請選擇停止位元",
        parity: "請選擇同位檢查",
        flowControl: "請選擇流量控制",
      },
      options: {
        parity: {
          none: "無",
          odd: "奇同位",
          even: "偶同位",
          mark: "標記",
          space: "空白",
        }
      }
    },
    theme: {
      font: "字型",
      fontSize: "字型大小",
      placeholders: {
        font: "請選擇字型",
        fontSize: "請輸入字型大小",
      },
      buttons: {
        connect: "連線",
        disconnect: "斷線",
      }
    }
  },
  message: {
    no_port: "請選擇連接埠",
    websocket_disconnected: "WebSocket 連線已中斷，請檢查網路或重新連線。",
    websocket_error: "WebSocket 發生錯誤，請檢查網路或重新連線。",
    websocket_params_changed: "表單參數已變更，WebSocket 連線已中斷。"
  },
  log: {
    api_key_not_available: "API 金鑰不可用，請先輸入 API 金鑰。",
    websocket_connected: "WebSocket 已連線。",
    websocket_disconnected: "WebSocket 已斷線。",
    websocket_error: "WebSocket 發生錯誤：",
    websocket_manually_disconnect: "WebSocket 已手動斷線。",
  }
}
