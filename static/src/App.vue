<template>
  <el-container style="min-width: 1024px; min-height: 768px; width: 100vw; height: 100vh">
    <el-header class="header" height="60px">Remote Serial Port Server Terminal</el-header>
    <el-container>
      <el-aside width="300px" class="config-aside">
        <el-form :model="form" label-width="100px" class="config-form">
          <el-form-item label="端口名">
            <el-select v-model="form.port" placeholder="请选择端口" @change="handleFormChange">
              <el-option v-for="port in ports" :key="port" :label="port" :value="port"></el-option>
            </el-select>
          </el-form-item>
          <el-form-item label="波特率">
            <el-input v-model.number="form.baudRate" placeholder="请输入波特率" @input="handleFormChange"></el-input>
          </el-form-item>
          <el-form-item label="数据位">
            <el-select v-model="form.dataBits" placeholder="请选择数据位" @change="handleFormChange">
              <el-option label="5" value="5"></el-option>
              <el-option label="6" value="6"></el-option>
              <el-option label="7" value="7"></el-option>
              <el-option label="8" value="8"></el-option>
            </el-select>
          </el-form-item>
          <el-form-item label="停止位">
            <el-select v-model="form.stopBits" placeholder="请选择停止位" @change="handleFormChange">
              <el-option label="1" value="1"></el-option>
              <el-option label="1.5" value="1.5"></el-option>
              <el-option label="2" value="2"></el-option>
            </el-select>
          </el-form-item>
          <el-form-item label="奇偶性">
            <el-select v-model="form.parity" placeholder="请选择奇偶性" @change="handleFormChange">
              <el-option label="无" value="none"></el-option>
              <el-option label="奇校验" value="odd"></el-option>
              <el-option label="偶校验" value="even"></el-option>
              <el-option label="标记" value="mark"></el-option>
              <el-option label="空" value="space"></el-option>
            </el-select>
          </el-form-item>
          <el-form-item>
            <el-button
              type="primary"
              @click="connectWebSocket"
              :disabled="isConnected"
            >
              连接
            </el-button>
            <el-button
              type="danger"
              @click="disconnectWebSocket"
              :disabled="!isConnected"
            >
              断开连接
            </el-button>
          </el-form-item>
        </el-form>
        <!-- 新增表单：字体和文字大小调整 -->
        <el-form label-width="100px" class="config-form">
          <el-form-item label="字体">
            <el-select v-model="fontFamily" placeholder="请选择字体" @change="updateTerminalStyle">
              <el-option label="JetBrains Mono" value="JetBrains Mono"></el-option>
              <el-option label="Monaco" value="Monaco"></el-option>
              <el-option label="Courier New" value="Courier New"></el-option>
              <el-option label="Consolas" value="Consolas"></el-option>
            </el-select>
          </el-form-item>
          <el-form-item label="文字大小">
            <el-input-number v-model="fontSize" :min="10" :max="50" @change="updateTerminalStyle"></el-input-number>
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

export default {
  setup() {
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
      console.log('apiKey:', apiKey.value);
    };

    const fetchPorts = async () => {
      if (!apiKey.value) {
        console.error('API Key is not available. Please ensure fetchApiKey is called first.');
        return;
      }
      const response = await axios.get(`${import.meta.env.VITE_API_BASE_URL}/api/port/?key=${apiKey.value}`);
      ports.value = response.data.data.ports;
      console.log('ports:', ports.value);
    };

    const connectWebSocket = () => {
      if (!form.value.port) {
        ElMessage('请选择端口');
        return;
      }

      const wsUrl = `${import.meta.env.VITE_API_BASE_URL}/api/port/open?key=${apiKey.value}&port=${form.value.port}&baudrate=${form.value.baudRate}&&stopbits=${form.value.stopBits}&parity=${form.value.parity}`;
      socket = new WebSocket(wsUrl);

      socket.onopen = () => {
        console.log('WebSocket connected');
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
        console.log('WebSocket disconnected');
        isConnected.value = false;
        ElMessage('WebSocket连接已断开，请检查网络或重新连接。');
      };

      socket.onerror = (error) => {
        console.error('WebSocket error:', error);
        isConnected.value = false;
        ElMessage('WebSocket发生错误，请检查网络或重新连接。', { type: 'error' });
      };
    };

    const disconnectWebSocket = () => {
      if (socket && socket.readyState === WebSocket.OPEN) {
        socket.close();
        console.log('WebSocket manually disconnected');
        isConnected.value = false;
      }
    };

    const handleFormChange = () => {
      if (socket && socket.readyState === WebSocket.OPEN) {
        disconnectWebSocket();
        ElMessage('表单参数已变更，WebSocket连接已断开。');
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
      updateTerminalStyle
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
  text-align: left;
  line-height: 60px;
  flex-shrink: 0; /* 防止 header 被压缩 */
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
