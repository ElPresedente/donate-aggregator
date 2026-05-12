const winterImages = {
  pointer: "https://raw.githubusercontent.com/ElPresedente/donate-aggregator/dev/widget_images/winter_roulette/pin.png",
  wrapper: "https://raw.githubusercontent.com/ElPresedente/donate-aggregator/dev/widget_images/winter_roulette/bg.png",
  frame: "https://raw.githubusercontent.com/ElPresedente/donate-aggregator/dev/widget_images/winter_roulette/decor.png",
  frontImagesCount: 3,
  backImagesCount: 3,
  frontImageV1: "https://raw.githubusercontent.com/ElPresedente/donate-aggregator/dev/widget_images/winter_roulette/com_cl.png",
  frontImageV2: "https://raw.githubusercontent.com/ElPresedente/donate-aggregator/dev/widget_images/winter_roulette/rar_cl.png",
  frontImageV3: "https://raw.githubusercontent.com/ElPresedente/donate-aggregator/dev/widget_images/winter_roulette/leg_cl.png",
  backImageV1: "https://raw.githubusercontent.com/ElPresedente/donate-aggregator/dev/widget_images/winter_roulette/com_op.png",
  backImageV2: "https://raw.githubusercontent.com/ElPresedente/donate-aggregator/dev/widget_images/winter_roulette/rar_op.png",
  backImageV3: "https://raw.githubusercontent.com/ElPresedente/donate-aggregator/dev/widget_images/winter_roulette/leg_op.png"
}

const normalImages = {
  pointer: "https://raw.githubusercontent.com/ElPresedente/donate-aggregator/dev/widget_images/common_roulette/pin.png",
  wrapper: "https://raw.githubusercontent.com/ElPresedente/donate-aggregator/dev/widget_images/common_roulette/bg.png",
  frame: "https://raw.githubusercontent.com/ElPresedente/donate-aggregator/dev/widget_images/common_roulette/dec.png",
  frontImagesCount: 3,
  backImagesCount: 3,
  frontImageV1: "https://raw.githubusercontent.com/ElPresedente/donate-aggregator/dev/widget_images/common_roulette/com_cl.png",
  frontImageV2: "https://raw.githubusercontent.com/ElPresedente/donate-aggregator/dev/widget_images/common_roulette/rar_cl.png",
  frontImageV3: "https://raw.githubusercontent.com/ElPresedente/donate-aggregator/dev/widget_images/common_roulette/leg_cl.png",
  backImageV1: "https://raw.githubusercontent.com/ElPresedente/donate-aggregator/dev/widget_images/common_roulette/com_op.png",
  backImageV2: "https://raw.githubusercontent.com/ElPresedente/donate-aggregator/dev/widget_images/common_roulette/rar_op.png",
  backImageV3: "https://raw.githubusercontent.com/ElPresedente/donate-aggregator/dev/widget_images/common_roulette/leg_op.png"
}

const blueClaireImages = {
  pointer: "https://raw.githubusercontent.com/ElPresedente/donate-aggregator/dev/widget_images/blue_claire_roulette/pin.png",
  pointerSize: "60px",
  pointerTop: "30px",
  wrapper: "https://raw.githubusercontent.com/ElPresedente/donate-aggregator/dev/widget_images/blue_claire_roulette/bg.png",
  frame: "https://raw.githubusercontent.com/ElPresedente/donate-aggregator/dev/widget_images/blue_claire_roulette/dec.png",
  frontImagesCount: 5,
  backImagesCount: 1,
  frontImageV1: "https://raw.githubusercontent.com/ElPresedente/donate-aggregator/dev/widget_images/blue_claire_roulette/com_cl.png",
  frontImageV2: "https://raw.githubusercontent.com/ElPresedente/donate-aggregator/dev/widget_images/blue_claire_roulette/uncom_cl.png",
  frontImageV3: "https://raw.githubusercontent.com/ElPresedente/donate-aggregator/dev/widget_images/blue_claire_roulette/rar_cl.png",
  frontImageV4: "https://raw.githubusercontent.com/ElPresedente/donate-aggregator/dev/widget_images/blue_claire_roulette/epic_cl.png",
  frontImageV5: "https://raw.githubusercontent.com/ElPresedente/donate-aggregator/dev/widget_images/blue_claire_roulette/leg_cl.png",
  backImageV1: "https://raw.githubusercontent.com/ElPresedente/donate-aggregator/dev/widget_images/blue_claire_roulette/bg_op.png",
}

