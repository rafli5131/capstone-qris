const http = require('http');
const fs = require('fs');
const path = require('path');
const { spawn } = require('child_process');

const PORT = Number(process.env.PORT || 8080);
const RESULTS_DIR = process.env.RESULTS_DIR || '/results';
const SCRIPTS_DIR = process.env.SCRIPTS_DIR || '/scripts';

const scenarios = [
  {
    id: 'normal',
    file: 'normal_test.js',
    name: 'Normal Network Test',
    users: '100 VU per sistem',
    duration: '180 detik',
    threshold: 'p95 < 500ms, error rate < 1%',
    description: 'Kondisi jaringan urban standar dengan delay 0.5 detik.'
  },
  {
    id: 'peak_rural',
    file: 'peak_rural_test.js',
    name: 'Peak Rural Network Test',
    users: '500 VU per sistem',
    duration: '180 detik',
    threshold: 'p95 < 5s, error rate < 10%',
    description: 'Peak load dengan latensi rural 0.5 sampai 1.5 detik.'
  },
  {
    id: 'peak',
    file: 'peak_test.js',
    name: 'Peak Load Test',
    users: '1000 VU per sistem',
    duration: '180 detik',
    threshold: 'p95 < 1s, error rate < 5%',
    description: 'Simulasi traffic tinggi jam sibuk dengan delay pendek.'
  },
  {
    id: 'rural',
    file: 'rural_test.js',
    name: 'Rural Network Test',
    users: '10 VU per sistem',
    duration: '180 detik',
    threshold: 'p95 < 3s, error rate < 5%',
    description: 'Bandwidth rendah dan latensi tinggi dengan delay acak 0.5 sampai 2 detik.'
  }
];

let runs = [];
let activeRun = null;

fs.mkdirSync(RESULTS_DIR, { recursive: true });
loadRuns();

function loadRuns() {
  const files = fs.readdirSync(RESULTS_DIR, { withFileTypes: true })
    .filter((entry) => entry.isFile() && entry.name.endsWith('-run.json'))
    .map((entry) => path.join(RESULTS_DIR, entry.name));

  runs = files.map((file) => {
    try {
      return JSON.parse(fs.readFileSync(file, 'utf8'));
    } catch (_) {
      return null;
    }
  }).filter(Boolean).sort((a, b) => new Date(b.startedAt) - new Date(a.startedAt));
}

function saveRun(run) {
  const file = path.join(RESULTS_DIR, `${run.id}-run.json`);
  fs.writeFileSync(file, JSON.stringify(run, null, 2));
}

function readSummary(run) {
  if (!run.summaryPath || !fs.existsSync(run.summaryPath)) {
    return null;
  }

  try {
    return JSON.parse(fs.readFileSync(run.summaryPath, 'utf8'));
  } catch (_) {
    return null;
  }
}

function metricValue(summary, metric, value) {
  const item = summary?.metrics?.[metric];
  if (!item) {
    return null;
  }

  return item.values?.[value] ?? item[value] ?? null;
}

function summarize(run) {
  const summary = readSummary(run);
  if (!summary) {
    return null;
  }

  const overall = {
    httpReqs: metricValue(summary, 'http_reqs', 'count'),
    httpReqRate: metricValue(summary, 'http_reqs', 'rate'),
    checksRate: metricValue(summary, 'checks', 'value'),
    failedRate: metricValue(summary, 'http_req_failed', 'value'),
    durationAvg: metricValue(summary, 'http_req_duration', 'avg'),
    durationP95: metricValue(summary, 'http_req_duration', 'p(95)'),
    durationMax: metricValue(summary, 'http_req_duration', 'max'),
    vusMax: metricValue(summary, 'vus_max', 'max'),
    iterations: metricValue(summary, 'iterations', 'count')
  };

  const go = serviceSummary(summary, 'go', 'Sistem Baru Go');
  const legacy = serviceSummary(summary, 'legacy', 'Legacy Java');

  return {
    ...overall,
    overall,
    services: { go, legacy }
  };
}

function serviceSummary(summary, prefix, name) {
  const fallback = legacyServiceFallback(summary, prefix);
  const httpReqs = metricValue(summary, `${prefix}_http_reqs`, 'count') ?? fallback.httpReqs;
  const durationAvg = metricValue(summary, `${prefix}_http_req_duration`, 'avg');
  const durationP95 = metricValue(summary, `${prefix}_http_req_duration`, 'p(95)');
  const durationMax = metricValue(summary, `${prefix}_http_req_duration`, 'max');

  return {
    name,
    httpReqs,
    httpReqRate: metricValue(summary, `${prefix}_http_reqs`, 'rate'),
    failedRate: metricValue(summary, `${prefix}_http_req_failed`, 'value'),
    durationAvg,
    durationP95,
    durationMax,
    hasData: httpReqs !== null || durationAvg !== null || durationP95 !== null || durationMax !== null,
    partialData: fallback.partialData && durationAvg === null && durationP95 === null && durationMax === null
  };
}

function legacyServiceFallback(summary, prefix) {
  const serviceKey = prefix === 'go' ? 'go-sistem-baru' : 'legacy-system-java';
  const groups = summary?.root_group?.groups || {};
  const group = Object.values(groups).find((item) => item.name?.startsWith(serviceKey));
  const checks = group?.checks || {};
  const responded = checks[`${serviceKey} responded`];
  if (!responded) {
    return { httpReqs: null, partialData: false };
  }

  return {
    httpReqs: Number(responded.passes || 0) + Number(responded.fails || 0),
    partialData: true
  };
}

