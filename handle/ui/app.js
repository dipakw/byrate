class SpeedTest {
    // Track downloaded bytes
    downloaded = 0;
    downloadStartedAt = 0;
    downloadStarted = false;

    // Callbacks to notify about changes
    callbacks = [];

    // Reader and controller
    reader = null;
    controller = null;

    constructor({ download }) {
        this.downloadUrl = download;
    }

    async download() {
        this.downloaded = 0;
        this.downloadStartedAt = Date.now();

        await this.downloadReal((done, e) => {
            let status = done ? 'completed' : 'progress';

            if (e) {
                status = 'errored';
            }

            this.notify({
                status,
                type: 'download',
                error: e,
                bytes: this.downloaded,
            });
        });
    }

    async downloadReal(cb) {
        try {
            const response = await fetch(this.downloadUrl);

            if (!response.body) {
                throw new Error("ReadableStream not supported.");
            }

            this.reader = response.body.getReader();

            const readChunk = async () => {
                const { done, value } = await this.reader.read();

                // Count bytes
                this.downloaded += value.byteLength;

                if (!this.downloadStarted && String.fromCharCode.apply(null, value.slice(0, 5)) === "start") {
                    this.downloadStartedAt = Date.now();
                    this.downloadStarted = true;
                }

                cb(done, null);

                await readChunk();
            };

            await readChunk();
        } catch (e) {
            cb(false, e);
        }
    }

    stop() {
        if (this.reader) {
            this.reader.cancel().catch(() => { });
            this.reader = null;
        }

        this.notify({
            type: 'download',
            status: 'stopped',
            bytes: this.downloaded,
        });
    }

    notify(obj) {
        const end = Date.now();
        const duration = (end - this.downloadStartedAt) / 1000;

        this.callbacks.forEach(cb => cb({
            ...obj,
            duration,
            startedAt: this.downloadStartedAt,
            endedAt: end,
            speed: this.format(this.downloaded, this.downloadStartedAt, end),
        }));
    }

    callback(cb) {
        this.callbacks.push(cb);
    }

    format(bytes, start, end) {
        const duration = (end - start) / 1000;

        if (duration <= 0) {
            return "0 bps";
        }

        const bytesPerSecond = bytes / duration;
        const bitsPerSecond = bytesPerSecond * 8;

        let unit = "bps";
        let value = bitsPerSecond;

        const sizes = [
            ["Kbps", 1e3],
            ["Mbps", 1e6],
            ["Gbps", 1e9],
            ["Tbps", 1e12],
            ["Pbps", 1e15],
            ["Ebps", 1e18],
        ];

        for (const [label, size] of sizes) {
            if (bitsPerSecond >= size) {
                unit = label;
                value = bitsPerSecond / size;
            }
        }

        return `${value.toFixed(2)} ${unit}`;
    }
}

class Timer {
    mmEl = null;
    ssEl = null;

    startedAt = 0;
    interval = null;

    constructor({ mmEl, ssEl }) {
        this.mmEl = mmEl;
        this.ssEl = ssEl;
    }

    start() {
        if (this.startedAt === 0) {
            this.startedAt = Date.now();
        }

        this.interval = setInterval(() => {
            let min = 0;
            let sec = Math.floor((Date.now() - this.startedAt) / 1000);

            if (sec >= 60) {
                min = Math.floor(sec / 60);
                sec = sec % 60;
            }

            this.mmEl.textContent = min.toString().padStart(2, '0');
            this.ssEl.textContent = sec.toString().padStart(2, '0');
        }, 50);

        return this;
    }

    pause() {
        clearInterval(this.interval);
        this.interval = null;

        return this;
    }

    reset() {
        clearInterval(this.interval);
        this.startedAt = 0;
        this.interval = null;
        this.mmEl.textContent = '00';
        this.ssEl.textContent = '00';

        return this;
    }
}

(() => {
    // Get necessary elements.
    const hand = document.querySelector('.hand');
    const speedOut = document.querySelector('.speed');
    const button = document.querySelector('#button');
    const meterTop = document.querySelector('.top');

    // Calculate radius and center coordinates for circles.
    const radius = (meterTop.offsetWidth - 56) / 2; // assuming width = diameter
    const centerX = radius + 23;
    const centerY = radius + 22; // because it's a semicircle pointing down

    const totalCircles = 9;
    const angleStart = Math.PI; // 180 degrees
    const angleEnd = 0;         // 0 degrees

    // Add circles
    for (let i = 0; i < totalCircles; i++) {
        const circle = document.createElement('span');
        circle.classList.add('circle');

        const angle = angleStart - (i / (totalCircles - 1)) * (angleStart - angleEnd);

        const x = centerX + radius * Math.cos(angle);
        const y = centerY - radius * Math.sin(angle);

        circle.style.position = 'absolute';
        circle.style.left = `${x}px`;
        circle.style.top = `${y}px`;

        meterTop.appendChild(circle);
    }

    let timeout = null;

    const timer = new Timer({
        mmEl: document.querySelector('.min'),
        ssEl: document.querySelector('.sec'),
    });

    const speedTest = new SpeedTest({
        download: '/download',
    });

    const startTest = () => {
        hand.classList.add('active');
        button.classList.add('active');
        button.textContent = 'Stop';
        timer.reset().start();
        speedTest.download();

        timeout = setTimeout(stopTest, 15000);
    }

    const stopTest = () => {
        clearTimeout(timeout);
        hand.classList.remove('active');
        button.classList.remove('active');
        button.textContent = 'Start';
        timer.pause();
        speedTest.stop();
    }

    const onButtonClick = () => {
        const isActive = hand.classList.contains('active');

        // Start.
        if (!isActive) {
            startTest();
        }

        // Stop.
        if (isActive) {
            stopTest();
        }
    }

    speedTest.callback((data) => {
        if (data.status === 'completed' || data.status === 'errored') {
            stopTest();
        }

        speedOut.textContent = data.speed;
    });

    button.addEventListener('click', onButtonClick);
})();