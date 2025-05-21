<template>
  <el-container class="terminal-container">
    <el-header class="header">Remote Serial Port Server Terminal</el-header>
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
      <el-main class="terminal-main">
        <div id="terminal" class="terminal"></div>
      </el-main>
    </el-container>
  </el-container>
</template>

<script>
import { ref, onMounted, onUnmounted, watch } from 'vue';
import axios from 'axios';
import { Terminal } from '@xterm/xterm';
import '@xterm/xterm/css/xterm.css';
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
    let term;
    let socket;
    const isConnected = ref(false);
    const fontFamily = ref('JetBrains Mono');
    const fontSize = ref(12);

    const fetchApiKey = async () => {
      const response = await axios.get('/api/system/key');
      apiKey.value = response.data.data.key;
      console.log('apiKey:', apiKey.value);
    };

    const fetchPorts = async () => {
      if (!apiKey.value) {
        console.error('API Key is not available. Please ensure fetchApiKey is called first.');
        return;
      }
      const response = await axios.get(`/api/port/?key=${apiKey.value}`);
      ports.value = response.data.data.ports;
      console.log('ports:', ports.value);
    };

    const connectWebSocket = () => {
      if (!form.value.port) {
        ElMessage('请选择端口');
        return;
      }

      const wsUrl = `/api/port/open?key=${apiKey.value}&port=${form.value.port}&baudrate=${form.value.baudRate}&&stopbits=${form.value.stopBits}&parity=${form.value.parity}`;
      socket = new WebSocket(wsUrl);

      socket.onopen = () => {
        console.log('WebSocket connected');
        isConnected.value = true;

        if (!term._keyListenerAttached) {
          term.onKey((key) => {
            const char = key.key;
            if (char === '\x7F') {
              socket.send('\x08');
            } else {
              socket.send(char);
            }
          });
          term._keyListenerAttached = true;
        }
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
      if (term) {
        term.options.fontFamily = fontFamily.value;
        term.options.fontSize = fontSize.value;
      }
    };

    const fitTerminalSize = () => {
      if (term) {
        const terminalContainer = document.getElementById('terminal');
        const cols = Math.floor(terminalContainer.offsetWidth / term.options.fontSize);
        const rows = Math.floor(terminalContainer.offsetHeight / term.options.fontSize);
        term.resize(cols, rows);
      }
    };

    onMounted(async () => {
      await fetchApiKey();
      await fetchPorts();

      term = new Terminal({
        fontFamily: fontFamily.value,
        fontSize: fontSize.value,
      });
      term.open(document.getElementById('terminal'));

      fitTerminalSize(); // 初始化时调整终端大小
    });

    onUnmounted(() => {
      disconnectWebSocket();
      term.offKey();
    });

    // 监听窗口大小变化
    watch(
      () => window.innerWidth,
      () => {
        fitTerminalSize();
      }
    );

    watch(
      () => window.innerHeight,
      () => {
        fitTerminalSize();
      }
    );

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
    };
  },
};
</script>

<style>
@import './assets/font/font.css';

html, body {
  margin: 0; /* 确保没有默认的外边距 */
  display: flex; /* 使用 Flex 布局 */
  flex-direction: column;
}

.terminal-container {
  display: flex;
  flex-direction: column;
  flex: 1; /* 允许终端容器根据父级容器扩展 */
  width: 100vw;
  height: 100vh; /* 确保容器占据整个视口高度 */
  min-width: 1024px; /* 最小宽度 */
  min-height: 768px; /* 最小高度 */
  overflow: hidden;
}

.config-aside {
  background-color: #f5f7fa;
  padding: 20px;
  box-sizing: border-box;
  flex-shrink: 0; /* 防止 config-aside 被压缩 */
  width: 300px; /* 固定宽度 */
}

.terminal-main {
  display: flex;
  flex-direction: column;
  flex: 1; /* 终端主区域占据剩余空间 */
  padding: 0;
  margin: 0;
  width: calc(100% - 300px); /* 动态计算宽度，确保与 aside 对齐 */
}

.header {
  background-color: #409EFF;
  color: white;
  text-align: left;
  line-height: 60px;
  flex-shrink: 0; /* 防止 header 被压缩 */
}

.terminal {
  flex: 1;
  display: flex; /* 确保终端内容也使用 Flex 布局 */
}

/* 添加媒体查询以适应不同屏幕尺寸 */
@media (max-width: 768px) {
  .config-aside {
    width: 100%; /* 在小屏幕上，配置面板占满整个宽度 */
  }

  .terminal-main {
    flex-direction: column; /* 在小屏幕上，终端主区域变为垂直布局 */
  }
}
</style>