const pinkImages = {
  pointer: "https://raw.githubusercontent.com/ElPresedente/donate-aggregator/dev/widget_images/pink_roulette/pin.png",
  pointerSize: "60px",
  pointerTop: "30px",
  wrapper: "https://raw.githubusercontent.com/ElPresedente/donate-aggregator/dev/widget_images/pink_roulette/bg.png",
  frame: "https://raw.githubusercontent.com/ElPresedente/donate-aggregator/dev/widget_images/pink_roulette/dec.png",
  frontImagesCount: 5,
  backImagesCount: 1,
  frontImageV1: "https://raw.githubusercontent.com/ElPresedente/donate-aggregator/dev/widget_images/pink_roulette/com_cl.png",
  frontImageV2: "https://raw.githubusercontent.com/ElPresedente/donate-aggregator/dev/widget_images/pink_roulette/uncom_cl.png",
  frontImageV3: "https://raw.githubusercontent.com/ElPresedente/donate-aggregator/dev/widget_images/pink_roulette/rar_cl.png",
  frontImageV4: "https://raw.githubusercontent.com/ElPresedente/donate-aggregator/dev/widget_images/pink_roulette/epic_cl.png",
  frontImageV5: "https://raw.githubusercontent.com/ElPresedente/donate-aggregator/dev/widget_images/pink_roulette/leg_cl.png",
  backImageV1: "https://raw.githubusercontent.com/ElPresedente/donate-aggregator/dev/widget_images/pink_roulette/bg_op.png",
}

let images = {}



const categoryMapping = {
  "Обычные": 0,
  "Необычные": 0,
  "Редкие": 1,
  "Эпические": 1,
  "Легендарные": 2,
  "Артефакты": 2
};

let isSpinning = false;
let globalSectorIdCounter = 0;
let ws;

const tickSound = new Audio("https://files.donationalerts.com/system/widgets/roulette/sounds/spin/Fix.mp3"); // можно заменить на локальный
tickSound.volume = 0.3;

const sectorWidth = 220;
const sectorHeight = 150;
const repeats = 50;
const targetOffset = 3;
const targetRepeats = repeats - targetOffset;
const rouletteTimeScroll = 6000;
const rouletteTimeScrollDelay = 500;
const rouletteAftrScrollTimeDelay = 4000;
const showRouletteTimeDelay = 2000;
const container = document.getElementById("roulette-container");
const donationQueue = [];

let previousOffset = 0; //для рассчета звуков при пересечении секторов

let isAnimated = false;

window.addEventListener('load', () => {
  resetTrack();
  connectWebSocket();
});

window.addEventListener('onWidgetLoad', function (obj) {
  initWidget(obj)
  initBGImages()
});

function initWidget(widgetLoadEventObject){
  const fieldData = widgetLoadEventObject.detail.fieldData;
  console.log(fieldData.rouletteType)
  switch(fieldData.rouletteType){
    case "normal":
      images = normalImages
      break
    case "winter":
      images = winterImages
      break
    case "pink":
      images = pinkImages
      break
    case "blueClaire":
      images = blueClaireImages
      break
    default:
      images = normalImages
      break
  }
}

function initBGImages(){
  let elements = document.querySelectorAll('.roulette-frame'); 

  elements.forEach(element => {
    element.style.backgroundImage = `url(${images.frame})`;
  });

  elements = document.querySelectorAll('.roulette-wrapper');

  elements.forEach(element => {
    element.style.backgroundImage = `url(${images.wrapper})`;
  });

  elements = document.querySelectorAll('.pointer');

  elements.forEach(element => {
    element.style.backgroundImage = `url(${images.pointer})`;
    if (images.pointerSize) {
      element.style.width = images.pointerSize;
      element.style.height = images.pointerSize;
    }
    if (images.pointerTop) {
      element.style.top = images.pointerTop;
    }
  });
}

function connectWebSocket() {
  const RETRY_INTERVAL = 5000;
  ws = new WebSocket('ws://localhost:8080/ws?type=roulette');

  ws.onopen = () => {
    console.log('✅ Подключено к серверу WebSocket');
  };

  ws.onmessage = (event) => {
    try {
      eventHandler(JSON.parse(event.data));
    } catch (error) {
      console.error('❌ Ошибка парсинга:', error);
    }
  };

  ws.onclose = () => {
    console.warn('⚠️ Соединение закрыто. Повторная попытка через 5 секунд...');
    setTimeout(connectWebSocket, RETRY_INTERVAL);
  };

  ws.onerror = (err) => {
    console.error('❌ Ошибка WebSocket:', err);
    ws.close(); // Принудительно закрываем, чтобы сработал onclose и началась повторная попытка
  };
}

function eventHandler(event) {
  console.log(event)
  switch (event.request)
  {
    case "enqueue-spins": return enqueueSpinsHandler(event.spins);
    case "reset":         return resetHandler();
  }
}

