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

let pinsArray = []
let render = false

let textInterval = null;
let currentActiveIndex = 0;
let maxHeight = 0;

window.addEventListener('onWidgetLoad', function (obj) {
  initWidget(obj)
  testCheck()
  connectWebSocket()
});

function initTextSizes() {
  const containers = document.querySelectorAll('.text-container');
  const wrapper = document.querySelector('#main-container');

  maxHeight = Math.max(...Array.from(containers).map(el => el.offsetHeight));
  wrapper.style.height = `${maxHeight}px`;
}

function error( err ){
  if( debugEnabled ){
    console.error(err)
  }
}

function handleEvent(event){
  log(event)
  switch( event.request ){
    case "set-text": return setText( event.pin )
    case "reset": return reset()
    case "reset-text": return resetText( event.pin )
  }
}

function startTextRotation() {
    if (textInterval) return;
    
    textInterval = setInterval(() => {
        showText(currentActiveIndex);
        currentActiveIndex = (currentActiveIndex + 1) % pinsArray.length;
    }, labelsScrollTime * 1000);
}

function stopTextRotation() {
    clearInterval(textInterval);
    textInterval = null;
}

function showText(index) {
  const els = document.querySelectorAll('.text-container');

  els.forEach((el, i) => {
    if (i === index) {
      el.classList.add('active')
    } else {
      el.classList.remove('active');
    }
  });
}

function animateAppear(element) {
  element.classList.add('active');
}

function renderWidget() {
  const widget = document.querySelector('#widget');

  if (render) {
    widget.classList.remove("hidden");
    void widget.offsetWidth
    widget.classList.add("visible");
    //widget.classList.add('visible');
    if (pinsArray.length === 1) {
      currentActiveIndex = 0;
      initTextSizes();
      showText(currentActiveIndex);
    } else {
      initTextSizes();
      startTextRotation();
    }
  } else {
    widget.classList.remove('visible');
    setTimeout(() => {
      widget.classList.add("hidden");
      stopTextRotation();
    }, 1000);
  }
}

function setText(obj){
  pinsArray.push(obj)
  addTextElem(obj.value)
  render = true
  renderWidget();
}

function addTextElem(text){
  const el = document.createElement("p")
  el.className = "text-container"
  //el.style.width = `${sectorWidth}px`;
  el.style.height = "auto"
  el.innerText = text
  
  document.getElementById("main-container").appendChild(el);
}

function resetText(obj){
  pinsArray = pinsArray.filter(el => el.value !== obj.value)
  currentActiveIndex = 0 //костыль, мб нормально потом делать

  switch (pinsArray.length)
  {
    case 0:
      render = false;
      renderWidget();
      setTimeout(() => {
        resetTextElem(obj.value)
      }, 1000);
      break;
    case 1:
      resetTextElem(obj.value)
      stopTextRotation();
      showText(0);
      break;
    default:
      resetTextElem(obj.value)
      break;
  }
}

function resetTextElem(text){
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
    } catch (err) {
      error('❌ Ошибка парсинга:', err);
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
  if(testEnabled === "true"){
    testStart()
  }
}

function testStart(){
  const testObj = {'value': testTextType}
  setText(testObj)
}

function warn( w ){
  if( debugEnabled === "true" ){
    console.warn(w)
  }
}

function log( logStr ){
  if( debugEnabled ){
      console.log( logStr )
  }
}