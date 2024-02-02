#!/bin/sh

cd /app
npm install --omit=dev
node index.js