function enqueueSpinsHandler(spins)
{
  console.log("Цикл объектов");
  for (const item of spins) {
    console.log(item);
    donationQueue.push({text: item.sector, category: item.category}); 
  }
  resetTrack();     // Очищаем старый трек
  setTimeout(() => showRoulette(), showRouletteTimeDelay);// Показываем рулетку
  processQueue();   // Стартуем очередь
}

function resetHandler()
{
  donationQueue.length = 0;
  resetTrack(); 
  hideRoulette();
}

function spinTo(sectorId) {
  console.log("spinTo");
  setTimeout(() => {
    const track = document.getElementById("track");
    const wrapperWidth = document.querySelector('.roulette-inner-wrapper').clientWidth;
    const centerOffset = wrapperWidth / 2 - sectorWidth / 2;

    const totalSectors = track.children.length;
    const targetIndex = totalSectors - targetOffset; //проверить постановку
    const totalOffset = targetIndex * sectorWidth - centerOffset; //возможно не надо будет минусовать 1

    track.style.transition = `transform ${rouletteTimeScroll}ms cubic-bezier(0.25, 0.1, 0.25, 1)`;
    track.style.transform = `translateX(-${totalOffset}px)`;

    const startTime = performance.now();
    const distanceDelta = totalOffset - previousOffset; // ✅ только разница
    let lastTickedSector = Math.floor(previousOffset / sectorWidth); // ✅ начинаем с текущего сектора

    const easing = cubicBezier(0.25, 0.1, 0.25, 1);

    function tickLoop(now) {
      const elapsed = now - startTime;
      const progress = Math.min(elapsed / rouletteTimeScroll, 1);
      const easedProgress = easing(progress);

      const currentOffset = previousOffset + distanceDelta * easedProgress;
      const currentSector = Math.floor(currentOffset / sectorWidth);

      if (currentSector !== lastTickedSector) {
        lastTickedSector = currentSector;
        tickSound.currentTime = 0;
        tickSound.play().catch(() => {});

        const pointer = document.querySelector('.pointer');
        if (pointer) {
          pointer.classList.remove('kick'); // сброс, если не успел закончиться
          void pointer.offsetWidth; // force reflow
          pointer.classList.add('kick');
        }
      }

      if (progress < 1) {
        requestAnimationFrame(tickLoop);
      }
    }

    requestAnimationFrame(tickLoop);
    previousOffset = totalOffset; // ✅ сохраняем для следующего spinTo
    
    setTimeout(() => {
      const coinSpanId = `coin-${sectorId}-${targetRepeats}`;
      const span = document.getElementById(coinSpanId);
      if (!span) {
        console.warn(`Не найден span с id ${coinSpanId}`);
        return;
      }
      const coinInner = span?.closest(".coin-inner");

      if (coinInner) {
        coinInner.classList.add("flipped");
      }
      setTimeout(() => {
        isSpinning = false;
        clearFlips();
        if (donationQueue.length > 0) {
          processQueue();
        } else {
          hideRoulette();
        }
      }, rouletteAftrScrollTimeDelay);
    }, rouletteTimeScroll + rouletteTimeScrollDelay);

  }, 1000 + showRouletteTimeDelay); // после прокрутки
}

function processQueue() {
  
  if (isSpinning || donationQueue.length === 0) return;
  console.log("processQueue");
  console.log(donationQueue);
  isSpinning = true;

  const {text, category} = donationQueue.shift();
  const sectorId = globalSectorIdCounter++;

  appendToTrack(text, sectorId, category);
  spinTo(sectorId);
}

function showRoulette() {
  console.log("showRoulette");
  container.classList.remove("hidden");
  void container.offsetWidth;
  container.classList.add("visible");
}

function hideRoulette() {
  isSpinning = false;
  container.classList.remove("visible");
  setTimeout(() => {
    container.classList.add("hidden");
    const reply = {
        request: "spins-done"
      }
      ws.send( JSON.stringify( reply ))
    resetTrack();
  }, 1000);
}

function clearFlips() {
  const track = document.getElementById("track");
  for (const sectorEl of track.children) {
    const coinInner = sectorEl.querySelector(".coin-inner");
    coinInner.classList.remove("flipped");
  }
}

function resetTrack(){
  console.log("resetTrack");
  isSpinning = false;
  const track = document.getElementById("track");
  track.innerHTML = "";
  track.style.transform = "translateX(0)";
  track.style.transition = "none";
  globalSectorIdCounter = 0;
  previousOffset = 0;
}

//Тут начнутся костыли. Мне лень думать, как обрабатывать дизайн с 3 и 5 картинками
//Поэтому разделим это на несколько функций

