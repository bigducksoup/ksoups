<script setup>
import 'xterm/css/xterm.css'
import {Terminal} from 'xterm'
import {FitAddon} from 'xterm-addon-fit';
import { CanvasAddon } from 'xterm-addon-canvas';
import {onMounted, ref} from "vue";
import { useMessage } from 'naive-ui';

const message = useMessage()

const props = defineProps(['sshId','width','height'])
const termEle = ref(null)


onMounted(()=>{

  const sid = window.localStorage.getItem('sid');

  if (!sid){
    message.error('请先登录！')
    return
  }

  const term = new Terminal({
    rendererType: "canvas", //渲染类型
    convertEol: true, //启用时，光标将设置为下一行的开头
    scrollback: 300, //终端中的回滚量
    disableStdin: false, //是否应禁用输入
    // cursorStyle: "underline", //光标样式
    cursorBlink: true, //光标闪烁
    theme: {
      foreground: "#C7C7C7", //字体
      background: "#000000", //背景色
      cursor: "help", //设置光标
      lineHeight: 20
    }
  });
  const fitAddon = new FitAddon()
  term.open(termEle.value);
  term.loadAddon(new CanvasAddon());
  term.loadAddon(fitAddon)
  fitAddon.fit()

  window.addEventListener('resize',()=>{
    fitAddon.fit()
  })

  let ws = new WebSocket("ws://localhost:8080/ws/ssh?sshInfoId="+props.sshId + "&sid=" + sid);
  ws.binaryType = "arraybuffer";
  
  ws.onmessage = async (e) => {
    term.write(new Uint8Array(e.data))
  }

  ws.onclose = async (e,c) => {
    message.error('连接已断开,请检查网络或连接参数')
  }

  ws.onerror = (err) => {
    console.log(err)
  }
  
  term.onKey((e) => {
    ws.send(e.key)
  })
})

</script>

<template>

  <div ref="termEle" class="bg-black"
       :style="{
          width: props.width ?? '100%',
          height:props.height ?? '100%'
       }"
  >

  </div>

</template>

<style scoped>

</style>