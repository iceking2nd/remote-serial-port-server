<template>
  <el-container style="min-width: 1024px; min-height: 768px; width: 100vw; height: 100vh">
    <el-header class="header" height="60px">
      <div class="header-title">
      Remote Serial Port Server Terminal
      </div>
      <el-select v-model="locale" class="language-select">
        <el-option key="en-us" label="English" value="en-us"></el-option>
        <el-option key="zh-cn" label="简体中文" value="zh-cn"></el-option>
      </el-select>
    </el-header>
    <el-container>
      <el-aside width="300px" class="config-aside">
        <el-form :model="form" label-width="100px" class="config-form">
          <el-form-item :label="$t('menu.serial.port')">
            <el-select v-model="form.port" :placeholder="$t('menu.serial.placeholders.port')" @change="handleFormChange">
              <el-option v-for="port in ports" :key="port" :label="port" :value="port"></el-option>
            </el-select>
          </el-form-item>
          <el-form-item :label="$t('menu.serial.baudRate')">
            <el-input v-model.number="form.baudRate" :placeholder="$t('menu.serial.placeholders.baudRate')" @input="handleFormChange"></el-input>
          </el-form-item>
          <el-form-item :label="$t('menu.serial.dataBits')">
            <el-select v-model="form.dataBits" :placeholder="$t('menu.serial.placeholders.dataBits')" @change="handleFormChange">
              <el-option label="5" value="5"></el-option>
              <el-option label="6" value="6"></el-option>
              <el-option label="7" value="7"></el-option>
              <el-option label="8" value="8"></el-option>
            </el-select>
          </el-form-item>
          <el-form-item :label="$t('menu.serial.stopBits')">
            <el-select v-model="form.stopBits" :placeholder="$t('menu.serial.placeholders.stopBits')" @change="handleFormChange">
              <el-option label="1" value="1"></el-option>
              <el-option label="1.5" value="1.5"></el-option>
              <el-option label="2" value="2"></el-option>
            </el-select>
          </el-form-item>
          <el-form-item :label="$t('menu.serial.parity')">
            <el-select v-model="form.parity" :placeholder="$t('menu.serial.placeholders.parity')" @change="handleFormChange">
              <el-option :label="$t('menu.serial.options.parity.none')" value="none"></el-option>
              <el-option :label="$t('menu.serial.options.parity.odd')" value="odd"></el-option>
              <el-option :label="$t('menu.serial.options.parity.even')" value="even"></el-option>
              <el-option :label="$t('menu.serial.options.parity.mark')" value="mark"></el-option>
              <el-option :label="$t('menu.serial.options.parity.space')" value="space"></el-option>
            </el-select>
          </el-form-item>
          <el-form-item class="button-group" label-width="0">
            <el-button
              type="primary"
              @click="connectWebSocket"
              :disabled="isConnected"
            >
              {{ $t('menu.theme.buttons.connect') }}
            </el-button>
            <el-button
              type="danger"
              @click="disconnectWebSocket"
              :disabled="!isConnected"
            >
              {{ $t('menu.theme.buttons.disconnect') }}
            </el-button>
          </el-form-item>
        </el-form>
        <!-- 新增表单：字体和文字大小调整 -->
        <el-form label-width="100px" class="config-form">
          <el-form-item :label="$t('menu.theme.font')">
            <el-select v-model="fontFamily" :placeholder="$t('menu.theme.placeholders.font')" @change="updateTerminalStyle">
              <el-option label="JetBrains Mono" value="JetBrains Mono"></el-option>
              <el-option label="Monaco" value="Monaco"></el-option>
              <el-option label="Courier New" value="Courier New"></el-option>
              <el-option label="Consolas" value="Consolas"></el-option>
            </el-select>
          </el-form-item>
          <el-form-item :label="$t('menu.theme.fontSize')">
            <el-input-number :placeholder="$t('menu.theme.placeholders.fontSize')" v-model="fontSize" :min="10" :max="50" @change="updateTerminalStyle"></el-input-number>
          </el-form-item>
        </el-form>
      </el-aside>
      <el-main class="terminal-main" id="terminal-main">
        <div id="terminal" class="terminal" ref="terminal"></div>
      </el-main>
    </el-container>
  </el-container>
</template>

<script>
import { ref, onMounted, onUnmounted, onBeforeUnmount } from 'vue'
import axios from 'axios';
import { Terminal } from '@xterm/xterm';
import { FitAddon } from '@xterm/addon-fit';
import { ClipboardAddon } from '@xterm/addon-clipboard';
import { ElMessage } from 'element-plus';
import { useI18n } from 'vue-i18n'

