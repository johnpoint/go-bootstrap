'use strict';

const fs = require('fs');
const https = require('https');

const MODELS_ENDPOINT = {
  hostname: 'models.inference.ai.azure.com',
  path: '/chat/completions',
};

const REVIEW_PROMPT =
  'You are a senior software engineer conducting a thorough code review. ' +
  'Analyze the provided PR diff and give actionable feedback. ' +
  'Structure your review with these sections:\n' +
  '## Summary\nBrief overview of the changes.\n\n' +
  '## Issues\nList any bugs, security vulnerabilities, or logic errors. If none, say "No issues found."\n\n' +
  '## Suggestions\nCode quality, performance, and maintainability improvements.\n\n' +
  '## Positives\nWell-written parts worth highlighting.\n\n' +
  'Be specific, cite line numbers when relevant, and use markdown.';

function buildPayload(model, diff) {
  return JSON.stringify({
    model,
    messages: [
      { role: 'system', content: REVIEW_PROMPT },
      { role: 'user', content: 'Please review the following PR diff:\n\n```diff\n' + diff + '\n```' },
    ],
    max_tokens: 2048,
    temperature: 0.3,
  });
}

function makeRequest(payload, token) {
  return new Promise((resolve, reject) => {
    const requestOptions = {
      hostname: MODELS_ENDPOINT.hostname,
      path: MODELS_ENDPOINT.path,
      method: 'POST',
      headers: {
        Authorization: 'Bearer ' + token,
        'Content-Type': 'application/json',
        'Content-Length': Buffer.byteLength(payload),
      },
    };

    const req = https.request(requestOptions, (res) => {
      let data = '';
      res.on('data', (chunk) => { data += chunk; });
      res.on('end', () => {
        if (res.statusCode === 429) {
          reject(new Error('rate_limited'));
          return;
        }
        if (res.statusCode >= 500) {
          console.error('Server error response body:', data);
          reject(new Error('server_error:' + res.statusCode));
          return;
        }
        if (res.statusCode !== 200) {
          console.error('Unexpected status ' + res.statusCode + ' response body:', data);
          reject(new Error('http_error:' + res.statusCode));
          return;
        }
        try {
          const response = JSON.parse(data);
          if (response.choices && response.choices[0]) {
            resolve(response.choices[0].message.content);
          } else {
            console.error('Unexpected response structure:', data);
            reject(new Error('unexpected_response'));
          }
        } catch (e) {
          console.error('Failed to parse response body:', data.substring(0, 500));
          reject(new Error('parse_error:' + e.message));
        }
      });
    });

    req.on('error', (e) => {
      console.error('Network error:', e.message);
      reject(e);
    });

    req.write(payload);
    req.end();
  });
}

async function generateReview(diff, model, token, maxRetries) {
  const payload = buildPayload(model, diff);

  for (let attempt = 1; attempt <= maxRetries; attempt++) {
    try {
      console.log('Attempt ' + attempt + '/' + maxRetries + ' using model: ' + model);
      return await makeRequest(payload, token);
    } catch (e) {
      console.error('Attempt ' + attempt + ' failed:', e.message);
      if (attempt < maxRetries) {
        const delayMs = Math.pow(2, attempt) * 1000;
        console.log('Retrying in ' + delayMs + 'ms...');
        await new Promise((r) => setTimeout(r, delayMs));
      }
    }
  }

  return '_Code review failed after ' + maxRetries + ' attempts. Please check the Actions logs for details._';
}

async function main() {
  const diffPath = process.env.DIFF_PATH || '/tmp/pr.diff';
  const reviewPath = process.env.REVIEW_PATH || '/tmp/review.txt';
  const modelPath = process.env.MODEL_PATH || '/tmp/model.txt';
  const model = process.env.REVIEW_MODEL || 'gpt-4o-mini';
  const token = process.env.GITHUB_TOKEN;
  const maxRetries = parseInt(process.env.MAX_RETRIES || '3', 10);
  const isTruncated = process.env.TRUNCATED === 'true';

  if (!token) {
    console.error('GITHUB_TOKEN environment variable is required');
    process.exit(1);
  }

  const diff = fs.readFileSync(diffPath, 'utf8');

  if (!diff.trim()) {
    fs.writeFileSync(reviewPath, 'No code changes found to review.');
    fs.writeFileSync(modelPath, model);
    return;
  }

  let review = await generateReview(diff, model, token, maxRetries);

  if (isTruncated) {
    review += '\n\n> ⚠️ **Note**: The diff was large and truncated. This review covers only a portion of the changes.';
  }

  fs.writeFileSync(reviewPath, review);
  fs.writeFileSync(modelPath, model);
  console.log('Review written to ' + reviewPath);
}

main().catch((e) => {
  console.error('Fatal error:', e.message);
  process.exit(1);
});
