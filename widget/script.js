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

let isSpinning = false;
let globalSectorIdCounter = 0;

const sectorWidth = 220;
const sectorHeight = 150;
const repeats = 50;
const rouletteTimeScroll = 6000;
const rouletteTimeDelay = 2000;
const container = document.getElementById("roulette-container");
const donationQueue = [];

let isAnimated = false;

function renderTrack() {
  const track = document.getElementById("track");
  track.innerHTML = "";


  for (let i = 0; i < repeats; i++) {
    const el = document.createElement("div");
    el.className = "sector";
    el.style.width = `${sectorWidth}px`;
    el.style.height = `${sectorHeight}px`;

    const frontImage = frontImages[Math.floor(Math.random() * frontImages.length)];
    const backImage = backImages[Math.floor(Math.random() * backImages.length)];

    el.innerHTML = `
      <div class="coin" style="width: ${sectorHeight}px; height: ${sectorHeight}px;">
        <div class="coin-inner">
          <div class="coin-front">
            <img src="${backImage}" alt="back" />
            <span class="coin-text" id="coin-${i}"></span>
          </div>
          <div class="coin-back">
            <img src="${frontImage}" alt="front" />
          </div>
        </div>
      </div>
    `;

    track.appendChild(el);
  }
  track.style.width = `${repeats * sectorWidth}px`;
}

window.addEventListener('load', () => {
  const ws = new WebSocket('ws://localhost:8080/ws');
  ws.onopen = () => {
    console.log('Подключено к серверу');
    ws.send('Тестовое сообщение');
  };
  ws.onmessage = (event) => {
    try {
      const obj = JSON.parse(event.data);
    } catch (error) {
      console.error('Ошибка парсинга:', error);
    }
    // Добавляем ВСЕ элементы в очередь
    for (const item of obj) {
      donationQueue.push(item.item);  // item.sector можно сохранить, если понадобится
    }
    resetTrack();     // Очищаем старый трек
    showRoulette();   // Показываем рулетку
    processQueue();   // Стартуем очередь
};
  ws.onclose = () => console.log('Соединение закрыто');
  renderTrack();
});

function spinTo(text = "") {
  //ПРОДУМАТЬ АНИМАЦИЮ
  clearFlips();
  showRoulette();
  const coinSpan = document.getElementById("coin-25");
  if (coinSpan) {
    coinSpan.innerText = text;
  } else {
    console.warn("coin-25 не найден!");
  }


  setTimeout(() => {
    const track = document.getElementById("track");
    const wrapperWidth = document.querySelector('.roulette-inner-wrapper').clientWidth;

    const centerOffset = wrapperWidth / 2 - sectorWidth / 2;
    const targetIndex = Math.floor(repeats / 2);
    const totalOffset = targetIndex * sectorWidth - centerOffset;
    track.style.transform = `translateX(-${totalOffset}px)`;



    track.style.transition = "none";
    track.style.transform = "translateX(0)";
    void track.offsetWidth;

    track.style.transition = `transform ${rouletteTimeScroll}ms cubic-bezier(0.25, 0.1, 0.25, 1)`;
    track.style.transform = `translateX(-${totalOffset}px)`;

    setTimeout(() => {
      const sectorEl = track.children[25];
      const coinInner = sectorEl.querySelector('.coin-inner');
      if (coinInner) coinInner.classList.add('flipped');

      setTimeout(() => {
        checkQueue();
      }, rouletteTimeDelay);
    }, rouletteTimeScroll + 100);

  }, 1000);
}

//НОВЫЙ СПИНТУ
function spinTo(sectorId) {
  const track = document.getElementById("track");
  const wrapperWidth = document.querySelector('.roulette-inner-wrapper').clientWidth;
  const centerOffset = wrapperWidth / 2 - sectorWidth / 2;

  const totalSectors = track.children.length;
  const targetIndex = totalSectors - (repeats - 2); // ставим последний нужный элемент в центр
  const totalOffset = (targetIndex - 1) * sectorWidth - centerOffset;

  track.style.transition = "transform 0.8s ease-out";
  track.style.transform = `translateX(-${totalOffset}px)`;

  setTimeout(() => {
    const coinSpanId = `coin-${sectorId}-${repeats - 2}`;
    const span = document.getElementById(coinSpanId);
    const coinInner = span?.closest(".coin-inner");

    if (coinInner) {
      coinInner.classList.add("flipped");
    }

    setTimeout(() => {
      if (donationQueue.length > 0) {
        processQueue();
      } else {
        setTimeout(() => {
          hideRoulette();
          isSpinning = false;
        }, 1000);
      }
    }, rouletteTimeDelay);

  }, 900); // после прокрутки
}


function checkQueue() {
  if(donationQueue.length != 0) {
    setTimeout(() => {
      processQueue();
    }, 1000);
  } else if (isSpinning) {
    hideRoulette();
  }
}

function handleDonation() {
  if (!isSpinning) {
    processQueue();
  }
}

function processQueue() {
  if (isSpinning || donationQueue.length === 0) return;

  isSpinning = true;

  const text = donationQueue.shift();
  const sectorId = globalSectorIdCounter++;

  appendToTrack(text, sectorId);
  spinTo(sectorId);
}


function clearFlips() {
  const track = document.getElementById("track");
  for (const sectorEl of track.children) {
    const coinInner = sectorEl.querySelector(".coin-inner");
    coinInner.classList.remove("flipped");
  }
}

function showRoulette() {
  isSpinning = true;
  container.classList.remove("hidden");
  void container.offsetWidth;
  container.classList.add("visible");
}

function hideRoulette() {
  isSpinning = false;
  container.classList.remove("visible");
  setTimeout(() => {
    container.classList.add("hidden");
    //clearFlips();
    resetTrack();
  }, 1000);
}

function resetTrack(){
  const track = document.getElementById("track");
  track.innerHTML = "";
  track.style.transform = "translateX(0)";
  track.style.transition = "none";
  globalSectorIdCounter = 0;
}

function appendToTrack(text, sectorId) {
  const track = document.getElementById("track");

  for (let i = 0; i < repeats; i++) {
    const el = document.createElement("div");
    el.className = "sector";
    el.style.width = `${sectorWidth}px`;
    el.style.height = `${sectorHeight}px`;

    const frontImage = frontImages[Math.floor(Math.random() * frontImages.length)];
    const backImage = backImages[Math.floor(Math.random() * backImages.length)];

    const id = `coin-${sectorId}-${i}`;
    const isTarget = i === (repeats - 2);

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


document.addEventListener("myCustomEvent", function(event) {
  const payload = event.detail;
  console.log("Custom event received with payload:", payload);
  // Process the payload, e.g., update UI
});

window.addEventListener("onEventReceived", function (obj) {
  if(obj.detail && obj.detail.event && obj.detail.event.isTestEvent) return;

  const { listener, event } = obj.detail;
  console.log(listener);
  if (listener === "tip-latest") {
    donationQueue.push("Текст"); 
    handleDonation();
  }
});

window.addEventListener('load', renderTrack);

