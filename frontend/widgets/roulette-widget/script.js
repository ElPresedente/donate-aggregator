const frontImages = [
  "https://images2.imgbox.com/01/c0/Rfkm3Nyn_o.png",
  "https://images2.imgbox.com/4e/a0/HuetaDBI_o.png",
  "https://images2.imgbox.com/28/b5/nJHFsdM9_o.png"
];

const backImages = [
  "https://images2.imgbox.com/83/dd/OHH2giFg_o.png",
  "https://images2.imgbox.com/d2/01/1tl7sMLf_o.png",
  "https://images2.imgbox.com/46/77/GpI2Pn23_o.png"
];

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


let isAnimated = false;

window.addEventListener('load', () => {
  resetTrack();
  connectWebSocket();
});

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
}

function getWeightedRandomIndex() {
  const rand = Math.random() * 100.0;
  if (rand < 69.0) return 0;
  if (rand < 99.0) return 1;
  return 2;
}

function appendToTrack(text, sectorId, categoryKey = null) {
  const track = document.getElementById("track");

  let index = 0;
  let frontImage = frontImages[index];
  let backImage = backImages[index];

  for (let i = 0; i < repeats; i++) {
    const el = document.createElement("div");
    el.className = "sector";
    el.style.width = `${sectorWidth}px`;
    el.style.height = `${sectorHeight}px`;

    const id = `coin-${sectorId}-${i}`;
    const isTarget = i === (targetRepeats);

    if (isTarget)
      index = categoryMapping[categoryKey];
    else
      index = getWeightedRandomIndex();

    frontImage = frontImages[index];
    backImage = backImages[index];

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

