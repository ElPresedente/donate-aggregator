const socketAddr = ""


let debugEnabled = false


window.addEventListener('onWidgetLoad', function (obj) {
  initWidget(obj)
  connectWebSocket()
});

function handleEvent(event){
  log(event)
  switch( event.request ){
    case "set-text": return setText( event.text )
    case "reset": return reset()
  }
}

function setText(text){
  const elem = document.getElementById('text-container')
  elem.innerText = text
}

function reset(){
  const elem = document.getElementById('text-container')
  elem.innerText = ''
}

function connectWebSocket() {
  const RETRY_INTERVAL = 5000;
  ws = new WebSocket('ws://localhost:8080/ws?type=reward');

  ws.onopen = () => {
    log('✅ Подключено к серверу WebSocket');
  };

  ws.onmessage = (event) => {
    try {
      log(event)
      handleEvent(JSON.parse(event.data));
    } catch (error) {
      error('❌ Ошибка парсинга:', error);
    }
  };

  ws.onclose = () => {
    warn('⚠️ Соединение закрыто. Повторная попытка через 5 секунд...');
    setTimeout(connectWebSocket, RETRY_INTERVAL);
  };

  ws.onerror = (err) => {
    error('❌ Ошибка WebSocket:', err);
    ws.close(); // Принудительно закрываем, чтобы сработал onclose и началась повторная попытка
  };
}

function initWidget(widgetLoadEventObject){
  const fieldData = widgetLoadEventObject.detail.fieldData;
  debugEnabled = fieldData.debugEnabled
}

function warn( w ){
  if( debugEnabled ){
    console.warn(w)
  }
}

function error( err ){
  if( debugEnabled ){
    console.error(err)
  }
}

function log( logStr ){
  if( debugEnabled ){
      console.log( logStr )
  }
}