export default {
  menu: {
    serial: {
      port: "端口",
      baudRate: "波特率",
      dataBits: "数据位",
      stopBits: "停止位",
      parity: "奇偶校验",
      flowControl: "流控制",
      placeholders: {
        port: "请选择端口",
        baudRate: "请填写波特率",
        dataBits: "请选择数据位",
        stopBits: "请选择停止位",
        parity: "请选择奇偶校验",
        flowControl: "请选择流控制",
      },
      options: {
        parity: {
          none: "无",
          odd: "奇校验",
          even: "偶校验",
          mark: "标记",
          space: "空",
        }
      }
    },
    theme: {
      font: "字体",
      fontSize: "字体大小",
      placeholders: {
        font: "请选择字体",
        fontSize: "请填写字体大小",
      },
      buttons: {
        connect: "连接",
        disconnect: "断开",
      }
    }
  },
  message: {
    no_port: "请选择端口",
    websocket_disconnected: "WebSocket连接已断开，请检查网络或重新连接。",
    websocket_error: "WebSocket发生错误，请检查网络或重新连接。",
    websocket_params_changed: "表单参数已变更，WebSocket连接已断开。"
  },
  log: {
    api_key_not_available: "API密钥不存在，请先填写API密钥。", //API Key is not available. Please ensure fetchApiKey is called first.
    websocket_connected: "WebSocket已连接。", //WebSocket connected
    websocket_disconnected: "WebSocket已断开。", //WebSocket disconnected
    websocket_error: "WebSocket发生错误：",
    websocket_manually_disconnect: "WebSocket已手动断开。",
  }
}
