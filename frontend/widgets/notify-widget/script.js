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
let render = false

let textInterval = null;
let currentTextIndex = 0;

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
    case "remove-text": return removeText(text)
  }
}

function startTextRotation() {
    if (textInterval) return;
    
    textInterval = setInterval(() => {
        showText(currentTextIndex);
        currentTextIndex = (currentTextIndex + 1) % textArray.length;
    }, labelsScrollTime * 1000);
}

function stopTextRotation() {
    clearInterval(textInterval);
    textInterval = null;
}

function showText(index) {
    const els = document.getElementsByClassName('text-container');
    
    els.forEach((el, i) => {
        if (i === index) {
            el.style.opacity = 1;
            el.style.visibility = 'visible';
            animateAppear(el);
        } else {
            el.style.opacity = 0;
            el.style.visibility = 'hidden';
        }
    });
}

function animateAppear(element) {
    element.style.opacity = 0;
    element.style.visibility = 'visible';
    
    setTimeout(() => {
        element.style.opacity = 1;
    }, 100);
}

function renderWidget() {
  const widget = document.getElementById('widget');

  if (render) {
    widget.classList.add('visible');
    if (textArray.length === 1) {
      showText(0);
    }
  } else {
    widget.classList.remove('visible');
    stopTextRotation();
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

  if (textArray.length > 1) {
    el.style.opacity = 0;
  }
  
  document.getElementById("main-container").appendChild(el);

  if (textArray.length === 2) {
    startTextRotation();
  }
}

function removeText(text){
  textArray = textArray.filter(el => el != text)
  removeTextElem(text)
  currentTextIndex = 0 //костыль, мб нормально потом делать

  if (textArray.length === 1) {
    stopTextRotation();
    showText(0);
  }
  
  if (textArray.length === 0) {
    render = false;
    renderWidget();
  }
}

function romoveTextElem(text){
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