function publicRun(run, includeLog = false) {
  const log = includeLog && run.logPath && fs.existsSync(run.logPath)
    ? fs.readFileSync(run.logPath, 'utf8')
    : undefined;

  return {
    ...run,
    log,
    summary: summarize(run)
  };
}

function sendJson(res, status, body) {
  const payload = JSON.stringify(body);
  res.writeHead(status, {
    'Content-Type': 'application/json; charset=utf-8',
    'Content-Length': Buffer.byteLength(payload)
  });
  res.end(payload);
}

function sendStatic(req, res) {
  const urlPath = req.url === '/' ? '/index.html' : req.url;
  const cleanPath = path.normalize(decodeURIComponent(urlPath)).replace(/^(\.\.[/\\])+/, '');
  const filePath = path.join(__dirname, 'public', cleanPath);

  if (!filePath.startsWith(path.join(__dirname, 'public'))) {
    res.writeHead(403);
    res.end('Forbidden');
    return;
  }

  fs.readFile(filePath, (err, content) => {
    if (err) {
      res.writeHead(404);
      res.end('Not found');
      return;
    }

    const ext = path.extname(filePath);
    const type = ext === '.css' ? 'text/css' : ext === '.js' ? 'text/javascript' : 'text/html';
    res.writeHead(200, { 'Content-Type': `${type}; charset=utf-8` });
    res.end(content);
  });
}

function readBody(req) {
  return new Promise((resolve, reject) => {
    let body = '';
    req.on('data', (chunk) => {
      body += chunk;
      if (body.length > 1024 * 1024) {
        reject(new Error('Body terlalu besar'));
        req.destroy();
      }
    });
    req.on('end', () => {
      try {
        resolve(body ? JSON.parse(body) : {});
      } catch (error) {
        reject(error);
      }
    });
  });
}

function startRun(scenarioId) {
  const scenario = scenarios.find((item) => item.id === scenarioId);
  if (!scenario) {
    const error = new Error('Scenario tidak dikenal');
    error.status = 404;
    throw error;
  }

  if (activeRun) {
    const error = new Error('Masih ada test yang berjalan');
    error.status = 409;
    throw error;
  }

  const now = new Date();
  const id = `${now.toISOString().replace(/[:.]/g, '-')}-${scenario.id}`;
  const summaryPath = path.join(RESULTS_DIR, `${id}-summary.json`);
  const logPath = path.join(RESULTS_DIR, `${id}.log`);
  const run = {
    id,
    scenarioId: scenario.id,
    scenarioName: scenario.name,
    file: scenario.file,
    status: 'running',
    startedAt: now.toISOString(),
    endedAt: null,
    exitCode: null,
    summaryPath,
    logPath
  };

  saveRun(run);
  runs.unshift(run);
  activeRun = run;

  const args = ['run', '--summary-export', summaryPath, path.join(SCRIPTS_DIR, scenario.file)];
  const child = spawn('k6', args, {
    env: {
      ...process.env,
      GO_BASE_URL: process.env.GO_BASE_URL || 'http://go-sistem-baru:3000',
      LEGACY_BASE_URL: process.env.LEGACY_BASE_URL || 'http://legacy-system-java:8081',
      K6_NO_USAGE_REPORT: 'true'
    }
  });

  const stream = fs.createWriteStream(logPath, { flags: 'a' });
  stream.write(`$ k6 ${args.join(' ')}\n\n`);

  child.stdout.on('data', (data) => stream.write(data));
  child.stderr.on('data', (data) => stream.write(data));
  child.on('error', (error) => {
    stream.write(`\nRunner error: ${error.message}\n`);
  });
  child.on('close', (code) => {
    run.status = code === 0 ? 'completed' : 'failed';
    run.exitCode = code;
    run.endedAt = new Date().toISOString();
    saveRun(run);
    stream.end();
    activeRun = null;
    loadRuns();
  });

  return publicRun(run, true);
}

const server = http.createServer(async (req, res) => {
  try {
    const url = new URL(req.url, `http://${req.headers.host}`);

    if (req.method === 'GET' && url.pathname === '/api/scenarios') {
      sendJson(res, 200, { scenarios });
      return;
    }

    if (req.method === 'GET' && url.pathname === '/api/runs') {
      sendJson(res, 200, { activeRun: activeRun ? publicRun(activeRun) : null, runs: runs.map((run) => publicRun(run)) });
      return;
    }

    if (req.method === 'GET' && url.pathname.startsWith('/api/runs/')) {
      const id = decodeURIComponent(url.pathname.split('/').pop());
      const run = runs.find((item) => item.id === id) || (activeRun?.id === id ? activeRun : null);
      if (!run) {
        sendJson(res, 404, { error: 'Run tidak ditemukan' });
        return;
      }
      sendJson(res, 200, { run: publicRun(run, true) });
      return;
    }

    if (req.method === 'POST' && url.pathname === '/api/runs') {
      const body = await readBody(req);
      const run = startRun(body.scenarioId);
      sendJson(res, 201, { run });
      return;
    }

    sendStatic(req, res);
  } catch (error) {
    sendJson(res, error.status || 500, { error: error.message || 'Internal server error' });
  }
});

server.listen(PORT, () => {
  console.log(`K6 web runner listening on ${PORT}`);
});