//Старая функция рандома для трёх редкостей
function getWeightedRandomIndex() {
  const rand = Math.random() * 100.0;
  if (rand < 69.0) return 0;
  if (rand < 99.0) return 1;
  return 2;
}

function getWeightedRandomIndexForFiveImgs() {
  const rand = Math.random() * 100.0;
  if (rand < 50.0) return 0;
  if (rand < 75.0) return 1;
  if (rand < 91.0) return 2;
  if (rand < 98.0) return 3;
  return 4;
}

function appendToTrack(text, sectorId, categoryKey = null) {
  const track = document.getElementById("track");

  let index = 0;
  let frontImage = images.frontImageV1
  let backImage = images.backImageV1

  for (let i = 0; i < repeats; i++) {
    const el = document.createElement("div");
    el.className = "sector";
    el.style.width = `${sectorWidth}px`;
    el.style.height = `${sectorHeight}px`;

    const id = `coin-${sectorId}-${i}`;
    const isTarget = i === (targetRepeats);

    //Клянусь, мне так похуй на хуйню ниже
    //Не я придуман комбинации сначала из 3 фронт 3 бэк картинки, а потом еще добавил 5 фронт 1 бэк
    //Откуда я знаю, что будет в будущем. Мне похуй

    if (isTarget)
      index = categoryMapping[categoryKey];
    else if(images.frontImagesCount == 3) {
      index = getWeightedRandomIndex();
      frontImage = getFrontImage(index);
    }
    else if(images.frontImagesCount == 5) {
      index = getWeightedRandomIndexForFiveImgs();
      frontImage = getFrontImageForFiveImgs(index);
    }

    if(images.backImage == 3)
      backImage = getBackImage(index);
    else if(images.backImage == 1)
      backImage = getBackImageForOneImgs();
    

    el.innerHTML = `
      <div class="coin" style="width: ${sectorHeight}px; height: ${sectorHeight}px;">
        <div class="coin-inner">
          <div class="coin-front">
            <img src="${backImage}" alt="back" />
            <span class="coin-text" id="${id}">${isTarget ? text : ""}</span>
          </div>
          <div class="coin-back">
            <img src="${frontImage}" alt="front" />
          </div>
        </div>
      </div>
    `;

    track.appendChild(el);
  }

  track.style.width = `${track.children.length * sectorWidth}px`;
}

//Функция для трёх картинок фронтовых
function getFrontImage(index){
  switch(index){
    case 0:
      return images.frontImageV1;
    case 1:
      return images.frontImageV2;
    case 2:
      return images.frontImageV3;
    default:
      return images.frontImageV1
  }
}

//Функция для пяти картинок фронтовых
function getFrontImageForFiveImgs(index){
  switch(index){
    case 0:
      return images.frontImageV1;
    case 1:
      return images.frontImageV2;
    case 2:
      return images.frontImageV3;
    case 3:
      return images.frontImageV4;
    case 4:
      return images.frontImageV5;
    default:
      return images.frontImageV1
  }
}

//Функция для трех картинок бэковых
function getBackImage(index){
  switch(index){
    case 0:
      return images.backImageV1;
    case 1:
      return images.backImageV2;
    case 2:
      return images.backImageV3;
    default:
      return images.backImageV1
  }
}

//Функция для одной картинки бэковой
function getBackImageForOneImgs(){
  return images.backImageV1
}

function cubicBezier(p1x, p1y, p2x, p2y) {
  // Бернулли-полиномы для bezier-кривой
  const cx = 3 * p1x;
  const bx = 3 * (p2x - p1x) - cx;
  const ax = 1 - cx - bx;

  const cy = 3 * p1y;
  const by = 3 * (p2y - p1y) - cy;
  const ay = 1 - cy - by;

  function bezierX(t) {
    return ((ax * t + bx) * t + cx) * t;
  }

  function bezierY(t) {
    return ((ay * t + by) * t + cy) * t;
  }

  // Ньютоно-Рафсон для поиска t по x (так как input – progress)
  function solve(x, epsilon = 1e-5) {
    let t = x;
    for (let i = 0; i < 8; i++) {
      const x2 = bezierX(t) - x;
      if (Math.abs(x2) < epsilon) return bezierY(t);
      const dX = (3 * ax * t + 2 * bx) * t + cx;
      if (Math.abs(dX) < epsilon) break;
      t -= x2 / dX;
    }

    // Фоллбэк бинарным поиском
    let t0 = 0, t1 = 1;
    t = x;
    while (t0 < t1) {
      const x2 = bezierX(t);
      if (Math.abs(x2 - x) < epsilon) return bezierY(t);
      if (x > x2) t0 = t;
      else t1 = t;
      t = (t1 - t0) * 0.5 + t0;
    }

    return bezierY(t);
  }

  return solve;
}

