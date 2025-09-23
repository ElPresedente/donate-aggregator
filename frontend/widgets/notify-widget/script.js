const socketAddr = ""

//FIELDS

let debugEnabled = false
let testEnabled = false
let testTextType = "Some text"
let titleText = ""
let titleSize = 20
let titleColor = "rgb(255, 255, 255)"
let titleFontName = "Roboto"
let titleMarginLeft = 20
let textSize = 20
let textColor = "rgb(255, 255, 255)"
let textFontName = "Roboto"
let widgetAppearanceTime = 3
let widgetDisappearanceTime = 3
let labelsScrollTime = 10
let backgroundColor = "rgb(255, 255, 255)"
let borderWidth = 1
let borderColor = "rgb(255, 255, 255)"
let borderRadius = 20

//VARIABLES

let textArray = []
let playAnim = false
let render = false

window.addEventListener('onWidgetLoad', function (obj) {
  initWidget(obj)
  testCheck()
  connectWebSocket()
});

function handleEvent(event){
  log(event)
  switch( event.request ){
    case "set-text": return setText( event.text )
    case "reset": return reset()
    case "delete": return delText(text)
  }
}

function renderWidget(){
  if(render){
    document.getElementById("widget").style.display = "block"
  } else {
    document.getElementById("widget").style.display = "none"
  }

  if(playAnim){
    renderPlayAnim()
  } else
  {
    
  }
}

function setText(text){
  textArray.push(text)
  addTextElem(text)
  render = true
  renderWidget();
}

function addTextElem(text){
  const el = document.createElement("p")
  el.className = "text-container"
  //el.style.width = `${sectorWidth}px`;
  el.style.height = "auto"
  el.innerText = text

  if(textArray.length>1)
  {
    el.style.visibility = "hidden"
    playAnim = true
  }
  
  document.getElementById("main-container").appendChild(el);
}

function delText(text){
  textArray = textArray.filter(el => el != text)
  delTextElem(text)
  if(textArray.length==0){
    render = false
  }
  renderWidget()
}

function delTextElem(text){
  const elems = document.getElementsByClassName("text-container");
  for(const el of elems)
  {
    if(el.textContent == text)
    {
      el.remove();
    }
  }
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

  debugEnabled            = fieldData.debugEnabled
  testEnabled             = fieldData.testEnabled
  testTextType            = fieldData.testTextType
  titleText               = fieldData.titleText
  titleSize               = fieldData.titleSize
  titleColor              = fieldData.titleColor
  titleFontName           = fieldData.titleFontName
  titleMarginLeft         = fieldData.titleMarginLeft
  textSize                = fieldData.textSize
  textColor               = fieldData.textColor
  textFontName            = fieldData.textFontName
  widgetAppearanceTime    = fieldData.widgetAppearanceTime
  widgetDisappearanceTime = fieldData.widgetDisappearanceTime
  labelsScrollTime        = fieldData.labelsScrollTime
  backgroundColor         = fieldData.backgroundColor
  borderWidth             = fieldData.borderWidth
  borderColor             = fieldData.borderColor
  borderRadius            = fieldData.borderRadius
}

function testCheck(){
  if(testEnabled){
    testStart()
  }
}

function testStart(){
  setText(testTextType)
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