export default {
  setup() {
    const { t, locale } = useI18n()

    const form = ref({
      port: '',
      baudRate: 9600,
      dataBits: '8',
      stopBits: '1',
      parity: 'none',
    });
    const ports = ref([]);
    const apiKey = ref('');
    let socket;
    let term;
    const fitAddon = new FitAddon();
    const clipboardAddon = new ClipboardAddon();
    const isConnected = ref(false);
    const fontFamily = ref('JetBrains Mono');
    const fontSize = ref(12);

    const CONTROL_CHAR_MAP = {
      '\x7F': '\x08',  // DEL -> BS
    };

    const fetchApiKey = async () => {
      const response = await axios.get(`${import.meta.env.VITE_API_BASE_URL}/api/system/key`);
      apiKey.value = response.data.data.key;
      console.debug('apiKey:', apiKey.value);
    };

    const fetchPorts = async () => {
      if (!apiKey.value) {
        console.error(t('log.api_key_not_available'));
        return;
      }
      const response = await axios.get(`${import.meta.env.VITE_API_BASE_URL}/api/port/?key=${apiKey.value}`);
      ports.value = response.data.data.ports;
      console.debug('ports:', ports.value);
    };

    const connectWebSocket = () => {
      if (!form.value.port) {
        ElMessage(t('message.no_port'));
        return;
      }

      const wsUrl = `${import.meta.env.VITE_API_BASE_URL}/api/port/open?key=${apiKey.value}&port=${form.value.port}&baudrate=${form.value.baudRate}&&stopbits=${form.value.stopBits}&parity=${form.value.parity}`;
      socket = new WebSocket(wsUrl);

      socket.onopen = () => {
        console.log(t('log.websocket_connected'));
        isConnected.value = true;
        term.onData((data) => {
          const processedData = data.replace(/[\x00-\x1F\x7F-\x9F]/g, m => CONTROL_CHAR_MAP[m] || m)
          socket.send(processedData);
        });
      };

      socket.onmessage = (event) => {
        term.write(event.data);
      };

      socket.onclose = () => {
        console.log(t('log.websocket_disconnected'));
        isConnected.value = false;
        ElMessage(t('message.websocket_disconnected'));
      };

      socket.onerror = (error) => {
        console.error(t('log.websocket_error'), error);
        isConnected.value = false;
        ElMessage(t('message.websocket_error'), { type: 'error' });
      };
    };

    const disconnectWebSocket = () => {
      if (socket && socket.readyState === WebSocket.OPEN) {
        socket.close();
        console.log(t('log.websocket_manually_disconnect'));
        isConnected.value = false;
      }
    };

    const handleFormChange = () => {
      if (socket && socket.readyState === WebSocket.OPEN) {
        disconnectWebSocket();
        ElMessage(t('message.websocket_params_changed'));
      }
    };

    const updateTerminalStyle = () => {
      term.options.fontFamily  = fontFamily.value;
      term.options.fontSize = fontSize.value;
      term.refresh(0, term.rows - 1);
      fitAddon.fit();
    };

    onMounted(async () => {
      await fetchApiKey();
      await fetchPorts();

      term = new Terminal({
        fontFamily: fontFamily.value,
        fontSize: fontSize.value,
        scrollback: 10000
      });

      term.loadAddon(fitAddon)
      term.loadAddon(clipboardAddon)
      term.open(document.getElementById('terminal'));
      fitAddon.fit()

      window.addEventListener('resize',fitAddon.fit);
    });

    onBeforeUnmount(() => {
      window.removeEventListener('resize',fitAddon.fit);
      term.dispose();
    })

    onUnmounted(() => {
      disconnectWebSocket();
    });

    return {
      form,
      ports,
      connectWebSocket,
      disconnectWebSocket,
      handleFormChange,
      isConnected,
      fontFamily,
      fontSize,
      updateTerminalStyle,
      t,
      locale
    };
  },
};
</script>

<style>
@import './assets/font/font.css';
@import '@xterm/xterm/css/xterm.css';

.config-aside {
  background-color: #f5f7fa;
  padding: 20px;
  box-sizing: border-box;
  flex-shrink: 0; /* 防止 config-aside 被压缩 */
  width: 300px; /* 固定宽度 */
}

.terminal-main {
  padding: 0;
}

.header {
  background-color: #409EFF;
  color: white;
  display: flex; /* 启用Flex布局 */
  justify-content: space-between; /* 子元素水平分布 */
  align-items: center; /* 子元素垂直居中 */
  flex-shrink: 0; /* 防止 header 被压缩 */
  padding: 0 20px; /* 左右内边距 */
}

.header-title {
  text-align: left;
}

.language-select {
  width: 120px; /* 设置 el-select 的宽度 */
}

.button-group .el-form-item__content {
  display: flex;
  justify-content: center;
  flex-wrap: nowrap; /* 禁止换行 */
  gap: 5px; /* 按钮之间的间距 */
}

.terminal {
  width: 100%;
  height: 100%;
}

.terminal ::-webkit-scrollbar {
  width: 12px;
}

.terminal ::-webkit-scrollbar-thumb {
  background-color: darkgrey;
  border-radius: 6px;
}
</style>
