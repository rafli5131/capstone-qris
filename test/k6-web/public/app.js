const scenarioGrid = document.getElementById('scenarioGrid');
const runList = document.getElementById('runList');
const metrics = document.getElementById('metrics');
const logOutput = document.getElementById('logOutput');
const activeStatus = document.getElementById('activeStatus');
const activeRun = document.getElementById('activeRun');
const refreshBtn = document.getElementById('refreshBtn');

let scenarios = [];
let selectedRunId = null;

const formatNumber = (value, digits = 2) => {
  if (value === null || value === undefined || Number.isNaN(Number(value))) return '-';
  return Number(value).toLocaleString('id-ID', { maximumFractionDigits: digits });
};

const formatMs = (value) => value === null || value === undefined ? '-' : `${formatNumber(value, 2)} ms`;
const formatPercent = (value) => value === null || value === undefined ? '-' : `${formatNumber(value * 100, 2)}%`;

async function api(path, options) {
  const response = await fetch(path, {
    headers: { 'Content-Type': 'application/json' },
    ...options
  });
  const data = await response.json();
  if (!response.ok) throw new Error(data.error || 'Request gagal');
  return data;
}

function renderScenarios(active) {
  scenarioGrid.innerHTML = scenarios.map((scenario) => `
    <article class="scenario-card">
      <h3>${scenario.name}</h3>
      <p>${scenario.description}</p>
      <div class="scenario-meta">
        <span>${scenario.file}</span>
        <span>${scenario.users}</span>
        <span>${scenario.duration}</span>
        <span>${scenario.threshold}</span>
      </div>
      <button class="primary" data-scenario="${scenario.id}" ${active ? 'disabled' : ''}>Jalankan</button>
    </article>
  `).join('');

  scenarioGrid.querySelectorAll('button[data-scenario]').forEach((button) => {
    button.addEventListener('click', () => startRun(button.dataset.scenario));
  });
}

function renderRuns(runs) {
  if (!runs.length) {
    runList.innerHTML = '<p class="eyebrow">Belum ada hasil test.</p>';
    return;
  }

  runList.innerHTML = runs.map((run) => `
    <button class="run-item" data-run="${run.id}">
      <span class="run-title">
        <strong>${run.scenarioName}</strong>
        <span class="badge ${run.status}">${run.status}</span>
      </span>
      <span class="eyebrow">${new Date(run.startedAt).toLocaleString('id-ID')}</span>
    </button>
  `).join('');

  runList.querySelectorAll('button[data-run]').forEach((button) => {
    button.addEventListener('click', () => selectRun(button.dataset.run));
  });
}

function renderMetrics(run) {
  const summary = run.summary;
  if (!summary) {
    metrics.innerHTML = '<p class="eyebrow">Summary belum tersedia. Tunggu sampai K6 selesai menulis hasil.</p>';
    return;
  }

  const services = summary.services || {};
  const go = services.go || {};
  const legacy = services.legacy || {};

  metrics.innerHTML = `
    ${renderServiceMetrics('Sistem Baru Go', go)}
    ${renderServiceMetrics('Legacy Java', legacy)}
    ${renderOverallMetrics(summary.overall || summary)}
  `;
}

function renderServiceMetrics(title, service) {
  const emptyNote = service.hasData === false
    ? '<p class="eyebrow">Run lama belum punya metric terpisah. Jalankan ulang scenario untuk mengisi bagian ini.</p>'
    : '';
  const partialNote = service.partialData
    ? '<p class="eyebrow">Run lama: hanya jumlah request yang bisa dipisah. Jalankan ulang scenario untuk latency dan error rate per sistem.</p>'
    : '';

  const items = [
    ['HTTP Requests', formatNumber(service.httpReqs, 0)],
    ['Req / Detik', formatNumber(service.httpReqRate, 2)],
    ['P95 Duration', formatMs(service.durationP95)],
    ['Avg Duration', formatMs(service.durationAvg)],
    ['Max Duration', formatMs(service.durationMax)],
    ['Error Rate', formatPercent(service.failedRate)]
  ];

  return `
    <section class="metric-group">
      <h3>${title}</h3>
      ${emptyNote}
      ${partialNote}
      <div class="metric-grid">
        ${items.map(([label, value]) => metricCard(label, value)).join('')}
      </div>
    </section>
  `;
}

function renderOverallMetrics(overall) {
  const items = [
    ['Total Requests', formatNumber(overall.httpReqs, 0)],
    ['Total Req / Detik', formatNumber(overall.httpReqRate, 2)],
    ['Overall P95', formatMs(overall.durationP95)],
    ['Overall Error Rate', formatPercent(overall.failedRate)],
    ['Check Rate', formatPercent(overall.checksRate)],
    ['Max VU', formatNumber(overall.vusMax, 0)]
  ];

  return `
    <section class="metric-group muted-group">
      <h3>Gabungan</h3>
      <div class="metric-grid">
        ${items.map(([label, value]) => metricCard(label, value)).join('')}
      </div>
    </section>
  `;
}

function metricCard(label, value) {
  return `
    <div class="metric">
      <span>${label}</span>
      <strong>${value}</strong>
    </div>
  `;
}

async function startRun(scenarioId) {
  try {
    const { run } = await api('/api/runs', {
      method: 'POST',
      body: JSON.stringify({ scenarioId })
    });
    selectedRunId = run.id;
    await refresh();
    await selectRun(run.id);
  } catch (error) {
    alert(error.message);
  }
}

async function selectRun(id) {
  selectedRunId = id;
  const { run } = await api(`/api/runs/${encodeURIComponent(id)}`);
  renderMetrics(run);
  logOutput.textContent = run.log || 'Log belum tersedia.';
}

async function refresh() {
  const [scenarioData, runData] = await Promise.all([
    api('/api/scenarios'),
    api('/api/runs')
  ]);

  scenarios = scenarioData.scenarios;
  activeStatus.textContent = runData.activeRun ? 'Sedang berjalan' : 'Siap';
  activeRun.textContent = runData.activeRun ? runData.activeRun.scenarioName : '-';
  renderScenarios(Boolean(runData.activeRun));
  renderRuns(runData.runs);

  if (!selectedRunId && runData.runs[0]) {
    selectedRunId = runData.runs[0].id;
  }

  if (selectedRunId) {
    const exists = runData.runs.some((run) => run.id === selectedRunId);
    if (exists) await selectRun(selectedRunId);
  }
}

refreshBtn.addEventListener('click', refresh);
refresh().catch((error) => {
  activeStatus.textContent = 'Error';
  logOutput.textContent = error.message;
});

setInterval(() => {
  refresh().catch(() => {});
}, 5000);
